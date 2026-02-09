package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var (
	_ resource.Resource                = &FirewallRuleResource{}
	_ resource.ResourceWithImportState = &FirewallRuleResource{}
)

type FirewallRuleResource struct {
	client *common.Client
}

type FirewallRuleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Zone     types.String `tfsdk:"zone"`
	RuleType types.String `tfsdk:"rule_type"`
	Rule     types.String `tfsdk:"rule"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Service  types.String `tfsdk:"service"`
}

func NewFirewallRuleResource() resource.Resource {
	return &FirewallRuleResource{}
}

func (r *FirewallRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.Client)
	if !ok {
		ProviderDataError(req.ProviderData, &resp.Diagnostics)
	}

	r.client = client
}

func (r *FirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *FirewallRuleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Firewall Rule (rich rule, port, or service)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the rule",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "Zone to add the rule to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of rule: 'rich', 'port', or 'service'",
				Validators: []validator.String{
					stringvalidator.OneOf("rich", "port", "service"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule": schema.StringAttribute{
				Optional:    true,
				Description: "Rich rule string (required if rule_type is 'rich')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "Port number (required if rule_type is 'port')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Protocol (tcp/udp, required if rule_type is 'port')",
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.StringAttribute{
				Optional:    true,
				Description: "Service name (required if rule_type is 'service')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallRuleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone := plan.Zone.ValueString()
	ruleType := plan.RuleType.ValueString()

	request := common.FirewallRuleRequest{
		Zone:     zone,
		RuleType: ruleType,
		Rule:     plan.Rule.ValueString(),
		Port:     plan.Port.ValueString(),
		Protocol: plan.Protocol.ValueString(),
		Service:  plan.Service.ValueString(),
	}

	// Validate required fields based on rule type
	if ruleType == "rich" && request.Rule == "" {
		resp.Diagnostics.AddError("Missing required field", "rule is required when rule_type is 'rich'")
		return
	}
	if ruleType == "port" && (request.Port == "" || request.Protocol == "") {
		resp.Diagnostics.AddError("Missing required fields", "port and protocol are required when rule_type is 'port'")
		return
	}
	if ruleType == "service" && request.Service == "" {
		resp.Diagnostics.AddError("Missing required field", "service is required when rule_type is 'service'")
		return
	}

	tflog.Debug(ctx, "Attempting to add firewall rule", map[string]any{"zone": zone, "rule_type": ruleType})

	ruleResp, err := r.client.FirewallAddRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add firewall rule", fmt.Sprintf("Failed to add rule. Unexpected error: %s", err.Error()))
		return
	}

	// Generate ID based on rule type and details
	var id string
	switch ruleType {
	case "rich":
		id = fmt.Sprintf("%s:rich:%s", zone, ruleResp.Rule)
	case "port":
		id = fmt.Sprintf("%s:port:%s/%s", zone, ruleResp.Port, ruleResp.Protocol)
	case "service":
		id = fmt.Sprintf("%s:service:%s", zone, ruleResp.Service)
	}

	plan = FirewallRuleResourceModel{
		ID:       types.StringValue(id),
		Zone:     types.StringValue(zone),
		RuleType: types.StringValue(ruleType),
		Rule:     types.StringValue(ruleResp.Rule),
		Port:     types.StringValue(ruleResp.Port),
		Protocol: types.StringValue(ruleResp.Protocol),
		Service:  types.StringValue(ruleResp.Service),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallRuleResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Firewall rules don't have a direct "get" method in firewalld
	// We assume if the resource exists in state, it exists in firewalld
	// In a production scenario, you'd query the zone to verify the rule exists
	tflog.Debug(ctx, "Reading firewall rule", map[string]any{"id": state.ID.ValueString()})

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Rules are immutable, so update requires replace
	resp.Diagnostics.AddWarning("Update not supported", "Firewall rules cannot be updated in place. Please recreate the resource.")
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallRuleResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := common.FirewallRuleRequest{
		Zone:     state.Zone.ValueString(),
		RuleType: state.RuleType.ValueString(),
		Rule:     state.Rule.ValueString(),
		Port:     state.Port.ValueString(),
		Protocol: state.Protocol.ValueString(),
		Service:  state.Service.ValueString(),
	}

	tflog.Debug(ctx, "Deleting firewall rule", map[string]any{"id": state.ID.ValueString()})

	err := r.client.FirewallRemoveRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete firewall rule", fmt.Sprintf("Unable to delete rule. Unexpected error: %s", err))
		return
	}
}

func (r *FirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
