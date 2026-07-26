# Quiz Header Component

Header/toolbar component untuk halaman tes (contoh: IQ Test). Menampilkan info tes, progress, sisa waktu, dan aksi akhiri tes. Digunakan sebagai referensi untuk agent saat generate komponen serupa.

## Layout

Container horizontal, full width, dibagi jadi 4 section sejajar (flex row, justify-between, align-items: center):

1. **Info Tes** (kiri)
2. **Progress Tes** (tengah-kiri)
3. **Waktu Tersisa** (tengah-kanan)
4. **Tombol Aksi** (kanan)

## 1. Info Tes (Left Section)

- Icon box: persegi rounded (rounded-xl), ukuran ±48x48px, background gradient ungu (dari ungu terang ke ungu tua, contoh `#7C6FF0` → `#5B4FE0`), berisi icon otak/brain berwarna putih di tengah.
- Di sebelah kanan icon, teks vertikal (stack):
  - **Judul**: "IQ Test" — font bold, ukuran ±16px, warna gelap (`#1A1A2E` / near-black).
  - **Subjudul**: "Ukur kemampuan kognitifmu" — font regular, ukuran ±13px, warna abu-abu (`#8A8A9E`).

## 2. Progress Tes (Section)

- Label kecil di atas: "Progress Tes" — font ±12px, warna abu-abu (`#8A8A9E`).
- Di bawahnya, baris horizontal berisi:
  - **Progress bar**: bentuk pill/rounded-full, lebar ±180-220px, tinggi ±8px.
    - Background track: abu-abu terang (`#E8E8F0`).
    - Fill: ungu solid (`#5B4FE0`), lebar sesuai persentase progress.
  - **Counter teks**: format `"X / Y"` (contoh: `11 / 20`) — font medium, ukuran ±13px, warna gelap, diletakkan di kanan progress bar dengan jarak (gap ±8px).

## 3. Waktu Tersisa (Timer Section)

- Card kecil dengan border rounded (rounded-xl), border tipis abu-abu (`#EAEAF2`), padding ±12px, background putih.
- Isi (flex row, align-center, gap kecil):
  - Icon box kecil: rounded, background ungu muda (`#EDEBFF`), berisi icon jam (clock) berwarna ungu.
  - Teks stack:
    - Label kecil: "Waktu Tersisa" — font ±11px, warna abu-abu.
    - Nilai waktu: format `MM:SS` (contoh: `08:10`) — font bold, ukuran ±16px, warna gelap.

## 4. Tombol Aksi (Right Section)

- Button dengan border rounded-full atau rounded-lg, border tipis abu-abu, background putih, padding horizontal ±16px, vertical ±10px.
- Isi: icon flag (bendera, outline) + label teks "Akhiri Tes".
- Warna teks & icon: gelap netral (`#1A1A2E`), hover bisa berubah jadi merah/warning untuk indikasi aksi mengakhiri.

## Container Styling (Overall)

- Background card: putih (`#FFFFFF`).
- Border radius: besar (±20-24px), full rounded card style.
- Padding: ±16-20px horizontal, ±14-16px vertical.
- Shadow: soft drop shadow tipis (opsional, `0 2px 8px rgba(0,0,0,0.04)`).
- Diletakkan di atas background halaman abu-abu muda (`#F5F5FA`).

## Dynamic Data (Props/Variables)

| Elemen | Variable | Contoh |
|---|---|---|
| Judul tes | `testTitle` | "IQ Test" |
| Subjudul | `testSubtitle` | "Ukur kemampuan kognitifmu" |
| Soal saat ini | `currentQuestion` | 11 |
| Total soal | `totalQuestions` | 20 |
| Waktu tersisa | `timeRemaining` | "08:10" |
| Label tombol aksi | `actionLabel` | "Akhiri Tes" |

## Behavior Notes

- Progress bar fill = `(currentQuestion / totalQuestions) * 100%`.
- Timer countdown berjalan real-time, format `MM:SS`, warna bisa berubah ke merah jika waktu hampir habis (misal < 1 menit).
- Tombol "Akhiri Tes" memicu konfirmasi (modal/alert) sebelum benar-benar mengakhiri tes.
- Komponen responsive: pada layar sempit, section bisa wrap atau disusun ulang jadi 2 baris (Info Tes di atas, Progress + Timer + Tombol di bawah).
