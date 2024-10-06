package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theoremoon/terraform-provider-pubsubhubbub-subscribe/internal/client"
)

type SubscribeResource struct {
	client *client.Client
}

type subscribeResourceModel struct {
	CallbackURL types.String `tfsdk:"username"`
	TopicURL    types.String `tfsdk:"owner"`
	HMACSecret  types.String `tfsdk:"apikey"`
}

// ensure
var (
	_ resource.Resource = &SubscribeResource{}
)

func NewSubscribeResource() resource.Resource {
	return &SubscribeResource{}
}

func (r *SubscribeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscribe"
}

func (r *SubscribeResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"callback_url": schema.StringAttribute{
				Description: "The URL of the callback server that will receive the feed updates.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// Validators
			},
			"topic_url": schema.StringAttribute{
				Description: "The URL of the topic feed to subscribe to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hmac_secret": schema.StringAttribute{
				Description: "The secret key used to sign the hub's requests to the callback server.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SubscribeResource) Configure(req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*subscribeProviderData).Client
}

func (r *SubscribeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subscribeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	callbackURL := plan.CallbackURL.ValueString()
	topicURL := plan.TopicURL.ValueString()
	hmacSecret := plan.HMACSecret.ValueString() // maybe empty

	err := r.client.Subscribe(callbackURL, topicURL, hmacSecret)
	if err != nil {
		resp.Diagnostics.AddError("Request Error", fmt.Sprintf("Failed to subscribe topic %s to %s: %s", topicURL, callbackURL, err))
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read actually does nothing. Because the hub does not provide a way to check the subscription status.
func (r *SubscribeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subscribeResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	return
}

// Update is not supported for this resource. Just replace the resource with a new one.
func (r *SubscribeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *SubscribeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subscribeResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	callbackURL := state.CallbackURL.ValueString()
	topicURL := state.TopicURL.ValueString()
	hmacSecret := state.HMACSecret.ValueString() // maybe empty

	err := r.client.Unsubscribe(callbackURL, topicURL, hmacSecret)
	if err != nil {
		resp.Diagnostics.AddError("Request Error", fmt.Sprintf("Failed to unsubscribe topic %s to %s: %s", topicURL, callbackURL, err))
		return
	}

	resp.State.RemoveResource(ctx)
}
