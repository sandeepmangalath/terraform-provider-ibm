---
layout: "ibm"
page_title: "IBM: ibm_container_cluster"
sidebar_current: "docs-ibm-datasource-container-cluster"
description: |-
  Get information about a Kubernetes cluster on IBM Cloud.
---

# ibm\_container_cluster


Import the details of a Kubernetes cluster on IBM Cloud as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```hcl
data "ibm_container_cluster" "cluster_foo" {
  cluster_name_id = "FOO"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name_id` - (Required, string) The name or ID of the cluster.
* `alb_type` - (Optional, string) Used to filter the albs based on type. Valid values are `private`, `public` and `all`. The default value is `all`.
* `org_guid` - (Deprecated, string) The GUID for the IBM Cloud organization associated with the cluster. You can retrieve the value from the `ibm_org` data source or by running the `ibmcloud iam orgs --guid` command in the [IBM Cloud CLI](https://cloud.ibm.com/docs/cli?topic=cloud-cli-getting-started).
* `space_guid` - (Deprecated, string) The GUID for the IBM Cloud space associated with the cluster. You can retrieve the value from the `ibm_space` data source or by running the `ibmcloud iam space <space-name> --guid` command in the IBM Cloud CLI.
* `account_guid` - (Deprecated, string) The GUID for the IBM Cloud account associated with the cluster. You can retrieve the value from the `ibm_account` data source or by running the `ibmcloud iam accounts` command in the IBM Cloud CLI.
* `region` - (Deprecated, string) The region where the cluster is provisioned. If the region is not specified it will be defaulted to provider region(IC_REGION/IBMCLOUD_REGION). To get the list of supported regions please access this [link](https://containers.bluemix.net/v1/regions) and use the alias.
* `resource_group_id` - (Optional, string) The ID of the resource group.  You can retrieve the value from data source `ibm_resource_group`. If not provided defaults to default resource group.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the cluster.
* `worker_count` - The number of workers that are attached to the cluster.
* `workers` - The IDs of the workers that are attached to the cluster.
* `ingress_hostname` - The Ingress hostname.
* `ingress_secret` - The Ingress secret.
* `bounded_services` - The services that are bounded to the cluster.
* `is_trusted` - Is trusted cluster feature enabled.
* `vlans` - The VLAN'S that are attached to the cluster. Nested `vlans` blocks have the following structure:
	* `id` - The VLAN id.
	* `subnets` - The list of subnets attached to VLAN belonging to the cluster. Nested `subnets` blocks have the following structure:
		* `id` - The Subnet Id.
		* `cidr` - The cidr range.
		* `ips` - The list of ip's in the subnet.
		* `is_byoip` - `true` if user provides a ip range else `false`.
		* `is_public` - `true` if VLAN is public else `false`.
* `worker_pools` - Worker pools attached to the cluster
  * `name` - The name of the worker pool.
  * `machine_type` - The machine type of the worker node.
  * `size_per_zone` - Number of workers per zone in this pool.
  * `hardware` - The level of hardware isolation for your worker node.
  * `id` - Worker pool id.
  * `state` - Worker pool state.
  * `labels` - Labels on all the workers in the worker pool.
	* `zones` - List of zones attached to the worker_pool.
		* `zone` - Zone name.
		* `private_vlan` - The ID of the private VLAN. 
		* `public_vlan` - The ID of the public VLAN.
		* `worker_count` - Number of workers attached to this zone.
* `albs` - Application load balancer (ALB)'s attached to the cluster
  * `id` - The application load balancer (ALB) id.
  * `name` - The name of the application load balancer (ALB).
  * `alb_type` - The application load balancer (ALB) type public or private.
  * `enable` -  Enable (true) or disable(false) application load balancer (ALB).
  * `state` - The status of the application load balancer (ALB)(enabled or disabled).
  * `num_of_instances` - Desired number of application load balancer (ALB) replicas.
  * `alb_ip` - BYOIP VIP to use for application load balancer (ALB). Currently supported only for private application load balancer (ALB).
  * `resize` - Indicate whether resizing should be done.
  * `disable_deployment` - Indicate whether to disable deployment only on disable application load balancer (ALB).
* `private_service_endpoint_url` - Private service endpoint url.
* `public_service_endpoint_url` - Public service endpoint url.
* `public_service_endpoint` - Is public service endpoint enabled to make the master publicly accessible.
* `private_service_endpoint` - Is private service endpoint enabled to make the master privately accessible.
* `crn` - CRN of the instance.
