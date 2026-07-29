# LAYOUT.md — Grid, Container, Spacing & Responsive

> Fokus dokumen ini: struktur layout halaman (bukan warna/komponen). Untuk styling visual lihat [`DESIGN.md`](context/DESIGN.md), untuk aturan komponen lihat [`COMPONENTS.md`](context/COMPONENTS.md). Spesifikasi hero lihat [`HERO.md`](context/HERO.md), spesifikasi quiz lihat [`QUIZUI.md`](context/QUIZUI.md).

---

## 1. Breakpoints

| Nama | Lebar Min | Keterangan |
|------|-----------|------------|
| `mobile` | 0px | Default, stack vertikal penuh |
| `sm` | 640px | Transisi kecil (flex-row dsb.) |
| `tablet` | 768px | Mulai 2 kolom di beberapa section |
| `desktop` | 1024px | Layout penuh multi-kolom |
| `wide` | 1200px | Container max-width bertambah ke 1290px |

---

## 2. Container

Container adalah wrapper tunggal yang dipakai berulang (`.container`) di **setiap section pada setiap halaman** — Hero, semua section, dan Footer. Tidak ada section yang punya max-width sendiri di luar aturan ini, supaya lebar semua bagian benar-benar sejajar vertikal dari atas sampai bawah halaman.

- **Max-width default:** `1200px`
- **Max-width pada layar lebar:** `1290px`, aktif mulai `min-width: 1200px`

```css
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px; /* mobile default */
}

@media (min-width: 768px) {
  .container { padding: 0 24px; }
}

@media (min-width: 1024px) {
  .container { padding: 0 40px; }
}

@media (min-width: 1200px) {
  .container { max-width: 1290px; }
}
```

### 2.1 Aturan Wajib Konsistensi
- **Navbar:** Navbar adalah pengecualian — ia **tidak** dibungkus `.container`. Navbar menggunakan wrapper sendiri (`.navbar-inner`) dengan `max-width: 1204px` dan `padding: 10px 14px`, posisi `fixed` dengan `top: 16px` dan `border-radius: 22px` (floating pill style). Background navbar `#FFFFFF` full-width.
- **Hero:** Hero section `padding` menggunakan nilai hardcoded (`80px 0`), bukan `.container` padding. Konten hero dibungkus `.container` di dalam `.hero-wrapper`.
- **Semua component section lainnya** (Trust Bar, Features, Tentang Tes, Cara Kerja, Pricing, CTA Banner, FAQ): semua memakai `.container` yang identik.
- **Footer:** Kolom footer dibungkus `.container` yang sama (`.footer-inner`), walau background footer (`#1a1917`) full-bleed selebar layar.
- **Quiz page:** Menggunakan `.container` standar yang sama, dibungkus dalam `.quiz-page` section.
- **Paywall page:** Menggunakan `.container` standar, dengan inner wrapper `max-width: 640px; margin: 0 auto` untuk konten.
- Jika suatu saat butuh section yang benar-benar full-width tanpa batas (mis. banner promo penuh warna), section tersebut tetap membungkus **konten teksnya** dengan `.container` di dalamnya — hanya background yang full-bleed.

### 2.2 Content Top Spacing
Karena navbar adalah floating fixed (`position: fixed; top: 16px`), konten halaman publik perlu padding-top:
```css
.public-content {
  padding-top: 88px;   /* desktop */
  background: #FFFFFF;
}

@media (max-width: 767px) {
  .public-content {
    padding-top: 80px; /* mobile */
    background: #FFFFFF;
  }
}
```

---

## 3. Grid System

- Gunakan CSS Grid / Flexbox dengan basis **12-column grid** untuk section yang butuh kolom.
- **Gutter (jarak antar kolom):** `24px` desktop, `16px` mobile.
- Section bebas grid (hero, section teks panjang) boleh pakai flexbox 2 kolom sederhana (45/55 atau 1/1).

---

## 4. Spacing Scale

Skala konsisten (kelipatan 4px) untuk margin/padding antar elemen & antar section:

```
--space-2xs: 4px    (jarak antar elemen kecil, mis. icon-teks)
--space-xs:  8px
--space-sm:  12px
--space-md:  16px   (padding dalam card/button, container mobile)
--space-lg:  24px   (gutter, gap card)
--space-xl:  32px   (padding dalam pricing card, jarak section heading)
--space-2xl: 48px   (jarak antar blok dalam satu section, section padding mobile)
--space-3xl: 80px   (jarak section heading ke konten, padding section mobile override)
--space-4xl: 120px  (jarak antar section besar, wide screen)
--space-5xl: 96px   (section padding desktop)
--space-6xl: 128px  (hero vertical padding)
```

### 4.1 Jarak Antar Section (Vertical Rhythm)
- **Mobile:** `48px` (`--space-2xl`), namun di-override di media query `max-width: 767px` menjadi `var(--space-3xl)` untuk `.section` dan `.section-alt`
- **Tablet:** `64px`
- **Desktop:** `96px` (`--space-5xl`)

```css
.section {
  padding: var(--space-2xl) 0;    /* 48px mobile */
}
.section-alt {
  padding: var(--space-2xl) 0;    /* 48px mobile */
  background: var(--bgSection);
}

@media (min-width: 768px) {
  .section, .section-alt { padding: 64px 0; }
}

@media (min-width: 1024px) {
  .section, .section-alt { padding: var(--space-5xl) 0; }  /* 96px desktop */
}
```

---

## 5. Struktur Halaman (Urutan Section Aktual)

Urutan yang di-render di [`index_page.templ`](templ/pages/index_page.templ):

```
├── Header / Navbar (floating, fixed, border-radius: 22px)
├── Hero Section (wrapper putih, multi-device mockup)
├── Trust Bar (statistik singkat, 4 item)
├── Features Grid / "Mengapa Tes Ini?" (6 card, 3×2 grid)
├── Section "Tentang Tes" (2 kolom: ilustrasi placeholder + teks)
├── Section "Cara Kerja" (3 langkah, angka 01/02/03)
├── Section Pricing (single card center)
├── CTA Banner (background `--bgSection`, centered CTA)
├── FAQ (accordion, max-width 800px)
└── Footer (5 kolom desktop, background #1a1917)
```

---

## 6. Layout per Section

### 6.1 Navbar (Header)
- **Style:** Floating pill — `position: fixed; top: 16px; max-width: 1204px; border-radius: 22px; background: #FFFFFF; box-shadow: 0 2px 8px rgba(0,0,0,0.15)`
- **Desktop:** 1 baris flex, `justify-content: space-between` dalam `.navbar-inner`:
  - Kiri: Logo (icon 44×44 + teks "IQ Test")
  - Tengah: 3 link navigasi (Beranda, Fitur, Tes)
  - Kanan: CTA button "Mulai Tes" + hamburger toggle (hidden di desktop)
- **Scrolled state:** Shadow mengecil (`0 4px 20px rgba(0,0,0,0.08)`), padding menyusut dari `10px 14px` ke `6px 14px`
- **Mobile (≤767px):** Navbar menyempit (`width: calc(100% - 24px)`, `top: 12px`), link navigasi dan CTA disembunyikan, hamburger toggle ditampilkan. Klik toggle membuka slide-in panel dari kanan.
- **Mobile slide panel:** Overlay gelap + panel dari kanan (lebar `80%`, max `320px`), berisi: close button (X) di pojok kanan atas → link navigasi vertikal → divider → CTA "Mulai Tes" di bagian bawah. Title "Menu" di-hidden via CSS (`display: none`). Lihat [`mobile-slide-panel.css`](assets/css/mobile-slide-panel.css).
- **Brand icon:** `linear-gradient(135deg, #7C6FF0, #5B4FE0)`, `border-radius: 12px`, ukuran `44px` (mobile `40px`)

### 6.2 Hero
- **Wrapper:** `.hero-wrapper` — background `#FFFFFF`, padding `80px 0`
- **Desktop:** Flexbox `justify-content: space-between`, `align-items: center`, gap `64px`
  - Kiri (`.hero-left`): `flex: 0 0 calc(45% - 32px)` — badge → headline → subheadline → CTA → 4 fitur → disclaimer
  - Kanan (`.hero-right`): `flex: 0 0 calc(55% - 32px)` — multi-device stack (laptop + tablet + phone)
- **Tablet (≤1024px):** Gap mengecil ke `40px`, headline font-size turun ke `48px`
- **Mobile (≤768px):** Stack vertikal — hero-left dan hero-right masing-masing `flex: 0 0 100%; max-width: 100%`. Headline font-size mengecil ke `36px`. Hero features berubah jadi grid `1fr 1fr`. Device stack height mengecil ke `400px`.
- **Detail lengkap:** Lihat [`HERO.md`](context/HERO.md)

### 6.3 Trust Bar
- **Desktop:** Flex row, `justify-content: center`, `align-items: center`, gap `--space-xl` (32px)
- **Mobile (≤767px):** Grid `1fr 1fr`, gap `--space-lg` (24px)
- **Anatomi per item:** Angka besar (`--text-3xl`, weight 800) + label kecil (`--text-sm`, color `--inkMuted`)
- **4 item:** 8M+ Users Worldwide, 30 Soal Bergambar, 4 Domain Kognitif, Rp14.900 Hasil Lengkap

### 6.4 Features Grid / "Mengapa Tes Ini?"
- **Section:** `.section-alt` (background `--bgSection`), dengan `.section-heading` terpusat di atas
- **Grid:** 3 kolom (`grid-template-columns: repeat(3, 1fr)`), gap `--space-lg` (24px)
- **Feature card:** `.feature-card` — background `--surface`, border `1px solid var(--bordLight)`, `border-radius: var(--radiusCard)`, padding `--space-xl` (32px), shadow `--shadow1`. Hover: border `--bord`, shadow `--shadow2`
- **Anatomi card:** Feature icon (44×44, background `--accentLight`, color `--accent`, `border-radius: 10px`) → h3 judul (`--text-base`, weight 600) → p deskripsi (`--text-sm`, color `--inkMuted`)
- **Mobile (≤767px):** 1 kolom
- **Di atas grid:** Paragraf penjelasan (`.why-text`) dengan `max-width: 700px`, center-aligned, teks `--inkMuted`
- **6 cards:** Berbasis Ilmiah, Non-Verbal, Anti-Cheat, 4 Domain, Bobot Kesulitan, Freemium

### 6.5 Section "Tentang Tes"
- **Desktop:** Grid 2 kolom `grid-template-columns: 45fr 55fr`, gap `--space-3xl` (80px), `align-items: center`
  - Kiri: Gambar/ilustrasi (placeholder `div` dengan background `--accentLight`, height `300px`, `border-radius: var(--radiusCard)`)
  - Kanan: 2 paragraf teks (`--inkMuted`, `--leading-relaxed`)
- **Mobile (≤767px):** 1 kolom stack. Gambar `order: 2` (di bawah teks).

### 6.6 Section "Cara Kerja" (How It Works)
- **Section:** `.section-alt` (background `--bgSection`)
- **Desktop:** Grid 3 kolom `grid-template-columns: repeat(3, 1fr)`, gap `--space-xl` (32px)
- **Dekorasi:** Garis horizontal `::before` di `.steps-container` menghubungkan ketiga step di bagian atas (`height: 1px`, background `--bordLight`, posisi `top: 48px`)
- **Anatomi step card (`.step-card`):**
  1. Angka besar "01"/"02"/"03" (`.step-number`): `font-size: clamp(3rem, 5vw, 4rem)`, weight 800, color `--bordLight`, `letter-spacing: -0.04em`
  2. Judul step (`h3`): `--text-base`, weight 600
  3. Deskripsi (`p`): `--text-sm`, color `--inkMuted`
  4. Kode/estimasi waktu (`.step-code`): font `JetBrains Mono`, `--text-xs`, background `--bgSection`, color `--accent`, `border-radius: 6px`, border `1px solid var(--bordLight)`
- **Mobile (≤767px):** 1 kolom stack, gap `--space-2xl`, garis dekorasi di-`display: none`

### 6.7 Section Pricing
- **Section:** `.section` (background `--bg`)
- **Layout:** Single card, center-aligned via `.pricing-card-wrapper` (flex, `justify-content: center`)
- **Card (`.pricing-card`):** `max-width: 480px`, width `100%`, background `--surface`, border `1px solid var(--bordLight)`, `border-radius: var(--radiusCard)`, padding `--space-xl` (mobile `--space-lg`), shadow `--shadow1`, `text-align: center`
- **Anatomi card:** Nama paket → Harga besar (color `--accent`) → Deskripsi → CTA button (`.cta-btn`)

### 6.8 CTA Banner
- **Section:** `.section-cta` — background `--bgSection`, padding `--space-5xl` 0, `text-align: center`
- **Anatomi:** Headline (`--text-2xl`, weight 700) → 3 benefit checks (✓ icon + teks) → CTA button besar (`.cta-btn`)
- **CTA button (`.cta-btn`):** `padding: 14px 36px`, `border-radius: var(--radiusBtn)`, background `--accent`, color white, weight 600, `font-size: 15px`. Hover: `--accentHover` + `translateY(-1px)`. Mobile: `width: 100%`

### 6.9 FAQ
- **Section:** Di-render via `@components.FaqSection()`
- **Container:** `.faq-accordion` — `max-width: 800px`, `margin: 0 auto`
- **Style:** List-style (bukan card), divider `1px solid var(--bordLight)` antar item
- **Accordion behavior:** Expand/collapse `max-height` dengan transisi `0.3s ease`. Trigger button: `justify-content: space-between`, icon `+` rotate `45deg` saat expanded. Hover trigger: color → `--accent`

### 6.10 Footer
- **Background:** `#1a1917` (gelap, full-bleed)
- **Dekorasi:** SVG logo besar sebagai background symbol (opacity `0.025`, posisi `bottom: -8%; right: -4%`), hidden di mobile
- **Grid desktop:** 5 kolom `grid-template-columns: repeat(5, 1fr)`, gap `--space-xl` (32px)
  - Kolom 1: **Brand** — logo icon + "IQ Test" + tagline
  - Kolom 2-5: Dibungkus `.footer-links-grid` (`display: contents` di desktop)
    - **Navigasi:** Beranda, Fitur, Mulai Tes
    - **Informasi:** Privasi, Syarat & Ketentuan, Kontak
    - **Sumber Daya:** Dokumentasi, API, Blog
    - **Ikuti Kami:** Twitter, Instagram, LinkedIn
- **Style link:** Default `rgba(255,255,255,0.4)`, hover `var(--accent)` + `padding-left: 4px`
- **Heading kolom:** `--text-sm`, weight 700, `rgba(255,255,255,0.7)`
- **Copyright:** `margin-top: var(--space-2xl)`, `padding-top: var(--space-lg)`, `border-top: 1px solid rgba(255,255,255,0.08)`, teks `rgba(255,255,255,0.3)`, `--text-sm`
- **Tablet (768-1023px):** 2 kolom grid
- **Mobile (≤767px):** Stack 1 kolom. `.footer-links-grid` berubah jadi 2 kolom grid (`1fr 1fr`). Brand section `text-align: center`. Padding footer mengecil ke `--space-2xl`. Background symbol di-`display: none`.

---

## 7. Aturan Responsive Umum

- Semua CTA button: `width: auto` di desktop, `width: 100%` di mobile (`≤767px`)
- Gambar/ilustrasi: `max-width: 100%; height: auto` — tidak ada gambar dengan lebar fixed px yang bisa overflow di mobile
- Navbar: Floating fixed di semua breakpoint. Mobile: link + CTA disembunyikan, toggle muncul.
- Back-to-top button: Muncul setelah scroll melewati tinggi 1 viewport, posisi `fixed`, `bottom: 24px; right: 24px` (mobile: `bottom: 16px; right: 16px`). Style: bulat `44px` (mobile `40px`), background `--accent`, shadow `--shadow2`, opacity/visibility transition `0.3s ease`.
- Hero visual (device mockup): Stack vertikal di mobile (tidak di-hide, tetap ditampilkan di bawah hero-left).
- Fitur grid: 3 kolom desktop → 1 kolom mobile.
- Steps container: 3 kolom desktop → 1 kolom mobile, garis dekorasi hidden.
- Reduced motion: Semua animasi di-override ke `0.01ms` via `@media (prefers-reduced-motion: reduce)`.
- Quiz page memiliki responsive behavior sendiri (lihat §8.1).
- Paywall page: `max-width: 640px` inner wrapper, padding adjustments di mobile.

---

## 8. Halaman Khusus

### 8.1 Quiz Page ([`quiz_page.templ`](templ/pages/quiz_page.templ))

- **Layout:** [`QuizLayout`](templ/layouts/quiz_layout.templ) — full-height, background `--qbg`, font `--qfont`
- **Container:** `.container` standar di dalam `.quiz-page` section
- **States:** `loading` → `question-active` → `identity` → `submitting` → `completed`

#### 8.1.1 Quiz Header
- **Container:** `.quiz-header` — background `--qwhite`, `border-radius: 22px`, `padding: 10px 14px`, flex `space-between`, `margin-bottom: 20px`, shadow `0 2px 8px rgba(0,0,0,0.04)`
- **4 sections (desktop, flex row):**
  1. **Info Tes** (`.quiz-header-brand`): icon `44×44` (gradient ungu `#7C6FF0 → #5B4FE0`, `border-radius: 12px`) + teks "IQ Test" / "Ukur kemampuan kognitifmu"
  2. **Progress Tes** (`.quiz-header-progress`): label "Progress Tes" + progress bar (`#E8E8F0` track, `#5B4FE0` fill, `height: 7px`) + counter teks
  3. **Waktu Tersisa** (`.quiz-timer-card`): card kecil dengan border `#D9D9E6`, icon timer (`#EDEBFF` bg) + waktu `MM:SS`
  4. **Tombol Aksi** (`.quiz-header-actions`): "Akhiri Tes" button (border `#D9D9E6`, hover merah)
- **Mobile (≤767px):** Header wraps, brand hidden, progress full-width, timer + end button aligned right

#### 8.1.2 Quiz Body
- **Layout:** `.quiz-body` — flex row (desktop), flex column (mobile), gap `20px`
- **Desktop:** Question card (`flex: 2`) + Answer card (`flex: 1.2`)
- **Mobile:** Stacked vertically

#### 8.1.3 Question Card
- **Container:** `.quiz-card` — background `--qwhite`, `border-radius: 18px`, padding `24px 20px`, shadow `--qshadow`
- **Anatomi:** Title "Bentuk manakah yang hilang?" → Divider (`2px solid --qgray-300`) → Matrix image (`max-width: 320px`, `max-height: 320px`, center)
- **Transition:** Alpine.js fade (150ms) antar soal

#### 8.1.4 Answer Card
- **Container:** `.quiz-card` (sama dengan question card)
- **Anatomi:** Title "Pilih jawaban:" → Divider → Answer grid (2 kolom × 2 baris, 4 opsi A-D)
- **Answer option (`.answer-option`):** Flex row, label huruf (bold, `--qnavy`) + shape image, `border-radius: 10px`, `border: 1px solid #e2e8f0`. Selected: `border: 2px solid --qorange` + pink background

#### 8.1.5 Action Buttons
- **Container:** `.quiz-actions` — flex row, center, gap `16px`
- **Buttons:** `.btn-quiz-action` — background `--qorange`, white text, pill `9999px`, padding `10px 28px`, shadow orange
- **Dua tombol:** "Kembali" + "Lewati Pertanyaan"

#### 8.1.6 Question Navigator
- **Container:** `.quiz-navigator` — background `#5b4fe0`, `border-radius: 18px`, padding `16px 20px`
- **Grid:** `.nav-grid` — flex wrap, gap `12px`, center
- **Nav item (`.nav-item`):** Circle `32px`, white text `19px`, weight 600. Active: white circle bg + navy text. Answered: superscript huruf jawaban (pink `#ffa5d1`)

#### 8.1.7 Identity Step (Post-Quiz)
- **Container:** `.identity-card` — `max-width: 520px`, center, quiz card style
- **Anatomi:** Heading "Simpan Hasil Tes Anda" → Divider → Summary (✓ + "Anda telah menjawab X dari Y soal") → Description → Form (Nama + Email inputs) → Submit button
- **Input fields:** `border-radius: 10px`, border `--qgray-300`, focus `--qorange`

#### 8.1.8 Finish Confirmation Modal
- **Overlay:** `rgba(0, 0, 0, 0.45)`, fixed, full screen
- **Card:** `max-width: 400px`, `border-radius: 22px`, padding `32px 28px`, white bg, shadow
- **Stats row:** Answered (green) vs Unanswered (red) dengan divider

#### 8.1.9 Submitting Overlay
- **Overlay:** `rgba(255,255,255,0.85)`, fixed, full screen, flex center
- **Spinner:** SVG rotating circle + "Menyimpan jawaban..."

### 8.2 Paywall Page ([`paywall_page.templ`](templ/pages/paywall_page.templ))

- **Layout:** [`PaywallLayout`](templ/layouts/paywall_layout.templ)
- **Section:** `.paywall-section` — `min-height: 100vh`, `padding-top: 96px` (mobile `80px`), background `--bg`
- **Inner container:** `.container` > `max-width: 640px; margin: 0 auto`
- **Style:** Lihat [`paywall.css`](assets/css/paywall.css) dan [`DESIGN.md`](context/DESIGN.md) §9

#### 8.2.1 Section Order
```
├── Completion Banner (waktu penyelesaian + domain terkuat)
├── Well-done Text (ajakan beli hasil)
├── Order Details (3 item: Skor IQ, Sertifikat, IQ Booster)
├── Payment Card
│   ├── Order Total (Rp9.900)
│   ├── Divider
│   └── Action Zone
│       ├── QRIS Container (placeholder QR code)
│       ├── Payment Methods (BCA, Mandiri, BNI, BRI, GoPay, Dana, OVO, LinkAja)
│       ├── Trust Signal (🔒 Pembayaran aman & terenkripsi)
│       ├── Agreement Checkbox (Syarat & Ketentuan + trial auto-renew)
│       ├── Name Input + CTA Button ("Lanjutkan ke Pembayaran")
│       └── Form Note
├── Disclaimer (subscription terms detail)
├── Customer Reviews (3 testimonials + rating summary)
└── Bottom CTA (duplicate pay button)
```

#### 8.2.2 Payment Card
- **Container:** `.paywall-payment-card` — `border-radius: 16px`, border `--bord`, shadow custom
- **Hover:** Shadow elevates + `translateY(-2px)`, transition `0.3s ease`

#### 8.2.3 QRIS Container
- **Container:** `.paywall-qris` — gradient accent background, border `--bordLight`, `border-radius: 12px`
- **Placeholder:** `176px × 176px` (mobile `140px`), dashed border `--bordLight`
- **Hover:** Border accent + solid, glow ring via `box-shadow`

### 8.3 Hasil Page ([`hasil_page.templ`](templ/pages/hasil_page.templ))
- Menampilkan hasil tes setelah pembayaran

### 8.4 Dashboard Page ([`dashboard_page.templ`](templ/pages/dashboard_page.templ))
- Admin dashboard dengan tabel data

---

## 9. Back-to-Top Button

```css
.back-to-top {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 40;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--accent);
  color: white;
  /* default: opacity 0, visibility hidden, translateY(8px) */
  transition: opacity 0.3s ease, visibility 0.3s ease, transform 0.3s ease, background 0.2s ease;
}

.back-to-top.visible {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

@media (max-width: 767px) {
  .back-to-top {
    bottom: 16px;
    right: 16px;
    width: 40px;
    height: 40px;
  }
}
```

---

## 10. CSS File Map (Layout-Related)

| File | Layout Content |
|------|---------------|
| [`assets/css/style.css`](assets/css/style.css) | Container, navbar, sections, trust bar, features grid, about grid, steps, pricing, FAQ, CTA banner, back-to-top, responsive overrides |
| [`assets/css/hero.css`](assets/css/hero.css) | Hero wrapper, hero-inner (flex 45/55), hero-left, hero-right, device stack, responsive |
| [`assets/css/footer.css`](assets/css/footer.css) | Footer grid (5 col → 2 col → 1 col), footer-links-grid, brand, copyright, responsive |
| [`assets/css/quiz.css`](assets/css/quiz.css) | Quiz page, header, body (flex row/col), cards, answer grid (2×2), navigator, identity, modal, responsive |
| [`assets/css/paywall.css`](assets/css/paywall.css) | Paywall section, order details, payment card, QRIS, form, reviews, responsive |
| [`assets/css/mobile-slide-panel.css`](assets/css/mobile-slide-panel.css) | Navbar overlay, slide panel, panel header, nav links, responsive |
