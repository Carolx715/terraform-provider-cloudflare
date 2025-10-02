package r2_bucket_lifecycle

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type listSortModifier struct {
	// The attribute name to sort by
	sortByAttribute string
}

// SortListByAttribute creates a plan modifier that sorts list elements by a specific attribute
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

	// Get the elements from the plan
	var planElements []attr.Value
	diags := req.PlanValue.ElementsAs(ctx, &planElements, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// // Sort the elements
	// sortedElements := m.sortElements(ctx, planElements, &resp.Diagnostics)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Create a new list with sorted elements
	// sortedList, diags := types.ListValue(req.PlanValue.ElementType(ctx), sortedElements)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// resp.PlanValue = req.PlanValue
}

func (m listSortModifier) sortElements(ctx context.Context, elements []attr.Value, diags *diag.Diagnostics) []attr.Value {
	// Create a sortable slice
	sortableElements := make([]sortableElement, len(elements))

	for i, elem := range elements {
		obj, ok := elem.(types.Object)
		if !ok {
			diags.AddError(
				"Invalid Element Type",
				"Expected object type in list",
			)
			return elements
		}

		// Extract the sort key
		sortKey := obj.Attributes()[m.sortByAttribute]
		if sortKey == nil {
			diags.AddWarning(
				"Missing Sort Attribute",
				"Element missing attribute: "+m.sortByAttribute,
			)
			sortableElements[i] = sortableElement{value: elem, sortKey: ""}
			continue
		}

		// Convert to string for sorting (adjust based on your needs)
		sortKeyStr := m.extractSortKey(sortKey)
		sortableElements[i] = sortableElement{
			value:   elem,
			sortKey: sortKeyStr,
		}
	}

	// Sort
	sort.Slice(sortableElements, func(i, j int) bool {
		return sortableElements[i].sortKey < sortableElements[j].sortKey
	})

	// Extract sorted values
	result := make([]attr.Value, len(sortableElements))
	for i, se := range sortableElements {
		result[i] = se.value
	}

	return result
}

func (m listSortModifier) extractSortKey(value attr.Value) string {
	// Handle different types - extend as needed
	switch v := value.(type) {
	case types.String:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return v.ValueString()
	case types.Int64:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return string(rune(v.ValueInt64())) // Simple conversion, might need better handling
	case types.Number:
		if v.IsNull() || v.IsUnknown() {
			return ""
		}
		return v.ValueBigFloat().String()
	default:
		return ""
	}
}

type sortableElement struct {
	value   attr.Value
	sortKey string
}
