resource "local_file" "foo" {
  content  = "foo!!!!"
  filename = "${path.module}/foo"
}

resource "local_file" "bar" {
  content  = "bar!"
  filename = "${path.module}/bar"
}
