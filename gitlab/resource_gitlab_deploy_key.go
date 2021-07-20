package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabDeployKey() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage [deploy keys](https://docs.gitlab.com/ee/user/project/deploy_keys/) for your GitLab projects.",

		Create: resourceGitlabDeployKeyCreate,
		Read:   resourceGitlabDeployKeyRead,
		Delete: resourceGitlabDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabDeployKeyStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the deploy key to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"title": {
				Description: "A title to describe the deploy key with.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Description: "The public ssh key body.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == strings.TrimSpace(new)
				},
			},
			"can_push": {
				Description: "Allow this deploy key to be used to push changes to the project.  Defaults to `false`. **NOTE::** this cannot currently be managed.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
		},
	}
}

func resourceGitlabDeployKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddDeployKeyOptions{
		Title:   gitlab.String(d.Get("title").(string)),
		Key:     gitlab.String(strings.TrimSpace(d.Get("key").(string))),
		CanPush: gitlab.Bool(d.Get("can_push").(bool)),
	}

	log.Printf("[DEBUG] create gitlab deployment key %s", *options.Title)

	deployKey, _, err := client.DeployKeys.AddDeployKey(project, options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", deployKey.ID))

	return resourceGitlabDeployKeyRead(d, meta)
}

func resourceGitlabDeployKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] read gitlab deploy key %s/%d", project, deployKeyID)

	deployKey, _, err := client.DeployKeys.GetDeployKey(project, deployKeyID)
	if err != nil {
		return err
	}

	d.Set("title", deployKey.Title)
	d.Set("key", deployKey.Key)
	d.Set("can_push", deployKey.CanPush)
	return nil
}

func resourceGitlabDeployKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete gitlab deploy key %s", d.Id())

	_, err = client.DeployKeys.DeleteDeployKey(project, deployKeyID)

	return err
}

func resourceGitlabDeployKeyStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Deploy Key import format; expected '{project_id}:{deploy_key_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
