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

terraform {
  required_version = "~> 1.5"

  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      # monitor_scheduled_query_alert 2.x requires 4.69, where the provider added
      # email_subject to the action block of the v2 scheduled query rule.
      version = ">= 4.69, < 5.0"
    }
  }
}
