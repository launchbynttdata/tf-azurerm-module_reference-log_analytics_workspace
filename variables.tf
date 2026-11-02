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

variable "sku" {
  type        = string
  description = "(Optional) Specifies the SKU of the Log Analytics Workspace. PerGB2018 and CapacityReservation are recommended. Other possible values are PerNode, Standalone, and Unlimited. Defaults to PerGB2018."
  default     = null
  validation {
    condition     = var.sku == null || can(regex("^(PerNode|Standalone|Unlimited|CapacityReservation|PerGB2018)$", var.sku))
    error_message = "SKU must be one of 'PerNode', 'Standalone', 'Unlimited', 'CapacityReservation', and 'PerGB2018'"
  }
}

variable "retention_in_days" {
  type        = number
  description = "(Optional) The workspace data retention in days. Possible values are in the range between 30 and 730."
  default     = 30
  validation {
    condition     = var.retention_in_days >= 30 && var.retention_in_days <= 730
    error_message = "retention_in_days should be between 30 to 730."
  }
}

variable "local_authentication_disabled" {
  type        = bool
  description = "(Optional) Boolean flag to specify whether local authentication should be disabled. Defaults to false."
  default     = false
}

variable "identity" {
  type = object({
    type         = string
    identity_ids = optional(list(string))
  })
  description = <<EOF
  (Optional) A identity block as defined below.
  type: Specifies the identity type of the Log Analytics Workspace. Possible values are SystemAssigned and UserAssigned.
  identity_ids: Specifies the list of User Assigned Identity IDs to be associated with the Log Analytics Workspace. This field is required when type is UserAssigned.
  EOF
  default     = null
}

variable "location" {
  type        = string
  description = "(Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created."
}

variable "resource_names_map" {
  description = "(Optional) A map of key to resource_name that will be used by tf-launch-module_library-resource_name to generate resource names"
  type = map(object({
    name       = string
    max_length = optional(number, 60)
  }))

  default = {
    log_analytics_workspace = {
      name       = "law"
      max_length = 60
    }
    resource_group = {
      name       = "rg"
      max_length = 60
    }
    monitor_action_group = {
      name       = "mag"
      max_length = 60
    }
    scheduled_query_alert = {
      name       = "sqa"
      max_length = 60
    }
  }
}

variable "instance_env" {
  type        = number
  description = "(Optional) Number that represents the instance of the environment."
  default     = 0

  validation {
    condition     = var.instance_env >= 0 && var.instance_env <= 999
    error_message = "Instance number should be between 0 to 999."
  }
}

variable "instance_resource" {
  type        = number
  description = "(Optional) Number that represents the instance of the resource."
  default     = 0

  validation {
    condition     = var.instance_resource >= 0 && var.instance_resource <= 100
    error_message = "Instance number should be between 0 to 100."
  }
}

variable "product_family" {
  type        = string
  description = "(Required) Name of the product family for which the resource is created"
  nullable    = false

  validation {
    condition     = can(regex("^[_\\-A-Za-z0-9]+$", var.product_family))
    error_message = "The variable must contain letters, numbers, -, _, and .."
  }

  default = "launch"
}

variable "product_service" {
  type        = string
  description = "(Required) Name of the product service for which the resource is created"
  nullable    = false

  validation {
    condition     = can(regex("^[_\\-A-Za-z0-9]+$", var.product_service))
    error_message = "The variable must contain letters, numbers, -, _, and .."
  }

  default = "law"
}

variable "class_env" {
  type        = string
  description = "(Required) Environment where resource is going to be deployed. For example. dev, qa, uat"
  nullable    = false
  default     = "dev"

  validation {
    condition     = length(regexall("\\b \\b", var.class_env)) == 0
    error_message = "Spaces between the words are not allowed."
  }
}

variable "use_azure_region_abbr" {
  description = "(Optional) Whether to use Azure region abbreviation e.g. eastus -> eus"
  type        = bool
  default     = true
}

variable "tags" {
  type        = map(string)
  description = "(Optional) A mapping of tags to assign to the resource."
  default     = {}
}

variable "query_alerts" {
  description = <<-EOT
    Map of scheduled query alerts to create against this workspace. Each key becomes a name suffix on
    the generated alert name. The shape follows azurerm_monitor_scheduled_query_rules_alert_v2; see the
    monitor_scheduled_query_alert primitive for the per-field detail, which this module passes through
    and does not re-validate.

    criteria.query must be valid against this workspace's schema, because Azure checks it on create
    unless skip_query_validation is true. A query that splits by a column must project that column, and
    one that aggregates a value must project the metric_measure_column.
  EOT

  type = map(object({
    description  = optional(string, "")
    display_name = optional(string)
    enabled      = optional(bool, true)
    severity     = optional(number, 1)

    evaluation_frequency = optional(string, "PT5M")
    window_duration      = optional(string, "PT30M")

    criteria = object({
      query                   = string
      threshold               = number
      operator                = optional(string, "GreaterThan")
      time_aggregation_method = optional(string, "Count")
      metric_measure_column   = optional(string)
      resource_id_column      = optional(string)
      dimensions = optional(list(object({
        name     = string
        values   = list(string)
        operator = optional(string, "Include")
      })), [])
      failing_periods = optional(object({
        minimum_failing_periods_to_trigger_alert = number
        number_of_evaluation_periods             = number
      }))
    })

    auto_mitigation_enabled           = optional(bool, false)
    mute_actions_after_alert_duration = optional(string)
    query_time_range_override         = optional(string)
    skip_query_validation             = optional(bool, false)
    target_resource_types             = optional(list(string), [])
    workspace_alerts_storage_enabled  = optional(bool, false)

    email_subject     = optional(string)
    custom_properties = optional(map(string), {})

    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))
  }))

  default = {}
}

variable "action_group_config" {
  description = "Configuration for action group (created only if alerts exist)"
  type = object({
    short_name         = string
    arm_role_receivers = optional(list(any), [])
    email_receivers    = optional(list(any), [])
  })
  default = null
}
