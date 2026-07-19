# TASK_PROGRESS.md — Migration Progress

## Phase 0 — Rename & Rebrand (Required) ✅ COMPLETE

| # | Task | Status |
|---|------|--------|
| 0.1 | Update `models/user.go` — rename struct fields | ✅ Done |
| 0.2 | Update `services/quiz.go` — rename functions, vars, comments | ✅ Done |
| 0.3 | Update `services/quiz.go` — rename Dark Triad mapping function | ✅ Done |
| 0.4 | Update `services/narasi.go` — rename function signatures | ✅ Done |
| 0.5 | Update `handlers/quiz.go` — rename references | ✅ Done |
| 0.6 | Update `handlers/admin.go` — rename field references | ✅ Done |
| 0.7 | Update `templ/types/dashboard_data.go` — rename fields | ✅ Done |
| 0.8 | Update `repositories/user.go` — rename column references | ✅ Done |
| 0.9 | Update `repositories/admin.go` — rename column references | ✅ Done |
| 0.10 | Verify build succeeds (`go build ./...`) | ✅ Done |

### Files Modified (9 files)
- `models/user.go` — 6 structs renamed, 15 fields renamed
- `services/quiz.go` — 4 functions renamed, all question IDs/dikotomi/poles renamed, axisOpposites/axisOpposite deleted, DeriveCognitiveProfile rewritten
- `services/narasi.go` — 1 comment updated
- `handlers/quiz.go` — question IDs updated, comments updated
- `handlers/admin.go` — 5 field references updated
- `repositories/user.go` — 3 SQL column references updated
- `repositories/admin.go` — 3 SQL column references updated
- `templ/types/dashboard_data.go` — 5 fields renamed
- `templ/pages/dashboard_page.templ` — header label and field reference updated
- `templ/pages/user_detail_page.templ` — labels and field references updated
- Generated files: `templ generate` re-ran successfully

### Build Status
- `go build ./...` — ✅ PASSES

---

## Phase 1 — Scoring Engine Rewrite (Required) ✅ COMPLETE

**Note:** Most tasks in Phase 1 were already completed during Phase 0 (rename phase). Only Task 1.6 required explicit action.

| # | Task | Status | Notes |
|---|------|--------|-------|
| 1.1 | Rewrite `questionDef` definitions | ✅ Done | IDs, dikotomi, pole values all renamed to IQ Test (LR/NA/SA/LV). Done in Phase 0. |
| 1.2 | Update Likert contribution mapping | ✅ Done | Mapping is identical per IQTEST.md §4.5. No change needed. |
| 1.3 | Delete MBTI cognitive function logic | ✅ Done | `axisOpposites` map and `axisOpposite()` function removed. Done in Phase 0. |
| 1.4 | Rewrite `CalculateIQResult()` | ✅ Done | Accumulators changed to "LR"/"NA"/"SA"/"LV", pole A/B logic updated. Done in Phase 0. |
| 1.5 | Rewrite `DeriveCognitiveProfile()` | ✅ Done | Implemented per IQTEST.md §3.4. Done in Phase 0. |
| 1.6 | Rename local variable `pci` → `sci` | ✅ Done | Variable renamed in `buildDimensionScore()`. |
| 1.7 | Strength label logic | ✅ Done | Labels match IQTEST.md §3.3: slight(≤25), moderate(≤50), clear(≤75), very_clear(>75). Done in Phase 0. |
| 1.8 | Update `ProcessQuizAnswers()` | ✅ Done | Uses new `CalculateIQResult()`, stores skorLR/NA/SA/LV. Done in Phase 0. |
| 1.9 | Update `GetQuizResult()` | ✅ Done | Uses "LR"/"NA"/"SA"/"LV" dimension names and `DeriveCognitiveProfile()`. Done in Phase 0. |
| 1.10 | Update `mapIQToDarkTriad()` | ✅ Done | L/R→Narcissism, N/A→Machiavellianism, S/A→Psychopathy per IQTEST.md §8.3. Done in Phase 0. |

### Key Changes in `services/quiz.go`
- Question bank uses LR/NA/SA/LV dimensions with correct weights and reverse-scored flags
- `CalculateIQResult()` replaces `CalculateMBTI()` — produces 4-letter IQ types from dimension preferences
- `DeriveCognitiveProfile()` replaces `DeriveCognitiveStack()` — maps 4-letter type to Dominant/Auxiliary/Complementary/Developing
- `mapIQToDarkTriad()` mapping updated per IQTEST.md §8.3
- All MBTI cognitive function theory (`axisOpposites`, `axisOpposite`) removed
- Local variable `pci` renamed to `sci` in `buildDimensionScore()`

### Build Status
- `go build ./...` — ✅ PASSES

---

## Next Phases (Not Started)

| Phase | Type | Status |
|-------|------|--------|
| Phase 2 — Question Bank & Frontend | Required | ⏳ Not started |
| Phase 3 — Database Schema Normalization | Required | ⏳ Not started |
| Phase 4 — Repository Layer Rewrite | Required | ⏳ Not started |
| Phase 5 — Handler & Template Alignment | Required | ⏳ Not started |
| Phase 6 — Narrative Engine Update | Required | ⏳ Not started |
| Phase 7 — Anti-Cheating & Reliability | Future Enhancement | ⏳ Not started |
| Phase 8 — Production Schema & Payments | Future Enhancement | ⏳ Not started |
| Phase 9 — Admin Panel Update | Future Enhancement | ⏳ Not started |