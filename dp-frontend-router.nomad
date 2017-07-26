job "dp-frontend-router" {
  datacenters = ["eu-west-1"]
  region      = "eu"
  type        = "service"

  update {
    stagger      = "20s"
    max_parallel = 1
  }

  constraint {
    distinct_hosts = true
  }

  group "web" {
    count = "{{WEB_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "web"
    }

    task "dp-frontend-router-web" {
      driver = "docker"

      config {
        image = "s3::https://s3-eu-west-1.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-frontend-router/{{REVISION}}.tar.gz"

        port_map {
          http = 8080
        }
      }

      service {
        name = "dp-frontend-router"
        port = "http"
        tags = ["web"]
      }

      resources {
        cpu    = "{{WEB_RESOURCE_CPU}}"
        memory = "{{WEB_RESOURCE_MEM}}"

        network {
          port "http" {}
        }
      }
    }
  }

  group "publising" {
    count = "{{PUBLISHING_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "publishing"
    }

    task "dp-frontend-router-publishing" {
      driver = "docker"

      config {
        image = "s3::https://s3-eu-west-1.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-frontend-router/{{REVISION}}.tar.gz"

        port_map {
          http = 8080
        }
      }

      service {
        name = "dp-frontend-router"
        port = "http"
        tags = ["publishing"]
      }

      resources {
        cpu    = "{{PUBLISHING_RESOURCE_CPU}}"
        memory = "{{PUBLISHING_RESOURCE_CPU}}"

        network {
          port "http" {}
        }
      }
    }
  }
}
