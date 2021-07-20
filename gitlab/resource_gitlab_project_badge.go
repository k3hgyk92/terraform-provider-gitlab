package gitlab

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectBadge() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage badges for your GitLab projects.\n" +
			"For further information on hooks, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ce/user/project/badges.html).",

		Create: resourceGitlabProjectBadgeCreate,
		Read:   resourceGitlabProjectBadgeRead,
		Update: resourceGitlabProjectBadgeUpdate,
		Delete: resourceGitlabProjectBadgeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The id of the project to add the badge to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"link_url": {
				Description: "The url linked with the badge.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"image_url": {
				Description: "The image url which will be presented on project overview.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"rendered_link_url": {
				Description: "The link_url argument rendered (in case of use of placeholders).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rendered_image_url": {
				Description: "The image_url argument rendered (in case of use of placeholders).",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGitlabProjectBadgeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	projectID := d.Get("project").(string)
	options := &gitlab.AddProjectBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] create gitlab project badge %q / %q", *options.LinkURL, *options.ImageURL)

	badge, _, err := client.ProjectBadges.AddProjectBadge(projectID, options)
	if err != nil {
		return err
	}

	badgeID := strconv.Itoa(badge.ID)

	d.SetId(buildTwoPartID(&projectID, &badgeID))

	return resourceGitlabProjectBadgeRead(d, meta)
}

func resourceGitlabProjectBadgeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	projectID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab project badge %s/%d", projectID, badgeID)

	badge, _, err := client.ProjectBadges.GetProjectBadge(projectID, badgeID)
	if err != nil {
		return err
	}

	resourceGitlabProjectBadgeSetToState(d, badge, &projectID)
	return nil
}

func resourceGitlabProjectBadgeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	projectID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return err
	}

	options := &gitlab.EditProjectBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] update gitlab project badge %s/%d", projectID, badgeID)

	_, _, err = client.ProjectBadges.EditProjectBadge(projectID, badgeID, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectBadgeRead(d, meta)
}

func resourceGitlabProjectBadgeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	projectID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Delete gitlab project badge %s/%d", projectID, badgeID)

	_, err = client.ProjectBadges.DeleteProjectBadge(projectID, badgeID)
	return err
}

func resourceGitlabProjectBadgeSetToState(d *schema.ResourceData, badge *gitlab.ProjectBadge, projectID *string) {
	d.Set("link_url", badge.LinkURL)
	d.Set("image_url", badge.ImageURL)
	d.Set("rendered_link_url", badge.RenderedLinkURL)
	d.Set("rendered_image_url", badge.RenderedImageURL)
	d.Set("project", projectID)
}
