# MIGRATION.md — MBTI → IQ Test v2.0 Migration Plan

## Version: 2.0 | Status: Draft | Last Updated: 2026-07-19

> **Peringatan:** IQTEST.md v2.0 mengubah fundamental metodologi dari tes kepribadian self-report Likert menjadi **tes kemampuan kognitif pilihan ganda bergambar dengan satu jawaban benar objektif**. Ini bukan migrasi inkremental — ini adalah *rewrite arsitektural* yang memerlukan perubahan pada hampir setiap layer. Namun, struktur Go/Gin/PostgreSQL/templ/Alpine.js tetap sama.

---

## TABLE OF CONTENTS

1. [Migration Overview](#1-migration-overview)
2. [Phase 0 — Models & Types Rewrite (Required)](#2-phase-0--models--types-rewrite-required)
3. [Phase 1 — Database Schema Rewrite (Required)](#3-phase-1--database-schema-rewrite-required)
4. [Phase 2 — Repository Layer Rewrite (Required)](#4-phase-2--repository-layer-rewrite-required)
5. [Phase 3 — Scoring Engine Rewrite (Required)](#5-phase-3--scoring-engine-rewrite-required)
6. [Phase 4 — Question Bank & Frontend Rewrite (Required)](#6-phase-4--question-bank--frontend-rewrite-required)
7. [Phase 5 — Handler & Template Rewrite (Required)](#7-phase-5--handler--template-rewrite-required)
8. [Phase 6 — Narrative Engine Removal & Cognitive Profile (Required)](#8-phase-6--narrative-engine-removal--cognitive-profile-required)
9. [Phase 7 — Anti-Cheating System (Required)](#9-phase-7--anti-cheating-system-required)
10. [Phase 8 — Payment & Production Schema (Future Enhancement)](#10-phase-8--payment--production-schema-future-enhancement)
11. [Phase 9 — Admin Panel Update (Future Enhancement)](#11-phase-9--admin-panel-update-future-enhancement)
12. [Dependency Graph](#12-dependency-graph)
13. [Rollback Strategy](#13-rollback-strategy)

---

## 1. MIGRATION OVERVIEW

### 1.1 Current State

The application currently implements an **MBTI (Myers-Briggs Type Indicator)** personality assessment using:
- 4 MBTI dichotomies: E/I, S/N, T/F, J/P (currently aliased as LR/NA/SA/LV)
- 20 Likert-scale (1–6) self-report personality statements
- Scoring based on preference intensity
- 4-letter IQ types (e.g., LNSL) derived from dimension preferences
- Cognitive profile hierarchy (Dominant, Auxiliary, Complementary, Developing)
- Dark Triad narrative engine (Narcissism, Machiavellianism, Psychopathy)
- Narrative engine generating relationship insights from personality type
- Flat `users_test` table storing raw scores per dimension

### 1.2 Target State

Per **IQTEST.md v2.0**, the target is a **Cognitive Ability Assessment** using:
- 4 cognitive domains: MTX (Matriks/Pola), SEQ (Deret Logis), SPA (Rotasi Spasial), ANL (Analogi Visual)
- 20 multiple-choice visual questions (A/B/C/D) with exactly one correct answer each
- Scoring based on correct/incorrect answers with difficulty-weighted scores (max raw = 30.5)
- Timer: 120 seconds per question, hard limit auto-advance
- Raw score + percentile (relative to platform users) as primary output
- Estimated IQ: **NULL until normative data (1,000+ participants) is collected**
- No 4-letter type, no cognitive profile hierarchy, no Dark Triad
- Domain-specific performance breakdown (% correct per domain)
- Anti-cheating: speed-guessing detection, tab-switch detection, pattern detection
- Reliability flags on results
- Honest disclaimer: *"Estimasi ini bersifat indikatif dan belum divalidasi secara klinis"*

### 1.3 What Changes (Summary)

| Aspek | Sebelum (v1.0 / migrasi parsial) | Sesudah (v2.0) |
|-------|----------------------------------|-----------------|
| **Jenis soal** | 20 pernyataan Likert 1–6 | 20 gambar pilihan ganda A/B/C/D |
| **Dasar penilaian** | Intensitas preferensi | Jawaban benar/salah objektif |
| **Dimensi** | LR/NA/SA/LV (kepribadian) | MTX/SEQ/SPA/ANL (kognitif figural) |
| **Skor utama** | 4-letter type (e.g., LNSL) | Raw score / 30.5 + persentil |
| **IQ** | Diturunkan dari tipe 4-huruf | NULL sampai data normatif ≥1.000 |
| **Dark Triad** | Dipetakan dari skor dimensi | **Dihapus total** (tak berdasar psikometri) |
| **Timer** | Tidak ada | 2 menit/soal, auto-advance |
| **Navigasi mundur** | Ada | **Dihapus** |
| **Reliabilitas** | Tidak diukur | Flag kecepatan, pola, tab-switch |
| **Narasi** | Hubungan interpersonal | Analisis kekuatan/area pengembangan kognitif |

### 1.4 Design Principles

| Principle | Description |
|-----------|-------------|
| **Buildable after every phase** | The Go code MUST compile after each phase. No broken builds. |
| **Independent phases** | Each phase can be completed and tested in isolation. |
| **Backward compatible where possible** | Old database records remain readable until Phase 1 migration. |
| **No code modification outside scope** | Each phase touches only the files explicitly listed. |
| **Incremental deployment** | Phases can be deployed to production incrementally. |

### 1.5 Phase Dependency Map

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

## 2. PHASE 0 — MODELS & TYPES REWRITE (Required)

**Objective:** Rewrite all Go model/type definitions to match IQTEST.md v2.0 data structures. This phase touches ONLY `models/user.go` and related type files. No logic changes. The application will NOT compile after this phase (because services/handlers reference old fields) — this is expected and resolved in later phases.

**Rationale:** Establishing the correct data structures first prevents cascading changes in later phases. Models are the foundation everything else depends on.

**Estimated effort:** 2–3 hours

---

### 2.1 Key Structural Changes

| Old (v1.x) | New (v2.0) | Reason |
|-------------|------------|--------|
| `User.SkorLR`, `User.SkorNA`, `User.SkorSA`, `User.SkorLV` | `User.RawScore`, `User.MaxPossibleScore`, `User.Percentile`, `User.EstimatedIQ` | Skor berbasis benar/salah, bukan preferensi |
| `User.IQTipe` (4-letter type) | `User.EstimatedIQ` (nullable DECIMAL) | Tidak ada tipe 4-huruf di v2.0 |
| `DimensionScore` (struct) | `DomainScore` (struct) | Domain kognitif, bukan dimensi kepribadian |
| `CognitiveProfile` (Dominant/Auxiliary/Complementary/Developing) | `CognitiveProfile` → **dihapus** | Tidak ada hierarki kognitif di v2.0 |
| `IQTestResult` | `IQTestResult` (refactored) | Tanpa CognitiveProfile, tanpa Type 4-huruf |
| `QuizResult` | `QuizResult` (refactored) | Tanpa CognitiveProfile, tambah reliability flags |
| — (baru) | `QuestionDef` struct | Metadata soal bergambar dengan correct answer |
| — (baru) | `SessionResponse` struct | Jawaban per soal + timing |
| — (baru) | `ReliabilityFlag` struct | Deteksi kecurangan |
| `mapIQToDarkTriad()` | **Dihapus** | Dark Triad dihapus total |
| Dark Triad narrative fields | **Dihapus** | Tidak ada narasi Dark Triad |

### 2.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 0.1 | Rewrite `User` struct | `models/user.go` | Remove `SkorLR`, `SkorNA`, `SkorSA`, `SkorLV`, `IQTipe`. Add `RawScore DECIMAL(5,2)`, `MaxPossibleScore DECIMAL(5,2)`, `MTXScorePct`, `SEQScorePct`, `SPAScorePct`, `ANLScorePct` (all `*float64` nullable), `Percentile *float64`, `EstimatedIQ *float64`, `AvgResponseMs *int`, `IsReliable bool`, `ReliabilityFlags *string` (JSON). Keep `ID`, `Nama`, `Email`, `StatusPembayaran`. |
| 0.2 | Create `QuestionDef` struct | `models/user.go` (or new `models/question.go`) | Per IQTEST.md §4.1: `ID string`, `QuestionCode string`, `Domain string` (MTX/SEQ/SPA/ANL), `ImageURL string`, `OptionImages [4]string`, `CorrectOption string` (A/B/C/D), `Difficulty string`, `Weight float64`, `PValue *float64`, `Discrimination *float64`. |
| 0.3 | Create `SessionResponse` struct | `models/user.go` (or new file) | `QuestionID string`, `SelectedOption *string` (nil if timeout), `IsCorrect bool`, `TimeTakenMs int`, `TimedOut bool`. |
| 0.4 | Rewrite `DimensionScore` → `DomainScore` | `models/user.go` | Rename to `DomainScore`. Fields: `Domain string`, `RawScore float64`, `MaxPossible float64`, `Percentage float64`. Remove `PoleAScore`, `PoleBScore`, `Preference`, `SCI`, `Strength` — ini semua untuk model kepribadian bipolar, tidak berlaku untuk domain kognitif. |
| 0.5 | Rewrite `IQTestResult` struct | `models/user.go` | Remove `Type string`, `CognitiveProfile`. Replace with: `RawScore float64`, `MaxPossible float64`, `DomainScores map[string]DomainScore`, `Percentile *float64`, `EstimatedIQ *float64`, `AvgResponseMs int`, `IsReliable bool`, `ReliabilityFlags []string`. |
| 0.6 | Rewrite `QuizResult` struct | `models/user.go` | Remove `IQTipe`, `SkorLR/NA/SA/LV`, `CognitiveProfile`. Remove Dark Triad narrative fields: `RelationshipProfile`, `RelationshipInsight`, `CompatibilityNotes`, `ReflectionQuestions`. Add: `RawScore float64`, `MaxPossible float64`, `Percentile *float64`, `EstimatedIQ *float64`, `DomainScores map[string]DomainScore`, `ReliabilityFlags []string`, `IsReliable bool`. Keep `Nama`, `ExecutiveSummary` (repurposed), `Kekuatan []string`, `AreaPerhatian []string`. |
| 0.7 | Create `PaywallData` update | `models/user.go` | Keep as is (ID, Nama) — no change needed. |
| 0.8 | Create `ReliabilityFlag` struct | `models/models.go` (new) or `models/user.go` | Per IQTEST.md §9.3: `IsReliable bool`, `Reasons []string`, `Recommendation string`. |
| 0.9 | Create new type file(s) if needed | `models/question.go`, `models/session.go` | Split if models/user.go becomes too large. |
| 0.10 | Remove Dark Triad type references | `models/user.go` | Ensure no Narsisme/Machiavellian/Psikopati fields remain in any struct. |

### 2.3 New Struct Definitions (Reference)

```go
// QuestionDef — metadata satu soal pilihan ganda bergambar
type QuestionDef struct {
    ID            string           // UUID
    QuestionCode  string           // e.g., "Q_MTX_001"
    Domain        string           // "MTX" | "SEQ" | "SPA" | "ANL"
    ImageURL      string           // gambar soal utama
    OptionImages  [4]string        // URL gambar opsi A, B, C, D
    CorrectOption string           // "A" | "B" | "C" | "D"
    Difficulty    string           // "easy" | "medium" | "hard" | "very_hard"
    Weight        float64          // 1.0 / 1.5 / 2.0 / 2.5
    PValue        *float64         // nullable — dikalibrasi dari data uji coba
    Discrimination *float64        // nullable — dikalibrasi dari data uji coba
}

// DomainScore — skor untuk satu domain kognitif
type DomainScore struct {
    Domain      string  // "MTX" | "SEQ" | "SPA" | "ANL"
    RawScore    float64 // skor tertimbang
    MaxPossible float64 // skor maksimum domain ini
    Percentage  float64 // (raw/max) * 100
}

// SessionResponse — jawaban user untuk satu soal
type SessionResponse struct {
    QuestionID    string
    SelectedOption *string // "A"/"B"/"C"/"D" atau nil jika timeout
    IsCorrect     bool
    TimeTakenMs   int
    TimedOut      bool
}

// ReliabilityFlag — indikator keandalan hasil tes
type ReliabilityFlag struct {
    IsReliable     bool
    Reasons        []string // "speed_guessing", "tab_switch_excessive", dll
    Recommendation string   // "hasil_valid" | "disarankan_mengulang"
}

// IQTestResult — output akhir kalkulasi (versi v2.0)
type IQTestResult struct {
    RawScore        float64
    MaxPossible     float64
    DomainScores    map[string]DomainScore
    Percentile      *float64 // NULL sampai data normatif tersedia
    EstimatedIQ     *float64 // NULL sampai data normatif tersedia
    AvgResponseMs   int
    IsReliable      bool
    ReliabilityFlags []string
}

// QuizResult — data yang dikirim ke template hasil (versi v2.0)
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

    // Narrative fields (repurposed untuk kognitif)
    ExecutiveSummary string
    Kekuatan         []string
    AreaPerhatian    []string
}
```

### 2.4 Completion Criteria

- [ ] All old structs (`DimensionScore` with Pole/Preference/SCI/Strength, `CognitiveProfile`, old `IQTestResult` with Type, old `QuizResult`) are removed or fully refactored
- [ ] New structs (`QuestionDef`, `DomainScore`, `SessionResponse`, `ReliabilityFlag`) exist with correct fields per IQTEST.md
- [ ] No Dark Triad fields remain in any model
- [ ] `User` struct no longer has SkorLR/NA/SA/LV or IQTipe fields
- [ ] File compiles syntactically (`go vet ./models/...` passes)

---

## 3. PHASE 1 — DATABASE SCHEMA REWRITE (Required)

**Objective:** Rewrite the database schema to support image-based multiple-choice questions, per-question response tracking, and cognitive ability scoring per IQTEST.md §10. This is a *new migration* — the existing `users_test` table schema is incompatible with v2.0.

**Estimated effort:** 2–4 hours (+ data archival)

---

### 3.1 Key Structural Changes

| v1.x Table | v2.0 Change |
|-------------|-------------|
| `users_test` — flat table with `skor_lr`, `skor_na`, `skor_sa`, `skor_lv`, `iq_tipe` | Keep `users_test` but rewrite columns: remove old score columns, add `raw_score`, `max_possible_score`, `mtx_score_pct`, `seq_score_pct`, `spa_score_pct`, `anl_score_pct`, `percentile`, `estimated_iq`, `avg_response_ms`, `is_reliable`, `reliability_flags` |
| — (baru) | `questions` — bank soal bergambar per IQTEST.md §10.1 |
| — (baru) | `session_responses` — jawaban per soal + timing per IQTEST.md §10.1 |
| — (baru) | `iq_results` — hasil kognitif lengkap per IQTEST.md §10.1 |

### 3.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 1.1 | Create migration SQL for `questions` table | `migrations/001_v2_schema.sql` | Per IQTEST.md §10.1: `id UUID`, `question_code VARCHAR(20) UNIQUE`, `domain VARCHAR(3)`, `difficulty VARCHAR(10)`, `weight DECIMAL(3,1)`, `image_url TEXT`, `option_a_image TEXT`, …, `correct_option CHAR(1)`, `p_value DECIMAL(4,3)`, `discrimination DECIMAL(4,3)`, `is_active BOOLEAN`, `created_at TIMESTAMPTZ`. CHECK constraints on domain, difficulty, correct_option. |
| 1.2 | Create migration SQL for `session_responses` table | `migrations/001_v2_schema.sql` | Per IQTEST.md §10.1: `id UUID`, `session_id UUID REFERENCES test_sessions(id)`, `question_id UUID REFERENCES questions(id)`, `selected_option CHAR(1)` (nullable), `is_correct BOOLEAN`, `time_taken_ms INTEGER`, `timed_out BOOLEAN`, `answered_at TIMESTAMPTZ`. UNIQUE(session_id, question_id). |
| 1.3 | Create migration SQL for `iq_results` table | `migrations/001_v2_schema.sql` | Per IQTEST.md §10.1: `id UUID`, `session_id UUID REFERENCES test_sessions(id) UNIQUE`, `raw_score DECIMAL(5,2)`, `max_possible_score DECIMAL(5,2) DEFAULT 30.5`, `mtx_score_pct DECIMAL(5,1)`, `seq_score_pct DECIMAL(5,1)`, `spa_score_pct DECIMAL(5,1)`, `anl_score_pct DECIMAL(5,1)`, `percentile DECIMAL(5,1)`, `estimated_iq DECIMAL(5,1)`, `avg_response_ms INTEGER`, `is_reliable BOOLEAN`, `reliability_flags JSONB`, `calculated_at TIMESTAMPTZ`. |
| 1.4 | Rewrite `users_test` table columns | `migrations/001_v2_schema.sql` | Remove columns: `skor_lr`, `skor_na`, `skor_sa`, `skor_lv`, `iq_tipe`. Add columns: `raw_score DECIMAL(5,2)`, `max_possible_score DECIMAL(5,2) DEFAULT 30.5`, `mtx_score_pct DECIMAL(5,1)`, `seq_score_pct DECIMAL(5,1)`, `spa_score_pct DECIMAL(5,1)`, `anl_score_pct DECIMAL(5,1)`, `percentile DECIMAL(5,1)`, `estimated_iq DECIMAL(5,1)`, `avg_response_ms INTEGER`, `is_reliable BOOLEAN DEFAULT TRUE`, `reliability_flags JSONB`. |
| 1.5 | Archive old data | `migrations/001_archive_old_data.sql` | Before dropping old columns: `CREATE TABLE users_test_archive_v1 AS SELECT * FROM users_test;` untuk preservasi data historis. |
| 1.6 | Write rollback migration | `migrations/001_rollback.sql` | Reverse all schema changes (drop new tables, restore old columns from archive). |
| 1.7 | Update `database/db.go` (if needed) | `database/db.go` | Verify no hardcoded schema references. Add migration runner if not present. |

### 3.3 SQL Migration (Reference)

```sql
-- migrations/001_v2_schema.sql

-- 1. Archive old data before schema changes
CREATE TABLE IF NOT EXISTS users_test_archive_v1 AS SELECT * FROM users_test;

-- 2. Create questions table
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

-- 3. Create session_responses table
CREATE TABLE session_responses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id        UUID NOT NULL REFERENCES test_sessions(id),
    question_id       UUID NOT NULL REFERENCES questions(id),
    selected_option   CHAR(1),
    is_correct        BOOLEAN NOT NULL,
    time_taken_ms     INTEGER NOT NULL,
    timed_out         BOOLEAN NOT NULL DEFAULT FALSE,
    answered_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, question_id)
);

-- 4. Create iq_results table
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

-- 5. Rewrite users_test columns
ALTER TABLE users_test DROP COLUMN IF EXISTS skor_lr;
ALTER TABLE users_test DROP COLUMN IF EXISTS skor_na;
ALTER TABLE users_test DROP COLUMN IF EXISTS skor_sa;
ALTER TABLE users_test DROP COLUMN IF EXISTS skor_lv;
ALTER TABLE users_test DROP COLUMN IF EXISTS iq_tipe;

ALTER TABLE users_test ADD COLUMN raw_score DECIMAL(5,2);
ALTER TABLE users_test ADD COLUMN max_possible_score DECIMAL(5,2) NOT NULL DEFAULT 30.5;
ALTER TABLE users_test ADD COLUMN mtx_score_pct DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN seq_score_pct DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN spa_score_pct DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN anl_score_pct DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN percentile DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN estimated_iq DECIMAL(5,1);
ALTER TABLE users_test ADD COLUMN avg_response_ms INTEGER;
ALTER TABLE users_test ADD COLUMN is_reliable BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE users_test ADD COLUMN reliability_flags JSONB;
```

### 3.4 Completion Criteria

- [ ] Migration runs without errors
- [ ] `questions` table created with CHECK constraints
- [ ] `session_responses` table created with FK references
- [ ] `iq_results` table created with FK references
- [ ] `users_test` table has new v2.0 columns (old columns dropped)
- [ ] Archive table `users_test_archive_v1` preserves old data
- [ ] Rollback script verified
- [ ] `go build ./...` passes (models must match new schema)

---

## 4. PHASE 2 — REPOSITORY LAYER REWRITE (Required)

**Objective:** Rewrite all repository functions to work with the new v2.0 schema. This includes new CRUD for questions, responses, and results, plus updated user repository.

**Estimated effort:** 2–3 hours

---

### 4.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 2.1 | Rewrite `InsertUser()` | `repositories/user.go` | Insert into `users_test` with new columns: `raw_score`, `max_possible_score`, `mtx_score_pct`, `seq_score_pct`, `spa_score_pct`, `anl_score_pct`, `percentile`, `estimated_iq`, `avg_response_ms`, `is_reliable`, `reliability_flags`. Remove old column inserts (`skor_lr`, `iq_tipe`, dll). |
| 2.2 | Rewrite `GetUserResult()` | `repositories/user.go` | SELECT new columns, scan into updated `QuizResult` struct. |
| 2.3 | Rewrite `GetAllUsers()` | `repositories/admin.go` | SELECT new columns. |
| 2.4 | Rewrite `GetUserByID()` | `repositories/admin.go` | SELECT new columns. |
| 2.5 | Create `InsertQuestion()` | `repositories/question.go` (new) | Insert a question into the `questions` table. |
| 2.6 | Create `GetActiveQuestions()` | `repositories/question.go` | SELECT * FROM questions WHERE is_active = TRUE ORDER BY question_code. |
| 2.7 | Create `InsertResponse()` | `repositories/response.go` (new) | Insert a session response into `session_responses`. |
| 2.8 | Create `GetSessionResponses()` | `repositories/response.go` | Get all responses for a session. |
| 2.9 | Create `InsertIQResult()` | `repositories/result.go` (new) | Insert into `iq_results`. |
| 2.10 | Create `GetIQResultBySession()` | `repositories/result.go` | SELECT from `iq_results` by session_id. |

### 4.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] All old SkorLR/NA/SA/LV references removed from repository files
- [ ] New repository files exist and are functional
- [ ] SQL queries match new v2.0 schema columns

---

## 5. PHASE 3 — SCORING ENGINE REWRITE (Required)

**Objective:** Replace the personality-based scoring engine with the cognitive ability scoring engine per IQTEST.md §6. This is the core behavioral change.

**Estimated effort:** 4–6 hours

---

### 5.1 Key Changes

| Aspek | v1.x (Old) | v2.0 (New) |
|-------|------------|-------------|
| **Input** | 20 Likert values (1–6) | 20 multiple-choice answers (A/B/C/D or timeout) |
| **Kebenaran** | Tidak ada — intensitas preferensi | Dicocokkan dengan `CorrectOption` |
| **Bobot** | 2.0 / 1.5 per pole | 1.0 / 1.5 / 2.0 / 2.5 per difficulty |
| **Skor dimensi** | LR/NA/SA/LV — akumulasi preferensi | MTX/SEQ/SPA/ANL — akumulasi jawaban benar × weight |
| **Output** | 4-letter type + CognitiveProfile | Raw score + DomainScores + percentile |
| **IQ** | Diturunkan dari tipe | NULL sampai norm data tersedia |
| **Dark Triad** | Narcissism/Machiavellianism/Psychopathy | **Dihapus total** |

### 5.2 Scoring Pipeline (per IQTEST.md §6.1)

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
Raw Score = Σ weighted_score (maks 30.5)
        │
        ▼
Hitung skor per domain (MTX, SEQ, SPA, ANL)
        │
        ▼
Konversi ke estimasi IQ (jika norm tersedia) + persentil
```

### 5.3 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 3.1 | Create question bank constant data | `services/quiz.go` | 20 `QuestionDef` entries per IQTEST.md §4.3 distribution (MTX:6, SEQ:5, SPA:5, ANL:4). Include `CorrectOption`, `Weight`, `Difficulty`. IDs: `Q_MTX_001`…`Q_ANL_004`. **`CorrectOption` hanya di server — tidak pernah dikirim ke client.** |
| 3.2 | Rewrite `CalculateIQResult()` | `services/quiz.go` | Take `[]SessionResponse` + `[]QuestionDef` as input. For each response: match question, check correctness, accumulate weighted score. Return `IQTestResult`. |
| 3.3 | Implement domain scoring | `services/quiz.go` | After raw score calculation, compute per-domain scores: group questions by domain, sum correct weight / sum max weight → percentage. |
| 3.4 | Implement percentile calculation | `services/quiz.go` | `CalculatePercentile(rawScore float64, allScores []float64) float64`. Count scores ≤ user score / total scores × 100. Return NULL if < 50 results in DB. |
| 3.5 | Implement IQ estimation (stub) | `services/quiz.go` | `EstimateIQ(rawScore float64, mean float64, stdDev float64) *float64`. Return NULL if mean/stdDev not available. Formula: `100 + (z_score × 15)`. |
| 3.6 | Remove `DeriveCognitiveProfile()` | `services/quiz.go` | Delete this function entirely — tidak ada hierarki kognitif di v2.0. |
| 3.7 | Remove `mapIQToDarkTriad()` | `services/quiz.go` | Delete this function entirely — Dark Triad dihapus total. |
| 3.8 | Remove `axisOpposites` / `axisOpposite()` | `services/quiz.go` | Already removed in Phase 0 (ensure no remnants). |
| 3.9 | Rewrite `ProcessQuizAnswers()` | `services/quiz.go` | Take map of question_code → selected_option. Build `[]SessionResponse`. Call `CalculateIQResult()`. Store results via repository. |
| 3.10 | Rewrite `GetQuizResult()` | `services/quiz.go` | Load from `iq_results` table via repository. No CognitiveProfile derivation. Compute percentile if enough data exists. |
| 3.11 | Add response time analytics | `services/quiz.go` | Calculate average response time from `SessionResponse.TimeTakenMs`. Include in `IQTestResult.AvgResponseMs`. |

### 5.4 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] `services/quiz.go` contains no CognitiveProfile, Dark Triad, or personality dimension logic
- [ ] `CalculateIQResult()` takes `[]SessionResponse` + `[]QuestionDef`, returns `IQTestResult`
- [ ] Domain scoring produces correct MTX/SEQ/SPA/ANL percentages
- [ ] Percentile calculation works (returns NULL when insufficient data)
- [ ] IQ estimation returns NULL when no normative data
- [ ] All old personality-scoring functions removed

---

## 6. PHASE 4 — QUESTION BANK & FRONTEND REWRITE (Required)

**Objective:** Rewrite the quiz frontend from Likert-scale personality statements to image-based multiple-choice cognitive ability questions with timer. This is the most visible user-facing change.

**Estimated effort:** 6–10 hours

---

### 6.1 Key Changes

| Elemen | v1.x | v2.0 |
|--------|------|------|
| Format soal | Teks pernyataan + skala 1–6 | Gambar soal + 4 opsi jawaban bergambar (A/B/C/D) |
| Timer | Tidak ada | Countdown 2 menit per soal, wajib tampil |
| Navigasi mundur | Ada | Dihapus (jawaban terkunci setelah dipilih) |
| Progress | Bar + nomor soal | Bar + nomor soal + timer countdown |
| State management | `answers[]` berisi nilai 1–6 per soal | `answers[]` berisi `{ option, elapsedMs, timedOut }` |

### 6.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 4.1 | Rewrite `quizApp` Alpine.js data | `templ/pages/quiz_page.templ` | Per IQTEST.md §5.2: implement `startQuestionTimer()`, `selectAnswer(option)`, `autoAdvance()`. State: `currentQuestion`, `timeRemaining`, `answers` (object with option/elapsedMs/timedOut). Timer: 120s countdown, auto-advance on 0. |
| 4.2 | Design new quiz UI layout | `templ/pages/quiz_page.templ` | Per IQTEST.md §4.5: soalnya gambar, opsi A/B/C/D bergambar dalam grid 2×2. Tombol "Sebelumnya" dihapus. Progress bar + timer countdown di atas. |
| 4.3 | Add timer visual warnings | `templ/pages/quiz_page.templ` | Countdown kuning di 30 detik, merah di 10 detik. Per IQTEST.md §5.1. |
| 4.4 | Fetch questions from API | `templ/pages/quiz_page.templ` | Load `questions[]` array from `GET /api/questions` endpoint (server provides 20 questions tanpa `correctOption`). |
| 4.5 | Submit answers as `{ questionId, option, elapsedMs }` | `templ/pages/quiz_page.templ` | `POST /submit-tes` sends array of `{ question_code, selected_option, time_taken_ms }` — not Likert values. |
| 4.6 | Add anti-cheat frontend measures | `templ/pages/quiz_page.templ` | Tab-switch detection via `visibilitychange` API (log timestamps). Disable right-click/copy-paste during quiz. Per IQTEST.md §9.2. |
| 4.7 | Update quiz identity form | `templ/pages/quiz_page.templ` | "Hasil asesmen IQ Test" (not MBTI/personality). |
| 4.8 | Update CSS for new quiz layout | `assets/css/` | Styles for image-based grid layout, timer warnings, progress indicators. |

### 6.3 Key Alpine.js Implementation (Reference)

```javascript
Alpine.data('quizApp', () => ({
    currentQuestion: 0,
    timeRemaining: 120,
    timerInterval: null,
    answers: [],
    questions: [],
    tabSwitches: 0,

    init() {
        fetch('/api/questions')
            .then(r => r.json())
            .then(qs => { this.questions = qs; this.startQuestionTimer(); });
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

    submitQuiz() { /* POST /submit-tes with this.answers */ }
}));
```

### 6.4 Completion Criteria

- [ ] `go build ./...` and `templ generate` pass without errors
- [ ] Quiz page loads 20 questions from API (no correct answers exposed)
- [ ] Each question shows an image + 4 image options in A/B/C/D grid
- [ ] Timer counts down from 120s, changes color at 30s and 10s
- [ ] Auto-advance works when timer hits 0 (records as timed_out)
- [ ] No backward navigation — selecting an answer advances immediately
- [ ] Tab-switch detection logs events
- [ ] Submission sends correct payload format

---

## 7. PHASE 5 — HANDLER & TEMPLATE REWRITE (Required)

**Objective:** Update all HTTP handlers and templ templates to serve the v2.0 cognitive assessment flow. Remove Dark Triad display, update result pages, update landing page branding.

**Estimated effort:** 4–6 hours

---

### 7.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 5.1 | Add `GET /api/questions` endpoint | `handlers/quiz.go` or `handlers/router.go` | Return JSON array of questions WITHOUT `correctOption`. Per IQTEST.md §11.1. |
| 5.2 | Rewrite `POST /submit-tes` handler | `handlers/quiz.go` | Accept `[]{question_code, selected_option, time_taken_ms}`. Call scoring engine. Store results. Return session ID. |
| 5.3 | Rewrite `quizResultToHasilData()` | `handlers/quiz.go` | Map `IQTestResult` to template data. Remove Dark Triad mapping. Add domain scores, reliability flags. |
| 5.4 | Remove duplicate `absInt()` | `handlers/quiz.go` | Ensure no duplicate utility functions. |
| 5.5 | Remove Dark Triad handler references | `handlers/quiz.go` | No more Dark Triad mapping calls. |
| 5.6 | Update landing page hero | `templ/pages/index_page.templ` | "Kenali Kemampuan Kognitifmu" (not kepribadian). Remove MBTI references. Add IQ Test branding. Remove mockup of E/I→L/R badges, 4-letter types. Replace with cognitive ability description. |
| 5.7 | Update landing page features | `templ/pages/index_page.templ` | "Mengapa IQ Test?" — feature cards about cognitive ability assessment. Remove MBTI/Jung personality references. |
| 5.8 | Update landing page how-it-works | `templ/pages/index_page.templ` | Step 3: "Profil Kognitifmu" instead of "tipe MBTI-mu". |
| 5.9 | Remove MBTI marquee/integration section | `templ/pages/index_page.templ` | Remove "MBTI", "Fungsi Kognitif", "16 Tipe" references. |
| 5.10 | Update landing page testimonials | `templ/pages/index_page.templ` | Remove INTJ/INTP type references. Use generic statements. |
| 5.11 | Rewrite result page (`hasil_page.templ`) | `templ/pages/hasil_page.templ` | Show: raw score / max, domain breakdown with progress bars (IQTEST.md §8.3), percentile (if available), estimated IQ (if available — with disclaimer), reliability flags, executive summary (cognitive focus), kekuatan/area perhatian. **Remove Dark Triad cards (Narsisme, Machiavellian, Psikopati).** Add disclaimer per IQTEST.md §7.1 & Appendix B. |
| 5.12 | Update paywall page | `templ/pages/paywall_page.templ` | No MBTI references. "Hasil lengkap IQ Test". |
| 5.13 | Update dashboard page | `templ/pages/dashboard_page.templ` | Show raw score, percentile, estimated IQ (if available), reliability status. Remove "IQ Tipe" column. |
| 5.14 | Update user detail page | `templ/pages/user_detail_page.templ` | Show domain scores with progress bars, reliability flags, response time stats. Remove all SkorLR/NA/SA/LV references. Remove IQTipe. Remove Dark Triad. |
| 5.15 | Update `templ/types/hasil_data.go` | `templ/types/hasil_data.go` | Remove Narsisme/Machiavellian/Psikopati fields. Add domain scores, reliability flags. |
| 5.16 | Update `templ/types/dashboard_data.go` | `templ/types/dashboard_data.go` | Remove IQTipe, SkorLR/NA/SA/LV. Add RawScore, Percentile, EstimatedIQ, DomainScores. |
| 5.17 | Update `handlers/admin.go` | `handlers/admin.go` | Update `ShowDashboard()` and `ShowUserDetail()` to use new field names. |
| 5.18 | Add result page disclaimer | `templ/pages/hasil_page.templ` | Per IQTEST.md §7.1 & Appendix B: *"Tes ini untuk tujuan hiburan dan pengembangan diri, bukan diagnosis klinis. Untuk asesmen resmi, konsultasikan psikolog berlisensi."* |
| 5.19 | Update `assets/js/app.js` (if exists) | `assets/js/app.js` | Check for MBTI references, update or remove. |

### 7.2 Domain Score Visualization (Reference — per IQTEST.md §8.3)

```
Penalaran Matriks   ████████████████░░░░  80%  (Sangat Baik)
Deret Logis         ████████████░░░░░░░░  60%  (Baik)
Rotasi Spasial      ██████░░░░░░░░░░░░░░  30%  (Perlu Latihan)
Analogi Visual      ██████████████░░░░░░  70%  (Baik)
```

### 7.3 Completion Criteria

- [ ] `go build ./...` and `templ generate` pass without errors
- [ ] `GET /api/questions` returns 20 questions without correct answers
- [ ] `POST /submit-tes` accepts new payload format
- [ ] Landing page shows IQ Test branding, no MBTI references
- [ ] Result page shows domain scores, raw score, percentile (no Dark Triad)
- [ ] Result page includes required disclaimer
- [ ] Dashboard shows raw score, no "IQ Tipe" column
- [ ] User detail shows domain breakdown, no Dark Triad
- [ ] Dark Triad display completely removed from all pages
- [ ] All template files have no MBTI/personality references

---

## 8. PHASE 6 — NARRATIVE ENGINE REMOVAL & COGNITIVE PROFILE (Required)

**Objective:** Strip out the old MBTI/personality narrative engine and replace with cognitive ability–focused narratives. Remove Dark Triad narratives entirely.

**Estimated effort:** 3–5 hours

---

### 8.1 What Gets Removed

| Function / Section | Action |
|--------------------|--------|
| `generateRelationshipProfile()` | **Remove entirely** — personality-based relationship analysis tidak relevan untuk tes kognitif |
| `generateRelationshipInsight()` | **Remove entirely** — 8 personality patterns tidak relevan |
| `generateCompatibilityNotes()` | **Remove entirely** — compatibility hanya relevan untuk tes kepribadian |
| `generateReflectionQuestions()` | **Remove entirely** — personal reflection questions tidak relevan |
| All Dark Triad narrative generation | **Remove entirely** — Narcissism/Machiavellianism/Psychopathy narratives dihapus total |

### 8.2 What Gets Updated

| Function / Section | Action |
|--------------------|--------|
| `generateExecutiveSummary()` | Rewrite to describe cognitive performance (raw score, domain strengths, percentile) instead of personality type |
| `generateKekuatan()` | Rewrite: derive from domain scores (highest domain = kekuatan utama) |
| `generateAreaPerhatian()` | Rewrite: derive from domain scores (lowest domain = area pengembangan) |
| `GenerateAllNarratives()` | Remove Dark Triad outputs, add cognitive profile section |

### 8.3 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 6.1 | Remove `generateRelationshipProfile()` | `services/narasi.go` | Delete entire function. |
| 6.2 | Remove `generateRelationshipInsight()` | `services/narasi.go` | Delete entire function. |
| 6.3 | Remove `generateCompatibilityNotes()` | `services/narasi.go` | Delete entire function. |
| 6.4 | Remove `generateReflectionQuestions()` | `services/narasi.go` | Delete entire function. |
| 6.5 | Remove Dark Triad narrative generation | `services/narasi.go` | Remove any code generating Narcissism/Machiavellianism/Psychopathy narratives. |
| 6.6 | Update `generateExecutiveSummary()` | `services/narasi.go` | New signature: `generateExecutiveSummary(rawScore float64, maxScore float64, domainScores map[string]DomainScore, percentile *float64, estimatedIQ *float64) string`. Generate summary about cognitive performance — overall ability, domain balance, relative standing. |
| 6.7 | Update `generateKekuatan()` | `services/narasi.go` | New signature: `generateKekuatan(domainScores map[string]DomainScore) []string`. Return strengths based on top-performing domains. |
| 6.8 | Update `generateAreaPerhatian()` | `services/narasi.go` | New signature: `generateAreaPerhatian(domainScores map[string]DomainScore) []string`. Return development areas based on lowest-performing domains. |
| 6.9 | Rewrite `GenerateAllNarratives()` | `services/narasi.go` | Return only: ExecutiveSummary, Kekuatan, AreaPerhatian. Remove all other narrative output. |
| 6.10 | Update narrative call in `services/quiz.go` | `services/quiz.go` | Update `GetQuizResult()` to call rewritten `GenerateAllNarratives()`. |

### 8.4 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] No `generateRelationshipProfile()`, `generateRelationshipInsight()`, `generateCompatibilityNotes()`, `generateReflectionQuestions()` functions exist
- [ ] No Dark Triad narrative generation code exists
- [ ] `generateExecutiveSummary()` produces cognitive performance summary
- [ ] `generateKekuatan()` and `generateAreaPerhatian()` derive from domain scores
- [ ] `GenerateAllNarratives()` returns only ExecutiveSummary, Kekuatan, AreaPerhatian
- [ ] Old narrative fields removed from `QuizResult` struct

---

## 9. PHASE 7 — ANTI-CHEATING SYSTEM (Required)

**Objective:** Implement anti-cheating mechanisms and reliability scoring per IQTEST.md §9. This ensures test result integrity.

**Estimated effort:** 3–5 days

---

### 9.1 Protections (Per IQTEST.md §9.1)

| Strategy | Implementation | Phase |
|----------|----------------|-------|
| Timer keras 2 menit/soal | Auto-advance frontend + server-side validation | 4 (Frontend) + 7 |
| Jawaban terkunci | Tidak bisa diubah setelah dipilih | 4 (Frontend) |
| Server-side validation | CorrectOption tidak pernah dikirim ke client | 5 (Handler) |
| Randomisasi urutan opsi | Posisi A/B/C/D diacak per sesi | 7 |

### 9.2 Detections (Per IQTEST.md §9.2)

| Detection | Method | Threshold | Action |
|-----------|--------|-----------|--------|
| **Speed-guessing** | time_taken_ms sangat rendah + banyak salah | < 3 detik/soal & akurasi < 25% | Flag "hasil tidak reliabel" |
| **Tab-switch detection** | visibility API | > 3 kali per tes | Flag di hasil |
| **Straight-pattern clicking** | Semua jawaban di opsi sama | ≥ 15/20 sama | Flag "pola respons tidak wajar" |
| **IP rate limiting** | Submission per IP | > 5/jam | 429 Too Many Requests |
| **Devtools/inspect tampering** | Deteksi perubahan DOM | Any | Invalidate sesi |

### 9.3 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 7.1 | Implement server-side answer validation | `services/quiz.go` | Validate that selected_option is A/B/C/D or null. Validate that time_taken_ms is ≥ 0 and ≤ 120,000. |
| 7.2 | Implement speed-guessing detection | `services/quiz.go` | `DetectSpeedGuessing(responses []SessionResponse) bool`. Returns true if >30% answers have time_taken_ms < 3000 and accuracy < 25%. |
| 7.3 | Implement straight-pattern detection | `services/quiz.go` | `DetectStraightPattern(responses []SessionResponse) bool`. Returns true if ≥ 15/20 answers are same option (all A, all B, etc.) or if variance ≤ 0.5. |
| 7.4 | Implement tab-switch tracking | `services/quiz.go` | Accept `tab_switch_count` in submission payload. Flag if > 3. |
| 7.5 | Implement reliability assessment | `services/quiz.go` | `AssessReliability(responses []SessionResponse, tabSwitches int) ReliabilityFlag`. Combine all detection results. Return flag with reasons. |
| 7.6 | Implement IP rate limiting middleware | `middleware/ratelimit.go` (new) | Per-IP rate limit: max 5 POST to `/submit-tes` per hour. Return 429. |
| 7.7 | Add random option shuffling | `services/quiz.go` + `handlers/quiz.go` | Before sending questions to client, shuffle the order of options (A/B/C/D) for each question per session. Store shuffled mapping in session data. |
| 7.8 | Update `handlers/router.go` | `handlers/router.go` | Register rate limiting middleware for `/submit-tes`. |
| 7.9 | Update submission handler | `handlers/quiz.go` | Accept and process `tab_switch_count` from client. Call `AssessReliability()`. Store reliability flags in result. |

### 9.4 ReliabilityFlag Struct (Reference)

```go
type ReliabilityFlag struct {
    IsReliable      bool
    Reasons         []string // "speed_guessing", "tab_switch_excessive", "straight_pattern"
    Recommendation  string   // "hasil_valid" | "disarankan_mengulang"
}
```

### 9.5 Completion Criteria

- [ ] Speed-guessing detection flags appropriately
- [ ] Straight-pattern detection flags appropriately
- [ ] Tab-switch events tracked and flagged
- [ ] `POST /submit-tes` rate limited (max 5/hour per IP)
- [ ] Option order randomized per session
- [ ] `CorrectOption` never reaches client
- [ ] Reliability flags stored with result
- [ ] Result page displays reliability notice if unreliable

---

## 10. PHASE 8 — PAYMENT & PRODUCTION SCHEMA (Future Enhancement)

**Objective:** Implement the full normalized database schema from IQTEST.md §10, including separate tables for users, test sessions, questions, responses, results, payments, and admins. Migrate data from the flat `users_test` table.

**Estimated effort:** 3–5 days

---

### 10.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 8.1 | Create `users` table migration | `migrations/002_production_schema.sql` | `id UUID`, `email`, `nama`, `phone`, `created_at`, `updated_at` |
| 8.2 | Create `test_sessions` table (if not exists) | `migrations/002_production_schema.sql` | Already partially exists — verify and update. Links user to session, tracks device/IP. |
| 8.3 | Create `payments` table | `migrations/002_production_schema.sql` | `user_id`, `session_id`, `amount`, `status`, `payment_method`, `paid_at`, `confirmed_by` |
| 8.4 | Create `admins` table | `migrations/002_production_schema.sql` | Admin accounts with `username`, `password_hash` |
| 8.5 | Data migration script | `migrations/002_migrate_data.sql` | Migrate existing `users_test` records to new normalized schema |
| 8.6 | Implement automated payment gateway | `services/payment.go`, `handlers/payment.go` | Integrate Midtrans/Xendit |
| 8.7 | Drop legacy `users_test` table | `migrations/003_drop_legacy.sql` | After migration verified |

### 10.2 Completion Criteria

- [ ] All new tables created with proper FKs and indexes
- [ ] Data migration preserves all historical test results
- [ ] Payment gateway integration works
- [ ] Legacy table safely archived and dropped

---

## 11. PHASE 9 — ADMIN PANEL UPDATE (Future Enhancement)

**Objective:** Update the admin panel to display v2.0 cognitive ability data (domain scores, raw scores, reliability flags). Add CSV export and analytics.

**Estimated effort:** 2–4 days

---

### 11.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 9.1 | Update admin dashboard statistics | `handlers/admin.go`, `templ/pages/dashboard_page.templ` | Show: average raw score, domain score distribution, reliability rate, paid/unpaid counts, total revenue |
| 9.2 | Update user detail page | `templ/pages/user_detail_page.templ` | Show domain scores with progress bars, reliability flags, response time stats, percentile |
| 9.3 | Add domain score distribution visualization | `templ/pages/dashboard_page.templ` | Bar chart of average % per domain (MTX, SEQ, SPA, ANL) |
| 9.4 | Add CSV export | `handlers/admin.go` | Export user data with v2.0 fields (raw_score, domain %, percentile, estimated_iq, reliability) |
| 9.5 | Add reliability rate dashboard card | `templ/pages/dashboard_page.templ` | % of results flagged as reliable |

### 11.2 Completion Criteria

- [ ] Dashboard shows v2.0 cognitive statistics
- [ ] User detail shows domain breakdown with progress bars
- [ ] Reliability flags visible in admin
- [ ] CSV export works with v2.0 fields

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

### Parallel Execution Possibilities

| Phase Set | Can run in parallel? | Reason |
|-----------|---------------------|--------|
| Phase 0 + Phase 1 | ❌ No | Models must define structs before schema migration |
| Phase 1 + Phase 2 | ❌ No | Repos depend on DB schema |
| Phase 3 + Phase 4 | ✅ Yes | Scoring engine and frontend are independent after models/schema/repos exist |
| Phase 3 + Phase 5 | ❌ No | Handlers depend on scoring engine |
| Phase 4 + Phase 5 | ✅ Yes | Frontend and handlers can be developed concurrently (different files) |
| Phase 6 + Phase 7 | ✅ Yes | Narrative rewrite and anti-cheat are independent |
| Phase 8, 9 | ✅ Yes | Independent of each other and of main migration |

---

## 13. ROLLBACK STRATEGY

### 13.1 Pre-Migration Requirements

Before beginning any phase:

- [ ] Database is backed up (`pg_dump`)
- [ ] Current application binary is tagged in Git (`git tag pre-migration-v2`)
- [ ] All Go dependencies are vendored (`go mod vendor`)
- [ ] Staging environment mirrors production

### 13.2 Per-Phase Rollback

| Phase | Rollback Action | Complexity |
|-------|----------------|------------|
| Phase 0 | `git revert` model changes | Low |
| Phase 1 | Run rollback SQL + restore archive | Medium |
| Phase 2 | `git revert` repository changes + rebuild | Low |
| Phase 3 | `git revert` scoring changes + rebuild | Medium |
| Phase 4 | `git revert` frontend changes + rebuild | Medium |
| Phase 5 | `git revert` handler/template changes + rebuild | Low |
| Phase 6 | `git revert` narrative changes + rebuild | Low |
| Phase 7 | `git revert` anti-cheat changes + rebuild | Low |
| Phase 8 | Run rollback migration + restore from backup | High |
| Phase 9 | `git revert` admin changes + rebuild | Low |

### 13.3 Full Rollback

```bash
# 1. Restore database from backup
pg_restore -d shadowself backup_2026-07-19.dump

# 2. Revert all code changes
git checkout pre-migration-v2

# 3. Rebuild and redeploy
go build -o shadowself .
./shadowself
```

---

## APPENDIX A — FILE CHANGE SUMMARY BY PHASE

| File | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Phase 6 | Phase 7 | Phase 8 | Phase 9 |
|------|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|
| `models/user.go` | ✅ | | | | | | | | ✅ | |
| `models/question.go` (new) | ✅ | | | | | | | | | |
| `models/session.go` (new) | ✅ | | | | | | | | | |
| `database/db.go` | | ✅ | | | | | | | ✅ | |
| `migrations/001_v2_schema.sql` (new) | | ✅ | | | | | | | | |
| `migrations/001_rollback.sql` (new) | | ✅ | | | | | | | | |
| `repositories/user.go` | | | ✅ | | | | | | ✅ | |
| `repositories/admin.go` | | | ✅ | | | ✅ | | | | ✅ |
| `repositories/question.go` (new) | | | ✅ | | | | | | | |
| `repositories/response.go` (new) | | | ✅ | | | | | | | |
| `repositories/result.go` (new) | | | ✅ | | | | | | | |
| `services/quiz.go` | | | | ✅ | | | ✅ | ✅ | | |
| `services/narasi.go` | | | | | | | ✅ | | | |
| `handlers/quiz.go` | | | | | | ✅ | | ✅ | | |
| `handlers/admin.go` | | | | | | ✅ | | | | ✅ |
| `handlers/router.go` | | | | | | | | ✅ | | |
| `middleware/ratelimit.go` (new) | | | | | | | | ✅ | | |
| `templ/types/hasil_data.go` | | | | | | ✅ | | | | |
| `templ/types/dashboard_data.go` | | | | | | ✅ | | | | ✅ |
| `templ/pages/index_page.templ` | | | | | | ✅ | | | | |
| `templ/pages/quiz_page.templ` | | | | | ✅ | ✅ | | ✅ | | |
| `templ/pages/hasil_page.templ` | | | | | | ✅ | | | | |
| `templ/pages/paywall_page.templ` | | | | | | ✅ | | | | |
| `templ/pages/dashboard_page.templ` | | | | | | ✅ | | | | ✅ |
| `templ/pages/user_detail_page.templ` | | | | | | ✅ | | | | ✅ |
| `assets/css/` | | | | | ✅ | | | | | |
| `assets/js/app.js` | | | | | | ✅ | | | | |

---

## APPENDIX B — V2.0 FIELD MAPPING

### Go Struct Fields: Old → New

| Old (v1.x) | New (v2.0) | Notes |
|------------|------------|-------|
| `User.SkorLR int` | `User.RawScore *float64` | Raw score is nullable until test completed |
| `User.SkorNA int` | `User.MaxPossibleScore float64` | Default 30.5 |
| `User.SkorSA int` | `User.MTXScorePct *float64` | Per-domain percentage |
| `User.SkorLV int` | `User.SEQScorePct *float64` | Per-domain percentage |
| `User.IQTipe string` | `User.SPAScorePct *float64` | Per-domain percentage |
| — | `User.ANLScorePct *float64` | New field |
| — | `User.Percentile *float64` | NULL until normative data |
| — | `User.EstimatedIQ *float64` | NULL until normative data |
| — | `User.AvgResponseMs *int` | New field |
| — | `User.IsReliable bool` | New field |
| — | `User.ReliabilityFlags *string` | JSON string |
| `DimensionScore.PoleAScore` | `DomainScore.RawScore` | No bipolar scoring |
| `DimensionScore.PoleBScore` | Removed | No bipolar scoring |
| `DimensionScore.Preference` | Removed | No 4-letter type |
| `DimensionScore.SCI` | Removed | No clarity index |
| `DimensionScore.Strength` | Removed | No strength label |
| `CognitiveProfile` (entire struct) | Removed | No hierarchy |
| `IQTestResult.Type` | Removed | No 4-letter type |
| `IQTestResult.CognitiveProfile` | `IQTestResult.DomainScores` | Domain breakdown |
| — | `IQTestResult.AvgResponseMs` | New field |
| — | `IQTestResult.IsReliable` | New field |
| — | `IQTestResult.ReliabilityFlags` | New field |
| `QuizResult.IQTipe` | Removed | |
| `QuizResult.SkorLR/NA/SA/LV` | `QuizResult.RawScore` | Single raw score |
| `QuizResult.CognitiveProfile` | `QuizResult.DomainScores` | Domain breakdown |
| `QuizResult.RelationshipProfile` | Removed | Dark Triad removed |
| `QuizResult.RelationshipInsight` | Removed | Dark Triad removed |
| `QuizResult.CompatibilityNotes` | Removed | Dark Triad removed |
| `QuizResult.ReflectionQuestions` | Removed | Dark Triad removed |

### Database Columns: Old → New

| Old (v1.x) | New (v2.0) |
|------------|------------|
| `users_test.skor_lr` | `users_test.raw_score` |
| `users_test.skor_na` | `users_test.mtx_score_pct` |
| `users_test.skor_sa` | `users_test.seq_score_pct` |
| `users_test.skor_lv` | `users_test.spa_score_pct` |
| `users_test.iq_tipe` | `users_test.anl_score_pct` |
| — | `users_test.percentile` |
| — | `users_test.estimated_iq` |
| — | `users_test.avg_response_ms` |
| — | `users_test.is_reliable` |
| — | `users_test.reliability_flags` |

---

## APPENDIX C — GLOSSARY (v2.0)

| Term | Definition |
|------|------------|
| **MTX** | Matrix/Pattern Reasoning — kemampuan penalaran induktif dan deteksi pola (mirip Raven's Progressive Matrices) |
| **SEQ** | Logical Sequence — kemampuan penalaran deduktif dan logika sekuensial bergambar |
| **SPA** | Spatial Rotation — kemampuan visualisasi spasial dan mental rotation |
| **ANL** | Visual Analogy — kemampuan penalaran analogis dan abstraksi relasional |
| **Raw Score** | Total skor tertimbang dari jawaban benar (maks. 30.5) |
| **Percentile** | Posisi relatif skor dibanding peserta lain di platform (bukan populasi umum) |
| **Deviation IQ** | Skor IQ dengan mean=100, SD=15, NULL sampai norm data tersedia |
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar suatu soal |
| **Item discrimination** | Seberapa baik soal membedakan peserta berkemampuan tinggi vs rendah |
| **Fluid intelligence** | Kemampuan menalar & memecahkan masalah baru tanpa bergantung pengetahuan yang dipelajari |
| **Reliability** | Konsistensi hasil tes (Cronbach's α, test-retest) — target α ≥ 0.80 |

---

## APPENDIX D — TEST MATRIX

| Test Scenario | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Phase 6 | Phase 7 | Phase 8 | Phase 9 |
|---------------|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|
| `go build ./...` compiles | ⚠️ (expected fail) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `go vet ./...` passes | ⚠️ (expected fail) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `templ generate` succeeds | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Landing page loads correctly | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| Quiz renders 20 image questions | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Timer counts down 120s per question | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Auto-advance on timeout | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Quiz submission correct/incorrect scoring | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Paywall page loads | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| Result shows raw score + domain breakdown | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| No Dark Triad display | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| Disclaimer present on result page | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| Estimated IQ = NULL (no norm data) | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Database migration runs clean | — | ✅ | ✅ | — | — | — | — | — | ✅ | — |
| Anti-cheat flags function | — | — | — | — | — | — | — | ✅ | ✅ | ✅ |
| IP rate limiting active | — | — | — | — | — | — | — | ✅ | ✅ | ✅ |
| Option randomization per session | — | — | — | — | — | — | — | ✅ | ✅ | ✅ |
| Admin dashboard loads | — | — | — | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| Payment gateway integration | — | — | — | — | — | — | — | — | ✅ | ✅ |
| CSV export with v2.0 fields | — | — | — | — | — | — | — | — | — | ✅ |

---

*End of MIGRATION.md v2.0 — Revised Migration Plan*