package r2_bucket_lifecycle

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type listSortModifier struct {
	sortByAttribute string
}

func SortListByAttribute(attributeName string) planmodifier.List {
	return listSortModifier{
		sortByAttribute: attributeName,
	}
}

func (m listSortModifier) Description(ctx context.Context) string {
	return "Sorts list elements by " + m.sortByAttribute + " to prevent spurious diffs"
}

func (m listSortModifier) MarkdownDescription(ctx context.Context) string {
	return "Sorts list elements by `" + m.sortByAttribute + "` to prevent spurious diffs"
}

func (m listSortModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If the plan value is null or unknown, nothing to sort
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Sort the plan value
	var planElements []basetypes.ObjectValue
	diags := req.PlanValue.ElementsAs(ctx, &planElements, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortedPlanElements := m.sortElements(ctx, planElements, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert to []attr.Value
	sortedPlanAttrValues := make([]attr.Value, len(sortedPlanElements))
	for i, elem := range sortedPlanElements {
		sortedPlanAttrValues[i] = elem
	}

	// Create sorted plan list
	sortedPlanList, diags := types.ListValue(req.PlanValue.ElementType(ctx), sortedPlanAttrValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = sortedPlanList

	// IMPORTANT: Also sort the config value if it exists
	// This ensures Terraform compares sorted plan against sorted config
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		var configElements []basetypes.ObjectValue
		diags := req.ConfigValue.ElementsAs(ctx, &configElements, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		sortedConfigElements := m.sortElements(ctx, configElements, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert to []attr.Value
		sortedConfigAttrValues := make([]attr.Value, len(sortedConfigElements))
		for i, elem := range sortedConfigElements {
			sortedConfigAttrValues[i] = elem
		}

		// Create sorted config list
		sortedConfigList, diags := types.ListValue(req.ConfigValue.ElementType(ctx), sortedConfigAttrValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		req.ConfigValue = sortedConfigList
	}

	// IMPORTANT: Also sort the state value if it exists
	// This ensures consistency across all values
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		var stateElements []basetypes.ObjectValue
		diags := req.StateValue.ElementsAs(ctx, &stateElements, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		sortedStateElements := m.sortElements(ctx, stateElements, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert to []attr.Value
		sortedStateAttrValues := make([]attr.Value, len(sortedStateElements))
		for i, elem := range sortedStateElements {
			sortedStateAttrValues[i] = elem
		}

		// Create sorted state list
		sortedStateList, diags := types.ListValue(req.StateValue.ElementType(ctx), sortedStateAttrValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		req.StateValue = sortedStateList
	}
}

func (m listSortModifier) sortElements(ctx context.Context, elements []basetypes.ObjectValue, diags *diag.Diagnostics) []basetypes.ObjectValue {
	if len(elements) == 0 {
		return elements
	}

	// Create a sortable slice
	sortableElements := make([]sortableElement, len(elements))

	for i, elem := range elements {
		// Extract the sort key from the object's attributes
		sortKey := elem.Attributes()[m.sortByAttribute]
		if sortKey == nil {
			diags.AddWarning(
				"Missing Sort Attribute",
				fmt.Sprintf("Element %d missing attribute: %s", i, m.sortByAttribute),
			)
			sortableElements[i] = sortableElement{value: elem, sortKey: ""}
			continue
		}

		// Convert to string for sorting
		sortKeyStr := m.extractSortKey(sortKey)
		sortableElements[i] = sortableElement{
			value:   elem,
			sortKey: sortKeyStr,
		}
	}

	// Sort (stable sort to maintain order for equal keys)
	sort.SliceStable(sortableElements, func(i, j int) bool {
		return sortableElements[i].sortKey < sortableElements[j].sortKey
	})

	// Extract sorted values
	result := make([]basetypes.ObjectValue, len(sortableElements))
	for i, se := range sortableElements {
		result[i] = se.value
	}

	return result
}

func (m listSortModifier) extractSortKey(value attr.Value) string {
	// Handle different types
	switch v := value.(type) {
	case basetypes.StringValue:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return v.ValueString()
	case basetypes.Int64Value:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		// Format with leading zeros for proper lexicographic sorting
		return fmt.Sprintf("%020d", v.ValueInt64())
	case basetypes.NumberValue:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return v.ValueBigFloat().String()
	default:
		return ""
	}
}

type sortableElement struct {
	value   basetypes.ObjectValue
	sortKey string
}
