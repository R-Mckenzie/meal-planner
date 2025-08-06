"use client"

import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { ThemeToggle } from "@/components/theme-toggle"
import { 
  Calendar, 
  ChefHat, 
  Clock, 
  Users, 
  ShoppingCart, 
  Star
} from "lucide-react"

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <ChefHat className="w-8 h-8" />
            <h1 className="text-2xl font-bold">Meal Planner</h1>
          </div>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <Link href="/auth/signin">
              <Button variant="ghost">Sign In</Button>
            </Link>
            <Link href="/auth/signup">
              <Button>Get Started</Button>
            </Link>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="py-20 px-4">
        <div className="container mx-auto text-center max-w-4xl">
          <h2 className="text-5xl font-bold mb-6">
            Plan Your Meals, 
            <span className="text-primary"> Simplify Your Life</span>
          </h2>
          <p className="text-xl text-muted-foreground mb-8">
            Organize your weekly meals, manage recipes, and generate shopping lists 
            automatically. Make meal planning effortless and enjoy more time with family.
          </p>
          <div className="flex items-center justify-center gap-4">
            <Link href="/auth/signup">
              <Button size="lg" className="text-lg px-8">
                Start Planning Free
              </Button>
            </Link>
            <Link href="/auth/signin">
              <Button size="lg" variant="outline" className="text-lg px-8">
                Sign In
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 px-4 bg-muted/30">
        <div className="container mx-auto max-w-6xl">
          <h3 className="text-3xl font-bold text-center mb-12">
            Everything You Need for Meal Planning
          </h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card>
              <CardContent className="p-6 text-center">
                <Calendar className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Weekly Planning</h4>
                <p className="text-muted-foreground">
                  Drag and drop recipes onto your calendar. Plan breakfast, lunch, 
                  dinner, and snacks for the entire week.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6 text-center">
                <ChefHat className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Recipe Management</h4>
                <p className="text-muted-foreground">
                  Store all your favorite recipes with ingredients, cooking times, 
                  and nutritional information in one place.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6 text-center">
                <ShoppingCart className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Smart Shopping Lists</h4>
                <p className="text-muted-foreground">
                  Automatically generate shopping lists from your meal plans. 
                  Never forget an ingredient again.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6 text-center">
                <Clock className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Save Time</h4>
                <p className="text-muted-foreground">
                  Spend less time thinking about what to cook and more time 
                  enjoying delicious meals with loved ones.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6 text-center">
                <Users className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Family Friendly</h4>
                <p className="text-muted-foreground">
                  Scale recipes for any number of people. Perfect for families, 
                  couples, or meal prep enthusiasts.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6 text-center">
                <Star className="w-12 h-12 mx-auto mb-4 text-primary" />
                <h4 className="text-xl font-semibold mb-2">Nutrition Tracking</h4>
                <p className="text-muted-foreground">
                  Keep track of calories, protein, carbs, and other nutrients 
                  to maintain a healthy diet.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-20 px-4">
        <div className="container mx-auto max-w-4xl">
          <h3 className="text-3xl font-bold text-center mb-12">
            How It Works
          </h3>
          <div className="space-y-12">
            <div className="flex items-center gap-8">
              <div className="flex-shrink-0 w-12 h-12 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-xl font-bold">
                1
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Add Your Recipes</h4>
                <p className="text-muted-foreground">
                  Start by adding your favorite recipes with ingredients, cooking times, 
                  and nutritional information.
                </p>
              </div>
            </div>

            <div className="flex items-center gap-8">
              <div className="flex-shrink-0 w-12 h-12 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-xl font-bold">
                2
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Plan Your Week</h4>
                <p className="text-muted-foreground">
                  Drag recipes from your collection onto your weekly calendar. 
                  Plan all meals and snacks ahead of time.
                </p>
              </div>
            </div>

            <div className="flex items-center gap-8">
              <div className="flex-shrink-0 w-12 h-12 bg-primary text-primary-foreground rounded-full flex items-center justify-center text-xl font-bold">
                3
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Generate Shopping List</h4>
                <p className="text-muted-foreground">
                  Automatically create a shopping list with all ingredients needed 
                  for your planned meals. Check items off as you shop.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 bg-primary/5">
        <div className="container mx-auto text-center max-w-3xl">
          <h3 className="text-3xl font-bold mb-4">
            Ready to Transform Your Meal Planning?
          </h3>
          <p className="text-xl text-muted-foreground mb-8">
            Join thousands of families who have simplified their meal planning 
            and reduced food waste with our intuitive platform.
          </p>
          <Link href="/auth/signup">
            <Button size="lg" className="text-lg px-8">
              Get Started Free Today
            </Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-8 px-4">
        <div className="container mx-auto text-center">
          <div className="flex items-center justify-center gap-2 mb-4">
            <ChefHat className="w-6 h-6" />
            <span className="text-lg font-semibold">Meal Planner</span>
          </div>
          <p className="text-muted-foreground">
            © 2024 Meal Planner. Simplifying meal planning for families everywhere.
          </p>
        </div>
      </footer>
    </div>
  )
}