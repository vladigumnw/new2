terraform {
    required_providers {
        docker = {
            source = "kreuzwerker/docker"
            version = "~> 3.0"
        }
    }
}

provider "docker" {}

resource "docker_image" "go-simple-project" {
    name = "my-app"
    build {
        context = "${path.module}"
    }
}

resource "docker_container" "go-simple-project" {
    name = "my-app"
    image = docker_image.go-simple-project.image_id
    ports {
        internal = 8080
        external = 8080
    }
}