client {
  enabled = true
  network_interface = "lo"
}

datacenter = "asia-southeast1-b"
data_dir = "/opt/nomad"
name = "wynyard"

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
