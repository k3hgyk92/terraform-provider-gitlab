package gitlab

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServicePipelinesEmail() *schema.Resource {
	return &schema.Resource{
		Description: "This resource manages a [Pipelines email integration](https://docs.gitlab.com/ee/user/project/integrations/overview.html#integrations-listing) that emails the pipeline status to a list of recipients.",

		Create: resourceGitlabServicePipelinesEmailCreate,
		Read:   resourceGitlabServicePipelinesEmailRead,
		Update: resourceGitlabServicePipelinesEmailCreate,
		Delete: resourceGitlabServicePipelinesEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"recipients": {
				Description: ") email addresses where notifications are sent.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"notify_only_broken_pipelines": {
				Description: "Notify only broken pipelines. Default is true.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"branches_to_be_notified": {
				Description:  "Branches to send notifications for. Valid options are `all`, `default`, `protected`, and `default_and_protected`. Default is `default`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "default", "protected", "default_and_protected"}, true),
				Default:      "default",
			},
		},
	}
}

func resourceGitlabServicePipelinesEmailSetToState(d *schema.ResourceData, service *gitlab.PipelinesEmailService) {
	d.Set("recipients", strings.Split(service.Properties.Recipients, ",")) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue
	d.Set("notify_only_broken_pipelines", service.Properties.NotifyOnlyBrokenPipelines)
	d.Set("branches_to_be_notified", service.Properties.BranchesToBeNotified)
}

func resourceGitlabServicePipelinesEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)
	options := &gitlab.SetPipelinesEmailServiceOptions{
		Recipients:                gitlab.String(strings.Join(*stringSetToStringSlice(d.Get("recipients").(*schema.Set)), ",")),
		NotifyOnlyBrokenPipelines: gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool)),
		BranchesToBeNotified:      gitlab.String(d.Get("branches_to_be_notified").(string)),
	}

	log.Printf("[DEBUG] create gitlab pipelines emails service for project %s", project)

	_, err := client.Services.SetPipelinesEmailService(project, options)
	if err != nil {
		return err
	}

	return resourceGitlabServicePipelinesEmailRead(d, meta)
}

func resourceGitlabServicePipelinesEmailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] read gitlab pipelines emails service for project %s", project)

	service, _, err := client.Services.GetPipelinesEmailService(project)
	if err != nil {
		return err
	}

	d.Set("project", project)
	resourceGitlabServicePipelinesEmailSetToState(d, service)
	return nil
}

func resourceGitlabServicePipelinesEmailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab pipelines email service for project %s", project)

	_, err := client.Services.DeletePipelinesEmailService(project)
	return err
}
