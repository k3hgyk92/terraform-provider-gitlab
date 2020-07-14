package gitlab

import (
	"fmt"
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
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"mirror_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"only_protected_branches": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"keep_divergent_refs": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGitlabMirrorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	projectID := d.Get("project_id").(string)
	URL := d.Get("url").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := &gitlab.AddProjectMirrorOptions{
		URL:                   &URL,
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}

	log.Printf("[DEBUG] create gitlab project mirror for project %v", projectID)

	mirror, _, err := client.ProjectMirrors.AddProjectMirror(projectID, options)
	if err != nil {
		return err
	}
	d.Set("mirror_id", mirror.ID)

	mirrorID := strconv.Itoa(mirror.ID)
	d.SetId(buildTwoPartID(&projectID, &mirrorID))
	return resourceGitlabProjectMirrorRead(d, meta)
}

func resourceGitlabMirrorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project_id").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := gitlab.EditProjectMirrorOptions{
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}
	log.Printf("[DEBUG] update gitlab project mirror %v for %s", mirrorID, projectID)

	_, _, err := client.ProjectMirrors.EditProjectMirror(projectID, mirrorID, &options)
	if err != nil {
		return err
	}
	return resourceGitlabProjectMirrorRead(d, meta)
}

// Documented remote mirrors API does not support a delete method, instead mirror is disabled.
func resourceGitlabMirrorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	enabled := false

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project_id").(string)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := gitlab.EditProjectMirrorOptions{
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}
	log.Printf("[DEBUG] Disable gitlab project mirror %v for %s", mirrorID, projectID)

	_, _, err := client.ProjectMirrors.EditProjectMirror(projectID, mirrorID, &options)

	return err
}

func resourceGitlabProjectMirrorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project_id").(string)
	log.Printf("[DEBUG] read gitlab project mirror %s id %v", projectID, mirrorID)

	mirrors, _, err := client.ProjectMirrors.ListProjectMirror(projectID)

	if err != nil {
		return err
	}

	var mirror *gitlab.ProjectMirror
	found := false

	for _, m := range mirrors {
		log.Printf("[DEBUG] project mirror found %v", m.ID)
		if m.ID == mirrorID {
			mirror = m
			found = true
			return nil
		}
	}

	if !found {
		return fmt.Errorf("unable to find mirror %v on project %s", mirrorID, projectID)
	}

	resourceGitlabProjectMirrorSetToState(d, mirror, &projectID)
	return nil
}

func resourceGitlabProjectMirrorSetToState(d *schema.ResourceData, projectMirror *gitlab.ProjectMirror, projectID *string) {

	mirrorID := strconv.Itoa(projectMirror.ID)
	d.Set("enabled", projectMirror.Enabled)
	d.Set("only_protected_branches", projectMirror.OnlyProtectedBranches)
	d.Set("keep_divergent_refs", projectMirror.KeepDivergentRefs)
	d.SetId(buildTwoPartID(projectID, &mirrorID))
}
