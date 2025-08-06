import { NextRequest, NextResponse } from "next/server"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/lib/auth"
import { prisma } from "@/lib/prisma"

export async function PUT(
  req: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const { id } = await params
    const data = await req.json()
    const { name, cookTime, servings, category, instructions, ingredients, nutrition } = data

    // Verify recipe ownership
    const existingRecipe = await prisma.recipe.findFirst({
      where: {
        id,    
        userId: session.user.id
      }
    })

    if (!existingRecipe) {
      return NextResponse.json({ error: "Recipe not found" }, { status: 404 })
    }

    // Update recipe with transaction
    const recipe = await prisma.$transaction(async (tx) => {
      // Delete existing ingredients and nutrition
      await tx.ingredient.deleteMany({
        where: { recipeId: id }
      })
      
      await tx.nutritionalInfo.deleteMany({
        where: { recipeId: id }
      })

      // Update recipe with new data
      return await tx.recipe.update({
        where: { id },
        data: {
          name,
          cookTime: parseInt(cookTime),
          servings: parseInt(servings),
          category,
          instructions,
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
    })

    return NextResponse.json({ recipe })
  } catch (error) {
    console.error("Error updating recipe:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}

export async function DELETE(
  req: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const { id } = await params

    // Verify recipe ownership
    const existingRecipe = await prisma.recipe.findFirst({
      where: {
        id,
        userId: session.user.id
      }
    })

    if (!existingRecipe) {
      return NextResponse.json({ error: "Recipe not found" }, { status: 404 })
    }

    await prisma.recipe.delete({
      where: { id }
    })

    return NextResponse.json({ message: "Recipe deleted successfully" })
  } catch (error) {
    console.error("Error deleting recipe:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}