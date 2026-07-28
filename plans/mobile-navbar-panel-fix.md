# Mobile Navbar Slide Panel Fix

## Problem Summary

The mobile slide-in panel on the main page has these issues:

1. **Missing X close button inside the panel** — [`context/COMPONENTS.md:77`](../context/COMPONENTS.md:77) requires: *"Icon close (X) di pojok kanan atas panel."* Currently the X icon only replaces the hamburger icon in the toggle button, but there's **no close button inside the slide panel itself**.

2. **Overlay display conflict** — The overlay uses both Alpine.js `x-show` AND a CSS class `.navbar-overlay.open { display: block; }` — redundant control.

3. **CSS not modular** — All mobile panel styles live in [`assets/css/style.css`](../assets/css/style.css:322-401) instead of a dedicated file.

## Refined Plan

### Files to Create
| File | Purpose |
|------|---------|
| [`assets/css/mobile-slide-panel.css`](../assets/css/mobile-slide-panel.css) | **New** — All mobile panel CSS extracted + new styles for close button |

### Files to Modify
| File | Change |
|------|--------|
| [`templ/components/navbar.templ`](../templ/components/navbar.templ:72) | Add X close button inside `.navbar-panel` header |
| [`templ/components/navbar.templ`](../templ/components/navbar.templ:68) | Remove `:class="{ 'open': mobileOpen }"` from overlay |
| [`templ/components/head.templ`](../templ/components/head.templ:16) | Add `<link rel="stylesheet" href="/assets/css/mobile-slide-panel.css"/>` |
| [`assets/css/style.css`](../assets/css/style.css:322-401) | Remove all `.navbar-panel*`, `.navbar-overlay*` rules (lines 322–401) |

---

## Detailed Steps

### Step 1 — Create `assets/css/mobile-slide-panel.css`

Extract all existing mobile panel CSS from [`style.css`](../assets/css/style.css:322-401) into this new file, then add new styles for the close button header.

The file will contain:

```css
/* ============================================================
   Mobile Slide Panel — extracted from style.css §navbar
   ============================================================ */

/* --- Overlay --- */
.navbar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  z-index: 49;
}

/* --- Slide Panel --- */
.navbar-panel {
  position: fixed;
  top: 0;
  right: 0;
  width: 80%;
  max-width: 320px;
  height: 100vh;
  background: #FFFFFF;
  z-index: 50;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  padding: 88px 24px 24px;
  overflow-y: auto;
  border-left: 1px solid #f1f5f9;
}

.navbar-panel.open {
  transform: translateX(0);
}

/* --- Panel Header (NEW) --- */
.navbar-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.navbar-panel-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  font-family: 'Inter', system-ui, sans-serif;
}

.navbar-panel-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #1A1A2E;
  cursor: pointer;
  transition: background 0.2s ease;
}

.navbar-panel-close:hover { background: #f1f5f9; }

.navbar-panel-close:focus-visible {
  outline: 2px solid #6366f1;
  outline-offset: 2px;
}

/* --- Panel Navigation Links --- */
.navbar-panel-links {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.navbar-panel-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  color: #475569;
  text-decoration: none;
  transition: background 0.2s ease, color 0.2s ease;
  font-family: 'Inter', system-ui, sans-serif;
}

.navbar-panel-link:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.navbar-panel-link:focus-visible {
  outline: 2px solid #6366f1;
  outline-offset: 2px;
}

.navbar-panel-cta {
  background: #eef2ff;
  color: #6366f1;
  font-weight: 600;
  margin-top: 12px;
}

.navbar-panel-divider {
  border: none;
  border-top: 1px solid #f1f5f9;
  margin: 12px 0;
}

/* --- Desktop: hide mobile elements --- */
@media (min-width: 768px) {
  .navbar-panel,
  .navbar-overlay {
    display: none;
  }
}
```

All pure CSS, no Tailwind. The panel hides via `transform: translateX(100%)` which pushes the **entire element including padding** off-screen — no visible padding when hidden.

### Step 2 — Add X Close Button in `navbar.templ`

Inside the `.navbar-panel` div, insert this before the `<nav class="navbar-panel-links">`:

```html
<div class="navbar-panel-header">
  <span class="navbar-panel-title">Menu</span>
  <button @click="mobileOpen = false"
          aria-label="Tutup menu navigasi"
          class="navbar-panel-close">
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
  </button>
</div>
```

Also remove `:class="{ 'open': mobileOpen }"` from the overlay element — `x-show` alone handles visibility.

### Step 3 — Load New CSS in `head.templ`

Add after the existing hero.css line:
```html
<link rel="stylesheet" href="/assets/css/mobile-slide-panel.css"/>
```

### Step 4 — Remove Old CSS from `style.css`

Delete lines 322–401 (`.navbar-overlay` through the `@media (min-width: 768px)` block).

### Step 5 — Regenerate & Build

```bash
templ generate
go build -o . ./...
```

## Visual Layout

```
Mobile view (<768px):

  ┌──────────────────┬───────────┐
  │                  │ [X] Menu  │  ← NEW: X button at top-right
  │   (overlay)      │ ────────  │      inside the panel
  │   semi-transp.   │ ● Beranda │
  │   black bg       │ ● Fitur   │
  │                  │ ● Tes     │
  │                  │ ───────── │
  │                  │ ▶ Mulai   │
  │                  │   Tes     │
  └──────────────────┴───────────┘
```

The X button sits at the top-right inside the panel, always visible regardless of scroll position within the panel. Tapping it or the overlay background closes the panel.
