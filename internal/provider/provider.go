package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theoremoon/terraform-provider-pubsubhubbub-subscribe/internal/client"
)

type subscribeProvider struct {
	version string
}

type subscribeProviderModel struct {
	HubUrl types.String `tfsdk:"hub_url"`
}

type subscribeProviderData struct {
	Client *client.Client
}

var (
	_ provider.Provider = &subscribeProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &subscribeProvider{
			version: version,
		}
	}
}

func (p *subscribeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pubsubhubbub-subscribe"
	resp.Version = p.version
}

func (p *subscribeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A terraform provider for subscribing/unsubscribing topic URLs to a PubSubHubbub hub.",
		Attributes: map[string]schema.Attribute{
			"hub_url": schema.StringAttribute{
				Description: "The URL of the PubSubHubbub hub to subscribe to. Defaults to 'https://pubsubhubbub.appspot.com/'.",
				Optional:    true,
			},
		},
	}
}

func (p *subscribeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config subscribeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var hubURL string
	if config.HubUrl.IsNull() {
		hubURL = "https://pubsubhubbub.appspot.com/"
	} else {
		hubURL = config.HubUrl.ValueString()
	}

	client := client.NewClient(p.version, hubURL)
	data := subscribeProviderData{
		Client: client,
	}
	resp.DataSourceData = &data
	resp.ResourceData = &data
}

func (p *subscribeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSubscribeResource,
	}
}

func (p *subscribeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
