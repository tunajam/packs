# Explain Code

Explain the following code clearly and concisely.

## Instructions

1. **Start with the big picture** - What does this code do overall? One sentence.

2. **Break down the structure** - What are the main parts/functions/classes?

3. **Explain the flow** - How does data/control move through the code?

4. **Highlight key concepts** - What patterns, algorithms, or techniques are used?

5. **Note any gotchas** - Are there non-obvious behaviors or edge cases?

## Output Format

```markdown
## Overview
[One sentence summary]

## How It Works
[Step-by-step explanation]

## Key Concepts
- [Concept 1]: [Brief explanation]
- [Concept 2]: [Brief explanation]

## Example Usage
[If applicable, show how to use this code]
```

## Guidelines

- Use plain language, avoid jargon unless necessary
- Relate concepts to familiar ideas when possible
- If the code has issues, mention them constructively
- Adjust depth based on the asker's apparent level
- Use analogies for complex concepts

## Example

**Input:**
```javascript
const debounce = (fn, delay) => {
  let timeoutId;
  return (...args) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
};
```

**Output:**

## Overview
This creates a debounced version of a function that only executes after the caller stops calling it for a specified time.

## How It Works
1. `debounce` takes a function and delay time
2. It returns a new function that wraps the original
3. Each time you call the wrapper, it:
   - Cancels any pending execution
   - Schedules a new execution after `delay` ms
4. The original function only runs if no new calls come in during the delay

## Key Concepts
- **Closure**: The inner function "remembers" `timeoutId` between calls
- **Rest/spread**: `...args` captures and forwards all arguments

## Example Usage
```javascript
const handleSearch = debounce((query) => {
  fetchResults(query);
}, 300);

// In a search input
input.addEventListener('input', (e) => handleSearch(e.target.value));
// Only searches 300ms after user stops typing
```
