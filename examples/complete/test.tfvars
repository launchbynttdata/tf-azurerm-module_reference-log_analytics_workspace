logical_product_service = "law"
logical_product_family  = "launch"
class_env               = "dev"
instance_env            = 0
instance_resource       = 0
location                = "eastus"
sku                     = "PerGB2018"
retention_in_days       = 30
resource_names_map = {
  resource_group = {
    name       = "rg"
    max_length = 80
  }
  log_analytics_workspace = {
    name       = "law"
    max_length = 80
  }
}
