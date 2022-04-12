terraform {
  required_providers {
    minanode = {
      source = "utibeabasi6/minanode"
      version = "0.0.7"
    }
  }
}

provider "minanode" {
  # Configuration options
  kubeconfig = "minanode-kubeconfig.yaml"
}

resource "minanode_node" "node1" {
    name = "utibe2"
    privkey = "test"
}

output "name" {
  value = minanode_node.node1.name
}