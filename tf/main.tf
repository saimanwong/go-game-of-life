variable "gcp_project" {
  type      = string
  sensitive = true
}

variable "gcp_image" {
  type      = string
  sensitive = true
}

variable "gcp_max_scale" {
  type = string
}

variable "gcp_location" {
  type    = string
  default = "europe-north1"
}

variable "gcp_service_name" {
  type    = string
  default = "go-game-of-life"
}

terraform {
  backend "gcs" {
    bucket  = "tf-state-go-game-of-life"
  }
}

provider "google" {
  project = var.gcp_project
}

resource "google_cloud_run_service" "default" {
  name     = var.gcp_service_name
  location = var.gcp_location

  template {
    spec {
      containers {
        image = "gcr.io/${var.gcp_image}"
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = var.gcp_max_scale
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name
  policy_data = data.google_iam_policy.noauth.policy_data
}
