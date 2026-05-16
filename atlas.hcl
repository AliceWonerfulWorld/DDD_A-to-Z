env "local" {
  url = getenv("DATABASE_URL")
  src = "file://db/schema"
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
  }

  schema {
    src = "file://db/schema"
  }
}

env "ci" {
  src = "file://db/schema"
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
  }

  schema {
    src = "file://db/schema"
  }

  lint {
    git {
      base = "origin/main"
    }
  }
}
