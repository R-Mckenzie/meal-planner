table "users" {
	schema = schema.public

	column "id" {
		null = false
		type = integer
		identity {
			generated = ALWAYS
			start = 0
			increment = 1
		}
	}

	column "email" {
		null = false
		type = varchar(255)
	}

	column "hash" {
		null = false
		type = char(60)
	}

	column "remember_hash" {
		null = false
		type = char(60)
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
		type = timestamp
		null = true
	}

	index "idx_email" {
		columns = [column.email]
		unique = true
	}

	index "idx_rememberHash" {
		columns = [column.remember_hash]
		unique = true
	}

	primary_key {
		columns = [column.id]
	}
}

table "recipes" {
	schema = schema.public

	column "id" {
		null = false
		type = integer
		identity {
			generated = ALWAYS
			start = 0
			increment = 1
		}
	}

	column "owner_id" {
		null = false
		type = integer
	}

	column "title" {
		null = false
		type = varchar(255)
	}

	column "ingredients" {
		null = false
		type = sql("varchar(50)[]")
	}

	column "method" {
		null = true
		type = varchar(500)
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
		type = timestamp
		null = true
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "recipe_owner_fk"{
		columns = [column.owner_id]	
		ref_columns = [table.users.column.id]
	}
}

table "meals" {
	schema = schema.public

	column "id" {
		null = false
		type = integer
		identity {
			generated = ALWAYS
			start = 0
			increment = 1
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
		type = timestamp
		null = true
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "meal_owner_fk"{
		columns = [column.owner_id]	
		ref_columns = [table.users.column.id]
	}
	
	foreign_key "recipe_fk" {
		columns = [column.recipe_id]
		ref_columns = [table.recipes.column.id]
	}
}

schema "public" {
}
