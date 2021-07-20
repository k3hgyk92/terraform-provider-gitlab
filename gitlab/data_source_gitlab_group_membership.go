package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabGroupMembership() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))
	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Description: "Provide details about a list of group members in the gitlab provider. The results include id, username, name and more about the requested members.",

		Read: dataSourceGitlabGroupMembershipRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"full_path",
				},
			},
			"full_path": {
				Description: "The full path of the group.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"group_id",
				},
			},
			"access_level": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
			},
			"members": {
				Description: "The list of group members.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The unique id assigned to the user by the gitlab server.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"username": {
							Description: "The username of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"state": {
							Description: "Whether the user is active or blocked.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"avatar_url": {
							Description: "The avatar URL of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"web_url": {
							Description: "User's website URL.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires_at": {
							Description: "Expiration date for the group membership.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var gm []*gitlab.GroupMember
	var group *gitlab.Group
	var err error

	log.Printf("[INFO] Reading Gitlab group")

	groupIDData, groupIDOk := d.GetOk("group_id")
	fullPathData, fullPathOk := d.GetOk("full_path")

	if groupIDOk {
		// Get group by id
		group, _, err = client.Groups.GetGroup(groupIDData.(int))
		if err != nil {
			return err
		}
	} else if fullPathOk {
		// Get group by full path
		group, _, err = client.Groups.GetGroup(fullPathData.(string))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("one and only one of group_id or full_path must be set")
	}

	log.Printf("[INFO] Reading Gitlab group memberships")

	// Get group memberships
	gm, _, err = client.Groups.ListGroupMembers(group.ID, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		return err
	}

	d.Set("group_id", group.ID)
	d.Set("full_path", group.FullPath)

	d.Set("members", flattenGitlabMembers(d, gm)) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue

	var optionsHash strings.Builder
	optionsHash.WriteString(strconv.Itoa(group.ID))

	if data, ok := d.GetOk("access_level"); ok {
		optionsHash.WriteString(data.(string))
	}

	id := schema.HashString(optionsHash.String())
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}

func flattenGitlabMembers(d *schema.ResourceData, members []*gitlab.GroupMember) []interface{} {
	membersList := []interface{}{}

	var filterAccessLevel gitlab.AccessLevelValue = gitlab.NoPermissions
	if data, ok := d.GetOk("access_level"); ok {
		filterAccessLevel = accessLevelID[data.(string)]
	}

	for _, member := range members {
		if filterAccessLevel != gitlab.NoPermissions && filterAccessLevel != member.AccessLevel {
			continue
		}

		values := map[string]interface{}{
			"id":           member.ID,
			"username":     member.Username,
			"name":         member.Name,
			"state":        member.State,
			"avatar_url":   member.AvatarURL,
			"web_url":      member.WebURL,
			"access_level": accessLevel[gitlab.AccessLevelValue(member.AccessLevel)],
		}

		if member.ExpiresAt != nil {
			values["expires_at"] = member.ExpiresAt.String()
		}

		membersList = append(membersList, values)
	}

	return membersList
}
