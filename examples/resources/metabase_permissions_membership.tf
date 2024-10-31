resource "metabase_permissions_membership" "example2" {
  group_id = metabase_permissions_group.example.id
  user_id  = metabase_user.example.id
  is_group_manager = false # can't be changed after creation if you don't have a Premium license
}
