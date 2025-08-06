import { NextRequest, NextResponse } from "next/server"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/lib/auth"
import { prisma } from "@/lib/prisma"

export async function GET() {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const recipes = await prisma.recipe.findMany({
      where: {
        userId: session.user.id
      },
      include: {
        ingredients: true,
        nutrition: true,
      },
      orderBy: {
        createdAt: 'desc'
      }
    })

    return NextResponse.json({ recipes })
  } catch (error) {
    console.error("Error fetching recipes:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}

export async function POST(req: NextRequest) {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const data = await req.json()
    const { name, cookTime, servings, category, instructions, ingredients, nutrition } = data

    const recipe = await prisma.recipe.create({
      data: {
        name,
        cookTime: parseInt(cookTime),
        servings: parseInt(servings),
        category,
        instructions,
        userId: session.user.id,
        ingredients: {
          create: ingredients.map((ingredient: { name: string; amount: string; unit: string }) => ({
            name: ingredient.name,
            amount: ingredient.amount,
            unit: ingredient.unit,
          }))
        },
        nutrition: {
          create: {
            calories: nutrition.calories || 0,
            protein: nutrition.protein || 0,
            carbs: nutrition.carbs || 0,
            fat: nutrition.fat || 0,
            fiber: nutrition.fiber || 0,
            sugar: nutrition.sugar || 0,
          }
        }
      },
      include: {
        ingredients: true,
        nutrition: true,
      }
    })

    return NextResponse.json({ recipe })
  } catch (error) {
    console.error("Error creating recipe:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}