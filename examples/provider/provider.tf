provider "metabase" {
  endpoint = "http://your-metabase-instance.com/api"
  username = "your-metabase-admin-username" # or use api_key
  password = "your-metabase-admin-password" # or use api_key
  api_key = "your-metabase-admin-api-key" # or use username/password
}
