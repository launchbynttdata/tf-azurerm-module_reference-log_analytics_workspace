# empty

query_alerts = {
  test_alert = {
    description       = "Test alert for AzureActivity"
    enabled           = true
    query             = "AzureActivity | take 5"
    severity          = 2
    frequency         = 5
    time_window       = 30
    trigger_operator  = "GreaterThan"
    trigger_threshold = 1
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
