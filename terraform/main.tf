terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  version = "3.5.0"

  credentials = file("./terraform-sa.json")

  project = var.project
  region = var.region
  zone = var.zone
}

resource "google_storage_default_object_access_control" "public_rule" {
  bucket = google_storage_bucket.bucket.name
  role   = "READER"
  entity = "allUsers"
}

resource "google_storage_bucket" "bucket" {
  name = "temp-read-large-file-bucket"
  location = "ASIA-SOUTHEAST1"
  force_destroy = true
}

resource "google_storage_bucket_object" "text_file" {
  name   = "big10.txt"
  source = "../assets/big10.txt"
  bucket = google_storage_bucket.bucket.name
}

resource "google_container_cluster" "primary" {
  name               = "playground-cluster"
  location           = var.zone
  initial_node_count = 1

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }

  cluster_autoscaling {
    enabled = false
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    metadata = {
      disable-legacy-endpoints = "true"
    }

  }
}


