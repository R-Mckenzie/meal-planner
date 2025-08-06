"use client"

import * as React from "react"
import { ThemeProvider as NextThemesProvider } from "next-themes"

interface ThemeProviderProps {
  children: React.ReactNode
  attribute?: string
  defaultTheme?: string
  enableSystem?: boolean
  disableTransitionOnChange?: boolean
  storageKey?: string
}

export function ThemeProvider({ 
  children, 
  attribute = "class",
  defaultTheme = "system", 
  enableSystem = true,
  disableTransitionOnChange = false,
  storageKey = "theme"
}: ThemeProviderProps) {
  const [mounted, setMounted] = React.useState(false)

  React.useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) {
    return <>{children}</>
  }

  const themeProps = {
    attribute: attribute as "class" | "data-theme",
    defaultTheme,
    enableSystem,
    disableTransitionOnChange,
    storageKey
  }

  return (
    <NextThemesProvider {...themeProps}>
      {children}
    </NextThemesProvider>
  )
}
