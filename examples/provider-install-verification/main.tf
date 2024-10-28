terraform {
  required_providers {
    coveo = {
      source  = "registry.terraform.io/rahulpatel596/coveo"
      version = "1.0.0"
    }
  }
}

provider "coveo" {
  organization_id = "pwcitlawbs3xfhj3dspx2yh2cbi"
  api_key = "xx6a366fe8-fe0e-49f0-9c51-d033ad1c2856"
}

resource "coveo_index" "example" {
  name = "Test Index 95"
}
