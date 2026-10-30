location = "eastus"

sku               = "PerGB2018"
retention_in_days = 30

# Usage is a standard table present in every Log Analytics workspace, so Azure can validate the
# query against a freshly created one. The query projects DataType because the alert splits the
# evaluation by it, and Quantity because the aggregation measures that column.
query_alerts = {
  ingestion = {
    description  = "Alert when a data type ingests more than the expected volume"
    display_name = "Ingestion volume by data type"
    enabled      = true
    severity     = 2

    evaluation_frequency = "PT5M"
    window_duration      = "PT30M"

    criteria = {
      query                   = "Usage | project TimeGenerated, DataType, Quantity"
      operator                = "GreaterThan"
      threshold               = 1000
      time_aggregation_method = "Total"
      metric_measure_column   = "Quantity"

      dimensions = [
        {
          name     = "DataType"
          values   = ["*"]
          operator = "Include"
        }
      ]

      failing_periods = {
        minimum_failing_periods_to_trigger_alert = 1
        number_of_evaluation_periods             = 1
      }
    }

    # Re-notify hourly while ingestion stays high. Azure rejects muting unless auto mitigation is off.
    auto_mitigation_enabled           = false
    mute_actions_after_alert_duration = "PT1H"

    email_subject = "Azure Alert: Log Analytics ingestion volume"

    custom_properties = {
      alert_source = "reference_log_analytics_workspace"
    }
  }
}

action_group_config = {
  short_name = "AlertTeam"
  email_receivers = [{
    name                    = "Email"
    email_address           = "alerts@example.com"
    use_common_alert_schema = true
  }]
}
