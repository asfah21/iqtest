# DESIGN.md — Warna, Tipografi, Radius, Shadow, Icon & Animation

> Fokus dokumen ini: bahasa visual (design tokens) yang digunakan di dalam project. Untuk struktur grid lihat [`LAYOUT.md`](context/LAYOUT.md), untuk aturan komponen lihat [`COMPONENTS.md`](context/COMPONENTS.md). Spesifikasi hero lihat [`HERO.md`](context/HERO.md), spesifikasi quiz lihat [`QUIZUI.md`](context/QUIZUI.md) dan [`quiz_header.md`](context/quiz_header.md).

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
| `--ink` | `#0f172a` | Teks body utama |
| `--inkMuted` | `#475569` | Teks sekunder, caption |
| `--inkSubtle` | `#94a3b8` | Teks tersier, placeholder |
| `--bord` | `#e2e8f0` | Border umum (card, divider) |
| `--bordLight` | `#f1f5f9` | Border halus untuk card |

### 1.3 Warna Data (Dark Triad)

Khusus untuk visualisasi data hasil tes di dashboard:

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--narcissus` | `#7c3aed` | Data domain ungu |
| `--machiavellian` | `#d97706` | Data domain amber |
| `--psychopath` | `#059669` | Data domain hijau |

### 1.4 Quiz Design Tokens (`assets/css/quiz.css`)

Quiz page memiliki design system terpisah (navy + pink). Token didefinisikan di `:root` dalam [`quiz.css`](assets/css/quiz.css:7):

| Token | Hex | Penggunaan |
|-------|-----|------------|
| `--qnavy` | `#16324F` | Warna header quiz, teks judul |
| `--qorange` | `#EC4899` | Aksen quiz (tombol, indikator selected) |
| `--qorange-hover` | `#f57bb8` | Hover state aksen quiz |
| `--qwhite` | `#FFFFFF` | Background card quiz |
| `--qbg` | `#F2F3F5` | Background halaman quiz |
| `--qgray-900` | `#333333` | Teks gelap quiz |
| `--qgray-500` | `#8A8F98` | Teks muted quiz |
| `--qgray-300` | `#D9DEE4` | Border quiz |
| `--qgray-light` | `#E8ECF0` | Background input quiz |
| `--qshadow` | `0 4px 12px rgba(0,0,0,0.06)` | Shadow card quiz |
| `--qradius-card` | `18px` | Radius card quiz |
| `--qradius-pill` | `9999px` | Radius pill button quiz |
| `--qfont` | `'Poppins', 'Nunito', system-ui, sans-serif` | Font quiz |

> **Catatan:** `--qorange` awalnya `#F5821F` (oranye) tetapi telah diubah ke `#EC4899` (pink). Hover state `--qorange-hover` juga berubah dari `#F7941D` ke `#f57bb8`.

### 1.5 Hero Colors (`assets/css/hero.css`)

Hero menggunakan hardcoded warna spesifik, tidak mereferensi CSS variables dari design system utama:

| Token (HERO.md) | Hex | Penggunaan |
|------------------|-----|------------|
| `--color-accent-blue` | `#4A5CF5` | Highlight teks headline |
| `--color-accent-blue-2` | `#5B6EF5` | CTA button gradient start, tombol "Next Question" |
| `--color-accent-purple` | `#6C5CE7` | CTA button gradient end |
| `--color-accent-yellow` | `#F5A623` | Sparkle dekorasi |
| `--color-accent-orange` | `#FB923C` | Kotak "?" di grid soal mockup |
| `--color-navy-header` | `#12294D` | Header bar di device mockup |
| `--color-text-dark` | `#111827` | Teks headline |
| `--color-text-secondary` | `#4B5563` | Teks subheadline, deskripsi |
| `--pastel-purple` | `#EDE9FE` | Background icon fitur "30 Questions" |
| `--pastel-green` | `#DCFCE7` | Background icon fitur "IQ Score" |
| `--pastel-orange` | `#FEF3C7` | Background icon fitur "Certificate" |
| `--pastel-blue` | `#DBEAFE` | Background icon fitur "7-Day Trial" |

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
- **Quiz page:** `'Poppins', 'Nunito', system-ui, sans-serif` (`--qfont`)

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
- **Navbar brand icon:** `border-radius: 12px`
- **Feature icon container:** `border-radius: 10px`
- **Quiz option:** `border-radius: 10px`
- **Quiz header:** `border-radius: 22px`
- **Quiz card:** `--qradius-card: 18px`
- **Quiz navigator:** `--qradius-card: 18px`
- **Paywall payment card:** `border-radius: 16px`
- **Paywall QRIS container:** `border-radius: 12px`
- **Step code badge:** `border-radius: 6px`

---

## 4. Shadow / Elevation

| Token | Value | Penggunaan |
|-------|-------|------------|
| `--shadow1` | `0 1px 2px rgba(15,23,42,0.04), 0 4px 12px rgba(15,23,42,0.06)` | Card default, navbar |
| `--shadow2` | `0 4px 16px rgba(15,23,42,0.08), 0 12px 40px rgba(15,23,42,0.06)` | Card hover, modal, dropdown, back-to-top button |

### Shadow Spesifik Komponen
- **CTA Hero:** `box-shadow: 0 12px 24px rgba(91,110,245,0.35)` ; hover: `0 16px 32px rgba(91,110,245,0.45)`
- **Navbar default:** `box-shadow: 0 2px 8px rgba(0,0,0,0.15)`
- **Navbar scrolled:** `box-shadow: 0 4px 20px rgba(0,0,0,0.08)`
- **Quiz shadow:** `--qshadow: 0 4px 12px rgba(0,0,0,0.06)`
- **Paywall payment card:** `box-shadow: 0 4px 24px rgba(15,23,42,0.08), 0 1px 2px rgba(15,23,42,0.04)` ; hover: `0 8px 40px rgba(15,23,42,0.1), 0 2px 4px rgba(15,23,42,0.06)`
- **Paywall CTA hover:** `box-shadow: 0 4px 16px rgba(99,102,241,0.3)`
- **Quiz button:** `box-shadow: 0 2px 6px rgba(245,130,31,0.3)` (menggunakan oranye meskipun tombol sekarang pink)
- **Paywall QRIS hover:** `box-shadow: 0 0 0 2px color-mix(in srgb, var(--accentLight) 60%, transparent)`
- **Paywall input focus:** `box-shadow: 0 0 0 3px color-mix(in srgb, var(--accentLight) 60%, transparent)`

### Backward Compatibility Aliases

| Token | Value |
|-------|-------|
| `.shadow-sm` | `0 1px 2px rgba(26,25,23,0.04)` |
| `.shadow-md` | `0 4px 6px rgba(26,25,23,0.06)` |
| `.shadow-xl` | `var(--shadow2)` |
| `.hover\:shadow-md:hover` | `0 4px 6px rgba(26,25,23,0.06)` |

---

## 5. Icon

- **Library:** Semua icon menggunakan **inline SVG** langsung di dalam HTML/templ — bukan icon library eksternal.
- **Style:** Stroke-based (outline), stroke-width bervariasi:
  - `1.5` — ikon dekoratif (navbar brand, footer, feature card)
  - `1.6` — ikon hero badge dan fitur
  - `1.8` — ikon fitur hero
  - `2.5` — ikon navigasi (navbar toggle, chevron, panel, CTA button)
- **Ukuran standar:**
  - Icon inline (dekat teks): `16px` atau `20px`
  - Icon di card fitur: `24px` (dalam container `44px × 44px`, dengan background `--accentLight`, color `--accent`)
  - Icon di Trust Bar: tidak menggunakan icon terpisah — hanya angka besar + label
  - Icon brand (navbar/footer): `24px` (dalam container `44px × 44px`, mobile `40px × 40px`)
  - Icon di quiz header: brand icon `44px × 44px`, timer icon `36px × 36px` (background `#EDEBFF`)
- **Warna icon:**
  - Default: mengikuti warna teks sekitarnya (`currentColor`)
  - Ikon fitur di hero: warna ungu `#7C3AED`, hijau `#16A34A`, amber `#F59E0B`, biru `#2563EB`
  - Ikon di panel mobile: `currentColor` (inherit dari link)
- Khusus untuk SVG dekoratif hero section, lihat [`HERO.md`](context/HERO.md) §6 untuk daftar lengkap placeholder SVG.

---

## 6. Ilustrasi & Gambar

- Hero menggunakan **multi-device mockup** (laptop + tablet + phone) dengan SVG inline, bukan foto stok generik. Lihat [`HERO.md`](context/HERO.md) §3.
- Ilustrasi "Tentang Tes" menggunakan placeholder `div` dengan background `--accentLight`, belum ada gambar final (tinggi `300px`).
- Semua gambar soal (`assets/images/q_*.svg`) adalah SVG pola matriks untuk quiz. Format: `q_{domain}_{seq}.svg` (mis. `q_mtx_001.svg`).
- Gambar opsi jawaban (`assets/images/opt_a.svg` s/d `opt_d2.svg`) adalah SVG pola untuk opsi A/B/C/D.
- Paywall order items menampilkan SVG placeholder (`32px`) untuk skor, sertifikat, dan IQ booster.

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
| Navbar panel slide-in | `0.3s` | `ease` |
| Quiz option selection | `0.2s` | `ease` |
| Quiz progress fill | `0.25s` | `ease` |
| Quiz step fade-in | `150ms` | `ease` (via Alpine.js `x-transition`) |
| Scroll reveal (fade-in-up) | `0.5s` | `ease` |
| Terminal cursor blink | `1s` | `step-end infinite` |
| Skeleton shimmer | `1.5s` | `ease-in-out infinite` |
| Spinner rotate | `1s` | `linear infinite` |
| Paywall card hover | `0.3s` | `ease` |
| Paywall CTA hover | `0.2s` | `ease` |
| Quiz navigator item hover | `0.15s` | `ease` |

### 7.2 Aturan Animasi
- **Hover button/card:** transisi `background-color`, `box-shadow`, `transform: translateY(-1px)` dengan durasi `0.2s ease`.
- **Accordion FAQ:** expand/collapse `max-height` dengan durasi `0.3s ease`. Icon rotate `45deg` saat expanded.
- **Scroll-reveal:** `.fade-in-up` + `.is-visible` menggunakan keyframe `fadeInUp` (opacity `0→1`, `translateY(12px → 0)`), durasi `0.5s`. Variasi `.fade-in-up-delayed` dengan delay `0.15s`.
- **Navbar:** transisi `box-shadow` dan `padding` saat scroll, durasi `0.3s ease`. Floating navbar dengan `top: 16px`, menyusut ke `padding: 6px 14px` saat scrolled.
- **Back-to-top button:** opacity + visibility + transform (`translateY(8px → 0)`), durasi `0.3s ease`. Muncul setelah scroll melewati 1 viewport.
- **Mobile slide panel:** transform `translateX` dengan durasi `0.3s ease`, overlay fade.
- **Quiz transitions:** Fade-in antar step (Alpine.js `x-transition`, 150ms), progress bar smooth (`0.25s ease`), option hover state (`0.15s ease`).
- **CTA button:** hover menambah `transform: translateY(-1px)`.
- **Paywall:** fade-in-up dengan delay bertahap per section (`0s`, `0.08s`, `0.16s`, `0.24s`, `0.32s`, `0.40s`, `0.48s`). Payment card hover `translateY(-2px)`.
- Hindari animasi berlebihan — semua animasi `@media (prefers-reduced-motion: reduce)` di-reset ke `0.01ms`.

### 7.3 Keyframes Defined

| Keyframe | Penggunaan | File |
|----------|------------|------|
| `fadeInUp` | Scroll reveal sections | [`style.css`](assets/css/style.css) |
| `blink` | Terminal cursor (jika digunakan) | [`style.css`](assets/css/style.css) |
| `shimmer` | Skeleton loading state quiz | [`quiz.css`](assets/css/quiz.css:502) |
| `spinner-rotate` | Submitting overlay spinner | [`quiz.css`](assets/css/quiz.css:686) |
| `spin` | General spinner (admin/paywall) | [`style.css`](assets/css/style.css:1246) |

---

## 8. Warna Komponen Spesifik

### 8.1 Navbar
- Background: `#FFFFFF`
- Shadow: `0 2px 8px rgba(0,0,0,0.15)` (default), `0 4px 20px rgba(0,0,0,0.08)` (scrolled)
- Brand icon background: `linear-gradient(135deg, #7C6FF0, #5B4FE0)`
- Brand icon size: `44px × 44px` (desktop), `40px × 40px` (mobile)
- Brand icon border-radius: `12px`
- Link text: `#475569` (default), `#0f172a` (hover)
- Link hover background: `#f1f5f9`
- CTA button: background `#6366f1`, hover `#4f46e5`, border-radius `10px`
- Navbar border-radius: `22px`
- Navbar padding: `10px 14px` (default), `6px 14px` (scrolled)
- Mobile: `top: 12px`, `width: calc(100% - 24px)`, `padding: 8px 12px`

### 8.2 Hero (lihat [`HERO.md`](context/HERO.md) untuk detail lengkap)
- CTA button: `linear-gradient(135deg, #5B6EF5, #6C5CE7)`
- Badge border: `#ECECF3`
- Disclaimer background: `#F1F3F9`
- Headline: `60px` (desktop), `48px` (1024px), `36px` (768px), font `'Poppins', 'Inter', system-ui, sans-serif`

### 8.3 Footer
- Background: `#1a1917`
- Teks heading: `rgba(255,255,255,0.7)`
- Teks link: `rgba(255,255,255,0.4)` (default), `var(--accent)` (hover)
- Teks tagline: `rgba(255,255,255,0.45)`
- Teks copyright: `rgba(255,255,255,0.3)`
- Divider: `rgba(255,255,255,0.08)`
- Logo icon: `linear-gradient(135deg, #7C6FF0, #5B4FE0)` (sama dengan navbar)
- Background symbol: SVG logo besar, opacity `0.025`, posisi `bottom: -8%; right: -4%`

### 8.4 Footer Mobile
- Teks link lebih terang: `rgba(255,255,255,0.55)` untuk contrast lebih baik
- Brand section: `text-align: center`

---

## 9. Paywall Design Tokens (`assets/css/paywall.css`)

Paywall page menggunakan tokens dari design system utama (`--accent`, `--ink`, `--bg`, dll.) ditambah kelas-kelas spesifik. Lihat [`paywall.css`](assets/css/paywall.css).

### 9.1 Komponen Paywall

| Komponen | Kelas CSS | Keterangan |
|----------|-----------|------------|
| Completion Banner | `.paywall-completion-banner` | Background `--accentLight`, radius `--radiusCard` |
| Banner Icon | `.paywall-banner-icon` | `36px × 36px`, background `--surface`, shadow accent |
| Well-done Text | `.paywall-well-done` | Center text, `max-width: --measure-narrow` |
| Order Heading | `.paywall-order-heading` | `--text-sm`, bold, uppercase |
| Order Item | `.paywall-order-item` | Background `--surface`, border `--bordLight`, radius `--radiusCard`, hover shadow `--shadow1` |
| Order Number | `.paywall-order-number` | `36px × 36px` circle, background `--accent` (item 1), `--narcissus` (item 2), `--machiavellian` (item 3) |
| Order Image | `.paywall-order-image` | `72px × 72px`, background `--bgSection`, radius `--radiusBtn` |
| Order Total | `.paywall-order-total` | Flex row, background `--surface`, border `--bord` |
| Payment Card | `.paywall-payment-card` | Background `--surface`, border `--bord`, radius `16px`, shadow custom |
| Price Zone | `.paywall-price-zone` | Center text, padding `--space-xl` |
| Price Badge | `.paywall-price-badge` | Pill badge, background `--accentLight`, color `--accent` |
| Price Amount | `.paywall-price-amount` | `clamp(2.75rem, 6vw, 3.75rem)`, weight 800 |
| Divider | `.paywall-divider` | `1px solid --bordLight` |
| QRIS Container | `.paywall-qris` | Gradient background `color-mix` accent, border `--bordLight`, radius `12px` |
| QRIS Placeholder | `.paywall-qris-placeholder` | `176px × 176px` (mobile `140px`), dashed border |
| Trust Signal | `.paywall-trust` | Flex center, background `--bgSection`, pill `999px` |
| Input | `.paywall-input` | Radius `--radiusBtn`, border `1.5px solid --bordLight`, focus accent |
| CTA Button | `.paywall-cta` | Background `--accent`, white text, radius `--radiusBtn`, hover `--accentHover` + `translateY(-2px)` |
| Payment Methods | `.paywall-payment-methods` | Flex wrap, gap `--space-sm` |
| Method Icon | `.paywall-payment-method-icon` | Background `--surface`, border `--bordLight`, radius `--radiusBtn` |
| Agreement | `.paywall-agreement` | `--text-xs`, checkbox accent |
| Disclaimer | `.paywall-disclaimer` | Background `--bgSection`, radius `--radiusCard`, border `--bordLight` |
| Reviews Heading | `.paywall-reviews-heading` | Center text, rating stars `#f59e0b` |
| Review Card | `.paywall-review-card` | Background `--surface`, border `--bordLight`, radius `--radiusCard` |
| Bottom CTA | `.paywall-bottom-cta` | Border top `--bordLight` |
| FAQ Mini | `.paywall-faq` | Background `--surface`, border `--bordLight`, radius `--radiusCard` |

---

## 10. Backward Compatibility

CSS mendefinisikan alias untuk variable lama di `:root` kedua (sekitar baris 1171 di [`style.css`](assets/css/style.css)):

| Alias | Maps to |
|-------|---------|
| `--textMain` | `var(--ink)` |
| `--textMuted` | `var(--inkMuted)` |
| `--textSubtle` | `var(--inkSubtle)` |
| `--warmBg` | `var(--bg)` |
| `--warmBord` | `var(--bordLight)` |
| `--bgAlt` | `#e8e4de` |

Serta color scale shorthand (semua indigo-based):
- `--accent-50: #eef2ff`
- `--accent-200: #c7d2fe`
- `--accent-500: #6366f1`
- `--accent-600: #4f46e5`

Status aliases:
- `--success: #059669`
- `--warning: #d97706`
- `--error: #dc2626`

---

## 11. Selection & Scrollbar

- **Selection highlight:** `background: var(--accentLight); color: var(--ink)`
- **Scrollbar:** width `6px`, thumb `var(--bord)` dengan `border-radius: 999px`, track transparan
- **Focus ring:** `outline: 2px solid var(--accent); outline-offset: 2px` via `:focus-visible`

---

## 12. Print & Reduced Motion

- **Print:** `.no-print` hidden, body background putih
- **Reduced motion:** Semua animasi dan transisi di-override menjadi `0.01ms`, mobile panel transition dinonaktifkan, skeleton shimmer dihentikan

---

## 13. CSS File Map

| File | Isi |
|------|-----|
| [`assets/css/style.css`](assets/css/style.css) | Design tokens utama, navbar, container, sections, trust bar, features, steps, pricing, FAQ, CTA, back-to-top, animations, utilities, backward compat |
| [`assets/css/hero.css`](assets/css/hero.css) | Hero section — wrapper, badge, headline, CTA, features row, device mockup, responsive |
| [`assets/css/footer.css`](assets/css/footer.css) | Footer — grid, brand, links, copyright, responsive, background symbol |
| [`assets/css/quiz.css`](assets/css/quiz.css) | Quiz design tokens, header, body layout, cards, matrix, answer grid, option states, navigator, skeleton, identity step, finish modal |
| [`assets/css/paywall.css`](assets/css/paywall.css) | Paywall — completion banner, well-done, order details, payment card, QRIS, trust, form, reviews, disclaimer, bottom CTA, FAQ mini |
| [`assets/css/mobile-slide-panel.css`](assets/css/mobile-slide-panel.css) | Mobile navbar panel — overlay, slide panel, header, nav links, CTA, responsive |
