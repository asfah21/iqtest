# MIGRATION.md — MBTI → IQ Test Engine Migration Plan

## Version: 1.0 | Status: Draft | Last Updated: 2026-07-19

---

## TABLE OF CONTENTS

1. [Migration Overview](#1-migration-overview)
2. [Phase 0 — Rename & Rebrand (Required)](#2-phase-0--rename--rebrand-required)
3. [Phase 1 — Scoring Engine Rewrite (Required)](#3-phase-1--scoring-engine-rewrite-required)
4. [Phase 2 — Question Bank & Frontend (Required)](#4-phase-2--question-bank--frontend-required)
5. [Phase 3 — Database Schema Normalization (Required)](#5-phase-3--database-schema-normalization-required)
6. [Phase 4 — Repository Layer Rewrite (Required)](#6-phase-4--repository-layer-rewrite-required)
7. [Phase 5 — Handler & Template Alignment (Required)](#7-phase-5--handler--template-alignment-required)
8. [Phase 6 — Narrative Engine Update (Required)](#8-phase-6--narrative-engine-update-required)
9. [Phase 7 — Anti-Cheating & Reliability (Future Enhancement)](#9-phase-7--anti-cheating--reliability-future-enhancement)
10. [Phase 8 — Production Schema & Payments (Future Enhancement)](#10-phase-8--production-schema--payments-future-enhancement)
11. [Phase 9 — Admin Panel Update (Future Enhancement)](#11-phase-9--admin-panel-update-future-enhancement)
12. [Dependency Graph](#12-dependency-graph)
13. [Rollback Strategy](#13-rollback-strategy)

---

## 1. MIGRATION OVERVIEW

### 1.1 Current State

The application currently implements an **MBTI (Myers-Briggs Type Indicator)** personality assessment using:
- 4 MBTI dichotomies: E/I, S/N, T/F, J/P
- MBTI cognitive function stack (Ni, Se, Fi, Te, etc.)
- 16 personality types (INTJ, ENFP, etc.)
- Dark Triad narrative engine mapped from MBTI raw scores

### 1.2 Target State

Per **IQTEST.md**, the target is an **IQ Test** cognitive assessment using:
- 4 cognitive dimensions: L/R, N/A, S/A, L/V
- Cognitive profile (Dominant, Auxiliary, Complementary, Developing)
- 4-letter IQ test types (LNSL, etc.)
- Score Clarity Index (SCI) per dimension
- Dark Triad narrative engine mapped from IQ raw scores

### 1.3 Design Principles

| Principle | Description |
|-----------|-------------|
| **Buildable after every phase** | The Go code MUST compile after each phase. No broken builds. |
| **Independent phases** | Each phase can be completed and tested in isolation. |
| **Backward compatible where possible** | Old database records remain readable until Phase 3. |
| **No code modification outside scope** | Each phase touches only the files explicitly listed. |
| **Incremental deployment** | Phases can be deployed to production incrementally. |

### 1.4 Phase Dependency Map

```
Phase 0 (Rename) ──► Phase 1 (Scoring) ──► Phase 2 (Frontend)
                          │
                          ▼
                    Phase 3 (Database) ──► Phase 4 (Repositories)
                                                  │
                                                  ▼
                                            Phase 5 (Handlers)
                                                  │
                                                  ▼
                                            Phase 6 (Narratives)

Phase 7, 8, 9 are independent future enhancements
```

---

## 2. PHASE 0 — RENAME & REBRAND (Required)

**Objective:** Rename all identifiers (variables, structs, fields, functions, types, comments) from MBTI terminology to IQ Test terminology. No logic changes. The application will compile and run identically to before, with the same behavior.

**Rationale:** Establishing the correct naming foundation early prevents confusion in later phases. This phase is purely mechanical — simple search-and-replace with zero behavioral change.

**Estimated effort:** 1–2 hours

---

### 2.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 0.1 | Update `models/user.go` — rename struct fields | `models/user.go` | `SkorEI`→`SkorLR`, `SkorSN`→`SkorNA`, `SkorTF`→`SkorSA`, `SkorJP`→`SkorLV`, `MBTITipe`→`IQTipe`, `CognitiveStack`→`CognitiveProfile` (with new field set: Dominant, Auxiliary, Complementary, Developing; remove Tertiary, Inferior), `MBTIResult`→`IQTestResult`, `DikotomiScore`→`DimensionScore`, `PCI`→`SCI` |
| 0.2 | Update `services/quiz.go` — rename functions, vars, comments | `services/quiz.go` | `CalculateMBTI()`→`CalculateIQResult()`, `buildDikotomiScore()`→`buildDimensionScore()`, `DeriveCognitiveStack()`→`DeriveCognitiveProfile()`, accumulators `"EI"`→`"LR"` etc., scores map keys, `mbtiType`→`iqType`, `cognitiveStack`→`cognitiveProfile`, all comments |
| 0.3 | Update `services/quiz.go` — rename Dark Triad mapping function | `services/quiz.go` | `mapMBTIToDarkTriad()`→`mapIQToDarkTriad()`, parameter names `skorEI`→`skorLR` etc. |
| 0.4 | Update `services/narasi.go` — rename function signatures | `services/narasi.go` | References to MBTI/parameter names in comments |
| 0.5 | Update `handlers/quiz.go` — rename references | `handlers/quiz.go` | `ProcessQuizAnswers()` comment, `quizResultToHasilData()` field references |
| 0.6 | Update `handlers/admin.go` — rename field references | `handlers/admin.go` | `u.MBTITipe`→`u.IQTipe`, `u.SkorEI`→`u.SkorLR` etc. |
| 0.7 | Update `templ/types/dashboard_data.go` — rename fields | `templ/types/dashboard_data.go` | `MBTITipe`→`IQTipe`, `SkorEI`→`SkorLR` etc. |
| 0.8 | Update `repositories/user.go` — rename column references | `repositories/user.go` | SQL queries: `skor_ei`→`skor_lr`, `skor_sn`→`skor_na`, `skor_tf`→`skor_sa`, `skor_jp`→`skor_lv`, `mbti_tipe`→`iq_tipe` |
| 0.9 | Update `repositories/admin.go` — rename column references | `repositories/admin.go` | Same SQL column renames as 0.8 |

### 2.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] All MBTI-related identifiers renamed to IQ Test equivalents
- [ ] Application runs with identical behavior (same input → same output)
- [ ] All template files compile (templ code generation succeeds)

---

## 3. PHASE 1 — SCORING ENGINE REWRITE (Required)

**Objective:** Replace the MBTI scoring algorithm with the IQ Test cognitive scoring algorithm per IQTEST.md §6. This is the core behavioral change of the migration. The question bank is NOT reworded in this phase — only the scoring mechanism changes.

**Warning:** This phase CHANGES how results are calculated. Existing database records will be affected if queries pull score data. The database schema columns are renamed in Phase 3, so old records with the old column names will need a migration script.

**Estimated effort:** 3–5 hours

---

### 3.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 1.1 | Rewrite `questionDef` definitions | `services/quiz.go` | Replace all 20 question definitions: `Dikotomi` values `"EI"`→`"LR"` etc., `PolePrimary` values `"E"|"I"`→`"L"|"R"` etc., IDs `Q_EI_001`→`Q_LR_001` etc. |
| 1.2 | Update Likert contribution mapping | `services/quiz.go` | Verify mapping matches IQTEST.md §4.5 (mapping is identical, no change needed) |
| 1.3 | Delete MBTI cognitive function logic | `services/quiz.go` | Remove `axisOpposites` map, `axisOpposite()` function entirely |
| 1.4 | Rewrite `CalculateIQResult()` | `services/quiz.go` | Change accumulators to `"LR"`, `"NA"`, `"SA"`, `"LV"`. Update pole A/B assignment logic per IQTEST.md §6.2 Step 2–3. Keep dimension scores as `map[string]DimensionScore`. |
| 1.5 | Rewrite `DeriveCognitiveProfile()` | `services/quiz.go` | Implement per IQTEST.md §3.4 derivation algorithm (not MBTI cognitive stack). Map 4-letter type to: Dominant, Auxiliary, Complementary, Developing abilities. |
| 1.6 | Add SCI calculation | `services/quiz.go` | `buildDimensionScore()`: rename `PCI`→`SCI`, ensure formula `|rawScore| / maxPossible * 100` per IQTEST.md §3.3 |
| 1.7 | Add Strength label logic | `services/quiz.go` | Ensure strength labels match IQTEST.md §3.3: `slight`(≤25), `moderate`(≤50), `clear`(≤75), `very_clear`(>75) |
| 1.8 | Update `ProcessQuizAnswers()` | `services/quiz.go` | Use new `CalculateIQResult()` return type. Store `skorLR`, `skorNA`, `skorSA`, `skorLV` and `result.IQTipe` |
| 1.9 | Update `GetQuizResult()` | `services/quiz.go` | Use new dimension names `"LR"`, `"NA"`, `"SA"`, `"LV"`. Use `DeriveCognitiveProfile()` instead of `DeriveCognitiveStack()` |
| 1.10 | Update `mapIQToDarkTriad()` | `services/quiz.go` | Per IQTEST.md §8.3: L/R→Narcissism, N/A→Machiavellianism, S/A→Psychopathy. **L/V dimension is NOT mapped** (was J/P→unused before, now L/V→unused). Parameter order: `(skorLR, skorNA, skorSA, skorLV)` |

### 3.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] `services/quiz.go` contains no MBTI dichotomies (EI/SN/TF/JP) in scoring logic
- [ ] `axisOpposites` map and `axisOpposite()` function are removed
- [ ] `CalculateIQResult()` produces 4-letter types using L/R, N/A, S/A, L/V dimensions
- [ ] `DeriveCognitiveProfile()` returns correct Dominant/Auxiliary/Complementary/Developing per IQTEST.md §3.4
- [ ] Unit tests pass for scoring with sample inputs (create test file if none exists)

---

## 4. PHASE 2 — QUESTION BANK & FRONTEND (Required)

**Objective:** Rewrite the 20 questions in the frontend (quiz_page.templ Alpine.js) and backend (services/quiz.go question bank) to reflect IQ Test cognitive dimensions. Questions must measure the 4 cognitive dimensions: Logical/Reasoning, Numerical/Analytical, Spatial/Abstract, Linguistic/Verbal.

**Note:** This is the user-facing content change. The actual question IDs and dikotomi labels change here.

**Estimated effort:** 4–6 hours

---

### 4.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 2.1 | Write new question bank definitions | `services/quiz.go` | 20 new `questionDef` entries with: IDs `Q_LR_001`…`Q_LV_004`, dikotomi `"LR"`/`"NA"`/`"SA"`/`"LV"`, weights per IQTEST.md §4.1 (L/R:5, N/A:6, S/A:5, L/V:4), reverse-scored per IQTEST.md Appendix B |
| 2.2 | Rewrite frontend question text | `templ/pages/quiz_page.templ` | Replace all 20 question objects in `questions: [...]` array inside `<script>`. Questions must measure the 4 cognitive dimensions. Dikotomi badges: `L/R`, `N/A`, `S/A`, `L/V`. IDs: `Q_LR_001`…`Q_LV_004`. |
| 2.3 | Update question ID list in handler | `handlers/quiz.go` | Replace question IDs list with `Q_LR_001`…`Q_LV_004` |
| 2.4 | Update quiz progress labels | `templ/pages/quiz_page.templ` | Update "20" count (remains 20). Update any dikotomi-related labels. |
| 2.5 | Add JSONB translations support (optional) | `services/quiz.go` | Per IQTEST.md §10.3, questions model should support translations JSONB. (Can be deferred to later phase — field reserved but unused.) |

### 4.2 Question Distribution (IQTEST.md §4.1)

| Dimension | Count | Weights | Max Score |
|-----------|-------|---------|-----------|
| L/R | 5 | 2.0, 2.0, 1.5, 1.5, 1.5 | 8.5 |
| N/A | 6 | 2.0, 2.0, 1.5, 1.5, 1.5, 1.5 | 10.0 |
| S/A | 5 | 2.0, 2.0, 1.5, 1.5, 1.5 | 8.5 |
| L/V | 4 | 2.0, 2.0, 1.5, 1.5 | 7.0 |
| **Total** | **20** | | |

### 4.3 Reverse-Scored Questions (IQTEST.md Appendix B)

| Question ID | Dimension |
|-------------|-----------|
| Q_LR_005 | L/R |
| Q_NA_005 | N/A |
| Q_NA_006 | N/A |
| Q_SA_005 | S/A |
| Q_LV_004 | L/V |

### 4.4 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] Frontend renders 20 questions with correct dikotomi badges (L/R, N/A, S/A, L/V)
- [ ] Question text reflects Logical/Reasoning, Numerical/Analytical, Spatial/Abstract, Linguistic/Verbal dimensions
- [ ] Form submission sends correct question IDs (`q_Q_LR_001` etc.)
- [ ] Scoring produces valid 4-letter IQ types (e.g., LNSL, RAVL, etc.)

---

## 5. PHASE 3 — DATABASE SCHEMA NORMALIZATION (Required)

**Objective:** Rename columns in the existing `users_test` table to match IQ Test terminology. This is a lightweight migration that keeps the single-table structure. The full normalized schema (6 tables from IQTEST.md §10.3) is deferred to Phase 8.

**Estimated effort:** 1 hour (+ data migration time depending on volume)

---

### 5.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 3.1 | Create database migration SQL | `migrations/001_rename_mbti_columns.sql` | `ALTER TABLE users_test RENAME COLUMN skor_ei TO skor_lr;` (same for sn→na, tf→sa, jp→lv, mbti_tipe→iq_tipe) |
| 3.2 | Update `database/db.go` (if needed) | `database/db.go` | No changes expected (connection string stays same, no schema references in code) |
| 3.3 | Write rollback migration | `migrations/001_rollback.sql` | Reverse all column renames for safe rollback |
| 3.4 | Test migration on staging | — | Run migration, verify data integrity, run application |

### 5.2 SQL Migration

```sql
-- migrations/001_rename_mbti_columns.sql
ALTER TABLE users_test RENAME COLUMN skor_ei TO skor_lr;
ALTER TABLE users_test RENAME COLUMN skor_sn TO skor_na;
ALTER TABLE users_test RENAME COLUMN skor_tf TO skor_sa;
ALTER TABLE users_test RENAME COLUMN skor_jp TO skor_lv;
ALTER TABLE users_test RENAME COLUMN mbti_tipe TO iq_tipe;
```

### 5.3 Rollback SQL

```sql
-- migrations/001_rollback.sql
ALTER TABLE users_test RENAME COLUMN skor_lr TO skor_ei;
ALTER TABLE users_test RENAME COLUMN skor_na TO skor_sn;
ALTER TABLE users_test RENAME COLUMN skor_sa TO skor_tf;
ALTER TABLE users_test RENAME COLUMN skor_lv TO skor_jp;
ALTER TABLE users_test RENAME COLUMN iq_tipe TO mbti_tipe;
```

### 5.4 Completion Criteria

- [ ] Migration runs without errors
- [ ] `SELECT * FROM users_test` returns columns: `skor_lr`, `skor_na`, `skor_sa`, `skor_lv`, `iq_tipe`
- [ ] Existing data is preserved (no data loss)
- [ ] Application reads/writes to renamed columns correctly
- [ ] Rollback script is verified to work

---

## 6. PHASE 4 — REPOSITORY LAYER REWRITE (Required)

**Objective:** Update all repository queries to use renamed columns and return updated struct types. This phase is tightly coupled with Phase 3 — the database columns must already be renamed.

**Estimated effort:** 1–2 hours

---

### 6.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 4.1 | Update `repositories/user.go` — `InsertUser()` | `repositories/user.go` | Update SQL: `skor_ei`→`skor_lr` etc., parameter names `skorEI`→`skorLR` etc., `mbtiTipe`→`iqTipe` |
| 4.2 | Update `repositories/user.go` — `GetUserResult()` | `repositories/user.go` | Update SQL column references and struct field scans |
| 4.3 | Update `repositories/admin.go` — `GetAllUsers()` | `repositories/admin.go` | Update SQL column references and struct field scans |
| 4.4 | Update `repositories/admin.go` — `GetUserByID()` | `repositories/admin.go` | Update SQL column references and struct field scans |

### 6.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] All repository functions use renamed columns
- [ ] `InsertUser()` stores data in correct columns
- [ ] `GetUserResult()` returns correct field values
- [ ] `GetAllUsers()` and `GetUserByID()` return data correctly

---

## 7. PHASE 5 — HANDLER & TEMPLATE ALIGNMENT (Required)

**Objective:** Update all handler logic and templ templates to use new field names, dimension labels, and display data correctly. No behavioral changes — only identifier updates and display label updates.

**Estimated effort:** 2–3 hours

---

### 7.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 5.1 | Update `handlers/quiz.go` — `quizResultToHasilData()` | `handlers/quiz.go` | Map `SkorLR`, `SkorNA`, `SkorSA`, `SkorLV` to Dark Triad values using updated `mapIQToDarkTriad()`. Update `MBTI`→`IQTipe` reference. |
| 5.2 | Remove duplicate `absInt()` | `handlers/quiz.go` | Delete lines 135-139 (duplicate of `services/quiz.go` version). Use `services.absInt()` instead. |
| 5.3 | Update `handlers/admin.go` — `ShowDashboard()` | `handlers/admin.go` | Update field references: `u.MBTITipe`→`u.IQTipe`, `SkorEI`→`SkorLR` etc. |
| 5.4 | Update `templ/types/hasil_data.go` — if needed | `templ/types/hasil_data.go` | Fields Narsisme, Machiavellian, Psikopati remain — source data changes only |
| 5.5 | Update `templ/types/dashboard_data.go` — field renames | `templ/types/dashboard_data.go` | `MBTITipe`→`IQTipe`, `SkorEI`→`SkorLR` etc. |
| 5.6 | Update `templ/pages/index_page.templ` — hero section | `templ/pages/index_page.templ` | Hero headline: "Kenali Tipe Kepribadianmu" → "Kenali Kemampuan Kognitifmu". Subheadline: remove MBTI reference, add IQ Test reference. Mockup: change E/I→L/R etc., Tipe: INTJ→LNSL. |
| 5.7 | Update `templ/pages/index_page.templ` — features section | `templ/pages/index_page.templ` | "Mengapa MBTI?" → "Mengapa IQ Test?". Feature cards: replace MBTI/Jung references with cognitive ability references. |
| 5.8 | Update `templ/pages/index_page.templ` — how it works | `templ/pages/index_page.templ` | Step 3: "tipe MBTI-mu" → "tipe IQ Test-mu" |
| 5.9 | Update `templ/pages/index_page.templ` — integration marquee | `templ/pages/index_page.templ` | Replace "MBTI", "Fungsi Kognitif", "16 Tipe", "4 Dikotomi" with IQ Test terminology |
| 5.10 | Update `templ/pages/index_page.templ` — testimonials | `templ/pages/index_page.templ` | Replace MBTI type references (INTJ, INTP) with IQ cognitive type references or generic statements |
| 5.11 | Update `templ/pages/quiz_page.templ` — identity form label | `templ/pages/quiz_page.templ` | "hasil asesmen MBTI" → "hasil asesmen IQ Test" |
| 5.12 | Update `templ/pages/dashboard_page.templ` — table header | `templ/pages/dashboard_page.templ` | "MBTI" → "IQ Tipe" |
| 5.13 | Update `templ/pages/user_detail_page.templ` — labels | `templ/pages/user_detail_page.templ` | "Tipe MBTI" → "IQ Tipe", "Skor E/I" → "Skor L/R", "Skor S/N" → "Skor N/A", "Skor T/F" → "Skor S/A", "Skor J/P" → "Skor L/V". Field references: `.MBTITipe`→`.IQTipe`, `.SkorEI`→`.SkorLR` etc. |
| 5.14 | Update `templ/pages/hasil_page.templ` — if needed | `templ/pages/hasil_page.templ` | Verify no MBTI references remain. Dark Triad display cards are acceptable (Narsisme, Machiavellian, Psikopati). |
| 5.15 | Update `assets/js/app.js` — if exists | `assets/js/app.js` | Check for MBTI references and update |

### 7.2 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] `templ` code generation succeeds for all templates
- [ ] Landing page (/) shows IQ Test branding, no MBTI references
- [ ] Dashboard page shows "IQ Tipe" column, correct dimension labels
- [ ] User detail page shows correct dimension labels
- [ ] Quiz page identity form says "IQ Test"
- [ ] Result page displays correctly

---

## 8. PHASE 6 — NARRATIVE ENGINE UPDATE (Required)

**Objective:** Update the narrative generator (`services/narasi.go`) to reference cognitive abilities and IQ Test terminology instead of personality-based framing. The Dark Triad mapping mechanism stays, but the narrative context shifts from "personality traits in relationships" to "cognitive abilities and their relational implications."

**Estimated effort:** 3–5 hours

---

### 8.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 6.1 | Update `services/narasi.go` — function parameter names | `services/narasi.go` | All parameter names using EI/SN/TF/JP→LR/NA/SA/LV (none found — parameters are `n, m, p int` generic) |
| 6.2 | Update executive summary narrative | `services/narasi.go` | `generateExecutiveSummary()`: Replace "kepribadian" framing with "kemampuan kognitif" framing. Add IQ Test type badge reference. Update to reflect cognitive ability context. |
| 6.3 | Update relationship profile sections | `services/narasi.go` | `generateRelationshipProfile()`: The 5-axis analysis (validasi, konflik, keputusan, pengaruh, emosional) remains structurally valid. Update any personality-specific phrasing to cognitive ability framing. |
| 6.4 | Update kekuatan and area perhatian | `services/narasi.go` | `generateKekuatan()` and `generateAreaPerhatian()`: Ensure references align with cognitive abilities vs personality traits. |
| 6.5 | Update relationship insight patterns | `services/narasi.go` | `generateRelationshipInsight()`: Update 8 personality patterns to reflect cognitive ability combinations. Add "cognitive profile" pattern. |
| 6.6 | Update compatibility notes | `services/narasi.go` | `generateCompatibilityNotes()`: Maintain structure, update terminology. |
| 6.7 | Update reflection questions | `services/narasi.go` | `generateReflectionQuestions()`: Maintain structure, update terminology. |
| 6.8 | Add cognitive profile narrative section | `services/narasi.go` | **New function** `generateCognitiveProfile(iqType string, scores map[string]DimensionScore) string`: Generate narrative describing the 4-letter IQ type, dominant/auxiliary/complementary/developing abilities per IQTEST.md §8.5. |
| 6.9 | Update `GenerateAllNarratives()` | `services/narasi.go` | Add cognitive profile narrative to output. Update signature if needed. |
| 6.10 | Update `GetQuizResult()` narrative call | `services/quiz.go` | Pass cognitive profile narrative data to template |

### 8.2 Narrative Context Shift

| Before (MBTI/Personality) | After (IQ Test/Cognitive) |
|---------------------------|--------------------------|
| "tipe kepribadian" | "profil kognitif" |
| "fungsi kognitif" | "kemampuan kognitif" |
| "INTJ, ENFP" | "LNSL, RAVL" |
| "introvert/ekstrovert" | "deduktif/induktif" |
| "sensing/intuition" | "numerik/analitis" |
| Persona-based analysis | Ability-based analysis |

### 8.3 Completion Criteria

- [ ] `go build ./...` passes without errors
- [ ] Narrative texts no longer reference MBTI personality traits
- [ ] Executive summary references cognitive abilities and IQ Test type
- [ ] Cognitive profile section is present in result output
- [ ] All narrative sections compile and render correctly
- [ ] Result page shows complete cognitive profile

---

## 9. PHASE 7 — ANTI-CHEATING & RELIABILITY (Future Enhancement)

**Objective:** Implement anti-cheating mechanisms and reliability scoring per IQTEST.md §9. This phase adds new functionality without breaking existing flows.

**Estimated effort:** 2–3 weeks

---

### 9.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 7.1 | Implement response time tracking | `templ/pages/quiz_page.templ`, `services/quiz.go` | Frontend: capture `time_taken_ms` per question per IQTEST.md §5.3. Backend: store in calculation pipeline. |
| 7.2 | Implement inconsistency scoring | `services/quiz.go` | `CalculateInconsistency()`: Compare anchor pairs per IQTEST.md §9.3. Define anchor pairs for cognitive dimensions. |
| 7.3 | Implement confidence scoring | `services/quiz.go` | `CalculateConfidence()`: Combine completion rate, response time, inconsistency score per IQTEST.md §9.4. Return `ConfidenceScore` struct. |
| 7.4 | Implement straight-line detection | `services/quiz.go` | Detect all-same-answer patterns. Flag if variance ≤ 0.5. |
| 7.5 | Implement tab-switch detection | `templ/pages/quiz_page.templ` | Frontend: use `visibilitychange` API per IQTEST.md §9.5. Send tab-switch timestamps. |
| 7.6 | Implement IP rate limiting | `handlers/router.go`, `middleware/ratelimit.go` | Per-IP rate limiting for POST `/submit-tes`. |
| 7.7 | Add frontend anti-cheat measures | `templ/pages/quiz_page.templ` | Disable right-click, text selection, copy-paste, back navigation during quiz per IQTEST.md §9.5. |

### 9.2 New Structures

```go
// ConfidenceScore — combined reliability indicator per IQTEST.md §9.4
type ConfidenceScore struct {
    Overall        float64           // 0–100
    PerDikotomi    map[string]float64 // SCI per dimension
    Flags          []string          // Warning flags
    Recommendation string            // "reliable" | "review_suggested" | "retest_recommended"
}
```

### 9.3 Completion Criteria

- [ ] Response times tracked and stored per question
- [ ] Inconsistency score calculated and stored
- [ ] Confidence score displayed in result (or stored for admin view)
- [ ] Straight-line answers flagged
- [ ] Tab-switch events logged
- [ ] IP rate limiting active
- [ ] Frontend anti-cheat measures in place

---

## 10. PHASE 8 — PRODUCTION SCHEMA & PAYMENTS (Future Enhancement)

**Objective:** Implement the full normalized database schema from IQTEST.md §10.3, including separate tables for users, test sessions, questions, responses, results, payments, and admins. Migrate data from the flat `users_test` table.

**Estimated effort:** 3–5 days

---

### 10.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 8.1 | Create `users` table migration | `migrations/002_production_schema.sql` | Per IQTEST.md §10.3: `id`, `email`, `nama`, `phone`, `created_at`, `updated_at` |
| 8.2 | Create `test_sessions` table | `migrations/002_production_schema.sql` | Per IQTEST.md §10.3: Links user to session, tracks device/IP |
| 8.3 | Create `questions` table | `migrations/002_production_schema.sql` | Dynamic question bank: `question_code`, `dikotomi`, `pole_primary`, `weight`, `reverse_scored`, `is_active`, `translations` |
| 8.4 | Create `session_responses` table | `migrations/002_production_schema.sql` | Per-question response storage: `session_id`, `question_id`, `answer_value`, `time_taken_ms` |
| 8.5 | Create `iq_results` table | `migrations/002_production_schema.sql` | Full result: `iq_type`, `lr_raw_score`, `na_raw_score`, `sa_raw_score`, `lv_raw_score`, `lr_sci`, `na_sci`, `sa_sci`, `lv_sci`, `cognitive_profile` JSONB, `completion_rate`, `avg_response_ms`, `inconsistency_score`, `is_reliable` |
| 8.6 | Create `payments` table | `migrations/002_production_schema.sql` | Payment tracking: `user_id`, `session_id`, `amount`, `status`, `payment_method`, `paid_at` |
| 8.7 | Create `admins` table | `migrations/002_production_schema.sql` | Admin accounts with `username`, `password_hash` |
| 8.8 | Create database indexes | `migrations/002_production_schema.sql` | Per IQTEST.md §10.3 index list |
| 8.9 | Data migration script | `migrations/002_migrate_data.sql` | Migrate existing `users_test` records to new schema: user data → `users`, result data → `iq_results` |
| 8.10 | Create new repository files | `repositories/user.go` (rewrite), `repositories/session.go`, `repositories/question.go`, `repositories/result.go`, `repositories/payment.go` | New data access layer for normalized schema |
| 8.11 | Update service layer | `services/quiz.go` | Update to use new repositories and session-based flow |
| 8.12 | Implement session-based quiz flow | `handlers/quiz.go` | Create session on quiz start, store responses per session |
| 8.13 | Implement automated payment gateway | `services/payment.go`, `handlers/payment.go` | Integrate Midtrans/Xendit for instant payment confirmation |
| 8.14 | Drop legacy `users_test` table | `migrations/003_drop_legacy.sql` | After migration verified, drop old table |

### 10.2 Entity Relationship (per IQTEST.md §10.4)

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
questions ──> questions (self-referencing or via translations)
```

### 10.3 Completion Criteria

- [ ] All new tables created and indexed
- [ ] Data migration preserves all historical test results
- [ ] Session-based quiz flow works end-to-end
- [ ] Dynamic question bank loads from database
- [ ] Payment tracking works with automated gateway
- [ ] Legacy `users_test` table is safely dropped
- [ ] All tests pass

---

## 11. PHASE 9 — ADMIN PANEL UPDATE (Future Enhancement)

**Objective:** Update the admin panel to work with the new schema and display IQ Test-specific data (SCI scores, cognitive profiles, confidence scores). Add CSV export and analytics.

**Estimated effort:** 2–4 days

---

### 11.1 Tasks

| # | Task | File(s) | Description |
|---|------|---------|-------------|
| 9.1 | Update admin dashboard statistics | `handlers/admin.go`, `templ/pages/dashboard_page.templ` | Show IQ type distribution, average SCI per dimension, paid/unpaid counts, total revenue |
| 9.2 | Update user detail page | `templ/pages/user_detail_page.templ` | Show dimension scores (LR/NA/SA/LV), SCI values, cognitive profile, confidence score |
| 9.3 | Add IQ type distribution chart | `templ/pages/dashboard_page.templ` | Bar chart or pie chart of most common IQ types |
| 9.4 | Add CSV export | `handlers/admin.go` | Export user data to CSV with IQ Test fields |
| 9.5 | Add admin notification for payments | `services/notification.go` (new) | Notify admin on new payment confirmations |

### 11.2 Completion Criteria

- [ ] Dashboard shows IQ Test statistics
- [ ] User detail shows SCI values and cognitive profile
- [ ] IQ type distribution chart renders
- [ ] CSV export works
- [ ] Payment notifications operational

---

## 12. DEPENDENCY GRAPH

```
Phase 0 (Rename) ─────────── required by ──► All later phases
       │
       ▼
Phase 1 (Scoring Engine) ─── required by ──► Phase 2, 3, 4, 5, 6
       │
       ▼
Phase 2 (Question Bank) ──── required by ──► Phase 1 (concurrent OK)
       │
       ▼
Phase 3 (DB Schema) ──────── required by ──► Phase 4
       │
       ▼
Phase 4 (Repositories) ───── required by ──► Phase 5
       │
       ▼
Phase 5 (Handlers+Template) ─ required by ──► Phase 6
       │
       ▼
Phase 6 (Narratives) ──────── leaf phase
```

### Parallel Execution Possibilities

| Phase Set | Can run in parallel? | Reason |
|-----------|---------------------|--------|
| Phase 0 + Phase 2 | ✅ Yes | Rename doesn't affect frontend questions |
| Phase 1 + Phase 2 | ✅ Yes | Scoring logic and question content are independent after rename |
| Phase 3 + Phase 4 | ❌ No | Phase 4 depends on Phase 3 DB changes |
| Phase 5 + Phase 6 | ✅ Yes | Handler changes and narrative changes are independent |
| Phase 7, 8, 9 | ✅ Yes | All independent of each other |

---

## 13. ROLLBACK STRATEGY

### 13.1 Pre-Migration Requirements

Before beginning any phase:

- [ ] Database is backed up (`pg_dump`)
- [ ] Current application binary is tagged in Git (`git tag pre-migration-v1`)
- [ ] All Go dependencies are vendored (`go mod vendor`)
- [ ] Staging environment mirrors production

### 13.2 Per-Phase Rollback

| Phase | Rollback Action | Complexity |
|-------|----------------|------------|
| Phase 0 | `git revert` the rename commit | Low |
| Phase 1 | `git revert` scoring changes + rebuild | Medium (data compatibility may break) |
| Phase 2 | `git revert` question changes | Low |
| Phase 3 | Run rollback SQL (`migrations/001_rollback.sql`) | Low |
| Phase 4 | `git revert` repository changes + rebuild | Low |
| Phase 5 | `git revert` template/handler changes + rebuild | Low |
| Phase 6 | `git revert` narrative changes + rebuild | Low |
| Phase 7 | `git revert` + rebuild | Low |
| Phase 8 | Run rollback migration + restore from backup | High (data normalization is irreversible) |
| Phase 9 | `git revert` + rebuild | Low |

### 13.3 Full Rollback

```bash
# 1. Restore database from backup
pg_restore -d shadowself backup_2026-07-19.dump

# 2. Revert all code changes
git checkout pre-migration-v1

# 3. Rebuild and redeploy
go build -o shadowself .
./shadowself
```

---

## APPENDIX A — FILE CHANGE SUMMARY BY PHASE

| File | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Phase 6 | Phase 7 | Phase 8 | Phase 9 |
|------|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|
| `models/user.go` | ✅ | | | | | | | | ✅ | |
| `services/quiz.go` | ✅ | ✅ | ✅ | | | | ✅ | ✅ | ✅ | |
| `services/narasi.go` | ✅ | | | | | | ✅ | | | |
| `handlers/quiz.go` | ✅ | | ✅ | | | ✅ | | | ✅ | |
| `handlers/admin.go` | ✅ | | | | | ✅ | | | | ✅ |
| `handlers/router.go` | | | | | | | | ✅ | | |
| `handlers/page.go` | | | | | | ✅ | | | | |
| `helpers/render.go` | | | | | | | | | | |
| `repositories/user.go` | ✅ | | | | ✅ | | | | ✅ | |
| `repositories/admin.go` | ✅ | | | | ✅ | | | | | |
| `templ/types/hasil_data.go` | | | | | | ✅ | | | | |
| `templ/types/dashboard_data.go` | ✅ | | | | | ✅ | | | | ✅ |
| `templ/pages/index_page.templ` | | | | | | ✅ | | | | |
| `templ/pages/quiz_page.templ` | | | ✅ | | | ✅ | | ✅ | | |
| `templ/pages/hasil_page.templ` | | | | | | ✅ | | | | |
| `templ/pages/dashboard_page.templ` | | | | | | ✅ | | | | ✅ |
| `templ/pages/user_detail_page.templ` | | | | | | ✅ | | | | ✅ |
| `assets/js/app.js` | | | | | | ✅ | | | | |
| `database/db.go` | | | | | | | | | ✅ | |
| `go.mod` / `go.sum` | | | | | | | | | ✅ | |
| `migrations/` (new) | | | | ✅ | | | | | ✅ | |

---

## APPENDIX B — TEST MATRIX

| Test Scenario | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Phase 6 | Phase 7 | Phase 8 | Phase 9 |
|---------------|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|:-------:|
| `go build ./...` compiles | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Landing page loads | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Quiz renders 20 questions | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Quiz submission returns user ID | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Scoring produces 4-letter type | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Paywall page loads | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Payment confirmation works | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Result page renders with narratives | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Admin dashboard loads | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Admin user detail loads | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Database migration runs clean | | | | ✅ | | | | | ✅ | |
| Old records remain readable | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Anti-cheat flags (if implemented) | | | | | | | | ✅ | ✅ | |
| Payment gateway integration | | | | | | | | | ✅ | |

---

## APPENDIX C — GLOSSARY

| Term | Definition |
|------|------------|
| **L/R** | Logical/Reasoning — cognitive dimension (replaces E/I) |
| **N/A** | Numerical/Analytical — cognitive dimension (replaces S/N) |
| **S/A** | Spatial/Abstract — cognitive dimension (replaces T/F) |
| **L/V** | Linguistic/Verbal — cognitive dimension (replaces J/P) |
| **SCI** | Score Clarity Index — measures strength of aptitude (0–100%) (replaces PCI) |
| **Cognitive Profile** | Hierarchy of 4 cognitive abilities: Dominant, Auxiliary, Complementary, Developing (replaces Cognitive Stack) |
| **IQ Type** | 4-letter identifier (e.g., LNSL) derived from dimension preferences (replaces MBTI type) |
| **Dimension** | A pair of opposing cognitive abilities (replaces Dikotomi) |
| **DimensionScore** | Scoring result for one cognitive dimension (replaces DikotomiScore) |

---

*End of MIGRATION.md — Comprehensive Migration Plan*