client {
  enabled = true
}

server {
  enabled = true
  bootstrap_expect = 1
}

datacenter = "a"
data_dir = "/opt/nomad"
name = "wynyard"
