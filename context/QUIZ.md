# QUIZ.md — Spesifikasi UI Halaman Quiz

> Dokumen ini **hanya** membahas UI halaman Quiz (desktop & mobile) pada platform ShadowSelf.
> Layout, komponen, responsive behavior, UI states, interaksi, dan reusable components.
> **Tidak** membahas backend, API, database, business logic, scoring, atau implementasi timer.
> Design tokens & komponen global merujuk ke `DESIGN.md`. Flow aplikasi & struktur konten merujuk ke `IQTEST.md`.

---

## 1. Tujuan & Gambaran Umum

Halaman quiz adalah inti dari tes kemampuan kognitif bergambar 20 soal.

| Atribut | Nilai |
|---------|-------|
| **Jumlah Soal** | 20 (pilihan ganda A/B/C/D, bergambar) |
| **Total Waktu** | 20 menit (1200 detik) untuk seluruh 20 soal, hitung mundur |
| **Navigasi** | Hanya maju — tidak ada tombol mundur |
| **Auto-advance** | Saat opsi dipilih |
| **Auto-submit** | Saat waktu habis, semua jawaban dikumpulkan otomatis |
| **Alur** | Quiz → Identity (nama & email) → Submit → Redirect ke `/paywall/:id` |

UI mengacu pada **`DESIGN.md`** untuk:
- **Primary accent**: IQ Indigo (#6366f1)
- **Base**: Surface #ffffff, Bg #f8fafc
- **Typography**: Inter (body/UI), JetBrains Mono (timer)
- **Elevation**: Level 1 untuk kartu
- **Rounded**: sm (8px) untuk tombol/input, md (12px) untuk kartu

Layout **responsive** (lihat §7):
- **Desktop (≥1024px)**: Opsi 4 kolom sejajar
- **Tablet (640–1024px)**: Opsi 2×2 grid
- **Mobile (<640px)**: Opsi 1 kolom

---

## 2. Alur Halaman (UI States)

### 2.1 States

| State | Deskripsi Visual |
|-------|------------------|
| `loading` | Skeleton placeholder: shimmer pada kartu soal & opsi. |
| `question-active` | Soal aktif, timer total berjalan (20 menit), opsi dapat dipilih. |
| `selected` | Opsi terpilih: border `--color-primary-border`, background `--color-primary-light`. |
| `answered` | Opsi tetap terlihat. Auto-advance ke soal berikutnya. |
| `time-warning-30` | Timer warna **warning** (#f59e0b) saat total sisa ≤30 detik. |
| `time-warning-10` | Timer warna **error** (#ef4444) saat total sisa ≤10 detik. |
| `time-up` | Timer mencapai 0 → auto-submit semua jawaban, soal belum terjawab mendapat skor 0. |
| `identity` | Form input nama & email (setelah quiz selesai). Tombol "Simpan & Lanjutkan". |
| `submitting` | Overlay + spinner, semua interaksi terkunci. |
| `completed` | Semua data terkirim → redirect ke `/paywall/:id`. |

### 2.2 Alur Lengkap

```
┌──────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Quiz    │───▶│   Identity   │───▶│   Submit     │───▶│ Redirect ke  │
│ (20 soal │    │  (Nama &     │    │  (Loading)   │    │ /paywall/:id │
│  20 min) │    │   Email)     │    │              │    │              │
└──────────┘    └──────────────┘    └──────────────┘    └──────────────┘
```

### 2.3 State Diagram Interaksi

```
[load halaman] ──▶ loading ──[data siap]──▶ question-active
                                                               │
                                                   ┌──────────┤
                                                   ▼          ▼
                                              selected   time-warning-30 (sisa ≤30dtk)
                                                   │          │
                                                   ▼          ▼
                                              answered   time-warning-10 (sisa ≤10dtk)
                                                   │          │
                                                   ▼          ▼
                                          [soal<20]    time-up (sisa 0dtk)
                                               │          │
                                               ▼          ▼
                                        auto-advance auto-submit seluruh jawaban
                                        (timer lanjut)     │
                                               │          │
                                               └────┬─────┘
                                                    ▼
                                            question-active (next, timer total lanjut)
                                                    │
                                        [soal=20 & semua terjawab ATAU time-up]
                                                    ▼
                                              identity (post-quiz)
                                                    │
                                          [isi nama & email]
                                                    ▼
                                              submitting ──▶ completed
```

Catatan penting:
- **Timer bersifat global**: 20 menit dihitung sejak soal pertama muncul, tidak di-reset per soal.
- **Saat time-up**: Semua soal yang belum terjawab otomatis dicatat dengan skor 0. Quiz berakhir dan langsung beralih ke identity step.
- **Identity step muncul setelah quiz selesai**: Baik karena semua soal terjawab maupun karena waktu habis.

---

## 3. Struktur Layout & Komponen

### 3.1 Hirarki Komponen

```
<QuizPage>
  ├─ <QuizSession>                          (state: quiz/loading)
  │   ├─ <QuizHeader>
  │   │   ├─ <ProgressBar current={n} total={20} />
  │   │   └─ <TimerDisplay seconds={remaining} variant={normal|warning|danger} />
  │   ├─ <QuestionCard>
  │   │   ├─ <QuestionCounter text="Soal {n} dari 20" />
  │   │   ├─ <Divider />
  │   │   └─ <QuestionImage src={url} alt="Soal nomor {n}" />
  │   └─ <OptionsCard>
  │       ├─ CardTitle subtle: "Pilih jawaban yang tepat:"
  │       └─ <OptionGrid>
  │           ├─ <OptionItem label="A" imageSrc={url} state={default|selected|disabled} />
  │           ├─ <OptionItem label="B" ... />
  │           ├─ <OptionItem label="C" ... />
  │           └─ <OptionItem label="D" ... />
  │           </OptionGrid>
  │
  ├─ <IdentityStep>                         (state: 'identity', post-quiz)
  │   └─ <Card>
  │       ├─ Heading "Simpan Hasil Tes Anda"
  │       ├─ <InputField label="Nama Lengkap" />
  │       ├─ <InputField label="Email" />
  │       ├─ Ringkasan jawaban (opsional: jumlah terjawab / total)
  │       └─ <ButtonPrimary label="Simpan & Lanjutkan" />
  │
  └─ <SubmittingOverlay>                    (state: 'submitting')
      └─ Spinner + "Menyimpan jawaban..."
```

Catatan:
- **Timer global** berada di `QuizHeader`, menyatu dengan ProgressBar — bukan di dalam QuestionCard.
- Tidak ada navigator soal (grid angka 1–20). Cukup progress bar + question counter.
- Tidak ada tombol navigasi maju/mundur. Pilih opsi → auto-advance; timer habis → auto-submit.
- Tidak ada tombol "Lewati". Soal yang tidak terjawab saat time-up mendapat skor 0.

### 3.2 Identity Step — Layout (setelah quiz)

```
[Page background: --color-bg]
┌─────────────────────────────────────────────┐
│  ┌─────────────────────────────────────┐    │
│  │  "Simpan Hasil Tes Anda"           │    │
│  │  ────────────────────────────────   │    │
│  │                                     │    │
│  │  ✅ Anda telah menjawab 18 dari 20  │    │
│  │     soal.                           │    │
│  │                                     │    │
│  │  Nama Lengkap  [________________]   │    │
│  │  Email         [________________]   │    │
│  │                                     │    │
│  │  Masukkan data diri untuk           │    │
│  │  menyimpan hasil tes Anda.          │    │
│  │                                     │    │
│  │  ┌───────────────────────────┐      │    │
│  │  │   Simpan & Lanjutkan      │      │    │
│  │  └───────────────────────────┘      │    │
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

- Kartu putih (`--color-surface`), shadow level 1, max-width 520px, center.
- Input fields: border `1px solid #e2e8f0`, radius 8px, padding 12px 16px. Focus: border `#6366f1`.
- Tombol primary indigo: `#6366f1` bg, white text, radius 8px, padding 12px 28px. Hover: `#4f46e5`.
- Ringkasan jumlah soal terjawab (opsional): Inter 0.875rem, `--color-ink-subtle`.

### 3.3 Quiz Step — Layout Mobile (<640px)

```
[Page background: --color-bg]
┌──────────────────────────────────────┐
│ QUIZ HEADER                          │
│  ████████░░░░░░░░░░░░  5/20         │
│  ⏱ 15:32                             │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ QUESTION CARD                        │
│  "Soal 3 dari 20"                    │
│  ───────────────────────────        │
│  [Gambar soal]                       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ OPTIONS CARD (single column)         │
│  ┌──────────────────────────────┐    │
│  │ [A] Gambar A                  │    │
│  ├──────────────────────────────┤    │
│  │ [B] Gambar B                  │    │
│  ├──────────────────────────────┤    │
│  │ [C] Gambar C                  │    │
│  ├──────────────────────────────┤    │
│  │ [D] Gambar D                  │    │
│  └──────────────────────────────┘    │
└──────────────────────────────────────┘
```

### 3.4 Quiz Step — Layout Tablet (640–1024px)

```
┌──────────────────────────────────────┐
│ QUIZ HEADER                          │
│  ████████░░░░░░░░░░░░  5/20         │
│  ⏱ 15:32                             │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ QUESTION CARD                        │
│  "Soal 3 dari 20"                    │
│  ───────────────────────────        │
│  [Gambar soal — center]             │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ OPTIONS CARD (2×2 grid)              │
│  ┌─────────┐  ┌─────────┐           │
│  │ [A] img │  │ [B] img │           │
│  └─────────┘  └─────────┘           │
│  ┌─────────┐  ┌─────────┐           │
│  │ [C] img │  │ [D] img │           │
│  └─────────┘  └─────────┘           │
└──────────────────────────────────────┘
```

### 3.5 Quiz Step — Layout Desktop (≥1024px)

```
[Container max-width: 900px, center]
┌──────────────────────────────────────┐
│ QUIZ HEADER                          │
│  ████████░░░░░░░░░░░░  5/20         │
│  ⏱ 15:32                             │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ QUESTION CARD                        │
│  "Soal 3 dari 20"                    │
│  ───────────────────────────        │
│  [Gambar soal — center, max 500px]  │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ OPTIONS CARD (4 kolom sejajar)       │
│  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐│
│  │ [A] │  │ [B] │  │ [C] │  │ [D] ││
│  │ img │  │ img │  │ img │  │ img ││
│  └─────┘  └─────┘  └─────┘  └─────┘│
└──────────────────────────────────────┘
```

---

## 4. Spesifikasi Detail Komponen

### 4.1 `<ProgressBar>`

- Posisi: sticky top (opsional), full width container.
- Layout: flex row, `align-items: center; gap: 12px`.
- Track: `height: 6px; background: #f1f5f9; border-radius: 9999px; flex: 1; overflow: hidden`.
- Fill: `height: 100%; background: #6366f1; border-radius: 9999px; transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1)`.
- Label: `"5/20"`, font JetBrains Mono 0.875rem, weight 600, color `--color-ink-muted`.
- Aksesibilitas: `role="progressbar"`, `aria-valuenow`, `aria-valuemin="0"`, `aria-valuemax="20"`, `aria-label="Progress soal: 5 dari 20"`.

### 4.2 `<QuizHeader>`

- Posisi: sticky top, full width container, z-index: 10.
- Layout: flex column, gap 8px.
- Berisi ProgressBar dan TimerDisplay secara horizontal.
- Background: `--color-surface` atau transparan dengan backdrop blur.
- Padding: `12px 16px` (mobile), `12px 24px` (desktop).

### 4.3 `<QuestionCard>`

- Background: `--color-surface`, radius `--rounded-md` (12px), shadow level 1, border `1px solid #f1f5f9`.
- Padding: `24px` (desktop), `16px` (mobile).
- Layout: flex column, gap `16px`.
- **QuestionCounter**: `"Soal {n} dari 20"` (Inter 0.875rem, 500, ink-muted). Timer TIDAK di sini — timer ada di QuizHeader global.
- **QuestionImage**: max-height `400px` atau `50vh` (mana lebih kecil), center dalam container. Loading: skeleton shimmer. Error: broken image icon + "Gambar tidak tersedia".

### 4.4 `<TimerDisplay>`

- Font: JetBrains Mono, size 0.875rem, weight 600.
- Format: `MM:SS` (e.g., `19:47`, `15:32`, `00:28`, `00:08`).
- Menampilkan **sisa waktu total** (20 menit awal → 0). Tidak di-reset per soal.
- 3 variant warna:
  - `normal` (>30s): color `--color-ink-muted` (#475569).
  - `warning` (≤30s): color `--color-warning` (#f59e0b).
  - `danger` (≤10s): color `--color-error` (#ef4444), weight 700.
- Ikon jam/stopwatch SVG kecil di kiri teks (optional, 16px, warna mengikuti variant).
- Aksesibilitas: `aria-live="polite"`, `aria-atomic="true"`, `aria-label="Sisa waktu: 15 menit 32 detik"`.

### 4.5 `<OptionsCard>`

- Background: `--color-surface`, radius `--rounded-md`, shadow level 1, border `1px solid #f1f5f9`.
- Padding: `24px` (desktop), `16px` (mobile).
- CardTitle subtle: `"Pilih jawaban yang tepat:"` — Inter 0.875rem, 500, ink-muted.

### 4.6 `<OptionGrid>` & `<OptionItem>`

- CSS Grid (mobile-first):
  ```css
  .option-grid {
    display: grid;
    grid-template-columns: 1fr;         /* Mobile <640px */
    gap: 12px;
  }
  @media (min-width: 640px) {
    .option-grid { grid-template-columns: repeat(2, 1fr); gap: 16px; }
  }
  @media (min-width: 1024px) {
    .option-grid { grid-template-columns: repeat(4, 1fr); gap: 16px; }
  }
  ```
- **OptionItem** container: border `1px solid #e2e8f0`, border-radius 8px, padding 16px, cursor pointer, transition all 0.15s ease.
- Label huruf: "A", "B", "C", "D" — bulat/kotak kecil di kiri atas, Inter semibold 0.875rem, `--color-ink-muted`.
- Gambar opsi: `max-width: 100%; height: auto; object-fit: contain; margin-top: 8px`.
- 4 state visual:
  - **default**: border `#e2e8f0`, bg `#ffffff`.
  - **hover** (desktop only): border `#c7d2fe`, bg `#f8fafc`.
  - **selected**: border `2px solid #c7d2fe`, bg `#eef2ff`; label huruf jadi `#6366f1`.
  - **disabled**: opacity 0.5, cursor not-allowed.
- Aksesibilitas: setiap item adalah `<button>` dengan `role="radio"`, `aria-checked`, `aria-label="Pilih jawaban A"`. Grid menggunakan `role="radiogroup"`.

### 4.7 Loading State (Skeleton)

- QuestionCard skeleton: 2 bar shimmer (header + gambar persegi ~300×300px).
- OptionsCard skeleton: 4 persegi panjang ~150×150px dengan shimmer.
- Animasi: linear gradient bergerak kiri→kanan, warna `--color-bg-alt` → `--color-bord-light`, durasi 1.5s infinite.
- `prefers-reduced-motion`: shimmer dimatikan.

### 4.8 Submitting State

- Overlay: `position: fixed; inset: 0; background: rgba(255,255,255,0.8); z-index: 50; display: flex; align-items: center; justify-content: center`.
- Spinner: SVG lingkaran berputar (indigo), 48px, + teks `"Menyimpan jawaban..."` (Inter 0.875rem, ink-muted).
- Semua interaksi terkunci (`pointer-events: none`).

---

## 5. Responsive Behavior

| Breakpoint | Layout Opsi | Padding Kartu | Gap Grid | Gambar Soal |
|------------|-------------|---------------|----------|-------------|
| <640px (mobile) | 1 kolom | 16px | 12px | Max lebar kartu, height auto |
| 640–1024px (tablet) | 2×2 grid | 20px | 16px | Max 400px width |
| ≥1024px (desktop) | 4 kolom sejajar | 24px | 16px | Max 500px width, center |

Container halaman: `max-width: 900px; margin: 0 auto; padding: 16px` (mobile), `padding: 24px` (desktop).
Progress bar + Timer: full width container (tidak terikat max-width 900px, jika ingin stretch).

---

## 6. Interaksi & Aksesibilitas

### 6.1 Interaksi per Komponen

| Komponen | Trigger | Response |
|----------|---------|----------|
| Load halaman quiz | DOMContentLoaded | Fetch data soal, tampilkan skeleton, transisi ke `question-active`. Timer 20 menit dimulai. |
| OptionItem | Click / Enter / Space | Set `selected` state pada item. Delay 300ms (feedback visual), lalu auto-advance ke soal berikutnya. Timer total tetap lanjut — tidak di-reset. |
| TimerDisplay | `timeRemaining === 0` | **Auto-submit seluruh jawaban**. Semua soal yang belum terjawab mendapat skor 0. Langsung beralih ke identity step. |
| TimerDisplay | `timeRemaining <= 30` | Timer berubah ke variant `warning` (kuning). |
| TimerDisplay | `timeRemaining <= 10` | Timer berubah ke variant `danger` (merah). |
| Transisi antar soal | Auto-advance | Fade-out kartu lama (150ms) → ganti data → fade-in kartu baru (150ms). Timer total tidak terpengaruh. |
| ProgressBar | Soal berubah | Fill width = `(currentIndex + 1) / 20 × 100%`, transisi 0.6s. |
| ButtonPrimary "Simpan & Lanjutkan" | Click | Validasi nama & email. Jika valid → `submitting` state, kirim data + jawaban, redirect. Jika invalid → tampilkan error pada field. |
| InputField (identity step) | Submit form | Validasi real-time: border merah (`--color-error`) jika kosong/salah format. Pesan error di bawah field. |

### 6.2 Aksesibilitas

- Semua click events juga merespon `keydown.Enter` dan `keydown.Space`.
- `OptionItem`: `role="radio"`, `aria-checked`.
- `OptionGrid`: `role="radiogroup"`.
- Timer: `aria-live="polite"` untuk update countdown.
- Focus management: saat soal berganti, fokus ke `OptionGrid` (first item).
- `prefers-reduced-motion`: semua animasi (shimmer, transisi kartu, hover) dinonaktifkan. Timer tetap update numerik tanpa animasi.
- Semua gambar soal & opsi memiliki `alt` text deskriptif (tanpa membocorkan jawaban benar).
- WCAG AA minimum pada semua elemen teks. Focus-visible: 2px indigo outline + 2px offset.

---

## 7. Reusable Components

| Komponen | Props | Catatan |
|----------|-------|---------|
| `<Card>` | `padding?`, `maxWidth?` | Container putih, border + shadow level 1, radius md. |
| `<CardTitle>` | `variant?: 'default' \| 'subtle'` | Default: lg, semibold, ink. Subtle: sm, 500, ink-muted. |
| `<Divider>` | — | `1px solid #e2e8f0`, margin vertical md. |
| `<InputField>` | `label`, `name`, `type?`, `placeholder?`, `error?` | Sesuai spec input DESIGN.md. |
| `<ButtonPrimary>` | `label`, `disabled?`, `loading?` | Indigo bg, white text, 8px radius, 12px 28px. Hover: #4f46e5. Focus: 2px indigo outline + 2px offset. |
| `<ButtonGhost>` | `label`, `disabled?` | Transparent, ink-muted text, 8px 16px. Hover: bg #f8fafc. |
| `<ProgressBar>` | `current`, `total` | 6px track, indigo fill, label X/total, role progressbar. |
| `<TimerDisplay>` | `seconds`, `variant?: 'normal' \| 'warning' \| 'danger'` | Mono font, MM:SS, 3 warna. aria-live polite. Menampilkan sisa waktu total 20 menit. |
| `<Skeleton>` | `width`, `height`, `borderRadius?`, `count?` | Shimmer animasi, reduced-motion: off. |
| `<QuestionImage>` | `src`, `alt` | Max-width 100%, contain, center. Loading: skeleton. Error: fallback. |
| `<OptionItem>` | `label`, `imageSrc`, `state?: 'default' \| 'selected' \| 'disabled'`, `onSelect` | 4 state visual, role radio, aria-checked. |

---

## 8. Urutan Implementasi (untuk AI Agent)

1. **Setup CSS variables**: pastikan design tokens dari `DESIGN.md` sudah terdefinisi di root stylesheet. Import Google Fonts: Inter, DM Serif Display (global), JetBrains Mono.

2. **Bangun komponen atomik** (urutan):
   a. `<Card>` — container reusable.
   b. `<CardTitle>` — heading dengan variant.
   c. `<Divider>` — garis horizontal.
   d. `<InputField>` — input + label + error state.
   e. `<ButtonPrimary>` & `<ButtonGhost>` — sesuai spec DESIGN.md.
   f. `<ProgressBar>` — bar + fill animasi.
   g. `<TimerDisplay>` — mono countdown, 3 variant.
   h. `<Skeleton>` — placeholder shimmer.

3. **Bangun `QuestionCard`**: QuestionCounter → Divider → QuestionImage. (Timer tidak ada di sini.)

4. **Bangun `QuizHeader`**: ProgressBar + TimerDisplay secara horizontal. (Timer global, sticky top.)

5. **Bangun `OptionItem`**: container border, label huruf, gambar, 4 state visual.

6. **Bangun `OptionsCard`**: CardTitle subtle + OptionGrid (responsive grid 1/2/4 kolom).

7. **Rakit `QuizSession`**: QuizHeader → QuestionCard → OptionsCard. (Timer global di QuizHeader.)

8. **Bangun `IdentityStep`** (post-quiz): Card → ringkasan jawaban → InputField (nama, email) → ButtonPrimary "Simpan & Lanjutkan".

9. **Rakit `QuizPage`**: conditional render QuizSession vs IdentityStep berdasarkan state.

10. **Implement state management** (Alpine.js): state `step` (`loading|question-active|identity|submitting|completed`), `currentQuestion`, `totalTimeRemaining` (1200 detik), `answers[]`, `questions[]`, `isSubmitting`. Actions: `loadQuiz()`, `selectAnswer()`, `autoAdvance()`, `nextQuestion()`, `submitQuiz()`. Timer global dimulai saat halaman di-load: `setInterval` 1 detik, decrement `totalTimeRemaining`, auto-submit jika mencapai 0. Lihat `IQTEST.md` §12.2.2 untuk struktur data Alpine.js.

11. **Implement responsive**: CSS media queries untuk breakpoint mobile/tablet/desktop. Test di 375px, 768px, 1024px.

12. **Tambahkan aksesibilitas**: ARIA labels, role, keyboard navigation, focus management, `prefers-reduced-motion`.

13. **Tambahkan loading & submitting states**: Skeleton untuk loading, overlay + spinner untuk submitting.

14. **Polish visual**: bandingkan dengan `DESIGN.md` — cek warna, radius, shadow, spacing, typography.

---

## 9. Aturan Visual (dari DESIGN.md) yang Relevan untuk Halaman Quiz

- **IQ Indigo (#6366f1)** adalah satu-satunya accent untuk elemen interaktif.
- **Tidak ada gradient** — background, text, atau button.
- **Tidak ada glassmorphism** (kecuali navbar scroll, yang tidak ada di halaman quiz).
- **Tidak ada ilustrasi, Lottie, atau SVG sketches**.
- **Tidak ada pure black (#000000)** — gunakan ink (#0f172a).
- **Surface selalu #ffffff** — bukan off-white.
- **Body text di ink-muted (#475569)** — WCAG AA.
- **Max body line length: 65–75ch**.
- **Card tidak flat** — gunakan Level 1 elevation.
- **prefers-reduced-motion** dihormati untuk semua animasi.
- **Color tidak pernah jadi sole differentiator** — selalu pair dengan shape/label/posisi.

---

*Dokumen ini hanya membahas UI halaman Quiz. Untuk design system lengkap, lihat `DESIGN.md`. Untuk flow aplikasi & spesifikasi teknis, lihat `IQTEST.md`.*