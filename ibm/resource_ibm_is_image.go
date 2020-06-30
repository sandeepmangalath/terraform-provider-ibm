package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcclassicv1"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

const (
	isImageHref                   = "href"
	isImageName                   = "name"
	isImageTags                   = "tags"
	isImageOperatingSystem        = "operating_system"
	isImageStatus                 = "status"
	isImageVisibility             = "visibility"
	isImageFile                   = "file"
	isImageMinimumProvisionedSize = "size"
	isImageFormat                 = "format"
	isImageArchitecure            = "architecture"
	isImageResourceGroup          = "resource_group"

	isImageProvisioning     = "provisioning"
	isImageProvisioningDone = "done"
	isImageDeleting         = "deleting"
	isImageDeleted          = "done"
)

func resourceIBMISImage() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISImageCreate,
		Read:     resourceIBMISImageRead,
		Update:   resourceIBMISImageUpdate,
		Delete:   resourceIBMISImageDelete,
		Exists:   resourceIBMISImageExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isImageHref: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Image Href value",
			},

			isImageName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateISName,
				Description:  "Image name",
			},

			isImageTags: {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              resourceIBMVPCHash,
				DiffSuppressFunc: applyOnce,
				Description:      "Tags for the image",
			},

			isImageOperatingSystem: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image Operating system",
			},

			isImageStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isImageArchitecure: {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field is removed",
			},

			isImageMinimumProvisionedSize: {
				Type:     schema.TypeInt,
				Computed: true,
			},

			isImageVisibility: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isImageFile: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isImageFormat: {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field is removed",
			},

			isImageResourceGroup: {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMISImageCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Image create")
	href := d.Get(isImageHref).(string)
	name := d.Get(isImageName).(string)
	operatingSystem := d.Get(isImageOperatingSystem).(string)

	if userDetails.generation == 1 {
		err := classicImgCreate(d, meta, href, name, operatingSystem)
		if err != nil {
			return err
		}
	} else {
		err := imgCreate(d, meta, href, name, operatingSystem)
		if err != nil {
			return err
		}
	}
	return resourceIBMISImageRead(d, meta)
}

func classicImgCreate(d *schema.ResourceData, meta interface{}, href, name, operatingSystem string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	imagePrototype := &vpcclassicv1.ImagePrototype{
		Name: &name,
		File: &vpcclassicv1.ImageFilePrototype{
			Href: &href,
		},
		OperatingSystem: &vpcclassicv1.OperatingSystemIdentity{
			Name: &operatingSystem,
		},
	}
	if rgrp, ok := d.GetOk(isImageResourceGroup); ok {
		rg := rgrp.(string)
		imagePrototype.ResourceGroup = &vpcclassicv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	options := &vpcclassicv1.CreateImageOptions{
		ImagePrototype: imagePrototype,
	}

	image, response, err := sess.CreateImage(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Image creation err %s\n%s", err, response)
	}
	d.SetId(*image.ID)
	log.Printf("[INFO] Floating IP : %s", *image.ID)
	_, err = isWaitForClassicImageAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isImageTags); ok || v != "" {
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc image (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func imgCreate(d *schema.ResourceData, meta interface{}, href, name, operatingSystem string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	imagePrototype := &vpcv1.ImagePrototype{
		Name: &name,
		File: &vpcv1.ImageFilePrototype{
			Href: &href,
		},
		OperatingSystem: &vpcv1.OperatingSystemIdentity{
			Name: &operatingSystem,
		},
	}
	if rgrp, ok := d.GetOk(isImageResourceGroup); ok {
		rg := rgrp.(string)
		imagePrototype.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	options := &vpcv1.CreateImageOptions{
		ImagePrototype: imagePrototype,
	}

	image, response, err := sess.CreateImage(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Image creation err %s\n%s", err, response)
	}
	d.SetId(*image.ID)
	log.Printf("[INFO] Floating IP : %s", *image.ID)
	_, err = isWaitForImageAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isImageTags); ok || v != "" {
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc image (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func isWaitForClassicImageAvailable(imageC *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageProvisioning},
		Target:     []string{isImageProvisioningDone, ""},
		Refresh:    isClassicImageRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicImageRefreshFunc(imageC *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getimgoptions := &vpcclassicv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}

		if *image.Status == "available" || *image.Status == "failed" {
			return image, isImageProvisioningDone, nil
		}

		return image, isImageProvisioning, nil
	}
}
func isWaitForImageAvailable(imageC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageProvisioning},
		Target:     []string{isImageProvisioningDone, ""},
		Refresh:    isImageRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isImageRefreshFunc(imageC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getimgoptions := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}

		if *image.Status == "available" || *image.Status == "failed" {
			return image, isImageProvisioningDone, nil
		}

		return image, isImageProvisioning, nil
	}
}

func resourceIBMISImageUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	id := d.Id()
	name := ""
	hasChanged := false

	if d.HasChange(isImageName) {
		name = d.Get(isImageName).(string)
		hasChanged = true
	}
	if userDetails.generation == 1 {
		err := classicImgUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	} else {
		err := imgUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	}
	return resourceIBMISImageRead(d, meta)
}

func classicImgUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isImageTags) {
		options := &vpcclassicv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := sess.GetImage(options)
		if err != nil {
			return fmt.Errorf("Error getting Image IP: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Image (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		options := &vpcclassicv1.UpdateImageOptions{
			ID:   &id,
			Name: &name,
		}
		_, response, err := sess.UpdateImage(options)
		if err != nil {
			return fmt.Errorf("Error on update of resource vpc Image: %s\n%s", err, response)
		}
	}
	return nil
}

func imgUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isImageTags) {
		options := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := sess.GetImage(options)
		if err != nil {
			return fmt.Errorf("Error getting Image IP: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Image (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		options := &vpcv1.UpdateImageOptions{
			ID:   &id,
			Name: &name,
		}
		_, response, err := sess.UpdateImage(options)
		if err != nil {
			return fmt.Errorf("Error on update of resource vpc Image: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISImageRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	id := d.Id()
	if userDetails.generation == 1 {
		err := classicImgGet(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := imgGet(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicImgGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcclassicv1.GetImageOptions{
		ID: &id,
	}
	image, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}
	d.Set("id", *image.ID)
	// d.Set(isImageArchitecure, image.Architecture)
	d.Set(isImageMinimumProvisionedSize, *image.MinimumProvisionedSize)
	d.Set(isImageName, *image.Name)
	d.Set(isImageOperatingSystem, *image.OperatingSystem)
	// d.Set(isImageFormat, image.Format)
	d.Set(isImageFile, *image.File)
	d.Set(isImageHref, *image.Href)
	d.Set(isImageStatus, *image.Status)
	d.Set(isImageVisibility, *image.Visibility)
	tags, err := GetTagsUsingCRN(meta, *image.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Image (%s) tags: %s", d.Id(), err)
	}
	d.Set(isImageTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/compute/image")
	d.Set(ResourceName, *image.Name)
	d.Set(ResourceStatus, *image.Status)
	d.Set(ResourceCRN, *image.Crn)
	if image.ResourceGroup != nil {
		d.Set(isImageResourceGroup, *image.ResourceGroup.ID)
		rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
		if err != nil {
			return err
		}
		grp, err := rsMangClient.ResourceGroup().Get(*image.ResourceGroup.ID)
		if err != nil {
			return err
		}
		d.Set(ResourceGroupName, grp.Name)
	}
	return nil
}

func imgGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetImageOptions{
		ID: &id,
	}
	image, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}
	d.Set("id", *image.ID)
	// d.Set(isImageArchitecure, image.Architecture)
	d.Set(isImageMinimumProvisionedSize, *image.MinimumProvisionedSize)
	d.Set(isImageName, *image.Name)
	d.Set(isImageOperatingSystem, *image.OperatingSystem)
	// d.Set(isImageFormat, image.Format)
	d.Set(isImageFile, *image.File)
	d.Set(isImageHref, *image.Href)
	d.Set(isImageStatus, *image.Status)
	d.Set(isImageVisibility, *image.Visibility)
	tags, err := GetTagsUsingCRN(meta, *image.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Image (%s) tags: %s", d.Id(), err)
	}
	d.Set(isImageTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/compute/image")
	d.Set(ResourceName, *image.Name)
	d.Set(ResourceStatus, *image.Status)
	d.Set(ResourceCRN, *image.Crn)
	if image.ResourceGroup != nil {
		d.Set(isImageResourceGroup, *image.ResourceGroup.ID)
	}
	return nil
}

func resourceIBMISImageDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicImgDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := imgDelete(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicImgDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getImageOptions := &vpcclassicv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(getImageOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}

	options := &vpcclassicv1.DeleteImageOptions{
		ID: &id,
	}
	response, err = sess.DeleteImage(options)
	if err != nil {
		return fmt.Errorf("Error Deleting Image : %s\n%s", err, response)
	}
	_, err = isWaitForClassicImageDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func imgDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getImageOptions := &vpcv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(getImageOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}

	options := &vpcv1.DeleteImageOptions{
		ID: &id,
	}
	response, err = sess.DeleteImage(options)
	if err != nil {
		return fmt.Errorf("Error Deleting Image : %s\n%s", err, response)
	}
	_, err = isWaitForImageDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForClassicImageDeleted(imageC *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageDeleting},
		Target:     []string{"", isImageDeleted},
		Refresh:    isClassicImageDeleteRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicImageDeleteRefreshFunc(imageC *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getimgoptions := &vpcclassicv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return image, isImageDeleted, nil
			}
			return image, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}
		return image, isImageDeleting, err
	}
}
func isWaitForImageDeleted(imageC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageDeleting},
		Target:     []string{"", isImageDeleted},
		Refresh:    isImageDeleteRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isImageDeleteRefreshFunc(imageC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getimgoptions := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return image, isImageDeleted, nil
			}
			return image, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}
		return image, isImageDeleting, err
	}
}
func resourceIBMISImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()

	if userDetails.generation == 1 {
		exists, err := classicImgExists(d, meta, id)
		return exists, err
	} else {
		exists, err := imgExists(d, meta, id)
		return exists, err
	}
}

func classicImgExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcclassicv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Image: %s\n%s", err, response)
	}
	return true, nil
}

func imgExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Image: %s\n%s", err, response)
	}
	return true, nil
}
