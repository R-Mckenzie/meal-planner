INSERT INTO recipes (owner_id, title, ingredients, method, created_at, updated_at)
VALUES
(1, 'Classic Margherita Pizza', 
 ARRAY['Pizza dough', 'Tomato sauce', 'Fresh mozzarella', 'Fresh basil', 'Olive oil'],
 'Roll out the dough, spread tomato sauce, add sliced mozzarella and bake. Top with fresh basil and a drizzle of olive oil after baking.',
 CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(1, 'Chocolate Chip Cookies', 
 ARRAY['Flour', 'Butter', 'Brown sugar', 'White sugar', 'Eggs', 'Vanilla extract', 'Baking soda', 'Salt', 'Chocolate chips'],
 'Cream butter and sugars, add eggs and vanilla. Mix in dry ingredients, fold in chocolate chips. Bake at 375°F for 9-11 minutes.',
 CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(1, 'Chicken Stir-Fry', 
 ARRAY['Chicken breast', 'Mixed vegetables', 'Soy sauce', 'Garlic', 'Ginger', 'Vegetable oil', 'Cornstarch'],
 'Stir-fry chicken in oil. Add vegetables, garlic, and ginger. Mix soy sauce with cornstarch and add to pan. Cook until sauce thickens.',
 CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(1, 'Greek Salad', 
 ARRAY['Cucumber', 'Tomatoes', 'Red onion', 'Feta cheese', 'Kalamata olives', 'Olive oil', 'Lemon juice', 'Dried oregano'],
 'Chop vegetables and combine in a bowl. Add crumbled feta and olives. Dress with olive oil, lemon juice, and oregano. Toss to combine.',
 CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
