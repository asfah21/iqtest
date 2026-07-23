# HERO SECTION SPEC — IQ Test Landing Page

Dokumen ini adalah spesifikasi teknis untuk membangun hero section landing page tes IQ. Tulis kode HTML/CSS (atau React) yang PERSIS mengikuti struktur, teks, warna, dan ukuran di bawah ini. Semua gambar/ikon menggunakan SVG inline placeholder yang sudah disediakan di Bagian 6 — jangan mencari file gambar eksternal, langsung pakai SVG yang ada.

---

## 1. Struktur Section

```
<section class="hero">
  <div class="hero-left">   <!-- konten teks, ~45% lebar -->
  <div class="hero-right">  <!-- mockup device, ~55% lebar -->
</section>
```

- Display: flex, `justify-content: space-between`, `align-items: center`.
- Padding section: `80px 64px` (desktop), `24px 16px` (mobile, stack vertikal — hero-right pindah ke bawah hero-left).
- Background section: `#FFFFFF`.
- `hero-right` posisinya `position: relative` karena berisi elemen dekoratif absolute (blob, dot grid, sparkle).
- Max-width container: `1440px`, margin auto.

---

## 2. HERO-LEFT — Konten & Urutan Elemen

### 2.1 Badge sosial proof
```html
<div class="badge">
  [ICON: avatar-group.svg]
  <span>Joined by 8+ million people worldwide</span>
</div>
```
- Style: `display:flex; align-items:center; gap:8px; background:#FFFFFF; border:1px solid #ECECF3; border-radius:999px; padding:8px 16px 8px 8px; box-shadow:0 2px 6px rgba(0,0,0,0.04); width:fit-content;`
- Icon di dalam lingkaran background `#EDEBFB`, diameter `36px`.
- Teks: `font-size:14px; font-weight:500; color:#1A1A2E;`

### 2.2 Headline (WAJIB PERSIS)
```html
<h1 class="headline">
  Discover Your<br>
  True<span class="quote">'</span><span class="highlight">Intelligence</span>
  [ICON: sparkle.svg — posisi absolute, di kanan atas kata "Intelligence"]
</h1>
```
- Teks baris 1 & "True": `color:#111827; font-weight:800;`
- `.highlight` ("Intelligence"): `color:#4A5CF5; font-weight:800;`
- `font-size:60px; line-height:1.05; letter-spacing:-1px;` (desktop). Mobile: `36px`.
- `.quote` (karakter `'`) hanya dekorasi pemisah, warna sama dengan teks gelap.
- Font-family: `'Poppins', 'Inter', sans-serif`.

### 2.3 Subheadline
```html
<p class="subheadline">
  The average IQ in Indonesia is 84.<br>
  Take the test and see where you stand.
</p>
```
- `font-size:19px; color:#4B5563; font-weight:400; line-height:1.5; margin:20px 0 32px;`

### 2.4 Tombol CTA
```html
<button class="cta-button">
  Start IQ Test Now
  [ICON: arrow-circle.svg]
</button>
```
- `background: linear-gradient(135deg,#5B6EF5,#6C5CE7); color:#fff; font-weight:700; font-size:18px; padding:18px 28px; border-radius:999px; border:none; display:flex; align-items:center; gap:12px; box-shadow:0 12px 24px rgba(91,110,245,0.35); cursor:pointer;`
- Icon arrow-circle: lingkaran putih diameter `28px`, panah biru di dalamnya, diletakkan di ujung kanan tombol.

### 2.5 Baris 4 fitur
```html
<div class="features">
  <div class="feature">
    [ICON: icon-questions.svg] bg:#EDE9FE
    <span>30 Questions</span>
  </div>
  <div class="feature">
    [ICON: icon-report.svg] bg:#DCFCE7
    <span>IQ Score & Detailed Report</span>
  </div>
  <div class="feature">
    [ICON: icon-certificate.svg] bg:#FEF3C7
    <span>Certificate Included</span>
  </div>
  <div class="feature">
    [ICON: icon-infinity.svg] bg:#DBEAFE
    <span>7-Day Free Trial IQ Booster</span>
  </div>
</div>
```
- `.features`: `display:flex; gap:32px; margin:40px 0 24px;` (mobile: `grid; grid-template-columns:1fr 1fr; gap:20px;`)
- `.feature`: `display:flex; flex-direction:column; align-items:center; text-align:center; max-width:120px;`
- Icon container: `width:44px; height:44px; border-radius:12px; display:flex; align-items:center; justify-content:center; margin-bottom:10px;` (background sesuai tabel di atas)
- Teks: `font-size:13px; font-weight:600; color:#374151; line-height:1.3;`

### 2.6 Disclaimer box
```html
<div class="disclaimer">
  [ICON: icon-shield.svg]
  <span>After trial, $29.99 billed every 28 days • Cancel anytime</span>
</div>
```
- `background:#F1F3F9; border-radius:10px; padding:12px 16px; display:flex; align-items:center; gap:10px; font-size:13px; color:#4B5563; width:fit-content;`

---

## 3. HERO-RIGHT — Mockup Device

Susun 3 elemen `<img>`/`<div>` bertumpuk pakai `position:absolute` di dalam container relative (`width:600px; height:600px` sebagai referensi desktop):

```html
<div class="device-stack">
  <div class="decor-blob"></div>      <!-- blob biru transparan, paling belakang -->
  <div class="decor-dots"></div>      <!-- dot grid pojok kanan atas -->
  [ICON: sparkle-outline.svg] x3      <!-- tersebar, posisi absolute -->

  [SVG: device-laptop.svg]   style="position:absolute; top:0; left:0; width:70%; z-index:1;"
  [SVG: device-tablet.svg]   style="position:absolute; bottom:0; right:5%; width:45%; z-index:2; transform:rotate(-3deg);"
  [SVG: device-phone.svg]    style="position:absolute; top:10%; right:0; width:22%; z-index:3; transform:rotate(4deg);"
</div>
```

- `.decor-blob`: lingkaran/ellipse besar, `background:#DCE4FF; opacity:0.4; filter:blur(40px); width:70%; height:70%; border-radius:50%; position:absolute; top:10%; right:0;`
- `.decor-dots`: grid titik-titik kecil (lihat SVG placeholder di Bagian 6), posisi pojok kanan-atas.
- Ketiga device pakai SVG yang SAMA isinya (UI tes IQ) hanya beda ukuran & frame (laptop = frame + keyboard deck, tablet/phone = frame rounded tanpa keyboard).

---

## 4. Isi Layar UI (di dalam tiap SVG device)

Setiap layar berisi:
1. **Header bar** navy (`#12294D`), rounded top corners: teks putih kiri "Question 12 / 30", kanan ikon jam ⏱ + "18:42".
2. **Panel kiri** (background putih): judul kecil "Which shape is missing?" lalu grid 3x3 kotak kecil (pattern kotak oranye/abu/putih acak), kotak terakhir (paling kanan-bawah) background oranye (`#FB923C`) berisi "?" putih. Di bawah grid: tombol outline abu "Back".
3. **Panel kanan** (background putih): judul kecil "Choose your answer" lalu grid 2x2 opsi berlabel A/B/C/D, masing-masing kotak berisi mini pattern grid serupa. Di bawah: tombol solid biru (`#5B6EF5`) rounded "Next Question", teks putih.

Placeholder SVG lengkap layar ini ada di Bagian 6 (`device-screen-content.svg`) — bisa dipakai identik di ketiga device, hanya di-scale.

---

## 5. Palet Warna & Tipografi (ringkasan)

| Token | Hex |
|---|---|
| `--color-text-dark` | `#111827` |
| `--color-accent-blue` | `#4A5CF5` |
| `--color-accent-blue-2` | `#5B6EF5` |
| `--color-accent-purple` | `#6C5CE7` |
| `--color-accent-yellow` | `#F5A623` |
| `--color-accent-orange` | `#FB923C` |
| `--color-navy-header` | `#12294D` |
| `--color-bg` | `#FFFFFF` |
| `--color-text-secondary` | `#4B5563` |
| `--pastel-purple` | `#EDE9FE` |
| `--pastel-green` | `#DCFCE7` |
| `--pastel-orange` | `#FEF3C7` |
| `--pastel-blue` | `#DBEAFE` |

Font: `'Poppins'` atau `'Inter'`, weight 400/500/600/700/800/900 sesuai elemen di atas.

---

## 6. SVG PLACEHOLDER — Pakai Langsung (copy-paste apa adanya)

### 6.1 `avatar-group.svg` (badge icon, 36x36)
```svg
<svg width="36" height="36" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="18" cy="18" r="18" fill="#EDEBFB"/>
  <circle cx="14" cy="15" r="4" stroke="#6C5CE7" stroke-width="1.6"/>
  <circle cx="22" cy="15" r="4" stroke="#6C5CE7" stroke-width="1.6"/>
  <path d="M8 27c0-3.3 2.7-6 6-6h0c3.3 0 6 2.7 6 6" stroke="#6C5CE7" stroke-width="1.6"/>
  <path d="M16 27c0-3.3 2.7-6 6-6h0c3.3 0 6 2.7 6 6" stroke="#6C5CE7" stroke-width="1.6"/>
</svg>
```

### 6.2 `sparkle.svg` (dekorasi headline, ~30x30)
```svg
<svg width="30" height="30" viewBox="0 0 30 30" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M6 3 Q10 8 6 13" stroke="#F5A623" stroke-width="2" stroke-linecap="round" fill="none"/>
  <path d="M14 1 Q18 8 13 16" stroke="#F5A623" stroke-width="2" stroke-linecap="round" fill="none"/>
  <path d="M22 4 Q25 9 21 14" stroke="#F5A623" stroke-width="2" stroke-linecap="round" fill="none"/>
</svg>
```

### 6.3 `arrow-circle.svg` (dalam tombol CTA, 28x28)
```svg
<svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="14" cy="14" r="14" fill="#FFFFFF"/>
  <path d="M10 14h8m0 0-3-3m3 3-3 3" stroke="#5B6EF5" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
```

### 6.4 Ikon fitur (44x44 masing-masing)

**icon-questions.svg**
```svg
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="4" y="3" width="16" height="18" rx="2" stroke="#7C3AED" stroke-width="1.6"/>
  <path d="M8 8h8M8 12h8M8 16h5" stroke="#7C3AED" stroke-width="1.6" stroke-linecap="round"/>
</svg>
```

**icon-report.svg**
```svg
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M4 20V10M10 20V4M16 20v-7M22 20V8" stroke="#16A34A" stroke-width="1.8" stroke-linecap="round"/>
</svg>
```

**icon-certificate.svg**
```svg
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="12" cy="9" r="6" stroke="#F59E0B" stroke-width="1.8"/>
  <path d="M9 14 L7.5 21 L12 18.5 L16.5 21 L15 14" stroke="#F59E0B" stroke-width="1.8" stroke-linejoin="round"/>
</svg>
```

**icon-infinity.svg**
```svg
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M7 15c-2.5 0-4-2-4-3.5S4.5 8 7 8c2.5 0 3.5 2 5 3.5C13.5 10 14.5 8 17 8c2.5 0 4 1.5 4 3.5S19.5 15 17 15c-2.5 0-3.5-2-5-3.5C10.5 13 9.5 15 7 15Z" stroke="#2563EB" stroke-width="1.8"/>
</svg>
```

**icon-shield.svg** (disclaimer, 20x20)
```svg
<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M10 2 L17 5v5c0 4.5-3 7.5-7 8-4-0.5-7-3.5-7-8V5l7-3Z" stroke="#5B6EF5" stroke-width="1.5"/>
  <path d="M7 10l2 2 4-4" stroke="#5B6EF5" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
```

### 6.5 `decor-dots.svg` (dot grid dekorasi, ~120x100)
```svg
<svg width="120" height="100" viewBox="0 0 120 100" xmlns="http://www.w3.org/2000/svg">
  <g fill="#C7D2FE">
    <!-- generate grid 8x6, dot r=2, spacing 15px -->
    <circle cx="5" cy="5" r="2"/><circle cx="20" cy="5" r="2"/><circle cx="35" cy="5" r="2"/><circle cx="50" cy="5" r="2"/><circle cx="65" cy="5" r="2"/><circle cx="80" cy="5" r="2"/>
    <circle cx="5" cy="20" r="2"/><circle cx="20" cy="20" r="2"/><circle cx="35" cy="20" r="2"/><circle cx="50" cy="20" r="2"/><circle cx="65" cy="20" r="2"/><circle cx="80" cy="20" r="2"/>
    <circle cx="5" cy="35" r="2"/><circle cx="20" cy="35" r="2"/><circle cx="35" cy="35" r="2"/><circle cx="50" cy="35" r="2"/><circle cx="65" cy="35" r="2"/><circle cx="80" cy="35" r="2"/>
    <circle cx="5" cy="50" r="2"/><circle cx="20" cy="50" r="2"/><circle cx="35" cy="50" r="2"/><circle cx="50" cy="50" r="2"/><circle cx="65" cy="50" r="2"/><circle cx="80" cy="50" r="2"/>
  </g>
</svg>
```

### 6.6 `sparkle-outline.svg` (dekorasi kecil sekitar device, ~16x16, diamond outline)
```svg
<svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="4" y="4" width="8" height="8" transform="rotate(45 8 8)" stroke="#A5B4FC" stroke-width="1.4"/>
</svg>
```

### 6.7 `device-screen-content.svg` (isi layar — pakai identik di laptop/tablet/phone, ukuran base 420x260, scale sesuai device)
```svg
<svg width="420" height="260" viewBox="0 0 420 260" xmlns="http://www.w3.org/2000/svg">
  <rect width="420" height="260" rx="8" fill="#FFFFFF"/>
  <!-- header bar -->
  <rect width="420" height="34" rx="8" fill="#12294D"/>
  <rect y="26" width="420" height="8" fill="#12294D"/>
  <text x="14" y="22" fill="#FFFFFF" font-size="12" font-family="Inter, sans-serif">Question 12 / 30</text>
  <text x="360" y="22" fill="#FFFFFF" font-size="12" font-family="Inter, sans-serif">⏱ 18:42</text>

  <!-- panel kiri -->
  <text x="20" y="56" fill="#111827" font-size="11" font-weight="600" font-family="Inter, sans-serif">Which shape is missing?</text>
  <g>
    <!-- grid 3x3, tiap kotak 34x34, gap 6, mulai x=20 y=64 -->
    <rect x="20" y="64" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="60" y="64" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="100" y="64" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="20" y="104" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="60" y="104" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="100" y="104" width="34" height="34" rx="4" fill="#FB923C"/>
    <rect x="20" y="144" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="60" y="144" width="34" height="34" rx="4" fill="#F3F4F6" stroke="#E5E7EB"/>
    <rect x="100" y="144" width="34" height="34" rx="4" fill="#FB923C"/>
    <text x="112" y="126" fill="#FFFFFF" font-size="14" font-weight="700" font-family="Inter, sans-serif">?</text>
  </g>
  <rect x="20" y="190" width="50" height="22" rx="11" fill="#FFFFFF" stroke="#D1D5DB"/>
  <text x="30" y="205" fill="#374151" font-size="10" font-family="Inter, sans-serif">Back</text>

  <!-- panel kanan -->
  <text x="220" y="56" fill="#111827" font-size="11" font-weight="600" font-family="Inter, sans-serif">Choose your answer</text>
  <g font-family="Inter, sans-serif" font-size="9">
    <rect x="220" y="64" width="70" height="44" rx="6" fill="#F9FAFB" stroke="#E5E7EB"/>
    <text x="226" y="76" fill="#374151">A)</text>
    <rect x="300" y="64" width="70" height="44" rx="6" fill="#F9FAFB" stroke="#E5E7EB"/>
    <text x="306" y="76" fill="#374151">B)</text>
    <rect x="220" y="116" width="70" height="44" rx="6" fill="#F9FAFB" stroke="#E5E7EB"/>
    <text x="226" y="128" fill="#374151">C)</text>
    <rect x="300" y="116" width="70" height="44" rx="6" fill="#F9FAFB" stroke="#E5E7EB"/>
    <text x="306" y="128" fill="#374151">D)</text>
  </g>
  <rect x="290" y="190" width="100" height="26" rx="13" fill="#5B6EF5"/>
  <text x="305" y="207" fill="#FFFFFF" font-size="10" font-weight="600" font-family="Inter, sans-serif">Next Question</text>
</svg>
```

### 6.8 `device-laptop.svg` (frame laptop membungkus screen di atas, viewBox 460x300)
```svg
<svg width="460" height="300" viewBox="0 0 460 300" xmlns="http://www.w3.org/2000/svg">
  <rect x="20" y="10" width="420" height="260" rx="10" fill="#0B0B0F"/>
  <foreignObject x="30" y="20" width="400" height="240">
    <!-- letakkan device-screen-content.svg (di-scale ke 400x240) di sini -->
  </foreignObject>
  <path d="M0 280h460l-20 20H20l-20-20Z" fill="#6B7280"/>
  <rect x="200" y="284" width="60" height="6" rx="3" fill="#4B5563"/>
</svg>
```

### 6.9 `device-tablet.svg` (frame tablet, viewBox 300x230)
```svg
<svg width="300" height="230" viewBox="0 0 300 230" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="300" height="230" rx="16" fill="#0B0B0F"/>
  <foreignObject x="12" y="14" width="276" height="202">
    <!-- letakkan device-screen-content.svg (di-scale ke 276x202) di sini -->
  </foreignObject>
</svg>
```

### 6.10 `device-phone.svg` (frame phone, viewBox 140x280)
```svg
<svg width="140" height="280" viewBox="0 0 140 280" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="140" height="280" rx="22" fill="#0B0B0F"/>
  <rect x="50" y="8" width="40" height="6" rx="3" fill="#1F2937"/>
  <foreignObject x="8" y="20" width="124" height="240">
    <!-- letakkan device-screen-content.svg (di-scale ke 124x240) di sini -->
  </foreignObject>
</svg>
```

> Catatan implementasi: karena `<foreignObject>` dalam SVG tidak selalu didukung penuh oleh semua renderer, alternatifnya adalah menyalin isi `<g>` dari `device-screen-content.svg` langsung ke dalam masing-masing frame device (laptop/tablet/phone) dengan `transform="translate(x,y) scale(s)"` sesuai ukuran area layar di tiap frame.

---

## 8. Catatan Khusus untuk AI yang Membaca Spesifikasi Ini

- Semua teks dalam Bahasa Inggris HARUS ditulis persis seperti tercantum (jangan diterjemahkan).
- Semua kode warna sudah final, jangan diganti dengan warna lain kecuali diminta.
- Semua SVG di Bagian 6 adalah placeholder siap pakai — cukup copy-paste ke dalam HTML/JSX, tidak perlu generate ulang dari nol atau mencari aset gambar eksternal.
- Ikuti urutan elemen di Bagian 2 dan 3 persis seperti tercantum, jangan menambah atau menghapus elemen tanpa instruksi tambahan.
- Jika target output adalah React/JSX, ubah atribut `class` menjadi `className`, dan SVG inline bisa langsung ditaruh sebagai komponen atau di-import sebagai string.
