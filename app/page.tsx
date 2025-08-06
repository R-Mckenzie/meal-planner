"use client"

import { useSession } from "next-auth/react"
import { ThemeProvider } from "@/components/theme-provider"
import MealPlanner from "../meal-planner"
import LandingPage from "@/components/landing-page"

export default function Page() {
  const { data: session, status } = useSession()

  if (status === "loading") {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  if (!session) {
    return (
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
        storageKey="meal-planner-theme"
      >
        <LandingPage />
      </ThemeProvider>
    )
  }

  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
      storageKey="meal-planner-theme"
    >
      <MealPlanner />
    </ThemeProvider>
  )
}