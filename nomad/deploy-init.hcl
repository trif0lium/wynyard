job "hello" {
  datacenters = ["asia-southeast1-a"]
  type = "service"

  group "hello" {
    count = 1

    task "volume-init" {
      lifecycle {
        hook = "prestart"
        sidecar = false
      }
      driver = "exec"
      config {
        command = "go"
        args = ["run", "main.go", "volume", "create", "-size", "3000", "vol_xyz"]
      }
    }

    task "ubuntu" {
      driver = "docker"
      config {
        image = "ubuntu:latest"
      }
    }
  }
}
