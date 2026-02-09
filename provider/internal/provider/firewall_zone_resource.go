package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var (
	_ resource.Resource                = &FirewallZoneResource{}
	_ resource.ResourceWithImportState = &FirewallZoneResource{}
)

type FirewallZoneResource struct {
	client *common.Client
}

type FirewallZoneResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Target      types.String `tfsdk:"target"`
}

func NewFirewallZoneResource() resource.Resource {
	return &FirewallZoneResource{}
}

func (r *FirewallZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.Client)
	if !ok {
		ProviderDataError(req.ProviderData, &resp.Diagnostics)
	}

	r.client = client
}

func (r *FirewallZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_zone"
}

func (r *FirewallZoneResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Firewall Zone",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the zone",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Zone name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Zone description",
			},
			"target": schema.StringAttribute{
				Optional:    true,
				Description: "Default target for the zone (default, ACCEPT, REJECT, DROP)",
			},
		},
	}
}

func (r *FirewallZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallZoneResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	target := plan.Target.ValueString()

	request := common.FirewallZoneCreateRequest{
		Name:        name,
		Description: description,
		Target:      target,
	}
	tflog.Debug(ctx, "Attempting to create firewall zone", map[string]any{"name": name})

	zone, err := r.client.FirewallCreateZone(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create firewall zone", fmt.Sprintf("Failed to create zone. Unexpected error: %s", err.Error()))
		return
	}

	plan = FirewallZoneResourceModel{
		ID:          types.StringValue(zone.Name),
		Name:        types.StringValue(zone.Name),
		Description: types.StringValue(zone.Description),
		Target:      types.StringValue(zone.Target),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallZoneResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneName := state.ID.ValueString()
	tflog.Debug(ctx, "Fetching firewall zone", map[string]any{"id": zoneName})
	err := r.doRead(ctx, zoneName, &state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read firewall zone", fmt.Sprintf("Unable to read zone. Unexpected error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallZoneResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Firewalld zones are immutable, so we need to recreate
	// This should be handled by RequiresReplace plan modifier
	resp.Diagnostics.AddWarning("Update not supported", "Firewall zones cannot be updated in place. Please recreate the resource.")
}

func (r *FirewallZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallZoneResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneName := state.Name.ValueString()
	tflog.Debug(ctx, "Deleting firewall zone", map[string]any{"name": zoneName})

	err := r.client.FirewallDeleteZone(ctx, zoneName)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete firewall zone", fmt.Sprintf("Unable to delete zone. Unexpected error: %s", err))
		return
	}
}

func (r *FirewallZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FirewallZoneResource) doRead(ctx context.Context, id string, data *FirewallZoneResourceModel) error {
	zone, err := r.client.FirewallGetZone(ctx, id)
	if err != nil {
		return err
	}

	data.ID = types.StringValue(zone.Name)
	data.Name = types.StringValue(zone.Name)
	data.Description = types.StringValue(zone.Description)
	data.Target = types.StringValue(zone.Target)
	return nil
}
