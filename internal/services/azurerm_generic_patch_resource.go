package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/validate"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/tf"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func ResourceAzureGenericPatchResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGenericPatchResourceCreateUpdate,
		Read:   resourceAzureGenericPatchResourceRead,
		Update: resourceAzureGenericPatchResourceCreateUpdate,
		Delete: resourceAzureGenericPatchResourceDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azurermg_patch_resource` doesn't support import")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AzureResourceID,
			},

			"api_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"body": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
			},

			"response_export_values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.HasChange("response_export_values") {
				d.SetNewComputed("output")
			}
			old, new := d.GetChange("body")
			if utils.NormalizeJson(old) != utils.NormalizeJson(new) {
				d.SetNewComputed("output")
			}
			return nil
		},
	}
}

func resourceAzureGenericPatchResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceID(d.Get("resource_id").(string), d.Get("api_version").(string))

	existing, _, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	if len(utils.GetId(existing)) == 0 {
		return fmt.Errorf("patch target does not exist %s", id)
	}

	var requestBody interface{}
	err = json.Unmarshal([]byte(d.Get("body").(string)), &requestBody)
	if err != nil {
		return err
	}

	requestBody = utils.GetMergedJson(existing, requestBody)
	requestBody = utils.GetIgnoredJson(requestBody, getUnsupportedProperties())

	j, _ := json.Marshal(requestBody)
	log.Printf("[INFO] request body: %v\n", string(j))
	_, _, err = client.CreateUpdate(ctx, id.AzureResourceId, id.ApiVersion, requestBody, http.MethodPut)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAzureGenericPatchResourceRead(d, meta)
}

func resourceAzureGenericPatchResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceID(d.Id())
	if err != nil {
		return err
	}

	responseBody, response, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Error reading %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("resource_id", id.AzureResourceId)
	d.Set("api_version", id.ApiVersion)

	bodyJson := d.Get("body").(string)
	var requestBody interface{}
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil {
		return err
	}
	data, err := json.Marshal(utils.GetUpdatedJson(requestBody, responseBody))
	if err != nil {
		return err
	}
	d.Set("body", string(data))

	paths := d.Get("response_export_values").([]interface{})
	var output interface{}
	if len(paths) != 0 {
		output = make(map[string]interface{})
		for _, path := range paths {
			part := utils.ExtractObject(responseBody, path.(string))
			if part == nil {
				continue
			}
			output = utils.GetMergedJson(output, part)
		}
	}
	if output == nil {
		output = make(map[string]interface{})
	}
	outputJson, _ := json.Marshal(output)
	d.Set("output", string(outputJson))
	return nil
}

func resourceAzureGenericPatchResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceID(d.Id())
	if err != nil {
		return err
	}

	existing, _, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	if len(utils.GetId(existing)) == 0 {
		return fmt.Errorf("patch target does not exist %s", id)
	}

	var requestBody interface{}
	err = json.Unmarshal([]byte(d.Get("body").(string)), &requestBody)
	if err != nil {
		return err
	}

	requestBody = utils.GetRemovedJson(existing, requestBody)
	requestBody = utils.GetIgnoredJson(requestBody, getUnsupportedProperties())
	j, _ := json.Marshal(requestBody)
	log.Printf("[INFO] request body: %v\n", string(j))
	_, _, err = client.CreateUpdate(ctx, id.AzureResourceId, id.ApiVersion, requestBody, http.MethodPut)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	return nil
}

func getUnsupportedProperties() []string {
	return []string{"provisioningState"}
}