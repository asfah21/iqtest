# DESIGN.md — Warna, Tipografi, Radius, Shadow, Icon & Animation

> Fokus dokumen ini: bahasa visual (design tokens). Untuk struktur grid lihat `LAYOUT.md`, untuk aturan komponen lihat `COMPONENTS.md`.

---

## 1. Palet Warna

> Sengaja berbeda dari palet biru khas situs referensi (wwiqtest.com) — gunakan kombinasi **teal-ungu** sebagai identitas visual pembeda.

### 1.1 Warna Utama (Brand)
| Token            | Hex       | Penggunaan                          |
| ----------------- | --------- | ------------------------------------ |
| `--color-primary`      | `#0EA5A0` | CTA utama, link aktif, highlight     |
| `--color-primary-dark` | `#0B7A76` | hover state tombol primary           |
| `--color-secondary`    | `#7C3AED` | aksen sekunder (badge, ikon, angka)  |
| `--color-secondary-dark`| `#5F27CD` | hover state tombol secondary        |

### 1.2 Warna Netral
| Token           | Hex       | Penggunaan                     |
| ---------------- | --------- | -------------------------------- |
| `--color-bg`          | `#FFFFFF` | background utama                |
| `--color-bg-alt`      | `#F7F8FA` | background section selang-seling |
| `--color-bg-dark`      | `#111827` | footer background               |
| `--color-text`        | `#1F2937` | teks body utama                 |
| `--color-text-muted`  | `#6B7280` | teks sekunder/caption           |
| `--color-border`      | `#E5E7EB` | border card, divider            |

### 1.3 Warna Status
| Token             | Hex       | Penggunaan            |
| ------------------ | --------- | ----------------------- |
| `--color-success`      | `#16A34A` | badge sukses, checklist |
| `--color-warning`      | `#F59E0B` | badge peringatan        |
| `--color-error`        | `#DC2626` | pesan error form        |

### 1.4 Aturan Penggunaan
- Warna `primary` hanya untuk elemen aksi (CTA, link penting) — jangan dipakai sebagai warna teks panjang.
- Warna `secondary` (ungu) dipakai untuk elemen aksen: nomor step "01/02/03", badge "Paling Populer", ikon di Trust Bar.
- Section background berselang-seling antara `--color-bg` dan `--color-bg-alt` untuk membedakan tiap section secara visual tanpa garis pembatas tegas.

---

## 2. Tipografi

### 2.1 Font Family
- **Heading:** `'Sora', sans-serif` (geometris, tegas — beda dari font default kebanyakan landing page tes IQ yang cenderung pakai font system standar)
- **Body:** `'Inter', sans-serif`

### 2.2 Skala Ukuran (Type Scale)

| Token        | Size (desktop) | Size (mobile) | Weight | Line-height |
| ------------- | --------------- | -------------- | ------ | ------------ |
| `--text-h1`  | 48px            | 32px           | 700    | 1.15         |
| `--text-h2`  | 36px            | 26px           | 700    | 1.2          |
| `--text-h3`  | 24px            | 20px           | 600    | 1.3          |
| `--text-lead`| 20px            | 17px           | 400    | 1.5          |
| `--text-body`| 16px            | 15px           | 400    | 1.6          |
| `--text-sm`  | 14px            | 13px           | 400    | 1.5          |
| `--text-xs`  | 12px            | 12px           | 500    | 1.4          |

### 2.3 Aturan
- H1 hanya boleh muncul satu kali per halaman (di Hero).
- Label kecil di atas H1 (mis. "Sudah dipakai lebih dari X juta orang") memakai `--text-xs`, uppercase, `letter-spacing: 0.08em`.
- Maksimal lebar baris teks paragraf: `70ch` agar mudah dibaca.

---

## 3. Border Radius

| Token          | Value  | Penggunaan                    |
| --------------- | ------ | -------------------------------- |
| `--radius-sm`  | 6px    | input, badge kecil               |
| `--radius-md`  | 12px   | button, card kecil                |
| `--radius-lg`  | 20px   | card besar (pricing card, modal)  |
| `--radius-full`| 999px  | pill button, badge bulat penuh    |

> Catatan: gunakan radius yang lebih besar (`lg`) dibanding referensi asli yang cenderung lebih kotak — memberi kesan lebih ramah & modern.

---

## 4. Shadow / Elevation

| Token           | Value                                      | Penggunaan               |
| ---------------- | -------------------------------------------- | --------------------------- |
| `--shadow-sm`   | `0 1px 2px rgba(17,24,39,0.06)`              | card default, input focus   |
| `--shadow-md`   | `0 4px 12px rgba(17,24,39,0.08)`             | card hover, dropdown         |
| `--shadow-lg`   | `0 12px 32px rgba(17,24,39,0.12)`            | pricing card, modal          |
| `--shadow-glow` | `0 0 0 4px rgba(14,165,160,0.15)`            | focus ring elemen interaktif |

---

## 5. Icon

- **Library:** gunakan icon set line-style konsisten (mis. Phosphor Icons atau Lucide) — hindari icon set default yang sama persis dengan situs referensi.
- **Ukuran standar:**
  - Icon inline (dekat teks): `16px`
  - Icon di card/Trust Bar: `24px`
  - Icon besar (dekorasi section): `40px`–`56px`
- **Warna icon:** default `--color-text-muted`, berubah `--color-primary` saat aktif/hover.
- **Style:** stroke-based (outline), stroke-width `1.5px`, bukan filled icon — agar terasa lebih ringan.

---

## 6. Ilustrasi & Gambar

- Gunakan ilustrasi flat/isometric custom (bukan foto stok generik atau ilustrasi identik dari situs referensi).
- Warna ilustrasi mengikuti palet brand (teal + ungu + netral), bukan biru dominan seperti referensi.
- Semua gambar dekoratif memakai `border-radius: var(--radius-lg)` bila berbentuk kotak.

---

## 7. Animation & Transition

### 7.1 Durasi & Easing Standar
| Token                | Value                              |
| --------------------- | ------------------------------------- |
| `--duration-fast`    | 150ms                                 |
| `--duration-normal`  | 250ms                                 |
| `--duration-slow`    | 400ms                                 |
| `--easing-default`   | `cubic-bezier(0.4, 0, 0.2, 1)`         |

### 7.2 Aturan Animasi
- **Hover button/card:** transisi `background-color`, `box-shadow`, `transform` (translateY -2px) dengan `--duration-fast`.
- **Accordion FAQ:** expand/collapse height dengan `--duration-normal`, easing default.
- **Scroll-reveal section:** fade-in + translateY(16px → 0) saat elemen masuk viewport, `--duration-slow`, trigger sekali saja (tidak berulang saat scroll naik-turun).
- **Sticky header saat scroll:** transisi tinggi & shadow halus `--duration-fast`.
- **Back-to-top button:** fade + scale saat muncul/hilang, `--duration-fast`.
- Hindari animasi berlebihan (parallax berat, autoplay video besar) yang memperlambat loading di mobile.

---

## 8. Dark Mode (Opsional)

Jika diperlukan versi dark mode di masa depan:
- `--color-bg` → `#0F172A`
- `--color-text` → `#F1F5F9`
- `--color-bg-alt` → `#1E293B`
- Warna primary/secondary tetap sama, hanya kontras background yang disesuaikan.
