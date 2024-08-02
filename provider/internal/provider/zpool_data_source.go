package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var (
	_ datasource.DataSource              = &zpoolDataSource{}
	_ datasource.DataSourceWithConfigure = &zpoolDataSource{}
)

type zpoolDataSource struct {
	client *common.Client
}

type zpoolDataSourceModel struct {
	Zpools []zpoolDataModel `tfsdk:"zpools"`
}

type zpoolDataModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewZpoolDataSource() datasource.DataSource {
	return &zpoolDataSource{}
}

func (d *zpoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.Client)
	if !ok {
		UnexpectedDataSourceConfigureType(ctx, req, resp)
	}
	d.client = client
}

func (d *zpoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zpools"
}

func (d *zpoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get configuration for a ZPool",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the zpool",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the zpool"},
		},
	}
}

func (d *zpoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zpoolDataSourceModel

	zpools, err := d.client.ZfsGetPools(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get zpools", err.Error())
		return
	}

	for _, pool := range zpools.Pools {
		zpoolState := zpoolDataModel{
			ID:   types.StringValue(pool.Name),
			Name: types.StringValue(pool.Name),
		}

		state.Zpools = append(state.Zpools, zpoolState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
