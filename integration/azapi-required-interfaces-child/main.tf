module "child" {
  source = "./child"
}

output "resource_id" {
  value = module.child.resource_id
}
