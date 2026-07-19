# MIGRATION.md — IQ Test Engine v2.0 Implementation Plan

## Version: 2.0 | Status: Draft | Last Updated: 2026-07-19

> Dokumen ini adalah panduan implementasi sistem **IQ Test Engine v2.0** — tes kemampuan kognitif pilihan ganda bergambar dengan satu jawaban benar objektif. Struktur Go/Gin/PostgreSQL/templ/Alpine.js digunakan sebagai basis.

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
10. [Phase 8 — Payment & Production Schema (Future)](#10-phase-8--payment--production-schema-future)
11. [Phase 9 — Admin Panel (Future)](#11-phase-9--admin-panel-future)
12. [Dependency Graph](#12-dependency-graph)

---

## 1. IMPLEMENTATION OVERVIEW

### 1.1 Target State

Per **IQTEST.md v2.0**, sistem adalah **Cognitive Ability Assessment** dengan:

- 4 domain kognitif: MTX (Matriks/Pola), SEQ (Deret Logis), SPA (Rotasi Spasial), ANL (Analogi Visual)
- 20 soal pilihan ganda bergambar (A/B/C/D) dengan tepat satu jawaban benar
- Skoring berdasarkan jawaban benar/salah dengan bobot kesulitan (max raw = 30.5)
- Timer: 120 detik per soal, hard limit dengan auto-advance
- Raw score + persentil (relatif terhadap pengguna platform)
- Estimasi IQ: **NULL sampai data normatif (1.000+ peserta) terkumpul**
- Breakdown performa per domain (% benar per domain)
- Anti-cheating: speed-guessing detection, tab-switch detection, pattern detection
- Reliability flags pada hasil
- Disclaimer jujur: *"Estimasi ini bersifat indikatif dan belum divalidasi secara klinis"*

### 1.2 Design Principles

| Principle | Description |
|-----------|-------------|
| **Buildable after every phase** | Go code MUST compile after each phase. No broken builds. |
| **Independent phases** | Each phase can be completed and tested in isolation. |
| **Fresh implementation** | Semua tabel dan kode dibuat baru untuk v2.0, tidak ada migrasi dari sistem lama. |
| **No code modification outside scope** | Each phase touches only the files explicitly listed. |

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

**Objective:** Mendefinisikan semua Go model/type struct untuk IQTEST.md v2.0. Ini adalah fondasi untuk semua fase lainnya.

**Estimated effort:** 1–2 hours

---

### 2.1 Struct Definitions

**`models/question.go`** — Metadata soal pilihan ganda bergambar:

```go
// QuestionDef — metadata satu soal pilihan ganda bergambar
type QuestionDef struct {
    ID            string           // UUID
    QuestionCode  string           // e.g., "Q_MTX_001"
    Domain        string           // "MTX" | "SEQ" | "SPA" | "ANL"
    ImageURL      string           // gambar soal utama
    OptionImages  [4]string        // URL gambar opsi A, B, C, D
    CorrectOption string           // "A" | "B" | "C" | "D" (HANYA di server!)
    Difficulty    string           // "easy" | "medium" | "hard" | "very_hard"
    Weight        float64          // 1.0 / 1.5 / 2.0 / 2.5
    PValue        *float64         // nullable — dikalibrasi dari data uji coba
    Discrimination *float64        // nullable — dikalibrasi dari data uji coba
}
```

**`models/session.go`** — Jawaban user per soal:

```go
// SessionResponse — jawaban user untuk satu soal
type SessionResponse struct {
    QuestionID    string
    QuestionCode  string
    SelectedOption *string // "A"/"B"/"C"/"D" atau nil jika timeout
    IsCorrect     bool
    TimeTakenMs   int
    TimedOut      bool
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

// IQTestResult — output akhir kalkulasi
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

// ReliabilityFlag — indikator keandalan hasil tes
type ReliabilityFlag struct {
    IsReliable     bool
    Reasons        []string // "speed_guessing", "tab_switch_excessive", dll
    Recommendation string   // "hasil_valid" | "disarankan_mengulang"
}
```

**`models/user.go`** — Data user dan QuizResult untuk template:

```go
// User — data pengguna yang mengikuti tes
type User struct {
    ID                 string
    Nama               string
    Email              string
    StatusPembayaran   string // "pending" | "paid"
    RawScore           *float64
    MaxPossibleScore   float64
    MTXScorePct        *float64
    SEQScorePct        *float64
    SPAScorePct        *float64
    ANLScorePct        *float64
    Percentile         *float64
    EstimatedIQ        *float64
    AvgResponseMs      *int
    IsReliable         bool
    ReliabilityFlags   *string // JSON string
}

// QuizResult — data yang dikirim ke template hasil
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
```

### 2.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 0.1 | Create `QuestionDef` struct | `models/question.go` | Per IQTEST.md §4.1: ID, QuestionCode, Domain, ImageURL, OptionImages, CorrectOption, Difficulty, Weight, PValue, Discrimination. |
| 0.2 | Create `SessionResponse` struct | `models/session.go` | QuestionID, QuestionCode, SelectedOption (nullable), IsCorrect, TimeTakenMs, TimedOut. |
| 0.3 | Create `DomainScore` struct | `models/result.go` | Domain, RawScore, MaxPossible, Percentage. |
| 0.4 | Create `IQTestResult` struct | `models/result.go` | RawScore, MaxPossible, DomainScores, Percentile, EstimatedIQ, AvgResponseMs, IsReliable, ReliabilityFlags. |
| 0.5 | Create `ReliabilityFlag` struct | `models/result.go` | IsReliable, Reasons, Recommendation. |
| 0.6 | Create `User` struct | `models/user.go` | ID, Nama, Email, StatusPembayaran, raw score fields, domain %, percentile, estimated_iq, avg_response_ms, is_reliable, reliability_flags. |
| 0.7 | Create `QuizResult` struct | `models/user.go` | ID, Nama, RawScore, MaxPossible, Percentile, EstimatedIQ, DomainScores, AvgResponseMs, IsReliable, ReliabilityFlags, ExecutiveSummary, Kekuatan, AreaPerhatian. |
| 0.8 | Create `PaywallData` struct | `models/user.go` | ID string, Nama string. |

### 2.3 Completion Criteria

- [ ] `go vet ./models/...` passes
- [ ] `QuestionDef`, `DomainScore`, `SessionResponse`, `ReliabilityFlag` structs exist
- [ ] All fields match IQTEST.md specification
- [ ] No old personality/MBTI fields exist (SkorLR, IQTipe, CognitiveProfile, Dark Triad, etc.)

---

## 3. PHASE 1 — DATABASE SCHEMA CREATION

**Objective:** Membuat skema database baru untuk v2.0 dari awal. Tidak ada migrasi dari tabel lama — semua tabel baru dibuat dengan `CREATE TABLE`.

**Estimated effort:** 1–2 hours

---

### 3.1 Tables

| Table | Description | Key Columns |
|-------|-------------|-------------|
| `users_test` | Data peserta tes + hasil | `id`, `nama`, `email`, `status_pembayaran`, `raw_score`, `mtx_score_pct`, `seq_score_pct`, `spa_score_pct`, `anl_score_pct`, `percentile`, `estimated_iq`, `avg_response_ms`, `is_reliable`, `reliability_flags` |
| `questions` | Bank soal bergambar | `id`, `question_code`, `domain`, `difficulty`, `weight`, `image_url`, `option_a_image`–`option_d_image`, `correct_option`, `p_value`, `discrimination`, `is_active` |
| `session_responses` | Jawaban per soal per sesi | `id`, `session_id`, `question_id`, `selected_option` (nullable), `is_correct`, `time_taken_ms`, `timed_out` |
| `iq_results` | Hasil kognitif lengkap | `id`, `session_id`, `raw_score`, `max_possible_score`, `mtx_score_pct`–`anl_score_pct`, `percentile`, `estimated_iq`, `avg_response_ms`, `is_reliable`, `reliability_flags` |

### 3.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 1.1 | Create `CREATE TABLE` SQL for `questions` | `migrations/001_init_schema.sql` | `id UUID PK`, `question_code VARCHAR(20) UNIQUE`, `domain VARCHAR(3)` CHECK IN ('MTX','SEQ','SPA','ANL'), `difficulty VARCHAR(10)` CHECK, `weight DECIMAL(3,1)`, `image_url TEXT`, `option_a_image`–`option_d_image TEXT`, `correct_option CHAR(1)` CHECK, `p_value DECIMAL(4,3)`, `discrimination DECIMAL(4,3)`, `is_active BOOLEAN DEFAULT TRUE`, `created_at TIMESTAMPTZ`. |
| 1.2 | Create `CREATE TABLE` SQL for `session_responses` | `migrations/001_init_schema.sql` | `id UUID PK`, `session_id UUID NOT NULL REFERENCES test_sessions(id)`, `question_id UUID NOT NULL REFERENCES questions(id)`, `selected_option CHAR(1)` (nullable), `is_correct BOOLEAN NOT NULL`, `time_taken_ms INTEGER NOT NULL`, `timed_out BOOLEAN DEFAULT FALSE`, `answered_at TIMESTAMPTZ`. UNIQUE(session_id, question_id). |
| 1.3 | Create `CREATE TABLE` SQL for `iq_results` | `migrations/001_init_schema.sql` | `id UUID PK`, `session_id UUID NOT NULL REFERENCES test_sessions(id) UNIQUE`, `raw_score DECIMAL(5,2)`, `max_possible_score DECIMAL(5,2) DEFAULT 30.5`, `mtx_score_pct`–`anl_score_pct DECIMAL(5,1)`, `percentile DECIMAL(5,1)`, `estimated_iq DECIMAL(5,1)`, `avg_response_ms INTEGER`, `is_reliable BOOLEAN DEFAULT TRUE`, `reliability_flags JSONB`, `calculated_at TIMESTAMPTZ`. |
| 1.4 | Create `CREATE TABLE` SQL for `users_test` | `migrations/001_init_schema.sql` | `id UUID PK`, `nama VARCHAR(100)`, `email VARCHAR(100)`, `status_pembayaran VARCHAR(20) DEFAULT 'pending'`, `raw_score DECIMAL(5,2)`, `max_possible_score DECIMAL(5,2) DEFAULT 30.5`, `mtx_score_pct DECIMAL(5,1)`, `seq_score_pct DECIMAL(5,1)`, `spa_score_pct DECIMAL(5,1)`, `anl_score_pct DECIMAL(5,1)`, `percentile DECIMAL(5,1)`, `estimated_iq DECIMAL(5,1)`, `avg_response_ms INTEGER`, `is_reliable BOOLEAN DEFAULT TRUE`, `reliability_flags JSONB`, `created_at TIMESTAMPTZ DEFAULT NOW()`. |
| 1.5 | Create `test_sessions` table | `migrations/001_init_schema.sql` | `id UUID PK`, `user_id UUID REFERENCES users_test(id)`, `ip_address VARCHAR(45)`, `user_agent TEXT`, `option_order JSONB` (shuffled option mapping per session), `started_at TIMESTAMPTZ`, `completed_at TIMESTAMPTZ`. |
| 1.6 | Migration runner | `database/db.go` | Ensure migration runner executes `001_init_schema.sql` on startup. |

### 3.3 SQL Schema (Reference)

```sql
-- migrations/001_init_schema.sql

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

CREATE TABLE users_test (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama                VARCHAR(100) NOT NULL,
    email               VARCHAR(100),
    status_pembayaran   VARCHAR(20) NOT NULL DEFAULT 'pending',
    raw_score           DECIMAL(5,2),
    max_possible_score  DECIMAL(5,2) NOT NULL DEFAULT 30.5,
    mtx_score_pct       DECIMAL(5,1),
    seq_score_pct       DECIMAL(5,1),
    spa_score_pct       DECIMAL(5,1),
    anl_score_pct       DECIMAL(5,1),
    percentile          DECIMAL(5,1),
    estimated_iq        DECIMAL(5,1),
    avg_response_ms     INTEGER,
    is_reliable         BOOLEAN NOT NULL DEFAULT TRUE,
    reliability_flags   JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE test_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users_test(id),
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    option_order    JSONB,     -- shuffled option mapping per question
    tab_switches    INTEGER DEFAULT 0,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ
);

CREATE TABLE session_responses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id        UUID NOT NULL REFERENCES test_sessions(id),
    question_id       UUID NOT NULL REFERENCES questions(id),
    selected_option   CHAR(1),          -- NULL jika timeout
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
```

### 3.4 Completion Criteria

- [ ] Migration runs without errors on empty database
- [ ] All 5 tables created with proper CHECK constraints and FK references
- [ ] `go build ./...` passes
- [ ] Models match schema columns

---

## 4. PHASE 2 — REPOSITORY LAYER

**Objective:** Membuat semua fungsi repository untuk CRUD tabel-tabel v2.0.

**Estimated effort:** 2–3 hours

---

### 4.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 2.1 | Create `InsertUser()` | `repositories/user.go` | INSERT INTO users_test (nama, email) VALUES ($1, $2) RETURNING id. |
| 2.2 | Create `UpdateUserResult()` | `repositories/user.go` | UPDATE users_test SET raw_score, mtx_score_pct, ..., is_reliable, reliability_flags WHERE id = $1. |
| 2.3 | Create `GetUserByID()` | `repositories/user.go` | SELECT all fields by ID. |
| 2.4 | Create `GetAllUsers()` | `repositories/admin.go` | SELECT all users for admin dashboard. |
| 2.5 | Create `InsertQuestion()` | `repositories/question.go` | INSERT into questions table. |
| 2.6 | Create `GetActiveQuestions()` | `repositories/question.go` | SELECT * FROM questions WHERE is_active = TRUE ORDER BY question_code. |
| 2.7 | Create `CreateSession()` | `repositories/session.go` | INSERT INTO test_sessions (user_id, ip_address, user_agent, option_order). |
| 2.8 | Create `UpdateSessionCompleted()` | `repositories/session.go` | SET completed_at, tab_switches WHERE id. |
| 2.9 | Create `InsertResponse()` | `repositories/response.go` | INSERT INTO session_responses. |
| 2.10 | Create `GetSessionResponses()` | `repositories/response.go` | SELECT all responses for a session JOIN questions for domain/weight. |
| 2.11 | Create `InsertIQResult()` | `repositories/result.go` | INSERT INTO iq_results. |
| 2.12 | Create `GetIQResultBySession()` | `repositories/result.go` | SELECT from iq_results by session_id. |
| 2.13 | Create `GetAllRawScores()` | `repositories/result.go` | SELECT raw_score FROM iq_results (for percentile calculation). |

### 4.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] All repository functions use correct v2.0 column names
- [ ] SQL queries match schema from Phase 1

---

## 5. PHASE 3 — SCORING ENGINE

**Objective:** Membuat scoring engine untuk cognitive ability test per IQTEST.md §6.

**Estimated effort:** 3–4 hours

---

### 5.1 Scoring Pipeline

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

### 5.2 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 3.1 | Create question bank constant data | `services/quiz.go` | 20 `QuestionDef` entries per IQTEST.md §4.3 distribution (MTX:6, SEQ:5, SPA:5, ANL:4). Include `CorrectOption`, `Weight`, `Difficulty`. IDs: `Q_MTX_001`…`Q_ANL_004`. **`CorrectOption` hanya di server — tidak pernah dikirim ke client.** |
| 3.2 | Create `CalculateIQResult()` | `services/quiz.go` | Take `[]SessionResponse` + `[]QuestionDef`. For each response: match question, check correctness, accumulate weighted score. Return `IQTestResult`. |
| 3.3 | Create domain scoring logic | `services/quiz.go` | Group questions by domain, sum correct weight / sum max weight → percentage. |
| 3.4 | Create `CalculatePercentile()` | `services/quiz.go` | `CalculatePercentile(rawScore float64, allScores []float64) float64`. Count scores ≤ user score / total scores × 100. Return 0 if no data. |
| 3.5 | Create `EstimateIQ()` (stub) | `services/quiz.go` | `EstimateIQ(rawScore float64, mean float64, stdDev float64) *float64`. Return NULL if mean/stdDev not available. Formula: `100 + (z_score × 15)`. |
| 3.6 | Create `ProcessQuizAnswers()` | `services/quiz.go` | Take map of question_code → selected_option + timing. Build `[]SessionResponse`. Call `CalculateIQResult()`. Store results via repository. |
| 3.7 | Create `GetQuizResult()` | `services/quiz.go` | Load from `iq_results` table via repository. Compute percentile if enough data exists. |

### 5.3 Completion Criteria

- [ ] `CalculateIQResult()` takes `[]SessionResponse` + `[]QuestionDef`, returns `IQTestResult`
- [ ] Domain scoring produces correct MTX/SEQ/SPA/ANL percentages
- [ ] Percentile calculation works (returns 0 when insufficient data)
- [ ] IQ estimation returns NULL when no normative data
- [ ] `go build ./...` passes without errors

---

## 6. PHASE 4 — QUESTION BANK & FRONTEND

**Objective:** Membuat halaman kuis frontend dengan soal bergambar, timer, dan pilihan ganda per IQTEST.md §4–5.

**Estimated effort:** 6–10 hours

---

### 6.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 4.1 | Create `quizApp` Alpine.js component | `templ/pages/quiz_page.templ` | Per IQTEST.md §5.2: `startQuestionTimer()`, `selectAnswer(option)`, `autoAdvance()`. State: `currentQuestion`, `timeRemaining`, `answers`. Timer: 120s countdown, auto-advance on 0. |
| 4.2 | Create quiz UI layout | `templ/pages/quiz_page.templ` | Per IQTEST.md §4.5: gambar soal di atas, 4 opsi bergambar A/B/C/D dalam grid 2×2. Tombol "Sebelumnya" tidak ada. Progress bar + timer countdown. |
| 4.3 | Add timer visual warnings | `templ/pages/quiz_page.templ` | Countdown kuning di 30 detik, merah di 10 detik. Per IQTEST.md §5.1. |
| 4.4 | Fetch questions from API | `templ/pages/quiz_page.templ` | Load `questions[]` from `GET /api/questions` (tanpa `correctOption`). |
| 4.5 | Submit answers | `templ/pages/quiz_page.templ` | `POST /submit-tes` sends array of `{ question_code, selected_option, time_taken_ms, timed_out }`. |
| 4.6 | Add tab-switch detection | `templ/pages/quiz_page.templ` | `visibilitychange` API listener, increment counter on each switch. Include `tab_switch_count` in submission. |
| 4.7 | Create quiz identity form | `templ/pages/quiz_page.templ` | Input nama + email (opsional). |
| 4.8 | Create CSS styles | `assets/css/quiz.css` | Styles for image grid, timer warnings, progress indicators. |

### 6.2 Alpine.js Implementation (Reference)

```javascript
Alpine.data('quizApp', () => ({
    step: 'identity', // 'identity' | 'quiz' | 'submitting'
    currentQuestion: 0,
    timeRemaining: 120,
    timerInterval: null,
    answers: [],
    questions: [],
    tabSwitches: 0,
    nama: '',
    email: '',

    init() {
        // listen for tab switches
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
            window.location.href = '/paywall/' + data.session_id;
        });
    }
}));
```

### 6.3 Completion Criteria

- [ ] `templ generate` passes without errors
- [ ] Quiz page loads 20 questions from API (no correct answers exposed)
- [ ] Each question shows an image + 4 image options in A/B/C/D grid
- [ ] Timer counts down from 120s, changes color at 30s and 10s
- [ ] Auto-advance works when timer hits 0 (records as timed_out)
- [ ] No backward navigation — selecting an answer advances immediately
- [ ] Tab-switch detection logs events
- [ ] Submission sends correct payload format

---

## 7. PHASE 5 — HANDLER & TEMPLATE

**Objective:** Membuat semua HTTP handlers dan template untuk flow IQ Test v2.0.

**Estimated effort:** 4–6 hours

---

### 7.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 5.1 | Create `GET /api/questions` | `handlers/quiz.go` | Return JSON array of 20 questions WITHOUT `correctOption`. Per IQTEST.md §11.1. |
| 5.2 | Create `POST /submit-tes` | `handlers/quiz.go` | Accept `{ nama, email, answers: [{question_code, selected_option, time_taken_ms, timed_out}], tab_switch_count }`. Create user + session. Call scoring engine. Store results. |
| 5.3 | Create `GET /paywall/:id` | `handlers/quiz.go` | Halaman pembayaran untuk membuka hasil. |
| 5.4 | Create `GET /hasil/:id` | `handlers/quiz.go` | Halaman hasil setelah bayar. Map IQTestResult → QuizResult → template data. |
| 5.5 | Create landing page handler | `handlers/page.go` | Serve index_page. |
| 5.6 | Create landing page template | `templ/pages/index_page.templ` | "Kenali Kemampuan Kognitifmu" — IQ Test branding. Fitur: 20 soal bergambar, 4 domain kognitif, timer, estimasi IQ jujur. |
| 5.7 | Create result page template | `templ/pages/hasil_page.templ` | Raw score / max, domain breakdown progress bars (IQTEST.md §8.3), percentile, estimated IQ (with disclaimer), reliability flags, executive summary, kekuatan/area perhatian. |
| 5.8 | Create paywall page template | `templ/pages/paywall_page.templ` | "Hasil lengkap IQ Test — Rp14.900". |
| 5.9 | Create dashboard template | `templ/pages/dashboard_page.templ` | Daftar user dengan raw_score, percentile, estimated IQ, status pembayaran. |
| 5.10 | Create user detail template | `templ/pages/user_detail_page.templ` | Domain scores progress bars, reliability flags, response time stats. |
| 5.11 | Add disclaimer to result page | `templ/pages/hasil_page.templ` | Per IQTEST.md §7.1 & Appendix B. |
| 5.12 | Create types for template data | `templ/types/hasil_data.go`, `templ/types/dashboard_data.go` | Struct definitions matching template needs. |

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
- [ ] `POST /submit-tes` accepts new payload and returns session_id
- [ ] Landing page shows IQ Test branding
- [ ] Result page shows domain scores, raw score, percentile, disclaimer
- [ ] Paywall page shows payment prompt
- [ ] Dashboard shows user results

---

## 8. PHASE 6 — NARRATIVE ENGINE

**Objective:** Membuat narrative engine untuk menghasilkan laporan kognitif berdasarkan skor domain.

**Estimated effort:** 2–3 hours

---

### 8.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 6.1 | Create `generateExecutiveSummary()` | `services/narasi.go` | Input: rawScore, maxScore, domainScores, percentile, estimatedIQ. Output: string deskripsi performa kognitif — kemampuan keseluruhan, keseimbangan domain, posisi relatif. |
| 6.2 | Create `generateKekuatan()` | `services/narasi.go` | Input: domainScores map[string]DomainScore. Output: []string — kekuatan berdasarkan domain dengan skor tertinggi. |
| 6.3 | Create `generateAreaPerhatian()` | `services/narasi.go` | Input: domainScores map[string]DomainScore. Output: []string — area pengembangan berdasarkan domain dengan skor terendah. |
| 6.4 | Create `GenerateAllNarratives()` | `services/narasi.go` | Call all three functions. Return structured narrative output. |

### 8.2 Completion Criteria

- [ ] `generateExecutiveSummary()` produces cognitive performance summary
- [ ] `generateKekuatan()` derives from top-performing domains
- [ ] `generateAreaPerhatian()` derives from lowest-performing domains
- [ ] `go build ./...` passes without errors

---

## 9. PHASE 7 — ANTI-CHEATING SYSTEM

**Objective:** Mengimplementasikan mekanisme anti-cheating dan reliability scoring per IQTEST.md §9.

**Estimated effort:** 3–5 days

---

### 9.1 Protections (Per IQTEST.md §9.1)

| Strategy | Implementation | Phase Source |
|----------|----------------|-------------|
| Timer keras 2 menit/soal | Auto-advance frontend + server-side validation | Phase 4 + 7 |
| Jawaban terkunci | Tidak bisa diubah setelah dipilih | Phase 4 |
| Server-side validation | `CorrectOption` tidak pernah dikirim ke client | Phase 5 + 7 |
| Randomisasi urutan opsi | Posisi A/B/C/D diacak per sesi | Phase 7 |

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
| 7.1 | Implement server-side answer validation | `services/quiz.go` | Validate selected_option is A/B/C/D or null. Validate time_taken_ms is ≥ 0 and ≤ 120,000. |
| 7.2 | Implement speed-guessing detection | `services/quiz.go` | `DetectSpeedGuessing(responses []SessionResponse) bool`. True if >30% answers have time_taken_ms < 3000 and accuracy < 25%. |
| 7.3 | Implement straight-pattern detection | `services/quiz.go` | `DetectStraightPattern(responses []SessionResponse) bool`. True if ≥ 15/20 answers are same option. |
| 7.4 | Implement tab-switch tracking | `services/quiz.go` | Accept `tab_switch_count` in submission. Flag if > 3. |
| 7.5 | Implement `AssessReliability()` | `services/quiz.go` | Combine all detections. Return `ReliabilityFlag`. |
| 7.6 | Implement IP rate limiting middleware | `middleware/ratelimit.go` | Per-IP: max 5 POST to `/submit-tes` per hour. Return 429. |
| 7.7 | Add random option shuffling | `services/quiz.go` + `handlers/quiz.go` | Shuffle option order per session. Store mapping in session. Un-shuffle on submission. |
| 7.8 | Register rate limiting | `handlers/router.go` | Apply rate limit middleware to `/submit-tes`. |

### 9.4 Completion Criteria

- [ ] Speed-guessing detection flags appropriately
- [ ] Straight-pattern detection flags appropriately
- [ ] Tab-switch events tracked and flagged
- [ ] `POST /submit-tes` rate limited (max 5/hour per IP)
- [ ] Option order randomized per session
- [ ] `CorrectOption` never reaches client
- [ ] Reliability flags stored with result
- [ ] `go build ./...` passes

---

## 10. PHASE 8 — PAYMENT & PRODUCTION SCHEMA (Future)

**Objective:** Implementasikan full production schema dengan tabel users, payments, dan integrasi payment gateway.

**Estimated effort:** 3–5 days

---

### 10.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 8.1 | Create `payments` table | `migrations/002_production.sql` | `user_id`, `session_id`, `amount`, `status`, `payment_method`, `paid_at`, `confirmed_by` |
| 8.2 | Implement payment gateway integration | `services/payment.go`, `handlers/payment.go` | Integrate Midtrans/Xendit |
| 8.3 | Payment verification flow | `handlers/payment.go` | Callback handler, status check |

### 10.2 Completion Criteria

- [ ] Payment gateway integration works
- [ ] Payment status tracked in database
- [ ] Result only accessible after payment confirmed

---

## 11. PHASE 9 — ADMIN PANEL (Future)

**Objective:** Lengkapi admin panel dengan statistik kognitif, CSV export, dan analytics.

**Estimated effort:** 2–4 days

---

### 11.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 9.1 | Admin dashboard statistics | `handlers/admin.go`, `templ/pages/dashboard_page.templ` | Average raw score, domain distribution, reliability rate, paid/unpaid counts. |
| 9.2 | User detail view | `templ/pages/user_detail_page.templ` | Domain scores progress bars, reliability flags, response time stats. |
| 9.3 | Domain score distribution chart | `templ/pages/dashboard_page.templ` | Average % per domain (MTX, SEQ, SPA, ANL). |
| 9.4 | CSV export | `handlers/admin.go` | Export with v2.0 fields. |

### 11.2 Completion Criteria

- [ ] Dashboard shows cognitive statistics
- [ ] User detail shows domain breakdown with progress bars
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

## APPENDIX A — FILE INVENTORY

| File | Phase | Action |
|------|-------|--------|
| `models/user.go` | 0 | Create |
| `models/question.go` | 0 | Create |
| `models/session.go` | 0 | Create |
| `models/result.go` | 0 | Create |
| `migrations/001_init_schema.sql` | 1 | Create |
| `database/db.go` | 1 | Update (migration runner) |
| `repositories/user.go` | 2 | Create |
| `repositories/admin.go` | 2 | Create |
| `repositories/question.go` | 2 | Create |
| `repositories/session.go` | 2 | Create |
| `repositories/response.go` | 2 | Create |
| `repositories/result.go` | 2 | Create |
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
| `templ/pages/dashboard_page.templ` | 5 | Create |
| `templ/pages/user_detail_page.templ` | 5 | Create |
| `assets/css/quiz.css` | 4 | Create |

---

## APPENDIX B — GLOSSARY

| Term | Definition |
|------|------------|
| **MTX** | Matrix/Pattern Reasoning — kemampuan penalaran induktif dan deteksi pola |
| **SEQ** | Logical Sequence — kemampuan penalaran deduktif dan logika sekuensial |
| **SPA** | Spatial Rotation — kemampuan visualisasi spasial dan mental rotation |
| **ANL** | Visual Analogy — kemampuan penalaran analogis dan abstraksi relasional |
| **Raw Score** | Total skor tertimbang dari jawaban benar (maks. 30.5) |
| **Percentile** | Posisi relatif skor dibanding peserta lain di platform |
| **Deviation IQ** | Skor IQ mean=100, SD=15, NULL sampai norm data tersedia |
| **Item difficulty (p-value)** | Proporsi peserta yang menjawab benar suatu soal |
| **Item discrimination** | Seberapa baik soal membedakan peserta berkemampuan tinggi vs rendah |
| **Fluid intelligence** | Kemampuan menalar & memecahkan masalah baru tanpa pengetahuan yang dipelajari |
| **Reliability** | Konsistensi hasil tes (Cronbach's α, test-retest) — target α ≥ 0.80 |

---

*End of MIGRATION.md v2.0 — Implementation Plan*