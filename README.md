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
v1 rule and creates a v2 one, losing alert history but nothing else. Because both types address the
same ARM resource, a replacement that keeps the same alert name will hit the provider's import check
and fail while the old rule still exists. The alert name is derived from `resource_names_map` plus
the `query_alerts` key, so **changing the key is the simplest way to get a new name**; otherwise
apply the removal of the 1.x alerts before adding the 2.x ones.

### Queries are validated against the workspace

Azure checks `criteria.query` against the workspace schema on create, so a query referencing a table
the workspace has no schema for is rejected at apply time. Set `skip_query_validation` on the alert
to bypass that. A query that splits by a column must project that column, and one that aggregates a
value must project its `metric_measure_column`.

## Pre-Commit hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to terraform, golang and common linting tasks. There are no custom hooks added.

`commitlint` hook enforces commit message in certain format. The commit contains the following structural elements, to communicate intent to the consumers of your commit messages:

- **fix**: a commit of the type `fix` patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- **feat**: a commit of the type `feat` introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- **BREAKING CHANGE**: a commit that has a footer `BREAKING CHANGE:`, or appends a `!` after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.
- **build**: a commit of the type `build` adds changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- **chore**: a commit of the type `chore` adds changes that don't modify src or test files
- **ci**: a commit of the type `ci` adds changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- **docs**: a commit of the type `docs` adds documentation only changes
- **perf**: a commit of the type `perf` adds code change that improves performance
- **refactor**: a commit of the type `refactor` adds code change that neither fixes a bug nor adds a feature
- **revert**: a commit of the type `revert` reverts a previous commit
- **style**: a commit of the type `style` adds code changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **test**: a commit of the type `test` adds missing tests or correcting existing tests

Base configuration used for this project is [commitlint-config-conventional (based on the Angular convention)](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional#type-enum)

If you are a developer using vscode, [this](https://marketplace.visualstudio.com/items?itemName=joshbolduc.commitlint) plugin may be helpful.

`detect-secrets-hook` prevents new secrets from being introduced into the baseline. TODO: INSERT DOC LINK ABOUT HOOKS

In order for `pre-commit` hooks to work properly

- You need to have the pre-commit package manager installed. [Here](https://pre-commit.com/#install) are the installation instructions.
- `pre-commit` would install all the hooks when commit message is added by default except for `commitlint` hook. `commitlint` hook would need to be installed manually using the command below

```
pre-commit install --hook-type commit-msg
```

## To test the resource group module locally

1. For development/enhancements to this module locally, you'll need to install all of its components. This is controlled by the `configure` target in the project's [`Makefile`](./Makefile). Before you can run `configure`, familiarize yourself with the variables in the `Makefile` and ensure they're pointing to the right places.

```
make configure
```

This adds in several files and directories that are ignored by `git`. They expose many new Make targets.

2. _THIS STEP APPLIES ONLY TO MICROSOFT AZURE. IF YOU ARE USING A DIFFERENT PLATFORM PLEASE SKIP THIS STEP._ The first target you care about is `env`. This is the common interface for setting up environment variables. The values of the environment variables will be used to authenticate with cloud provider from local development workstation.

`make configure` command will bring down `azure_env.sh` file on local workstation. Devloper would need to modify this file, replace the environment variable values with relevant values.

These environment variables are used by `terratest` integration suit.

Service principle used for authentication(value of ARM_CLIENT_ID) should have below privileges on resource group within the subscription.

```
"Microsoft.Resources/subscriptions/resourceGroups/write"
"Microsoft.Resources/subscriptions/resourceGroups/read"
"Microsoft.Resources/subscriptions/resourceGroups/delete"
```

Then run this make target to set the environment variables on developer workstation.

```
make env
```

3. The first target you care about is `check`.

**Pre-requisites**
Before running this target it is important to ensure that, developer has created files mentioned below on local workstation under root directory of git repository that contains code for primitives/segments. Note that these files are `azure` specific. If primitive/segment under development uses any other cloud provider than azure, this section may not be relevant.

- A file named `provider.tf` with contents below

```
provider "azurerm" {
  features {}
}
```

- A file named `terraform.tfvars` which contains key value pair of variables used.

Note that since these files are added in `gitignore` they would not be checked in into primitive/segment's git repo.

After creating these files, for running tests associated with the primitive/segment, run

```
make check
```

If `make check` target is successful, developer is good to commit the code to primitive/segment's git repo.

`make check` target

- runs `terraform commands` to `lint`,`validate` and `plan` terraform code.
- runs `conftests`. `conftests` make sure `policy` checks are successful.
- runs `terratest`. This is integration test suit.
- runs `opa` tests
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
