# COMPONENTS.md — Aturan Komponen UI

> Fokus dokumen ini: perilaku & anatomi tiap komponen. Token warna/tipografi/radius/shadow merujuk ke `DESIGN.md`. Posisi & grid merujuk ke `LAYOUT.md`.

---

## 1. Button

### 1.1 Varian
| Varian       | Background            | Teks                 | Border                  | Penggunaan            |
| ------------- | ----------------------- | ----------------------- | -------------------------- | ------------------------ |
| Primary      | `--color-primary`      | putih                  | none                       | CTA utama ("Mulai Tes IQ") |
| Secondary    | transparan              | `--color-primary`      | `1px solid --color-primary`| CTA sekunder ("Masuk Member") |
| Ghost        | transparan              | `--color-text`         | none                       | link aksi ringan, mis. "Lihat Semua" |
| Danger       | `--color-error`        | putih                  | none                       | konfirmasi batal langganan |

### 1.2 Ukuran
| Size | Padding (v/h)   | Font size       | Radius            |
| ----- | ----------------- | ------------------ | -------------------- |
| `sm` | 8px / 16px        | `--text-sm`         | `--radius-md`         |
| `md` | 12px / 24px       | `--text-body`       | `--radius-md`         |
| `lg` | 16px / 32px       | `--text-lead`       | `--radius-lg`         |

### 1.3 State
- **Default → Hover:** background jadi `--color-primary-dark`, tambah `transform: translateY(-2px)`, shadow naik ke `--shadow-md`.
- **Active/Pressed:** `transform: translateY(0)`, shadow turun ke `--shadow-sm`.
- **Disabled:** opacity `0.5`, `cursor: not-allowed`, tanpa hover effect.
- **Loading:** teks diganti spinner kecil (16px), button tetap ukuran sama (tidak "jump").

### 1.4 Aturan
- Setiap button wajib punya state `:focus-visible` dengan `--shadow-glow` untuk aksesibilitas keyboard.
- Icon dalam button (jika ada) selalu di kiri teks, gap `8px`.
- Mobile: button CTA utama full-width (`width: 100%`), CTA sekunder tetap auto-width.

---

## 2. Card

### 2.1 Card Umum (mis. card statistik, card "Cara Kerja")
- Background: `--color-bg`
- Border: `1px solid --color-border`
- Radius: `--radius-md`
- Padding: `24px`
- Shadow default: `--shadow-sm`, hover: `--shadow-md` (jika card interaktif/clickable)

### 2.2 Pricing Card (khusus, lebih besar)
- Radius: `--radius-lg`
- Shadow: `--shadow-lg`
- Border: gradient tipis 1px (`linear-gradient(135deg, var(--color-primary), var(--color-secondary))`) — pembeda visual dari card biasa
- Badge "Paling Populer" menempel di pojok kanan atas card, sedikit overflow ke luar card (`position: absolute; top: -12px; right: 24px`)
- Struktur internal (urutan dari atas ke bawah):
  1. Nama paket
  2. Harga besar + keterangan harga lanjutan (font kecil di bawah harga)
  3. Divider tipis
  4. List checklist fitur (icon check `--color-success` + teks)
  5. Button CTA full-width di dalam card
  6. Teks kecil disclaimer/link detail harga

---

## 3. Navbar / Header

### 3.1 Anatomi (desktop)
```
[Logo]   [Menu: Beranda | Contoh Laporan | Harga | Bantuan]   [Globe-Bahasa] [Button Secondary: Masuk] [Button Primary: Mulai Tes]
```

### 3.2 Perilaku
- **Sticky:** menempel di atas saat scroll, background solid + `--shadow-sm` muncul setelah scroll > 10px (saat di top, background transparan/blend dengan hero jika hero punya background berwarna).
- **Menu item aktif:** underline tipis `2px` warna `--color-primary` di bawah teks.
- **Dropdown bahasa:** trigger dengan klik (bukan hover saja, untuk ramah mobile/touch), list muncul sebagai dropdown card kecil dengan shadow `--shadow-md`, radius `--radius-md`.

### 3.3 Mobile
- Header hanya menampilkan: Logo + Icon Hamburger.
- Klik hamburger → panel slide-in dari **kanan** (bukan dari kiri seperti kebanyakan template default), lebar `80%`/max `320px`, overlay gelap semi-transparan di belakang panel (`rgba(0,0,0,0.4)`).
- Di dalam panel: menu vertikal → divider → pilihan bahasa (accordion) → 2 button CTA di bagian bawah panel (sticky di dasar panel).
- Icon close (X) di pojok kanan atas panel.

---

## 4. Footer

### 4.1 Struktur Kolom
1. **Brand:** logo versi putih/terang + tagline + jam layanan support
2. **Dukungan:** Hubungi Support, Batalkan Langganan
3. **Legal:** Privasi, Syarat & Ketentuan, Kebijakan Langganan, Harga
4. **Eksplorasi:** Tentang Kami, FAQ, Blog
5. **Program Lanjutan:** Login, Tentang Program

### 4.2 Style
- Background: `--color-bg-dark`
- Teks: `--color-text-muted` versi terang (mis. `#94A3B8`), heading kolom pakai putih `#FFFFFF` `--text-sm` bold uppercase letter-spacing kecil
- Link footer: default `#CBD5E1`, hover → `--color-primary` + underline
- Baris disclaimer & copyright: dipisahkan `border-top: 1px solid rgba(255,255,255,0.1)`, padding atas `24px`, font `--text-xs`, rata tengah

---

## 5. Input & Form (untuk form kontak/pencarian jika ada)

### 5.1 Text Input
- Height: `44px`
- Padding horizontal: `16px`
- Border: `1px solid --color-border`, radius `--radius-sm`
- Focus: border `--color-primary`, tambah `--shadow-glow`
- Error: border `--color-error`, teks bantuan error di bawah input warna `--color-error`, `--text-xs`

### 5.2 Select / Dropdown (mis. filter "Urutkan berdasarkan Benua")
- Sama seperti text input, tambah icon chevron-down di kanan (`16px`, `--color-text-muted`)
- Dropdown list: card mengambang, `--shadow-md`, radius `--radius-sm`, max-height `240px` dengan scroll jika opsi banyak

---

## 6. Badge

| Jenis                 | Style                                                        |
| ---------------------- | --------------------------------------------------------------- |
| Badge status ("Populer")| background `--color-secondary`, teks putih, radius `--radius-full`, padding `4px 12px`, `--text-xs` bold |
| Badge sosial proof kecil (label hero) | background transparan, teks `--color-text-muted`, uppercase, letter-spacing lebar, tanpa background |
| Badge checklist item   | icon check bulat kecil (`16px`) warna `--color-success`, tanpa background |

---

## 7. Modal

### 7.1 Penggunaan
Dipakai untuk: konfirmasi pembatalan langganan, preview sertifikat, atau popup bahasa (opsional).

### 7.2 Anatomi
- Overlay: `rgba(17,24,39,0.5)`, klik di luar modal → menutup modal
- Container modal: max-width `480px`, radius `--radius-lg`, shadow `--shadow-lg`, padding `32px`
- Header modal: judul (`--text-h3`) + icon close (kanan atas)
- Footer modal: 2 button (Ghost untuk "Batal", Primary/Danger untuk aksi konfirmasi), rata kanan

### 7.3 Animasi
- Muncul: fade overlay + scale modal dari `0.95 → 1` dalam `--duration-normal`
- Tutup: kebalikannya, `--duration-fast`

---

## 8. FAQ (Accordion)

### 8.1 Anatomi per Item
```
[Icon chevron kiri/kanan] [Pertanyaan — --text-lead bold]
   └── (saat expand) [Jawaban — --text-body, --color-text-muted]
```

### 8.2 Perilaku
- Hanya **satu item terbuka dalam satu waktu** (accordion eksklusif) — beda dari kebanyakan FAQ yang membolehkan banyak terbuka sekaligus, dipakai di sini agar halaman tetap ringkas.
- Klik pertanyaan yang sedang terbuka → menutup kembali.
- Transisi height memakai `--duration-normal`, `--easing-default`.
- Divider tipis `1px solid --color-border` antar item, tanpa card/box per item (list style, bukan card style).

---

## 9. Pricing Section (gabungan komponen)

- Menggunakan **Pricing Card** (lihat §2.2) sebagai elemen utama, diletakkan center dalam container `max-width: 480px` (lihat `LAYOUT.md` §6.8).
- Di bawah card: 1 baris micro-testimonial (avatar inisial bulat + kutipan singkat + nama fiktif/generik), font `--text-sm italic`, warna `--color-text-muted`.
- Link kecil "Lihat detail harga →" di bawah testimonial, style Ghost button.

---

## 10. Trust Bar Item

- Anatomi: `[Icon 24px] [Angka besar --text-h3 bold] [Label kecil --text-xs --color-text-muted di bawah angka]`
- Layout per item: vertikal, rata tengah (`text-align: center`, `flex-direction: column`)
- Tidak ada border/card pemisah antar item — cukup gap `48px` (desktop) / grid 2x2 gap `24px` (mobile), sesuai `LAYOUT.md` §6.3

---

## 11. Tabel Perbandingan (Tabel IQ Negara)

- Header tabel: background `--color-bg-alt`, teks bold `--text-sm`, padding `12px 16px`
- Baris data: padding `12px 16px`, border-bottom `1px solid --color-border`
- Baris berselang (zebra stripe): baris genap background `--color-bg-alt` tipis (opsional, untuk keterbacaan)
- Kolom "Peringkat": badge bulat kecil angka (bukan teks polos) — pembeda dari tabel referensi yang teksnya polos
- Hover baris (desktop): background `--color-bg-alt`, cursor default (tabel tidak clickable)

---

## 12. Ikon Numerik "Cara Kerja" (01/02/03)

- Angka besar `--text-h1` size tapi weight `800`, warna `--color-secondary` dengan opacity `0.15` sebagai watermark besar di belakang icon/judul step
- Icon step (di depan angka watermark) ukuran `40px`, warna `--color-primary`
- Susunan per step: angka watermark (background) → icon (depan, overlap sedikit) → judul step (`--text-h3`) → deskripsi singkat (`--text-body`, `--color-text-muted`)
