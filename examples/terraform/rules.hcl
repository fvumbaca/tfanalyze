resource "null_resource" "foo" {
  on_destroy {
    severity = "warn"
    message = "A message describing why foo being deleted is bad."
  }
}

resource "null_resource" "bar" {
  on_create {
    severity = "warn"
    message = "A message explaining why bar being created is bad."
  }
}

resource "null_resource" "i_dont_exist" {
  on_create {
    severity = "inform"
    message = "I do not exist"
  }
}

list_resource "null_resource" "im_a_count" "1" {
  on_create {
    severity = "error"
    message = <<EOF
A Message about counting
    EOF
  }
}