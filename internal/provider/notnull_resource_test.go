package provider

import (
	"fmt"
	"testing"

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
			// Test with only default value
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
