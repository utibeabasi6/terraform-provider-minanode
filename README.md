# Terraform Provider Minanode

This repository contains the code for the Minanode Terraform provider

## How to use

In your terraform code, declare the provider thus

```hcl

terraform {
  required_providers {
    minanode = {
      source = "utibeabasi6/minanode"
    }
  }
}

```

Create a provider block and pass in the path to your Kubeconfig file. Leave blank to use the default kubeconfig location of `~/.kube/config`

```hcl
provider "minanode" {
  # Configuration options
  kubeconfig = # Path to kubeconfig file
}
```

Create a `minanode_node` resource and pass in the name, private key, namespace and number of replicas for the node. By default, the namespace is set to `default` and the replicas is set to `1` if left empty. The name and privkey are required

```hcl
resource "minanode_node" "node" {
    name = "minaprotocol"
    privkey = "key"
}
```
Run `terraform apply` to create the resource and finally run `kubectl get pod -n <namespace>` to view the created resources.