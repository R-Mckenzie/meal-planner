import { NextResponse } from "next/server"
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/lib/auth"
import { seedDefaultRecipesForUser } from "@/lib/seed-defaults"

export async function POST() {
  try {
    const session = await getServerSession(authOptions)
    
    if (!session?.user?.id) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const count = await seedDefaultRecipesForUser(session.user.id)

    return NextResponse.json({ 
      message: `Added ${count} default recipes`,
      count 
    })
  } catch (error) {
    console.error("Error seeding default recipes:", error)
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    )
  }
}