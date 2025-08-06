// Default recipes that new users get when they sign up

export const defaultRecipes = [
  {
    name: "Classic Scrambled Eggs",
    cookTime: 5,
    servings: 2,
    category: "Breakfast",
    instructions: "Beat eggs with salt and pepper. Heat butter in non-stick pan over medium-low heat. Add eggs and gently stir with spatula, pushing eggs from edges to center. Continue until eggs are just set but still creamy. Remove from heat and serve immediately.",
    ingredients: [
      { name: "Large eggs", amount: "4", unit: "" },
      { name: "Butter", amount: "1", unit: "tbsp" },
      { name: "Salt", amount: "1/4", unit: "tsp" },
      { name: "Black pepper", amount: "1/8", unit: "tsp" },
      { name: "Chives", amount: "1", unit: "tbsp" }
    ],
    nutrition: {
      calories: 180,
      protein: 14,
      carbs: 2,
      fat: 13,
      fiber: 0,
      sugar: 1
    }
  },
  {
    name: "Simple Pasta with Marinara",
    cookTime: 15,
    servings: 4,
    category: "Italian",
    instructions: "Cook pasta according to package directions. Meanwhile, heat olive oil in large pan and sauté garlic until fragrant. Add marinara sauce and herbs, simmer 5 minutes. Drain pasta and toss with sauce. Serve with parmesan cheese.",
    ingredients: [
      { name: "Pasta", amount: "1", unit: "lb" },
      { name: "Marinara sauce", amount: "24", unit: "oz" },
      { name: "Garlic", amount: "3", unit: "cloves" },
      { name: "Olive oil", amount: "2", unit: "tbsp" },
      { name: "Italian herbs", amount: "1", unit: "tsp" },
      { name: "Parmesan cheese", amount: "1/2", unit: "cup" }
    ],
    nutrition: {
      calories: 420,
      protein: 15,
      carbs: 75,
      fat: 8,
      fiber: 4,
      sugar: 8
    }
  },
  {
    name: "Grilled Chicken Breast",
    cookTime: 20,
    servings: 4,
    category: "American",
    instructions: "Season chicken breasts with salt, pepper, and garlic powder. Preheat grill to medium-high. Cook chicken 6-7 minutes per side until internal temperature reaches 165°F. Let rest 5 minutes before slicing.",
    ingredients: [
      { name: "Chicken breasts", amount: "4", unit: "" },
      { name: "Olive oil", amount: "2", unit: "tbsp" },
      { name: "Salt", amount: "1", unit: "tsp" },
      { name: "Black pepper", amount: "1/2", unit: "tsp" },
      { name: "Garlic powder", amount: "1", unit: "tsp" },
      { name: "Paprika", amount: "1/2", unit: "tsp" }
    ],
    nutrition: {
      calories: 185,
      protein: 35,
      carbs: 0,
      fat: 4,
      fiber: 0,
      sugar: 0
    }
  },
  {
    name: "Garden Salad with Vinaigrette",
    cookTime: 10,
    servings: 4,
    category: "Salad",
    instructions: "Wash and chop all vegetables. In a small bowl, whisk together olive oil, vinegar, mustard, salt and pepper. Combine lettuce, tomatoes, cucumber, and onion in large bowl. Drizzle with dressing and toss gently.",
    ingredients: [
      { name: "Mixed greens", amount: "6", unit: "cups" },
      { name: "Cherry tomatoes", amount: "1", unit: "cup" },
      { name: "Cucumber", amount: "1", unit: "" },
      { name: "Red onion", amount: "1/4", unit: "" },
      { name: "Olive oil", amount: "3", unit: "tbsp" },
      { name: "Balsamic vinegar", amount: "1", unit: "tbsp" },
      { name: "Dijon mustard", amount: "1", unit: "tsp" }
    ],
    nutrition: {
      calories: 120,
      protein: 3,
      carbs: 8,
      fat: 10,
      fiber: 3,
      sugar: 5
    }
  },
  {
    name: "Beef Tacos",
    cookTime: 20,
    servings: 4,
    category: "Mexican",
    instructions: "Brown ground beef in large skillet over medium-high heat. Add onion and cook until soft. Stir in taco seasoning and water, simmer 5 minutes. Warm tortillas and fill with beef mixture. Top with desired toppings.",
    ingredients: [
      { name: "Ground beef", amount: "1", unit: "lb" },
      { name: "Taco seasoning", amount: "1", unit: "packet" },
      { name: "Yellow onion", amount: "1/2", unit: "" },
      { name: "Corn tortillas", amount: "8", unit: "" },
      { name: "Shredded cheese", amount: "1", unit: "cup" },
      { name: "Lettuce", amount: "2", unit: "cups" },
      { name: "Tomato", amount: "1", unit: "" },
      { name: "Sour cream", amount: "1/2", unit: "cup" }
    ],
    nutrition: {
      calories: 380,
      protein: 25,
      carbs: 28,
      fat: 18,
      fiber: 4,
      sugar: 4
    }
  },
  {
    name: "Vegetable Stir Fry",
    cookTime: 15,
    servings: 3,
    category: "Asian",
    instructions: "Heat oil in large wok or skillet over high heat. Add garlic and ginger, stir-fry 30 seconds. Add vegetables in order of cooking time needed (hard vegetables first). Stir-fry 3-5 minutes until crisp-tender. Add sauce and toss to coat.",
    ingredients: [
      { name: "Vegetable oil", amount: "2", unit: "tbsp" },
      { name: "Mixed vegetables", amount: "4", unit: "cups" },
      { name: "Garlic", amount: "3", unit: "cloves" },
      { name: "Fresh ginger", amount: "1", unit: "tbsp" },
      { name: "Soy sauce", amount: "3", unit: "tbsp" },
      { name: "Sesame oil", amount: "1", unit: "tsp" },
      { name: "Green onions", amount: "2", unit: "" }
    ],
    nutrition: {
      calories: 140,
      protein: 4,
      carbs: 12,
      fat: 10,
      fiber: 4,
      sugar: 8
    }
  },
  {
    name: "Banana Smoothie",
    cookTime: 5,
    servings: 2,
    category: "Breakfast",
    instructions: "Add all ingredients to blender. Blend on high speed until smooth and creamy, about 60 seconds. Add more milk if needed for desired consistency. Pour into glasses and serve immediately.",
    ingredients: [
      { name: "Ripe bananas", amount: "2", unit: "" },
      { name: "Greek yogurt", amount: "1/2", unit: "cup" },
      { name: "Milk", amount: "1", unit: "cup" },
      { name: "Honey", amount: "1", unit: "tbsp" },
      { name: "Vanilla extract", amount: "1/2", unit: "tsp" },
      { name: "Ice cubes", amount: "1/2", unit: "cup" }
    ],
    nutrition: {
      calories: 160,
      protein: 8,
      carbs: 32,
      fat: 2,
      fiber: 3,
      sugar: 24
    }
  },
  {
    name: "Baked Salmon",
    cookTime: 25,
    servings: 4,
    category: "American",
    instructions: "Preheat oven to 400°F. Place salmon fillets on baking sheet lined with parchment paper. Brush with olive oil and season with salt, pepper, and lemon juice. Bake 12-15 minutes until fish flakes easily with a fork.",
    ingredients: [
      { name: "Salmon fillets", amount: "4", unit: "" },
      { name: "Olive oil", amount: "2", unit: "tbsp" },
      { name: "Lemon", amount: "1", unit: "" },
      { name: "Salt", amount: "1/2", unit: "tsp" },
      { name: "Black pepper", amount: "1/4", unit: "tsp" },
      { name: "Fresh dill", amount: "2", unit: "tbsp" }
    ],
    nutrition: {
      calories: 280,
      protein: 40,
      carbs: 1,
      fat: 12,
      fiber: 0,
      sugar: 1
    }
  },
  {
    name: "Rice and Beans",
    cookTime: 25,
    servings: 6,
    category: "American",
    instructions: "Cook rice according to package directions. Heat oil in large pan and sauté onion until soft. Add garlic, cumin, and paprika, cook 1 minute. Add beans, tomatoes, and seasonings. Simmer 10 minutes. Serve over rice.",
    ingredients: [
      { name: "White rice", amount: "1.5", unit: "cups" },
      { name: "Black beans", amount: "2", unit: "cans" },
      { name: "Yellow onion", amount: "1", unit: "" },
      { name: "Garlic", amount: "3", unit: "cloves" },
      { name: "Diced tomatoes", amount: "1", unit: "can" },
      { name: "Cumin", amount: "1", unit: "tsp" },
      { name: "Paprika", amount: "1", unit: "tsp" },
      { name: "Olive oil", amount: "2", unit: "tbsp" }
    ],
    nutrition: {
      calories: 285,
      protein: 11,
      carbs: 58,
      fat: 3,
      fiber: 10,
      sugar: 4
    }
  },
  {
    name: "Chocolate Chip Cookies",
    cookTime: 15,
    servings: 24,
    category: "Dessert",
    instructions: "Preheat oven to 375°F. Cream butter and sugars until light and fluffy. Beat in eggs and vanilla. Mix in flour, baking soda, and salt. Stir in chocolate chips. Drop rounded tablespoons onto ungreased baking sheets. Bake 9-11 minutes until golden brown.",
    ingredients: [
      { name: "Butter", amount: "1", unit: "cup" },
      { name: "Brown sugar", amount: "3/4", unit: "cup" },
      { name: "White sugar", amount: "1/4", unit: "cup" },
      { name: "Eggs", amount: "2", unit: "" },
      { name: "Vanilla extract", amount: "2", unit: "tsp" },
      { name: "All-purpose flour", amount: "2.25", unit: "cups" },
      { name: "Baking soda", amount: "1", unit: "tsp" },
      { name: "Salt", amount: "1", unit: "tsp" },
      { name: "Chocolate chips", amount: "2", unit: "cups" }
    ],
    nutrition: {
      calories: 185,
      protein: 2,
      carbs: 26,
      fat: 9,
      fiber: 1,
      sugar: 16
    }
  }
]