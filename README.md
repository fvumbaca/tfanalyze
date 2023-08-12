# Terraform Summarize

> 🚧 Under Construction 🚧

A tool to sumarize and enforce rules around Terraform plan files.

## Example

Say, we are deploying a simple app with a database through Terraform:

```hcl
resource "google_sql_database" "database" {
  name     = "my-database"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_database_instance" "instance" {
  name             = "my-database-instance"
  region           = "us-central1"
  database_version = "MYSQL_8_0"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = "true"
}

# .... more deployment resources ....

```

When we are deploying, we want to make sure we do not make a horrible mistake.
We are responcible engineers, so we run a terraform plan before applying on our
pipelines:

```text
null_resource.foo: Refreshing state... [id=923362935913696576]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  +   create
  -   destroy
  -/+ recreate

Terraform will perform the following actions:

  # google_sql_database.database will be deleted and recreated
  -/+ resource "null_resource" "bar" {
    ...
    }

[ .... say there is 9 more changes ... ]

Plan: 10 to add, 0 to change, 1 to destroy.

─────────────────────────────────────────────────────────────────────────────
```

And this time around, Tom - a new member on the team, is the one approving merge
requests. For whatever reason (ie. being new to Terraform, being in a rush,
pipeline too slow and set to auto-merge, etc) approves the merge request.

Sometimes, there are resources you know should never be deleted:

```toml
# rules.hcl

resource "google_sql_database" "database" {
  on_delete {
    action = "error"
    message = "We should not be deleting the database! Even if a change re-creates it!"
  }
}
```

By running `tfanalyze` you can stop a dangerous apply in it's tracks:

```
ERROR: We should not be deleting the database! Even if a change re-creates it!

exit-status 1
```
