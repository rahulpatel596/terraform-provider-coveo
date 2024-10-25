terraform {
  required_providers {
    coveo = {
      source = "registry.terraform.io/rahulpatel596/coveo"
    }
  }
}

provider "coveo" {}

data "coveo_coffees" "example" {}
