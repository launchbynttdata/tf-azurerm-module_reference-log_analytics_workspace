# tf-azurerm-module_reference-log_analytics_workspace

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

Manages a Log Analytics (formally Operational Insights) Workspace.

## Upgrading from 1.x

2.0.0 moves the composed `monitor_scheduled_query_alert` primitive from 1.x to 2.x, which replaces
`azurerm_monitor_scheduled_query_rules_alert` with
`azurerm_monitor_scheduled_query_rules_alert_v2`. The azurerm floor moves from `~> 3.117` to
`>= 4.69, < 5.0`, because 4.69 is where the provider added `email_subject` to the v2 `action` block.

The composed `log_analytics_workspace` primitive also moves from `~> 1.2` to `~> 2.0`. That bump is
not required by the provider floor — the 1.x line already declares `azurerm >= 3.77, < 5.0` — it
just moves onto the current major. Its interface is unchanged across the bump.

### `query_alerts` is a different shape

Every alert field changed. The map now takes a nested `criteria` object instead of flat
`query` / `trigger_*` fields, ISO 8601 durations instead of minutes, and `custom_properties`
instead of a webhook JSON string:

| 1.x | 2.x |
|-----|-----|
| `query`, `trigger_operator`, `trigger_threshold` | `criteria.query`, `criteria.operator`, `criteria.threshold`, plus `criteria.time_aggregation_method` |
| `frequency` (minutes) | `evaluation_frequency` (ISO 8601, e.g. `PT5M`) |
| `time_window` (minutes) | `window_duration` (ISO 8601, e.g. `PT30M`) |
| `custom_webhook_payload` (JSON string) | `custom_properties` (`map(string)`) |
| — | `criteria.dimensions`, `criteria.failing_periods`, `criteria.metric_measure_column`, `criteria.resource_id_column` |
| — | `display_name`, `auto_mitigation_enabled`, `mute_actions_after_alert_duration`, `query_time_range_override`, `skip_query_validation`, `target_resource_types`, `workspace_alerts_storage_enabled`, `identity` |

`custom_properties` is a replacement rather than a rename: the 1.x value was the entire webhook JSON
body, while `custom_properties` is a flat map of properties added to Azure's own alert payload.

### The alerts are replaced, not updated

The underlying resource type changes, so no `moved` block can carry state — Terraform destroys each
v1 rule and creates a v2 one, losing the alert history. The workspace, resource group and action
group keep their module addresses, and the Log Analytics primitive's resource definition is the same
across 1.x and 2.x. Moving from azurerm 3.x to 4.x can still change other resources, so read your own
plan before applying. Because both types address the same ARM resource, a replacement that keeps the
same alert name will hit the provider's import check and fail while the old rule still exists. The
alert name is derived from `resource_names_map` plus the `query_alerts` key, so **changing the key is
the simplest way to get a new name**; otherwise apply the removal of the 1.x alerts before adding the
2.x ones.

### Queries are validated against the workspace

Azure checks `criteria.query` against the workspace schema on create, so a query referencing a table
the workspace has no schema for is rejected at apply time. Set `skip_query_validation` on the alert
to bypass that. A query that splits by a column must project that column, and one that aggregates a
value must project its `metric_measure_column`.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 4.69, < 5.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_azurerm_log_analytics_workspace"></a> [azurerm\_log\_analytics\_workspace](#module\_azurerm\_log\_analytics\_workspace) | terraform.registry.launch.nttdata.com/module_primitive/log_analytics_workspace/azurerm | ~> 2.0 |
| <a name="module_monitor_action_group"></a> [monitor\_action\_group](#module\_monitor\_action\_group) | terraform.registry.launch.nttdata.com/module_primitive/monitor_action_group/azurerm | ~> 1.0 |
| <a name="module_resource_group"></a> [resource\_group](#module\_resource\_group) | terraform.registry.launch.nttdata.com/module_primitive/resource_group/azurerm | ~> 1.0 |
| <a name="module_resource_names"></a> [resource\_names](#module\_resource\_names) | terraform.registry.launch.nttdata.com/module_library/resource_name/launch | ~> 2.0 |
| <a name="module_scheduled_query_alert"></a> [scheduled\_query\_alert](#module\_scheduled\_query\_alert) | terraform.registry.launch.nttdata.com/module_primitive/monitor_scheduled_query_alert/azurerm | ~> 2.0 |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_action_group_config"></a> [action\_group\_config](#input\_action\_group\_config) | Configuration for action group (created only if alerts exist) | <pre>object({<br/>    short_name         = string<br/>    arm_role_receivers = optional(list(any), [])<br/>    email_receivers    = optional(list(any), [])<br/>  })</pre> | `null` | no |
| <a name="input_class_env"></a> [class\_env](#input\_class\_env) | (Required) Environment where resource is going to be deployed. For example. dev, qa, uat | `string` | `"dev"` | no |
| <a name="input_identity"></a> [identity](#input\_identity) | (Optional) A identity block as defined below.<br/>  type: Specifies the identity type of the Log Analytics Workspace. Possible values are SystemAssigned and UserAssigned.<br/>  identity\_ids: Specifies the list of User Assigned Identity IDs to be associated with the Log Analytics Workspace. This field is required when type is UserAssigned. | <pre>object({<br/>    type         = string<br/>    identity_ids = optional(list(string))<br/>  })</pre> | `null` | no |
| <a name="input_instance_env"></a> [instance\_env](#input\_instance\_env) | (Optional) Number that represents the instance of the environment. | `number` | `0` | no |
| <a name="input_instance_resource"></a> [instance\_resource](#input\_instance\_resource) | (Optional) Number that represents the instance of the resource. | `number` | `0` | no |
| <a name="input_local_authentication_disabled"></a> [local\_authentication\_disabled](#input\_local\_authentication\_disabled) | (Optional) Boolean flag to specify whether local authentication should be disabled. Defaults to false. | `bool` | `false` | no |
| <a name="input_location"></a> [location](#input\_location) | (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created. | `string` | n/a | yes |
| <a name="input_product_family"></a> [product\_family](#input\_product\_family) | (Required) Name of the product family for which the resource is created | `string` | `"launch"` | no |
| <a name="input_product_service"></a> [product\_service](#input\_product\_service) | (Required) Name of the product service for which the resource is created | `string` | `"law"` | no |
| <a name="input_query_alerts"></a> [query\_alerts](#input\_query\_alerts) | Map of scheduled query alerts to create against this workspace. Each key becomes a name suffix on<br/>the generated alert name. The shape follows azurerm\_monitor\_scheduled\_query\_rules\_alert\_v2; see the<br/>monitor\_scheduled\_query\_alert primitive for the per-field detail, which this module passes through<br/>and does not re-validate.<br/><br/>criteria.query must be valid against this workspace's schema, because Azure checks it on create<br/>unless skip\_query\_validation is true. A query that splits by a column must project that column, and<br/>one that aggregates a value must project the metric\_measure\_column. | <pre>map(object({<br/>    description  = optional(string, "")<br/>    display_name = optional(string)<br/>    enabled      = optional(bool, true)<br/>    severity     = optional(number, 1)<br/><br/>    evaluation_frequency = optional(string, "PT5M")<br/>    window_duration      = optional(string, "PT30M")<br/><br/>    criteria = object({<br/>      query                   = string<br/>      threshold               = number<br/>      operator                = optional(string, "GreaterThan")<br/>      time_aggregation_method = optional(string, "Count")<br/>      metric_measure_column   = optional(string)<br/>      resource_id_column      = optional(string)<br/>      dimensions = optional(list(object({<br/>        name     = string<br/>        values   = list(string)<br/>        operator = optional(string, "Include")<br/>      })), [])<br/>      failing_periods = optional(object({<br/>        minimum_failing_periods_to_trigger_alert = number<br/>        number_of_evaluation_periods             = number<br/>      }))<br/>    })<br/><br/>    auto_mitigation_enabled           = optional(bool, false)<br/>    mute_actions_after_alert_duration = optional(string)<br/>    query_time_range_override         = optional(string)<br/>    skip_query_validation             = optional(bool, false)<br/>    target_resource_types             = optional(list(string), [])<br/>    workspace_alerts_storage_enabled  = optional(bool, false)<br/><br/>    email_subject     = optional(string)<br/>    custom_properties = optional(map(string), {})<br/><br/>    identity = optional(object({<br/>      type         = string<br/>      identity_ids = optional(list(string), [])<br/>    }))<br/>  }))</pre> | `{}` | no |
| <a name="input_resource_names_map"></a> [resource\_names\_map](#input\_resource\_names\_map) | (Optional) A map of key to resource\_name that will be used by tf-launch-module\_library-resource\_name to generate resource names | <pre>map(object({<br/>    name       = string<br/>    max_length = optional(number, 60)<br/>  }))</pre> | <pre>{<br/>  "log_analytics_workspace": {<br/>    "max_length": 60,<br/>    "name": "law"<br/>  },<br/>  "monitor_action_group": {<br/>    "max_length": 60,<br/>    "name": "mag"<br/>  },<br/>  "resource_group": {<br/>    "max_length": 60,<br/>    "name": "rg"<br/>  },<br/>  "scheduled_query_alert": {<br/>    "max_length": 60,<br/>    "name": "sqa"<br/>  }<br/>}</pre> | no |
| <a name="input_retention_in_days"></a> [retention\_in\_days](#input\_retention\_in\_days) | (Optional) The workspace data retention in days. Possible values are in the range between 30 and 730. | `number` | `30` | no |
| <a name="input_sku"></a> [sku](#input\_sku) | (Optional) Specifies the SKU of the Log Analytics Workspace. PerGB2018 and CapacityReservation are recommended. Other possible values are PerNode, Standalone, and Unlimited. Defaults to PerGB2018. | `string` | `null` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | (Optional) A mapping of tags to assign to the resource. | `map(string)` | `{}` | no |
| <a name="input_use_azure_region_abbr"></a> [use\_azure\_region\_abbr](#input\_use\_azure\_region\_abbr) | (Optional) Whether to use Azure region abbreviation e.g. eastus -> eus | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_action_group_id"></a> [action\_group\_id](#output\_action\_group\_id) | The ID of the Monitor Action Group |
| <a name="output_action_group_name"></a> [action\_group\_name](#output\_action\_group\_name) | The name of the Monitor Action Group |
| <a name="output_id"></a> [id](#output\_id) | The Log Analytics resource ID. |
| <a name="output_name"></a> [name](#output\_name) | The Log Analytics Workspace name. |
| <a name="output_primary_shared_key"></a> [primary\_shared\_key](#output\_primary\_shared\_key) | Value of the primary shared key for the Log Analytics Workspace. |
| <a name="output_resource_group_name"></a> [resource\_group\_name](#output\_resource\_group\_name) | The Log Analytics resource group name |
| <a name="output_scheduled_query_alerts"></a> [scheduled\_query\_alerts](#output\_scheduled\_query\_alerts) | Map of scheduled query alert IDs, names and identity principal IDs. principal\_id is null unless the alert was given a SystemAssigned identity. |
| <a name="output_secondary_shared_key"></a> [secondary\_shared\_key](#output\_secondary\_shared\_key) | Value of the secondary shared key for the Log Analytics Workspace. |
| <a name="output_workspace_id"></a> [workspace\_id](#output\_workspace\_id) | The Workspace (or Customer) ID for the Log Analytics Workspace. |
<!-- END_TF_DOCS -->

## Module Development

### Pre-Requisites

The following commands should be available on your system:

- `asdf` or `mise`
- `make`
- `python3` (for pre-commit)

Additionally, your `git` user and email must be configured. Run the `make configure` command from the root of the repository to ensure that you meet these requirements.

### Pre-Commit hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to Terraform and Golang, as well as some common linting tasks. These will be configured for you when you run `make configure`.

### Local Validation

You should validate the changes you make to any module locally, prior to pushing your changes in a branch to GitHub.

1. Ensure that you have run `make configure` successfully.

2. Ensure you are signed into the appropriate cloud provider (e.g. AWS or Azure) for the module under test in your current console session.

3. Run the Terraform and Golang linters with the following command:

```
make lint
```

4. Once you have satisfied the linters, the following command will build example infrastructure in your configured cloud, run the tests, and then tear down the infrastructure it created:

```
make test
```

The pre-commit validations, as well as the `make lint` and `make test` targets, will all be performed in CI. Running these validations locally prior to opening a PR helps ensure a smooth review and merge process.

### Review & Merge Process

Once your change has been tested locally and your branch pushed up, open a new Pull Request for your branch to the default (main) branch of this repository.

The title of your Pull Request will determine the version bump for this change, and the title must be in [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) format in order to merge. A breaking change will trigger a major version bump, a feature will trigger a minor version bump, and all other types will trigger a patch version bump.

Ensure your CI workflows are passing; seek approval from teammates and address any feedback; seek any explicit approvals required by the CODEOWNERS file. You may merge the PR as soon as all requirements are met, and a new release and tag will be automatically created for you.

### Automatic Updates

The shared configuration and workflow files in this repository are largely managed through the [launch-terraform-skeleton](https://github.com/launchbynttdata/launch-terraform-skeleton) repository. Outside of perhaps the `.gitignore` to account for specific files being generated by certain Terraform modules (e.g. Lambda functions), there should not be much cause to update these files on a per-repo basis, and making changes to them individually is discouraged.

If desired, you can check for and run these updates locally in a branch if you have the `copier` tool installed. Some example commands are included below:

```
# Check for updates, optionally checking prerelease versions
copier check-update [--prereleases]

# Run an update, using default answers if there are any. We use tasks, which requires --trust to be set.
copier update --defaults --trust [--prereleases]

# Recopy from the source, and --overwrite all templated files in the process
copier recopy --defaults --trust --overwrite [--prereleases]
```

Automatic updates will run through a scheduled workflow, and if the post-update tests are successful, the Pull Request created will automatically merge. Conflicts in the update or failures to test may leave a Pull Request outstanding, which needs to be addressed by a Launch Engineer.
