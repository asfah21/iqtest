# Prompt: Refactor UI Halaman Paywall — ShadowSelf

## Konteks project

Ini halaman `/paywall/:id` di project ShadowSelf (Go + Gin + templ + Alpine.js).
Dokumentasi lengkap ada di:
- `templ/pages/paywall_page.templ` — komponen utama
- `templ/layouts/paywall_layout.templ` — layout wrapper
- `models/user.go:44` — struct `PaywallData`
- `assets/css/style.css` — semua design token & utility class
- `context/DESIGN.md`, `context/LAYOUT.md`, `context/COMPONENTS.md` — panduan desain project

**Tujuan:** halaman ini adalah titik konversi utama funnel (bayar Rp14.900 untuk buka
hasil tes). Struktur & fungsinya sudah lengkap dan benar (lihat §9.1 di dokumentasi
project — greeting personal, blur preview, QRIS, tombol konfirmasi, animasi
fade-in-up sudah ada). **Masalahnya bukan kurang fitur, tapi tampilannya masih
terasa seperti template generik / "AI slop"** — bukan produk psikometri premium
yang orang percaya untuk keluarkan uang.

## Batasan wajib (jangan dilanggar)

- **Jangan ubah struktur data / handler / route.** `PaywallData`, `ShowPaywall`,
  `KonfirmasiBayar` di `handlers/quiz.go` tetap sama.
- **Jangan ubah logic JS** di fungsi `konfirmasiBayar` (validasi nama pengirim,
  fetch ke `/konfirmasi-bayar/{id}`, redirect ke `/hasil/{id}`) — hanya boleh
  ubah markup/kelas di sekitarnya.
- **Pakai token CSS yang sudah ada** (`--color-primary`, `--color-text`, `--color-text-muted`, `--color-bg`,
  `--shadow1`, `--shadow2`, dll — lihat tabel §8 dokumentasi). Jangan
  tambah warna baru sembarangan; kalau butuh warna tambahan, definisikan sebagai
  token baru di `style.css` mengikuti pola penamaan yang ada, jangan hardcode hex
  di `.templ`.
- **Container tetap `max-width: 640px; margin: 0 auto`** sesuai konvensi project.
  Jangan ubah ke lebar lain tanpa alasan kuat.
- Semua perubahan hanya di file `.templ` + `style.css`. **Jangan sentuh file
  `*_templ.go`** (generated) — jalankan `templ generate` setelah edit.
- Stack tetap Alpine.js untuk interaktivitas ringan — jangan tambah framework
  frontend baru.

## Kenapa terasa "AI slop" — hal spesifik yang perlu dibenahi

1. **Preview skor terlalu sederhana untuk jadi curiosity gap.**
   Saat ini cuma skor total + persentil yang di-blur, dengan progress bar warna
   netral (`#e7e5e4`) yang tidak terhubung ke identitas visual brand. Tambahkan
   preview skor **per domain** (kalau data domain sudah tersedia dari model —
   cek apakah `PaywallData` atau data quiz sudah punya breakdown per domain;
   kalau belum ada di model, cukup desain UI-nya dengan data dummy/placeholder
   dan tandai TODO untuk wiring data-nya). Progress bar per domain pakai warna
   `--color-primary` dengan opacity berbeda per baris, bukan abu-abu generik, supaya
   preview terasa "menyimpan insight" bukan sekadar dikaburkan.

2. **Badge & label terasa template default.** Badge "Akses Premium" dengan dot
   kuning + badge "Terkunci" adalah pola yang sangat umum dipakai starter
   template SaaS. Pertimbangkan pendekatan yang lebih halus dan sesuai brand
   ShadowSelf (mis. label teks kecil dengan letter-spacing dan warna `--color-text`
   alih-alih pill badge berwarna).

3. **Card pembayaran kurang hierarki.** Harga, QRIS, dan tombol saat ini semua
   dalam satu card besar tanpa pemisahan visual yang jelas antara "informasi
   harga" dan "aksi pembayaran". Beri jarak/divider tipis (`--color-border`)
   antara blok harga dan blok QRIS, dan pastikan tombol "Saya Sudah Bayar"
   adalah elemen paling menonjol di halaman (bukan bersaing dengan besar font
   harga).

4. **QRIS container pakai `--color-bg-alt` (#F7F8FA) — warna hangat/beige yang tidak
   konsisten dengan sisa palet dingin (`--color-bg`, `--color-bg`, `--color-text`).** Ini
   salah satu sumber rasa "template" karena warnanya nyeleneh dari sistem token
   lain. Evaluasi apakah token ini perlu diselaraskan atau memang sengaja untuk
   membedakan area QRIS — kalau sengaja, pastikan konsisten dipakai di tempat
   lain juga.

5. **Tidak ada trust signal sama sekali** (lihat §9.2 — sudah diketahui sebagai
   gap). Tambahkan salah satu, secukupnya, tanpa berlebihan:
   - Micro-copy soal keamanan pembayaran (ikon lock kecil + teks singkat)
   - ATAU social proof singkat (jumlah user yang sudah tes) — hanya jika ada
     data real untuk itu; jangan buat angka fiktif.

6. **Animasi & micro-interaction generik.** `fade-in-up` dengan delay bertahap
   sudah ada dan itu bagus — jangan tambah animasi baru yang berlebihan
   (parallax, gradient animation, dll — itu justru bikin makin terasa "AI slop").
   Yang perlu dibenahi cukup: transisi hover pada tombol utama sedikit lebih
   halus (`transform: scale/translateY` kecil, bukan cuma shadow), dan pastikan
   `shadow-glow` pada tombol CTA tidak terlalu mencolok/neon.

7. **Tipografi hierarki price.** Harga `Rp14.900` sudah `text-4xl font-bold`,
   tapi label "Harga Spesial" dan "/ sekali" perlu dibedakan lebih jelas
   (ukuran, warna `--color-text-muted`) supaya angka harga jadi fokus tunggal, bukan
   tiga elemen teks yang ukurannya berdekatan.

## Yang TIDAK perlu dikerjakan

- Jangan tambah urgency/scarcity palsu (mis. countdown timer harga) kalau
  tidak ada mekanisme real di backend untuk itu.
- Jangan tambah testimoni/rating fiktif.
- Jangan redesign total dari nol — ini refinement, bukan rebuild. Struktur
  section (banner → header → preview → payment card → link kembali) sudah
  tepat, cukup perbaiki eksekusi visualnya.

## Output yang diharapkan

1. Perubahan pada `templ/pages/paywall_page.templ` (markup + kelas)
2. Penambahan/penyesuaian token & kelas di `assets/css/style.css` jika perlu
   (dengan penamaan konsisten mengikuti pola token yang sudah ada)
3. Ringkasan singkat: token/kelas apa saja yang baru ditambahkan, dan bagian
   mana yang masih butuh data real (misal breakdown skor per domain) sebelum
   bisa dianggap selesai sepenuhnya.
