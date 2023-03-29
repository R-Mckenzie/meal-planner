INSERT INTO users (email, hash, remember_hash, created_at, updated_at) 
VALUES ('defaultrecipeuser', '134asf@£1sa@£$%!df£@£$1', '123ASF@123$!1341!@£$T$@£', current_timestamp, current_timestamp);

INSERT INTO recipes (owner_id, title, ingredients, method, created_at, updated_at) 
VALUES (0, 'Chilli Con Carne','{"Mince", "Chopped Tomatoes", "Kidney Beans", "Chillis", "Onions", "Seasoning"}', 'Fry mince and onions \nAdd chopped tomatoes \nAdd kidney beans and seasoning and simmer', current_timestamp, current_timestamp),
(0, 'Spaghetti Bolognese','{"Mince", "Chopped Tomatoes", "Kidney Beans", "Chillis", "Onions", "Seasoning"}', 'Fry mince and onions \nAdd chopped tomatoes \nAdd kidney beans and seasoning and simmer', current_timestamp, current_timestamp),
(0, 'Fajitas','{"Chicken", "Peppers", "onions", "Cayenne Pepper", "Salt", "Paprika", "Tortilla Wraps"}', 'Chicken and vegetables\nAdd seasoning \nServe in warm tortilla wraps', current_timestamp, current_timestamp);
