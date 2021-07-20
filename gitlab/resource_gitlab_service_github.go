package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceGithub() *schema.Resource {
	return &schema.Resource{
		Description: "**NOTE**: requires either EE (self-hosted) or Silver and above (GitLab.com).\n\n" +
			"This resource manages a [GitHub integration](https://docs.gitlab.com/ee/user/project/integrations/github.html) that updates pipeline statuses on a GitHub repo's pull requests.",

		Create: resourceGitlabServiceGithubCreate,
		Read:   resourceGitlabServiceGithubRead,
		Update: resourceGitlabServiceGithubUpdate,
		Delete: resourceGitlabServiceGithubDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabServiceGithubImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "A GitHub personal access token with at least `repo:status` scope.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"repository_url": {
				Description: "The URL of the GitHub repo to integrate with, e,g, https://github.com/gitlabhq/terraform-provider-gitlab.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"static_context": {
				Description: "Append instance name instead of branch to the status. Must enable to set a GitLab status check as _required_ in GitHub. See [Static / dynamic status check names] to learn more.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			// Computed from the GitLab API. Omitted event fields because they're always true in Github.
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGitlabServiceGithubSetToState(d *schema.ResourceData, service *gitlab.GithubService) {
	d.SetId(fmt.Sprintf("%d", service.ID))
	d.Set("repository_url", service.Properties.RepositoryURL)
	d.Set("static_context", service.Properties.StaticContext)

	d.Set("title", service.Title)
	d.Set("created_at", service.CreatedAt.String())
	d.Set("updated_at", service.UpdatedAt.String())
	d.Set("active", service.Active)
}

func resourceGitlabServiceGithubCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] create gitlab github service for project %s", project)

	opts := &gitlab.SetGithubServiceOptions{
		Token:         gitlab.String(d.Get("token").(string)),
		RepositoryURL: gitlab.String(d.Get("repository_url").(string)),
		StaticContext: gitlab.Bool(d.Get("static_context").(bool)),
	}

	_, err := client.Services.SetGithubService(project, opts)
	if err != nil {
		return err
	}

	return resourceGitlabServiceGithubRead(d, meta)
}

func resourceGitlabServiceGithubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] read gitlab github service for project %s", project)

	service, _, err := client.Services.GetGithubService(project)
	if err != nil {
		return err
	}

	resourceGitlabServiceGithubSetToState(d, service)

	return nil
}

func resourceGitlabServiceGithubUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServiceGithubCreate(d, meta)
}

func resourceGitlabServiceGithubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] delete gitlab github service for project %s", project)

	_, err := client.Services.DeleteGithubService(project)
	return err
}

func resourceGitlabServiceGithubImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())

	return []*schema.ResourceData{d}, nil
}
