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

resource "coveo_document" "example" {
  title = "Test Index 95"
  content = "This is a test index"
  source_id = "pwcitlawbs3xfhj3dspx2yh2cbi-wndwuafzh7ml7oiyacjka62lw4"
}
