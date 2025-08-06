import { NextRequest, NextResponse } from "next/server"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/lib/auth"
import { prisma } from "@/lib/prisma"

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

    // Verify ownership through meal plan
    const mealPlanRecipe = await prisma.mealPlanRecipe.findFirst({
      where: {
        id,
        mealPlan: {
          userId: session.user.id
        }
      }
    })

    if (!mealPlanRecipe) {
      return NextResponse.json({ error: "Meal plan recipe not found" }, { status: 404 })
    }

    await prisma.mealPlanRecipe.delete({
      where: { id }
    })

    return NextResponse.json({ message: "Recipe removed from meal plan successfully" })
  } catch (error) {
    console.error("Error removing recipe from meal plan:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}