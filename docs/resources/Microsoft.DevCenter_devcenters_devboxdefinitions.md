---
subcategory: "Microsoft.DevCenter - Microsoft Dev Box"
page_title: "devcenters/devboxdefinitions"
description: |-
  Manages a Dev Center Dev Box Definition.
---

# Microsoft.DevCenter/devcenters/devboxdefinitions - Dev Center Dev Box Definition

This article demonstrates how to use `azapi` provider to manage the Dev Center Dev Box Definition resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

provider "azurerm" {
  features {
  }
}

data "azurerm_client_config" "current" {
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devcenters@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "SystemAssigned"
      userAssignedIdentities = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "devBoxDefinition" {
  type      = "Microsoft.DevCenter/devcenters/devboxdefinitions@2024-10-01-preview"
  parent_id = azapi_resource.devCenter.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hibernateSupport = "Enabled"
      imageReference = {
        id = "${azapi_resource.devCenter.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
      },
      sku = {
        name = "general_i_8c32gb256ssd_v2"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DevCenter/devcenters/devboxdefinitions@api-version`. The available api-versions for this resource are: [`2022-08-01-preview`, `2022-09-01-preview`, `2022-10-12-preview`, `2022-11-11-preview`, `2023-01-01-preview`, `2023-04-01`, `2023-08-01-preview`, `2023-10-01-preview`, `2024-02-01`, `2024-05-01-preview`, `2024-06-01-preview`, `2024-07-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`, `2025-02-01`, `2025-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DevCenter/devcenters/devboxdefinitions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{resourceName}/devboxdefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{resourceName}/devboxdefinitions/{resourceName}?api-version=2025-04-01-preview
 ```
