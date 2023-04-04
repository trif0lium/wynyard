client {
  enabled = true
  network_interface = "lo"
  server_join {
    retry_join = ["wynyard-0.asia-southeast1-a.c.railway-infra-dev.internal:4648"]
  }
}

datacenter = "asia-southeast1-b"
data_dir = "/opt/nomad"
name = "wynyard-1"

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
