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
      driver = "raw_exec"
      config {
        command = "/root/wynyard/go/build/wynyard"
        args = ["volume", "create", "-size", "3000", "-snapshot-host", "wynyard-1", "-snapshot-location", "asia-southeast1-b", "vol_xyz"]
      }
    }

    task "ubuntu" {
      driver = "docker"
      config {
        image = "ubuntu:latest"
        command = "sleep"
        args = ["infinity"]
        mount {
          type = "volume"
          target = "/mnt/external"
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
