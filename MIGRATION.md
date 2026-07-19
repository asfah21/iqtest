# MIGRATION.md — IQ Test Engine Implementation Plan

## Status: Draft | Last Updated: 2026-07-19

> Dokumen ini adalah panduan implementasi sistem **IQ Test Engine** — tes kemampuan kognitif pilihan ganda bergambar dengan satu jawaban benar objektif. Struktur Go/Gin/PostgreSQL/templ/Alpine.js digunakan sebagai basis.
>
> **Sumber kebenaran seluruhnya di IQTEST.md.** Jika terjadi perbedaan antara MIGRATION.md dan IQTEST.md, IQTEST.md yang berlaku.

---

## TABLE OF CONTENTS

1. [Implementation Overview](#1-implementation-overview)
2. [Phase 0 — Models & Types Definition](#2-phase-0--models--types-definition)
3. [Phase 1 — Database Schema Creation](#3-phase-1--database-schema-creation)
4. [Phase 2 — Repository Layer](#4-phase-2--repository-layer)
5. [Phase 3 — Scoring Engine](#5-phase-3--scoring-engine)
6. [Phase 4 — Question Bank & Frontend](#6-phase-4--question-bank--frontend)
7. [Phase 5 — Handler & Template](#7-phase-5--handler--template)
8. [Phase 6 — Narrative Engine](#8-phase-6--narrative-engine)
9. [Phase 7 — Anti-Cheating System](#9-phase-7--anti-cheating-system)
10. [Phase 8 — Automated Payment Gateway (Future)](#10-phase-8--automated-payment-gateway-future)
11. [Phase 9 — Admin Analytics & CSV Export (Future)](#11-phase-9--admin-analytics--csv-export-future)
12. [Dependency Graph](#12-dependency-graph)
13. [Appendices](#13-appendices)

---

## 1. IMPLEMENTATION OVERVIEW

### 1.1 Target State

Per **IQTEST.md**, sistem adalah **Cognitive Ability Assessment** dengan:

- 4 domain kognitif: MTX (Matriks/Pola), SEQ (Deret Logis), SPA (Rotasi Spasial), ANL (Analogi Visual)
- 20 soal pilihan ganda bergambar (A/B/C/D) dengan tepat satu jawaban benar
- Skoring berdasarkan jawaban benar/salah dengan bobot kesulitan (max raw = 30.5)
- Timer: 120 detik per soal, hard limit dengan auto-advance (IQTEST.md §5)
- Raw score + persentil (relatif terhadap pengguna platform)
- Estimasi IQ: **NULL sampai data normatif (1.000+ peserta) terkumpul** (IQTEST.md §7.4)
- Breakdown performa per domain (% benar per domain)
- Anti-cheating: timer keras, jawaban terkunci, server-side validation, randomisasi opsi, speed-guessing detection, tab-switch detection, pattern detection (IQTEST.md §9)
- Reliability flags pada hasil
- Disclaimer jujur: *"Estimasi ini bersifat indikatif dan belum divalidasi secara klinis"*
- Payment flow: freemium — tes gratis, hasil lengkap IDR 14.900 (IQTEST.md §1.1)
- Database: 7 tabel (`users`, `test_sessions`, `questions`, `session_responses`, `iq_results`, `payments`, `admins`) dengan DDL per IQTEST.md §10.2

### 1.2 Design Principles

| Principle | Description |
|-----------|-------------|
| **Buildable after every phase** | Go code MUST compile after each phase. No broken builds. |
| **Independent phases** | Each phase can be completed and tested in isolation. |
| **Fresh implementation** | Semua tabel dan kode dibuat baru, tidak ada migrasi dari sistem lama. |
| **No code modification outside scope** | Each phase touches only the files explicitly listed. |
| **CorrectOption is SERVER-ONLY** | Kunci jawaban tidak boleh bocor ke client via API/DOM (IQTEST.md §9.1). |

### 1.3 Phase Dependency Map

```
Phase 0 (Models) ──► Phase 1 (DB Schema) ──► Phase 2 (Repos) ──► Phase 3 (Scoring)
                                                                         │
                                                                         ▼
                                                                   Phase 4 (Frontend)
                                                                         │
                                                                         ▼
                                                                   Phase 5 (Handlers)
                                                                         │
                                                                         ▼
                                                                   Phase 6 (Narratives)
                                                                         │
                                                                         ▼
                                                                   Phase 7 (Anti-Cheat)

Phase 8, 9 are independent future enhancements
```

---

## 2. PHASE 0 — MODELS & TYPES DEFINITION

**Objective:** Mendefinisikan semua Go model/type struct per IQTEST.md. Ini adalah fondasi untuk semua fase lainnya.

**Source:** IQTEST.md §4.1 (QuestionDef), §8.2 (QuizResult), §9.3 (ReliabilityFlag), §10.2 (User), §10.2 DDL (semua tabel)

**Estimated effort:** 1–2 hours

---

### 2.1 Struct Definitions

**`models/question.go`** — Metadata soal pilihan ganda bergambar (per IQTEST.md §4.1):

```go
// QuestionDef — metadata satu soal pilihan ganda bergambar
type QuestionDef struct {
    ID             string           // UUID dari database
    QuestionCode   string           // e.g., "Q_MTX_001"
    Domain         string           // "MTX" | "SEQ" | "SPA" | "ANL"
    ImageURL       string           // gambar soal utama
    OptionImages   [4]string        // URL gambar opsi A, B, C, D
    CorrectOption  string           // "A" | "B" | "C" | "D" (HANYA di server!)
    Difficulty     string           // "easy" | "medium" | "hard" | "very_hard"
     Weight         float64          // 1.0 / 1.5 / 2.0 / 2.5 sesuai kesulitan
    PValue         *float64          // nullable — dikalibrasi dari data uji coba
    Discrimination *float64          // nullable — dikalibrasi dari data uji coba
}
```

**`models/session.go`** — Jawaban user per soal:

```go
// SessionResponse — jawaban user untuk satu soal
type SessionResponse struct {
    QuestionID     string
    QuestionCode   string
    SelectedOption *string // "A"/"B"/"C"/"D" atau nil jika timeout
    IsCorrect      bool
    TimeTakenMs    int
    TimedOut       bool
}
```

**`models/result.go`** — Struktur hasil dan skor:

```go
// DomainScore — skor untuk satu domain kognitif
type DomainScore struct {
    Domain      string  // "MTX" | "SEQ" | "SPA" | "ANL"
    RawScore    float64 // skor tertimbang
    MaxPossible float64 // skor maksimum domain ini
    Percentage  float64 // (raw/max) * 100
}

// IQTestResult — output akhir kalkulasi (per IQTEST.md §6)
type IQTestResult struct {
    RawScore         float64
    MaxPossible      float64
    DomainScores     map[string]DomainScore
    Percentile       *float64               // NULL sampai data normatif tersedia
    EstimatedIQ      *float64               // NULL sampai data normatif tersedia
    AvgResponseMs    int
    IsReliable       bool
    ReliabilityFlags []string
}

// ReliabilityFlag — indikator keandalan hasil tes (per IQTEST.md §9.3)
type ReliabilityFlag struct {
    IsReliable     bool
    Reasons        []string // "speed_guessing", "tab_switch_excessive", dll
    Recommendation string   // "hasil_valid" | "disarankan_mengulang"
}
```

**`models/user.go`** — Data user dan QuizResult untuk template:

```go
// User — data pengguna (per IQTEST.md §10.2 users table)
type User struct {
    ID        string
    Email     string
    Nama      string
    Phone     *string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// QuizResult — data yang dikirim ke template hasil (per IQTEST.md §8.2)
type QuizResult struct {
    ID    string
    Nama  string

    // Skor
    RawScore    float64
    MaxPossible float64
    Percentile  *float64
    EstimatedIQ *float64

    // Skor per domain
    DomainScores map[string]DomainScore

    // Waktu
    AvgResponseMs int

    // Reliabilitas
    IsReliable       bool
    ReliabilityFlags []string

    // Narrative fields
    ExecutiveSummary string
    Kekuatan         []string
    AreaPerhatian    []string
}

// PaywallData — data untuk rendering paywall page
type PaywallData struct {
    ID   string
    Nama string
}
```

### 2.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 0.1 | Create `QuestionDef` struct | `models/question.go` | Per IQTEST.md §4.1: ID, QuestionCode, Domain, ImageURL, OptionImages, CorrectOption, Difficulty, Weight, PValue (*float64), Discrimination (*float64). |
| 0.2 | Create `SessionResponse` struct | `models/session.go` | QuestionID, QuestionCode, SelectedOption (nullable), IsCorrect, TimeTakenMs, TimedOut. |
| 0.3 | Create `DomainScore` struct | `models/result.go` | Domain, RawScore, MaxPossible, Percentage. |
| 0.4 | Create `IQTestResult` struct | `models/result.go` | RawScore, MaxPossible, DomainScores, Percentile, EstimatedIQ, AvgResponseMs, IsReliable, ReliabilityFlags. |
| 0.5 | Create `ReliabilityFlag` struct | `models/result.go` | IsReliable, Reasons, Recommendation. |
| 0.6 | Create `User` struct | `models/user.go` | Per IQTEST.md §10.2 `users` table: ID, Email, Nama, Phone (*string), CreatedAt, UpdatedAt. **Score fields removed** — skor disimpan di `iq_results`. |
| 0.7 | Create `QuizResult` struct | `models/user.go` | ID, Nama, RawScore, MaxPossible, Percentile, EstimatedIQ, DomainScores, AvgResponseMs, IsReliable, ReliabilityFlags, ExecutiveSummary, Kekuatan, AreaPerhatian. |
| 0.8 | Create `PaywallData` struct | `models/user.go` | ID string, Nama string. |
| 0.9 | Create `Session` struct | `models/session.go` | Per IQTEST.md §10.2 `test_sessions` table: ID, UserID, SessionToken, StartedAt, CompletedAt, DeviceType, IPAddress, IsCompleted, Metadata (JSONB). |

### 2.3 Completion Criteria

- [ ] `go vet ./models/...` passes
- [ ] Semua struct di atas didefinisikan dengan field yang sesuai IQTEST.md
- [ ] `QuestionDef` memiliki `PValue` dan `Discrimination` sebagai `*float64`
- [ ] `User` tidak memiliki score fields — skor hanya di `IQTestResult` / `QuizResult`
- [ ] `SessionResponse` memiliki `TimeTakenMs` (int, bukan string/float)
- [ ] Tidak ada field personality/MBTI (SkorLR, IQTipe, CognitiveProfile, Dark Triad, dll.)

---

## 3. PHASE 1 — DATABASE SCHEMA CREATION

**Objective:** Membuat skema database baru per IQTEST.md §10.2. Tidak ada migrasi dari tabel lama — semua tabel baru dibuat dengan `CREATE TABLE`.

**Source:** IQTEST.md §10.1 (ERD), §10.2 (DDL)

**Estimated effort:** 1–2 hours

---

### 3.1 Tables (7 tabel per IQTEST.md §10.2)

| Table | Description | Key Columns |
|-------|-------------|-------------|
| `users` | Data peserta tes | `id`, `email` (UNIQUE), `nama`, `phone`, `created_at`, `updated_at` |
| `test_sessions` | Sesi tes per user | `id`, `user_id`, `session_token` (UNIQUE), `started_at`, `completed_at`, `device_type`, `ip_address`, `is_completed`, `metadata` |
| `questions` | Bank soal bergambar | `id`, `question_code` (UNIQUE), `domain`, `difficulty`, `weight`, `image_url`, `option_a_image`–`option_d_image`, `correct_option`, `p_value`, `discrimination`, `is_active` |
| `session_responses` | Jawaban per soal per sesi | `id`, `session_id`, `question_id`, `selected_option` (nullable), `is_correct`, `time_taken_ms`, `timed_out`, UNIQUE(session_id, question_id) |
| `iq_results` | Hasil kognitif lengkap | `id`, `session_id` (UNIQUE), `raw_score`, `max_possible_score`, `mtx_score_pct`–`anl_score_pct`, `percentile`, `estimated_iq`, `avg_response_ms`, `is_reliable`, `reliability_flags`, `calculated_at` |
| `payments` | Pembayaran hasil tes | `id`, `user_id`, `session_id`, `amount`, `currency`, `status`, `payment_method`, `paid_at`, `confirmed_by` (FK ke admins) |
| `admins` | Admin panel credentials | `id`, `username` (UNIQUE), `password_hash` |

### 3.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 1.1 | Create `users` table | `migrations/001_init_schema.sql` | DDL persis per IQTEST.md §10.2: `id UUID PK DEFAULT gen_random_uuid()`, `email VARCHAR(255) UNIQUE NOT NULL`, `nama VARCHAR(255) NOT NULL`, `phone VARCHAR(20)`, `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`, `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`. |
| 1.2 | Create `test_sessions` table | `migrations/001_init_schema.sql` | `id UUID PK`, `user_id UUID REFERENCES users(id)`, `session_token VARCHAR(64) UNIQUE NOT NULL`, `started_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`, `completed_at TIMESTAMPTZ`, `device_type VARCHAR(20)`, `ip_address INET`, `is_completed BOOLEAN NOT NULL DEFAULT FALSE`, `metadata JSONB`. |
| 1.3 | Create `questions` table | `migrations/001_init_schema.sql` | `id UUID PK`, `question_code VARCHAR(20) UNIQUE NOT NULL`, `domain VARCHAR(3)` CHECK IN ('MTX','SEQ','SPA','ANL'), `difficulty VARCHAR(10)` CHECK, `weight DECIMAL(3,1)`, `image_url TEXT`, `option_a_image`–`option_d_image TEXT`, `correct_option CHAR(1)` CHECK, `p_value DECIMAL(4,3)`, `discrimination DECIMAL(4,3)`, `is_active BOOLEAN DEFAULT TRUE`, `created_at TIMESTAMPTZ`. |
| 1.4 | Create `session_responses` table | `migrations/001_init_schema.sql` | `id UUID PK`, `session_id UUID NOT NULL REFERENCES test_sessions(id)`, `question_id UUID NOT NULL REFERENCES questions(id)`, `selected_option CHAR(1)` (nullable), `is_correct BOOLEAN NOT NULL`, `time_taken_ms INTEGER NOT NULL`, `timed_out BOOLEAN DEFAULT FALSE`, `answered_at TIMESTAMPTZ`. UNIQUE(session_id, question_id). |
| 1.5 | Create `iq_results` table | `migrations/001_init_schema.sql` | `id UUID PK`, `session_id UUID NOT NULL REFERENCES test_sessions(id) UNIQUE`, `raw_score DECIMAL(5,2) NOT NULL`, `max_possible_score DECIMAL(5,2) NOT NULL DEFAULT 30.5`, `mtx_score_pct DECIMAL(5,1)`, `seq_score_pct DECIMAL(5,1)`, `spa_score_pct DECIMAL(5,1)`, `anl_score_pct DECIMAL(5,1)`, `percentile DECIMAL(5,1)`, `estimated_iq DECIMAL(5,1)`, `avg_response_ms INTEGER`, `is_reliable BOOLEAN NOT NULL DEFAULT TRUE`, `reliability_flags JSONB`, `calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`. |
| 1.6 | Create `payments` table | `migrations/001_init_schema.sql` | `id UUID PK`, `user_id UUID NOT NULL REFERENCES users(id)`, `session_id UUID NOT NULL REFERENCES test_sessions(id)`, `amount DECIMAL(12,2) NOT NULL`, `currency VARCHAR(3) NOT NULL DEFAULT 'IDR'`, `status VARCHAR(20) NOT NULL DEFAULT 'PENDING'`, `payment_method VARCHAR(50)`, `paid_at TIMESTAMPTZ`, `confirmed_by UUID REFERENCES admins(id)`, `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`. |
| 1.7 | Create `admins` table | `migrations/001_init_schema.sql` | `id UUID PK`, `username VARCHAR(50) UNIQUE NOT NULL`, `password_hash VARCHAR(255) NOT NULL`, `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`. |
| 1.8 | Create indexes | `migrations/001_init_schema.sql` | Per IQTEST.md §10.2: `idx_sessions_user` (user_id WHERE NOT NULL), `idx_sessions_token`, `idx_responses_session`, `idx_results_session`, `idx_payments_user`, `idx_payments_status`, `idx_questions_domain` (WHERE is_active). |
| 1.9 | Migration runner | `database/db.go` | Ensure migration runner executes `001_init_schema.sql` on startup. |

### 3.3 DDL (Per IQTEST.md §10.2 — Source of Truth)

```sql
-- migrations/001_init_schema.sql

CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) UNIQUE NOT NULL,
    nama              VARCHAR(255) NOT NULL,
    phone             VARCHAR(20),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE questions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_code     VARCHAR(20) UNIQUE NOT NULL,
    domain            VARCHAR(3) NOT NULL CHECK (domain IN ('MTX','SEQ','SPA','ANL')),
    difficulty        VARCHAR(10) NOT NULL CHECK (difficulty IN ('easy','medium','hard','very_hard')),
    weight            DECIMAL(3,1) NOT NULL,
    image_url         TEXT NOT NULL,
    option_a_image    TEXT NOT NULL,
    option_b_image    TEXT NOT NULL,
    option_c_image    TEXT NOT NULL,
    option_d_image    TEXT NOT NULL,
    correct_option    CHAR(1) NOT NULL CHECK (correct_option IN ('A','B','C','D')),
    p_value           DECIMAL(4,3),
    discrimination    DECIMAL(4,3),
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE iq_results (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id            UUID NOT NULL REFERENCES test_sessions(id) UNIQUE,
    raw_score             DECIMAL(5,2) NOT NULL,
    max_possible_score    DECIMAL(5,2) NOT NULL DEFAULT 30.5,
    mtx_score_pct         DECIMAL(5,1),
    seq_score_pct         DECIMAL(5,1),
    spa_score_pct         DECIMAL(5,1),
    anl_score_pct         DECIMAL(5,1),
    percentile            DECIMAL(5,1),
    estimated_iq          DECIMAL(5,1),
    avg_response_ms       INTEGER,
    is_reliable           BOOLEAN NOT NULL DEFAULT TRUE,
    reliability_flags     JSONB,
    calculated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE admins (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username          VARCHAR(50) UNIQUE NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sessions_user ON test_sessions(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_sessions_token ON test_sessions(session_token);
CREATE INDEX idx_responses_session ON session_responses(session_id);
CREATE INDEX idx_results_session ON iq_results(session_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_questions_domain ON questions(domain) WHERE is_active = TRUE;
```

### 3.4 Completion Criteria

- [ ] Migration runs without errors on empty database
- [ ] All 7 tables created (`users`, `test_sessions`, `questions`, `session_responses`, `iq_results`, `payments`, `admins`)
- [ ] All CHECK constraints, FK references, and UNIQUE constraints per IQTEST.md §10.2
- [ ] All 7 indexes created
- [ ] `go build ./...` passes
- [ ] Models match schema columns

---

## 4. PHASE 2 — REPOSITORY LAYER

**Objective:** Membuat semua fungsi repository untuk CRUD tabel-tabel per IQTEST.md §10.2.

**Estimated effort:** 2–3 hours

---

### 4.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 2.1 | Create `InsertUser()` | `repositories/user.go` | INSERT INTO users (email, nama, phone) VALUES ($1, $2, $3) RETURNING id. |
| 2.2 | Create `GetUserByID()` | `repositories/user.go` | SELECT id, email, nama, phone, created_at, updated_at FROM users WHERE id = $1. |
| 2.3 | Create `GetUserByEmail()` | `repositories/user.go` | SELECT by email (untuk handle retake / user existing). |
| 2.4 | Create `GetAllUsers()` | `repositories/admin.go` | SELECT users.* dengan LEFT JOIN ke iq_results untuk raw_score, percentile, estimated_iq. |
| 2.5 | Create `InsertQuestion()` | `repositories/question.go` | INSERT into questions table. |
| 2.6 | Create `GetActiveQuestions()` | `repositories/question.go` | SELECT * FROM questions WHERE is_active = TRUE ORDER BY question_code. |
| 2.7 | Create `GetQuestionByCode()` | `repositories/question.go` | SELECT by question_code (untuk scoring). |
| 2.8 | Create `CreateSession()` | `repositories/session.go` | INSERT INTO test_sessions (user_id, session_token, device_type, ip_address) VALUES ($1, $2, $3, $4) RETURNING id. |
| 2.9 | Create `UpdateSessionCompleted()` | `repositories/session.go` | UPDATE test_sessions SET completed_at = NOW(), is_completed = TRUE WHERE id = $1. |
| 2.10 | Create `GetSessionByID()` | `repositories/session.go` | SELECT by id (untuk result lookup). |
| 2.11 | Create `GetSessionsByUserID()` | `repositories/session.go` | SELECT all sessions for a user (untuk riwayat). |
| 2.12 | Create `InsertResponse()` | `repositories/response.go` | INSERT INTO session_responses (session_id, question_id, selected_option, is_correct, time_taken_ms, timed_out). |
| 2.13 | Create `InsertResponsesBatch()` | `repositories/response.go` | Batch INSERT semua jawaban dalam satu transaksi. |
| 2.14 | Create `GetSessionResponses()` | `repositories/response.go` | SELECT all responses for a session JOIN questions for domain/weight/correct_option (for scoring re-calculation). |
| 2.15 | Create `InsertIQResult()` | `repositories/result.go` | INSERT INTO iq_results. |
| 2.16 | Create `GetIQResultBySession()` | `repositories/result.go` | SELECT from iq_results by session_id. |
| 2.17 | Create `GetAllRawScores()` | `repositories/result.go` | SELECT raw_score FROM iq_results (for percentile calculation). |
| 2.18 | Create `InsertPayment()` | `repositories/payment.go` | INSERT INTO payments (user_id, session_id, amount, currency, status). |
| 2.19 | Create `UpdatePaymentStatus()` | `repositories/payment.go` | UPDATE payments SET status = 'PAID', paid_at = NOW() WHERE session_id = $1. |
| 2.20 | Create `GetPaymentBySession()` | `repositories/payment.go` | SELECT by session_id (untuk verifikasi pembayaran). |

### 4.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] All repository functions use correct table/column names per IQTEST.md §10.2
- [ ] `GetAllUsers()` JOINs with iq_results for score display
- [ ] `InsertResponsesBatch()` uses transaction
- [ ] SQL queries match schema from Phase 1

---

## 5. PHASE 3 — SCORING ENGINE

**Objective:** Membuat scoring engine untuk cognitive ability test per IQTEST.md §6 (Scoring Algorithm) dan §7 (IQ Score Conversion).

**Source:** IQTEST.md §3.3 (Bobot skor), §6 (Scoring Algorithm), §7 (IQ Score Conversion)

**Estimated effort:** 3–4 hours

---

### 5.1 Scoring Pipeline (Per IQTEST.md §6.1)

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
Raw Score = Σ weighted_score (maks 30.5, per §3.3)
        │
        ▼
Hitung skor per domain (MTX, SEQ, SPA, ANL)
        │
        ▼
Konversi ke estimasi IQ (jika norm tersedia) + persentil
```

### 5.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 3.1 | Create question bank constant data | `services/quiz.go` | 20 `QuestionDef` entries per IQTEST.md §4.3 distribution (MTX:6, SEQ:5, SPA:5, ANL:4). Include `CorrectOption`, `Weight`, `Difficulty`. IDs: `Q_MTX_001`…`Q_ANL_004`. **`CorrectOption` hanya di server — tidak pernah dikirim ke client.** |
| 3.2 | Create `CalculateIQResult()` | `services/quiz.go` | Per IQTEST.md §6.2 Step 1–3: Take `[]SessionResponse` + `[]QuestionDef`. For each response: match question by QuestionCode, check correctness, accumulate weighted score per domain + total. Return `IQTestResult`. |
| 3.3 | Create domain scoring logic | `services/quiz.go` | Per IQTEST.md §6.2 Step 2 + Step 4: Group questions by domain, sum correct weight / sum max weight → percentage. Domain map: MTX, SEQ, SPA, ANL. |
| 3.4 | Create `CalculatePercentile()` | `services/quiz.go` | Per IQTEST.md §6.2 implicit: `CalculatePercentile(rawScore float64, allScores []float64) float64`. Count scores ≤ user score / total scores × 100. Return 0 if no data. |
| 3.5 | Create `EstimateIQ()` (stub) | `services/quiz.go` | Per IQTEST.md §7.2: `EstimateIQ(rawScore float64, mean float64, stdDev float64) *float64`. Return NULL if mean/stdDev not available (0 or nil). Formula: `100 + (z_score × 15)` dengan z_score = (raw − mean) / stdDev. |
| 3.6 | Create `ProcessQuizAnswers()` | `services/quiz.go` | Per IQTEST.md §11.1 flow: Take map of question_code → selected_option + timing. Build `[]SessionResponse`. Call `CalculateIQResult()`. Compute percentile from all existing iq_results. Call `EstimateIQ()` (will be NULL until norm data). Store results via repository. |
| 3.7 | Create `GetQuizResult()` | `services/quiz.go` | Per IQTEST.md §11.1: Load from `iq_results` table via repository. Join with session to get user nama. Return `QuizResult`. |

### 5.3 Scoring Reference (Per IQTEST.md §3.3 & Appendix B)

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

### 5.4 IQ Conversion Model (Per IQTEST.md §7.2 — Hanya setelah norm data tersedia)

```go
// EstimateIQ — Deviation IQ dengan mean=100, SD=15 (Wechsler-style)
// Hanya valid jika population_mean dan population_std_dev telah dikalibrasi
// dari ≥1.000 responden (per §7.4)
func EstimateIQ(rawScore float64, populationMean float64, populationStdDev float64) *float64 {
    if populationStdDev == 0 {
        return nil // belum ada data normatif
    }
    zScore := (rawScore - populationMean) / populationStdDev
    iq := 100 + (zScore * 15)
    return &iq
}
```

### 5.5 Completion Criteria

- [ ] `CalculateIQResult()` takes `[]SessionResponse` + `[]QuestionDef`, returns `IQTestResult`
- [ ] Domain scoring produces correct MTX/SEQ/SPA/ANL percentages (max per §3.3)
- [ ] Percentile calculation works (returns 0 when insufficient data)
- [ ] IQ estimation returns NULL when no normative data (per §7.1)
- [ ] `go build ./...` passes without errors

---

## 6. PHASE 4 — QUESTION BANK & FRONTEND

**Objective:** Membuat halaman kuis frontend dengan soal bergambar, timer, dan pilihan ganda per IQTEST.md §4–5.

**Source:** IQTEST.md §4 (Question Structure), §5 (Timer Rules), §12.2.2 (Quiz Page)

**Estimated effort:** 6–10 hours

---

### 6.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 4.1 | Create `quizApp` Alpine.js component | `templ/pages/quiz_page.templ` | Per IQTEST.md §5.2: `startQuestionTimer()`, `selectAnswer(option)`, `autoAdvance()`. State: `currentQuestion`, `timeRemaining`, `answers`. Timer: 120s countdown, auto-advance on 0. |
| 4.2 | Create quiz UI layout | `templ/pages/quiz_page.templ` | Per IQTEST.md §4.5 & §12.2.2: gambar soal di atas, 4 opsi bergambar A/B/C/D dalam grid. Progress bar "Soal N dari 20". Timer countdown. Tombol "Sebelumnya" tidak ada. |
| 4.3 | Add timer visual warnings | `templ/pages/quiz_page.templ` | Per IQTEST.md §5.1: Countdown berubah warna kuning di 30 detik tersisa, merah di 10 detik. |
| 4.4 | Fetch questions from API | `templ/pages/quiz_page.templ` | Load `questions[]` from `GET /api/questions` (tanpa `correctOption` — per IQTEST.md §11.2). |
| 4.5 | Submit answers | `templ/pages/quiz_page.templ` | `POST /submit-tes` sends array of `{ question_code, selected_option, time_taken_ms, timed_out }`. Per IQTEST.md §11.1. |
| 4.6 | Add tab-switch detection | `templ/pages/quiz_page.templ` | Per IQTEST.md §9.2: `visibilitychange` API listener, increment counter on each switch. Include `tab_switch_count` in submission. |
| 4.7 | Create quiz identity form | `templ/pages/quiz_page.templ` | Per IQTEST.md §12.2.2: Input nama + email (ditampilkan sebelum quiz dimulai, step 'identity'). |
| 4.8 | Create CSS styles | `assets/css/quiz.css` | Styles for image grid, timer warnings, progress indicators. Responsive per IQTEST.md §12.4 (mobile < 640px stacking). |
| 4.9 | Add accessibility attributes | `templ/pages/quiz_page.templ` | Per IQTEST.md §12.5: Alt text deskriptif pada gambar (tanpa bocorkan jawaban), `aria-live="polite"` pada timer, focus outlines, semantic HTML. |

### 6.2 Alpine.js Implementation (Per IQTEST.md §5.2 + §12.2.2)

```javascript
Alpine.data('quizApp', () => ({
    step: 'identity',       // 'identity' | 'quiz' | 'submitting' | 'done'
    currentQuestion: 0,
    timeRemaining: 120,
    timerInterval: null,
    answers: [],
    questions: [],
    tabSwitches: 0,
    nama: '',
    email: '',

    get progress() {
        return Object.keys(this.answers).length;
    },

    init() {
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) this.tabSwitches++;
        });
    },

    startQuiz() {
        fetch('/api/questions')
            .then(r => r.json())
            .then(qs => {
                this.questions = qs;
                this.step = 'quiz';
                this.startQuestionTimer();
            });
    },

    startQuestionTimer() {
        this.timeRemaining = 120;
        clearInterval(this.timerInterval);
        this.timerInterval = setInterval(() => {
            this.timeRemaining--;
            if (this.timeRemaining <= 0) this.autoAdvance();
        }, 1000);
    },

    selectAnswer(option) {
        const elapsedMs = (120 - this.timeRemaining) * 1000;
        this.answers[this.currentQuestion] = {
            question_code: this.questions[this.currentQuestion].question_code,
            selected_option: option,
            time_taken_ms: elapsedMs,
            timed_out: false,
        };
        clearInterval(this.timerInterval);
        this.nextQuestion();
    },

    autoAdvance() {
        this.answers[this.currentQuestion] = {
            question_code: this.questions[this.currentQuestion].question_code,
            selected_option: null,
            time_taken_ms: 120000,
            timed_out: true,
        };
        this.nextQuestion();
    },

    nextQuestion() {
        if (this.currentQuestion + 1 < this.questions.length) {
            this.currentQuestion++;
            this.startQuestionTimer();
        } else {
            this.submitQuiz();
        }
    },

    submitQuiz() {
        this.step = 'submitting';
        fetch('/submit-tes', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                nama: this.nama,
                email: this.email,
                answers: this.answers,
                tab_switch_count: this.tabSwitches,
            }),
        })
        .then(r => r.json())
        .then(data => {
            window.location.href = '/paywall/' + data.id;
        });
    }
}));
```

### 6.3 Completion Criteria

- [ ] `templ generate` passes without errors
- [ ] Quiz page loads 20 questions from API (no correct answers exposed — per §11.2)
- [ ] Each question shows an image + 4 image options in A/B/C/D grid (per §4.5)
- [ ] Timer counts down from 120s, changes color at 30s (kuning) and 10s (merah) — per §5.1
- [ ] Auto-advance works when timer hits 0 (records as timed_out, skor 0 — per §5.1)
- [ ] No backward navigation — selecting an answer advances immediately (per §5.1)
- [ ] Tab-switch detection logs events (per §9.2)
- [ ] Submission sends correct payload format (per §11.1)
- [ ] Accessibility: alt text, aria-live on timer, focus outlines, semantic HTML (per §12.5)

---

## 7. PHASE 5 — HANDLER & TEMPLATE

**Objective:** Membuat semua HTTP handlers dan template untuk flow IQ Test.

**Source:** IQTEST.md §11.1 (API Flow), §11.3 (Route Table), §11.4 (Response Types), §12 (UI/UX Flow)

**Estimated effort:** 4–6 hours

---

### 7.1 Route Table (Per IQTEST.md §11.3)

| Method | Path | Handler | Auth | Description |
|--------|------|---------|------|-------------|
| GET | `/` | ShowHome | None | Landing page (per §12.2.1) |
| GET | `/quiz` | ShowQuiz | None | Assessment page (per §12.2.2) |
| GET | `/api/questions` | GetQuestions | None | Ambil 20 soal (tanpa correct_option — per §11.2) |
| POST | `/submit-tes` | SubmitTest | None | Submit jawaban, hitung skor server-side |
| GET | `/paywall/:id` | ShowPaywall | None | Payment gate (per §12.2.3) |
| POST | `/konfirmasi-bayar/:id` | KonfirmasiBayar | None | Payment confirm manual |
| GET | `/hasil/:id` | ShowResult | None | View results (PAID only — per §12.2.4) |
| GET | `/tentang` | ShowTentang | None | About page |
| GET | `/admin/login` | ShowLogin | None | Admin login form |
| POST | `/admin/login` | LoginProcess | None | Admin login action |
| GET | `/admin/dashboard` | ShowDashboard | Admin cookie | Admin panel |
| GET | `/admin/user/:id` | ShowUserDetail | Admin cookie | User detail |
| GET | `/admin/logout` | LogoutProcess | Admin cookie | Logout |

### 7.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 5.1 | Create `GET /api/questions` | `handlers/quiz.go` | Return JSON array of 20 questions WITHOUT `correctOption`. Per IQTEST.md §11.2: hanya id, question_code, domain, image_url, options{A,B,C,D}. |
| 5.2 | Create `POST /submit-tes` | `handlers/quiz.go` | Accept `{ nama, email, answers: [{question_code, selected_option, time_taken_ms, timed_out}], tab_switch_count }`. Per §11.1 flow: create/update user, create session, call scoring engine, store results. Return `{ id: sessionID }` atau `{ id: userID }` (sesuai flow paywall). |
| 5.3 | Create `GET /paywall/:id` | `handlers/quiz.go` | Per §12.2.3: Halaman pembayaran IDR 14.900. "Halo, {nama}!" — ambil nama dari user/session. |
| 5.4 | Create `POST /konfirmasi-bayar/:id` | `handlers/quiz.go` | Per §11.1 + §12.2.3: Set payment status to PAID. Redirect ke /hasil/:id. |
| 5.5 | Create `GET /hasil/:id` | `handlers/quiz.go` | Per §12.2.4: Halaman hasil setelah bayar. Cek payment status. Map IQTestResult → QuizResult → template data. |
| 5.6 | Create landing page handler | `handlers/page.go` | Per §12.2.1: Serve index_page dengan tagline "Kenali Dirimu Lebih Dalam". |
| 5.7 | Create landing page template | `templ/pages/index_page.templ` | Per §12.2.1: Hero dengan DM Serif Display, CTA "Mulai Tes Gratis", trust pills, features 3-card, FAQ accordion. |
| 5.8 | Create quiz page template | `templ/pages/quiz_page.templ` | Per §12.2.2: Alpine.js quizApp component. Step 'identity' (nama + email) → step 'quiz' (soal + timer) → submit. |
| 5.9 | Create result page template | `templ/pages/hasil_page.templ` | Per §12.2.4 + §8.3: Raw score / max + persentil, domain breakdown progress bars (MTX/SEQ/SPA/ANL), estimated IQ (dengan disclaimer §7.1), reliability flags, executive summary, kekuatan/area perhatian. |
| 5.10 | Create paywall page template | `templ/pages/paywall_page.templ` | Per §12.2.3: "Hasil lengkap IQ Test — Rp14.900", instruksi transfer manual, tombol "Saya sudah bayar". |
| 5.11 | Create tentang page template | `templ/pages/tentang_page.templ` | Halaman statis tentang tes. |
| 5.12 | Add disclaimer to result page | `templ/pages/hasil_page.templ` | Per IQTEST.md §7.1 + Appendix E: "Estimasi ini bersifat indikatif dan belum divalidasi secara klinis. Bukan pengganti tes IQ terstandar oleh psikolog berlisensi." |
| 5.13 | Create types for template data | `templ/types/hasil_data.go` | Struct definitions matching QuizResult + template needs. |
| 5.14 | Create admin handlers | `handlers/admin.go` | Login, session cookie check, dashboard, user detail, logout. |
| 5.15 | Create admin templates | `templ/pages/dashboard_page.templ`, `templ/pages/user_detail_page.templ` | Dashboard: daftar user + scores. User detail: domain breakdown + responses. |
| 5.16 | Register all routes | `handlers/router.go` | Per §11.3 route table. |

### 7.3 Response Contract (Per IQTEST.md §11.4)

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

#### Question API Response (Per IQTEST.md §11.2 — WITHOUT correctOption)

```json
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
}
```

### 7.4 Completion Criteria

- [ ] `go build ./...` and `templ generate` pass without errors
- [ ] `GET /api/questions` returns 20 questions without correct answers (per §11.2)
- [ ] `POST /submit-tes` accepts new payload and returns ID
- [ ] Landing page shows IQ Test branding (per §12.2.1)
- [ ] Result page shows domain scores, raw score, percentile, disclaimer (per §12.2.4)
- [ ] Paywall page shows payment prompt (per §12.2.3)
- [ ] Dashboard shows user results with scores
- [ ] Route table matches §11.3 exactly

---

## 8. PHASE 6 — NARRATIVE ENGINE

**Objective:** Membuat narrative engine untuk menghasilkan laporan kognitif berdasarkan skor domain, per IQTEST.md §8 (Result Interpretation).

**Source:** IQTEST.md §8.2 (Result Components), §8.3 (Visualisasi Skor Domain)

**Estimated effort:** 2–3 hours

---

### 8.1 Narrative Output Structure (Per IQTEST.md §8.2)

| Section | Source | Description |
|---------|--------|--------------|
| **Skor Total** | Scoring Engine | Raw score / max, persentil |
| **Profil per Domain** | Scoring Engine | % benar di MTX, SEQ, SPA, ANL |
| **Klasifikasi** | Konversi / persentil | Kategori kemampuan |
| **Kekuatan Kognitif** | Narrative Generator | Domain dengan skor tertinggi |
| **Area Pengembangan** | Narrative Generator | Domain dengan skor terendah |
| **Kecepatan Pemrosesan** | Response time analytics | Rata-rata waktu jawab vs benar/salah |
| **Rekomendasi Latihan** | Narrative Generator | Saran domain untuk dilatih |

### 8.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 6.1 | Create `generateExecutiveSummary()` | `services/narasi.go` | Input: rawScore, maxScore, domainScores, percentile, estimatedIQ. Output: string deskripsi performa kognitif — kemampuan keseluruhan, keseimbangan domain, posisi relatif. |
| 6.2 | Create `generateKekuatan()` | `services/narasi.go` | Input: domainScores map[string]DomainScore. Output: []string — kekuatan berdasarkan domain dengan skor tertinggi. 3–5 bullet points. |
| 6.3 | Create `generateAreaPerhatian()` | `services/narasi.go` | Input: domainScores map[string]DomainScore. Output: []string — area pengembangan berdasarkan domain dengan skor terendah. 3–5 bullet points. |
| 6.4 | Create `generateRekomendasi()` | `services/narasi.go` | Input: domainScores. Output: []string — saran latihan spesifik per domain rendah. |
| 6.5 | Create `GenerateAllNarratives()` | `services/narasi.go` | Call all four functions. Return structured narrative output (ExecutiveSummary, Kekuatan, AreaPerhatian, Rekomendasi). |

### 8.3 Domain Score Visualization (Per IQTEST.md §8.3)

```
Penalaran Matriks   ████████████████░░░░  80%  (Sangat Baik)
Deret Logis         ████████████░░░░░░░░  60%  (Baik)
Rotasi Spasial      ██████░░░░░░░░░░░░░░  30%  (Perlu Latihan)
Analogi Visual      ██████████████░░░░░░  70%  (Baik)
```

### 8.4 Completion Criteria

- [ ] `generateExecutiveSummary()` produces cognitive performance summary
- [ ] `generateKekuatan()` derives from top-performing domains
- [ ] `generateAreaPerhatian()` derives from lowest-performing domains
- [ ] `generateRekomendasi()` produces actionable training suggestions
- [ ] No personality trait mapping (per §8.2 note: tidak menurunkan Dark Triad dari skor kognitif)
- [ ] `go build ./...` passes without errors

---

## 9. PHASE 7 — ANTI-CHEATING SYSTEM

**Objective:** Mengimplementasikan mekanisme anti-cheating dan reliability scoring per IQTEST.md §9.

**Source:** IQTEST.md §9.1 (Current Protections), §9.2 (Planned Detections), §9.3 (Reliability Flag)

**Estimated effort:** 3–5 days

---

### 9.1 Protections (Per IQTEST.md §9.1)

| Strategy | Implementation | Phase Source |
|----------|----------------|-------------|
| Timer keras 2 menit/soal | Auto-advance frontend (120s, hard limit) | Phase 4 |
| Jawaban terkunci | Tidak bisa diubah setelah dipilih | Phase 4 |
| Server-side validation | `CorrectOption` tidak pernah dikirim ke client (per §11.2) | Phase 5 + 7 |
| Randomisasi urutan opsi | Posisi A/B/C/D diacak per sesi | Phase 7 |
| Paywall protection | Hasil terkunci hingga pembayaran dikonfirmasi | Phase 5 |

### 9.2 Detections (Per IQTEST.md §9.2)

| Detection | Method | Threshold | Action |
|-----------|--------|-----------|--------|
| **Speed-guessing** | time_taken_ms sangat rendah + banyak salah | < 3 detik/soal & akurasi < 25% | Flag "hasil tidak reliabel" |
| **Tab-switch detection** | visibility API | > 3 kali per tes | Flag di hasil (bukan blokir) |
| **Straight-pattern clicking** | Semua jawaban di opsi sama (mis. selalu "A") | ≥ 15/20 sama | Flag "pola respons tidak wajar" |
| **IP rate limiting** | Submission per IP | > 5/jam | 429 Too Many Requests |
| **Devtools/inspect tampering** | Deteksi perubahan DOM pada opsi jawaban | Any | Invalidate sesi |

### 9.3 Reliability Flag Structure (Per IQTEST.md §9.3)

```go
type ReliabilityFlag struct {
    IsReliable      bool
    Reasons         []string // "speed_guessing", "tab_switch_excessive", dll
    Recommendation  string   // "hasil_valid" | "disarankan_mengulang"
}
```

### 9.4 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 7.1 | Implement server-side answer validation | `services/quiz.go` | Validate selected_option is A/B/C/D or null. Validate time_taken_ms is ≥ 0 and ≤ 120,000. |
| 7.2 | Implement speed-guessing detection | `services/quiz.go` | `DetectSpeedGuessing(responses []SessionResponse) bool`. True if >30% answers have time_taken_ms < 3000 and accuracy < 25%. |
| 7.3 | Implement straight-pattern detection | `services/quiz.go` | `DetectStraightPattern(responses []SessionResponse) bool`. True if ≥ 15/20 answers are same option. |
| 7.4 | Implement tab-switch tracking | `services/quiz.go` | Accept `tab_switch_count` in submission. Flag if > 3. |
| 7.5 | Implement `AssessReliability()` | `services/quiz.go` | Combine all detections. Return `ReliabilityFlag`. Store in iq_results.reliability_flags. |
| 7.6 | Implement IP rate limiting middleware | `middleware/ratelimit.go` | Per-IP: max 5 POST to `/submit-tes` per hour. Return 429. |
| 7.7 | Add random option shuffling | `services/quiz.go` + `handlers/quiz.go` | Shuffle option order per session. Store mapping in session metadata. Un-shuffle on submission before scoring. |
| 7.8 | Add devtools tamper detection (frontend) | `templ/pages/quiz_page.templ` | Detect DOM changes on option elements. Invalidate session on detection. |
| 7.9 | Register rate limiting | `handlers/router.go` | Apply rate limit middleware to `/submit-tes`. |

### 9.5 Completion Criteria

- [ ] Speed-guessing detection flags when >30% answers < 3s + accuracy < 25% (per §9.2)
- [ ] Straight-pattern detection flags when ≥ 15/20 same option (per §9.2)
- [ ] Tab-switch events tracked (> 3 = flagged)
- [ ] `POST /submit-tes` rate limited (max 5/hour per IP, return 429)
- [ ] Option order randomized per session (store mapping in session metadata)
- [ ] `CorrectOption` never reaches client (server-side validation only — per §11.2)
- [ ] Reliability flags stored with result in iq_results.reliability_flags (JSONB)
- [ ] `go build ./...` passes

---

## 10. PHASE 8 — AUTOMATED PAYMENT GATEWAY (Future)

**Objective:** Implementasikan integrasi payment gateway otomatis (Midtrans/Xendit) untuk menggantikan konfirmasi manual.

**Source:** IQTEST.md §13.3 (Payment & Monetization)

**Estimated effort:** 3–5 days

---

### 10.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 8.1 | Payment gateway integration | `services/payment.go`, `handlers/payment.go` | Integrate Midtrans/Xendit API untuk pembayaran instan. |
| 8.2 | Payment callback handler | `handlers/payment.go` | Webhook handler untuk update status pembayaran otomatis. |
| 8.3 | Payment status check endpoint | `handlers/payment.go` | Endpoint check status pembayaran (untuk polling frontend). |
| 8.4 | Multiple price tiers | `services/payment.go` | Basic (gratis) / Premium (detail, IDR 14.900) / Pro (konsultasi). |
| 8.5 | Discount codes | `services/payment.go`, `handlers/admin.go` | Kode promo dikelola admin. |

### 10.2 Completion Criteria

- [ ] Payment gateway integration works (Midtrans/Xendit)
- [ ] Payment status auto-updated via callback
- [ ] Result only accessible after payment confirmed
- [ ] Discount codes functional

---

## 11. PHASE 9 — ADMIN ANALYTICS & CSV EXPORT (Future)

**Objective:** Lengkapi admin panel dengan statistik kognitif, CSV export, dan analytics.

**Source:** IQTEST.md §12.2.5 (Admin Dashboard), §13.4 (User Experience)

**Estimated effort:** 2–4 days

---

### 11.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 9.1 | Admin dashboard statistics | `handlers/admin.go`, `templ/pages/dashboard_page.templ` | Average raw score, domain distribution (avg MTX/SEQ/SPA/ANL %), reliability rate, paid/unpaid counts, total revenue. |
| 9.2 | User detail view enhancement | `templ/pages/user_detail_page.templ` | Domain scores progress bars (per §8.3), reliability flags, response time stats per question. |
| 9.3 | Domain score distribution chart | `templ/pages/dashboard_page.templ` | Average % per domain (MTX, SEQ, SPA, ANL) — bisa bar chart CSS. |
| 9.4 | CSV export | `handlers/admin.go` | Export user results with v2.0 fields: nama, email, raw_score, persentil, domain %, estimated_iq, reliability. |
| 9.5 | Result history view | `templ/pages/user_detail_page.templ` | Track score changes across retake sessions (cooldown 30 hari per §2.4). |

### 11.2 Completion Criteria

- [ ] Dashboard shows cognitive statistics (averages, distributions)
- [ ] User detail shows domain breakdown with progress bars
- [ ] CSV export works with v2.0 fields
- [ ] Result history viewable

---

## 12. DEPENDENCY GRAPH

```
Phase 0 (Models) ──────► Phase 1 (DB Schema) ──► Phase 2 (Repos)
                                                          │
                                                          ▼
                                                    Phase 3 (Scoring Engine)
                                                          │
                                               ┌──────────┴──────────┐
                                               ▼                     ▼
                                         Phase 4 (Frontend)    Phase 5 (Handlers)
                                               │                     │
                                               └──────────┬──────────┘
                                                          ▼
                                                    Phase 6 (Narratives)
                                                          │
                                                          ▼
                                                    Phase 7 (Anti-Cheat)

                                         Phase 8, 9 (Independent futures)
```

### Parallel Execution

| Phase Set | Can run in parallel? | Reason |
|-----------|---------------------|--------|
| Phase 0 + Phase 1 | ❌ No | Models must be defined before schema |
| Phase 1 + Phase 2 | ❌ No | Repos depend on DB schema |
| Phase 3 + Phase 4 | ✅ Yes | Scoring and frontend independent after models/schema/repos |
| Phase 3 + Phase 5 | ❌ No | Handlers depend on scoring engine |
| Phase 4 + Phase 5 | ✅ Yes | Different files, can develop concurrently |
| Phase 6 + Phase 7 | ✅ Yes | Independent of each other |
| Phase 8, 9 | ✅ Yes | Independent future features |

---

## 13. APPENDICES

### Appendix A — File Inventory

| File | Phase | Action |
|------|-------|--------|
| `models/user.go` | 0 | Create |
| `models/question.go` | 0 | Create |
| `models/session.go` | 0 | Create |
| `models/result.go` | 0 | Create |
| `migrations/001_init_schema.sql` | 1 | Create (7 tabel + indexes per IQTEST.md §10.2) |
| `database/db.go` | 1 | Update (migration runner) |
| `repositories/user.go` | 2 | Create |
| `repositories/admin.go` | 2 | Create |
| `repositories/question.go` | 2 | Create |
| `repositories/session.go` | 2 | Create |
| `repositories/response.go` | 2 | Create |
| `repositories/result.go` | 2 | Create |
| `repositories/payment.go` | 2 | Create |
| `services/quiz.go` | 3 | Create |
| `services/narasi.go` | 6 | Create |
| `handlers/quiz.go` | 5 | Create |
| `handlers/admin.go` | 5 | Create |
| `handlers/page.go` | 5 | Create |
| `handlers/router.go` | 5 | Create |
| `middleware/ratelimit.go` | 7 | Create |
| `templ/types/hasil_data.go` | 5 | Create |
| `templ/types/dashboard_data.go` | 5 | Create |
| `templ/pages/index_page.templ` | 5 | Create |
| `templ/pages/quiz_page.templ` | 4 | Create |
| `templ/pages/hasil_page.templ` | 5 | Create |
| `templ/pages/paywall_page.templ` | 5 | Create |
| `templ/pages/tentang_page.templ` | 5 | Create |
| `templ/pages/dashboard_page.templ` | 5 | Create |
| `templ/pages/user_detail_page.templ` | 5 | Create |
| `assets/css/quiz.css` | 4 | Create |

### Appendix B — Glossary (Per IQTEST.md Appendix C)

| Term | Definition |
|------|------------|
| **Domain** | Kategori kemampuan kognitif yang diuji (MTX, SEQ, SPA, ANL) |
| **MTX** | Matrix/Pattern Reasoning — kemampuan penalaran induktif dan deteksi pola |
| **SEQ** | Logical Sequence — kemampuan penalaran deduktif dan logika sekuensial |
| **SPA** | Spatial Rotation — kemampuan visualisasi spasial dan mental rotation |
| **ANL** | Visual Analogy — kemampuan penalaran analogis dan abstraksi relasional |
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar suatu soal |
| **Item discrimination** | Seberapa baik soal membedakan peserta berkemampuan tinggi vs rendah |
| **Raw score** | Total skor tertimbang dari jawaban benar (maks. 30.5) |
| **Deviation IQ** | Skor IQ dengan mean=100, SD=15, dihitung dari z-score terhadap populasi norma (Wechsler-style) |
| **Percentile** | Posisi relatif skor dibanding peserta lain di platform |
| **Fluid intelligence** | Kemampuan menalar & memecahkan masalah baru tanpa bergantung pengetahuan yang dipelajari |
| **Reliability** | Konsistensi hasil tes (Cronbach's α, test-retest) — target α ≥ 0.80 |
| **Construct validity** | Sejauh mana tes benar-benar mengukur apa yang diklaim |
| **Norming / Normative data** | Data populasi referensi untuk mengonversi raw score ke skala IQ standar |

### Appendix C — Environment Variables (Per IQTEST.md Appendix D)

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | HTTP server port | `8080` |
| `ADMIN_USERNAME` | Admin login username | `admin` |
| `ADMIN_PASSWORD` | Admin login password | `admin360` |
| `QUESTION_TIMER_SECONDS` | Waktu per soal (detik) | `120` |
| `GIN_MODE` | Gin framework mode | `release` |

### Appendix D — Future Roadmap Reference (Per IQTEST.md §13.6)

```
Phase 1 (Q3 2026) — Fondasi Psikometri
├── Uji coba item bank ke 300+ responden
├── Kalibrasi p-value & discrimination
├── Randomisasi urutan soal
├── CDN untuk gambar soal

Phase 2 (Q4 2026) — Validasi & Monetisasi
├── Studi validasi konstruk (vs Raven's/CFIT)
├── Automated payment gateway (Phase 8)
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

### Appendix E — Peringatan Etis & Legal (Per IQTEST.md Appendix E)

- Sistem **tidak boleh menampilkan angka "IQ"** yang diklaim setara tes klinis tanpa validasi psikometri (per §3.4 & §7.4).
- Disarankan mencantumkan disclaimer di setiap halaman hasil: *"Tes ini untuk tujuan hiburan dan pengembangan diri, bukan diagnosis klinis. Untuk asesmen resmi, konsultasikan psikolog berlisensi."*
- Hindari klaim pemasaran seperti "tes IQ tervalidasi ilmiah" sebelum studi validasi (per §13.1) selesai.
- Sistem tidak memetakan skor kemampuan kognitif ke trait kepribadian (narsisme, machiavellianism, psikopati) — tidak ada dasar psikometri untuk ini (per §8.2).

---

*End of MIGRATION.md — Implementation Plan*
