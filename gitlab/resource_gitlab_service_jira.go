package gitlab

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceJira() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to manage Jira integration.",

		Create: resourceGitlabServiceJiraCreate,
		Read:   resourceGitlabServiceJiraRead,
		Update: resourceGitlabServiceJiraUpdate,
		Delete: resourceGitlabServiceJiraDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabServiceJiraImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
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
			"url": {
				Description:  "The URL to the JIRA project which is being linked to this GitLab project. For example, https://jira.example.com.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateURLFunc,
			},
			"project_key": {
				Description: "The short identifier for your JIRA project, all uppercase, e.g., PROJ.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"username": {
				Description: "The username of the user created to be used with GitLab/JIRA.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password of the user created to be used with GitLab/JIRA.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"jira_issue_transition_id": {
				Description: "The ID of a transition that moves issues to a closed state. You can find this number under the JIRA workflow administration (Administration > Issues > Workflows) by selecting View under Operations of the desired workflow of your project. By default, this ID is set to 2.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"commit_events": {
				Description: "Enable notifications for commit events",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"merge_requests_events": {
				Description: "Enable notifications for merge request events",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"comment_on_event_enabled": {
				Description: "Enable comments inside Jira issues on each GitLab event (commit / merge request)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceGitlabServiceJiraCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)

	jiraOptions, err := expandJiraOptions(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Create Gitlab Jira service")

	if _, err := client.Services.SetJiraService(project, jiraOptions); err != nil {
		return fmt.Errorf("couldn't create Gitlab Jira service: %w", err)
	}

	d.SetId(project)

	return resourceGitlabServiceJiraRead(d, meta)
}

func resourceGitlabServiceJiraRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	p, resp, err := client.Projects.GetProject(project, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Removing Gitlab Jira service %s because project %s not found", d.Id(), p.Name)
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Read Gitlab Jira service %s", d.Id())

	jiraService, _, err := client.Services.GetJiraService(project)
	if err != nil {
		return err
	}

	if v := jiraService.Properties.URL; v != "" {
		d.Set("url", v)
	}
	if v := jiraService.Properties.Username; v != "" {
		d.Set("username", v)
	}
	if v := jiraService.Properties.ProjectKey; v != "" {
		d.Set("project_key", v)
	}
	if v := jiraService.Properties.JiraIssueTransitionID; v != "" {
		d.Set("jira_issue_transition_id", v)
	}

	d.Set("title", jiraService.Title)
	d.Set("created_at", jiraService.CreatedAt.String())
	d.Set("updated_at", jiraService.UpdatedAt.String())
	d.Set("active", jiraService.Active)
	d.Set("push_events", jiraService.PushEvents)
	d.Set("issues_events", jiraService.IssuesEvents)
	d.Set("commit_events", jiraService.CommitEvents)
	d.Set("merge_requests_events", jiraService.MergeRequestsEvents)
	d.Set("comment_on_event_enabled", jiraService.CommentOnEventEnabled)
	d.Set("tag_push_events", jiraService.TagPushEvents)
	d.Set("note_events", jiraService.NoteEvents)
	d.Set("pipeline_events", jiraService.PipelineEvents)
	d.Set("job_events", jiraService.JobEvents)

	return nil
}

func resourceGitlabServiceJiraUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServiceJiraCreate(d, meta)
}

func resourceGitlabServiceJiraDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)

	log.Printf("[DEBUG] Delete Gitlab Jira service %s", d.Id())

	_, err := client.Services.DeleteJiraService(project)

	return err
}

func expandJiraOptions(d *schema.ResourceData) (*gitlab.SetJiraServiceOptions, error) {
	setJiraServiceOptions := gitlab.SetJiraServiceOptions{}

	// Set required properties
	setJiraServiceOptions.URL = gitlab.String(d.Get("url").(string))
	setJiraServiceOptions.ProjectKey = gitlab.String(d.Get("project_key").(string))
	setJiraServiceOptions.Username = gitlab.String(d.Get("username").(string))
	setJiraServiceOptions.Password = gitlab.String(d.Get("password").(string))
	setJiraServiceOptions.CommitEvents = gitlab.Bool(d.Get("commit_events").(bool))
	setJiraServiceOptions.MergeRequestsEvents = gitlab.Bool(d.Get("merge_requests_events").(bool))
	setJiraServiceOptions.CommentOnEventEnabled = gitlab.Bool(d.Get("comment_on_event_enabled").(bool))

	// Set optional properties
	if val := d.Get("jira_issue_transition_id"); val != nil {
		setJiraServiceOptions.JiraIssueTransitionID = gitlab.String(val.(string))
	}

	return &setJiraServiceOptions, nil
}

func resourceGitlabServiceJiraImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())

	return []*schema.ResourceData{d}, nil
}
