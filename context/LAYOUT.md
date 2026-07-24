# LAYOUT.md — Grid, Container, Spacing & Responsive

> Fokus dokumen ini: struktur layout halaman (bukan warna/komponen). Untuk styling visual lihat `DESIGN.md`, untuk aturan komponen lihat `COMPONENTS.md`.

---

## 1. Breakpoints

| Nama       | Lebar Min | Keterangan                     |
| ---------- | --------- | ------------------------------- |
| `mobile`   | 0px       | default, stack vertikal penuh   |
| `tablet`   | 768px     | mulai 2 kolom di beberapa section |
| `desktop`  | 1024px    | layout penuh multi-kolom        |
| `wide`     | 1440px    | container max-width terkunci    |

---

## 2. Container

Container adalah wrapper tunggal yang dipakai berulang (`.container`) di **setiap section pada setiap halaman** — Navbar, Hero, semua Component section, dan Footer. Tidak ada section yang punya max-width sendiri di luar aturan ini, supaya lebar semua bagian benar-benar sejajar/rata (searasi) secara vertikal dari atas sampai bawah halaman.

- **Max-width default:** `1200px`
- **Max-width pada layar besar:** `1290px`, aktif mulai `min-width: 1200px`

```css
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding-inline: 16px; /* mobile default, lihat tabel padding di bawah */
}

@media (min-width: 1200px) {
  .container {
    max-width: 1290px;
  }
}
```

- **Padding horizontal container:**
  - Mobile: `16px`
  - Tablet: `24px`
  - Desktop (`≥1024px`): `40px`
  - Wide (`≥1200px`): `40px` (padding tetap, hanya max-width yang bertambah jadi `1290px`)
- **Centering:** `margin: 0 auto`

### 2.1 Aturan Wajib Konsistensi
- **Navbar:** isi navbar (logo, menu, CTA) dibungkus `.container` yang sama persis — bukan full-bleed lebih lebar dari section lain. Background navbar boleh full-width, tapi kontennya tetap terkunci di `.container`.
- **Hero, semua Component section (Trust Bar, Tabel IQ, Tentang, Cara Kerja, Kenapa Pilih, Pricing, FAQ):** semua memakai `.container` yang identik. Tidak boleh ada section dengan max-width custom (mis. 1100px atau 1320px) — supaya tepi kiri-kanan konten selalu sejajar antar section saat discroll.
- **Footer:** kolom-kolom footer (§6.10) juga dibungkus `.container` yang sama, walau background footer (`--color-bg-dark`) full-bleed selebar layar.
- **Semua halaman lain** (blog, about, pricing detail, dsb.): wajib pakai `.container` yang sama ini, bukan didefinisikan ulang per halaman, supaya lebar konten konsisten di seluruh situs.
- Jika suatu saat butuh section yang benar-benar full-width tanpa batas (mis. banner promo penuh warna), section tersebut tetap membungkus **konten teksnya** dengan `.container` di dalamnya — hanya background yang full-bleed.

---

## 3. Grid System

- Gunakan CSS Grid / Flexbox dengan basis **12-column grid** untuk section yang butuh kolom (tabel negara, pricing, footer).
- **Gutter (jarak antar kolom):** `24px` desktop, `16px` mobile.
- Section bebas grid (hero, section teks panjang) boleh pakai flexbox 2 kolom sederhana (60/40 atau 50/50).

---

## 4. Spacing Scale

Skala konsisten (kelipatan 4px) untuk margin/padding antar elemen & antar section:

```
4px   - xs   (jarak antar elemen kecil, mis. icon-teks)
8px   - sm
16px  - md   (padding dalam card/button)
24px  - lg
32px  - xl
48px  - 2xl  (jarak antar blok dalam satu section)
80px  - 3xl  (jarak antar section besar, desktop)
120px - 4xl  (jarak antar section besar, wide screen)
```

- **Jarak antar section (vertical rhythm):**
  - Mobile: `48px`
  - Tablet: `64px`
  - Desktop: `96px`

---

## 5. Struktur Halaman (Urutan Section)

```
├── Top Bar (opsional, statis 1 baris, tinggi ~32px)
├── Header / Navbar (sticky, tinggi ~72px)
├── Hero Section
├── Trust Bar (statistik singkat)
├── Tabel Perbandingan IQ per Negara
├── Section "Tentang Tes"
├── Section "Cara Kerja" (3 langkah)
├── Section "Kenapa Pilih Tes Ini"
├── Section Pricing
├── FAQ
├── Footer
└── Back-to-top button (floating, posisi fixed)
```

Isi konten tiap section (copywriting, gambar) tidak dibahas di file ini — dokumen ini hanya mengatur **struktur grid & spacing-nya**. Aturan visual tiap komponen ada di `COMPONENTS.md`.

---

## 6. Layout per Section

### 6.1 Header
- Desktop: 1 baris flex, `justify-content: space-between` (logo kiri — menu tengah — CTA + bahasa kanan)
- Mobile: logo + hamburger icon saja, menu masuk ke slide-in panel dari kanan (lebar panel: `80%` dari viewport, max `320px`)

### 6.2 Hero
- Desktop: grid 2 kolom, rasio `55% teks / 45% gambar`
- Tablet: rasio `50/50`
- Mobile: stack vertikal, urutan → label kecil → judul → sub-judul → CTA → bullet list → gambar di paling bawah

### 6.3 Trust Bar
- Desktop: flex row 4 kolom sejajar, rata tengah
- Mobile: grid 2x2

### 6.4 Tabel IQ Negara
- Desktop: 3 kolom tabel sejajar (grid 3 kolom, gutter `24px`)
- Tablet: 2 kolom + 1 kolom di bawah
- Mobile: 1 kolom, tabel stacked, atau scrollable horizontal per tabel dengan `overflow-x: auto`

### 6.5 Section "Tentang Tes"
- Desktop: 2 kolom (gambar 45% - teks 55%), gambar di kiri (posisi dibalik dari hero agar ritme visual bervariasi)
- Mobile: stack, gambar di bawah teks

### 6.6 Section "Cara Kerja"
- Desktop: 3 kolom sejajar, lebar sama rata (`grid-template-columns: repeat(3, 1fr)`)
- Tablet: 2 kolom + 1 di bawah
- Mobile: stack vertikal, garis vertikal tipis penghubung antar step (opsional)

### 6.7 Section "Kenapa Pilih Tes Ini"
- Desktop: gambar (bell curve) di atas, teks di bawah, max-width teks `700px` center-aligned

### 6.8 Pricing
- Single card, center-aligned, `max-width: 480px`
- Padding dalam card: `32px` desktop / `24px` mobile

### 6.9 FAQ
- Accordion full-width, `max-width: 800px`, center-aligned

### 6.10 Footer
- Desktop: grid 5 kolom (`repeat(5, 1fr)`), gutter `32px`
- Tablet: grid 2 kolom, wrap
- Mobile: stack 1 kolom, tiap kelompok link collapsible (opsional)

---

## 7. Aturan Responsive Umum

- Semua CTA button: `width: auto` di desktop, `width: 100%` di mobile (`<768px`)
- Gambar/ilustrasi: `max-width: 100%; height: auto` — tidak ada gambar dengan lebar fixed px yang bisa overflow di mobile
- Section dengan background image/illustration besar: sediakan fallback warna solid di mobile agar loading tidak berat
- Sticky header: aktif di semua breakpoint, tinggi mengecil sedikit saat scroll (opsional: `72px → 56px`)
- Back-to-top button: muncul setelah scroll melewati tinggi 1 viewport, posisi `fixed`, `bottom: 24px; right: 24px` (mobile: `bottom: 16px; right: 16px`)