package gitlab

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectMirror() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabMirrorCreate,
		Read:   resourceGitlabMirrorRead,
		Update: resourceGitlabMirrorUpdate,
		Delete: resourceGitlabMirrorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"url": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"only_protected_branches": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"keep_divergent_refs": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceGitlabMirrorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	projectId := d.Get("project_id").(string)
	accessLevelId := accessLevelID[d.Get("access_level").(string)]

	options := &gitlab.AddProjectMemberOptions{
		UserID:      &userId,
		AccessLevel: &accessLevelId,
	}
	log.Printf("[DEBUG] create gitlab project membership for %d in %s", options.UserID, projectId)

	_, _, err := client.ProjectMembers.AddProjectMember(projectId, options)
	if err != nil {
		return err
	}
	userIdString := strconv.Itoa(userId)
	d.SetId(buildTwoPartID(&projectId, &userIdString))
	return resourceGitlabProjectMembershipRead(d, meta)
}
