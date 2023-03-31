client {
  enabled = true
  network_interface = "lo"
}

server {
  enabled = true
  bootstrap_expect = 1
}

datacenter = "asia-southeast1-a"
data_dir = "/opt/nomad"
name = "wynyard-0"

plugin "raw_exec" {
  config {
    enabled = true
    no_cgroups = true
  }
}

plugin "docker" {
  config {
    volumes {
      enabled = true
    }
  }
}
