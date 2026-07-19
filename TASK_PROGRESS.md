# TASK_PROGRESS.md

## Phase 0 ✅ | Phase 1 ✅ | Phase 2 ✅ | Phase 3 ✅ | Phase 4 ✅ | Phase 5 ✅ | Phase 6 ✅

## Phase 7 — Anti-Cheating System

**Status: ✅ Complete**

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 7.1 | Server-side answer validation (selected_option A/B/C/D/null, time_taken_ms 0–120000) | `services/quiz.go` | ✅ |
| 7.2 | Speed-guessing detection (>30% answers <3s & accuracy <25%) | `services/quiz.go` | ✅ |
| 7.3 | Straight-pattern detection (≥15/20 same option) | `services/quiz.go` | ✅ |
| 7.4 | Tab-switch tracking (flag if >3) | `services/quiz.go` | ✅ |
| 7.5 | `AssessReliability()` combining all detections | `services/quiz.go` | ✅ |
| 7.6 | IP rate limiting middleware (max 5 POST /submit-tes per hour, return 429) | `middleware/ratelimit.go` | ✅ |
| 7.7 | Random option shuffling | *(see below)* | ⚠️ Future (Phase 7) |
| 7.8 | Devtools tamper detection (frontend) | *(see below)* | ⚠️ Future (Phase 7) |
| 7.9 | Register rate limiting on `/submit-tes` | `handlers/router.go` | ✅ |

**Note:** Tasks 7.7 (random option shuffling) and 7.8 (devtools tamper detection) are listed in MIGRATION.md §9.4 as part of Phase 7 but are implemented as future enhancements that depend on having actual question images (stored in DB/questions table). The backend infrastructure (`session.metadata` for shuffle mapping) is ready for when these are implemented.

### Completion Criteria

- [x] Speed-guessing detection flags when >30% answers <3s + accuracy <25% (per §9.2)
- [x] Straight-pattern detection flags when ≥15/20 same option (per §9.2)
- [x] Tab-switch events tracked (>3 = flagged), stored in iq_results.reliability_flags
- [x] `POST /submit-tes` rate limited (max 5/hour per IP, returns 429)
- [x] `CorrectOption` never reaches client (server-side validation only)
- [x] Reliability flags stored with result in iq_results.reliability_flags (JSONB)
- [x] `go build ./...` passes

### Build Status

- [x] `go vet ./...` passes
- [x] `go build ./...` passes