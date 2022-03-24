package dnacenter

import (
	"context"
	"errors"
	"reflect"
	"time"

	"log"

	dnacentersdkgo "github.com/cisco-en-programmability/dnacenter-go-sdk/v3/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBusinessSdaWirelessControllerCreate() *schema.Resource {
	return &schema.Resource{
		Description: `It performs create operation on Fabric Wireless.
- Add WLC to Fabric Domain
		Missing.
`,

		CreateContext: resourceBusinessSdaWirelessControllerCreateCreate,
		ReadContext:   resourceBusinessSdaWirelessControllerCreateRead,
		UpdateContext: resourceBusinessSdaWirelessControllerCreateUpdate,
		DeleteContext: resourceBusinessSdaWirelessControllerCreateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"item": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"execution_id": &schema.Schema{
							Description: `Status of the job for wireless state change in fabric domain
`,
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_status_url": &schema.Schema{
							Description: `executionStatusURL`,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"message": &schema.Schema{
							Description: `message`,
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": &schema.Schema{
							Description: `EWLC Device Name
			`,
							Type:     schema.TypeString,
							Optional: true,
						},
						"site_name_hierarchy": &schema.Schema{
							Description: `Site Name Hierarchy
			`,
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceBusinessSdaWirelessControllerCreateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dnacentersdkgo.Client)

	var diags diag.Diagnostics

	resourceItem := *getResourceItem(d.Get("parameters"))
	vDevice_name := resourceItem["device_name"]
	vvDevice_name := interfaceToString(vDevice_name)
	vSite_name_hierarchy := resourceItem["site_name_hierarchy"]
	vvSite_name_hierarchy := interfaceToString(vSite_name_hierarchy)

	request1 := expandRequestBusinessSdaWirelessControllerCreateCreateAddWLCToFabricDomain(ctx, "parameters.0", d)
	log.Printf("[DEBUG] request sent => %v", responseInterfaceToString(*request1))

	response1, restyResp1, err := client.FabricWireless.AddWLCToFabricDomain(request1)
	if err != nil || response1 == nil {
		if restyResp1 != nil {
			log.Printf("[DEBUG] Retrieved error response %s", restyResp1.String())
		}
		diags = append(diags, diagErrorWithAlt(
			"Failure when executing AddWLCToFabricDomain", err,
			"Failure at AddWLCToFabricDomain, unexpected response", ""))
		return diags
	}
	log.Printf("[DEBUG] Retrieved response %+v", responseInterfaceToString(*response1))

	vItems1 := flattenFabricWirelessAddWLCToFabricDomainItems(response1)
	if err := d.Set("item", vItems1); err != nil {
		diags = append(diags, diagError(
			"Failure when setting AddWLCToFabricDomain response",
			err))
		return diags
	}

	executionId := response1.ExecutionID
	log.Printf("[DEBUG] ExecutionID => %s", executionId)
	if executionId != "" {
		time.Sleep(5 * time.Second)
		response2, restyResp2, err := client.Task.GetBusinessAPIExecutionDetails(executionId)
		if err != nil || response2 == nil {
			if restyResp2 != nil {
				log.Printf("[DEBUG] Retrieved error response %s", restyResp2.String())
			}
			diags = append(diags, diagErrorWithAlt(
				"Failure when executing GetBusinessAPIExecutionDetails", err,
				"Failure at GetBusinessAPIExecutionDetails, unexpected response", ""))
			return diags
		}
		for response2.Status == "IN_PROGRESS" {
			time.Sleep(10 * time.Second)
			response2, restyResp1, err = client.Task.GetBusinessAPIExecutionDetails(executionId)
			if err != nil || response2 == nil {
				if restyResp1 != nil {
					log.Printf("[DEBUG] Retrieved error response %s", restyResp1.String())
				}
				diags = append(diags, diagErrorWithAlt(
					"Failure when executing GetExecutionByID", err,
					"Failure at GetExecutionByID, unexpected response", ""))
				return diags
			}
		}
		if response2.Status == "FAILURE" {
			bapiError := response2.BapiError
			diags = append(diags, diagErrorWithAlt(
				"Failure when executing AddWLCToFabricDomain", err,
				"Failure at AddWLCToFabricDomain execution", bapiError))
			return diags
		}
	}

	resourceMap := make(map[string]string)
	resourceMap["device_name"] = vvDevice_name
	resourceMap["site_name_hierarchy"] = vvSite_name_hierarchy
	d.SetId(joinResourceID(resourceMap))
	return resourceBusinessSdaWirelessControllerCreateRead(ctx, d, m)
}

func resourceBusinessSdaWirelessControllerCreateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := m.(*dnacentersdkgo.Client)
	var diags diag.Diagnostics
	return diags
}

func resourceBusinessSdaWirelessControllerCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceBusinessSdaWirelessControllerCreateRead(ctx, d, m)
}

func resourceBusinessSdaWirelessControllerCreateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*dnacentersdkgo.Client)

	var diags diag.Diagnostics

	resourceID := d.Id()
	resourceMap := separateResourceID(resourceID)
	vSiteID := resourceMap["site_id"]
	vDeviceFamilyIDentifier := resourceMap["device_family_identifier"]
	vDeviceRole := resourceMap["device_role"]
	vImageID := resourceMap["image_id"]

	selectedMethod := 1
	//var vvID string
	//var vvName string
	if selectedMethod == 1 {
		//vvID = vID
		getResp, _, err := client.SoftwareImageManagementSwim.GetGoldenTagStatusOfAnImage(vSiteID, vDeviceFamilyIDentifier, vDeviceRole, vImageID)
		if err != nil || getResp == nil {
			// Assume that element it is already gone
			return diags
		}
	}
	response1, restyResp1, err := client.SoftwareImageManagementSwim.RemoveGoldenTagForImage(vSiteID, vDeviceFamilyIDentifier, vDeviceRole, vImageID)
	if err != nil || response1 == nil {
		if restyResp1 != nil {
			log.Printf("[DEBUG] resty response for delete operation => %v", restyResp1.String())
			diags = append(diags, diagErrorWithAltAndResponse(
				"Failure when executing RemoveGoldenTagForImage", err, restyResp1.String(),
				"Failure at RemoveGoldenTagForImage, unexpected response", ""))
			return diags
		}
		diags = append(diags, diagErrorWithAlt(
			"Failure when executing RemoveGoldenTagForImage", err,
			"Failure at RemoveGoldenTagForImage, unexpected response", ""))
		return diags
	}
	taskId := response1.Response.TaskID
	log.Printf("[DEBUG] TASKID => %s", taskId)
	if taskId != "" {
		time.Sleep(5 * time.Second)
		response2, restyResp2, err := client.Task.GetTaskByID(taskId)
		if err != nil || response2 == nil {
			if restyResp2 != nil {
				log.Printf("[DEBUG] Retrieved error response %s", restyResp2.String())
			}
			diags = append(diags, diagErrorWithAlt(
				"Failure when executing GetTaskByID", err,
				"Failure at GetTaskByID, unexpected response", ""))
			return diags
		}
		if response2.Response != nil && response2.Response.IsError != nil && *response2.Response.IsError {
			log.Printf("[DEBUG] Error reason %s", response2.Response.FailureReason)
			restyResp3, err := client.CustomCall.GetCustomCall(response2.Response.AdditionalStatusURL, nil)
			if err != nil {
				diags = append(diags, diagErrorWithAlt(
					"Failure when executing GetCustomCall", err,
					"Failure at GetCustomCall, unexpected response", ""))
				return diags
			}
			var errorMsg string
			if restyResp3 == nil {
				errorMsg = response2.Response.Progress + "\nFailure Reason: " + response2.Response.FailureReason
			} else {
				errorMsg = restyResp3.String()
			}
			err1 := errors.New(errorMsg)
			diags = append(diags, diagError(
				"Failure when executing RemoveGoldenTagForImage", err1))
			return diags
		}
	}
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func expandRequestBusinessSdaWirelessControllerCreateCreateAddWLCToFabricDomain(ctx context.Context, key string, d *schema.ResourceData) *dnacentersdkgo.RequestFabricWirelessAddWLCToFabricDomain {
	request := dnacentersdkgo.RequestFabricWirelessAddWLCToFabricDomain{}
	if v, ok := d.GetOkExists(fixKeyAccess(key + ".device_name")); !isEmptyValue(reflect.ValueOf(d.Get(fixKeyAccess(key+".device_name")))) && (ok || !reflect.DeepEqual(v, d.Get(fixKeyAccess(key+".device_name")))) {
		request.DeviceName = interfaceToString(v)
	}
	if v, ok := d.GetOkExists(fixKeyAccess(key + ".site_name_hierarchy")); !isEmptyValue(reflect.ValueOf(d.Get(fixKeyAccess(key+".site_name_hierarchy")))) && (ok || !reflect.DeepEqual(v, d.Get(fixKeyAccess(key+".site_name_hierarchy")))) {
		request.SiteNameHierarchy = interfaceToString(v)
	}
	if isEmptyValue(reflect.ValueOf(request)) {
		return nil
	}

	return &request
}

func flattenFabricWirelessAddWLCToFabricDomainItems(items *dnacentersdkgo.ResponseFabricWirelessAddWLCToFabricDomain) []map[string]interface{} {
	if items == nil {
		return nil
	}
	var respItems []map[string]interface{}
	respItem := make(map[string]interface{})
	respItem["execution_id"] = items.ExecutionID
	respItem["execution_status_url"] = items.ExecutionStatusURL
	respItem["message"] = items.Message
	respItems = append(respItems, respItem)
	return respItems
}
