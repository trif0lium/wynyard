job "hello" {
  datacenters = ["asia-southeast1-a"]
  type = "service"

  group "hello" {
    count = 1

    task "ubuntu" {
      driver = "docker"
      config {
        image = "ubuntu:latest"
      }
    }
  }
}
