env "local" {
  url = getenv("DATABASE_URL")
  src = "file://db/schema.sql"
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
  }

  schema {
    src = "file://db/schema.sql"
  }
}

env "ci" {
  src = "file://db/schema.sql"
  url = getenv("CI_DATABASE_URL")
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
  }

  schema {
    src = "file://db/schema.sql"
  }

  lint {
    git {
      base = "origin/main"
    }
  }
}

env "prod" {
  url = getenv("NEON_DATABASE_URL")
  schemas = ["public"]
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
    revisions_schema = "public"
  }
}
