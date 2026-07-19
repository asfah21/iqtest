# IQTEST.md — IQ Test Engine
## Complete Technical & Functional Specification
### Version: 2.0 | Status: Draft (Revisi) | Last Updated: 2026-07-19

> **Catatan revisi:** Tes IQ yang sah mengukur *kemampuan kognitif* melalui soal dengan **satu jawaban benar objektif**, dinilai berdasarkan performa aktual, bukan preferensi yang dilaporkan sendiri. Versi 2.0 ini merombak metodologi inti (Bagian 2–9) agar sesuai standar seperti Raven's Progressive Matrices, WAIS, dan CFIT, sambil mempertahankan arsitektur sistem (Go/Gin/PostgreSQL) yang sudah relevan.

---

## TABLE OF CONTENTS

1. [System Overview](#1-system-overview)
2. [Assessment Methodology](#2-assessment-methodology)
3. [Psychometric Foundation](#3-psychometric-foundation)
4. [Question Structure](#4-question-structure)
5. [Timer Rules](#5-timer-rules)
6. [Scoring Algorithm](#6-scoring-algorithm)
7. [IQ Score Conversion](#7-iq-score-conversion)
8. [Result Interpretation](#8-result-interpretation)
9. [Anti-Cheating Strategy](#9-anti-cheating-strategy)
10. [Database Model](#10-database-model)
11. [API Flow](#11-api-flow)
12. [UI/UX Flow](#12-uiux-flow)
13. [Future Improvements](#13-future-improvements)
14. [Appendices](#14-appendices)

---

## 1. SYSTEM OVERVIEW

### 1.1 Platform Identity

| Attribute | Value |
|-----------|-------|
| **Product Name** | ShadowSelf |
| **Domain** | Cognitive Ability Assessment |
| **Core Test** | 20-item visual multiple-choice cognitive ability test |
| **Framework** | Go (Gin) + PostgreSQL + templ |
| **Target Audience** | General public, individuals seeking self-understanding |
| **Monetization** | Freemium — tes gratis, hasil lengkap berbayar (IDR 14.900) |

### 1.2 Architecture Diagram (High-Level)

*(tidak berubah dari v1.0 — lihat lampiran)*

```
Browser (Alpine.js) → Gin Router → Handlers (Go) → Services → Repositories → PostgreSQL
```

### 1.3 Technology Stack

Tidak berubah dari v1.0 (Go 1.25, Gin v1.12, templ v0.3, PostgreSQL, Alpine.js, Docker).

---

## 2. ASSESSMENT METHODOLOGY

### 2.1 Perubahan Mendasar dari v1.0

| Aspek | v1.0 (Salah) | v2.0 (Standar) |
|-------|--------------|-----------------|
| Format soal | Pernyataan Likert 6-poin self-report | **Pilihan ganda A/B/C/D bergambar** |
| Dasar penilaian | Preferensi subjektif pengguna | **Jawaban benar/salah objektif** |
| Output | Tipe 4-huruf (mirip MBTI) | **Skor performa kognitif → estimasi IQ** |
| Validitas konstruk | Mengukur kepribadian, dilabeli "IQ" | Mengukur kemampuan kognitif nyata |

### 2.2 Domain Kemampuan Kognitif yang Diuji

Setiap soal bergambar termasuk dalam salah satu dari 4 domain kemampuan non-verbal/figural (dipilih karena format bergambar cocok untuk pengujian *fluid intelligence*, minim bias bahasa/budaya — mirip pendekatan Raven's Progressive Matrices & Cattell Culture Fair Test):

| Domain | Kode | Contoh Soal | Kemampuan yang Diukur |
|--------|------|--------------|------------------------|
| **Penalaran Matriks/Pola** | MTX | Lengkapi pola gambar 3×3 yang hilang | Penalaran induktif, deteksi pola |
| **Deret Logis Bergambar** | SEQ | Urutan bentuk/angka bergambar berikutnya | Penalaran deduktif, logika sekuensial |
| **Rotasi & Penalaran Spasial** | SPA | Bentuk mana hasil rotasi objek 3D ini | Visualisasi spasial, mental rotation |
| **Analogi Visual** | ANL | Gambar A:B seperti C:? | Penalaran analogis, abstraksi relasional |

### 2.3 Prinsip Asesmen

| Prinsip | Implementasi |
|---------|---------------|
| **Performance-based** | Skor dihitung dari jumlah jawaban benar, bukan preferensi |
| **Objektivitas** | Setiap soal punya tepat satu jawaban benar yang telah divalidasi |
| **Kalibrasi kesulitan** | Soal diurutkan dari mudah → sulit (item difficulty meningkat) |
| **Time-boxed** | Setiap soal dibatasi waktu untuk mengukur kecepatan pemrosesan kognitif |
| **Norm-referenced** | Skor akhir dibandingkan terhadap distribusi populasi (bukan skor absolut) |

### 2.4 Test Length & Duration

| Metric | Value |
|--------|-------|
| **Total Soal** | 20 (pilihan ganda A/B/C/D, bergambar) |
| **Waktu per Soal** | 2 menit (120 detik), hard limit dengan auto-advance |
| **Total Durasi Maksimum** | 40 menit (20 × 2 menit) |
| **Break Policy** | Tidak didukung (single session, timer berjalan terus) |
| **Retake Policy** | Cooldown 30 hari sebelum boleh mengulang (mencegah practice effect merusak validitas) |

---

## 3. PSYCHOMETRIC FOUNDATION

### 3.1 Teori yang Digunakan

Model ini mengikuti **Classical Test Theory (CTT)** dengan elemen kalibrasi kesulitan item, selaras dengan pendekatan tes kemampuan umum (*general cognitive ability, g-factor*, Spearman) yang diuji lewat penalaran figural non-verbal.

### 3.2 Struktur Item

| Elemen | Deskripsi |
|--------|-----------|
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar pada item tsb (dikalibrasi dari data uji coba) |
| **Item discrimination** | Seberapa baik item membedakan peserta berkemampuan tinggi vs rendah (point-biserial correlation) |
| **Bobot skor per item** | Item lebih sulit → bobot lebih tinggi saat dijawab benar |

### 3.3 Distribusi Kesulitan 20 Soal

| Level Kesulitan | Jumlah Soal | Bobot Skor per Soal | Target p-value |
|------------------|-------------|----------------------|-----------------|
| Mudah | 5 | 1.0 | 0.80–0.90 |
| Sedang | 8 | 1.5 | 0.50–0.79 |
| Sulit | 5 | 2.0 | 0.25–0.49 |
| Sangat Sulit | 2 | 2.5 | 0.10–0.24 |
| **Total** | **20** | **Max raw score = 30.5** | — |

### 3.4 Reliabilitas & Validitas (Wajib Sebelum Rilis Publik)

| Indikator | Target Minimum | Metode |
|-----------|------------------|--------|
| **Internal consistency** | Cronbach's α ≥ 0.80 | Dihitung dari data uji coba (min. 300 responden) |
| **Test-retest reliability** | r ≥ 0.75 | Subsample mengulang tes setelah 2–4 minggu |
| **Item discrimination** | Point-biserial ≥ 0.20 per item | Item di bawah ini dibuang/direvisi |
| **Construct validity** | Korelasi dengan tes kemampuan kognitif tervalidasi (mis. Raven's) | Studi validasi terpisah, disarankan sebelum klaim "IQ" digunakan secara publik |

> ⚠️ **Penting:** Tanpa proses kalibrasi dan validasi di atas, sistem **tidak boleh menampilkan angka "skor IQ"** kepada pengguna — lihat Bagian 7 untuk penjelasan dan alternatif sementara yang jujur (skor persentil relatif, bukan skala IQ resmi).

---

## 4. QUESTION STRUCTURE

### 4.1 Skema Metadata Soal

```go
type questionDef struct {
    ID            string   // e.g., "Q_MTX_001"
    Domain        string   // "MTX" | "SEQ" | "SPA" | "ANL"
    ImageURL      string   // gambar soal utama
    OptionImages  [4]string // gambar untuk opsi A, B, C, D
    CorrectOption string   // "A" | "B" | "C" | "D"
    Difficulty    string   // "easy" | "medium" | "hard" | "very_hard"
    Weight        float64  // 1.0 / 1.5 / 2.0 / 2.5 sesuai kesulitan
    PValue        float64  // dikalibrasi dari data uji coba, nullable saat awal
    Discrimination float64 // dikalibrasi, nullable saat awal
}
```

### 4.2 Question ID Naming Convention

```
Q_{DOMAIN}_{SEQUENCE}

Contoh:
  Q_MTX_001 — Penalaran Matriks, soal #1
  Q_SEQ_003 — Deret Logis, soal #3
  Q_SPA_002 — Rotasi Spasial, soal #2
  Q_ANL_004 — Analogi Visual, soal #4
```

### 4.3 Distribusi Domain (20 Soal)

| Domain | Jumlah Soal | Rentang Kesulitan |
|--------|-------------|---------------------|
| MTX (Matriks/Pola) | 6 | Mudah → Sangat Sulit |
| SEQ (Deret Logis) | 5 | Mudah → Sulit |
| SPA (Rotasi Spasial) | 5 | Sedang → Sangat Sulit |
| ANL (Analogi Visual) | 4 | Mudah → Sulit |

### 4.4 Format Respons

Setiap soal punya **4 opsi jawaban bergambar (A/B/C/D)**, tepat satu benar. Tidak ada skala setuju/tidak setuju — murni pilihan objektif.

### 4.5 Format Soal (UI)

```
┌─────────────────────────────────────────────┐
│   Soal 3 dari 20              ⏱ 01:47       │
│                                             │
│   [Gambar soal / pola yang harus dilengkapi] │
│                                             │
│   ┌───┐  ┌───┐  ┌───┐  ┌───┐                │
│   │ A │  │ B │  │ C │  │ D │                │
│   └───┘  └───┘  └───┘  └───┘                │
│  (masing-masing berisi gambar opsi jawaban)  │
│                                             │
│        [Pilih jawaban untuk lanjut →]       │
└─────────────────────────────────────────────┘
```

Catatan desain: tombol "Sebelumnya" **tidak tersedia** — dalam tes kemampuan kognitif standar, navigasi mundur untuk mengubah jawaban setelah waktu berjalan dapat merusak validitas pengukuran kecepatan pemrosesan.

---

## 5. TIMER RULES

### 5.1 Spesifikasi (Sesuai Permintaan)

| Rule | Value | Rationale |
|------|-------|-----------|
| **Waktu per soal** | 2 menit (120 detik), hard limit | Standar untuk soal penalaran figural tingkat menengah–sulit |
| **Total durasi tes** | Maks. 40 menit (20 × 2 menit) | Ceiling total jika semua soal terpakai penuh |
| **Auto-advance** | Ya — saat 120 detik habis, soal otomatis dianggap **tidak dijawab** (skor 0) dan lanjut ke soal berikutnya | Mencegah stalling, menjaga integritas timing |
| **Warning visual** | Countdown berubah warna kuning di 30 detik tersisa, merah di 10 detik | Memberi sinyal tanpa mengganggu konsentrasi |
| **Pause/Resume** | Tidak didukung | Sesuai standar administrasi tes terstandar (timed, tanpa jeda) |
| **Navigasi mundur** | Dinonaktifkan | Jawaban terkunci begitu dipilih atau waktu habis |

### 5.2 Implementasi Frontend (Alpine.js)

```javascript
Alpine.data('quizApp', () => ({
    currentQuestion: 0,
    timeRemaining: 120, // detik, reset setiap soal baru
    timerInterval: null,

    startQuestionTimer() {
        this.timeRemaining = 120;
        clearInterval(this.timerInterval);
        this.timerInterval = setInterval(() => {
            this.timeRemaining--;
            if (this.timeRemaining <= 0) {
                this.autoAdvance(); // submit sebagai unanswered
            }
        }, 1000);
    },

    selectAnswer(option) {
        const elapsedMs = (120 - this.timeRemaining) * 1000;
        this.answers[this.currentQuestion] = { option, elapsedMs };
        clearInterval(this.timerInterval);
        this.nextQuestion();
    },

    autoAdvance() {
        this.answers[this.currentQuestion] = { option: null, elapsedMs: 120000, timedOut: true };
        this.nextQuestion();
    }
}));
```

### 5.3 Response Time sebagai Data Tambahan

`time_taken_ms` tetap dicatat per soal (bukan untuk skor utama, tapi untuk analitik dan deteksi anomali — lihat Bagian 9).

---

## 6. SCORING ALGORITHM

### 6.1 Pipeline Skoring (Baru)

```
Jawaban User (A/B/C/D atau timeout)
        │
        ▼
Cocokkan dengan CorrectOption
        │
        ▼
Benar? → weighted_score += item.Weight
Salah/Timeout → weighted_score += 0
        │
        ▼
Raw Score = Σ weighted_score (maks 30.5, lihat 3.3)
        │
        ▼
Hitung skor per domain (MTX, SEQ, SPA, ANL)
        │
        ▼
Konversi ke estimasi IQ (Bagian 7) + persentil
```

### 6.2 Step-by-Step

**Step 1 — Cek jawaban tiap soal:**
```go
isCorrect := userAnswer.Option == question.CorrectOption
score := 0.0
if isCorrect {
    score = question.Weight
}
```

**Step 2 — Akumulasi per domain:**
```go
domainScores := map[string]float64{"MTX": 0, "SEQ": 0, "SPA": 0, "ANL": 0}
domainMax := map[string]float64{"MTX": 0, "SEQ": 0, "SPA": 0, "ANL": 0}

for _, q := range answeredQuestions {
    domainMax[q.Domain] += q.Weight
    if q.IsCorrect {
        domainScores[q.Domain] += q.Weight
    }
}
```

**Step 3 — Raw score total:**
```go
rawScore := domainScores["MTX"] + domainScores["SEQ"] + domainScores["SPA"] + domainScores["ANL"]
maxPossible := 30.5 // lihat tabel bobot 3.3
```

**Step 4 — Skor per domain sebagai persentase:**
```go
domainPercent := domainScores[d] / domainMax[d] * 100
```

### 6.3 Contoh Perhitungan

| Soal | Kesulitan | Bobot | Jawaban User | Benar? | Skor |
|------|-----------|-------|----------------|--------|------|
| Q_MTX_001 | Mudah | 1.0 | B | ✅ | 1.0 |
| Q_MTX_002 | Sedang | 1.5 | A | ❌ | 0 |
| Q_SEQ_001 | Sulit | 2.0 | C | ✅ | 2.0 |
| Q_SPA_003 | Sangat Sulit | 2.5 | (timeout) | ❌ | 0 |

Raw Score (contoh 4 soal) = 1.0 + 0 + 2.0 + 0 = **3.0 / 7.0 maksimum**

---

## 7. IQ SCORE CONVERSION

### 7.1 Status Saat Ini — Kejujuran ke Pengguna

Selama sistem **belum memiliki data normatif tervalidasi** (populasi referensi ≥1.000 peserta, lihat 3.4), aplikasi **tidak menampilkan angka "IQ" resmi**. Sebagai gantinya, tampilkan:

- **Skor mentah** (raw score / 30.5)
- **Persentil relatif** terhadap peserta lain di platform (bukan populasi umum tervalidasi)
- Disclaimer eksplisit: *"Estimasi ini bersifat indikatif dan belum divalidasi secara klinis. Bukan pengganti tes IQ terstandar oleh psikolog berlisensi."*

### 7.2 Model Konversi IQ (Setelah Data Normatif Tersedia)

```
IQ Score = 100 + (z_score × 15)

z_score = (raw_score_user − population_mean) / population_std_dev
```

Deviation IQ dengan mean=100, SD=15 ini adalah konvensi standar (Wechsler-style), bukan rasio usia mental seperti Stanford-Binet lama.

### 7.3 Tabel Konversi & Klasifikasi (Standar Wechsler, untuk referensi setelah norming)

| IQ Range | Klasifikasi |
|----------|-------------|
| 130+ | Very Superior |
| 120–129 | Superior |
| 110–119 | High Average |
| 90–109 | Average |
| 80–89 | Low Average |
| 70–79 | Borderline |
| < 70 | Extremely Low |

### 7.4 Syarat Normative Data

1. Minimum **1.000 sesi tes lengkap** untuk menghitung mean & SD awal.
2. Idealnya stratifikasi demografis (usia, pendidikan) — tanpa ini, skor IQ hasil hanya berlaku relatif terhadap basis pengguna platform, bukan populasi umum.
3. Re-norming tiap 6–12 bulan.
4. Kalibrasi ulang p-value & discrimination index item setiap penambahan data signifikan.

---

## 8. RESULT INTERPRETATION

### 8.1 Result Delivery Flow

```
User menyelesaikan 20 soal (atau waktu habis per soal)
            │
            ▼
    POST /submit-tes
            │
            ▼
    Server menghitung raw score & skor per domain
    Simpan jawaban + correctness
            │
            ▼
    Redirect ke /paywall/{id}
            │
            ▼
    User bayar → GET /hasil/{id}
            │
            ▼
    Tampilkan skor + interpretasi domain + narasi kekuatan/perkembangan
```

### 8.2 Result Components

| Section | Source | Description |
|---------|--------|--------------|
| **Skor Total** | Scoring Engine | Raw score / max, persentil |
| **Profil per Domain** | Scoring Engine | % benar di MTX, SEQ, SPA, ANL |
| **Klasifikasi** | Konversi (jika norm tersedia) / persentil (jika belum) | Kategori kemampuan |
| **Kekuatan Kognitif** | Narrative Generator | Domain dengan skor tertinggi |
| **Area Pengembangan** | Narrative Generator | Domain dengan skor terendah |
| **Kecepatan Pemrosesan** | Response time analytics | Rata-rata waktu jawab vs benar/salah |
| **Rekomendasi Latihan** | Narrative Generator | Saran domain untuk dilatih |

> **Dark Triad mapping pada v1.0 dihapus sepenuhnya.** Memetakan skor kemampuan kognitif ke narsisme/machiavellianism/psikopati tidak memiliki dasar psikometri apa pun dan berisiko menyesatkan serta merugikan pengguna secara psikologis. Jika ingin fitur kepribadian, itu harus jadi **tes terpisah** dengan instrumen kepribadian yang tervalidasi (mis. Big Five/HEXACO), bukan diturunkan dari skor tes kemampuan.

### 8.3 Visualisasi Skor Domain

```
Penalaran Matriks   ████████████████░░░░  80%  (Sangat Baik)
Deret Logis         ████████████░░░░░░░░  60%  (Baik)
Rotasi Spasial      ██████░░░░░░░░░░░░░░  30%  (Perlu Latihan)
Analogi Visual      ██████████████░░░░░░  70%  (Baik)
```

---

## 9. ANTI-CHEATING STRATEGY

### 9.1 Current Protections

| Strategy | Implementation | Status |
|----------|------------------|--------|
| Timer keras 2 menit/soal | Auto-advance saat habis | ✅ Wajib diimplementasikan (bukan opsional lagi) |
| Jawaban terkunci | Tidak bisa diubah setelah dipilih | ✅ Active |
| Server-side validation | Kunci jawaban benar tidak pernah dikirim ke client | ✅ **Kritis** — pastikan `CorrectOption` tidak bocor lewat API/DOM |
| Randomisasi urutan opsi | Posisi A/B/C/D diacak per sesi | ✅ Baru — mencegah pola hafalan |
| Paywall protection | Hasil terkunci hingga pembayaran dikonfirmasi | ✅ Active |

### 9.2 Planned Detections

| Detection | Method | Threshold | Action |
|-----------|--------|-----------|--------|
| **Speed-guessing** | time_taken_ms sangat rendah + banyak salah | < 3 detik/soal & akurasi < 25% | Flag "hasil tidak reliabel" |
| **Tab-switch detection** | visibility API | > 3 kali per tes | Flag di hasil (bukan blokir) |
| **Straight-pattern clicking** | Semua jawaban di opsi sama (mis. selalu "A") | ≥ 15/20 sama | Flag "pola respons tidak wajar" |
| **IP rate limiting** | Submission per IP | > 5/jam | 429 Too Many Requests |
| **Devtools/inspect tampering** | Deteksi perubahan DOM pada opsi jawaban | Any | Invalidate sesi |

### 9.3 Reliability Flag pada Hasil

```go
type ReliabilityFlag struct {
    IsReliable      bool
    Reasons         []string // "speed_guessing", "tab_switch_excessive", dll
    Recommendation  string   // "hasil_valid" | "disarankan_mengulang"
}
```

---

## 10. DATABASE MODEL

### 10.1 Skema yang Diperbarui

```sql
-- =============================================
-- Questions bank (BARU — dengan jawaban benar)
-- =============================================
CREATE TABLE questions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_code     VARCHAR(20) UNIQUE NOT NULL,   -- e.g. Q_MTX_001
    domain            VARCHAR(3) NOT NULL CHECK (domain IN ('MTX','SEQ','SPA','ANL')),
    difficulty        VARCHAR(10) NOT NULL CHECK (difficulty IN ('easy','medium','hard','very_hard')),
    weight            DECIMAL(3,1) NOT NULL,
    image_url         TEXT NOT NULL,                 -- gambar soal utama
    option_a_image    TEXT NOT NULL,
    option_b_image    TEXT NOT NULL,
    option_c_image    TEXT NOT NULL,
    option_d_image    TEXT NOT NULL,
    correct_option    CHAR(1) NOT NULL CHECK (correct_option IN ('A','B','C','D')),
    p_value           DECIMAL(4,3),                  -- dikalibrasi dari data uji coba
    discrimination    DECIMAL(4,3),                  -- dikalibrasi dari data uji coba
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- Session responses (per soal)
-- =============================================
CREATE TABLE session_responses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id        UUID NOT NULL REFERENCES test_sessions(id),
    question_id       UUID NOT NULL REFERENCES questions(id),
    selected_option   CHAR(1),                       -- NULL jika timeout
    is_correct        BOOLEAN NOT NULL,
    time_taken_ms     INTEGER NOT NULL,
    timed_out         BOOLEAN NOT NULL DEFAULT FALSE,
    answered_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, question_id)
);

-- =============================================
-- IQ Test results (BARU)
-- =============================================
CREATE TABLE iq_results (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id            UUID NOT NULL REFERENCES test_sessions(id) UNIQUE,
    raw_score             DECIMAL(5,2) NOT NULL,
    max_possible_score    DECIMAL(5,2) NOT NULL DEFAULT 30.5,
    mtx_score_pct         DECIMAL(5,1),
    seq_score_pct         DECIMAL(5,1),
    spa_score_pct         DECIMAL(5,1),
    anl_score_pct         DECIMAL(5,1),
    percentile            DECIMAL(5,1),               -- relatif terhadap basis pengguna
    estimated_iq          DECIMAL(5,1),                -- NULL sampai norm data tersedia
    avg_response_ms       INTEGER,
    is_reliable           BOOLEAN NOT NULL DEFAULT TRUE,
    reliability_flags     JSONB,
    calculated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Tabel `users`, `test_sessions`, `payments`, `admins` tidak berubah dari v1.0.

### 10.2 Perubahan Kunci vs v1.0

| v1.0 | v2.0 |
|------|------|
| `skor_lr`, `skor_na`, `skor_sa`, `skor_lv` (Likert) | `raw_score`, skor per domain (%) berbasis benar/salah |
| `iq_tipe VARCHAR(4)` (tipe 4-huruf) | `estimated_iq` (nullable sampai norm tersedia) + `percentile` |
| Tidak ada `correct_option` di soal | `correct_option` wajib, **tidak boleh dikirim ke client** |
| Tidak ada timer per soal di DB | `time_taken_ms`, `timed_out` per respons |

---

## 11. API FLOW

Alur endpoint tetap sama seperti v1.0 (`/quiz`, `/submit-tes`, `/paywall/:id`, `/hasil/:id`, dst), dengan satu perubahan kritis keamanan:

### 11.1 Kontrak Data Soal ke Client (Penting)

```json
// GET /api/questions — HANYA field ini yang boleh dikirim ke browser
{
  "id": "uuid",
  "question_code": "Q_MTX_003",
  "domain": "MTX",
  "image_url": "...",
  "options": {
    "A": "image_url_a",
    "B": "image_url_b",
    "C": "image_url_c",
    "D": "image_url_d"
  }
  // correct_option TIDAK PERNAH disertakan
}
```

Validasi jawaban **selalu dilakukan di server** saat `POST /submit-tes`, tidak pernah di client — mencegah manipulasi lewat devtools.

### 11.2 Route Table

Tidak berubah dari v1.0 (lihat Bagian 11.2 versi sebelumnya) — hanya payload request/response yang disesuaikan dengan skema pilihan ganda.

---

## 12. UI/UX FLOW

### 12.1 Quiz Page — Perubahan Utama

| Elemen | v1.0 | v2.0 |
|--------|------|------|
| Format soal | Teks pernyataan + skala 1–6 | **Gambar soal + 4 opsi jawaban bergambar** |
| Timer | Tidak ada | **Countdown 2 menit per soal, wajib tampil** |
| Navigasi mundur | Ada | **Dihapus** (integritas timing) |
| Progress bar | Ada | Dipertahankan, ditambah indikator waktu tersisa |

### 12.2 Alpine.js State (Diperbarui)

```javascript
Alpine.data('quizApp', () => ({
    step: 'identity',
    currentQuestion: 0,
    timeRemaining: 120,
    answers: {},            // { questionId: { option, elapsedMs, timedOut } }
    questions: [...],       // 20 soal, tanpa correct_option

    get progress() {
        return Object.keys(this.answers).length;
    },

    selectAnswer(option) { /* lihat Bagian 5.2 */ },
    autoAdvance() { /* lihat Bagian 5.2 */ },
    submitQuiz() { /* POST semua jawaban ke /submit-tes */ }
}));
```

### 12.3 Halaman lain (Landing, Paywall, Hasil, Admin)

Struktur tidak berubah signifikan dari v1.0 — hanya konten "Hasil" disesuaikan dengan skor domain baru (Bagian 8) dan disclaimer kejujuran skor IQ (Bagian 7.1).

---

## 13. FUTURE IMPROVEMENTS

| Improvement | Priority | Description |
|-------------|----------|-------------|
| **Kalibrasi item bank (p-value, discrimination)** | **Kritis** | Wajib sebelum klaim skor IQ ditampilkan ke publik |
| **Studi validasi konstruk** | **Kritis** | Bandingkan dengan tes tervalidasi (Raven's/CFIT) pada sampel independen |
| **Perluasan item bank 60–100 soal** | Tinggi | Untuk mendukung randomisasi & adaptive testing |
| **Computerized Adaptive Testing (CAT)** | Sedang | Pilih soal berikutnya berdasarkan performa real-time |
| **Stratifikasi norma demografis** | Sedang | Usia, pendidikan — meningkatkan akurasi konversi IQ |
| **Automated payment gateway** | Tinggi | Midtrans/Xendit |
| **Audit keamanan soal (anti-leak)** | Tinggi | Pastikan correct_option tidak pernah bocor ke client/log |

---

## 14. APPENDICES

### Appendix A — Glossary (Diperbarui)

| Term | Definition |
|------|------------|
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar suatu soal |
| **Item discrimination** | Seberapa baik soal membedakan peserta berkemampuan tinggi vs rendah |
| **Raw score** | Total skor tertimbang dari jawaban benar |
| **Deviation IQ** | Skor IQ dengan mean=100, SD=15, dihitung dari z-score terhadap populasi norma |
| **Percentile** | Posisi relatif skor dibanding peserta lain (bukan skala IQ absolut) |
| **Fluid intelligence** | Kemampuan menalar & memecahkan masalah baru tanpa bergantung pengetahuan yang dipelajari — inti dari soal bergambar non-verbal |
| **Reliability** | Konsistensi hasil tes (Cronbach's α, test-retest) |
| **Construct validity** | Sejauh mana tes benar-benar mengukur apa yang diklaim (kemampuan kognitif, bukan kepribadian) |

### Appendix B — Peringatan Etis & Legal

- Sistem **tidak boleh menampilkan angka "IQ"** yang diklaim setara tes klinis tanpa validasi psikometri (Bagian 3.4 & 7.4).
- Disarankan mencantumkan disclaimer di setiap halaman hasil: *"Tes ini untuk tujuan hiburan dan pengembangan diri, bukan diagnosis klinis. Untuk asesmen resmi, konsultasikan psikolog berlisensi."*
- Hindari klaim pemasaran seperti "tes IQ tervalidasi ilmiah" sebelum studi validasi (Bagian 13) selesai.

---

*End of IQTEST.md v2.0 — Revised Specification*
