# Tailwind CSS

Utility-first CSS patterns and common recipes.

## Layout

### Flexbox

```html
<!-- Center everything -->
<div class="flex items-center justify-center">

<!-- Space between items -->
<div class="flex justify-between items-center">

<!-- Vertical stack -->
<div class="flex flex-col gap-4">

<!-- Wrap items -->
<div class="flex flex-wrap gap-2">

<!-- Grow/shrink -->
<div class="flex">
  <div class="flex-1">Takes remaining space</div>
  <div class="flex-none">Fixed width</div>
</div>
```

### Grid

```html
<!-- Basic grid -->
<div class="grid grid-cols-3 gap-4">

<!-- Responsive grid -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

<!-- Auto-fit columns -->
<div class="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">

<!-- Span columns -->
<div class="col-span-2">
```

### Container

```html
<div class="container mx-auto px-4">
<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
```

## Spacing

```html
<!-- Margin -->
<div class="m-4">     <!-- all sides -->
<div class="mx-4">    <!-- horizontal -->
<div class="my-4">    <!-- vertical -->
<div class="mt-4">    <!-- top only -->

<!-- Padding -->
<div class="p-4">     <!-- all sides -->
<div class="px-4">    <!-- horizontal -->
<div class="py-4">    <!-- vertical -->

<!-- Gap (for flex/grid) -->
<div class="gap-4">   <!-- all -->
<div class="gap-x-4"> <!-- horizontal -->
<div class="gap-y-4"> <!-- vertical -->
```

## Typography

```html
<!-- Size -->
<p class="text-sm">Small</p>
<p class="text-base">Base (16px)</p>
<p class="text-lg">Large</p>
<p class="text-2xl">2XL</p>

<!-- Weight -->
<p class="font-medium">Medium</p>
<p class="font-semibold">Semibold</p>
<p class="font-bold">Bold</p>

<!-- Color -->
<p class="text-gray-600">Muted</p>
<p class="text-gray-900">Dark</p>

<!-- Line height -->
<p class="leading-tight">Tight</p>
<p class="leading-relaxed">Relaxed</p>

<!-- Truncate -->
<p class="truncate">Long text that gets truncated...</p>
<p class="line-clamp-2">Multi-line truncation...</p>
```

## Colors

```html
<!-- Background -->
<div class="bg-white">
<div class="bg-gray-100">
<div class="bg-blue-500">

<!-- Text -->
<p class="text-gray-900">
<p class="text-blue-600">

<!-- Border -->
<div class="border border-gray-200">

<!-- Hover states -->
<button class="bg-blue-500 hover:bg-blue-600">

<!-- Opacity -->
<div class="bg-black/50">  <!-- 50% opacity -->
```

## Components

### Button

```html
<button class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed">
  Click me
</button>
```

### Card

```html
<div class="bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-2">Title</h2>
  <p class="text-gray-600">Content</p>
</div>
```

### Input

```html
<input 
  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
  placeholder="Enter text..."
/>
```

### Badge

```html
<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
  Active
</span>
```

### Avatar

```html
<img class="h-10 w-10 rounded-full object-cover" src="..." alt="" />
```

## Responsive Design

```html
<!-- Mobile-first breakpoints -->
<div class="text-sm md:text-base lg:text-lg">
<!-- sm: 640px, md: 768px, lg: 1024px, xl: 1280px, 2xl: 1536px -->

<!-- Hide/show at breakpoints -->
<div class="hidden md:block">Desktop only</div>
<div class="md:hidden">Mobile only</div>

<!-- Responsive grid -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
```

## Dark Mode

```html
<!-- With class strategy -->
<div class="bg-white dark:bg-gray-900">
<p class="text-gray-900 dark:text-gray-100">

<!-- Toggle dark mode -->
document.documentElement.classList.toggle('dark');
```

## Animations

```html
<!-- Transitions -->
<button class="transition-colors duration-200 hover:bg-blue-600">

<!-- Transform -->
<div class="transform hover:scale-105 transition-transform">

<!-- Animate -->
<div class="animate-spin">Loading...</div>
<div class="animate-pulse">Skeleton</div>
<div class="animate-bounce">Attention</div>
```

## Common Patterns

### Centered Page

```html
<div class="min-h-screen flex items-center justify-center bg-gray-100">
  <div class="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
    <!-- content -->
  </div>
</div>
```

### Sticky Header

```html
<header class="sticky top-0 z-50 bg-white/80 backdrop-blur-sm border-b">
```

### Aspect Ratio

```html
<div class="aspect-video">16:9 container</div>
<div class="aspect-square">1:1 container</div>
```

### Divider

```html
<hr class="border-t border-gray-200 my-4" />
```

## Tips

1. **Use arbitrary values** when needed: `w-[347px]`
2. **Group hover states**: `group` on parent, `group-hover:` on child
3. **Use @apply sparingly** - prefer utility classes
4. **Install Tailwind CSS IntelliSense** VS Code extension
