# Terraform Provider Data

Simple data manipulation resources for Terraform.

## Description

This provider implements data manipulation resources for Terraform using the [HashiCorp Plugin Framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework).

## Resources

### `data_notnull`

The `data_notnull` resource provides a way to handle nullable values with fallback logic and state preservation.

#### Schema

- **Inputs:**
  - `value` (string, optional) - The primary value to use for the result will keep last non-null value when changed to null
  - `default_value` (string, optional) - The default value to use when value is and no previous value exists

- **Outputs:**
  - `result` (string, computed) - The computed result value
  - `id` (string, computed) - Internal identifier

#### Behavior

The resource determines the `result` value using the following logic:

1. If `value` is not null, `result` equals `value`
2. If `value` is null and there is stored state (meaning value changed from non-null to null), `result` preserves the previous stored value
3. If `value` is null and there is no stored state, `result` equals `default_value`
4. If all inputs are null, `result` is an empty string

#### Example Usage

```hcl
terraform {
  required_providers {
    data = {
      source = "plain-insure/data"
    }
  }
}

provider "data" {}

# Example with value provided
resource "data_notnull" "example1" {
  value         = "hello"
  default_value = "default"
}

output "result1" {
  value = data_notnull.example1.result  # Output: "hello"
}

# Example with only default_value (no value)
resource "data_notnull" "example2" {
  default_value = "fallback"
}

output "result2" {
  value = data_notnull.example2.result  # Output: "fallback"
}

# Example showing state preservation
# If you first apply with value = "initial", then change to value = null,
# the result will still be "initial" (preserved from state)
resource "data_notnull" "example3" {
  value         = "initial"
  default_value = "default"
}
```

## Building the Provider

```bash
go build -o terraform-provider-data
```

## Testing

Run the unit tests:

```bash
go test ./internal/provider -v
```

Run acceptance tests (requires `TF_ACC=1`):

```bash
TF_ACC=1 go test ./internal/provider -v
```

## Development

To use a local development version of the provider, create a `~/.terraformrc` file with:

```hcl
provider_installation {
  dev_overrides {
    "plain-insure/data" = "/path/to/terraform-provider-data"
  }

  direct {}
}
```

Then build the provider and run Terraform commands as usual.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
