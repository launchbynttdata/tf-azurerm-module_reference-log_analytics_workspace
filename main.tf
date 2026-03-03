// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

module "resource_names" {
  source  = "terraform.registry.launch.nttdata.com/module_library/resource_name/launch"
  version = "~> 2.0"

  for_each = var.resource_names_map

  logical_product_family  = var.product_family
  logical_product_service = var.product_service
  region                  = var.location
  class_env               = var.class_env
  cloud_resource_type     = each.value.name
  instance_env            = var.instance_env
  maximum_length          = each.value.max_length
  instance_resource       = var.instance_resource
  use_azure_region_abbr   = var.use_azure_region_abbr
}

module "resource_group" {
  source  = "terraform.registry.launch.nttdata.com/module_primitive/resource_group/azurerm"
  version = "~> 1.0"

  name     = module.resource_names["resource_group"].standard
  location = var.location

  tags = merge(var.tags, { resource_name = module.resource_names["resource_group"].standard })
}

module "azurerm_log_analytics_workspace" {
  source  = "terraform.registry.launch.nttdata.com/module_primitive/log_analytics_workspace/azurerm"
  version = "~> 1.2"

  name                          = module.resource_names["log_analytics_workspace"].standard
  location                      = var.location
  resource_group_name           = module.resource_names["resource_group"].standard
  sku                           = var.sku
  retention_in_days             = var.retention_in_days
  local_authentication_disabled = var.local_authentication_disabled
  identity                      = var.identity
  tags                          = local.tags
  depends_on                    = [module.resource_group]
}

module "monitor_action_group" {
  count = length(var.query_alerts) > 0 && var.action_group_config != null ? 1 : 0

  source  = "terraform.registry.launch.nttdata.com/module_primitive/monitor_action_group/azurerm"
  version = "~> 1.0"

  action_group_name   = module.resource_names["monitor_action_group"].standard
  resource_group_name = module.resource_group.name
  short_name          = var.action_group_config.short_name
  arm_role_receivers  = var.action_group_config.arm_role_receivers
  email_receivers     = var.action_group_config.email_receivers

  tags = merge(local.tags, {
    resource_name = module.resource_names["monitor_action_group"].standard
  })

  depends_on = [module.resource_group]
}
module "scheduled_query_alert" {
  for_each = var.query_alerts

  source = "git::ssh://git@github.com/launchbynttdata/tf-azurerm-module_primitive-monitor_scheduled_query_alert.git?ref=feature/10489-query-alert"

  resource_group_name = module.resource_group.name
  location            = var.location

  alert_name = "${module.resource_names["scheduled_query_alert"].standard}-${each.key}"

  data_source_id = module.azurerm_log_analytics_workspace.id

  description = each.value.description
  enabled     = each.value.enabled
  query       = each.value.query
  severity    = each.value.severity
  frequency   = each.value.frequency
  time_window = each.value.time_window

  trigger_operator  = each.value.trigger_operator
  trigger_threshold = each.value.trigger_threshold

  action_group_ids = length(module.monitor_action_group) > 0 ? [module.monitor_action_group[0].action_group_id] : []

  tags = merge(local.tags, {
    resource_name = "${module.resource_names["scheduled_query_alert"].standard}-${each.key}"
  })

  depends_on = [
    module.resource_group,
    module.azurerm_log_analytics_workspace,
    module.monitor_action_group
  ]
}
