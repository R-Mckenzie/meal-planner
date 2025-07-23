"use client"

import { ThemeProvider } from "@/components/theme-provider"
import MealPlanner from "../meal-planner"

export default function Page() {
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
