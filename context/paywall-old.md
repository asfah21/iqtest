# Paywall Page — IQTest

> Dokumentasi untuk halaman [`/paywall/:id`](handlers/quiz.go) — titik konversi utama dari funel IQ Test.

---

## 1. Tujuan Halaman

Halaman paywall adalah **gerbang pembayaran** satu kali (Rp14.900) untuk membuka hasil tes kognitif. User tiba di sini setelah menyelesaikan 20 soal (`POST /submit-tes`) dan diarahkan ke `/paywall/{id}`. Jika user mencoba akses `/hasil/{id}` langsung tanpa bayar, akan di-redirect kembali ke paywall.

**Flow:**

```
Submit Tes → POST /submit-tes → redirect ke /paywall/{id}
     ↓
User bayar (manual) → POST /konfirmasi-bayar/{id} → redirect ke /hasil/{id}
     ↓
Akses hasil lengkap
```

---

## 2. Struktur File

| File | Peran |
|------|-------|
| [`templ/pages/paywall_page.templ`](templ/pages/paywall_page.templ) | Templ component utama halaman paywall |
| [`templ/layouts/paywall_layout.templ`](templ/layouts/paywall_layout.templ) | Layout wrapper (head, public-content, footer, back-to-top) |
| [`models/user.go`](models/user.go:44) | Definisi `PaywallData` struct |
| [`handlers/quiz.go`](handlers/quiz.go) | Handler `ShowPaywall`, `KonfirmasiBayar` |
| [`assets/css/style.css`](assets/css/style.css) | Semua kelas CSS pendukung (card, button, utility) |

---

## 3. Data Model

Definisi di [`models/user.go`](models/user.go:44):

```go
type PaywallData struct {
    ID              string   // UUID sesi user
    Nama            string   // Nama user
    DurationMinutes int      // Lama pengerjaan tes (menit)
    DomainMessage   string   // Pesan tentang domain terkuat (opsional)
}
```

---

## 4. Layout & Navigasi

Halaman menggunakan [`PaywallLayout`](templ/layouts/paywall_layout.templ) yang membungkus konten dengan:

- **`Head` component** — title tag "Bayar"
- **`public-content`** — wrapper dengan padding-top untuk floating navbar
- **`Footer` component** — footer standar
- **Back-to-top button** — floating fixed di pojok kanan bawah (Alpine.js)

Konten paywall dibatasi lebar `max-width: 640px` dan rata tengah.

---

## 5. Komponen UI (Urutan dari Atas ke Bawah)

### 5.1 Banner Hasil Tes

```
┌────────────────────────────────────────────────┐
│ ✅ Kamu menyelesaikan tes dalam X menit        │
│    {DomainMessage (jika ada)}                  │
└────────────────────────────────────────────────┘
```

- Kelas: `.card`, `.p-5`, `.mb-6`, `.fade-in-up`, `.border-l-4`, `.border-l-accent`, `bg-[var(--bgAlt)]`
- Animasi `fadeInUp` dengan `animation-delay: 0s`
- Ikon check-circle SVG di kiri, teks durasi + domain message di kanan
- Warna aksen: `var(--accent)` (#6366f1, indigo)

### 5.2 Header & Value Proposition

```
┌────────────────────────────────────────────────┐
│ [Akses Premium] badge                           │
│ Selangkah Lagi, {Nama}                         │
│ Dapatkan skor lengkap, analisis per domain...   │
└────────────────────────────────────────────────┘
```

- Badge "Akses Premium" dengan titik kuning (`var(--warning)`)
- Heading `h1` dengan `text-3xl md:text-4xl font-bold tracking-tight`
- Nama user di-highlight dengan `text-accent` (indigo)
- Subtitle `text-[var(--textMuted)]` dengan `max-w-md` center

### 5.3 Preview Hasil (Blurred)

```
┌────────────────────────────────────────────────┐
│ Hasil Premium                    [Terkunci]     │
│                                                 │
│ Skor Total          •• / 30.5                  │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░                             │
│                                                 │
│ Persentil           ••%                         │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░                             │
└────────────────────────────────────────────────┘
```

- Kelas: `.card`, `.card-elevated`, `.p-6`
- Badge "Terkunci" dengan kelas `.badge`, `.badge-warning`
- Konten di-blur dengan `opacity-50 blur-sm select-none`
- Progress bar menggunakan `.progress-bar` + `.progress-fill`
- Skor total ditampilkan sebagai `•• / 30.5` (disensor)
- Persentil ditampilkan sebagai `••%` (disensor)
- Warna progress bar: `#e7e5e4` (abu-abu netral, bukan accent)

**Catatan desain:** Saat ini preview hanya menampilkan skor total + persentil yang di-blur. Belum ada visualisasi per domain (MTX/SEQ/SPA/ANL) — ini area potensial untuk *curiosity gap* yang lebih kuat.

### 5.4 Card Pembayaran

```
┌────────────────────────────────────────────────┐
│            Harga Spesial                        │
│            Rp14.900                             │
│            / sekali                             │
│                                                 │
│  ┌──────────────────────────────────────┐       │
│  │          [QR Code Placeholder]       │       │
│  │          Scan QRIS                   │       │
│  └──────────────────────────────────────┘       │
│  Scan QR code di atas menggunakan aplikasi      │
│  mobile banking atau e-wallet.                  │
│                                                 │
│  [Nama pengirim transfer]                       │
│                                                 │
│  ┌──────────────────────────────────────┐       │
│  │  ✅ Saya Sudah Bayar                 │       │
│  └──────────────────────────────────────┘       │
│                                                 │
│  Hasil akan langsung bisa diakses setelah       │
│  pembayaran terverifikasi.                      │
└────────────────────────────────────────────────┘
```

- Kelas card: `.card`, `.card-elevated`, `.p-8`, `.text-center`
- Label "Harga Spesial" dengan `text-[var(--textSubtle)]`
- Harga: `text-4xl font-bold text-[var(--textMain)]` — **Rp14.900**
- QRIS container terpisah dengan `bg-[var(--bgAlt)]`, `rounded-2xl`, `border`
- Placeholder QR code berupa ikon SVG (kotak) + teks "Scan QRIS"
- Input field: `.input-field` untuk `namaPengirim`
- Tombol CTA: `.btn-primary`, `.w-full`, `.shadow-glow` — dengan SVG icon + teks "Saya Sudah Bayar"

### 5.5 Link Kembali

- Tombol `.btn-ghost` dengan ikon panah kiri + teks "Kembali ke Beranda"
- Navigasi ke `/`

---

## 6. Animasi & Transisi

| Elemen | Animasi | Delay | Detail |
|--------|---------|-------|--------|
| Banner hasil | `fade-in-up` | 0s | `animation: fadeInUp 0.5s ease both` |
| Card preview skor | `fade-in-up` | 0.1s | — |
| Card pembayaran | `fade-in-up` | 0.2s | — |
| Semua card | `transition-all duration-300 hover:shadow-xl` | — | Efek shadow naik saat hover |
| Tombol "Saya Sudah Bayar" | `.shadow-glow`, `transition-all duration-300` | — | Glow subtle |
| QRIS container | `transition-all duration-300 hover:shadow-sm` | — | Shadow naik saat hover |
| Badge "Akses Premium" | `transition-all duration-300 hover:shadow-sm` | — | — |
| Input field | `transition-all duration-200` | — | Fokus border-color berubah |

---

## 7. JavaScript (Client-side)

Fungsi [`konfirmasiBayar`](templ/pages/paywall_page.templ:102) — inline script di halaman:

```javascript
function konfirmasiBayar(btn) {
    const nama = document.getElementById('namaPengirim').value.trim();
    if (!nama) { alert('Silakan masukkan nama pengirim transfer.'); return; }
    const id = btn.getAttribute('data-paywall-id');
    btn.disabled = true;
    btn.innerHTML = '... Memverifikasi...';
    fetch('/konfirmasi-bayar/' + id, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nama_pengirim: nama })
    })
    .then(r => r.json())
    .then(data => {
        if (data.success) { window.location.href = '/hasil/' + data.id; }
        else { /* tampilkan error + restore tombol */ }
    })
    .catch(() => { /* tampilkan error + restore tombol */ });
}
```

**Alur:**
1. Validasi `namaPengirim` tidak boleh kosong
2. Disable tombol + tampilkan spinner "Memverifikasi..."
3. `POST /konfirmasi-bayar/{id}` dengan body `{ nama_pengirim }`
4. Sukses → redirect ke `/hasil/{id}`
5. Gagal → tampilkan error + restore tombol ke state semula

---

## 8. CSS Design Tokens yang Digunakan

Token dari [`assets/css/style.css`](assets/css/style.css) yang relevan dengan halaman paywall:

| Token | Value | Penggunaan |
|-------|-------|------------|
| `--accent` | `#6366f1` | Warna aksen indigo (teks, border, icon) |
| `--accentHover` | `#4f46e5` | Hover state tombol primary |
| `--accentLight` | `#eef2ff` | Background ikon check di banner |
| `--ink` | `#0f172a` | Teks utama (judul, harga) |
| `--inkMuted` | `#475569` | Teks sekunder |
| `--inkSubtle` | `#94a3b8` | Teks tersier (label, sensor) |
| `--bg` | `#f8fafc` | Background halaman (`bg-gradient-paywall`) |
| `--surface` | `#ffffff` | Background card |
| `--bgAlt` | `#e8e4de` | Background QRIS container |
| `--bordLight` | `#f1f5f9` | Border card subtle |
| `--bord` | `#e2e8f0` | Border card default |
| `--warning` | `#d97706` | Warna badge "Akses Premium" dot |
| `--warm-50` | `#fffbeb` | Background badge akses premium |
| `--warm-200` | `#fde68a` | Border badge akses premium |
| `--warm-600` | `#d97706` | Teks badge akses premium |
| `--shadow1` | `0 1px 2px rgba(15,23,42,0.04), 0 4px 12px rgba(15,23,42,0.06)` | Shadow card default |
| `--shadow2` | `0 4px 16px rgba(15,23,42,0.08), 0 12px 40px rgba(15,23,42,0.06)` | Shadow card hover |

---

## 9. Catatan Implementasi & Area Pengembangan

### 9.1 Sudah Diimplementasi

- ✅ Personal greeting ("Selangkah Lagi, {nama}")
- ✅ Banner durasi + domain message
- ✅ Preview skor di-blur (opacity + blur CSS)
- ✅ Harga Rp14.900 dengan label "Harga Spesial"
- ✅ Container QRIS dengan border + instruksi
- ✅ Input nama pengirim transfer
- ✅ Tombol CTA "Saya Sudah Bayar" dengan loading state
- ✅ Animasi fade-in-up dengan staggered delay
- ✅ Card elevation dengan shadow halus
- ✅ Responsive (max-width 640px, utility classes)
- ✅ Link kembali ke beranda

### 9.2 Belum Diimplementasi (Potensi Improvisasi)

- ❌ Trust signals (testimoni, jumlah user, badge keamanan)
- ❌ Social proof micro-copy (mis. "12.847 orang sudah mengetahui skor IQ mereka")
- ❌ Curiosity gap visual per domain (progress bar per domain yang di-blur)
- ❌ Urgency/scarcity (batas waktu harga spesial)
- ❌ Micro-interactions pada tombol (scale hover, fade-in yang lebih halus)
- ❌ Aksen warna premium konsisten (saat ini masih indigo netral)

### 9.3 Konvensi Proyek

- **Brand:** IQTest — produk asesmen kemampuan kognitif
- **Stack:** Go + Gin + templ + Alpine.js + PostgreSQL
- **Design tokens:** Prefix `--ink`, `--accent`, `--bg`, `--surface`, dll (lihat DESIGN.md)
- **Layout:** Container `max-width: 1200px` (1290px di layar ≥1200px), padding 16–40px (lihat LAYOUT.md)
- **Paywall content wrapper:** `max-width: 640px; margin: 0 auto` inline
- **Semua perubahan UI:** hanya di file `.templ` + CSS, regenerate dengan `templ generate`
- **Jangan edit** `*_templ.go` — file generated

---

## 10. Referensi

- [`IQTEST.md`](context/IQTEST.md) — Spesifikasi lengkap engine tes (scoring, API flow, UI/UX flow §12.2.3)
- [`DESIGN.md`](context/DESIGN.md) — Panduan desain visual (warna, tipografi, shadow)
- [`LAYOUT.md`](context/LAYOUT.md) — Panduan grid, container, spacing & responsive
- [`COMPONENTS.md`](context/COMPONENTS.md) — Aturan komponen reusable
