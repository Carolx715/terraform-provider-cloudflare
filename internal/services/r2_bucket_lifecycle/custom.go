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

	// Get the elements from the plan - use concrete type basetypes.ObjectValue
	var planElements []basetypes.ObjectValue
	diags := req.PlanValue.ElementsAs(ctx, &planElements, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Sort the elements
	sortedElements := m.sortElements(ctx, planElements, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert back to []attr.Value for creating the new list
	sortedAttrValues := make([]attr.Value, len(sortedElements))
	for i, elem := range sortedElements {
		sortedAttrValues[i] = elem
	}

	// Create a new list with sorted elements
	sortedList, diags := types.ListValue(req.PlanValue.ElementType(ctx), sortedAttrValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = sortedList
}

func (m listSortModifier) sortElements(ctx context.Context, elements []basetypes.ObjectValue, diags *diag.Diagnostics) []basetypes.ObjectValue {
	// Create a sortable slice
	sortableElements := make([]sortableElement, len(elements))

	for i, elem := range elements {
		// Extract the sort key from the object's attributes
		sortKey := elem.Attributes()[m.sortByAttribute]
		if sortKey == nil {
			diags.AddWarning(
				"Missing Sort Attribute",
				"Element missing attribute: "+m.sortByAttribute,
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

	// Sort
	sort.Slice(sortableElements, func(i, j int) bool {
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
		return fmt.Sprint(v.ValueInt64())
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
