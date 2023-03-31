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
        command = "/bin/bash"
        args = ["-c" "cd /root/wynyard/go && go run main.go volume create -size 3000 vol_xyz"]
      }
    }

    task "ubuntu" {
      driver = "docker"
      config {
        image = "ubuntu:latest"
        mount {
          type = "volume"
          target = "/mnt/external"
          readonly = false
          volume_options {
            driver_config {
              name = "local"
              options {
                device = "/dev/mapper/vg0-vol_xyz"
                type = "ext4"
              }
            }
          }
        }
      }
    }
  }
}