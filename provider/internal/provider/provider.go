package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nickrobison/terraform-linux-provider/common"
)

var _ provider.Provider = &LinuxProvider{}

type LinuxProvider struct {
	version string
}

type LinuxProviderModel struct {
	Host types.String `tfsdk:"host"`
	Port types.Int32  `tfsdk:"port"`
}

const (
	providerName = "linux"
)

func (p *LinuxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerName
	resp.Version = p.version
}

func (p *LinuxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
				Description: "This is the hostname for the API connection." +
					" May also be provided via " + EnvHost + " environment variable.",
			},
			"port": schema.Int32Attribute{
				Optional: true,
				Description: "This is the tcp port for the API connection." +
					" May also be provided via " + EnvPort + " environment variable.",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
		},
	}
}

func (p *LinuxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Initializing Linux provider")
	var config LinuxProviderModel

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	unknownValueErrorMessage := "The provider cannot create the Linux client as there is an unknown configuration value "
	instructionUnknownMessage := " Either target apply the source of the value first, " +
		"set the value statically in the configuration, or use the %s environment variable."

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("address"), "Unknown Linux Host", fmt.Sprintf("%s for the Linux Host. %s", unknownValueErrorMessage, fmt.Sprintf(instructionUnknownMessage, EnvHost)))
	}

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("port"), "Unknown Linux Port", fmt.Sprintf("%s for the Linux Port. %s", unknownValueErrorMessage, fmt.Sprintf(instructionUnknownMessage, EnvPort)))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv(EnvHost)
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("address"),
			"Missing Linux address target",
			"The provider cannot create the Linux client as there is a missing or empty value for the Linux address."+
				" Set the value in the configuration or use the "+EnvHost+" environment variable."+
				" If either is already set, ensure the value is not empty.",
		)

		return
	}

	client := common.NewClient(host)

	if !config.Port.IsNull() {
		client.WithPort(int(config.Port.ValueInt32()))
	} else if v := os.Getenv(EnvPort); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("port"),
				"Failed to parse "+EnvPort,
				fmt.Sprintf("Error to parse value in "+EnvPort+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithPort(d)
		}
	}

	ctx = tflog.SetField(ctx, "host", host)
	tflog.Info(ctx, "Created client")
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LinuxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewZpoolResource,
		NewFirewallZoneResource,
		NewFirewallRuleResource,
	}
}

func (p *LinuxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewZpoolDataSource,
		NewFirewallZoneDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LinuxProvider{
			version: version,
		}
	}
}
