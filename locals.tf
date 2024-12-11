locals {
  default_tags = {
    provisioner = "terraform"
  }
  tags = merge(local.default_tags, var.tags)
}
