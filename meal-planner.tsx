"use client"

import type React from "react"

import { useState, useEffect } from "react"
import {
  Calendar,
  ChefHat,
  Clock,
  Users,
  ChevronDown,
  Zap,
  ChevronLeft,
  ChevronRight,
  X,
  Menu,
  ShoppingCart,
} from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ThemeToggle } from "@/components/theme-toggle"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Edit, Trash2, Save } from "lucide-react"
import { Textarea } from "@/components/ui/textarea"
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet"
import { Checkbox } from "@/components/ui/checkbox"

interface Ingredient {
  id: string
  name: string
  amount: string
  unit: string
}

interface NutritionalInfo {
  calories: number
  protein: number
  carbs: number
  fat: number
  fiber: number
  sugar: number
}

interface Recipe {
  id: string
  name: string
  cookTime: number
  servings: number
  category: string
  ingredients: Ingredient[]
  nutrition: NutritionalInfo
  instructions?: string
}

interface MealPlan {
  [date: string]: Recipe[]
}

interface ShoppingListItem {
  id: string
  name: string
  amount: number
  unit: string
  recipes: string[]
  checked: boolean
}

const defaultRecipes: Recipe[] = [
  {
    id: "1",
    name: "Spaghetti Carbonara",
    cookTime: 20,
    servings: 4,
    category: "Italian",
    ingredients: [
      { id: "1", name: "Spaghetti", amount: "400", unit: "g" },
      { id: "2", name: "Eggs", amount: "4", unit: "large" },
      { id: "3", name: "Pancetta", amount: "150", unit: "g" },
      { id: "4", name: "Parmesan cheese", amount: "100", unit: "g" },
      { id: "5", name: "Black pepper", amount: "1", unit: "tsp" },
    ],
    nutrition: {
      calories: 520,
      protein: 28,
      carbs: 45,
      fat: 24,
      fiber: 3,
      sugar: 2,
    },
  },
  {
    id: "2",
    name: "Chicken Stir Fry",
    cookTime: 15,
    servings: 3,
    category: "Asian",
    ingredients: [
      { id: "1", name: "Chicken breast", amount: "500", unit: "g" },
      { id: "2", name: "Mixed vegetables", amount: "300", unit: "g" },
      { id: "3", name: "Soy sauce", amount: "3", unit: "tbsp" },
      { id: "4", name: "Garlic", amount: "3", unit: "cloves" },
      { id: "5", name: "Ginger", amount: "1", unit: "tbsp" },
    ],
    nutrition: {
      calories: 285,
      protein: 35,
      carbs: 12,
      fat: 8,
      fiber: 4,
      sugar: 8,
    },
  },
  {
    id: "3",
    name: "Caesar Salad",
    cookTime: 10,
    servings: 2,
    category: "Salad",
    ingredients: [
      { id: "1", name: "Romaine Lettuce", amount: "1", unit: "head" },
      { id: "2", name: "Croutons", amount: "1", unit: "cup" },
      { id: "3", name: "Parmesan Cheese", amount: "0.5", unit: "cup" },
      { id: "4", name: "Caesar Dressing", amount: "2", unit: "tbsp" },
    ],
    nutrition: {
      calories: 200,
      protein: 10,
      carbs: 15,
      fat: 12,
      fiber: 5,
      sugar: 3,
    },
  },
  {
    id: "4",
    name: "Beef Tacos",
    cookTime: 25,
    servings: 4,
    category: "Mexican",
    ingredients: [
      { id: "1", name: "Ground Beef", amount: "500", unit: "g" },
      { id: "2", name: "Taco Shells", amount: "8", unit: "" },
      { id: "3", name: "Salsa", amount: "1", unit: "cup" },
      { id: "4", name: "Shredded Cheese", amount: "1", unit: "cup" },
      { id: "5", name: "Lettuce", amount: "1", unit: "cup" },
    ],
    nutrition: {
      calories: 350,
      protein: 25,
      carbs: 20,
      fat: 20,
      fiber: 5,
      sugar: 3,
    },
  },
  {
    id: "5",
    name: "Mushroom Risotto",
    cookTime: 35,
    servings: 4,
    category: "Italian",
    ingredients: [
      { id: "1", name: "Arborio Rice", amount: "300", unit: "g" },
      { id: "2", name: "Mushrooms", amount: "250", unit: "g" },
      { id: "3", name: "Vegetable Broth", amount: "1.5", unit: "L" },
      { id: "4", name: "Parmesan Cheese", amount: "50", unit: "g" },
      { id: "5", name: "Onion", amount: "1", unit: "" },
    ],
    nutrition: {
      calories: 400,
      protein: 15,
      carbs: 60,
      fat: 10,
      fiber: 5,
      sugar: 3,
    },
  },
  {
    id: "6",
    name: "Greek Salad",
    cookTime: 5,
    servings: 2,
    category: "Mediterranean",
    ingredients: [
      { id: "1", name: "Cucumber", amount: "1", unit: "" },
      { id: "2", name: "Tomato", amount: "2", unit: "" },
      { id: "3", name: "Red Onion", amount: "0.5", unit: "" },
      { id: "4", name: "Feta Cheese", amount: "100", unit: "g" },
      { id: "5", name: "Olives", amount: "50", unit: "g" },
    ],
    nutrition: {
      calories: 250,
      protein: 10,
      carbs: 15,
      fat: 18,
      fiber: 5,
      sugar: 3,
    },
  },
  {
    id: "7",
    name: "Beef Tacos",
    cookTime: 25,
    servings: 4,
    category: "Mexican",
    ingredients: [
      { id: "1", name: "Ground Beef", amount: "500", unit: "g" },
      { id: "2", name: "Taco Shells", amount: "8", unit: "" },
      { id: "3", name: "Lettuce", amount: "1", unit: "head" },
      { id: "4", name: "Tomato", amount: "2", unit: "" },
      { id: "5", name: "Cheese", amount: "100", unit: "g" },
    ],
    nutrition: {
      calories: 380,
      protein: 22,
      carbs: 28,
      fat: 20,
      fiber: 4,
      sugar: 3,
    },
  },
  {
    id: "8",
    name: "Mushroom Risotto",
    cookTime: 35,
    servings: 4,
    category: "Italian",
    ingredients: [
      { id: "1", name: "Arborio Rice", amount: "300", unit: "g" },
      { id: "2", name: "Mushrooms", amount: "400", unit: "g" },
      { id: "3", name: "Vegetable Stock", amount: "1", unit: "L" },
      { id: "4", name: "Parmesan", amount: "100", unit: "g" },
      { id: "5", name: "White Wine", amount: "125", unit: "ml" },
    ],
    nutrition: {
      calories: 420,
      protein: 16,
      carbs: 65,
      fat: 12,
      fiber: 3,
      sugar: 4,
    },
  },
  {
    id: "9",
    name: "Thai Green Curry",
    cookTime: 30,
    servings: 4,
    category: "Thai",
    ingredients: [
      { id: "1", name: "Chicken Breast", amount: "500", unit: "g" },
      { id: "2", name: "Green Curry Paste", amount: "3", unit: "tbsp" },
      { id: "3", name: "Coconut Milk", amount: "400", unit: "ml" },
      { id: "4", name: "Bell Peppers", amount: "2", unit: "" },
      { id: "5", name: "Thai Basil", amount: "1", unit: "bunch" },
    ],
    nutrition: {
      calories: 320,
      protein: 28,
      carbs: 12,
      fat: 18,
      fiber: 3,
      sugar: 8,
    },
  },
  {
    id: "10",
    name: "Caesar Salad",
    cookTime: 15,
    servings: 2,
    category: "Salad",
    ingredients: [
      { id: "1", name: "Romaine Lettuce", amount: "2", unit: "heads" },
      { id: "2", name: "Croutons", amount: "1", unit: "cup" },
      { id: "3", name: "Parmesan", amount: "50", unit: "g" },
      { id: "4", name: "Caesar Dressing", amount: "4", unit: "tbsp" },
      { id: "5", name: "Anchovies", amount: "4", unit: "fillets" },
    ],
    nutrition: {
      calories: 280,
      protein: 12,
      carbs: 18,
      fat: 20,
      fiber: 6,
      sugar: 4,
    },
  },
  {
    id: "11",
    name: "Beef Burger",
    cookTime: 20,
    servings: 4,
    category: "American",
    ingredients: [
      { id: "1", name: "Ground Beef", amount: "600", unit: "g" },
      { id: "2", name: "Burger Buns", amount: "4", unit: "" },
      { id: "3", name: "Cheese Slices", amount: "4", unit: "" },
      { id: "4", name: "Lettuce", amount: "4", unit: "leaves" },
      { id: "5", name: "Tomato", amount: "1", unit: "" },
    ],
    nutrition: {
      calories: 520,
      protein: 32,
      carbs: 35,
      fat: 28,
      fiber: 3,
      sugar: 5,
    },
  },
  {
    id: "12",
    name: "Vegetable Soup",
    cookTime: 40,
    servings: 6,
    category: "Soup",
    ingredients: [
      { id: "1", name: "Mixed Vegetables", amount: "500", unit: "g" },
      { id: "2", name: "Vegetable Stock", amount: "1.5", unit: "L" },
      { id: "3", name: "Onion", amount: "1", unit: "" },
      { id: "4", name: "Garlic", amount: "3", unit: "cloves" },
      { id: "5", name: "Herbs", amount: "2", unit: "tbsp" },
    ],
    nutrition: {
      calories: 120,
      protein: 4,
      carbs: 22,
      fat: 2,
      fiber: 6,
      sugar: 8,
    },
  },
  {
    id: "13",
    name: "Pancakes",
    cookTime: 20,
    servings: 4,
    category: "Breakfast",
    ingredients: [
      { id: "1", name: "Flour", amount: "200", unit: "g" },
      { id: "2", name: "Eggs", amount: "2", unit: "" },
      { id: "3", name: "Milk", amount: "300", unit: "ml" },
      { id: "4", name: "Butter", amount: "50", unit: "g" },
      { id: "5", name: "Maple Syrup", amount: "4", unit: "tbsp" },
    ],
    nutrition: {
      calories: 350,
      protein: 12,
      carbs: 45,
      fat: 14,
      fiber: 2,
      sugar: 15,
    },
  },
  {
    id: "14",
    name: "Fish and Chips",
    cookTime: 30,
    servings: 2,
    category: "British",
    ingredients: [
      { id: "1", name: "White Fish", amount: "400", unit: "g" },
      { id: "2", name: "Potatoes", amount: "500", unit: "g" },
      { id: "3", name: "Flour", amount: "100", unit: "g" },
      { id: "4", name: "Beer", amount: "200", unit: "ml" },
      { id: "5", name: "Oil", amount: "500", unit: "ml" },
    ],
    nutrition: {
      calories: 680,
      protein: 35,
      carbs: 55,
      fat: 35,
      fiber: 4,
      sugar: 3,
    },
  },
]

const RecipeCard = ({
  recipe,
  onDragStart,
  onEdit,
  onDelete,
  onMobileSelect,
  isSelected,
}: {
  recipe: Recipe
  onDragStart: (recipe: Recipe) => void
  onEdit: (recipe: Recipe) => void
  onDelete: (recipeId: string) => void
  onMobileSelect?: (recipe: Recipe) => void
  isSelected?: boolean
}) => {
  const [isExpanded, setIsExpanded] = useState(false)

  const handleCardClick = () => {
    if (onMobileSelect && window.innerWidth < 1024) {
      onMobileSelect(recipe)
    }
  }

  return (
    <Card
      className={`cursor-grab active:cursor-grabbing hover:shadow-md transition-shadow group lg:cursor-grab mb-2:${
        isSelected ? "ring-2 ring-primary bg-primary/5" : ""
      } ${onMobileSelect ? "lg:cursor-grab cursor-pointer lg:cursor-grab" : ""}`}
      onClick={handleCardClick}
    >
      <CardContent className="px-3">
        <div className="flex gap-2" draggable onDragStart={() => onDragStart(recipe)}>
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <h3 className="font-medium text-sm truncate">{recipe.name}</h3>
                <Badge variant="secondary" className="text-xs mb-2">
                  {recipe.category}
                </Badge>
              </div>
              <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-6 w-6 p-0"
                  onClick={(e) => {
                    e.stopPropagation()
                    setIsExpanded(!isExpanded)
                  }}
                >
                  <ChevronDown className={`w-3 h-3 transition-transform ${isExpanded ? "rotate-180" : ""}`} />
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-6 w-6 p-0"
                  onClick={(e) => {
                    e.stopPropagation()
                    onEdit(recipe)
                  }}
                >
                  <Edit className="w-3 h-3" />
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                  onClick={(e) => {
                    e.stopPropagation()
                    onDelete(recipe.id)
                  }}
                >
                  <Trash2 className="w-3 h-3" />
                </Button>
              </div>
            </div>
            <div className="flex items-center gap-2 sm:gap-3 text-xs text-muted-foreground mb-2 flex-wrap">
              <div className="flex items-center gap-1">
                <Clock className="w-3 h-3" />
                {recipe.cookTime}m
              </div>
              <div className="flex items-center gap-1">
                <Users className="w-3 h-3" />
                {recipe.servings}
              </div>
              <div className="flex items-center gap-1">
                <Zap className="w-3 h-3" />
                {recipe.nutrition.calories} cal
              </div>
            </div>
          </div>
        </div>

        {isExpanded && (
          <div className="mt-3 pt-3 border-t space-y-3" onClick={(e) => e.stopPropagation()}>
            <div>
              <h4 className="text-xs font-medium mb-2">Ingredients:</h4>
              <div className="space-y-1">
                {recipe.ingredients.map((ingredient) => (
                  <div key={ingredient.id} className="text-xs text-muted-foreground">
                    {ingredient.amount} {ingredient.unit} {ingredient.name}
                  </div>
                ))}
              </div>
            </div>

            <div>
              <h4 className="text-xs font-medium mb-2">Nutrition (per serving):</h4>
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 text-xs">
                <div>Protein: {recipe.nutrition.protein}g</div>
                <div>Carbs: {recipe.nutrition.carbs}g</div>
                <div>Fat: {recipe.nutrition.fat}g</div>
                <div>Fiber: {recipe.nutrition.fiber}g</div>
                <div>Sugar: {recipe.nutrition.sugar}g</div>
                <div>Calories: {recipe.nutrition.calories}</div>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

const CalendarDay = ({
  date,
  recipes,
  onDrop,
  onDragOver,
  onDragLeave,
  isDragOver,
  onRemoveRecipe,
  onClearDay,
  onMobileAdd,
  isMobileAddMode,
}: {
  date: string
  recipes: Recipe[]
  onDrop: (date: string) => void
  onDragOver: (e: React.DragEvent) => void
  onDragLeave: () => void
  isDragOver: boolean
  onRemoveRecipe: (date: string, recipeIndex: number) => void
  onClearDay: (date: string) => void
  onMobileAdd?: (date: string) => void
  isMobileAddMode?: boolean
}) => {
  const dayNumber = new Date(date).getDate()
  const isToday = false // Disable today highlighting to prevent hydration issues

  const handleDayClick = () => {
    if (onMobileAdd && isMobileAddMode && window.innerWidth < 1024) {
      onMobileAdd(date)
    }
  }

  return (
    <div
      className={`min-h-24 sm:min-h-32 p-2 border rounded-lg transition-colors ${
        isDragOver ? "border-primary bg-primary/5" : "border-border"
      } ${isToday ? "bg-primary/5 border-primary" : ""} ${
        isMobileAddMode ? "lg:cursor-auto cursor-pointer hover:bg-primary/10 lg:hover:bg-transparent" : ""
      }`}
      onDrop={() => onDrop(date)}
      onDragOver={onDragOver}
      onDragLeave={onDragLeave}
      onClick={handleDayClick}
    >
      <div className="flex items-center justify-between mb-2">
        <span className={`text-sm font-medium ${isToday ? "text-primary" : ""}`}>{dayNumber}</span>
        <div className="flex items-center gap-1">
          {isToday && (
            <Badge variant="default" className="text-xs">
              Today
            </Badge>
          )}
          {recipes.length > 0 && (
            <Button
              size="sm"
              variant="ghost"
              className="h-5 w-5 p-0 opacity-0 group-hover:opacity-100 transition-opacity text-muted-foreground hover:text-destructive"
              onClick={(e) => {
                e.stopPropagation()
                onClearDay(date)
              }}
              title="Clear all meals"
            >
              <X className="w-3 h-3" />
            </Button>
          )}
        </div>
      </div>
      <div className="space-y-1 group">
        {recipes.map((recipe, index) => (
          <div
            key={`${recipe.id}-${index}`}
            className="bg-secondary/50 rounded p-1 text-xs truncate group/recipe hover:bg-secondary/70 transition-colors relative"
            title={recipe.name}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1 flex-1 min-w-0">
                <ChefHat className="w-3 h-3 flex-shrink-0" />
                <span className="truncate">{recipe.name}</span>
              </div>
              <Button
                size="sm"
                variant="ghost"
                className="h-4 w-4 p-0 opacity-0 group-hover/recipe:opacity-100 transition-opacity text-muted-foreground hover:text-destructive flex-shrink-0"
                onClick={(e) => {
                  e.stopPropagation()
                  onRemoveRecipe(date, index)
                }}
                title="Remove meal"
              >
                <X className="w-2.5 h-2.5" />
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

const RecipeSidebar = ({
  recipes,
  onDragStart,
  onEdit,
  onDelete,
  onCreateRecipe,
  onMobileSelect,
  selectedRecipe,
  isMobileAddMode,
}: {
  recipes: Recipe[]
  onDragStart: (recipe: Recipe) => void
  onEdit: (recipe: Recipe) => void
  onDelete: (recipeId: string) => void
  onCreateRecipe: () => void
  onMobileSelect?: (recipe: Recipe) => void
  selectedRecipe?: Recipe | null
  isMobileAddMode?: boolean
}) => (
  <div className="max-h-screen flex flex-col">
    <div className="p-4 border-b">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <ChefHat className="w-5 h-5" />
            Recipes
          </h2>
          <p className="text-sm text-muted-foreground mt-1">
            <span className="lg:hidden">Tap recipe, then tap day to add</span>
            <span className="hidden lg:inline">Drag recipes to your calendar</span>
          </p>
          {isMobileAddMode && (
            <p className="text-xs text-primary mt-1">Now tap a day to add "{selectedRecipe?.name}"</p>
          )}
        </div>
        <Button size="sm" onClick={onCreateRecipe}>
          <Plus className="w-4 h-4" />
        </Button>
      </div>
    </div>
    <ScrollArea className="flex-1 overflow-scroll" style="scrollbar-width: none">
      <div className="p-4">
        {recipes.map((recipe) => (
          <RecipeCard
            key={recipe.id}
            recipe={recipe}
            onDragStart={onDragStart}
            onEdit={onEdit}
            onDelete={onDelete}
            onMobileSelect={onMobileSelect}
            isSelected={selectedRecipe?.id === recipe.id}
          />
        ))}
      </div>
    </ScrollArea>
  </div>
)

export default function MealPlanner() {
  const [mealPlan, setMealPlan] = useState<MealPlan>({})
  const [draggedRecipe, setDraggedRecipe] = useState<Recipe | null>(null)
  const [dragOverDate, setDragOverDate] = useState<string | null>(null)
  const [recipes, setRecipes] = useState<Recipe[]>(defaultRecipes)
  const [isRecipeDialogOpen, setIsRecipeDialogOpen] = useState(false)
  const [editingRecipe, setEditingRecipe] = useState<Recipe | null>(null)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [isShoppingListOpen, setIsShoppingListOpen] = useState(false)
  const [peopleCount, setPeopleCount] = useState(4)
  const [shoppingList, setShoppingList] = useState<ShoppingListItem[]>([])
  const [recipeForm, setRecipeForm] = useState({
    name: "",
    cookTime: "",
    servings: "",
    category: "",
    instructions: "",
    ingredients: [] as Ingredient[],
    nutrition: {
      calories: 0,
      protein: 0,
      carbs: 0,
      fat: 0,
      fiber: 0,
      sugar: 0,
    },
  })

  const [selectedWeekStart, setSelectedWeekStart] = useState(() => {
    const today = new Date()
    const startOfWeek = new Date(today)
    startOfWeek.setDate(today.getDate() - today.getDay())
    return startOfWeek
  })

  const today = new Date()
  const weekDates = Array.from({ length: 7 }, (_, i) => {
    const date = new Date(selectedWeekStart)
    date.setDate(selectedWeekStart.getDate() + i)
    return date.toISOString().split("T")[0]
  })

  const weekDays = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]

  const [selectedRecipeForMobile, setSelectedRecipeForMobile] = useState<Recipe | null>(null)
  const [isMobileAddMode, setIsMobileAddMode] = useState(false)

  const generateShoppingList = () => {
    const ingredientMap = new Map<string, { amount: number; unit: string; recipes: string[] }>()

    // Collect all recipes from the meal plan
    Object.values(mealPlan).forEach((dayRecipes) => {
      dayRecipes.forEach((recipe) => {
        const servingMultiplier = peopleCount / recipe.servings

        recipe.ingredients.forEach((ingredient) => {
          const key = `${ingredient.name.toLowerCase()}-${ingredient.unit}`
          const amount = Number.parseFloat(ingredient.amount) * servingMultiplier

          if (ingredientMap.has(key)) {
            const existing = ingredientMap.get(key)!
            existing.amount += amount
            if (!existing.recipes.includes(recipe.name)) {
              existing.recipes.push(recipe.name)
            }
          } else {
            ingredientMap.set(key, {
              amount,
              unit: ingredient.unit,
              recipes: [recipe.name],
            })
          }
        })
      })
    })

    // Convert to shopping list items
    const items: ShoppingListItem[] = Array.from(ingredientMap.entries()).map(([key, data]) => {
      const name = key.split("-")[0]
      return {
        id: key,
        name: name.charAt(0).toUpperCase() + name.slice(1),
        amount: Math.round(data.amount * 100) / 100, // Round to 2 decimal places
        unit: data.unit,
        recipes: data.recipes,
        checked: false,
      }
    })

    // Sort by name
    items.sort((a, b) => a.name.localeCompare(b.name))

    setShoppingList(items)
    setIsShoppingListOpen(true)
  }

  const toggleShoppingListItem = (itemId: string) => {
    setShoppingList((prev) => prev.map((item) => (item.id === itemId ? { ...item, checked: !item.checked } : item)))
  }

  const clearShoppingList = () => {
    setShoppingList([])
  }

  const exportShoppingList = () => {
    const text = shoppingList
      .map((item) => {
        const checkbox = item.checked ? "☑" : "☐"
        return `${checkbox} ${item.amount} ${item.unit} ${item.name}`
      })
      .join("\n")

    const blob = new Blob([text], { type: "text/plain" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `shopping-list-${new Date().toISOString().split("T")[0]}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  const handleDragStart = (recipe: Recipe) => {
    setDraggedRecipe(recipe)
  }

  const handleDragOver = (e: React.DragEvent, date: string) => {
    e.preventDefault()
    setDragOverDate(date)
  }

  const handleDragLeave = () => {
    setDragOverDate(null)
  }

  const handleDrop = (date: string) => {
    if (draggedRecipe) {
      setMealPlan((prev) => ({
        ...prev,
        [date]: [...(prev[date] || []), draggedRecipe],
      }))
    }
    setDraggedRecipe(null)
    setDragOverDate(null)
  }

  const handleRemoveRecipe = (date: string, recipeIndex: number) => {
    setMealPlan((prev) => ({
      ...prev,
      [date]: prev[date]?.filter((_, index) => index !== recipeIndex) || [],
    }))
  }

  const handleClearDay = (date: string) => {
    setMealPlan((prev) => ({
      ...prev,
      [date]: [],
    }))
  }

  const handleCreateRecipe = () => {
    setEditingRecipe(null)
    setRecipeForm({
      name: "",
      cookTime: "",
      servings: "",
      category: "",
      instructions: "",
      ingredients: [],
      nutrition: {
        calories: 0,
        protein: 0,
        carbs: 0,
        fat: 0,
        fiber: 0,
        sugar: 0,
      },
    })
    setIsRecipeDialogOpen(true)
    setIsMobileMenuOpen(false)
  }

  const handleEditRecipe = (recipe: Recipe) => {
    setEditingRecipe(recipe)
    setRecipeForm({
      name: recipe.name,
      cookTime: recipe.cookTime.toString(),
      servings: recipe.servings.toString(),
      category: recipe.category,
      instructions: recipe.instructions || "",
      ingredients: [...recipe.ingredients],
      nutrition: { ...recipe.nutrition },
    })
    setIsRecipeDialogOpen(true)
    setIsMobileMenuOpen(false)
  }

  const handleDeleteRecipe = (recipeId: string) => {
    setRecipes((prev) => prev.filter((recipe) => recipe.id !== recipeId))
    setMealPlan((prev) => {
      const updated = { ...prev }
      Object.keys(updated).forEach((date) => {
        updated[date] = updated[date].filter((recipe) => recipe.id !== recipeId)
      })
      return updated
    })
  }

  const handleSaveRecipe = () => {
    if (!recipeForm.name || !recipeForm.cookTime || !recipeForm.servings || !recipeForm.category) {
      return
    }

    const newRecipe: Recipe = {
      id: editingRecipe?.id || Date.now().toString(),
      name: recipeForm.name,
      cookTime: Number.parseInt(recipeForm.cookTime),
      servings: Number.parseInt(recipeForm.servings),
      category: recipeForm.category,
      ingredients: recipeForm.ingredients.filter((ing) => ing.name && ing.amount),
      nutrition: recipeForm.nutrition,
      instructions: recipeForm.instructions,
    }

    if (editingRecipe) {
      setRecipes((prev) => prev.map((recipe) => (recipe.id === editingRecipe.id ? newRecipe : recipe)))
      setMealPlan((prev) => {
        const updated = { ...prev }
        Object.keys(updated).forEach((date) => {
          updated[date] = updated[date].map((recipe) => (recipe.id === editingRecipe.id ? newRecipe : recipe))
        })
        return updated
      })
    } else {
      setRecipes((prev) => [...prev, newRecipe])
    }

    setIsRecipeDialogOpen(false)
    setEditingRecipe(null)
  }

  const handlePreviousWeek = () => {
    setSelectedWeekStart((prev) => {
      const newDate = new Date(prev)
      newDate.setDate(prev.getDate() - 7)
      return newDate
    })
  }

  const handleNextWeek = () => {
    setSelectedWeekStart((prev) => {
      const newDate = new Date(prev)
      newDate.setDate(prev.getDate() + 7)
      return newDate
    })
  }

  const handleGoToToday = () => {
    const today = new Date()
    const startOfWeek = new Date(today)
    startOfWeek.setDate(today.getDate() - today.getDay())
    setSelectedWeekStart(startOfWeek)
  }

  const isCurrentWeek = () => {
    return false // Disable current week check to prevent hydration issues
  }

  const handleMobileRecipeSelect = (recipe: Recipe) => {
    if (selectedRecipeForMobile?.id === recipe.id) {
      // Deselect if clicking the same recipe
      setSelectedRecipeForMobile(null)
      setIsMobileAddMode(false)
    } else {
      // Select new recipe and enter add mode
      setSelectedRecipeForMobile(recipe)
      setIsMobileAddMode(true)
      setIsMobileMenuOpen(false) // Close mobile menu
    }
  }

  const handleMobileAddToDay = (date: string) => {
    if (selectedRecipeForMobile) {
      setMealPlan((prev) => ({
        ...prev,
        [date]: [...(prev[date] || []), selectedRecipeForMobile],
      }))
      setSelectedRecipeForMobile(null)
      setIsMobileAddMode(false)
    }
  }

  const hasPlannedMeals = Object.values(mealPlan).some((dayRecipes) => dayRecipes.length > 0)

  return (
    <div className="flex flex-col lg:flex-row h-screen bg-background">
      {/* Desktop Sidebar */}
      <div className="hidden lg:block w-80 border-r bg-muted/30">
        <RecipeSidebar
          recipes={recipes}
          onDragStart={handleDragStart}
          onEdit={handleEditRecipe}
          onDelete={handleDeleteRecipe}
          onCreateRecipe={handleCreateRecipe}
          onMobileSelect={handleMobileRecipeSelect}
          selectedRecipe={selectedRecipeForMobile}
          isMobileAddMode={isMobileAddMode}
        />
      </div>

      {/* Mobile Sheet for Recipes */}
      <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
        <SheetContent side="left" className="w-80 p-0">
          <SheetHeader className="p-4 border-b">
            <SheetTitle className="flex items-center gap-2">
              <ChefHat className="w-5 h-5" />
              Recipes
            </SheetTitle>
          </SheetHeader>
            <RecipeSidebar
              recipes={recipes}
              onDragStart={handleDragStart}
              onEdit={handleEditRecipe}
              onDelete={handleDeleteRecipe}
              onCreateRecipe={handleCreateRecipe}
              onMobileSelect={handleMobileRecipeSelect}
              selectedRecipe={selectedRecipeForMobile}
              isMobileAddMode={isMobileAddMode}
            />
        </SheetContent>
      </Sheet>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-h-0">
        <div className="p-4 border-b">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              {/* Mobile Menu Button */}
              <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
                <SheetTrigger asChild>
                  <Button variant="outline" size="sm" className="lg:hidden bg-transparent">
                    <Menu className="w-4 h-4" />
                  </Button>
                </SheetTrigger>
              </Sheet>

              <div>
                <h1 className="text-xl sm:text-2xl font-bold flex items-center gap-2">
                  <Calendar className="w-5 h-5 sm:w-6 sm:h-6" />
                  Meal Planner
                </h1>
                <p className="text-sm text-muted-foreground">
                  <span suppressHydrationWarning>
                    Week of{" "}
                    {selectedWeekStart.toLocaleDateString("en-US", {
                      month: "long",
                      day: "numeric",  
                      year: "numeric",
                    })}
                  </span>
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={generateShoppingList}
                disabled={!hasPlannedMeals}
                className="hidden sm:flex bg-transparent"
              >
                <ShoppingCart className="w-4 h-4 mr-2" />
                Shopping List
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={generateShoppingList}
                disabled={!hasPlannedMeals}
                className="sm:hidden bg-transparent"
              >
                <ShoppingCart className="w-4 h-4" />
              </Button>
              <ThemeToggle />
              <div className="flex items-center gap-1">
                <Button variant="outline" size="sm" onClick={handlePreviousWeek}>
                  <ChevronLeft className="w-4 h-4" />
                </Button>
                <Button
                  variant={isCurrentWeek() ? "default" : "outline"}
                  size="sm"
                  onClick={handleGoToToday}
                  disabled={isCurrentWeek()}
                  className="hidden sm:flex"
                >
                  Today
                </Button>
                <Button variant="outline" size="sm" onClick={handleNextWeek}>
                  <ChevronRight className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <div className="flex-1 p-2 sm:p-4 overflow-auto">
          {/* Mobile: Stack days vertically, Desktop: 7-column grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-7 gap-2 sm:gap-4">
            {weekDates.map((date, index) => (
              <div key={date} className="flex flex-col">
                <div className="text-center mb-2">
                  <div className="font-medium text-sm">{weekDays[index]}</div>
                  <div className="text-xs text-muted-foreground" suppressHydrationWarning>
                    {new Date(date).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                  </div>
                </div>
                <CalendarDay
                  date={date}
                  recipes={mealPlan[date] || []}
                  onDrop={handleDrop}
                  onDragOver={(e) => handleDragOver(e, date)}
                  onDragLeave={handleDragLeave}
                  isDragOver={dragOverDate === date}
                  onRemoveRecipe={handleRemoveRecipe}
                  onClearDay={handleClearDay}
                  onMobileAdd={handleMobileAddToDay}
                  isMobileAddMode={isMobileAddMode}
                />
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Shopping List Dialog */}
      <Dialog open={isShoppingListOpen} onOpenChange={setIsShoppingListOpen}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <ShoppingCart className="w-5 h-5" />
              Weekly Shopping List
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-6">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <Label htmlFor="people-count">People eating:</Label>
                <Input
                  id="people-count"
                  type="number"
                  min="1"
                  max="20"
                  value={peopleCount}
                  onChange={(e) => setPeopleCount(Number(e.target.value) || 1)}
                  className="w-20"
                />
              </div>
              <Button onClick={generateShoppingList} size="sm">
                Recalculate
              </Button>
            </div>

            {shoppingList.length > 0 && (
              <>
                <div className="flex items-center justify-between">
                  <p className="text-sm text-muted-foreground">
                    {shoppingList.filter((item) => item.checked).length} of {shoppingList.length} items checked
                  </p>
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" onClick={exportShoppingList}>
                      Export List
                    </Button>
                    <Button variant="outline" size="sm" onClick={clearShoppingList}>
                      Clear All
                    </Button>
                  </div>
                </div>

                <ScrollArea className="h-96">
                  <div className="space-y-2">
                    {shoppingList.map((item) => (
                      <div
                        key={item.id}
                        className={`flex items-center gap-3 p-3 rounded-lg border ${
                          item.checked ? "bg-muted/50 text-muted-foreground" : "bg-background"
                        }`}
                      >
                        <Checkbox checked={item.checked} onCheckedChange={() => toggleShoppingListItem(item.id)} />
                        <div className="flex-1">
                          <div className={`font-medium ${item.checked ? "line-through" : ""}`}>
                            {item.amount} {item.unit} {item.name}
                          </div>
                          <div className="text-xs text-muted-foreground">Used in: {item.recipes.join(", ")}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              </>
            )}

            {shoppingList.length === 0 && (
              <div className="text-center py-8 text-muted-foreground">
                <ShoppingCart className="w-12 h-12 mx-auto mb-4 opacity-50" />
                <p>No ingredients found. Add some recipes to your meal plan first!</p>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {/* Recipe Dialog */}
      <Dialog open={isRecipeDialogOpen} onOpenChange={setIsRecipeDialogOpen}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{editingRecipe ? "Edit Recipe" : "Create New Recipe"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-6">
            <div className="space-y-4">
              <h3 className="text-sm font-medium">Basic Information</h3>
              <div>
                <Label htmlFor="name">Recipe Name</Label>
                <Input
                  id="name"
                  value={recipeForm.name}
                  onChange={(e) => setRecipeForm((prev) => ({ ...prev, name: e.target.value }))}
                  placeholder="Enter recipe name"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="cookTime">Cook Time (minutes)</Label>
                  <Input
                    id="cookTime"
                    type="number"
                    value={recipeForm.cookTime}
                    onChange={(e) => setRecipeForm((prev) => ({ ...prev, cookTime: e.target.value }))}
                    placeholder="30"
                  />
                </div>
                <div>
                  <Label htmlFor="servings">Servings</Label>
                  <Input
                    id="servings"
                    type="number"
                    value={recipeForm.servings}
                    onChange={(e) => setRecipeForm((prev) => ({ ...prev, servings: e.target.value }))}
                    placeholder="4"
                  />
                </div>
              </div>
              <div>
                <Label htmlFor="category">Category</Label>
                <Select
                  value={recipeForm.category}
                  onValueChange={(value) => setRecipeForm((prev) => ({ ...prev, category: value }))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Italian">Italian</SelectItem>
                    <SelectItem value="Asian">Asian</SelectItem>
                    <SelectItem value="Mexican">Mexican</SelectItem>
                    <SelectItem value="Mediterranean">Mediterranean</SelectItem>
                    <SelectItem value="American">American</SelectItem>
                    <SelectItem value="Indian">Indian</SelectItem>
                    <SelectItem value="Salad">Salad</SelectItem>
                    <SelectItem value="Dessert">Dessert</SelectItem>
                    <SelectItem value="Other">Other</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium">Ingredients</h3>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    const newIngredient: Ingredient = {
                      id: Date.now().toString(),
                      name: "",
                      amount: "",
                      unit: "",
                    }
                    setRecipeForm((prev) => ({
                      ...prev,
                      ingredients: [...prev.ingredients, newIngredient],
                    }))
                  }}
                >
                  <Plus className="w-4 h-4 mr-1" />
                  Add Ingredient
                </Button>
              </div>
              <div className="space-y-2 max-h-40 overflow-y-auto">
                {recipeForm.ingredients.map((ingredient, index) => (
                  <div key={ingredient.id} className="flex gap-2 items-center">
                    <Input
                      placeholder="Amount"
                      value={ingredient.amount}
                      onChange={(e) => {
                        const updated = [...recipeForm.ingredients]
                        updated[index] = { ...updated[index], amount: e.target.value }
                        setRecipeForm((prev) => ({ ...prev, ingredients: updated }))
                      }}
                      className="w-20"
                    />
                    <Input
                      placeholder="Unit"
                      value={ingredient.unit}
                      onChange={(e) => {
                        const updated = [...recipeForm.ingredients]
                        updated[index] = { ...updated[index], unit: e.target.value }
                        setRecipeForm((prev) => ({ ...prev, ingredients: updated }))
                      }}
                      className="w-20"
                    />
                    <Input
                      placeholder="Ingredient name"
                      value={ingredient.name}
                      onChange={(e) => {
                        const updated = [...recipeForm.ingredients]
                        updated[index] = { ...updated[index], name: e.target.value }
                        setRecipeForm((prev) => ({ ...prev, ingredients: updated }))
                      }}
                      className="flex-1"
                    />
                    <Button
                      type="button"
                      size="sm"
                      variant="ghost"
                      onClick={() => {
                        setRecipeForm((prev) => ({
                          ...prev,
                          ingredients: prev.ingredients.filter((_, i) => i !== index),
                        }))
                      }}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-4">
              <h3 className="text-sm font-medium">Nutritional Information (per serving)</h3>
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="calories">Calories</Label>
                  <Input
                    id="calories"
                    type="number"
                    value={recipeForm.nutrition.calories}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, calories: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div>
                  <Label htmlFor="protein">Protein (g)</Label>
                  <Input
                    id="protein"
                    type="number"
                    value={recipeForm.nutrition.protein}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, protein: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div>
                  <Label htmlFor="carbs">Carbs (g)</Label>
                  <Input
                    id="carbs"
                    type="number"
                    value={recipeForm.nutrition.carbs}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, carbs: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div>
                  <Label htmlFor="fat">Fat (g)</Label>
                  <Input
                    id="fat"
                    type="number"
                    value={recipeForm.nutrition.fat}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, fat: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div>
                  <Label htmlFor="fiber">Fiber (g)</Label>
                  <Input
                    id="fiber"
                    type="number"
                    value={recipeForm.nutrition.fiber}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, fiber: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
                <div>
                  <Label htmlFor="sugar">Sugar (g)</Label>
                  <Input
                    id="sugar"
                    type="number"
                    value={recipeForm.nutrition.sugar}
                    onChange={(e) =>
                      setRecipeForm((prev) => ({
                        ...prev,
                        nutrition: { ...prev.nutrition, sugar: Number(e.target.value) || 0 },
                      }))
                    }
                    placeholder="0"
                  />
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <h3 className="text-sm font-medium">Instructions (optional)</h3>
              <Textarea
                value={recipeForm.instructions}
                onChange={(e) => setRecipeForm((prev) => ({ ...prev, instructions: e.target.value }))}
                placeholder="Enter cooking instructions..."
                rows={4}
              />
            </div>

            <div className="flex gap-2 justify-end">
              <Button variant="outline" onClick={() => setIsRecipeDialogOpen(false)}>
                <X className="w-4 h-4 mr-2" />
                Cancel
              </Button>
              <Button onClick={handleSaveRecipe}>
                <Save className="w-4 h-4 mr-2" />
                {editingRecipe ? "Update" : "Create"}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
