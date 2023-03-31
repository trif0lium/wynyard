client {
  enabled = true
}

server {
  enabled = true
  bootstrap_expect = 1
}

datacenter = "asia-southeast1-a"
data_dir = "/opt/nomad"
name = "wynyard"
