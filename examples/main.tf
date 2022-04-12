terraform {
  required_providers {
    minanode = {
        source  = "hashicorp.com/edu/minanode"
    }
  }
}

provider "minanode" {
}

resource "minanode_node" "node1" {
    name = "utibe2"
    privkey = "test"
    replicas = 1
    namespace = "default"
}

output "name" {
  value = minanode_node.node1.name
}