# tf-azurerm-module_reference-log_analytics_workspace

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.0 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | ~> 3.117 |

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

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The Log Analytics resource ID. |
| <a name="output_workspace_id"></a> [workspace\_id](#output\_workspace\_id) | The Workspace (or Customer) ID for the Log Analytics Workspace. |
| <a name="output_workspace_name"></a> [workspace\_name](#output\_workspace\_name) | The Workspace name for the Log Analytics Workspace. |
| <a name="output_resource_group_name"></a> [resource\_group\_name](#output\_resource\_group\_name) | The Log Analytics Workspace resource group name |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
