variable "database_url" {
  type    = string
  default = "postgres://postgres:zcicd123@localhost:5432/zcicd?sslmode=disable"
}

env "local" {
  src = "file://migrations"
  url = var.database_url
  dev = "docker://postgres/16/dev?search_path=public"
  migration {
    dir = "file://migrations"
  }
}

env "production" {
  src = "file://migrations"
  url = var.database_url
  migration {
    dir    = "file://migrations"
    format = golang-migrate
  }
}
