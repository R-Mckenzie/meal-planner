table "meals" {
  schema = schema.public
  column "id" {
    null = false
    type = integer
    identity {
      generated = ALWAYS
    }
  }
  column "owner_id" {
    null = false
    type = integer
  }
  column "recipe_id" {
    null = false
    type = integer
  }
  column "date" {
    null = false
    type = date
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "meal_owner_fk" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "recipe_fk" {
    columns     = [column.recipe_id]
    ref_columns = [table.recipes.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "recipes" {
  schema = schema.public
  column "id" {
    null = false
    type = integer
    identity {
      generated = ALWAYS
    }
  }
  column "owner_id" {
    null = false
    type = integer
  }
  column "title" {
    null = false
    type = character_varying(255)
  }
  column "ingredients" {
    null = false
    type = sql("character varying(50)[]")
  }
  column "method" {
    null = true
    type = character_varying(500)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "recipe_owner_fk" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = integer
    identity {
      generated = ALWAYS
    }
  }
  column "email" {
    null = false
    type = character_varying(255)
  }
  column "hash" {
    null = false
    type = character(60)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "remember_hash" {
    null = false
    type = character(60)
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_email" {
    unique  = true
    columns = [column.email]
  }
  index "idx_rememberHash" {
    unique  = true
    columns = [column.remember_hash]
  }
}
schema "public" {
}
