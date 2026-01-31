# React Patterns

Modern React patterns and best practices for building maintainable applications.

## Component Patterns

### Composition over Props

```tsx
// ❌ Prop drilling
<Card
  title="User"
  subtitle="Details"
  action={<Button>Edit</Button>}
  footer={<Link>View</Link>}
/>

// ✅ Composition
<Card>
  <Card.Header>
    <Card.Title>User</Card.Title>
    <Card.Subtitle>Details</Card.Subtitle>
  </Card.Header>
  <Card.Body>...</Card.Body>
  <Card.Footer>
    <Button>Edit</Button>
  </Card.Footer>
</Card>
```

### Render Props

```tsx
function MouseTracker({ children }) {
  const [position, setPosition] = useState({ x: 0, y: 0 });
  
  return (
    <div onMouseMove={(e) => setPosition({ x: e.clientX, y: e.clientY })}>
      {children(position)}
    </div>
  );
}

// Usage
<MouseTracker>
  {({ x, y }) => <p>Mouse at {x}, {y}</p>}
</MouseTracker>
```

### Compound Components

```tsx
const Toggle = ({ children }) => {
  const [on, setOn] = useState(false);
  return (
    <ToggleContext.Provider value={{ on, toggle: () => setOn(!on) }}>
      {children}
    </ToggleContext.Provider>
  );
};

Toggle.Button = () => {
  const { on, toggle } = useContext(ToggleContext);
  return <button onClick={toggle}>{on ? 'ON' : 'OFF'}</button>;
};

Toggle.Status = () => {
  const { on } = useContext(ToggleContext);
  return <span>{on ? 'Active' : 'Inactive'}</span>;
};
```

## Hook Patterns

### Custom Hooks

Extract reusable logic:

```tsx
function useLocalStorage<T>(key: string, initial: T) {
  const [value, setValue] = useState<T>(() => {
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) : initial;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue] as const;
}
```

### useCallback for Callbacks

```tsx
// ❌ New function every render
<Button onClick={() => handleClick(id)} />

// ✅ Stable reference
const handleClick = useCallback(() => {
  doSomething(id);
}, [id]);
<Button onClick={handleClick} />
```

### useMemo for Expensive Computations

```tsx
const sortedItems = useMemo(() => {
  return items.slice().sort((a, b) => a.name.localeCompare(b.name));
}, [items]);
```

## State Patterns

### Derived State

```tsx
// ❌ Synchronized state
const [items, setItems] = useState([]);
const [count, setCount] = useState(0);

// ✅ Derived value
const [items, setItems] = useState([]);
const count = items.length;
```

### Reducer for Complex State

```tsx
type State = { items: Item[]; loading: boolean; error: string | null };
type Action = 
  | { type: 'FETCH_START' }
  | { type: 'FETCH_SUCCESS'; items: Item[] }
  | { type: 'FETCH_ERROR'; error: string };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'FETCH_START':
      return { ...state, loading: true, error: null };
    case 'FETCH_SUCCESS':
      return { ...state, loading: false, items: action.items };
    case 'FETCH_ERROR':
      return { ...state, loading: false, error: action.error };
  }
}
```

### Context for Shared State

```tsx
const ThemeContext = createContext<Theme | null>(null);

function useTheme() {
  const theme = useContext(ThemeContext);
  if (!theme) throw new Error('useTheme must be within ThemeProvider');
  return theme;
}
```

## Performance Patterns

### React.memo for Pure Components

```tsx
const ExpensiveList = memo(function ExpensiveList({ items }) {
  return items.map(item => <Item key={item.id} {...item} />);
});
```

### Virtualization for Long Lists

```tsx
import { useVirtualizer } from '@tanstack/react-virtual';

function VirtualList({ items }) {
  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50,
  });
  
  return (
    <div ref={parentRef} style={{ height: 400, overflow: 'auto' }}>
      <div style={{ height: virtualizer.getTotalSize() }}>
        {virtualizer.getVirtualItems().map(virtualRow => (
          <div key={virtualRow.key} style={{
            position: 'absolute',
            top: virtualRow.start,
            height: virtualRow.size,
          }}>
            {items[virtualRow.index].name}
          </div>
        ))}
      </div>
    </div>
  );
}
```

## Error Boundaries

```tsx
class ErrorBoundary extends Component {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error, info) {
    logError(error, info);
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback />;
    }
    return this.props.children;
  }
}
```

## Form Patterns

### Controlled Components

```tsx
function Form() {
  const [email, setEmail] = useState('');
  return <input value={email} onChange={(e) => setEmail(e.target.value)} />;
}
```

### Uncontrolled with Refs

```tsx
function Form() {
  const inputRef = useRef<HTMLInputElement>(null);
  const handleSubmit = () => console.log(inputRef.current?.value);
  return <input ref={inputRef} defaultValue="" />;
}
```

## Key Principles

1. **Lift state up** only as needed
2. **Colocate state** with components that use it
3. **Prefer composition** over configuration
4. **Memoize expensive operations**, not everything
5. **Keep components small** and focused
