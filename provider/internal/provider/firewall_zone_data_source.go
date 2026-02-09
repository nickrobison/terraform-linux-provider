package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var (
	_ datasource.DataSource              = &firewallZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallZoneDataSource{}
)

type firewallZoneDataSource struct {
	client *common.Client
}

type firewallZoneDataSourceModel struct {
	Zones []firewallZoneDataModel `tfsdk:"zones"`
}

type firewallZoneDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Target      types.String `tfsdk:"target"`
	Services    types.List   `tfsdk:"services"`
	Ports       types.List   `tfsdk:"ports"`
	RichRules   types.List   `tfsdk:"rich_rules"`
}

func NewFirewallZoneDataSource() datasource.DataSource {
	return &firewallZoneDataSource{}
}

func (d *firewallZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.Client)
	if !ok {
		UnexpectedDataSourceConfigureType(ctx, req, resp)
	}
	d.client = client
}

func (d *firewallZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_zones"
}

func (d *firewallZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get configuration for firewall zones",
		Attributes: map[string]schema.Attribute{
			"zones": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of firewall zones",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "An identifier for the zone",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the zone",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the zone",
						},
						"target": schema.StringAttribute{
							Computed:    true,
							Description: "Default target for the zone",
						},
						"services": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Services allowed in the zone",
						},
						"ports": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Ports allowed in the zone",
						},
						"rich_rules": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Rich rules for the zone",
						},
					},
				},
			},
		},
	}
}

func (d *firewallZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewallZoneDataSourceModel

	zones, err := d.client.FirewallGetZones(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get firewall zones", err.Error())
		return
	}

	for _, zone := range zones {
		services, diags := types.ListValueFrom(ctx, types.StringType, zone.Services)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ports, diags := types.ListValueFrom(ctx, types.StringType, zone.Ports)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		richRules, diags := types.ListValueFrom(ctx, types.StringType, zone.RichRules)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		zoneState := firewallZoneDataModel{
			ID:          types.StringValue(zone.Name),
			Name:        types.StringValue(zone.Name),
			Description: types.StringValue(zone.Description),
			Target:      types.StringValue(zone.Target),
			Services:    services,
			Ports:       ports,
			RichRules:   richRules,
		}

		state.Zones = append(state.Zones, zoneState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
