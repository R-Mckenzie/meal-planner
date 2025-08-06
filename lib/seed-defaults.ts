// Utility function to seed default recipes for existing users
import { prisma } from "./prisma"
import { defaultRecipes } from "./default-recipes"

export async function seedDefaultRecipesForUser(userId: string) {
  // Check if user already has recipes
  const existingRecipes = await prisma.recipe.count({
    where: { userId }
  })

  // Only add default recipes if user has no recipes
  if (existingRecipes === 0) {
    const recipesData = defaultRecipes.map(recipe => ({
      name: recipe.name,
      cookTime: recipe.cookTime,
      servings: recipe.servings,
      category: recipe.category,
      instructions: recipe.instructions,
      userId: userId,
      ingredients: {
        create: recipe.ingredients.map(ingredient => ({
          name: ingredient.name,
          amount: ingredient.amount,
          unit: ingredient.unit,
        }))
      },
      nutrition: {
        create: {
          calories: recipe.nutrition.calories,
          protein: recipe.nutrition.protein,
          carbs: recipe.nutrition.carbs,
          fat: recipe.nutrition.fat,
          fiber: recipe.nutrition.fiber,
          sugar: recipe.nutrition.sugar,
        }
      }
    }))

    // Create all default recipes
    for (const recipeData of recipesData) {
      await prisma.recipe.create({
        data: recipeData
      })
    }

    console.log(`Added ${defaultRecipes.length} default recipes for user ${userId}`)
    return defaultRecipes.length
  }

  return 0
}