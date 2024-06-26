package services_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ListDataSource struct{}

func TestAccListDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output").IsJson(),
			),
		},
	})
}

func TestAccListDataSource_hclOutput(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.hclOutput(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.value.#").Exists(),
			),
		},
	})
}

func TestAccListDataSource_paging(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.paging(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ListDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Resources/resourceGroups@2024-03-01"
  parent_id              = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}
`
}

func (r ListDataSource) hclOutput() string {
	return `
provider "azurerm" {
  features {}
}

provider "azapi" {
  enable_hcl_output_for_data_source = true
}

data "azurerm_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Resources/resourceGroups@2024-03-01"
  parent_id              = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}
`
}

func (r ListDataSource) paging() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Authorization/policyDefinitions@2021-06-01"
  parent_id              = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}
`
}
