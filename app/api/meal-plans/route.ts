import { NextRequest, NextResponse } from "next/server"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/lib/auth"
import { prisma } from "@/lib/prisma"

export async function GET(req: NextRequest) {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const { searchParams } = new URL(req.url)
    const startDate = searchParams.get('startDate')
    const endDate = searchParams.get('endDate')

    const where: { userId: string; date?: { gte: Date; lte: Date } } = {
      userId: session.user.id
    }

    if (startDate && endDate) {
      where.date = {
        gte: new Date(startDate),
        lte: new Date(endDate)
      }
    }

    const mealPlans = await prisma.mealPlan.findMany({
      where,
      include: {
        recipes: {
          include: {
            recipe: {
              include: {
                ingredients: true,
                nutrition: true,
              }
            }
          }
        }
      },
      orderBy: {
        date: 'asc'
      }
    })

    return NextResponse.json({ mealPlans })
  } catch (error) {
    console.error("Error fetching meal plans:", error)
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
    const { date, recipeId, mealType } = data

    // Find or create meal plan for the date
    let mealPlan = await prisma.mealPlan.findFirst({
      where: {
        date: new Date(date),
        userId: session.user.id
      }
    })

    if (!mealPlan) {
      mealPlan = await prisma.mealPlan.create({
        data: {
          date: new Date(date),
          userId: session.user.id
        }
      })
    }

    // Add recipe to meal plan
    const mealPlanRecipe = await prisma.mealPlanRecipe.create({
      data: {
        mealPlanId: mealPlan.id,
        recipeId,
        mealType: mealType.toUpperCase()
      },
      include: {
        recipe: {
          include: {
            ingredients: true,
            nutrition: true,
          }
        }
      }
    })

    return NextResponse.json({ mealPlanRecipe })
  } catch (error) {
    console.error("Error adding recipe to meal plan:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}

export async function DELETE(req: NextRequest) {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const { searchParams } = new URL(req.url)
    const date = searchParams.get('date')

    if (!date) {
      return NextResponse.json({ error: "Date is required" }, { status: 400 })
    }

    // Find meal plan for the date
    const mealPlan = await prisma.mealPlan.findFirst({
      where: {
        date: new Date(date),
        userId: session.user.id
      }
    })

    if (mealPlan) {
      // Delete all meal plan recipes for this date
      await prisma.mealPlanRecipe.deleteMany({
        where: {
          mealPlanId: mealPlan.id
        }
      })

      // Delete the meal plan itself
      await prisma.mealPlan.delete({
        where: {
          id: mealPlan.id
        }
      })
    }

    return NextResponse.json({ message: "Meal plan cleared successfully" })
  } catch (error) {
    console.error("Error clearing meal plan:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}