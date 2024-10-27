terraform {
  required_providers {
    coveo = {
      source = "registry.terraform.io/rahulpatel596/coveo"
      version = "1.0.0"
    }
  }
}

provider "coveo" {
  api_key         = "xx6a366fe8-fe0e-49f0-9c51-d033ad1c2856"
  organization_id = "pwcitlawbs3xfhj3dspx2yh2cbi"
}

resource "coveo_index" "example" {
  name = "Test Index"
}
