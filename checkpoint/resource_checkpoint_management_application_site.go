package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	
	
)

func resourceManagementApplicationSite() *schema.Resource {   
    return &schema.Resource{
        Create: createManagementApplicationSite,
        Read:   readManagementApplicationSite,
        Update: updateManagementApplicationSite,
        Delete: deleteManagementApplicationSite,
        Schema: map[string]*schema.Schema{ 
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Object name.",
            },
            "additional_categories": {
                Type:        schema.TypeSet,
                Optional:    true,
                Description: "Used to configure or edit the additional categories of a custom application / site used in the Application and URL Filtering or Threat Prevention.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "description": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "A description for the application.",
            },
            "primary_category": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Each application is assigned to one primary category based on its most defining aspect.",
            },
            "tags": {
                Type:        schema.TypeSet,
                Optional:    true,
                Description: "Collection of tag identifiers.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "url_list": {
                Type:        schema.TypeSet,
                Optional:    true,
                Description: "URLs that determine this particular application.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "application_signature": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Application signature generated by <a href=\"https://supportcenter.checkpoint.com/supportcenter/portal?eventSubmit_doGoviewsolutiondetails=&solutionid=sk103051\">Signature Tool</a>.",
            },
            "urls_defined_as_regular_expression": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: "States whether the URL is defined as a Regular Expression or not.",
                Default:     false,
            },
            "color": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Color of the object. Should be one of existing colors.",
                Default:     "black",
            },
            "comments": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Comments string.",
            },
            "groups": {
                Type:        schema.TypeSet,
                Optional:    true,
                Description: "Collection of group identifiers.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "ignore_warnings": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: "Apply changes ignoring warnings.",
                Default:     false,
            },
            "ignore_errors": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: "Apply changes ignoring errors. You won't be able to publish such a changes. If ignore-warnings flag was omitted - warnings will also be ignored.",
                Default:     false,
            },
        },
    }
}

func createManagementApplicationSite(d *schema.ResourceData, m interface{}) error {
    client := m.(*checkpoint.ApiClient)

    applicationSite := make(map[string]interface{})

    if v, ok := d.GetOk("name"); ok {
        applicationSite["name"] = v.(string)
    }

    if v, ok := d.GetOk("additional_categories"); ok {
        applicationSite["additional-categories"] = v.(*schema.Set).List()
    }

    if v, ok := d.GetOk("description"); ok {
        applicationSite["description"] = v.(string)
    }

    if v, ok := d.GetOk("primary_category"); ok {
        applicationSite["primary-category"] = v.(string)
    }

    if v, ok := d.GetOk("tags"); ok {
        applicationSite["tags"] = v.(*schema.Set).List()
    }

    if v, ok := d.GetOk("url_list"); ok {
        applicationSite["url-list"] = v.(*schema.Set).List()
    }

    if v, ok := d.GetOk("application_signature"); ok {
        applicationSite["application-signature"] = v.(string)
    }

    if v, ok := d.GetOkExists("urls_defined_as_regular_expression"); ok {
        applicationSite["urls-defined-as-regular-expression"] = v.(bool)
    }

    if v, ok := d.GetOk("color"); ok {
        applicationSite["color"] = v.(string)
    }

    if v, ok := d.GetOk("comments"); ok {
        applicationSite["comments"] = v.(string)
    }

    if v, ok := d.GetOk("groups"); ok {
        applicationSite["groups"] = v.(*schema.Set).List()
    }

    if v, ok := d.GetOkExists("ignore_warnings"); ok {
        applicationSite["ignore-warnings"] = v.(bool)
    }

    if v, ok := d.GetOkExists("ignore_errors"); ok {
        applicationSite["ignore-errors"] = v.(bool)
    }

    log.Println("Create ApplicationSite - Map = ", applicationSite)

    addApplicationSiteRes, err := client.ApiCall("add-application-site", applicationSite, client.GetSessionID(), true, false)
    if err != nil || !addApplicationSiteRes.Success {
        if addApplicationSiteRes.ErrorMsg != "" {
            return fmt.Errorf(addApplicationSiteRes.ErrorMsg)
        }
        return fmt.Errorf(err.Error())
    }

    d.SetId(addApplicationSiteRes.GetData()["uid"].(string))

    return readManagementApplicationSite(d, m)
}

func readManagementApplicationSite(d *schema.ResourceData, m interface{}) error {

    client := m.(*checkpoint.ApiClient)

    payload := map[string]interface{}{
        "uid": d.Id(),
    }

    showApplicationSiteRes, err := client.ApiCall("show-application-site", payload, client.GetSessionID(), true, false)
    if err != nil {
		return fmt.Errorf(err.Error())
	}
    if !showApplicationSiteRes.Success {
		if objectNotFound(showApplicationSiteRes.GetData()["code"].(string)) {
			d.SetId("")
			return nil
		}
        return fmt.Errorf(showApplicationSiteRes.ErrorMsg)
    }

    applicationSite := showApplicationSiteRes.GetData()

    log.Println("Read ApplicationSite - Show JSON = ", applicationSite)

	if v := applicationSite["name"]; v != nil {
		_ = d.Set("name", v)
	}

    if applicationSite["additional_categories"] != nil {
        additionalCategoriesJson, ok := applicationSite["additional_categories"].([]interface{})
        if ok {
            additionalCategoriesIds := make([]string, 0)
            if len(additionalCategoriesJson) > 0 {
                for _, additional_categories := range additionalCategoriesJson {
                    additional_categories := additional_categories.(map[string]interface{})
                    additionalCategoriesIds = append(additionalCategoriesIds, additional_categories["name"].(string))
                }
            }
        _ = d.Set("additional_categories", additionalCategoriesIds)
        }
    } else {
        _ = d.Set("additional_categories", nil)
    }

	if v := applicationSite["description"]; v != nil {
		_ = d.Set("description", v)
	}

	if v := applicationSite["primary-category"]; v != nil {
		_ = d.Set("primary_category", v)
	}

    if applicationSite["tags"] != nil {
        tagsJson, ok := applicationSite["tags"].([]interface{})
        if ok {
            tagsIds := make([]string, 0)
            if len(tagsJson) > 0 {
                for _, tags := range tagsJson {
                    tags := tags.(map[string]interface{})
                    tagsIds = append(tagsIds, tags["name"].(string))
                }
            }
        _ = d.Set("tags", tagsIds)
        }
    } else {
        _ = d.Set("tags", nil)
    }

    if applicationSite["url_list"] != nil {
        urlListJson, ok := applicationSite["url_list"].([]interface{})
        if ok {
            urlListIds := make([]string, 0)
            if len(urlListJson) > 0 {
                for _, url_list := range urlListJson {
                    url_list := url_list.(map[string]interface{})
                    urlListIds = append(urlListIds, url_list["name"].(string))
                }
            }
        _ = d.Set("url_list", urlListIds)
        }
    } else {
        _ = d.Set("url_list", nil)
    }

	if v := applicationSite["application-signature"]; v != nil {
		_ = d.Set("application_signature", v)
	}

	if v := applicationSite["urls-defined-as-regular-expression"]; v != nil {
		_ = d.Set("urls_defined_as_regular_expression", v)
	}

	if v := applicationSite["color"]; v != nil {
		_ = d.Set("color", v)
	}

	if v := applicationSite["comments"]; v != nil {
		_ = d.Set("comments", v)
	}

    if applicationSite["groups"] != nil {
        groupsJson, ok := applicationSite["groups"].([]interface{})
        if ok {
            groupsIds := make([]string, 0)
            if len(groupsJson) > 0 {
                for _, groups := range groupsJson {
                    groups := groups.(map[string]interface{})
                    groupsIds = append(groupsIds, groups["name"].(string))
                }
            }
        _ = d.Set("groups", groupsIds)
        }
    } else {
        _ = d.Set("groups", nil)
    }

	if v := applicationSite["ignore-warnings"]; v != nil {
		_ = d.Set("ignore_warnings", v)
	}

	if v := applicationSite["ignore-errors"]; v != nil {
		_ = d.Set("ignore_errors", v)
	}

	return nil

}

func updateManagementApplicationSite(d *schema.ResourceData, m interface{}) error {

	client := m.(*checkpoint.ApiClient)
    applicationSite := make(map[string]interface{})

    if ok := d.HasChange("name"); ok {
        oldName, newName := d.GetChange("name")
        applicationSite["name"] = oldName
        applicationSite["new-name"] = newName
    } else {
        applicationSite["name"] = d.Get("name")
    }

    if d.HasChange("additional_categories") {
        if v, ok := d.GetOk("additional_categories"); ok {
            applicationSite["additional_categories"] = v.(*schema.Set).List()
        } else {
            oldAdditional_Categories, _ := d.GetChange("additional_categories")
	           applicationSite["additional_categories"] = map[string]interface{}{"remove": oldAdditional_Categories.(*schema.Set).List()}
        }
    }

    if ok := d.HasChange("description"); ok {
	       applicationSite["description"] = d.Get("description")
    }

    if ok := d.HasChange("primary_category"); ok {
	       applicationSite["primary-category"] = d.Get("primary_category")
    }

    if d.HasChange("tags") {
        if v, ok := d.GetOk("tags"); ok {
            applicationSite["tags"] = v.(*schema.Set).List()
        } else {
            oldTags, _ := d.GetChange("tags")
	           applicationSite["tags"] = map[string]interface{}{"remove": oldTags.(*schema.Set).List()}
        }
    }

    if d.HasChange("url_list") {
        if v, ok := d.GetOk("url_list"); ok {
            applicationSite["url_list"] = v.(*schema.Set).List()
        } else {
            oldUrl_List, _ := d.GetChange("url_list")
	           applicationSite["url_list"] = map[string]interface{}{"remove": oldUrl_List.(*schema.Set).List()}
        }
    }

    if ok := d.HasChange("application_signature"); ok {
	       applicationSite["application-signature"] = d.Get("application_signature")
    }

    if v, ok := d.GetOkExists("urls_defined_as_regular_expression"); ok {
	       applicationSite["urls-defined-as-regular-expression"] = v.(bool)
    }

    if ok := d.HasChange("color"); ok {
	       applicationSite["color"] = d.Get("color")
    }

    if ok := d.HasChange("comments"); ok {
	       applicationSite["comments"] = d.Get("comments")
    }

    if d.HasChange("groups") {
        if v, ok := d.GetOk("groups"); ok {
            applicationSite["groups"] = v.(*schema.Set).List()
        } else {
            oldGroups, _ := d.GetChange("groups")
	           applicationSite["groups"] = map[string]interface{}{"remove": oldGroups.(*schema.Set).List()}
        }
    }

    if v, ok := d.GetOkExists("ignore_warnings"); ok {
	       applicationSite["ignore-warnings"] = v.(bool)
    }

    if v, ok := d.GetOkExists("ignore_errors"); ok {
	       applicationSite["ignore-errors"] = v.(bool)
    }

    log.Println("Update ApplicationSite - Map = ", applicationSite)

    updateApplicationSiteRes, err := client.ApiCall("set-application-site", applicationSite, client.GetSessionID(), true, false)
    if err != nil || !updateApplicationSiteRes.Success {
        if updateApplicationSiteRes.ErrorMsg != "" {
            return fmt.Errorf(updateApplicationSiteRes.ErrorMsg)
        }
        return fmt.Errorf(err.Error())
    }

    return readManagementApplicationSite(d, m)
}

func deleteManagementApplicationSite(d *schema.ResourceData, m interface{}) error {

    client := m.(*checkpoint.ApiClient)

    applicationSitePayload := map[string]interface{}{
        "uid": d.Id(),
    }

    log.Println("Delete ApplicationSite")

    deleteApplicationSiteRes, err := client.ApiCall("delete-application-site", applicationSitePayload , client.GetSessionID(), true, false)
    if err != nil || !deleteApplicationSiteRes.Success {
        if deleteApplicationSiteRes.ErrorMsg != "" {
            return fmt.Errorf(deleteApplicationSiteRes.ErrorMsg)
        }
        return fmt.Errorf(err.Error())
    }
    d.SetId("")

    return nil
}

