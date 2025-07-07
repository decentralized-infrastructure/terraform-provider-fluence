# Provider Configuration
terraform {
  required_providers {
    fluence = {
      source = "hashicorp.com/0xthresh/fluence"
    }
  }
}

# Configure the Fluence provider
provider "fluence" {
  # API key can be set via environment variable FLUENCE_API_KEY
  # or specified here directly:
  # api_key = "your-api-key-here"
  
  # Optional: Override the default API host
  # host = "https://api.fluence.dev"
}
