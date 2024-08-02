package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var (
	_ resource.Resource                = &ZpoolResource{}
	_ resource.ResourceWithImportState = &ZpoolResource{}
)

type ZpoolResource struct {
	client *common.Client
}

type ZpoolResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewZpoolResource() resource.Resource {
	return &ZpoolResource{}
}

func (r *ZpoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.Client)
	if !ok {
		ProviderDataError(req.ProviderData, &resp.Diagnostics)
	}

	r.client = client
}

func (r *ZpoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zpool"
}

func (r *ZpoolResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Zpool",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the zpool",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Zpool name",
			},
		},
	}
}

func (r *ZpoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ZpoolResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	request := common.ZpoolCreateRequest{
		Name: name,
	}
	tflog.Debug(ctx, "Attempting to create zpool", map[string]any{"name": name})

	pool, err := r.client.ZfsCreatePool(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create zpool", fmt.Sprintf("Failed to create zpool. Unexpected error: %s", err.Error()))
		return
	}
	plan = ZpoolResourceModel{
		ID:   types.StringValue(pool.Name),
		Name: types.StringValue(pool.Name),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ZpoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ZpoolResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zpoolName := state.ID.ValueString()
	tflog.Debug(ctx, "Fetching zpool", map[string]any{"id": zpoolName})
	err := r.doRead(ctx, zpoolName, &state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read zpool", fmt.Sprintf("Unable to read zpool. Unexpected error: %a", err))
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ZpoolResource) Update(_ context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: Implement
}

func (r *ZpoolResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO: Implement
}

func (r *ZpoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ZpoolResource) doRead(ctx context.Context, id string, data *ZpoolResourceModel) error {
	zpool, err := r.client.ZfsGetPool(ctx, id)
	if err != nil {
		return err
	}

	data.ID = types.StringValue(zpool.Name)
	data.Name = types.StringValue(zpool.Name)
	return nil
}
