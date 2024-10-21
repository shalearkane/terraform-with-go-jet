terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = "messiers-virgo"
}

variable "service_name" {
  description = "The name of the Cloud Run service"
  default     = "my-cloud-run-service"
}

variable "POSTGRES_HOST" {
  type = string
}

variable "POSTGRES_PORT" {
  type = string
}

variable "POSTGRES_USER" {
  type = string
}

variable "POSTGRES_PASSWORD" {
  type = string
}

variable "POSTGRES_DATABASE" {
  type = string
}



resource "google_cloud_run_v2_service" "default" {
  name     = "soumiks-sql"
  location = "asia-south2"
  client   = "terraform"


  template {
    containers {
      image = "asia-south2-docker.pkg.dev/messiers-virgo/cloud-run-source-deploy/soumik-sql:v0.1"
      env {
        name  = "POSTGRES_HOST"
        value = var.POSTGRES_HOST
      }
      env {
        name  = "POSTGRES_PORT"
        value = var.POSTGRES_PORT
      }
      env {
        name  = "POSTGRES_USER"
        value = var.POSTGRES_USER
      }
      env {
        name  = "POSTGRES_PASSWORD"
        value = var.POSTGRES_PASSWORD
      }
      env {
        name  = "POSTGRES_DATABASE"
        value = var.POSTGRES_DATABASE
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
