
resource "vercel_project" "my_project" {
  // ...
}

resource "vercel_project_domain" "my_domain" {
  project_id = vercel_project.my_project.id // or use a hardcoded value of an existing project
  name       = "domain.test"
  git_branch = "test"
}
