package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabBranch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabBranchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"merged": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"developers_can_push": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"developers_can_merge": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"can_push": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_short_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_author_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_author_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_authored_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_committed_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_committer_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_committer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_parent_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	name := d.Get("name").(string)
	project := d.Get("project").(string)

	branch, _, err := client.Branches.GetBranch(project, name)
	if err != nil {
		return err
	}

	d.Set("merged", branch.Merged)
	d.Set("protected", branch.Protected)
	d.Set("default", branch.Default)
	d.Set("developers_can_push", branch.DevelopersCanPush)
	d.Set("developers_can_merge", branch.DevelopersCanMerge)
	d.Set("can_push", branch.CanPush)
	d.Set("web_url", branch.WebURL)
	d.Set("commit_id", branch.Commit.ID)
	d.Set("commit_short_id", branch.Commit.ShortID)
	d.Set("commit_author_email", branch.Commit.AuthorEmail)
	d.Set("commit_author_name", branch.Commit.AuthorName)
	d.Set("commit_authored_date", branch.Commit.AuthoredDate)
	d.Set("commit_committed_date", branch.Commit.CommittedDate)
	d.Set("commit_committer_email", branch.Commit.CommitterEmail)
	d.Set("commit_committer_name", branch.Commit.CommitterName)
	d.Set("commit_title", branch.Commit.Title)
	d.Set("commit_message", branch.Commit.Message)
	d.Set("commit_parent_ids", branch.Commit.ParentIDs)

	d.SetId(fmt.Sprintf("%s:%s", project, name))

	return nil
}
