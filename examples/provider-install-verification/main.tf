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

resource "coveo_document" "example6" {
  title = "Test Index 01"
  content = "This is a test index"
  source_id = "pwcitlawbs3xfhj3dspx2yh2cbi-wndwuafzh7ml7oiyacjka62lw4"
  document_id = "product://006-green1"
}

resource "coveo_document" "example5" {
  title = "Test Index 02"
  content = "This is a test index"
  source_id = "pwcitlawbs3xfhj3dspx2yh2cbi-wndwuafzh7ml7oiyacjka62lw4"
  document_id = "product://004-yellow2"
}
