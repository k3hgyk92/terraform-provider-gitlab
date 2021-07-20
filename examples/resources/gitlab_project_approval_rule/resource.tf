resource "gitlab_project_approval_rule" "example-one" {
  project            = 5
  name               = "Example Rule 1"
  approvals_required = 3
  user_ids           = [50, 500]
  group_ids          = [51]
}

resource "gitlab_branch_protection" "example" {
  project            = 5
  branch             = "main"
  push_access_level  = "developer"
  merge_access_level = "developer"
}

resource "gitlab_project_approval_rule" "example-two" {
  project              = 5
  name                 = "Example Rule 2"
  approvals_required   = 1
  user_ids             = []
  group_ids            = [52]
  protected_branch_ids = [gitlab_branch_protection.example.id]
}
