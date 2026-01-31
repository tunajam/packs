# shadcn/ui

Copy-paste component patterns built on Radix UI and Tailwind CSS.

## Installation

```bash
npx shadcn@latest init
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add dialog
```

## Component Patterns

### Button Variants

```tsx
import { Button } from "@/components/ui/button"

<Button>Default</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>
<Button variant="destructive">Delete</Button>

<Button size="sm">Small</Button>
<Button size="lg">Large</Button>
<Button size="icon"><IconPlus /></Button>

<Button disabled>Disabled</Button>
<Button asChild>
  <Link href="/login">Login</Link>
</Button>
```

### Card Layout

```tsx
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

<Card>
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>Card description text</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Main content goes here</p>
  </CardContent>
  <CardFooter className="flex justify-between">
    <Button variant="outline">Cancel</Button>
    <Button>Save</Button>
  </CardFooter>
</Card>
```

### Dialog/Modal

```tsx
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

<Dialog>
  <DialogTrigger asChild>
    <Button>Open Dialog</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Are you sure?</DialogTitle>
      <DialogDescription>
        This action cannot be undone.
      </DialogDescription>
    </DialogHeader>
    <DialogFooter>
      <Button variant="outline">Cancel</Button>
      <Button>Confirm</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

### Controlled Dialog

```tsx
const [open, setOpen] = useState(false);

<Dialog open={open} onOpenChange={setOpen}>
  <DialogTrigger asChild>
    <Button>Open</Button>
  </DialogTrigger>
  <DialogContent>
    <Button onClick={() => setOpen(false)}>Close</Button>
  </DialogContent>
</Dialog>
```

### Form with Validation

```tsx
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
})

function LoginForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { email: "", password: "" },
  })

  function onSubmit(values: z.infer<typeof formSchema>) {
    console.log(values)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="email@example.com" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}
```

### Data Table

```tsx
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Email</TableHead>
      <TableHead className="text-right">Amount</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {users.map((user) => (
      <TableRow key={user.id}>
        <TableCell className="font-medium">{user.name}</TableCell>
        <TableCell>{user.email}</TableCell>
        <TableCell className="text-right">${user.amount}</TableCell>
      </TableRow>
    ))}
  </TableBody>
</Table>
```

### Dropdown Menu

```tsx
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

<DropdownMenu>
  <DropdownMenuTrigger asChild>
    <Button variant="outline">Options</Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent>
    <DropdownMenuLabel>My Account</DropdownMenuLabel>
    <DropdownMenuSeparator />
    <DropdownMenuItem>Profile</DropdownMenuItem>
    <DropdownMenuItem>Settings</DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem className="text-red-600">
      Delete Account
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

### Toast Notifications

```tsx
import { useToast } from "@/hooks/use-toast"
import { Toaster } from "@/components/ui/toaster"

// In layout
<Toaster />

// In component
const { toast } = useToast()

toast({
  title: "Success",
  description: "Your changes have been saved.",
})

toast({
  variant: "destructive",
  title: "Error",
  description: "Something went wrong.",
})
```

## Customization

### Extending Variants

```tsx
// components/ui/button.tsx
const buttonVariants = cva(
  "...",
  {
    variants: {
      variant: {
        // Add custom variant
        brand: "bg-purple-600 text-white hover:bg-purple-700",
      },
    },
  }
)
```

### Custom Colors

```css
/* globals.css */
@layer base {
  :root {
    --primary: 240 5.9% 10%;
    --primary-foreground: 0 0% 98%;
    /* Add custom colors */
    --brand: 262 83% 58%;
  }
}
```

## Tips

1. **Use `asChild`** to pass props to child elements
2. **Compose components** rather than customizing heavily
3. **Check Radix docs** for advanced behavior options
4. **Use `cn()`** helper for conditional classes
