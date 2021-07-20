package gitlab

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Description: "Provide details about a specific project in the gitlab provider. The results include the name of the project, path, description, default branch, etc.",

		Read: dataSourceGitlabProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Description: "The path of the repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"path_with_namespace": {
				Description: "The path of the repository with namespace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "A description of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_branch": {
				Description: "The default branch for the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"request_access_enabled": {
				Description: "Allow users to request member access.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"issues_enabled": {
				Description: "Enable issue tracking for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"merge_requests_enabled": {
				Description: "Enable merge requests for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"pipelines_enabled": {
				Description: "Enable pipelines for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"wiki_enabled": {
				Description: "Enable wiki for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"snippets_enabled": {
				Description: "Enable snippets for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"lfs_enabled": {
				Description: "Enable LFS for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"visibility_level": {
				Description: "Repositories are created as private by default.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"namespace_id": {
				Description: "The namespace (group or user) of the project. Defaults to your user.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"ssh_url_to_repo": {
				Description: "URL that can be provided to `git clone` to clone the",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"http_url_to_repo": {
				Description: "URL that can be provided to `git clone` to clone the",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"web_url": {
				Description: "URL that can be used to find the project in a browser.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"runners_token": {
				Description: "Registration token to use during runner setup.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"archived": {
				Description: "Whether the project is in read-only mode (archived).",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"remove_source_branch_after_merge": {
				Description: "Enable `Delete source branch` option by default for all new merge requests",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			// lintignore: S031 // TODO: Resolve this tfproviderlint issue
			"push_rules": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"author_email_regex": {
							Description: "All commit author emails must match this regex, e.g. `@my-company.com$`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"branch_name_regex": {
							Description: "All branch names must match this regex, e.g. `(feature|hotfix)\\/*`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"commit_message_regex": {
							Description: "All commit messages must match this regex, e.g. `Fixed \\d+\\..*`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"commit_message_negative_regex": {
							Description: "No commit message is allowed to match this regex, for example `ssh\\:\\/\\/`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"file_name_regex": {
							Description: "All commited filenames must not match this regex, e.g. `(jar|exe)$`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"commit_committer_check": {
							Description: "Users can only push commits to this repository that were committed with one of their own verified emails.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"deny_delete_tag": {
							Description: "Deny deleting a tag.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"member_check": {
							Description: "Restrict commits by author (email) to existing GitLab users.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"prevent_secrets": {
							Description: "GitLab will reject any files that are likely to contain secrets.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"reject_unsigned_commits": {
							Description: "Reject commit when it’s not signed through GPG.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"max_file_size": {
							Description: "Maximum file size (MB).",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab project")

	v, _ := d.GetOk("id")

	found, _, err := client.Projects.GetProject(v, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", found.ID))
	d.Set("name", found.Name)
	d.Set("path", found.Path)
	d.Set("path_with_namespace", found.PathWithNamespace)
	d.Set("description", found.Description)
	d.Set("default_branch", found.DefaultBranch)
	d.Set("request_access_enabled", found.RequestAccessEnabled)
	d.Set("issues_enabled", found.IssuesEnabled)
	d.Set("merge_requests_enabled", found.MergeRequestsEnabled)
	d.Set("pipelines_enabled", found.JobsEnabled)
	d.Set("wiki_enabled", found.WikiEnabled)
	d.Set("snippets_enabled", found.SnippetsEnabled)
	d.Set("visibility_level", string(found.Visibility))
	d.Set("namespace_id", found.Namespace.ID)
	d.Set("ssh_url_to_repo", found.SSHURLToRepo)
	d.Set("http_url_to_repo", found.HTTPURLToRepo)
	d.Set("web_url", found.WebURL)
	d.Set("runners_token", found.RunnersToken)
	d.Set("archived", found.Archived)
	d.Set("remove_source_branch_after_merge", found.RemoveSourceBranchAfterMerge)

	log.Printf("[DEBUG] Reading Gitlab project %q push rules", d.Id())

	pushRules, _, err := client.Projects.GetProjectPushRules(d.Id())
	var httpError *gitlab.ErrorResponse
	if errors.As(err, &httpError) && httpError.Response.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] Failed to get push rules for project %q: %v", d.Id(), err)
	} else if err != nil {
		return fmt.Errorf("Failed to get push rules for project %q: %w", d.Id(), err)
	}

	d.Set("push_rules", flattenProjectPushRules(pushRules)) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue

	return nil
}
