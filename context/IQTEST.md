# IQTEST.md — IQ Test Engine
## Complete Technical & Functional Specification
### Version: 1.0 | Status: Draft | Last Updated: 2026-07-19

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

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Browser   │────▶│  Gin Router  │────▶│  Handlers   │
│  (Alpine.js)│     │  (HTTP/1.1)  │     │  (Go)       │
└─────────────┘     └──────────────┘     └──────┬──────┘
                                                │
                                        ┌───────▼──────┐
                                        │   Services   │
                                        │  (Business   │
                                        │   Logic)     │
                                        └───────┬──────┘
                                                │
                                        ┌───────▼──────┐
                                        │ Repositories │
                                        │  (Data Access)│
                                        └───────┬──────┘
                                                │
                                        ┌───────▼──────┐
                                        │  PostgreSQL  │
                                        │  (Database)  │
                                        └──────────────┘
```

### 1.3 Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Backend | Go 1.25 | HTTP server, business logic |
| HTTP Framework | Gin v1.12 | Routing, middleware, request handling |
| Templating | templ v0.3 | Type-safe HTML components |
| Database | PostgreSQL | Persistent storage for users, sessions, results |
| Frontend | Alpine.js | Client-side interactivity (quiz, timer, navigation) |
| Styling | Custom CSS | Design system based on DESIGN.md |
| Containerization | Docker | Development & deployment consistency |

---

## 2. ASSESSMENT METHODOLOGY

### 2.1 Prinsip Dasar

Tes ini mengukur **kemampuan kognitif nyata** (bukan preferensi atau kepribadian) melalui **20 soal pilihan ganda bergambar (A/B/C/D)** dengan tepat satu jawaban benar per soal. Pendekatan ini selaras dengan tes penalaran non-verbal standar seperti Raven's Progressive Matrices dan Cattell Culture Fair Test — format bergambar dipilih karena meminimalkan bias bahasa/budaya dan cocok untuk mengukur *fluid intelligence* (kemampuan menalar terhadap masalah baru).

### 2.2 Domain Kemampuan Kognitif yang Diuji

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

Model ini mengikuti **Classical Test Theory (CTT)** dengan kalibrasi kesulitan item, selaras dengan pendekatan tes kemampuan umum (*general cognitive ability, g-factor*, Spearman) yang diuji lewat penalaran figural non-verbal.

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
    ID             string    // e.g., "Q_MTX_001"
    Domain         string    // "MTX" | "SEQ" | "SPA" | "ANL"
    ImageURL       string    // gambar soal utama
    OptionImages   [4]string // gambar untuk opsi A, B, C, D
    CorrectOption  string    // "A" | "B" | "C" | "D"
    Difficulty     string    // "easy" | "medium" | "hard" | "very_hard"
    Weight         float64   // 1.0 / 1.5 / 2.0 / 2.5 sesuai kesulitan
    PValue         float64   // dikalibrasi dari data uji coba, nullable saat awal
    Discrimination float64   // dikalibrasi dari data uji coba, nullable saat awal
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

Setiap soal memiliki **4 opsi jawaban bergambar (A/B/C/D)**, tepat satu benar. Tidak ada skala setuju/tidak setuju — murni pilihan objektif.

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

### 5.1 Spesifikasi

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

### 6.1 Scoring Pipeline

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

### 6.2 Step-by-Step Algorithm

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

### 6.3 Scoring Example

| Soal | Kesulitan | Bobot | Jawaban User | Benar? | Skor |
|------|-----------|-------|----------------|--------|------|
| Q_MTX_001 | Mudah | 1.0 | B | ✅ | 1.0 |
| Q_MTX_002 | Sedang | 1.5 | A | ❌ | 0 |
| Q_SEQ_001 | Sulit | 2.0 | C | ✅ | 2.0 |
| Q_SPA_003 | Sangat Sulit | 2.5 | (timeout) | ❌ | 0 |

Raw Score (contoh 4 soal) = 1.0 + 0 + 2.0 + 0 = **3.0 / 7.0 maksimum**

---

## 7. IQ SCORE CONVERSION

### 7.1 Kebijakan Tampilan Skor

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

> Sistem tidak memetakan skor kemampuan kognitif ke trait kepribadian apa pun (mis. narsisme, machiavellianism, psikopati) karena tidak ada dasar psikometri untuk menurunkan trait kepribadian dari skor tes kemampuan. Jika di masa depan ingin menambahkan fitur kepribadian, itu harus menjadi **modul tes terpisah** dengan instrumen kepribadian yang tervalidasi (mis. Big Five/HEXACO), bukan diturunkan dari skor tes kemampuan kognitif ini.

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
| Timer keras 2 menit/soal | Auto-advance saat habis | ✅ Wajib diimplementasikan |
| Jawaban terkunci | Tidak bisa diubah setelah dipilih | ✅ Active |
| Server-side validation | Kunci jawaban benar tidak pernah dikirim ke client | ✅ **Kritis** — pastikan `CorrectOption` tidak bocor lewat API/DOM |
| Randomisasi urutan opsi | Posisi A/B/C/D diacak per sesi | ✅ Active — mencegah pola hafalan |
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

### 10.1 Entity Relationship Diagram

```
users
  │
  ├──< test_sessions
  │      │
  │      ├──< session_responses
  │      │
  │      └──> iq_results
  │
  └──< payments

admins ──> payments (confirmed_by)
questions ──< session_responses
```

### 10.2 DDL (PostgreSQL)

```sql
-- =============================================
-- Users table
-- =============================================
CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) UNIQUE NOT NULL,
    nama              VARCHAR(255) NOT NULL,
    phone             VARCHAR(20),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- Test sessions
-- =============================================
CREATE TABLE test_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID REFERENCES users(id),
    session_token     VARCHAR(64) UNIQUE NOT NULL,
    started_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at      TIMESTAMPTZ,
    device_type       VARCHAR(20),
    ip_address        INET,
    is_completed      BOOLEAN NOT NULL DEFAULT FALSE,
    metadata          JSONB
);

-- =============================================
-- Questions bank
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
-- IQ Test results
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

-- =============================================
-- Payments tracking
-- =============================================
CREATE TABLE payments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES users(id),
    session_id        UUID NOT NULL REFERENCES test_sessions(id),
    amount            DECIMAL(12,2) NOT NULL,
    currency          VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status            VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_method    VARCHAR(50),
    paid_at           TIMESTAMPTZ,
    confirmed_by      UUID REFERENCES admins(id),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- Admins
-- =============================================
CREATE TABLE admins (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username          VARCHAR(50) UNIQUE NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- Indexes
-- =============================================
CREATE INDEX idx_sessions_user ON test_sessions(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_sessions_token ON test_sessions(session_token);
CREATE INDEX idx_responses_session ON session_responses(session_id);
CREATE INDEX idx_results_session ON iq_results(session_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_questions_domain ON questions(domain) WHERE is_active = TRUE;
```

---

## 11. API FLOW

### 11.1 Complete Request/Response Flow

```
┌────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Browser │     │  Gin     │     │ Services │     │ Database │
└────┬───┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │               │                │                │
     │  GET /        │                │                │
     │──────────────▶│                │                │
     │◀──────────────│  IndexPage     │                │
     │               │                │                │
     │  GET /quiz    │                │                │
     │──────────────▶│                │                │
     │◀──────────────│  QuizPage (20 soal, tanpa correct_option) │
     │               │                │                │
     │  POST /submit-tes              │                │
     │  {email, nama, answers[{question_id, option, elapsed_ms}]} │
     │──────────────▶│                │                │
     │               │  ProcessQuizAnswers()           │
     │               │───────────────▶│                │
     │               │                │  ScoreAgainstCorrectOption │
     │               │                │  InsertUser()  │
     │               │                │───────────────▶│
     │               │                │◀───────────────│
     │               │◀───────────────│ {userID, err}  │
     │◀──────────────│  {id: userID}  │                │
     │               │                │                │
     │  GET /paywall/{id}            │                │
     │──────────────▶│                │                │
     │               │  GetPaywallData(id)             │
     │               │───────────────▶│                │
     │               │                │  GetUserName()  │
     │               │                │───────────────▶│
     │               │                │◀───────────────│
     │               │◀───────────────│ {nama, err}    │
     │◀──────────────│  PaywallPage   │                │
     │               │                │                │
     │  POST /konfirmasi-bayar/{id}  │                │
     │  {nama_pengirim}              │                │
     │──────────────▶│                │                │
     │               │  ConfirmPayment(id)             │
     │               │───────────────▶│                │
     │               │                │ UpdatePaymentStatus│
     │               │                │───────────────▶│
     │               │                │◀───────────────│
     │               │◀───────────────│ {err}          │
     │◀──────────────│  {success:true, id}             │
     │               │                │                │
     │  GET /hasil/{id}              │                │
     │──────────────▶│                │                │
     │               │  GetQuizResult(id)              │
     │               │───────────────▶│                │
     │               │                │ GetUserResult()│
     │               │                │───────────────▶│
     │               │                │◀───────────────│
     │               │                │ GenerateAllNarratives()│
     │               │◀───────────────│ {QuizResult}   │
     │◀──────────────│  HasilPage     │                │
```

### 11.2 Kontrak Data Soal ke Client (Penting)

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

### 11.3 Route Table

| Method | Path | Handler | Auth | Description |
|--------|------|---------|------|-------------|
| GET | `/` | ShowHome | None | Landing page |
| GET | `/quiz` | ShowQuiz | None | Assessment page |
| GET | `/api/questions` | GetQuestions | None | Ambil 20 soal (tanpa correct_option) |
| POST | `/submit-tes` | SubmitTest | None | Submit jawaban, hitung skor server-side |
| GET | `/paywall/:id` | ShowPaywall | None | Payment gate |
| POST | `/konfirmasi-bayar/:id` | KonfirmasiBayar | None | Payment confirm |
| GET | `/hasil/:id` | ShowResult | None | View results (PAID only) |
| GET | `/tentang` | ShowTentang | None | About page |
| GET | `/admin/login` | ShowLogin | None | Admin login form |
| POST | `/admin/login` | LoginProcess | None | Admin login action |
| GET | `/admin/dashboard` | ShowDashboard | Admin cookie | Admin panel |
| GET | `/admin/user/:id` | ShowUserDetail | Admin cookie | User detail |
| GET | `/admin/logout` | LogoutProcess | Admin cookie | Logout |

### 11.4 Response Types

#### Success Response (Submit Test)

```json
{
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

#### Error Response

```json
{
    "error": "Gagal menyimpan data tes: [reason]"
}
```

### 11.5 HTTP Status Codes

| Code | Condition |
|------|-----------|
| 200 | Success |
| 303 | Redirect (payment / paywall) |
| 400 | Bad request (invalid form data) |
| 404 | User/session not found |
| 429 | Rate limit exceeded |
| 500 | Internal server error |

---

## 12. UI/UX FLOW

### 12.1 User Journey Map

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Landing  │───▶│  Quiz    │───▶│ Paywall  │───▶│ Results  │
│  Page    │    │ (20 Soal │    │ (Payment)│    │  (PAID)  │
│          │    │  Timed)  │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
      │              │               │               │
      ▼              ▼               ▼               ▼
  IndexPage     QuizPage       PaywallPage      HasilPage
  (Static)      (Alpine.js)    (Static)         (Narratives)
```

### 12.2 Page Descriptions

#### 12.2.1 Landing Page (`GET /`)

| Element | Description |
|---------|-------------|
| **Hero** | Tagline: "Kenali Dirimu Lebih Dalam" dengan DM Serif Display |
| **CTA** | "Mulai Tes Gratis" — tombol utama ke `/quiz` |
| **Trust Pills** | "Anonim · Gratis · Hasil Instan" |
| **Features** | 3-card section menjelaskan asesmen |
| **FAQ** | Accordion FAQ |
| **Footer** | Links, copyright, brand info |

#### 12.2.2 Quiz Page (`GET /quiz`)

| Element | Description |
|---------|-------------|
| **Progress Bar** | Top of page, shows completion (X/20) |
| **Question Counter** | "Soal N dari 20" |
| **Timer Countdown** | 2:00 menghitung mundur per soal, warna berubah di 30s & 10s tersisa |
| **Gambar Soal** | Gambar pola/matriks/deret yang harus dilengkapi |
| **4 Opsi Bergambar** | Kartu A/B/C/D, masing-masing berisi gambar opsi jawaban |
| **Navigasi** | Hanya maju otomatis setelah memilih atau waktu habis (tidak ada tombol mundur) |
| **Form Fields** | Email dan Nama (ditampilkan sebelum quiz dimulai) |

**Alpine.js State Management:**

```javascript
Alpine.data('quizApp', () => ({
    step: 'identity',       // 'identity' | 'quiz' | 'submitting' | 'done'
    nama: '',
    email: '',
    currentQuestion: 0,
    timeRemaining: 120,
    timerInterval: null,
    answers: {},             // { questionId: { option, elapsedMs, timedOut } }
    questions: [...],        // 20 soal (tanpa correct_option)

    get progress() {
        return Object.keys(this.answers).length;
    },

    startQuiz() { this.step = 'quiz'; this.startQuestionTimer(); },
    startQuestionTimer() { /* lihat Bagian 5.2 */ },
    selectAnswer(option) { /* lihat Bagian 5.2 */ },
    autoAdvance() { /* lihat Bagian 5.2 */ },
    nextQuestion() {
        this.currentQuestion++;
        if (this.currentQuestion < 20) this.startQuestionTimer();
        else this.submitQuiz();
    },
    submitQuiz() { /* POST semua jawaban ke /submit-tes */ }
}));
```

#### 12.2.3 Paywall Page (`GET /paywall/:id`)

| Element | Description |
|---------|-------------|
| **Greeting** | "Halo, {nama}!" |
| **Value Prop** | Penjelasan hasil premium (skor per domain, narasi lengkap) |
| **Pricing** | IDR 14.900 (one-time payment) |
| **Payment Instructions** | Transfer manual ke rekening bank |
| **Confirmation Button** | "Saya sudah bayar" — POST ke `/konfirmasi-bayar/:id` |
| **Error State** | "belum_bayar" query param → pesan informasi |

#### 12.2.4 Results Page (`GET /hasil/:id`)

| Section | Content |
|---------|---------|
| **Header** | "Hasil Asesmen {nama}" |
| **Skor Total** | Raw score / 30.5 + persentil |
| **Profil per Domain** | Bar chart MTX/SEQ/SPA/ANL (% benar) |
| **Klasifikasi** | Kategori kemampuan berdasar persentil (atau estimasi IQ jika norm tersedia) |
| **Kekuatan Kognitif** | 3–5 bullet points |
| **Area Pengembangan** | 3–5 bullet points |
| **Kecepatan Pemrosesan** | Rata-rata waktu jawab per soal |
| **Disclaimer** | Catatan bahwa hasil bersifat indikatif, bukan diagnosis klinis |
| **Share/Print** | Action buttons (future) |

#### 12.2.5 Admin Dashboard (`GET /admin/dashboard`)

| Element | Description |
|---------|-------------|
| **Statistics** | Total users, paid/unpaid counts, total revenue |
| **User Table** | ID, Name, Email, Raw Score, Payment Status |
| **Question Bank Manager** | CRUD soal: upload gambar, set correct_option, difficulty, weight |
| **Search/Filter** | By name, email, status pembayaran |
| **User Detail** | Link ke `/admin/user/:id` — lihat jawaban per soal |
| **Logout** | Clear session cookie |

### 12.3 Design System Integration

| Token | Value |
|-------|-------|
| Surface | #FCF9F6 |
| Ink | #1a1917 |
| Primary | #0d7377 |
| Rounded (sm/md) | 8px / 12px |
| Body font | Inter, 1rem, 1.65 line-height |
| Display font | DM Serif Display |

### 12.4 Responsive Breakpoints

| Breakpoint | Width | Layout |
|------------|-------|--------|
| Mobile | < 640px | Single column, gambar soal & opsi ditumpuk 2×2 |
| Tablet | 640–1024px | 2-column grid pada opsi jawaban |
| Desktop | > 1024px | Full layout, opsi jawaban 4-kolom sejajar |

### 12.5 Accessibility Requirements

- WCAG AA minimum pada semua elemen teks
- Alt text deskriptif pada semua gambar soal & opsi (tanpa membocorkan jawaban benar)
- Focus-visible outlines pada elemen interaktif
- `prefers-reduced-motion` dihormati untuk animasi timer
- Semantic HTML structure (nav, main, section, footer)
- ARIA labels pada seluruh komponen interaktif, termasuk countdown timer (`aria-live="polite"`)

---

## 13. FUTURE IMPROVEMENTS

### 13.1 Psikometri & Item Bank

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Kalibrasi item bank (p-value, discrimination)** | Kritis | Wajib sebelum klaim skor IQ ditampilkan ke publik | 3 minggu |
| **Studi validasi konstruk** | Kritis | Bandingkan dengan tes tervalidasi (Raven's/CFIT) pada sampel independen | 4 minggu |
| **Perluasan item bank 60–100 soal** | Tinggi | Mendukung randomisasi & mencegah kebocoran soal | 3 minggu |
| **Computerized Adaptive Testing (CAT)** | Sedang | Pilih soal berikutnya berdasarkan performa real-time | 4 minggu |
| **Stratifikasi norma demografis** | Sedang | Usia, pendidikan — meningkatkan akurasi konversi IQ | 2 minggu |

### 13.2 Timer & Anti-Cheating

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Randomisasi urutan soal** | Tinggi | Urutan 20 soal diacak per sesi (domain tetap seimbang) | 3 hari |
| **Tab-switch detection** | Sedang | Log dan flag saat tab dialihkan selama tes | 2 hari |
| **IP rate limiting** | Sedang | Cegah submission berulang dari IP sama | 2 hari |
| **CAPTCHA integration** | Rendah | Google reCAPTCHA v3 saat submit | 2 hari |

### 13.3 Payment & Monetization

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Automated payment gateway** | Tinggi | Integrasi Midtrans/Xendit untuk pembayaran instan | 4 minggu |
| **Multiple price tiers** | Sedang | Basic (gratis) / Premium (detail) / Pro (konsultasi) | 2 minggu |
| **Discount codes** | Sedang | Kode promo dikelola admin | 1 minggu |

### 13.4 User Experience

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Email delivery hasil (PDF)** | Tinggi | Kirim PDF hasil ke email setelah pembayaran | 1 minggu |
| **Progress save/restore** | Sedang | Simpan progres parsial (dengan catatan: timer tetap berjalan sesuai kebijakan single-session) | 2 minggu |
| **Social sharing** | Sedang | Bagikan hasil ke medsos (link saja, tanpa skor mentah) | 3 hari |
| **Result history** | Rendah | Lacak perubahan skor antar sesi tes ulang | 2 minggu |

### 13.5 Technical Infrastructure

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **CDN untuk gambar soal** | Tinggi | Load gambar cepat, krusial untuk timer 2 menit/soal | 3 hari |
| **Redis caching** | Sedang | Cache question bank, kurangi beban DB | 1 minggu |
| **CI/CD pipeline** | Sedang | Automated testing & deployment | 1 minggu |
| **Load testing** | Sedang | Benchmark concurrent users (target: 10.000) | 1 minggu |

### 13.6 Recommended Roadmap

```
Phase 1 (Q3 2026) — Fondasi Psikometri
├── Uji coba item bank ke 300+ responden
├── Kalibrasi p-value & discrimination
├── Randomisasi urutan soal
├── CDN untuk gambar soal

Phase 2 (Q4 2026) — Validasi & Monetisasi
├── Studi validasi konstruk (vs Raven's/CFIT)
├── Automated payment gateway
├── Email delivery hasil (PDF)
├── Discount codes

Phase 3 (Q1 2027) — Skala & Normatif
├── Kumpulkan 1.000+ data untuk norming
├── Aktifkan estimated_iq resmi
├── Perluasan item bank 60–100 soal
├── Stratifikasi demografis

Phase 4 (Q2 2027) — Fitur Lanjutan
├── Computerized Adaptive Testing (CAT)
├── Result comparison tool
├── CI/CD pipeline
├── Load testing & optimization
```

---

## 14. APPENDICES

### Appendix A — File Map

| File | Purpose |
|------|---------|
| `main.go` | Application entry point, server initialization |
| `database/db.go` | PostgreSQL connection setup |
| `handlers/router.go` | Route registration |
| `handlers/page.go` | Static page handlers (home, quiz, about, error) |
| `handlers/quiz.go` | Quiz submission, paywall, result display handlers |
| `handlers/admin.go` | Admin login, dashboard, question bank management |
| `helpers/render.go` | Templ component rendering helper |
| `models/question.go` | Data models (Question, Answer, IQResult, ReliabilityFlag) |
| `repositories/user.go` | User data access (insert, query, update payment) |
| `repositories/question.go` | Question bank data access |
| `services/quiz.go` | Scoring algorithm, domain aggregation |
| `services/narasi.go` | Narrative generation engine |
| `templ/components/` | Reusable UI components (head, navbar, footer, timer) |
| `templ/layouts/` | Page layouts (public, quiz, auth, dashboard) |
| `templ/pages/` | Page templates (index, quiz, paywall, hasil, admin) |
| `assets/css/` | Stylesheets |
| `assets/js/` | JavaScript (Alpine.js modules, timer logic) |
| `assets/images/questions/` | Gambar soal & opsi jawaban |

### Appendix B — Scoring Reference Card

```
┌──────────────────────────────────────────────────────────────────┐
│                    SCORING REFERENCE CARD                        │
├────────────┬────────┬──────────┬────────┬─────────┬─────────────┤
│ Domain     │ Kode   │ Jumlah   │ Bobot  │ Max     │ Rentang     │
│            │        │ Soal     │ Range  │ Skor    │ Kesulitan   │
├────────────┼────────┼──────────┼────────┼─────────┼─────────────┤
│ Matriks    │ MTX    │ 6        │ 1.0–2.5│ ~9.5    │ Mudah→S.Sulit│
│ Deret Logis│ SEQ    │ 5        │ 1.0–2.0│ ~7.5    │ Mudah→Sulit │
│ Rotasi     │ SPA    │ 5        │ 1.5–2.5│ ~9.0    │ Sedang→S.Sulit│
│ Analogi    │ ANL    │ 4        │ 1.0–2.0│ ~5.5    │ Mudah→Sulit │
├────────────┴────────┴──────────┴────────┴─────────┴─────────────┤
│ Total Soal: 20  │  Waktu: 2 menit/soal  │  Max Raw Score: 30.5  │
└──────────────────────────────────────────────────────────────────┘
```

### Appendix C — Glossary

| Term | Definition |
|------|------------|
| **Domain** | Kategori kemampuan kognitif yang diuji (MTX, SEQ, SPA, ANL) |
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar suatu soal |
| **Item discrimination** | Seberapa baik soal membedakan peserta berkemampuan tinggi vs rendah |
| **Raw score** | Total skor tertimbang dari jawaban benar |
| **Deviation IQ** | Skor IQ dengan mean=100, SD=15, dihitung dari z-score terhadap populasi norma |
| **Percentile** | Posisi relatif skor dibanding peserta lain (bukan skala IQ absolut) |
| **Fluid intelligence** | Kemampuan menalar & memecahkan masalah baru tanpa bergantung pengetahuan yang dipelajari |
| **Reliability** | Konsistensi hasil tes (Cronbach's α, test-retest) |
| **Construct validity** | Sejauh mana tes benar-benar mengukur apa yang diklaim |
| **Norming / Normative data** | Data populasi referensi untuk mengonversi raw score ke skala IQ standar |

### Appendix D — Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | HTTP server port | `8080` |
| `ADMIN_USERNAME` | Admin login username | `admin` |
| `ADMIN_PASSWORD` | Admin login password | `admin360` |
| `QUESTION_TIMER_SECONDS` | Waktu per soal (detik) | `120` |
| `GIN_MODE` | Gin framework mode | `release` |

### Appendix E — Peringatan Etis & Legal

- Sistem **tidak boleh menampilkan angka "IQ"** yang diklaim setara tes klinis tanpa validasi psikometri (Bagian 3.4 & 7.4).
- Disarankan mencantumkan disclaimer di setiap halaman hasil: *"Tes ini untuk tujuan hiburan dan pengembangan diri, bukan diagnosis klinis. Untuk asesmen resmi, konsultasikan psikolog berlisensi."*
- Hindari klaim pemasaran seperti "tes IQ tervalidasi ilmiah" sebelum studi validasi (Bagian 13.1) selesai.

---

*End of IQTEST.md — Complete Specification Document*
