---
layout: "gitlab"
page_title: "GitLab: gitlab_group"
sidebar_current: "docs-gitlab-resource-group"
description: |-
  Creates and manages GitLab groups
---

# gitlab\_group

This resource allows you to create and manage GitLab groups.
Note your provider will need to be configured with admin-level access for this resource to work.

## Example Usage

```hcl
resource "gitlab_group" "example" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

// Create a project in the example group
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = "${gitlab_group.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this group.

* `path` - (Required) The path of the group.

* `description` - (Optional) The description of the group.

* `lfs_enabled` - (Optional) Boolean, defaults to true.  Whether to enable LFS
support for projects in this group.

* `request_access_enabled` - (Optional) Boolean, defaults to false.  Whether to
enable users to request access to the group.

* `visibility_level` - (Optional) Set to `public` to create a public group.
  Valid values are `private`, `internal`, `public`.
  Groups are created as private by default.

* `parent_id` - (Optional) Integer, id of the parent group (creates a nested group).

* `membership_lock` - (Optional) Boolean, defaults to false. Prevent adding new members to project membership within this group. Requires GitLab Starter/Bronze.

* `share_with_group_lock` (Optional) Boolean, defaults to false. Prevent sharing a project with another group within this group.

* `require_two_factor_authentication` (Optional) Boolean, defaults to false. Require all users in this group to setup Two-factor authentication.

* `two_factor_grace_period` (Optional) Integer, time before Two-factor authentication is enforced (in hours). The default grace period is 48 hours.

* `project_creation_level` (Optional) Determine if developers can create projects in the group. Can be `noone` (No one), `maintainer` (Maintainers), or `developer` (Developers + Maintainers). By default Developers and Maintainers can create projects. 

* `auto_devops_enabled` (Optional) Boolean, defaults to false. Default to Auto DevOps pipeline for all projects within this group.

* `subgroup_creation_level` (Optional) Allowed to create subgroups. Can be `owner` (Owners), or `maintainer` (Maintainers). By default only Owners can create subgroups.

* `emails_disabled` (Optional) Boolean, defaults to false. 	Disable email notifications.

* `mentions_disabled` (Optional) Boolean, defaults to false. Disable the capability of a group from getting mentioned.

* `shared_runners_minutes_limit` (Optional) Integer, pipeline minutes quota for this group (included in plan). Can be nil (default; inherit system default), 0 (unlimited) or > 0. Requires GitLab Starter.

* `extra_shared_runners_minutes_limit` (Optional) Integer, extra pipeline minutes quota for this group (purchased in addition to the minutes included in the plan). Requires GitLab Starter.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the group by the GitLab server.  Serves as a
  namespace id where one is needed.
  
* `full_path` - The full path of the group.

* `full_name` - The full name of the group.

* `web_url` - Web URL of the group.

* `runners_token` - The group level registration token to use during runner setup.

## Importing groups

You can import a group state using `terraform import <resource> <id>`.  The
`id` can be whatever the [details of a group][details_of_a_group] api takes for
its `:id` value, so for example:

    terraform import gitlab_group.example example

[details_of_a_group]: https://docs.gitlab.com/ee/api/groups.html#details-of-a-group
