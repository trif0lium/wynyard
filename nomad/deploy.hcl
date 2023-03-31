job "hello" {
  datacenters = ["asia-southeast-1"]
  type = "service"

  group "xyz" {
    count = 1

    task "ubuntu" {
      driver = "docker"
      config {
        image = "ubuntu:latest"
      }
    }
  }
}
