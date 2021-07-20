# Example Usage - Project
resource "gitlab_deploy_token" "example" {
  project    = "example/deploying"
  name       = "Example deploy token"
  username   = "example-username"
  expires_at = "2020-03-14T00:00:00.000Z"
  
  scopes = [ "read_repository", "read_registry" ]
}

# Example Usage - Group
resource "gitlab_deploy_token" "example" {
  group      = "example/deploying"
  name       = "Example group deploy token"
  
  scopes = [ "read_repository" ]
}
