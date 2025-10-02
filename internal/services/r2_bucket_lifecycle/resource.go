// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package r2_bucket_lifecycle

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/r2"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/apijson"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*R2BucketLifecycleResource)(nil)
var _ resource.ResourceWithModifyPlan = (*R2BucketLifecycleResource)(nil)

func NewResource() resource.Resource {
	return &R2BucketLifecycleResource{}
}

// R2BucketLifecycleResource defines the resource implementation.
type R2BucketLifecycleResource struct {
	client *cloudflare.Client
}

func (r *R2BucketLifecycleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_r2_bucket_lifecycle"
}

func (r *R2BucketLifecycleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cloudflare.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("Expected *cloudflare.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *R2BucketLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *R2BucketLifecycleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sortedRules := sortRulesByID(data.Rules)
	data.Rules = &sortedRules

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	env := R2BucketLifecycleResultEnvelope{*data}
	_, err = r.client.R2.Buckets.Lifecycle.Update(
		ctx,
		data.BucketName.ValueString(),
		r2.BucketLifecycleUpdateParams{
			AccountID: cloudflare.F(data.AccountID.ValueString()),
		},
		option.WithHeader(consts.R2JurisdictionHTTPHeaderName, data.Jurisdiction.ValueString()),
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &env)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data = &env.Result
	// sortedRules := sortRulesByID(data.Rules)
	// data.Rules = &sortedRules

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *R2BucketLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *R2BucketLifecycleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sortedRules := sortRulesByID(data.Rules)
	data.Rules = &sortedRules

	var state *R2BucketLifecycleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	env := R2BucketLifecycleResultEnvelope{*data}
	_, err = r.client.R2.Buckets.Lifecycle.Update(
		ctx,
		data.BucketName.ValueString(),
		r2.BucketLifecycleUpdateParams{
			AccountID: cloudflare.F(data.AccountID.ValueString()),
		},
		option.WithHeader(consts.R2JurisdictionHTTPHeaderName, data.Jurisdiction.ValueString()),
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &env)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data = &env.Result

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *R2BucketLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *R2BucketLifecycleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sortedRules := sortRulesByID(data.Rules)
	data.Rules = &sortedRules

	res := new(http.Response)
	env := R2BucketLifecycleResultEnvelope{*data}
	_, err := r.client.R2.Buckets.Lifecycle.Get(
		ctx,
		data.BucketName.ValueString(),
		r2.BucketLifecycleGetParams{
			AccountID: cloudflare.F(data.AccountID.ValueString()),
		},
		option.WithHeader(consts.R2JurisdictionHTTPHeaderName, data.Jurisdiction.ValueString()),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if res != nil && res.StatusCode == 404 {
		resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &env)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data = &env.Result

	// sortedRules := sortRulesByID(data.Rules)
	// data.Rules = &sortedRules

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *R2BucketLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *R2BucketLifecycleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"This resource cannot be destroyed from Terraform. If you create this resource, it will be "+
				"present in the API until manually deleted.",
		)
	}
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will remove the resource from the Terraform state "+
				"but will not change it in the API. If you would like to destroy or reset this resource "+
				"in the API, refer to the documentation for how to do it manually.",
		)
	}

	var planApp, stateApp *R2BucketLifecycleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planApp)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateApp)...)

	if resp.Diagnostics.HasError() || planApp == nil {
		return
	}

	// Check if Rules is set in both plan and state
	if stateApp != nil && planApp.Rules != nil && len(*planApp.Rules) > 0 &&
		stateApp.Rules != nil && len(*stateApp.Rules) > 0 {

		planAppRules := sortRulesByID(planApp.Rules)
		stateAppRules := sortRulesByID(stateApp.Rules)

		// If lists are equal (ignoring order), use the state's version to prevent spurious diffs
		if rulesAreEqual(planAppRules, stateAppRules) {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), stateApp.Rules)...)
		} else {
			// If they're different, set the sorted plan version
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), &planAppRules)...)
		}
	} else if planApp.Rules != nil && len(*planApp.Rules) > 0 {
		// If only plan has rules, sort them for consistency
		sortedRules := sortRulesByID(planApp.Rules)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), &sortedRules)...)
	}
}

// sortRulesByID returns a sorted copy of the rules slice, sorted by ID
func sortRulesByID(rules *[]*R2BucketLifecycleRulesModel) []*R2BucketLifecycleRulesModel {
	if rules == nil || len(*rules) == 0 {
		return []*R2BucketLifecycleRulesModel{}
	}

	// Create a copy to avoid modifying the original
	sortedRules := make([]*R2BucketLifecycleRulesModel, len(*rules))
	copy(sortedRules, *rules)

	// Sort by ID
	sort.SliceStable(sortedRules, func(i, j int) bool {
		// Handle nil cases
		if sortedRules[i] == nil || sortedRules[i].ID.IsNull() {
			return true
		}
		if sortedRules[j] == nil || sortedRules[j].ID.IsNull() {
			return false
		}

		return sortedRules[i].ID.ValueString() < sortedRules[j].ID.ValueString()
	})

	return sortedRules
}

// rulesAreEqual checks if two sorted rule lists contain the same rules (ignoring order)
func rulesAreEqual(rules1, rules2 []*R2BucketLifecycleRulesModel) bool {
	if len(rules1) != len(rules2) {
		return false
	}

	// Since both are sorted by ID, we can compare element by element
	for i := range rules1 {
		if !ruleEquals(rules1[i], rules2[i]) {
			return false
		}
	}

	return true
}

// ruleEquals compares two individual rules for semantic equality
func ruleEquals(rule1, rule2 *R2BucketLifecycleRulesModel) bool {
	if rule1 == nil && rule2 == nil {
		return true
	}
	if rule1 == nil || rule2 == nil {
		return false
	}

	// Compare ID
	if !rule1.ID.Equal(rule2.ID) {
		return false
	}

	// Compare Enabled
	if !rule1.Enabled.Equal(rule2.Enabled) {
		return false
	}

	// Compare Conditions
	if !conditionsEqual(rule1.Conditions, rule2.Conditions) {
		return false
	}

	// Compare AbortMultipartUploadsTransition
	if !abortMultipartTransitionEqual(rule1.AbortMultipartUploadsTransition, rule2.AbortMultipartUploadsTransition) {
		return false
	}

	// Compare DeleteObjectsTransition
	if !deleteObjectsTransitionEqual(rule1.DeleteObjectsTransition, rule2.DeleteObjectsTransition) {
		return false
	}

	// Compare StorageClassTransitions
	if !storageClassTransitionsEqual(rule1.StorageClassTransitions, rule2.StorageClassTransitions) {
		return false
	}

	return true
}

// conditionsEqual compares two Conditions objects
func conditionsEqual(c1, c2 *R2BucketLifecycleRulesConditionsModel) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}

	return c1.Prefix.Equal(c2.Prefix)
}

// abortMultipartTransitionEqual compares two AbortMultipartUploadsTransition objects
func abortMultipartTransitionEqual(t1, t2 *R2BucketLifecycleRulesAbortMultipartUploadsTransitionModel) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}

	return transitionConditionEqual(t1.Condition, t2.Condition)
}

// deleteObjectsTransitionEqual compares two DeleteObjectsTransition objects
func deleteObjectsTransitionEqual(t1, t2 *R2BucketLifecycleRulesDeleteObjectsTransitionModel) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}

	return transitionConditionEqual(t1.Condition, t2.Condition)
}

// transitionConditionEqual compares transition conditions (works for both abort and delete)
func transitionConditionEqual(c1, c2 interface{}) bool {
	// This is a generic comparison - adjust based on your actual condition model structure
	// You might need to type assert to the specific condition model type

	// If both are nil, they're equal
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}

	// Type assert to your specific condition models and compare fields
	// For example, if you have a shared condition model:
	switch cond1 := c1.(type) {
	case *R2BucketLifecycleRulesAbortMultipartUploadsTransitionConditionModel:
		cond2, ok := c2.(*R2BucketLifecycleRulesAbortMultipartUploadsTransitionConditionModel)
		if !ok {
			return false
		}
		return cond1.MaxAge.Equal(cond2.MaxAge) && cond1.Type.Equal(cond2.Type)

	case *R2BucketLifecycleRulesDeleteObjectsTransitionConditionModel:
		cond2, ok := c2.(*R2BucketLifecycleRulesDeleteObjectsTransitionConditionModel)
		if !ok {
			return false
		}
		return cond1.MaxAge.Equal(cond2.MaxAge) &&
			cond1.Type.Equal(cond2.Type) &&
			cond1.Date.Equal(cond2.Date)
	}

	return false
}

// storageClassTransitionsEqual compares two StorageClassTransitions lists
func storageClassTransitionsEqual(t1, t2 *[]*R2BucketLifecycleRulesStorageClassTransitionsModel) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	if len(*t1) != len(*t2) {
		return false
	}

	// Sort both lists by storage class for comparison
	sorted1 := sortStorageClassTransitions(t1)
	sorted2 := sortStorageClassTransitions(t2)

	for i := range sorted1 {
		if !storageClassTransitionEqual(sorted1[i], sorted2[i]) {
			return false
		}
	}

	return true
}

// sortStorageClassTransitions sorts storage class transitions by storage class name
func sortStorageClassTransitions(transitions *[]*R2BucketLifecycleRulesStorageClassTransitionsModel) []*R2BucketLifecycleRulesStorageClassTransitionsModel {
	if transitions == nil || len(*transitions) == 0 {
		return []*R2BucketLifecycleRulesStorageClassTransitionsModel{}
	}

	sorted := make([]*R2BucketLifecycleRulesStorageClassTransitionsModel, len(*transitions))
	copy(sorted, *transitions)

	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i] == nil || sorted[i].StorageClass.IsNull() {
			return true
		}
		if sorted[j] == nil || sorted[j].StorageClass.IsNull() {
			return false
		}
		return sorted[i].StorageClass.ValueString() < sorted[j].StorageClass.ValueString()
	})

	return sorted
}

// storageClassTransitionEqual compares two individual storage class transitions
func storageClassTransitionEqual(t1, t2 *R2BucketLifecycleRulesStorageClassTransitionsModel) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}

	if !t1.StorageClass.Equal(t2.StorageClass) {
		return false
	}

	// Compare conditions
	if t1.Condition == nil && t2.Condition == nil {
		return true
	}
	if t1.Condition == nil || t2.Condition == nil {
		return false
	}

	return t1.Condition.MaxAge.Equal(t2.Condition.MaxAge) &&
		t1.Condition.Type.Equal(t2.Condition.Type) &&
		t1.Condition.Date.Equal(t2.Condition.Date)
}
