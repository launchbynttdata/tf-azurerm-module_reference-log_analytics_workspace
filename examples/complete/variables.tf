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

variable "location" {
  type        = string
  description = "Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created."
  default     = "eastus"
}

variable "sku" {
  type        = string
  description = "Specifies the SKU of the Log Analytics Workspace. Possible values are PerNode, Standalone, Unlimited, CapacityReservation, and PerGB2018. Defaults to PerGB2018. Premium and Standard do not work as per testing."
  default     = null
}

variable "retention_in_days" {
  type        = number
  description = "The workspace data retention in days. Possible values are of the range between 30 and 730."
  default     = 30
}

variable "query_alerts" {
  description = "Optional scheduled query alerts. Mirrors the root module's type; see there for field detail."
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
  description = "Optional action group config"
  type = object({
    short_name         = string
    arm_role_receivers = optional(list(any), [])
    email_receivers    = optional(list(any), [])
  })
  default = null
}
