# DESIGN.md — Warna, Tipografi, Radius, Shadow, Icon & Animation

> Fokus dokumen ini: bahasa visual (design tokens) yang digunakan di dalam project. Untuk struktur grid lihat [`LAYOUT.md`](context/LAYOUT.md), untuk aturan komponen lihat [`COMPONENTS.md`](context/COMPONENTS.md).

---

## 1. Palet Warna

> Design system utama project menggunakan palet **indigo** (`#6366f1`) sebagai accent color. Design terinspirasi dari 9router.com — modern, elegan, developer-focused. Lihat `:root` di [`style.css`](assets/css/style.css).

### 1.1 Warna Utama (Main Design System — `assets/css/style.css`)

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--accent` | `#6366f1` | CTA utama, link aktif, highlight, focus ring |
| `--accentHover` | `#4f46e5` | Hover state tombol primary |
| `--accentLight` | `#eef2ff` | Background ringan untuk ikon fitur, elemen aksen |

### 1.2 Warna Netral

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--bg` | `#f8fafc` | Background halaman utama |
| `--surface` | `#ffffff` | Background card, navbar, elemen permukaan |
| `--bgSection` | `#f1f5f9` | Background section selang-seling (`.section-alt`) |
| `--bgFooter` | `#e2e8f0` | (Tersedia tapi tidak dipakai — footer menggunakan background gelap `#1a1917`) |
| `--ink` | `#0f172a` | Teks body utama |
| `--inkMuted` | `#475569` | Teks sekunder, caption |
| `--inkSubtle` | `#94a3b8` | Teks tersier, placeholder |
| `--bord` | `#e2e8f0` | Border umum (card, divider) |
| `--bordLight` | `#f1f5f9` | Border halus untuk card |

### 1.3 Warna Status & Data (Dark Triad)

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--narcissus` | `#7c3aed` | Data narsisistik / domain ungu |
| `--machiavellian` | `#d97706` | Data machiavellian / domain amber |
| `--psychopath` | `#059669` | Data psychopath / domain hijau |
| `--success` | `#059669` | (Backward compat alias untuk psychopath) |
| `--warning` | `#d97706` | (Backward compat alias untuk machiavellian) |
| `--error` | `#dc2626` | Pesan error, badge error |

### 1.4 Quiz Design Tokens (`assets/css/quiz.css`)

Quiz page memiliki design system terpisah (navy + pink):

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--qnavy` | `#16324F` | Warna header quiz |
| `--qorange` | `#EC4899` | Aksen quiz (tombol, indikator) |
| `--qorange-hover` | `#f57bb8` | Hover state aksen quiz |
| `--qwhite` | `#FFFFFF` | Background card quiz |
| `--qbg` | `#F2F3F5` | Background halaman quiz |
| `--qgray-900` | `#333333` | Teks gelap quiz |
| `--qgray-500` | `#8A8F98` | Teks muted quiz |
| `--qgray-300` | `#D9DEE4` | Border quiz |
| `--qgray-light` | `#E8ECF0` | Background input quiz |

### 1.5 Hero Colors (`assets/css/hero.css`)

Hero menggunakan hardcoded warna spesifik, tidak mereferensi CSS variables dari design system utama:

| Token (HERO.md) | Hex |
|------------------|-----|
| `--color-accent-blue` | `#4A5CF5` |
| `--color-accent-blue-2` | `#5B6EF5` |
| `--color-accent-purple` | `#6C5CE7` |
| `--color-accent-yellow` | `#F5A623` |
| `--color-accent-orange` | `#FB923C` |
| `--color-navy-header` | `#12294D` |
| `--color-text-dark` | `#111827` |
| `--color-text-secondary` | `#4B5563` |
| `--pastel-purple` | `#EDE9FE` |
| `--pastel-green` | `#DCFCE7` |
| `--pastel-orange` | `#FEF3C7` |
| `--pastel-blue` | `#DBEAFE` |

### 1.6 Aturan Penggunaan
- Warna `--accent` hanya untuk elemen aksi (CTA, link penting) — jangan dipakai sebagai warna teks panjang.
- Dark triad colors (`--narcissus`, `--machiavellian`, `--psychopath`) khusus untuk visualisasi data hasil tes di dashboard.
- Section background berselang-seling antara `--bg` (`.section`) dan `--bgSection` (`.section-alt`) untuk membedakan tiap section secara visual tanpa garis pembatas tegas.
- Footer menggunakan background gelap `#1a1917` (hardcoded di [`footer.css`](assets/css/footer.css)).
- Hero section memiliki design system warna sendiri yang didefinisikan di [`HERO.md`](context/HERO.md).

---

## 2. Tipografi

### 2.1 Font Family
- **Body default:** `'Inter', system-ui, -apple-system, sans-serif`
- **Hero headlines:** `'Poppins', 'Inter', system-ui, sans-serif`
- **Navbar brand:** `'Poppins', 'Nunito', system-ui, sans-serif`
- **Font monospace (terminal/code):** `'JetBrains Mono', monospace`

### 2.2 Skala Ukuran (Type Scale)

> Menggunakan `clamp()` untuk fluid typography — tidak ada ukuran fixed px.

| Token | Value (desktop) | Weight | Line-height | Penggunaan |
|-------|-----------------|--------|-------------|------------|
| `--text-xs` | `0.75rem` (12px) | 500 | `--leading-normal` (1.5) | Label kecil, badge, caption |
| `--text-sm` | `0.875rem` (14px) | 400 | `--leading-relaxed` (1.65) | Teks card, footer link |
| `--text-base` | `1rem` (16px) | 400 | `--leading-relaxed` (1.65) | Body paragraf |
| `--text-lg` | `clamp(1.125rem, 1.5vw, 1.25rem)` | 400 | `--leading-normal` (1.5) | Lead text, subheadline |
| `--text-xl` | `clamp(1.5rem, 2.5vw, 2rem)` | 600 | `--leading-tight` (1.1) | H3 |
| `--text-2xl` | `clamp(1.75rem, 3vw, 2.5rem)` | 700 | `--leading-tight` (1.1) | H2, section heading |
| `--text-3xl` | `clamp(2rem, 4vw, 3.25rem)` | 800 | `--leading-tight` (1.1) | Angka besar (trust bar) |
| `--text-4xl` | `clamp(2.25rem, 5vw, 4rem)` | 800 | `--leading-tight` (1.1) | H1 |

### 2.3 Leading (Line Height) Tokens

| Token | Value |
|-------|-------|
| `--leading-tight` | `1.1` |
| `--leading-snug` | `1.25` |
| `--leading-normal` | `1.5` |
| `--leading-relaxed` | `1.65` |

### 2.4 Letter Spacing Tokens

| Token | Value |
|-------|-------|
| `--tracking-tight` | `-0.03em` |
| `--tracking-normal` | `0em` |
| `--tracking-wide` | `0.04em` |

### 2.5 Measure (Max Width untuk Teks)

| Token | Value | Penggunaan |
|-------|-------|------------|
| `--measure-narrow` | `45ch` | Section heading terpusat |
| `--measure-body` | `65ch` | Paragraf body |
| `--measure-wide` | `75ch` | Area teks lebar |

### 2.6 Aturan
- H1 hanya boleh muncul satu kali per halaman (di Hero).
- Heading menggunakan `text-wrap: balance` untuk mencegah orphans.
- Body text menggunakan `text-wrap: pretty`.
- Label kecil di atas H1 (badge sosial proof) memakai `font-size:14px`, `font-weight:500`, warna `#1A1A2E`.
- Font quiz page menggunakan `'Poppins', 'Nunito', system-ui, sans-serif` sebagai `--qfont`.

---

## 3. Border Radius

| Token | Value | Penggunaan |
|-------|-------|------------|
| `--radiusBtn` | `8px` | Button, input, badge kecil |
| `--radiusCard` | `12px` | Card, mockup, modal umum |
| `rounded-lg` (utility) | `8px` | Alternatif untuk radius kecil |
| `rounded-xl` (utility) | `12px` | Alternatif untuk radius card |
| `rounded-2xl` (utility) | `16px` | Card besar |
| `rounded-full` (utility) | `9999px` | Pill button, badge bulat, CTA hero |

### Nilai Hardcoded Spesifik
- **CTA Hero:** `border-radius: 999px` (pill penuh)
- **Navbar:** `border-radius: 22px` (floating pill)
- **Quiz option:** `border-radius: 14px`
- **Quiz header:** `border-radius: 22px`
- **Quiz card:** `--qradius-card: 18px`

---

## 4. Shadow / Elevation

| Token | Value | Penggunaan |
|-------|-------|------------|
| `--shadow1` | `0 1px 2px rgba(15,23,42,0.04), 0 4px 12px rgba(15,23,42,0.06)` | Card default, navbar |
| `--shadow2` | `0 4px 16px rgba(15,23,42,0.08), 0 12px 40px rgba(15,23,42,0.06)` | Card hover, modal, dropdown |

### Backward Compatibility Aliases

| Token | Value |
|-------|-------|
| `.shadow-sm` | `0 1px 2px rgba(26,25,23,0.04)` |
| `.shadow-md` | `0 4px 6px rgba(26,25,23,0.06)` |
| `.shadow-xl` | `var(--shadow2)` |
| `.hover\:shadow-md:hover` | `0 4px 6px rgba(26,25,23,0.06)` |

### Shadow Spesifik Komponen
- **CTA Hero:** `box-shadow: 0 12px 24px rgba(91,110,245,0.35)` ; hover: `0 16px 32px rgba(91,110,245,0.45)`
- **Navbar scrolled:** `box-shadow: 0 4px 20px rgba(0,0,0,0.08)`
- **Quiz shadow:** `--qshadow: 0 4px 12px rgba(0,0,0,0.06)`
- **Focus ring input:** `box-shadow: 0 0 0 3px rgba(13,115,119,0.1)` (menggunakan `--accent-500` teal legacy)

---

## 5. Icon

- **Library:** Semua icon menggunakan **inline SVG** langsung di dalam HTML/templ — bukan icon library eksternal.
- **Style:** Stroke-based (outline), stroke-width bervariasi:
  - `1.5` — ikon dekoratif (navbar brand, footer)
  - `1.6` — ikon hero badge dan fitur
  - `1.8` — ikon fitur hero
  - `2.5` — ikon navigasi (navbar toggle, chevron, panel)
- **Ukuran standar:**
  - Icon inline (dekat teks): `16px` atau `20px`
  - Icon di card fitur: `24px` (dalam container `44px × 44px`)
  - Icon di Trust Bar: tidak menggunakan icon terpisah — hanya angka besar + label
  - Icon brand (navbar/footer): `24px` (dalam container `44px × 44px`)
- **Warna icon:**
  - Default: mengikuti warna teks sekitarnya (`currentColor`)
  - Ikon fitur di hero: warna ungu `#7C3AED`, hijau `#16A34A`, amber `#F59E0B`, biru `#2563EB`
  - Ikon di panel mobile: `currentColor` (inherit dari link)
- Khusus untuk SVG dekoratif hero section, lihat [`HERO.md`](context/HERO.md) §6 untuk daftar lengkap placeholder SVG.

---

## 6. Ilustrasi & Gambar

- Hero menggunakan **multi-device mockup** (laptop + tablet + phone) dengan SVG inline, bukan foto stok generik. Lihat [`HERO.md`](context/HERO.md) §3.
- Ilustrasi "Tentang Tes" menggunakan placeholder `div` dengan background `--accentLight`, belum ada gambar final (tinggi `300px`).
- Semua gambar soal (`assets/images/q_*.svg`) adalah SVG pola matriks untuk quiz.

---

## 7. Animation & Transition

### 7.1 Durasi & Easing (Hardcoded)

Project **tidak** menggunakan CSS variables untuk durasi/easing animasi. Semua nilai hardcoded di masing-masing komponen:

| Komponen | Durasi | Easing |
|----------|--------|--------|
| Hover button/card | `0.2s` | `ease` |
| Accordion FAQ expand | `0.3s` | `ease` |
| Back-to-top button | `0.3s` | `ease` |
| Navbar scroll transition | `0.3s` | `ease` |
| Quiz option selection | `0.2s` | `ease` |
| Quiz progress fill | `0.6s` | `cubic-bezier(0.22, 1, 0.36, 1)` |
| Quiz step fade-in | `0.35s` | `ease` |
| Scroll reveal (fade-in-up) | `0.5s` | `ease` |
| Marquee integration logos | `30s` | `linear infinite` |
| Marquee testimonials | `40s` | `linear infinite` |
| Terminal cursor blink | `1s` | `step-end infinite` |
| Loading dots pulse | `1.4s` | `ease-in-out infinite` |
| Loading bar indeterminate | `1.5s` | `ease-in-out infinite` |

### 7.2 Aturan Animasi
- **Hover button/card:** transisi `background-color`, `box-shadow`, `transform: translateY(-1px)` dengan durasi `0.2s ease`.
- **Accordion FAQ:** expand/collapse `max-height` dengan durasi `0.3s ease`. Icon rotate `45deg` saat expanded.
- **Scroll-reveal:** `.fade-in-up` + `.is-visible` menggunakan keyframe `fadeInUp` (opacity `0→1`, `translateY(12px → 0)`), durasi `0.5s`. Variasi `.fade-in-up-delayed` dengan delay `0.15s`.
- **Navbar:** transisi `box-shadow` dan `padding` saat scroll, durasi `0.3s ease`. Floating navbar dengan `top: 16px`, menyusut ke `padding: 6px 14px` saat scrolled.
- **Back-to-top button:** opacity + visibility + transform (`translateY(8px → 0)`), durasi `0.3s ease`. Muncul setelah scroll melewati 1 viewport.
- **Marquee (integration logos & testimonials):** animasi `translateX` infinite, pause on hover.
- **Quiz transitions:** Fade-in antar step (`quizFadeIn`), progress bar smooth (`cubic-bezier`), option hover `translateX(4px)`.
- **CTA button:** hover menambah `transform: translateY(-1px)`.
- Hindari animasi berlebihan — semua animasi `@media (prefers-reduced-motion: reduce)` di-reset ke `0.01ms`.

### 7.3 Keyframes Defined

| Keyframe | Penggunaan |
|----------|------------|
| `fadeInUp` | Scroll reveal sections |
| `lineReveal` | (Tersedia, tidak aktif digunakan) |
| `marqueeScroll` | Integration logos marquee |
| `testimonialScroll` / `testimonialScrollReverse` | Testimonial marquee |
| `blink` | Terminal cursor |
| `quizFadeIn` | Quiz step transition |
| `quizPulse` | Loading dots |
| `quizIndeterminate` | Loading bar indeterminate |
| `pulseDot` | General loading dots |
| `spin` | Spinner (admin/payment) |

---

## 8. Warna Komponen Spesifik

### 8.1 Navbar
- Background: `#FFFFFF`
- Shadow: `0 2px 8px rgba(0,0,0,0.15)` (default), `0 4px 20px rgba(0,0,0,0.08)` (scrolled)
- Brand icon background: `linear-gradient(135deg, #7C6FF0, #5B4FE0)`
- Link text: `#475569` (default), `#0f172a` (hover)
- CTA button: background `#6366f1`, hover `#4f46e5`

### 8.2 Hero (lihat [`HERO.md`](context/HERO.md) untuk detail lengkap)
- CTA button: `linear-gradient(135deg, #5B6EF5, #6C5CE7)`
- Badge border: `#ECECF3`
- Disclaimer background: `#F1F3F9`

### 8.3 Footer
- Background: `#1a1917`
- Teks heading: `rgba(255,255,255,0.7)`
- Teks link: `rgba(255,255,255,0.4)` (default), `var(--accent)` (hover)
- Teks tagline: `rgba(255,255,255,0.45)`
- Teks copyright: `rgba(255,255,255,0.3)`
- Divider: `rgba(255,255,255,0.08)`
- Logo icon: `linear-gradient(135deg, #7C6FF0, #5B4FE0)` (sama dengan navbar)

### 8.4 Footer Mobile
- Teks link lebih terang: `rgba(255,255,255,0.55)` untuk contrast lebih baik

---

## 9. Backward Compatibility

CSS mendefinisikan alias untuk variable lama di `:root` kedua (baris 2202-2239 di [`style.css`](assets/css/style.css)):

| Alias | Maps to |
|-------|---------|
| `--textMain` | `var(--ink)` |
| `--textMuted` | `var(--inkMuted)` |
| `--textSubtle` | `var(--inkSubtle)` |
| `--warmBg` | `var(--bg)` |
| `--warmBord` | `var(--bordLight)` |

Serta color scale lengkap `--accent-50` sampai `--accent-700` (teal `#0d7377` based), `--brand-50` sampai `--brand-700`, `--warm-50` sampai `--warm-600`, dan light variants untuk dark triad.

---

## 10. Selection & Scrollbar

- **Selection highlight:** `background: var(--accentLight); color: var(--ink)`
- **Scrollbar:** width `6px`, thumb `var(--bord)` dengan `border-radius: 999px`, track transparan
- **Focus ring:** `outline: 2px solid var(--accent); outline-offset: 2px` via `:focus-visible`

---

## 11. Print & Reduced Motion

- **Print:** `.no-print` hidden, body background putih
- **Reduced motion:** Semua animasi dan transisi di-override menjadi `0.01ms`, marquee dihentikan, terminal cursor dihilangkan, navbar panel transition dinonaktifkan
