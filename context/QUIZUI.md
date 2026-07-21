# QUIZUI.md — Spesifikasi UI Halaman Quiz (Pattern Recognition Test)

> Spesifikasi ini harus diikuti **secara literal**. Jangan gunakan default theme/komponen bawaan library UI (warna primary bawaan, card style bawaan, badge bawaan) — semua warna, layout, dan komponen harus eksplisit sesuai dokumen ini.

---

## 1. Gambaran Umum

Halaman quiz pola visual (pattern/shape recognition), satu soal per halaman, dengan format:
- Grid **3x3** berisi bentuk geometris yang membentuk pola, sel terakhir kosong ditandai `?`.
- Pengguna memilih jawaban dari **6 opsi (A–F)**.
- Ada timer hitung mundur, nomor soal ("4/22"), tombol navigasi, dan navigator soal di bawah berupa grid nomor 1–22.

Layout **responsive**:
- **Desktop (≥1024px)**: dua kartu putih berdampingan — kartu soal (kiri, lebih lebar) & kartu jawaban (kanan, lebih sempit).
- **Mobile (<768px)**: dua kartu ditumpuk vertikal — kartu soal di atas, kartu jawaban di bawah.

---

## 2. Design Tokens

### 2.1 Warna

| Token | Hex | Penggunaan |
|---|---|---|
| `--color-navy-900` | `#16324F` | Background header, teks judul, **semua outline shape**, label huruf jawaban |
| `--color-orange-500` | `#F5821F` | Tombol aksi, lingkaran placeholder `?`, aksen selected |
| `--color-orange-400` | `#F7941D` | Hover state tombol orange |
| `--color-white` | `#FFFFFF` | Background kartu |
| `--color-bg-page` | `#F2F3F5` | Background halaman (abu sangat muda netral) |
| `--color-gray-900` | `#333333` | Background navigator bawah |
| `--color-gray-300` | `#D9DEE4` | Border tipis / divider |
| `--color-gray-500` | `#8A8F98` | Teks sekunder |

Hanya 4 warna yang boleh dominan tampil: **navy, orange, putih, abu**. Tidak ada warna lain (ungu, biru terang, hijau, dll) di mana pun pada UI.

### 2.2 Tipografi

- Font: rounded sans-serif bold — `'Poppins', 'Nunito', system-ui, sans-serif`.
- Heading kartu & label huruf jawaban: **Bold/700**.
- Angka navigator & teks tombol: **SemiBold/600–700**.
- Ukuran: header ~15px, judul kartu ~17px, label opsi ~16px, angka navigator ~14px.

### 2.3 Spacing & Radius

- Radius kartu besar: `16–20px`.
- Radius header & navigator bawah: `16–20px`.
- Radius tombol: `9999px` (pill).
- Shape jawaban/soal: **sudut asli bentuk** (kotak = sudut lancip, bukan rounded).
- Padding kartu: desktop `24–28px`, mobile `16–20px`.
- Gap grid 3x3 (soal): `20–28px`.
- Gap grid jawaban (2 kolom x 3 baris): gap kolom `32px`, gap baris `20px`.

### 2.4 Shadow

- Kartu putih: `box-shadow: 0 4px 12px rgba(0,0,0,0.06)`.
- Header & navigator bawah: flat, tanpa shadow mencolok.

---

## 3. Struktur Layout

```
┌──────────────────────────────────────────────────────────┐
│ HEADER (navy solid, rounded, full width)                  │
│  "Pertanyaan No: {n}/{total}"        🕐 {mm:ss}            │
└──────────────────────────────────────────────────────────┘

DESKTOP (≥1024px):                     MOBILE (<768px):
┌───────────────┐ ┌──────────────┐     ┌────────────────────┐
│ Kartu Soal     │ │ Kartu Jawaban│     │ Kartu Soal          │
│ "Bentuk        │ │ "Pilih       │     │ (3x3 grid)          │
│  manakah yang  │ │  jawaban:"   │     ├────────────────────┤
│  hilang?"      │ │ (2x3 grid,   │     │ Kartu Jawaban       │
│ (3x3 grid)     │ │  A-F)        │     │ (2x3 grid, A-F)     │
└───────────────┘ └──────────────┘     └────────────────────┘

        [Kembali]  [Lewati Pertanyaan]      ← center, pill orange

┌──────────────────────────────────────────────────────────┐
│ NAVIGATOR (dark bar, rounded, grid nomor 1..N)             │
└──────────────────────────────────────────────────────────┘
```

**Aturan wajib:**
- Header **selalu** solid navy, rounded, tanpa progress bar tipis di bawahnya.
- Info soal cukup **satu baris** teks: `"Pertanyaan No: {n}/{total}"` — tidak ada subtitle duplikat, tidak ada badge kategori/tag.
- Timer menyatu dalam header (background sama, tanpa pill/border terpisah), icon stopwatch putih kecil + teks putih bold.
- Kartu soal berisi **grid 3x3 penuh (9 sel)** — bukan 1 shape tunggal di dalam box abu-abu.
- Tidak ada box/container abu-abu tambahan yang membungkus shape di dalam kartu.
- Tidak ada caption/ID teknis (seperti nama file atau kode data) ditampilkan di UI.
- Judul kartu kanan persis: **"Pilih jawaban:"**.
- Opsi jawaban selalu **6 buah (A–F)**, layout **2 kolom x 3 baris**, bukan 4 opsi 1 baris.
- Setiap opsi = label huruf (`A)`, `B)`, dst) sejajar horizontal di kiri shape — **tanpa card/border pembungkus per item** kecuali untuk state hover/selected.
- Tombol "Kembali" & "Lewati Pertanyaan" **selalu tampil**, center, di bawah kedua kartu.
- Navigator bawah (grid nomor soal) **selalu tampil** di bagian paling bawah halaman.

---

## 4. Hierarki Komponen

```
<QuizPage>
 ├─ <QuizHeader>
 │   ├─ <QuestionCounter current={4} total={22} />
 │   └─ <Timer seconds={remainingSeconds} />
 │
 ├─ <QuizBody>  (flex row di desktop, flex column di mobile)
 │   ├─ <QuestionCard>
 │   │   ├─ <CardTitle text="Bentuk manakah yang hilang?" />
 │   │   ├─ <Divider />
 │   │   └─ <ShapeGrid columns={3} rows={3}>
 │   │        ├─ <ShapeIcon ... />        (x8, pola berbeda tiap sel)
 │   │        └─ <UnknownPlaceholder />   (sel terakhir, ikon "?")
 │   │       </ShapeGrid>
 │   │
 │   └─ <AnswerCard>
 │       ├─ <CardTitle text="Pilih jawaban:" />
 │       ├─ <Divider />
 │       └─ <AnswerOptionGrid columns={2} rows={3}>
 │            ├─ <AnswerOption label="A" shapeProps={...} state="default|hover|selected" />
 │            ├─ ... sampai F
 │           </AnswerOptionGrid>
 │
 ├─ <ActionButtons>
 │   ├─ <ButtonSolid label="Kembali" />
 │   └─ <ButtonSolid label="Lewati Pertanyaan" />
 │
 └─ <QuestionNavigator>
     └─ <NavGrid>
          ├─ <NavItem number={1} status="answered" answerLabel="E" />
          ├─ <NavItem number={4} status="active" />
          ├─ <NavItem number={5..22} status="unanswered" />
         </NavGrid>
```

---

## 5. Komponen Shape — `<ShapeIcon>`

Shape **tidak selalu pie chart** — bisa berupa kombinasi bentuk geometris dasar. Komponen harus fleksibel:

```ts
type ShapeType = "circle" | "square" | "diamond" | "pie";

interface ShapeIconProps {
  type: ShapeType;
  size: number;              // px, 48-72
  hasCenterDot?: boolean;    // titik solid navy di tengah
  crossHatch?: boolean;      // motif garis silang diagonal di area shape
  nested?: ShapeIconProps[]; // shape lain yang ditumpuk konsentris di dalamnya
  filledIndexes?: number[];  // khusus type "pie": index segmen berwarna orange
}
```

### Aturan render:

- **`circle`**: outline navy `stroke-width: 2px`, fill putih/transparan, sudut penuh melingkar.
- **`square`**: outline navy, sudut lancip (radius 0, **bukan rounded**).
- **`diamond`**: sama seperti `square` tapi dirotasi 45°.
- **`hasCenterDot: true`**: tambahkan lingkaran kecil solid navy (~18–22% dari ukuran shape) tepat di tengah.
- **Shape bertumpuk (`nested`)**: dua atau lebih shape konsentris dalam satu ikon, misal diamond di dalam circle, atau circle+diamond+dot bertumpuk 3 lapis — semua outline navy tipis, saling center-align.
- **`crossHatch: true`**: area shape diisi pola garis diagonal tipis navy (`stroke-width: 1px`, jarak antar garis ~4px, sudut 45°) menggunakan SVG `<pattern>` — dipakai untuk salah satu varian opsi jawaban bermotif arsir.
- **`pie`**: lingkaran dibagi 10 juring (36°/juring), outline navy tipis antar juring, juring pada `filledIndexes` diisi solid orange, sisanya kosong.

### Prinsip:
- Shape **tidak pernah** solid fill warna cerah (orange/ungu/dsb) kecuali dot kecil di tengah dan tipe `pie`.
- Semua shape dalam satu grid soal berukuran sama; shape opsi jawaban sedikit lebih kecil dari shape grid soal.
- Warna satu-satunya elemen orange solid di seluruh halaman: placeholder `?`, tombol aksi, state selected, dan juring pie (jika tipe soal pie).

---

## 6. Spesifikasi Detail Komponen

### 6.1 `QuizHeader`
- Background `--color-navy-900`, `border-radius: 16px`, padding `14px 20px`, flex `justify-content: space-between; align-items: center`.
- Kiri: `"Pertanyaan No: {current}/{total}"` — putih, bold, ~15px.
- Kanan: icon stopwatch putih ~16px + teks waktu format `MM:SS`, putih bold.

### 6.2 `QuestionCard` & `AnswerCard`
- Background putih, `border-radius: 16–20px`, padding `20–28px`, shadow lembut.
- `CardTitle`: navy bold ~17px + divider tipis (`height:2px; background:#DDE2E8`) di bawahnya, margin-bottom ~16–20px sebelum grid.

### 6.3 `ShapeGrid` (3x3, kartu soal)
- CSS grid `repeat(3, 1fr)`, gap `20–28px`, item center align.
- Sel terakhir → `<UnknownPlaceholder>`: lingkaran solid orange, ukuran sama dengan shape lain, isi `?` putih tebal center, tanpa outline.

### 6.4 `AnswerOptionGrid` & `AnswerOption`
- Grid 2 kolom x 3 baris (berlaku di desktop maupun mobile).
- Tiap `AnswerOption`: flex row, label huruf (`"A)"`, dst — bold, navy/dark gray, ~16px) di kiri, `ShapeIcon` di kanan, gap ~12px.
- Tanpa card/border/background di state default.
- **Hover** (desktop): background sangat muda (mis. `#F5F7FA`), radius `10px`, cursor pointer.
- **Selected**: border orange 2px rounded di sekeliling label+shape, atau background orange sangat muda.

### 6.5 `ActionButtons`
- Flex row, gap `16px`, center, margin-top setelah kartu ~20px.
- Kedua tombol: `<ButtonSolid>` — background `orange-500`, teks putih bold, padding `10px 28px`, radius `9999px`, shadow tipis. Hover → `orange-400`.
- Label persis: `"Kembali"` dan `"Lewati Pertanyaan"`.

### 6.6 `QuestionNavigator`
- Background `--color-gray-900`, radius `16–20px`, padding `16–20px`.
- Grid nomor: `flex-wrap: wrap`, gap `12–16px`.
- `NavItem` default: teks putih/abu terang bold ~14px, tanpa background.
- `NavItem` answered: nomor + superscript huruf jawaban kecil di kanan-atas (mis. "1" + "ᴬ").
- `NavItem` active (soal saat ini): dibungkus lingkaran putih solid, teks navy, `border-radius: 50%`, ~32px.
- Semua item clickable untuk lompat ke soal terkait.

---

## 7. State & Interaksi

| State | Perilaku |
|---|---|
| `question-active` | Tampilan normal sesuai spesifikasi di atas. |
| `hover` (opsi jawaban, desktop) | Background muda muncul saat kursor di atas opsi. |
| `selected` | Opsi yang diklik mendapat border/background orange; opsi lain kembali ke default (single-select). |
| `answered` (navigator) | Nomor soal menampilkan superscript huruf jawaban. |
| `active` (navigator) | Nomor soal saat ini dilingkari putih solid. |
| `loading` | Skeleton abu-abu berbentuk kartu & grid saat data belum siap. |

Interaksi:
- Klik opsi jawaban → set `selected`, update state single-select.
- Klik `"Kembali"` → soal index-1, update header & navigator.
- Klik `"Lewati Pertanyaan"` → soal index+1 tanpa menandai jawaban.
- Klik nomor di navigator → lompat langsung ke soal tsb, update `active`.
- Timer berjalan mundur tiap detik, format `MM:SS`.

---

## 8. Checklist Validasi Visual

- [ ] Background halaman abu muda netral (`#F2F3F5`), bukan warna lain.
- [ ] Hanya 4 warna dominan: navy, orange, putih, abu.
- [ ] Header navy solid rounded, tanpa progress bar, berisi info soal + timer menyatu.
- [ ] Desktop: kartu soal & jawaban sejajar 2 kolom. Mobile: ditumpuk vertikal.
- [ ] Kartu soal = grid 3x3 penuh, tanpa box pembungkus abu-abu tambahan.
- [ ] Sel terakhir grid soal = lingkaran orange solid berisi "?".
- [ ] Tidak ada caption/ID teknis di UI.
- [ ] Judul kartu kanan persis "Pilih jawaban:".
- [ ] 4 opsi jawaban (A–F), grid 2 kolom x 2 baris, tanpa card per item.
- [ ] Shape berupa outline navy (bukan solid warna terang), mendukung circle/square/diamond/dot/nested/crosshatch.
- [ ] 2 tombol pill orange ("Kembali", "Lewati Pertanyaan") tampil & center.
- [ ] Navigator bawah dark bar tampil dengan grid nomor lengkap, status active & answered berbeda visual.
- [ ] Font bold & rounded konsisten di semua teks.
- [ ] Radius kartu besar (16–20px) dengan shadow lembut terlihat jelas.
