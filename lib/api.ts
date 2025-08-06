// API helper functions for meal planner

export interface Recipe {
  id: string
  name: string
  cookTime: number
  servings: number
  category: string
  instructions?: string
  ingredients: Ingredient[]
  nutrition: NutritionalInfo
}

export interface Ingredient {
  id: string
  name: string
  amount: string
  unit: string
}

export interface NutritionalInfo {
  id: string
  calories: number
  protein: number
  carbs: number
  fat: number
  fiber: number
  sugar: number
}

export interface MealPlan {
  id: string
  date: string
  recipes: MealPlanRecipe[]
}

export interface MealPlanRecipe {
  id: string
  mealType: 'BREAKFAST' | 'LUNCH' | 'DINNER' | 'SNACK'
  recipe: Recipe
}

export async function fetchRecipes(): Promise<Recipe[]> {
  const response = await fetch('/api/recipes')
  if (!response.ok) {
    throw new Error('Failed to fetch recipes')
  }
  const data = await response.json()
  return data.recipes
}

interface CreateRecipeData {
  name: string
  cookTime: number
  servings: number
  category: string
  instructions?: string
  ingredients: { name: string; amount: string; unit: string }[]
  nutrition: {
    calories: number
    protein: number
    carbs: number
    fat: number
    fiber: number
    sugar: number
  }
}

export async function createRecipe(recipe: CreateRecipeData): Promise<Recipe> {
  const response = await fetch('/api/recipes', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(recipe),
  })
  if (!response.ok) {
    throw new Error('Failed to create recipe')
  }
  const data = await response.json()
  return data.recipe
}

export async function updateRecipe(id: string, recipe: CreateRecipeData): Promise<Recipe> {
  const response = await fetch(`/api/recipes/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(recipe),
  })
  if (!response.ok) {
    throw new Error('Failed to update recipe')
  }
  const data = await response.json()
  return data.recipe
}

export async function deleteRecipe(id: string): Promise<void> {
  const response = await fetch(`/api/recipes/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to delete recipe')
  }
}

export async function fetchMealPlans(startDate: string, endDate: string): Promise<MealPlan[]> {
  const response = await fetch(`/api/meal-plans?startDate=${startDate}&endDate=${endDate}`)
  if (!response.ok) {
    throw new Error('Failed to fetch meal plans')
  }
  const data = await response.json()
  return data.mealPlans
}

export async function addRecipeToMealPlan(date: string, recipeId: string, mealType: string): Promise<MealPlanRecipe> {
  const response = await fetch('/api/meal-plans', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ date, recipeId, mealType }),
  })
  if (!response.ok) {
    throw new Error('Failed to add recipe to meal plan')
  }
  const data = await response.json()
  return data.mealPlanRecipe
}

export async function removeRecipeFromMealPlan(mealPlanRecipeId: string): Promise<void> {
  const response = await fetch(`/api/meal-plans/${mealPlanRecipeId}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to remove recipe from meal plan')
  }
}

export async function clearMealPlan(date: string): Promise<void> {
  const response = await fetch(`/api/meal-plans?date=${date}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to clear meal plan')
  }
}