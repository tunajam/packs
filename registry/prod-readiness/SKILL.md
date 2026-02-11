---
name: prod-readiness
description: |
  Production readiness checklist and setup. Run when ready to ship to real users.
  Handles Clerk prod, Convex prod, Vercel, Stripe live, monitoring, legal.
  Separate from project-init which is dev-only.
---

# Prod Readiness

Ship to production with confidence. Run through this when you're ready for real users.

---

## When to Use

- All P0 features complete
- E2E tests passing
- Ready for real users / App Store submission
- Moving from dev to production environment

---

## Prerequisites

Before running through this checklist:

- [ ] All P0 user stories marked done in MC
- [ ] `bun ready` passes (typecheck + lint)
- [ ] E2E tests passing locally
- [ ] App tested manually by human
- [ ] PRD acceptance criteria verified

---

## The Checklist

### 1. Authentication (Clerk)

**Create Production Instance:**

1. Go to [clerk.com](https://clerk.com)
2. Create new application: `projectname-prod`
3. Copy production keys

```bash
# Add to Vercel production environment
vercel env add NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY production
vercel env add CLERK_SECRET_KEY production
```

**Configure:**
- [ ] Production domain added to Clerk
- [ ] Social providers configured (if using)
- [ ] Webhook endpoint configured for user sync

---

### 2. Database (Convex)

**Deploy to Production:**

```bash
# Create production deployment
bunx convex deploy --prod

# Get production URL
bunx convex dashboard  # Copy prod URL
```

**Configure:**
- [ ] Production URL in environment variables
- [ ] Auth (Clerk) configured in Convex dashboard
- [ ] Indexes optimized for production queries

```bash
# Vercel
vercel env add NEXT_PUBLIC_CONVEX_URL production
```

---

### 3. Hosting (Vercel) — Web Only

**Setup:**

```bash
# Link project (if not already)
vercel link

# Set environment to production
vercel env pull .env.production.local
```

**Configure:**
- [ ] Custom domain added
- [ ] SSL certificate active
- [ ] Environment variables set for production
- [ ] Preview deployments use dev keys (not prod!)

**Environment Variable Strategy:**
| Vercel Environment | Uses |
|--------------------|------|
| Production | Prod keys (Clerk prod, Convex prod, Stripe live) |
| Preview | Dev keys (safe for PR testing) |
| Development | Dev keys (local) |

---

### 4. Payments (if applicable)

**Stripe (Web):**

```bash
# Switch to live mode in Stripe dashboard
# Copy live keys

vercel env add STRIPE_SECRET_KEY production
vercel env add NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY production
```

- [ ] Live webhook endpoint configured
- [ ] Products/prices created in live mode
- [ ] Test purchase completed with real card

**RevenueCat (Mobile):**

1. Configure production app in RevenueCat
2. Link to App Store Connect / Play Console
3. Set up entitlements for live products

- [ ] Production API key configured
- [ ] Products created in App Store Connect / Play Console
- [ ] Sandbox testing completed

---

### 5. Analytics (PostHog)

- [ ] Verify events tracking in production
- [ ] Session replay enabled (if desired)
- [ ] Funnels configured for key flows
- [ ] Alerts set up for errors/drops

```bash
# Verify PostHog is receiving events
# Check https://app.posthog.com project dashboard
```

---

### 6. Error Tracking (Sentry)

**Setup Production Project:**

```bash
# Web
bunx @sentry/wizard@latest -i nextjs

# Mobile
bunx @sentry/wizard@latest -i reactNative
```

- [ ] Production DSN configured
- [ ] Source maps uploaded
- [ ] Alert rules configured
- [ ] Test error sent and received

---

### 7. Domain & SSL

**Web:**
- [ ] Domain purchased/configured
- [ ] DNS pointing to Vercel
- [ ] SSL certificate active (automatic with Vercel)
- [ ] WWW redirect configured

**Mobile:**
- [ ] Deep linking domain configured
- [ ] Universal links / App links set up
- [ ] Associated domains file hosted

---

### 8. Legal (if needed)

- [ ] Privacy Policy page/link
- [ ] Terms of Service page/link
- [ ] Cookie consent (if EU users)
- [ ] GDPR data export/delete capability (if EU users)

**For App Store:**
- [ ] Privacy policy URL in App Store Connect
- [ ] Data collection disclosures accurate

---

### 9. App Store / Play Store (Mobile Only)

**Apple App Store:**
- [ ] App Store Connect app created
- [ ] Screenshots uploaded
- [ ] App description written
- [ ] Keywords optimized
- [ ] Privacy policy URL added
- [ ] Build uploaded via EAS Submit
- [ ] TestFlight testing completed

```bash
eas build --platform ios --profile production
eas submit --platform ios
```

**Google Play Store:**
- [ ] Play Console app created
- [ ] Store listing complete
- [ ] Content rating questionnaire
- [ ] Build uploaded
- [ ] Internal testing completed

```bash
eas build --platform android --profile production
eas submit --platform android
```

---

### 10. Monitoring & Alerts

- [ ] Uptime monitoring (e.g., BetterStack, Pingdom)
- [ ] Error rate alerts in Sentry
- [ ] Revenue alerts in Stripe/RevenueCat
- [ ] Performance monitoring enabled

---

### 11. Backup & Recovery

- [ ] Convex data export tested
- [ ] Recovery procedure documented
- [ ] Team access configured (if applicable)

---

### 12. Documentation Update

- [ ] README updated for production
- [ ] CLAUDE.md updated with prod commands
- [ ] Runbook for common issues created
- [ ] On-call procedure documented (if applicable)

---

## Final Verification

### Smoke Test Checklist

Run through core flows on production:

- [ ] New user can sign up
- [ ] Existing user can sign in
- [ ] Core action works (the main thing)
- [ ] Payment flow works (if applicable)
- [ ] Errors are caught in Sentry
- [ ] Analytics events appear in PostHog

### Performance Check

- [ ] Page load < 3s on 3G
- [ ] Core Web Vitals passing
- [ ] No console errors
- [ ] No memory leaks on mobile

---

## Launch!

Once all checks pass:

1. **Flip the switch** — Enable production access
2. **Monitor closely** — Watch Sentry, PostHog, Stripe for first few hours
3. **Announce** — Post to relevant communities, social media
4. **Update MC** — Mark launch task as done
5. **Celebrate** 🎉

---

## Post-Launch

**First 24 hours:**
- Monitor error rates
- Respond to user feedback
- Hot-fix critical bugs

**First week:**
- Review analytics
- Identify drop-off points
- Plan first iteration

---

## Quick Reference

```bash
# Production deploy commands

# Convex
bunx convex deploy --prod

# Vercel (auto on merge to main)
vercel --prod

# iOS
eas build --platform ios --profile production
eas submit --platform ios

# Android
eas build --platform android --profile production
eas submit --platform android
```

---

*Dev is done. Time to ship to real humans.*
