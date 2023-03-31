client {
  enabled = true
}

datacenter = "asia-southeast1-b"
data_dir = "/opt/nomad"
name = "wynyard"

plugin "raw_exec" {
  config {
    enabled = true
  }
}
