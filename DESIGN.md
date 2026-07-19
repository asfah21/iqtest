---
name: IQ Test
description: IQ & cognitive assessment — intelligence metrics, aptitude scoring, performance analytics
colors:
  primary: "#6366f1"
  primary-hover: "#4f46e5"
  primary-light: "#eef2ff"
  primary-border: "#c7d2fe"
  ink: "#0f172a"
  ink-muted: "#475569"
  ink-subtle: "#94a3b8"
  bord: "#e2e8f0"
  bord-light: "#f1f5f9"
  surface: "#ffffff"
  bg: "#f8fafc"
  bg-alt: "#f1f5f9"
  iq-low: "#f59e0b"
  iq-average: "#10b981"
  iq-high: "#6366f1"
  iq-superior: "#8b5cf6"
  success: "#10b981"
  warning: "#f59e0b"
  error: "#ef4444"
typography:
  display:
    fontFamily: "DM Serif Display, Georgia, serif"
    fontSize: "clamp(2.25rem, 5vw, 4rem)"
    fontWeight: 400
    lineHeight: 1.1
    letterSpacing: "-0.02em"
  headline:
    fontFamily: "DM Serif Display, Georgia, serif"
    fontSize: "clamp(1.75rem, 3vw, 2.5rem)"
    fontWeight: 400
    lineHeight: 1.2
    letterSpacing: "-0.01em"
  title:
    fontFamily: "Inter, system-ui, sans-serif"
    fontSize: "clamp(1.25rem, 2vw, 1.5rem)"
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: "-0.01em"
  body:
    fontFamily: "Inter, system-ui, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.65
    letterSpacing: "0em"
  label:
    fontFamily: "Inter, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 500
    lineHeight: 1.5
    letterSpacing: "0em"
  mono:
    fontFamily: "JetBrains Mono, monospace"
    fontSize: "0.875rem"
    fontWeight: 600
    lineHeight: 1.5
    letterSpacing: "0em"
rounded:
  sm: "8px"
  md: "12px"
  lg: "16px"
  full: "9999px"
spacing:
  "2xs": "4px"
  xs: "8px"
  sm: "12px"
  md: "16px"
  lg: "24px"
  xl: "32px"
  "2xl": "48px"
  "3xl": "64px"
  "4xl": "80px"
  "5xl": "96px"
  "6xl": "128px"
elevation:
  level-1:
    boxShadow: "0 1px 2px rgba(15,23,42,0.04), 0 4px 12px rgba(15,23,42,0.06)"
    border: "1px solid #f1f5f9"
  level-2:
    boxShadow: "0 4px 16px rgba(15,23,42,0.08), 0 12px 40px rgba(15,23,42,0.06)"
  navbar-scroll:
    boxShadow: "0 1px 12px rgba(15,23,42,0.06)"
    background: "rgba(255,255,255,0.85)"
    backdropFilter: "blur(16px)"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "#ffffff"
    rounded: "{rounded.sm}"
    padding: "12px 28px"
    typography: "{typography.label}"
  button-primary-hover:
    backgroundColor: "{colors.primary-hover}"
    textColor: "#ffffff"
    rounded: "{rounded.sm}"
    padding: "12px 28px"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.ink-muted}"
    rounded: "{rounded.sm}"
    padding: "12px 28px"
    border: "1px solid {colors.bord}"
  button-secondary-hover:
    backgroundColor: "transparent"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "12px 28px"
    border: "1px solid {colors.ink}"
  button-ghost:
    backgroundColor: "transparent"
    textColor: "{colors.ink-muted}"
    rounded: "{rounded.sm}"
    padding: "8px 16px"
  button-ghost-hover:
    backgroundColor: "{colors.bg}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "8px 16px"
  card:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.md}"
    border: "1px solid {colors.bord-light}"
    boxShadow: "{elevation.level-1.boxShadow}"
  card-flat:
    backgroundColor: "{colors.bg}"
    rounded: "{rounded.md}"
    border: "none"
  card-hover:
    border: "1px solid {colors.bord}"
    boxShadow: "0 4px 16px rgba(15,23,42,0.08)"
    transition: "all 0.2s ease"
  input-field:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "12px 16px"
    border: "1px solid {colors.bord}"
  input-field-focus:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "12px 16px"
    border: "1px solid {colors.primary}"
  progress-bar:
    backgroundColor: "{colors.bord-light}"
    rounded: "{rounded.full}"
    height: "6px"
  progress-fill:
    rounded: "{rounded.full}"
    height: "100%"
    transition: "width 0.6s cubic-bezier(0.4, 0, 0.2, 1)"
---

# Design System: IQ Test

## 1. Overview

**Creative North Star: "Measure Your Potential"**

IQ Test's visual system is a gateway to self-discovery — it
feels sharp, confident, and energizing. The interface creates
a sense of challenge and achievement, not anxiety. Every pixel
serves clarity and cognitive engagement.

The system is **bright and vibrant**: modern indigo accents
on a clean white base create a crisp, professional feel.
Generous spacing and subtle elevation signal quality.
The palette is optimistic — blues and purples evoke
intelligence, growth, and clarity.

**Key Characteristics:**
- Clean white base — pure, minimal, focused
- Indigo primary (#6366f1) — intelligence, trust, energy
- Brighter neutrals — slate-based for clarity
- Vibrant data colors — amber, emerald, indigo, violet
- Multiple accent usage allowed — IQ bands need range
- Modern and crisp — like a premium assessment platform

## 2. Colors

### Primary
- **IQ Indigo** (#6366f1): The primary accent. Interactive
  elements, focus rings, active states, data highlights,
  brand mark. Bold and confident.
- **IQ Indigo Hover** (#4f46e5): Darkened for hover states.
- **IQ Indigo Light** (#eef2ff): Subtle backgrounds,
  badge fills, icon containers.
- **IQ Indigo Border** (#c7d2fe): Bordered containers
  needing a hint of brand.

### Neutral (Clean)
- **Ink** (#0f172a): Primary text. Slate near-black —
  sharp, authoritative.
- **Ink Muted** (#475569): Secondary text, body copy.
  WCAG AA 4.5:1 against surface.
- **Ink Subtle** (#94a3b8): Placeholders, metadata, captions.
  WCAG AA 4.5:1 against surface.
- **Bord** (#e2e8f0): Standard borders, dividers. Clean slate.
- **Bord Light** (#f1f5f9): Subtle borders, card outlines.
- **Surface** (#ffffff): Primary background. Pure white —
  crisp, focused, modern.
- **Bg** (#f8fafc): Section alternation, flat card backgrounds.
- **Bg Alt** (#f1f5f9): Deeper alternation, subtle fills.

### Data (IQ Score Bands)
- **IQ Low** (#f59e0b): Amber. Scores below 85.
  Used in score displays and progress bars only.
- **IQ Average** (#10b981): Emerald. Scores 85–115.
  Used in score displays and progress bars only.
- **IQ High** (#6366f1): Indigo. Scores 115–130.
  Used in score displays and progress bars only.
- **IQ Superior** (#8b5cf6): Violet. Scores above 130.
  Used in score displays and progress bars only.

### Semantic
- **Success** (#10b981): Confirmation, positive states.
- **Warning** (#f59e0b): Cautionary indicators.
- **Error** (#ef4444): Error states, destructive actions.

### Named Rules

**The Intentional Accent Rule.** IQ Indigo is the primary
accent for interactive elements. IQ band colors (Amber,
Emerald, Indigo, Violet) are exclusive to score displays
and progress bars. Never mix in navigation or CTAs.

**The Clean Base Rule.** Surface is always #ffffff.
Off-white (#f8fafc) for section alternation only.
No warm undertones — this is a precision instrument,
not a cozy corner.

**The Modern Shadow Rule.** Shadows use slate ink tone
(#0f172a), never pure black. Opacity stays below 0.08
— shadows suggest depth, they don't announce it.

## 3. Typography

**Display/Headline:** DM Serif Display — editorial gravity,
intellectual weight. Used sparingly for impact.

**Body/UI:** Inter — clean, sharp, highly readable.
Recedes into comfortable readability.

**Mono:** JetBrains Mono — for scores, percentiles,
and data points where precision matters.

### Hierarchy
- **Display** (DM Serif Display, 400, clamp(2.25rem, 5vw,
  4rem), 1.1, -0.02em): Hero headlines only. text-wrap: balance.
- **Headline** (DM Serif Display, 400, clamp(1.75rem, 3vw,
  2.5rem), 1.2, -0.01em): Section headings. text-wrap: balance.
- **Title** (Inter, 600, clamp(1.25rem, 2vw, 1.5rem), 1.3,
  -0.01em): Card titles, subheadings. text-wrap: balance.
- **Body** (Inter, 400, 1rem, 1.65): Paragraphs, descriptions.
  Max 65–75ch. text-wrap: pretty.
- **Label** (Inter, 500, 0.875rem, 1.5): Buttons, nav, labels.
- **Mono** (JetBrains Mono, 600, 0.875rem, 1.5): Score values,
  percentiles, data points.

### Named Rules

**The Serif-for-Gravity Rule.** DM Serif Display for display
and headline only. Inter for everything else. Never mix
within the same hierarchy level.

**The Weight-Only Hierarchy Rule.** Hierarchy through weight,
size, and font family — not color, tracking, or case changes.
No uppercase eyebrows. No wide-tracked labels.

## 4. Elevation

Three levels of intentional depth. Used to indicate
interactivity and layering — never decoration.

**Level 0 — Flat:** Page background, section containers.
No shadow, no border. Base surface (#ffffff).

**Level 1 — Raised:** Cards, input fields.
box-shadow: 0 1px 2px rgba(15,23,42,0.04),
0 4px 12px rgba(15,23,42,0.06);
border: 1px solid #f1f5f9;

**Level 2 — Floating:** Dropdowns, tooltips, modals.
box-shadow: 0 4px 16px rgba(15,23,42,0.08),
0 12px 40px rgba(15,23,42,0.06);

**Navbar on scroll:**
background: rgba(255,255,255,0.85);
backdrop-filter: blur(16px);
box-shadow: 0 1px 12px rgba(15,23,42,0.06);

### Named Rules

**The Intentional Elevation Rule.** Three levels only.
Cool-toned, low opacity, purposeful. No colored shadows,
no neon glows, no material-design-style hard shadows.

**The No-Pure-Flat Rule.** Cards and interactive surfaces
use Level 1 elevation — not fully flat. Flatness is for
backgrounds and non-interactive containers only.

## 5. Components

### Buttons
- **Primary:** IQ Indigo (#6366f1) background, white text,
  8px radius, 12px 28px padding.
  Hover: #4f46e5. Transition: background 0.2s ease.
  Focus: 2px indigo outline, 2px offset. No shadow, no glow.
- **Secondary:** Transparent, inkMuted text, 1px bord border.
  Hover: border → ink, text → ink.
- **Ghost:** Transparent, inkMuted text, 8px 16px padding.
  Hover: bg background, ink text.

### Navbar
- **Default:** background #ffffff, 1px bord-light bottom border.
  No blur, no shadow.
- **On scroll:** background rgba(255,255,255,0.85),
  backdrop-filter blur(16px),
  box-shadow 0 1px 12px rgba(15,23,42,0.06).
- Transition: all 0.3s ease on scroll trigger.
- Height: 56–64px max. Slim, invisible when possible.
- Layout: logo left, nav links center, CTA button right.
- Links: ghost style. No active indicator beyond URL.
- Respects prefers-reduced-motion: no blur/transition
  if reduced motion enabled.

### Hero
- Background: #ffffff — always light, never dark.
- Headline: DM Serif Display, large, ink (#0f172a).
- One word or phrase in IQ Indigo (#6366f1).
- Typographic decoration: oversized symbol or word at
  3–5% opacity (e.g. "IQ", "Σ", "≥") — NOT illustration.
- Subtle animated gradient shimmer on headline
  (CSS only, respects prefers-reduced-motion).
- Two CTAs: primary indigo + secondary ghost, side by side.
- Trust pills below CTAs: "Anonim" · "Gratis" · "5 Menit"
- NO photos, NO Lottie, NO SVG illustrations, NO dark hero.

### Cards
- Background: surface (#ffffff).
- Border: 1px solid bord-light (#f1f5f9).
- Shadow: Level 1 elevation.
- Radius: 12px (rounded.md).
- Padding: 24–32px.
- Hover: border → bord (#e2e8f0),
  shadow → 0 4px 16px rgba(15,23,42,0.08).
  Transition: all 0.2s ease.

### Inputs
- Border: 1px solid bord (#e2e8f0).
- Background: surface (#ffffff).
- Radius: 8px. Padding: 12px 16px.
- Focus: border → primary (#6366f1). No glow.
- Placeholder: inkSubtle (#94a3b8) — WCAG AA.
- Error: border → error (#ef4444).

### Progress Bars
- Height: 6px. Track: bord-light. Full-radius caps.
- Fill: IQ band color — no gradients.
- Transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1).

### Footer
- Background: bg (#f8fafc) — slightly deeper than surface.
- Top border: 1px solid bord-light only.
- No dark inversion, no color change.
- Layout: logo + tagline left, 2–3 link columns right.
- Copyright: inkSubtle, small, bottom.

## 6. Accessibility

- WCAG AA minimum on all elements.
- AAA contrast target for body text and CTAs where palette allows.
- prefers-reduced-motion: all animations and transitions
  must be disabled or reduced when enabled.
- Color blindness: never use color as sole differentiator.
  Always pair with shape, label, or position.
- ARIA labels on all interactive elements.
- Focus-visible: 2px indigo outline, 2px offset on all
  interactive elements.

## 7. Do's and Don'ts

### Do:
- Use IQ Indigo (#6366f1) as the primary accent
- Use pure white (#ffffff) as base — crisp and modern
- Use DM Serif Display for hero and section headlines
- Use Level 1 elevation on cards — not fully flat
- Use glass navbar effect on scroll
- Use tonal layering for section depth
- Keep body text at inkMuted (#475569) — WCAG AA
- Cap body line length at 65–75ch
- Use IQ band colors (Amber, Emerald, Indigo, Violet) exclusively for score data
- Respect prefers-reduced-motion on all animations
- Pair color with shape/label for color blindness support
- Use slate-based neutrals for a clean, professional feel

### Don't:
- Don't use pure black (#000000) — use ink (#0f172a)
- Don't use Inter for display or headline — use DM Serif Display
- Don't use black/ink for primary button — use indigo
- Don't use gradients — backgrounds, text, or buttons
- Don't use glassmorphism outside navbar scroll state
- Don't use illustrations, Lottie, or SVG sketches
- Don't use uppercase tracked eyebrows above headings
- Don't use numbered section markers (01, 02, 03)
- Don't use colored shadows or neon glows
- Don't use IQ band colors outside score displays
- Don't gamify — no badges, achievements, celebrations
- Don't use warm neutrals — this is a precision tool, not cozy
- Don't use border-radius larger than 16px on cards
- Don't pair border + heavy shadow on the same element