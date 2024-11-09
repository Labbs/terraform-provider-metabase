resource "metabase_database" "mysql" {
  name   = "mysql"
  engine = "mysql"
  mysql_details = {
    host     = "127.0.0.1"
    port     = 3306
    database = "mysql"
    user     = "root"
    password = "SuperSecret"
  }
}

resource "metabase_database" "postgres" {
  name   = "postgres"
  engine = "postgres"
  postgresql_details = {
    host     = "127.0.0.1"
    port     = 5432
    database = "mydatabase"
    user     = "myuser"
    password = "mypassword"
  }
}