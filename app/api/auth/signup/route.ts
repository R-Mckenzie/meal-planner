import { NextRequest, NextResponse } from "next/server"
import bcrypt from "bcryptjs"
import { prisma } from "@/lib/prisma"
import { defaultRecipes } from "@/lib/default-recipes"

export async function POST(req: NextRequest) {
  try {
    const { name, email, password } = await req.json()

    if (!name || !email || !password) {
      return NextResponse.json(
        { error: "Missing required fields" },
        { status: 400 }
      )
    }

    // Check if user already exists
    const existingUser = await prisma.user.findUnique({
      where: { email }
    })

    if (existingUser) {
      return NextResponse.json(
        { error: "User already exists" },
        { status: 400 }
      )
    }

    // Hash password
    const hashedPassword = await bcrypt.hash(password, 12)

    // Create user and default recipes in a transaction
    const result = await prisma.$transaction(async (tx) => {
      // Create user
      const user = await tx.user.create({
        data: {
          name,
          email,
          password: hashedPassword,
        },
      })

      // Create default recipes for the new user
      const recipesData = defaultRecipes.map(recipe => ({
        name: recipe.name,
        cookTime: recipe.cookTime,
        servings: recipe.servings,
        category: recipe.category,
        instructions: recipe.instructions,
        userId: user.id,
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
        await tx.recipe.create({
          data: recipeData
        })
      }

      return user
    })

    return NextResponse.json({
      user: {
        id: result.id,
        name: result.name,
        email: result.email,
      },
    })
  } catch (error) {
    console.error("Signup error:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}