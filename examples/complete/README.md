# tf-azurerm-module_reference-log_analytics_workspace

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 4.69, < 5.0 |

## Providers

No providers.

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_log_analytics_workspace"></a> [log\_analytics\_workspace](#module\_log\_analytics\_workspace) | ../.. | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_location"></a> [location](#input\_location) | Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created. | `string` | `"eastus"` | no |
| <a name="input_sku"></a> [sku](#input\_sku) | Specifies the SKU of the Log Analytics Workspace. Possible values are PerNode, Standalone, Unlimited, CapacityReservation, and PerGB2018. Defaults to PerGB2018. Premium and Standard do not work as per testing. | `string` | `null` | no |
| <a name="input_retention_in_days"></a> [retention\_in\_days](#input\_retention\_in\_days) | The workspace data retention in days. Possible values are of the range between 30 and 730. | `number` | `30` | no |
| <a name="input_query_alerts"></a> [query\_alerts](#input\_query\_alerts) | Optional scheduled query alerts. Mirrors the root module's type; see there for field detail. | <pre>map(object({<br/>    description  = optional(string, "")<br/>    display_name = optional(string)<br/>    enabled      = optional(bool, true)<br/>    severity     = optional(number, 1)<br/><br/>    evaluation_frequency = optional(string, "PT5M")<br/>    window_duration      = optional(string, "PT30M")<br/><br/>    criteria = object({<br/>      query                   = string<br/>      threshold               = number<br/>      operator                = optional(string, "GreaterThan")<br/>      time_aggregation_method = optional(string, "Count")<br/>      metric_measure_column   = optional(string)<br/>      resource_id_column      = optional(string)<br/>      dimensions = optional(list(object({<br/>        name     = string<br/>        values   = list(string)<br/>        operator = optional(string, "Include")<br/>      })), [])<br/>      failing_periods = optional(object({<br/>        minimum_failing_periods_to_trigger_alert = number<br/>        number_of_evaluation_periods             = number<br/>      }))<br/>    })<br/><br/>    auto_mitigation_enabled           = optional(bool, false)<br/>    mute_actions_after_alert_duration = optional(string)<br/>    query_time_range_override         = optional(string)<br/>    skip_query_validation             = optional(bool, false)<br/>    target_resource_types             = optional(list(string), [])<br/>    workspace_alerts_storage_enabled  = optional(bool, false)<br/><br/>    email_subject     = optional(string)<br/>    custom_properties = optional(map(string), {})<br/><br/>    identity = optional(object({<br/>      type         = string<br/>      identity_ids = optional(list(string), [])<br/>    }))<br/>  }))</pre> | `{}` | no |
| <a name="input_action_group_config"></a> [action\_group\_config](#input\_action\_group\_config) | Optional action group config | <pre>object({<br/>    short_name         = string<br/>    arm_role_receivers = optional(list(any), [])<br/>    email_receivers    = optional(list(any), [])<br/>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The Log Analytics resource ID. |
| <a name="output_workspace_id"></a> [workspace\_id](#output\_workspace\_id) | The Workspace (or Customer) ID for the Log Analytics Workspace. |
| <a name="output_workspace_name"></a> [workspace\_name](#output\_workspace\_name) | The Workspace name for the Log Analytics Workspace. |
| <a name="output_resource_group_name"></a> [resource\_group\_name](#output\_resource\_group\_name) | The Log Analytics Workspace resource group name |
| <a name="output_action_group_id"></a> [action\_group\_id](#output\_action\_group\_id) | The ID of the Monitor Action Group |
| <a name="output_action_group_name"></a> [action\_group\_name](#output\_action\_group\_name) | The name of the Monitor Action Group |
| <a name="output_scheduled_query_alerts"></a> [scheduled\_query\_alerts](#output\_scheduled\_query\_alerts) | Map of scheduled query alert IDs, names and identity principal IDs |
<!-- END_TF_DOCS -->
