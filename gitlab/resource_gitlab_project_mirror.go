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
		Read:   resourceGitlabProjectMirrorRead,
		Update: resourceGitlabMirrorUpdate,
		Delete: resourceGitlabMirrorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"mirror_id": {
				Type:	schema.TypeInt,
				Computed: true,
			}
			"url": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"only_protected_branches": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"keep_divergent_refs": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
		},
	}
}

func resourceGitlabMirrorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	projectID := d.Get("project").(string)
	URL := d.Get("url").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches")
	keepDivergentRefs := d.Get("keep_divergent_refs")

	options := &gitlab.AddProjectMirrorOptions{
		URL: &URL,
		Enabled: &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs: &keepDivergentRefs
	}

	log.Printf("[DEBUG] create gitlab project mirror for project %d", projectId)

	mirror, _, err := client.ProjectMirror.AddProjectMirror(projectId, options)
	if err != nil {
		return err
	}
	d.Set("mirror_id", mirror.ID)
	d.SetId(buildTwoPartID(&projectId, mirror.ID))
	return resourceGitlabProjectMembershipRead(d, meta)
}



func resourceGitlabProjectMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)
	URL := d.Get("url").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches")
	keepDivergentRefs := d.Get("keep_divergent_refs")

	options := gitlab.EditProjectMirrorOptions{
		URL: &URL,
		Enabled: &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs: &keepDivergentRefs
	}
	log.Printf("[DEBUG] update gitlab project mirror %v for %s", userId, projectId)

	_, _, err := client.ProjectMirror.EditProjectMirror(projectID, mirroID, &options)
	if err != nil {
		return err
	}
	return resourceGitlabProjectMembershipRead(d, meta)
}

// Documented remote mirrors API does not support a delete method, instead mirror is disabled. 
func resourceGitlabProjectMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)


	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)
	URL := d.Get("url").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches")
	keepDivergentRefs := d.Get("keep_divergent_refs")

	options := gitlab.EditProjectMirrorOptions{
		URL: &URL,
		Enabled: false,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs: &keepDivergentRefs
	}
	log.Printf("[DEBUG] Disable gitlab project mirror %v for %s", userId, projectId)

	_, _, err := client.ProjectMirror.EditProjectMirror(projectID, mirroID, &options)

	return err
}

func resourceGitlabProjectMirrorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)
	projectID := d.Get("URL").(string)
	log.Printf("[DEBUG] read gitlab project mirror %s id %v", projectID, mirrorID)

	mirrors := gitlab.projectMirror.ListProjectMirror(projectID)

	var mirror *ProjectMirror

	for i, m := range mirrors {
		if m.ID == mirrorID {
			mirror = m
		}
		else {
			set("")
			return nil
		}
	}

	resourceGitlabProjectMirrorSetToState(d, mirror, &projectId)
	return nil
}

func resourceGitlabProjectMirrorSetToState(d *schema.ResourceData, projectMirror *gitlab.ProjectMirror, projectId *string) {
	d.Set("enabled", projectMirror.Enabled)
	d.Set("only_protected_branches", protectMirror.OnlyProtectedBranches)
	d.Set("keep_divergent_refs", projectMirror.KeepDivergentRefs)
	d.SetId(buildTwoPartID(&projectId, projectMirror.ID))
}
