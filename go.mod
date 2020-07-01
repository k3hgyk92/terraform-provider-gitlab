module github.com/gitlabhq/terraform-provider-gitlab

require (
	github.com/hashicorp/terraform-plugin-sdk v1.13.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/xanzy/go-gitlab v0.34.1
)

go 1.14

replace github.com/xanzy/go-gitlab v0.32.1 => github.com/xanzy/go-gitlab v0.32.2-0.20200701195523-9d3a87a48b01
