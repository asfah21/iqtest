# TASK_PROGRESS.md — Migration Progress to v2.0

> **⚠️ Perubahan Fundamental:** IQTEST.md v2.0 (2026-07-19) mengubah metodologi dari tes kepribadian self-report Likert menjadi **tes kemampuan kognitif pilihan ganda bergambar dengan satu jawaban benar objektif**. Semua progres migrasi sebelumnya (berdasarkan v1.0) **tidak lagi relevan** dan telah dihapus. Lihat MIGRATION.md v2.0 untuk rencana lengkap.

---

## Status Keseluruhan

| Fase | Deskripsi | Status | Prioritas |
|------|-----------|--------|-----------|
| **Phase 0** | Models & Types Rewrite | ⏳ Belum dimulai | Required |
| **Phase 1** | Database Schema Rewrite | ⏳ Belum dimulai | Required |
| **Phase 2** | Repository Layer Rewrite | ⏳ Belum dimulai | Required |
| **Phase 3** | Scoring Engine Rewrite | ⏳ Belum dimulai | Required |
| **Phase 4** | Question Bank & Frontend Rewrite | ⏳ Belum dimulai | Required |
| **Phase 5** | Handler & Template Rewrite | ⏳ Belum dimulai | Required |
| **Phase 6** | Narrative Engine Removal & Cognitive Profile | ⏳ Belum dimulai | Required |
| **Phase 7** | Anti-Cheating System | ⏳ Belum dimulai | Required |
| **Phase 8** | Payment & Production Schema | ⏳ Belum dimulai | Future Enhancement |
| **Phase 9** | Admin Panel Update | ⏳ Belum dimulai | Future Enhancement |

---

## Phase 0 — Models & Types Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Rewrite all Go model/type definitions untuk mencocokkan IQTEST.md v2.0.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 0.1 | Rewrite `User` struct — remove SkorLR/NA/SA/LV, IQTipe. Add RawScore, MaxPossibleScore, MTXScorePct, SEQScorePct, SPAScorePct, ANLScorePct, Percentile, EstimatedIQ, AvgResponseMs, IsReliable, ReliabilityFlags | `models/user.go` | ⏳ |
| 0.2 | Create `QuestionDef` struct — metadata soal bergambar dengan correct answer | `models/question.go` (new) | ⏳ |
| 0.3 | Create `SessionResponse` struct — jawaban per soal + timing | `models/user.go` or new file | ⏳ |
| 0.4 | Rewrite `DimensionScore` → `DomainScore` — remove bipolar fields, add Percentage | `models/user.go` | ⏳ |
| 0.5 | Rewrite `IQTestResult` struct — remove Type, CognitiveProfile. Add DomainScores, reliability fields | `models/user.go` | ⏳ |
| 0.6 | Rewrite `QuizResult` struct — remove IQTipe, SkorLR/NA/SA/LV, CognitiveProfile, Dark Triad narrative fields. Add domain scores, reliability | `models/user.go` | ⏳ |
| 0.7 | Create `ReliabilityFlag` struct | `models/user.go` or new `models/models.go` | ⏳ |
| 0.8 | Remove Dark Triad type references entirely | `models/user.go` | ⏳ |

**Completion Criteria:** `go vet ./models/...` passes. All old personality structs removed.

---

## Phase 1 — Database Schema Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Rewrite DB schema untuk soal bergambar, tracking respons per-soal, dan skor kognitif.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1.1 | Create `questions` table migration (with CHECK constraints) | `migrations/001_v2_schema.sql` | ⏳ |
| 1.2 | Create `session_responses` table migration (with FK references) | `migrations/001_v2_schema.sql` | ⏳ |
| 1.3 | Create `iq_results` table migration (with FK references) | `migrations/001_v2_schema.sql` | ⏳ |
| 1.4 | Rewrite `users_test` columns — drop old, add new v2.0 columns | `migrations/001_v2_schema.sql` | ⏳ |
| 1.5 | Archive old data before migration | `migrations/001_archive_old_data.sql` | ⏳ |
| 1.6 | Write rollback migration | `migrations/001_rollback.sql` | ⏳ |
| 1.7 | Update `database/db.go` if needed | `database/db.go` | ⏳ |

**Completion Criteria:** Migration runs clean. `go build ./...` passes.

---

## Phase 2 — Repository Layer Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Rewrite all repository functions untuk schema v2.0.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 2.1 | Rewrite `InsertUser()` — insert new v2.0 columns | `repositories/user.go` | ⏳ |
| 2.2 | Rewrite `GetUserResult()` — SELECT new columns | `repositories/user.go` | ⏳ |
| 2.3 | Rewrite `GetAllUsers()` — SELECT new columns | `repositories/admin.go` | ⏳ |
| 2.4 | Rewrite `GetUserByID()` — SELECT new columns | `repositories/admin.go` | ⏳ |
| 2.5 | Create `InsertQuestion()` | `repositories/question.go` (new) | ⏳ |
| 2.6 | Create `GetActiveQuestions()` | `repositories/question.go` (new) | ⏳ |
| 2.7 | Create `InsertResponse()` | `repositories/response.go` (new) | ⏳ |
| 2.8 | Create `GetSessionResponses()` | `repositories/response.go` (new) | ⏳ |
| 2.9 | Create `InsertIQResult()` | `repositories/result.go` (new) | ⏳ |
| 2.10 | Create `GetIQResultBySession()` | `repositories/result.go` (new) | ⏳ |

**Completion Criteria:** `go build ./...` passes. All old SkorLR/NA/SA/LV references removed.

---

## Phase 3 — Scoring Engine Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Replace personality-based scoring dengan cognitive ability scoring.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 3.1 | Create 20-question bank constant data (QuestionDef) | `services/quiz.go` | ⏳ |
| 3.2 | Rewrite `CalculateIQResult()` — correct/incorrect matching + weighted scoring | `services/quiz.go` | ⏳ |
| 3.3 | Implement domain scoring (MTX/SEQ/SPA/ANL percentages) | `services/quiz.go` | ⏳ |
| 3.4 | Implement percentile calculation (NULL if < 50 results) | `services/quiz.go` | ⏳ |
| 3.5 | Implement IQ estimation stub (returns NULL) | `services/quiz.go` | ⏳ |
| 3.6 | Remove `DeriveCognitiveProfile()` | `services/quiz.go` | ⏳ |
| 3.7 | Remove `mapIQToDarkTriad()` | `services/quiz.go` | ⏳ |
| 3.8 | Remove `axisOpposites` / `axisOpposite()` (verify no remnants) | `services/quiz.go` | ⏳ |
| 3.9 | Rewrite `ProcessQuizAnswers()` — accept multiple-choice responses | `services/quiz.go` | ⏳ |
| 3.10 | Rewrite `GetQuizResult()` — load from iq_results, no CognitiveProfile | `services/quiz.go` | ⏳ |
| 3.11 | Add response time analytics (avg response ms) | `services/quiz.go` | ⏳ |

**Completion Criteria:** `go build ./...` passes. No CognitiveProfile/Dark Triad logic remains.

---

## Phase 4 — Question Bank & Frontend Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Rewrite quiz frontend dari Likert scale menjadi image-based multiple-choice dengan timer.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 4.1 | Rewrite `quizApp` Alpine.js data — timer, auto-advance, answer tracking | `templ/pages/quiz_page.templ` | ⏳ |
| 4.2 | Design new quiz UI — image grid layout, no backward nav | `templ/pages/quiz_page.templ` | ⏳ |
| 4.3 | Add timer visual warnings (kuning @30s, merah @10s) | `templ/pages/quiz_page.templ` | ⏳ |
| 4.4 | Fetch questions from `GET /api/questions` | `templ/pages/quiz_page.templ` | ⏳ |
| 4.5 | Submit answers as `{ question_code, selected_option, time_taken_ms }` | `templ/pages/quiz_page.templ` | ⏳ |
| 4.6 | Add anti-cheat frontend measures (tab-switch detection, disable right-click) | `templ/pages/quiz_page.templ` | ⏳ |
| 4.7 | Update quiz identity form label | `templ/pages/quiz_page.templ` | ⏳ |
| 4.8 | Update CSS for new quiz layout | `assets/css/` | ⏳ |

**Completion Criteria:** Quiz renders 20 image questions with 120s timer. Auto-advance works.

---

## Phase 5 — Handler & Template Rewrite (Required) ⏳ BELUM DIMULAI

**Objective:** Update handlers dan templates untuk v2.0. Remove Dark Triad, update branding.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 5.1 | Add `GET /api/questions` endpoint (no correctOption) | `handlers/quiz.go` / `handlers/router.go` | ⏳ |
| 5.2 | Rewrite `POST /submit-tes` handler (new payload format) | `handlers/quiz.go` | ⏳ |
| 5.3 | Rewrite `quizResultToHasilData()` — remove Dark Triad, add domain scores | `handlers/quiz.go` | ⏳ |
| 5.4 | Remove duplicate `absInt()` | `handlers/quiz.go` | ⏳ |
| 5.5 | Remove Dark Triad handler references | `handlers/quiz.go` | ⏳ |
| 5.6 | Update landing page hero — cognitive ability branding | `templ/pages/index_page.templ` | ⏳ |
| 5.7 | Update landing page features — cognitive ability focus | `templ/pages/index_page.templ` | ⏳ |
| 5.8 | Update landing page how-it-works | `templ/pages/index_page.templ` | ⏳ |
| 5.9 | Remove MBTI marquee/integration section | `templ/pages/index_page.templ` | ⏳ |
| 5.10 | Update landing page testimonials — remove INTJ/INTP references | `templ/pages/index_page.templ` | ⏳ |
| 5.11 | Rewrite hasil_page.templ — domain scores, no Dark Triad, disclaimer | `templ/pages/hasil_page.templ` | ⏳ |
| 5.12 | Update paywall page — "Hasil lengkap IQ Test" | `templ/pages/paywall_page.templ` | ⏳ |
| 5.13 | Update dashboard page — remove "IQ Tipe" column, show raw score | `templ/pages/dashboard_page.templ` | ⏳ |
| 5.14 | Update user detail page — domain progress bars, no Dark Triad | `templ/pages/user_detail_page.templ` | ⏳ |
| 5.15 | Update `templ/types/hasil_data.go` — remove Dark Triad, add domain scores | `templ/types/hasil_data.go` | ⏳ |
| 5.16 | Update `templ/types/dashboard_data.go` — v2.0 fields | `templ/types/dashboard_data.go` | ⏳ |
| 5.17 | Update `handlers/admin.go` — new field names | `handlers/admin.go` | ⏳ |
| 5.18 | Add result page disclaimer | `templ/pages/hasil_page.templ` | ⏳ |
| 5.19 | Update `assets/js/app.js` (if needed) | `assets/js/app.js` | ⏳ |

**Completion Criteria:** No MBTI/Dark Triad references anywhere. Result shows domain scores + disclaimer.

---

## Phase 6 — Narrative Engine Removal & Cognitive Profile (Required) ⏳ BELUM DIMULAI

**Objective:** Hapus narrative engine lama (personality/relationship) dan ganti dengan kognitif.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 6.1 | Remove `generateRelationshipProfile()` | `services/narasi.go` | ⏳ |
| 6.2 | Remove `generateRelationshipInsight()` | `services/narasi.go` | ⏳ |
| 6.3 | Remove `generateCompatibilityNotes()` | `services/narasi.go` | ⏳ |
| 6.4 | Remove `generateReflectionQuestions()` | `services/narasi.go` | ⏳ |
| 6.5 | Remove Dark Triad narrative generation code | `services/narasi.go` | ⏳ |
| 6.6 | Update `generateExecutiveSummary()` — cognitive performance summary | `services/narasi.go` | ⏳ |
| 6.7 | Update `generateKekuatan()` — derive from domain scores | `services/narasi.go` | ⏳ |
| 6.8 | Update `generateAreaPerhatian()` — derive from domain scores | `services/narasi.go` | ⏳ |
| 6.9 | Rewrite `GenerateAllNarratives()` — return only ExecutiveSummary, Kekuatan, AreaPerhatian | `services/narasi.go` | ⏳ |
| 6.10 | Update narrative call in `services/quiz.go` | `services/quiz.go` | ⏳ |

**Completion Criteria:** No relationship/Dark Triad narrative functions exist. Narratives are cognitive-focused.

---

## Phase 7 — Anti-Cheating System (Required) ⏳ BELUM DIMULAI

**Objective:** Implement anti-cheating mechanisms dan reliability scoring.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 7.1 | Implement server-side answer validation | `services/quiz.go` | ⏳ |
| 7.2 | Implement speed-guessing detection | `services/quiz.go` | ⏳ |
| 7.3 | Implement straight-pattern detection | `services/quiz.go` | ⏳ |
| 7.4 | Implement tab-switch tracking | `services/quiz.go` | ⏳ |
| 7.5 | Implement reliability assessment (combine all detections) | `services/quiz.go` | ⏳ |
| 7.6 | Implement IP rate limiting middleware | `middleware/ratelimit.go` (new) | ⏳ |
| 7.7 | Add random option shuffling per session | `services/quiz.go` + `handlers/quiz.go` | ⏳ |
| 7.8 | Update `handlers/router.go` — register rate limiter | `handlers/router.go` | ⏳ |
| 7.9 | Update submission handler — accept tab_switch_count, store reliability | `handlers/quiz.go` | ⏳ |

**Completion Criteria:** Speed-guessing, pattern, tab-switch detection functional. Rate limiting active.

---

## Phase 8 — Payment & Production Schema (Future Enhancement) ⏳ BELUM DIMULAI

**Objective:** Full normalized database schema dan payment gateway.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 8.1 | Create `users` table migration | `migrations/002_production_schema.sql` | ⏳ |
| 8.2 | Create/update `test_sessions` table | `migrations/002_production_schema.sql` | ⏳ |
| 8.3 | Create `payments` table | `migrations/002_production_schema.sql` | ⏳ |
| 8.4 | Create `admins` table | `migrations/002_production_schema.sql` | ⏳ |
| 8.5 | Data migration script (users_test → normalized) | `migrations/002_migrate_data.sql` | ⏳ |
| 8.6 | Implement automated payment gateway (Midtrans/Xendit) | `services/payment.go`, `handlers/payment.go` | ⏳ |
| 8.7 | Drop legacy `users_test` table | `migrations/003_drop_legacy.sql` | ⏳ |

**Completion Criteria:** Normalized schema operational. Payment gateway integrated.

---

## Phase 9 — Admin Panel Update (Future Enhancement) ⏳ BELUM DIMULAI

**Objective:** Update admin panel untuk v2.0 cognitive data.

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 9.1 | Update admin dashboard statistics — raw score, domain distribution, reliability rate | `handlers/admin.go`, `templ/pages/dashboard_page.templ` | ⏳ |
| 9.2 | Update user detail — domain progress bars, reliability | `templ/pages/user_detail_page.templ` | ⏳ |
| 9.3 | Add domain score distribution visualization | `templ/pages/dashboard_page.templ` | ⏳ |
| 9.4 | Add CSV export with v2.0 fields | `handlers/admin.go` | ⏳ |
| 9.5 | Add reliability rate dashboard card | `templ/pages/dashboard_page.templ` | ⏳ |

**Completion Criteria:** Admin shows v2.0 cognitive stats. CSV export works.

---

## Dependency Graph Ringkas

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

                                        Phase 8, 9 (Independent futures)
```

---

## Catatan Penting

1. **Semua progres sebelumnya dihapus** — IQTEST.md v2.0 mengubah fundamental aplikasi.
2. **Phase 0 (Models)** akan memecah build — ini ekspektasi normal karena service/handler belum diperbarui.
3. **Phase 3 (Scoring)** adalah perubahan inti — dari Likert preference → correct/incorrect objective.
4. **Dark Triad dihapus total** — tidak ada dasar psikometri untuk memetakan skor kognitif ke narsisme/dark triad.
5. **IQ Score = NULL** — sampai data normatif ≥1.000 peserta terkumpul, aplikasi tidak menampilkan angka IQ.
6. **Disclaimer wajib** — setiap halaman hasil harus menyertakan peringatan bahwa tes bersifat indikatif.
7. **Lihat MIGRATION.md v2.0** untuk detail lengkap setiap fase, termasuk struktur kode referensi dan test matrix.

---

*End of TASK_PROGRESS.md v2.0 — Updated for revised IQTEST.md*
