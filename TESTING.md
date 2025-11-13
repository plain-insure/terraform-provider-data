# Testing Results

This document shows the testing performed on the terraform-provider-data provider.

## Manual Testing

### Test Setup

The provider was tested using Terraform's dev_overrides feature to use the locally built binary.

### Test Scenarios

#### Scenario 1: Value Provided

```hcl
resource "data_notnull" "test" {
  value         = "from_value"
  default_value = "from_default"
}
```

**Result:** `result = "from_value"` ✅

**Explanation:** When value is provided, it is used as the result.

---

#### Scenario 2: Only Default Value

```hcl
resource "data_notnull" "test" {
  default_value = "only_default"
}
```

**Result:** `result = "only_default"` ✅

**Explanation:** When value is null and there's no prior state, default_value is used.

---

#### Scenario 3: State Preservation (Critical Test)

**Initial Apply:**
```hcl
resource "data_notnull" "test" {
  value         = "initial_value"
  default_value = "backup_value"
}
```

**Result:** `result = "initial_value"` ✅

**Updated Configuration (value removed):**
```hcl
resource "data_notnull" "test" {
  # value removed - becomes null
  default_value = "backup_value"
}
```

**Result:** `result = "initial_value"` ✅

**Explanation:** This is the key requirement! When value changes from non-null to null, the resource preserves the previous value from state instead of falling back to default_value. This demonstrates that the state management logic is working correctly.

---

## Unit Tests

Unit tests are located in `internal/provider/notnull_resource_test.go`.

To run tests:
```bash
go test ./internal/provider -v
```

For acceptance tests (requires TF_ACC=1):
```bash
TF_ACC=1 go test ./internal/provider -v
```

## Security Scanning

### CodeQL Analysis
- **Status:** ✅ PASSED
- **Alerts:** 0 vulnerabilities found

### Dependency Scanning
- **Status:** ✅ PASSED  
- **Key Dependencies Scanned:**
  - github.com/hashicorp/terraform-plugin-framework v1.16.1
  - github.com/hashicorp/terraform-plugin-go v0.29.0
  - github.com/hashicorp/terraform-plugin-testing v1.13.3
- **Vulnerabilities:** None found

### Go Vet
- **Status:** ✅ PASSED
- **Issues:** None found

### Code Formatting
- **Status:** ✅ PASSED
- **All files properly formatted:** Yes

## Build Verification

```bash
$ go build -o terraform-provider-data
# Build successful - no errors
```

## Conclusion

All tests passed successfully. The provider correctly implements the required behavior:

1. ✅ Uses `value` when provided
2. ✅ Uses `default_value` when value is null and no prior state exists
3. ✅ **Preserves previous state value when value changes from non-null to null** (key requirement)
4. ✅ No security vulnerabilities detected
5. ✅ Code follows Go best practices
