---
subcategory: "microsoft.alertsManagement - Azure Monitor"
page_title: "smartDetectorAlertRules"
description: |-
  Manages a Monitor Smart Detector Alert Rule.
---

# microsoft.alertsManagement/smartDetectorAlertRules - Monitor Smart Detector Alert Rule

This article demonstrates how to use `azapi` provider to manage the Monitor Smart Detector Alert Rule resource in Azure.

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

resource "azapi_resource" "actionGroup" {
  type      = "Microsoft.Insights/actionGroups@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      armRoleReceivers = [
      ]
      automationRunbookReceivers = [
      ]
      azureAppPushReceivers = [
      ]
      azureFunctionReceivers = [
      ]
      emailReceivers = [
      ]
      enabled = true
      eventHubReceivers = [
      ]
      groupShortName = "acctestag"
      itsmReceivers = [
      ]
      logicAppReceivers = [
      ]
      smsReceivers = [
      ]
      voiceReceivers = [
      ]
      webhookReceivers = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "smartDetectorAlertRule" {
  type      = "microsoft.alertsManagement/smartDetectorAlertRules@2019-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      actionGroups = {
        customEmailSubject   = ""
        customWebhookPayload = ""
        groupIds = [
          azapi_resource.actionGroup.id,
        ]
      }
      description = ""
      detector = {
        id = "FailureAnomaliesDetector"
      }
      frequency = "PT1M"
      scope = [
        azapi_resource.component.id,
      ]
      severity = "Sev0"
      state    = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `microsoft.alertsManagement/smartDetectorAlertRules@api-version`. The available api-versions for this resource are: [`2019-03-01`, `2019-06-01`, `2021-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/microsoft.alertsManagement/smartDetectorAlertRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{resourceName}?api-version=2021-04-01
 ```
