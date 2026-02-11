---
name: reddit-business-research
description: |
  Mine Reddit for solo founder business ideas and pain points. Use for morning briefs,
  idea validation, or when hunting for problems worth solving. Finds real complaints
  from real people willing to pay.
---

# Reddit Business Research

Find validated business ideas by mining Reddit for pain points, frustrations, and "I wish there was an app for this" posts.

## When to Use

- Morning brief founder idea section
- Validating a specific business concept
- Finding problems in a particular niche
- Researching competition sentiment

## High-Signal Subreddits

### For SaaS/Tech Ideas
| Subreddit | Signal Type |
|-----------|-------------|
| r/SaaS | Direct pain points, what's working, case studies |
| r/Entrepreneur | Scaling problems, tool needs, workflow gaps |
| r/startups | Early-stage validation signals |
| r/smallbusiness | Non-tech business owner frustrations |
| r/SideProject | What solos are building, gaps in market |

### For Niche Discovery
| Subreddit | Why It's Gold |
|-----------|---------------|
| r/ADHD | Most detailed feature requests — current tools fail them |
| r/personalfinance | High willingness to pay, specific tool needs |
| r/homelab / r/homeassistant | Privacy-focused tech users with money |
| r/freelance | Service business pain points |
| r/ecommerce | Shopify/Amazon seller tool gaps |

## Search Queries That Work

### Finding Pain Points
```
site:reddit.com "I wish there was" OR "someone should build" OR "I'd pay for"
site:reddit.com "frustrated with" OR "hate using" OR "why is there no"
site:reddit.com "switched from" OR "looking for alternative to"
```

### Subreddit-Specific (use old.reddit.com for better scraping)
```
https://old.reddit.com/r/SaaS/search?q=problem+OR+pain+point&restrict_sr=on&sort=top&t=month
https://old.reddit.com/r/Entrepreneur/search?q=frustrated+OR+struggling&restrict_sr=on&sort=top&t=month
```

## Key Metrics to Extract

### 1. Willingness to Pay Signals
Look for keywords: "buy", "pay", "premium", "subscription", "worth it"

**High-paying niches (from 9,300+ post analysis):**
- Finance: 193 pay signals (highest)
- E-commerce: 76 pay signals
- Travel: 42 pay signals

### 2. Frustration Score (Post Length)
Longer posts = deeper pain = better opportunity

**Highest frustration niches:**
- Developer Platforms: 229 avg length
- Cooking/Recipes: 223 avg length  
- Parenting: 221 avg length

### 3. Trending Categories
Check what's spiking recently:
- Health/Fitness spikes January (New Year resolutions)
- Smart Home spikes after holidays (new devices)
- Productivity spikes Monday/Tuesday (work week begins)

## The "Anti-Cloud" Trend

~7% of all requests specifically want offline-first or local-only tools. Subscription fatigue is real. One-time purchase + privacy = differentiator.

## Research Workflow

### Quick Scan (5 min)
1. Open `old.reddit.com/r/SaaS/top?t=week`
2. Look for complaint posts, "I built X" posts with comments
3. Note recurring themes

### Deep Dive (15-30 min)
1. Search specific pain point across multiple subs
2. Read comment threads (gold is in replies)
3. Look for "I'd pay for this" or "shut up and take my money"
4. Note specific feature requests

### Validation Check
Before recommending an idea, verify:
- [ ] **Check our existing ideas first** — Search Linear (TUN/TEN) and ideas.tunajam.com
- [ ] **Check market** — Is there already a dominant solution? (e.g., "Just the Recipe" for recipe cleaning)
- [ ] Multiple people complaining (not just one person)
- [ ] Recent posts (problem still exists)
- [ ] No obvious dominant solution OR clear differentiation angle
- [ ] Scope fits solo founder (buildable in weeks, not years)

## Output Format for Morning Brief

```markdown
## 💡 Solo Founder Idea

**Problem:** [One sentence pain point]

**Signal:** [Where you found it, how many complaints]

**Solution Sketch:** [Simple MVP description]

**Why Now:** [Trend, timing, or gap in market]

**Buildability:** [Can we ship this in days/weeks?]
```

## Example Finds

### Recipe Cleaner
- **Source:** r/cooking, r/recipes complaints
- **Signal:** 223 avg post length (high frustration), daily complaints about recipe blog bloat
- **Idea:** Browser extension that strips recipes to just ingredients + steps

### ADHD Task Manager  
- **Source:** r/ADHD feature requests
- **Signal:** Most detailed user feedback of any niche — current tools fail neurodivergent workflows
- **Idea:** Task manager built for urgency-based prioritization, dopamine-friendly design

### Local Finance Dashboard
- **Source:** r/personalfinance, r/SaaS
- **Signal:** 193 pay signals (highest), 7% want offline-first
- **Idea:** One-time purchase portfolio tracker with beautiful data viz

## Tips

1. **Read the comments** — The post might be generic, but replies contain specific pain
2. **Sort by "top" for validation, "new" for freshness**
3. **Cross-reference complaints** — Same pain in 3+ subs = real problem
4. **Look for failed products** — "Why did X shut down?" reveals demand + execution gaps
5. **"I switched from X to Y because..."** — Reveals feature gaps in existing tools
