package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func UnexpectedDataSourceConfigureType(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	resp.Diagnostics.AddError(
		"Unexpected Data Source Configure Type",
		fmt.Sprintf(
			"Expected *Client, got: %T. Please report this issue to the provider developers.",
			req.ProviderData,
		),
	)
}

func ProviderDataError(data any, diags *diag.Diagnostics) {
	diags.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *common.Client, got: %T. Please report this issue to the provider developers.", data))
}
