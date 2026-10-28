# avm_interface_location

Checks that a module's `location` variable uses the AVM location interface.

## Telemetry-only location

Modules that use `location` only for telemetry may make it optional:

```hcl
variable "location" {
  type    = string
  default = null

  validation {
    condition     = !var.enable_telemetry || var.location != null
    error_message = "The location must not be null when telemetry is enabled."
  }
}
```

`nullable = true` is also accepted, but is unnecessary because it is Terraform's default. Modules that use `location` for resources or other non-telemetry behavior must retain the required, non-nullable location interface.
