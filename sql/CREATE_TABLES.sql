CREATE TABLE users (
	id SERIAL PRIMARY KEY NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_pass VARCHAR(60) NOT NULL,
	role SMALLINT DEFAULT 0,
	created TIMESTAMP NOT NULL,
	updated TIMESTAMP NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE(email);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE recipes (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    ingredients VARCHAR(50)[] NOT NULL,
    method VARCHAR(500),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT recipe_owner_fk FOREIGN KEY (owner_id)
        REFERENCES public.users (id) ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE TABLE meals (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    recipe_id INTEGER NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT meal_owner_fk FOREIGN KEY (owner_id)
        REFERENCES public.users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT recipe_fk FOREIGN KEY (recipe_id)
        REFERENCES public.recipes (id) ON DELETE CASCADE ON UPDATE NO ACTION
);
