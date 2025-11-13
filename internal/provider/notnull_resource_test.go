// Copyright (c) Plain Technologies Aps

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotNullResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with value provided
			{
				Config: testAccNotNullResourceConfig("test_value", "default_value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data_notnull.test", "value", "test_value"),
					resource.TestCheckResourceAttr("data_notnull.test", "default_value", "default_value"),
					resource.TestCheckResourceAttr("data_notnull.test", "result", "test_value"),
					resource.TestCheckResourceAttr("data_notnull.test", "id", "notnull"),
				),
			},
			// Update testing - change value
			{
				Config: testAccNotNullResourceConfig("updated_value", "default_value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data_notnull.test", "value", "updated_value"),
					resource.TestCheckResourceAttr("data_notnull.test", "result", "updated_value"),
				),
			},
			// Test state preservation - remove value, should preserve previous result
			{
				Config: testAccNotNullResourceConfigDefaultOnly("default_only"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data_notnull.test", "default_value", "default_only"),
					// Result should preserve "updated_value" from previous state when value becomes null
					resource.TestCheckResourceAttr("data_notnull.test", "result", "updated_value"),
				),
			},
		},
	})
}

func TestAccNotNullResourceDefaultOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test creating resource with only default value (no prior state)
			{
				Config: testAccNotNullResourceConfigDefaultOnly("default_only"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data_notnull.test", "default_value", "default_only"),
					resource.TestCheckResourceAttr("data_notnull.test", "result", "default_only"),
				),
			},
		},
	})
}

func testAccNotNullResourceConfig(value, defaultValue string) string {
	return fmt.Sprintf(`
resource "data_notnull" "test" {
  value         = %[1]q
  default_value = %[2]q
}
`, value, defaultValue)
}

func testAccNotNullResourceConfigDefaultOnly(defaultValue string) string {
	return fmt.Sprintf(`
resource "data_notnull" "test" {
  default_value = %[1]q
}
`, defaultValue)
}

// Unit tests for computeResult logic
func TestComputeResult(t *testing.T) {
	r := &NotNullResource{}

	tests := []struct {
		name        string
		data        *NotNullResourceModel
		priorState  *NotNullResourceModel
		expected    string
		description string
	}{
		{
			name: "value provided",
			data: &NotNullResourceModel{
				Value:        types.StringValue("test_value"),
				DefaultValue: types.StringValue("default"),
			},
			priorState:  nil,
			expected:    "test_value",
			description: "When value is provided, use value",
		},
		{
			name: "value is null with prior state",
			data: &NotNullResourceModel{
				Value:        types.StringNull(),
				DefaultValue: types.StringValue("default"),
			},
			priorState: &NotNullResourceModel{
				Result: types.StringValue("previous_result"),
			},
			expected:    "previous_result",
			description: "When value is null and prior state exists, preserve previous result",
		},
		{
			name: "value is unknown with prior state",
			data: &NotNullResourceModel{
				Value:        types.StringUnknown(),
				DefaultValue: types.StringValue("default"),
			},
			priorState: &NotNullResourceModel{
				Result: types.StringValue("previous_result"),
			},
			expected:    "previous_result",
			description: "When value is unknown and prior state exists, preserve previous result",
		},
		{
			name: "value is null without prior state",
			data: &NotNullResourceModel{
				Value:        types.StringNull(),
				DefaultValue: types.StringValue("default"),
			},
			priorState:  nil,
			expected:    "default",
			description: "When value is null and no prior state, use default_value",
		},
		{
			name: "value is unknown without prior state",
			data: &NotNullResourceModel{
				Value:        types.StringUnknown(),
				DefaultValue: types.StringValue("default"),
			},
			priorState:  nil,
			expected:    "default",
			description: "When value is unknown and no prior state, use default_value",
		},
		{
			name: "everything is null",
			data: &NotNullResourceModel{
				Value:        types.StringNull(),
				DefaultValue: types.StringNull(),
			},
			priorState:  nil,
			expected:    "",
			description: "When everything is null, return empty string",
		},
		{
			name: "value unknown, default null, no prior state",
			data: &NotNullResourceModel{
				Value:        types.StringUnknown(),
				DefaultValue: types.StringNull(),
			},
			priorState:  nil,
			expected:    "",
			description: "When value is unknown, default is null, and no prior state, return empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.computeResult(tt.data, tt.priorState)
			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.description, tt.expected, result)
			}
		})
	}
}
