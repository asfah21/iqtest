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
| **Domain** | Personality & Cognitive Assessment |
| **Core Test** | MBTI-Based Personality Assessment (20-question variant) |
| **Framework** | Go (Gin) + PostgreSQL + templ |
| **Target Audience** | General public, individuals seeking self-understanding |
| **Monetization** | Freemium — free assessment, paid full results (IDR 14,900) |

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
| Frontend | Alpine.js | Client-side interactivity (quiz, navigation) |
| Styling | Custom CSS | Design system based on DESIGN.md |
| Containerization | Docker | Development & deployment consistency |

---

## 2. ASSESSMENT METHODOLOGY

### 2.1 Assessment Type

This system implements a **personality assessment** based on the **Jungian Cognitive Function theory**, operationalized through the **Myers-Briggs Type Indicator (MBTI)** framework. While the user interface advertises "IQ-like" or "personality insight", the underlying engine measures:

- **4 Dichotomies**: Extraversion/Introversion (E/I), Sensing/Intuition (S/N), Thinking/Feeling (T/F), Judging/Perceiving (J/P)
- **Cognitive Function Stack**: Derivation of 8 Jungian cognitive functions from the 4-letter type
- **Dark Triad Correlation**: Mapping of MBTI raw scores to Narcissism, Machiavellianism, and Psychopathy dimensions for narrative generation

### 2.2 Assessment Principles

| Principle | Implementation |
|-----------|---------------|
| **Self-Report** | User responds to Likert-scale questions about their own preferences |
| **Forced-Preference** | 6-point Likert scale omits neutral midpoint, forcing directional choice |
| **Weighted Scoring** | Questions have varying weights (1.5–2.0) based on discriminative power |
| **Multi-Dimensional** | Each dichotomy independently scored; type derived from combination |
| **Reliability Check** | Completion rate, inconsistency detection, response time analysis |

### 2.3 Test Length & Duration

| Metric | Value |
|--------|-------|
| **Total Questions** | 20 |
| **Estimated Completion Time** | 5–10 minutes |
| **Time Limit** | None (unlimited) |
| **Break Policy** | Not supported (single session) |
| **Retake Policy** | Not implemented (future improvement) |

---

## 3. PSYCHOMETRIC FOUNDATION

### 3.1 Theoretical Basis

The assessment is grounded in **Carl Jung's theory of psychological types** (1921), later developed into the MBTI by Isabel Briggs Myers and Katharine Cook Briggs. The core concepts:

```
Jungian Cognitive Functions:
  ┌─────────────────────────────────────────────────────┐
  │                    Perceiving (S/N)                 │
  │  ┌──────────────┐              ┌──────────────┐    │
  │  │ Se (Ext.     │              │ Ne (Ext.     │    │
  │  │  Sensing)    │              │  Intuition)  │    │
  │  └──────┬───────┘              └──────┬───────┘    │
  │         │                              │            │
  │  ┌──────▼───────┐              ┌──────▼───────┐    │
  │  │ Si (Int.     │              │ Ni (Int.     │    │
  │  │  Sensing)    │              │  Intuition)  │    │
  │  └──────────────┘              └──────────────┘    │
  ├─────────────────────────────────────────────────────┤
  │                    Judging (T/F)                    │
  │  ┌──────────────┐              ┌──────────────┐    │
  │  │ Te (Ext.     │              │ Fe (Ext.     │    │
  │  │  Thinking)   │              │  Feeling)    │    │
  │  └──────┬───────┘              └──────┬───────┘    │
  │         │                              │            │
  │  ┌──────▼───────┐              ┌──────▼───────┐    │
  │  │ Ti (Int.     │              │ Fi (Int.     │    │
  │  │  Thinking)   │              │  Feeling)    │    │
  │  └──────────────┘              └──────────────┘    │
  └─────────────────────────────────────────────────────┘
```

### 3.2 The 4 Dichotomies

| Dichotomy | Pole A | Pole B | Psychological Construct |
|-----------|--------|--------|------------------------|
| **E/I** | Extraversion | Introversion | Direction of psychic energy (outer vs inner world) |
| **S/N** | Sensing | Intuition | Information gathering preference (concrete vs abstract) |
| **T/F** | Thinking | Feeling | Decision-making preference (logic vs values) |
| **J/P** | Judging | Perceiving | Lifestyle orientation (structured vs flexible) |

### 3.3 Preference Clarity Index (PCI)

The PCI measures how strongly a user prefers one pole over another on each dichotomy:

```
PCI Formula:
  PCI = |raw_score| / max_possible_score × 100

  Where:
    raw_score = pole_a_score - pole_b_score
    max_possible_score = Σ(weight_i) for all questions in that dichotomy
```

| PCI Range | Label | Interpretation |
|-----------|-------|----------------|
| 0% – 25% | Slight | Weak preference, flexible |
| 26% – 50% | Moderate | Moderate preference |
| 51% – 75% | Clear | Strong preference |
| 76% – 100% | Very Clear | Very strong preference |

### 3.4 Cognitive Function Stack Derivation

Each 4-letter MBTI type maps to a unique stack of 4 cognitive functions:

```
Rules:
  1. Dominant + Inferior are from the same axis (e.g., Se↔Ni)
  2. Auxiliary + Tertiary are from the same axis
  3. Stack contains both Perceiving and Judging functions
  4. Stack contains both Extraverted and Introverted functions
```

**Derivation Algorithm (simplified):**

```
Given: [E/I] [S/N] [T/F] [J/P]

If E:
  If J:  Dominant = Te|Fe,  Auxiliary = Ni|Si
  If P:  Dominant = Se|Ne,  Auxiliary = Ti|Fi
If I:
  If J:  Dominant = Ti|Fi,  Auxiliary = Se|Ne
  If P:  Dominant = Si|Ni,  Auxiliary = Te|Fe
```

### 3.5 Reliability Indicators

| Indicator | Calculation | Threshold |
|-----------|-------------|-----------|
| **Completion Rate** | (answered / total) × 100 | ≥ 90% for valid result |
| **Avg Response Time** | Σ(time_per_question) / count | ≥ 1.5s to detect speed-clicking |
| **Inconsistency Score** | Deviation across anchor pairs | ≤ 35 for reliable result |

---

## 4. QUESTION STRUCTURE

### 4.1 Question Bank Composition

The current question bank contains **20 questions** distributed across 4 dichotomies:

| Dichotomy | Count | Weight Distribution |
|-----------|-------|-------------------|
| E/I | 5 | 2× weight 2.0, 3× weight 1.5 |
| S/N | 6 | 2× weight 2.0, 4× weight 1.5 |
| T/F | 5 | 2× weight 2.0, 3× weight 1.5 |
| J/P | 4 | 2× weight 2.0, 2× weight 1.5 |
| **Total** | **20** | **Max score per dichotomy varies** |

### 4.2 Question Metadata Schema

```go
// questionDef — metadata for each MBTI question
type questionDef struct {
    ID            string  // e.g., "Q_EI_001"
    Dikotomi      string  // "EI" | "SN" | "TF" | "JP"
    PolePrimary   string  // "E"|"I"|"S"|"N"|"T"|"F"|"J"|"P"
    Weight        float64 // 1.5 or 2.0
    ReverseScored bool    // whether scoring direction is inverted
}
```

### 4.3 Question ID Naming Convention

```
Q_{DIKOTOMI}_{SEQUENCE}

Examples:
  Q_EI_001  — Extraversion/Introversion, question #1
  Q_SN_003  — Sensing/Intuition, question #3
  Q_TF_005  — Thinking/Feeling, question #5
  Q_JP_002  — Judging/Perceiving, question #2
```

### 4.4 Response Format

All questions use a **6-point Likert scale**:

| Value | Meaning | Contribution |
|-------|---------|-------------|
| 1 | Strongly agree with the statement | 100% toward Pole A |
| 2 | Agree with the statement | 67% toward Pole A |
| 3 | Slightly agree with the statement | 33% toward Pole A |
| 4 | Slightly disagree with the statement | 33% toward Pole B |
| 5 | Disagree with the statement | 67% toward Pole B |
| 6 | Strongly disagree with the statement | 100% toward Pole B |

**For Reverse-Scored Questions (3 of 20):**
The scale is inverted: `adjusted = 7 - raw`

### 4.5 Likert Contribution Mapping

```go
var likertContribution = map[int]float64{
    1: 1.00,  // Full contribution to pole_primary
    2: 0.67,  // Strong contribution to pole_primary
    3: 0.33,  // Weak contribution to pole_primary
    4: 0.33,  // Weak contribution to pole_opposite
    5: 0.67,  // Strong contribution to pole_opposite
    6: 1.00,  // Full contribution to pole_opposite
}
```

### 4.6 Question Format (UI)

Each question is rendered in the frontend as:

```
┌─────────────────────────────────────────────┐
│   Pertanyaan 3 dari 20                       │
│                                             │
│   [Soal teks appears here]                   │
│                                             │
│   1   2   3   4   5   6                      │
│   ○───○───○───○───○───○                     │
│   │   │   │   │   │   │                     │
│  SS   S   AS  AD   D   SD                   │
│                                             │
│        ◀ Sebelumnya    Selanjutnya ▶        │
└─────────────────────────────────────────────┘

SS = Sangat Setuju
S  = Setuju
AS = Agak Setuju
AD = Agak Tidak Setuju
D  = Tidak Setuju
SD = Sangat Tidak Setuju
```

---

## 5. TIMER RULES

### 5.1 Current Implementation

The current system does **not** implement a hard time limit per question or per test. However, the architecture supports response time tracking for anti-cheating analysis.

### 5.2 Timer Architecture (Planned)

| Feature | Specification |
|---------|--------------|
| **Total Test Timer** | Optional countdown (future improvement) |
| **Per-Question Timer** | Not implemented (planned for v2) |
| **Warning Threshold** | N/A |
| **Auto-Submit** | Not implemented |
| **Pause/Resume** | Not supported |

### 5.3 Response Time Tracking (Current)

The frontend captures `time_taken_ms` per question via Alpine.js:

```javascript
// timestamp when question first rendered
const questionStartTime = Date.now();

// when user selects an answer, capture elapsed time
onAnswerSelect: function() {
    const elapsed = Date.now() - questionStartTime;
    // elapsed is stored per-question in the response payload
}
```

### 5.4 Recommended Timer Rules (for Future Implementation)

| Rule | Value | Rationale |
|------|-------|-----------|
| **Max duration per question** | 60 seconds | Prevents indecision looping |
| **Max total duration** | 20 minutes | 20 questions × 60s = 20 min ceiling |
| **Warning at** | 5 minutes remaining | Visual countdown in navbar |
| **Auto-submit at** | 0 minutes remaining | Force submit current answers |
| **Speed flag threshold** | < 1.5s avg per question | Flags potential random clicking |
| **Dwell time minimum** | 0.5s per question | Ignores accidental clicks |

---

## 6. SCORING ALGORITHM

### 6.1 Scoring Pipeline

```
┌──────────┐    ┌────────────┐    ┌───────────┐    ┌──────────────┐
│ Raw      │───▶│ Likert     │───▶│ Weighted  │───▶│ Accumulate   │
│ Answer   │    │ Adjustment │    │ Score     │    │ per Pole     │
│ (1-6)    │    │ + Reverse  │    │           │    │              │
└──────────┘    └────────────┘    └───────────┘    └──────┬───────┘
                                                           │
┌──────────┐    ┌────────────┐    ┌───────────┐           │
│ 4-Letter │◀───│ Preference │◀───│ Raw Score │◀──────────┘
│ Type     │    │ Selection  │    │ (A - B)   │
└─────┬────┘    └────────────┘    └───────────┘
      │
      ▼
┌──────────────┐
│ Cognitive    │
│ Stack Derive │
└──────────────┘
```

### 6.2 Step-by-Step Algorithm

#### Step 1: Initialize Accumulators

```go
accumulators := map[string]*acc{
    "EI": {poleA: 0, poleB: 0, max: 0},
    "SN": {poleA: 0, poleB: 0, max: 0},
    "TF": {poleA: 0, poleB: 0, max: 0},
    "JP": {poleA: 0, poleB: 0, max: 0},
}
```

For each dichotomy:
- **poleA** = first letter (E, S, T, J)
- **poleB** = second letter (I, N, F, P)
- **max** = sum of weights for questions answered in that dichotomy

#### Step 2: Process Each Answer

For each answered question `q`:

```
1. Look up question definition (dikotomi, pole_primary, weight, reverse_scored)

2. Apply reverse scoring if needed:
   adjusted = (q.ReverseScored) ? 7 - raw_value : raw_value

3. Get Likert contribution:
   contribution = likertContribution[adjusted]

4. Calculate weighted score:
   weighted = contribution × q.weight

5. Determine direction:
   If adjusted <= 3 → contribution goes to pole_primary
   If adjusted > 3  → contribution goes to pole_opposite

6. Accumulate:
   If answer favors poleA → acc.poleA += weighted
   If answer favors poleB → acc.poleB += weighted
   acc.max += q.weight
```

#### Step 3: Calculate Raw Scores

```go
rawScore_EI = accumulators["EI"].poleA - accumulators["EI"].poleB
rawScore_SN = accumulators["SN"].poleA - accumulators["SN"].poleB
rawScore_TF = accumulators["TF"].poleA - accumulators["TF"].poleB
rawScore_JP = accumulators["JP"].poleA - accumulators["JP"].poleB
```

**Interpretation:**
- Positive raw score → Preference for Pole A (E, S, T, J)
- Negative raw score → Preference for Pole B (I, N, F, P)

#### Step 4: Derive 4-Letter Type

```go
mbtiType := ""
mbtiType += (rawScore_EI >= 0) ? "E" : "I"
mbtiType += (rawScore_SN >= 0) ? "S" : "N"
mbtiType += (rawScore_TF >= 0) ? "T" : "F"
mbtiType += (rawScore_JP >= 0) ? "J" : "P"
// Example output: "INTJ"
```

#### Step 5: Calculate PCI for Each Dichotomy

```go
pci = |rawScore| / maxPossible × 100
pci = round(pci × 10) / 10  // Round to 1 decimal place
```

#### Step 6: Determine Strength Label

```go
strength := "very_clear"
switch {
case pci <= 25: strength = "slight"
case pci <= 50: strength = "moderate"
case pci <= 75: strength = "clear"
}
```

#### Step 7: Derive Cognitive Stack

Using the `DeriveCognitiveStack` function that maps 4-letter type to:
- Dominant function
- Auxiliary function
- Tertiary function
- Inferior function

### 6.3 Scoring Example

| Question | Answer | Adjusted | Contribution | Weight | Weighted | Direction |
|----------|--------|----------|--------------|--------|----------|-----------|
| Q_EI_001 | 2 | 2 | 0.67 | 2.0 | 1.34 | E |
| Q_EI_002 | 5 | 5 | 0.67 | 2.0 | 1.34 | I |
| Q_EI_003 | 1 | 1 | 1.00 | 1.5 | 1.50 | E |
| Q_EI_004 | 3 | 3 | 0.33 | 1.5 | 0.50 | I |
| Q_EI_005 (R) | 6 | 1 | 1.00 | 1.5 | 1.50 | E |

**Result:**
- E total = 1.34 + 1.50 + 1.50 = 4.34
- I total = 1.34 + 0.50 = 1.84
- Raw Score = 4.34 - 1.84 = 2.50 (positive → E preference)
- Max possible = 2.0 + 2.0 + 1.5 + 1.5 + 1.5 = 8.5
- PCI = |2.50| / 8.5 × 100 = 29.4% → "moderate"

---

## 7. IQ SCORE CONVERSION

### 7.1 Current Status

**Important Note:** This system does **not** measure IQ. The "IQ" label is **not applicable** to this personality assessment engine. The results produce MBTI types, cognitive function stacks, and Dark Triad correlation narratives — not IQ scores.

### 7.2 Planned IQ Module (Future Enhancement)

If IQ-style scoring is desired in the future, the following framework is proposed:

#### 7.2.1 IQ Scoring Model

```
IQ Score = 100 + (z_score × 15)

Where:
  z_score = (user_score - population_mean) / population_std_dev
  15 = standard deviation for typical IQ scales
```

#### 7.2.2 Conversion Table (Illustrative)

| Raw Score Range | Z-Score | IQ Equivalent | Classification |
|-----------------|---------|---------------|----------------|
| > 1.5σ above mean | > 1.5 | > 122 | Superior |
| 0.5σ to 1.5σ above mean | 0.5 to 1.5 | 107–122 | High Average |
| ±0.5σ of mean | -0.5 to 0.5 | 92–107 | Average |
| 0.5σ to 1.5σ below mean | -1.5 to -0.5 | 77–92 | Low Average |
| > 1.5σ below mean | < -1.5 | < 77 | Below Average |

#### 7.2.3 Normative Data Requirements

To implement IQ conversion, the following data must be collected:

1. **Population mean** — calculated from minimum 1,000 completed assessments
2. **Population standard deviation** — calculated from the same cohort
3. **Demographic stratification** — age, education level, gender (optional)
4. **Regular re-norming** — every 6–12 months to maintain validity

---

## 8. RESULT INTERPRETATION

### 8.1 Result Delivery Flow

```
User completes 20 questions
            │
            ▼
    POST /submit-tes
            │
            ▼
    Server calculates MBTI
    Stores raw scores
            │
            ▼
    Returns user_id
            │
            ▼
    Redirect to /paywall/{id}
            │
            ▼
    User pays IDR 14,900
            │
            ▼
    POST /konfirmasi-bayar/{id}
    (by admin confirmation)
            │
            ▼
    GET /hasil/{id}
            │
            ▼
    Generates narratives
    via GenerateAllNarratives()
            │
            ▼
    Renders hasil_page.templ
```

### 8.2 Result Components

The final results page displays the following sections:

| Section | Source | Description |
|---------|--------|-------------|
| **MBTI Type** | Scoring Algorithm | 4-letter type (e.g., INTJ) |
| **Cognitive Stack** | Cognitive Function Engine | 4 functions with descriptions |
| **Executive Summary** | Narrative Generator | Personalized introduction |
| **Relationship Profile** | Narrative Generator | 5-axis relationship analysis |
| **Kekuatan (Strengths)** | Narrative Generator | 3–5 strength items |
| **Area Perhatian (Concerns)** | Narrative Generator | 3–5 development areas |
| **Relationship Insight** | Narrative Generator | Pattern identification |
| **Compatibility Notes** | Narrative Generator | Partner type recommendations |
| **Reflection Questions** | Narrative Generator | 3 guided introspection questions |

### 8.3 Narrative Generation Engine

The narrative generator (`services/narasi.go`) uses the **Dark Triad mapping** of MBTI raw scores:

```go
func mapMBTIToDarkTriad(skorEI, skorSN, skorTF, skorJP int) (narsisme, machiavellian, psikopati int) {
    narsisme = absInt(skorEI) * 5      // E/I → Narcissism
    machiavellian = absInt(skorSN) * 5 // S/N → Machiavellianism
    psikopati = absInt(skorTF) * 5     // T/F → Psychopathy
}
```

**Thresholds applied:**
- Values capped at 100 (score × 5, maxed at 100)
- Labels: ≤25 = rendah, ≤50 = sedang, ≤75 = cukup tinggi, >75 = tinggi

### 8.4 Score Visualization

Each dichotomy's PCI is visualized as a progress bar:

```
E ────────────────────────────── I
       ████████████████
       PCI: 32.2% (moderate)
       Preference: I

S ────────────────────────────── N
          ████████████████████████
          PCI: 37.5% (moderate)
          Preference: N

T ────────────────────────────── F
    ██████████████████████
    PCI: 45.8% (moderate)
    Preference: T

J ────────────────────────────── P
               ██████████
               PCI: 22.2% (slight)
               Preference: J
```

### 8.5 Cognitive Stack Display

```
┌────────────────────────────────────────────────┐
│  Cognitive Function Stack                       │
│                                                │
│  Dominant    │  Ni  │ Introverted iNtuition    │
│  ────────────────────────────────────────────── │
│  Auxiliary   │  Te  │ Extraverted Thinking     │
│  ────────────────────────────────────────────── │
│  Tertiary    │  Fi  │ Introverted Feeling      │
│  ────────────────────────────────────────────── │
│  Inferior    │  Se  │ Extraverted Sensing      │
│                                                │
│  Temperament: NT (Rational)                    │
│  Nickname: The Architect (Sang Arsitek)        │
└────────────────────────────────────────────────┘
```

---

## 9. ANTI-CHEATING STRATEGY

### 9.1 Current Protections

| Strategy | Implementation | Status |
|----------|---------------|--------|
| No neutral midpoint | 6-point Likert scale (no 3.5) | ✅ Active |
| Reverse-coded questions | 3 of 20 questions are reverse-scored | ✅ Active |
| Weighted scoring | Variable weights prevent trivial gaming | ✅ Active |
| Server-side validation | All scoring is server-side | ✅ Active |
| Paywall protection | Results locked until payment confirmed | ✅ Active |

### 9.2 Planned Detections (Not Yet Implemented)

| Detection | Method | Threshold | Action |
|-----------|--------|-----------|--------|
| **Response time analysis** | Track ms per question | Avg < 1.5s → flag | Warning in results |
| **Straight-line detection** | Check answer variance | Variance ≤ 0.5 → flag | Invalidate session |
| **Inconsistency scoring** | Compare anchor pairs | Score > 35 → warn | Review suggested |
| **Completion rate check** | Count answered / total | < 80% → reject | Force completion |
| **Session replay** | Check for duplicate interactions | Multiple entries → flag | Block with CAPTCHA |
| **IP rate limiting** | Count submissions per IP | > 5 per hour → block | 429 Too Many Requests |

### 9.3 Anchor Pairs for Inconsistency Detection

Anchor pairs are pairs of questions that should theoretically yield consistent answers:

```
Anchor Pair Types:
  "same"     — Both should trend toward the same pole
  "opposite" — Should trend toward opposite poles
```

**Example Anchor Pair:**

```go
AnchorPair{
    QuestionAID: "Q_EI_001", // Measures E preference
    QuestionBID: "Q_EI_003", // Also measures E preference
    Expected:    "same",     // Both should be E or both I
}
```

**Inconsistency Calculation:**

```go
// For "same" pairs:
inconsistent := (responseA < 3.5) != (responseB < 3.5)
// If inconsistent, magnitude = |responseA - responseB| × 10

// For "opposite" pairs:
consistent := (responseA < 3.5) != (responseB < 3.5)
// If inconsistent, magnitude = |responseA - (7 - responseB)| × 10
```

### 9.4 Confidence Score

The overall confidence score combines multiple reliability indicators:

```go
type ConfidenceScore struct {
    Overall        float64           // 0–100
    PerDikotomi    map[string]float64 // PCI per dichotomy
    Flags          []string          // Warning flags
    Recommendation string            // "reliable" | "review_suggested" | "retest_recommended"
}
```

**Scoring penalties:**
- PCI < 5 (near-zero preference): −25 per dichotomy
- PCI < 15 (slight preference): −10 per dichotomy
- Inconsistency > 50: −30
- Inconsistency > 30: −15
- Completion < 90%: −20
- Avg response time < 1.5s: −20

**Recommendation thresholds:**
- ≥ 70: "reliable"
- 40–69: "review_suggested"
- < 40: "retest_recommended"

### 9.5 Frontend Anti-Cheat Measures

```
┌─────────────────────────────────────────────────────┐
│ Anti-Cheat Layer (Frontend)                         │
├─────────────────────────────────────────────────────┤
│ ✅ Disable right-click context menu (optional)      │
│ ✅ Disable text selection on quiz page              │
│ ✅ Prevent copy-paste in answer areas               │
│ ✅ Track tab-switching events (visibility API)      │
│ ✅ Disable browser back navigation during quiz      │
│ ✅ Warn on page refresh attempt                     │
│ ✅ Capture viewport focus/blur timestamps            │
└─────────────────────────────────────────────────────┘
```

```javascript
// Tab switch detection
document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
        logTabSwitch(Date.now());
    }
});
```

---

## 10. DATABASE MODEL

### 10.1 Current Schema (Simplified)

```
users_test
┌──────────────┬─────────────┬──────────────────────────┐
│ Column       │ Type        │ Description              │
├──────────────┼─────────────┼──────────────────────────┤
│ id           │ SERIAL/UUID │ Primary key              │
│ nama         │ VARCHAR     │ User name                │
│ email        │ VARCHAR     │ User email               │
│ skor_ei      │ INTEGER     │ Raw E/I score            │
│ skor_sn      │ INTEGER     │ Raw S/N score            │
│ skor_tf      │ INTEGER     │ Raw T/F score            │
│ skor_jp      │ INTEGER     │ Raw J/P score            │
│ mbti_tipe    │ VARCHAR(4)  │ 4-letter MBTI type       │
│ status_pembayaran │ VARCHAR │ 'UNPAID' or 'PAID'      │
└──────────────┴─────────────┴──────────────────────────┘
```

### 10.2 Current DDL (PostgreSQL)

```sql
CREATE TABLE users_test (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) NOT NULL,
    skor_ei           INTEGER NOT NULL,
    skor_sn           INTEGER NOT NULL,
    skor_tf           INTEGER NOT NULL,
    skor_jp           INTEGER NOT NULL,
    mbti_tipe         VARCHAR(4) NOT NULL,
    status_pembayaran VARCHAR(10) DEFAULT 'UNPAID',
    created_at        TIMESTAMPTZ DEFAULT NOW()
);
```

### 10.3 Recommended Production Schema (Target)

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
-- Question responses (per question)
-- =============================================
CREATE TABLE session_responses (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id        UUID NOT NULL REFERENCES test_sessions(id),
    question_id       UUID NOT NULL REFERENCES questions(id),
    answer_value      DECIMAL(4,2) NOT NULL,
    answered_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    time_taken_ms     INTEGER,
    UNIQUE(session_id, question_id)
);

-- =============================================
-- Questions bank
-- =============================================
CREATE TABLE questions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_code     VARCHAR(20) UNIQUE NOT NULL,
    dikotomi          VARCHAR(2) NOT NULL CHECK (dikotomi IN ('EI','SN','TF','JP')),
    pole_primary      VARCHAR(1) NOT NULL CHECK (pole_primary IN ('E','I','S','N','T','F','J','P')),
    weight            DECIMAL(3,1) NOT NULL DEFAULT 1.0,
    format            VARCHAR(20) NOT NULL DEFAULT 'likert_6',
    reverse_scored    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    translations      JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- MBTI results
-- =============================================
CREATE TABLE mbti_results (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id          UUID NOT NULL REFERENCES test_sessions(id) UNIQUE,
    mbti_type           VARCHAR(4) NOT NULL,
    ei_raw_score        DECIMAL(8,2) NOT NULL,
    sn_raw_score        DECIMAL(8,2) NOT NULL,
    tf_raw_score        DECIMAL(8,2) NOT NULL,
    jp_raw_score        DECIMAL(8,2) NOT NULL,
    ei_pci              DECIMAL(5,1) NOT NULL,
    sn_pci              DECIMAL(5,1) NOT NULL,
    tf_pci              DECIMAL(5,1) NOT NULL,
    jp_pci              DECIMAL(5,1) NOT NULL,
    cognitive_stack     JSONB NOT NULL,
    completion_rate     DECIMAL(5,1) NOT NULL,
    avg_response_ms     INTEGER,
    inconsistency_score DECIMAL(5,1) DEFAULT 0,
    is_reliable         BOOLEAN NOT NULL DEFAULT TRUE,
    calculated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
CREATE INDEX idx_results_session ON mbti_results(session_id);
CREATE INDEX idx_results_type ON mbti_results(mbti_type);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
```

### 10.4 Entity Relationship Diagram

```
users
  │
  ├──< test_sessions
  │      │
  │      ├──< session_responses
  │      │
  │      └──> mbti_results
  │
  └──< payments

admins ──> payments (confirmed_by)

questions ──> questions_translations (via JSONB or junction table)
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
     │◀──────────────│  QuizPage      │                │
     │               │                │                │
     │  POST /submit-tes              │                │
     │  {email, nama, answers...}     │                │
     │──────────────▶│                │                │
     │               │  ProcessQuizAnswers()           │
     │               │───────────────▶│                │
     │               │                │  CalculateMBTI │
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

### 11.2 Route Table

| Method | Path | Handler | Auth | Description |
|--------|------|---------|------|-------------|
| GET | `/` | ShowHome | None | Landing page |
| GET | `/quiz` | ShowQuiz | None | Assessment page |
| POST | `/submit-tes` | SubmitTest | None | Submit answers |
| GET | `/paywall/:id` | ShowPaywall | None | Payment gate |
| POST | `/konfirmasi-bayar/:id` | KonfirmasiBayar | None | Payment confirm |
| GET | `/hasil/:id` | ShowResult | None | View results (PAID only) |
| GET | `/tentang` | ShowTentang | None | About page |
| GET | `/admin/login` | ShowLogin | None | Admin login form |
| POST | `/admin/login` | LoginProcess | None | Admin login action |
| GET | `/admin/dashboard` | ShowDashboard | Admin cookie | Admin panel |
| GET | `/admin/user/:id` | ShowUserDetail | Admin cookie | User detail |
| GET | `/admin/logout` | LogoutProcess | Admin cookie | Logout |

### 11.3 Response Types

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

#### Payment Confirmation Response

```json
{
    "success": true,
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 11.4 HTTP Status Codes

| Code | Condition |
|------|-----------|
| 200 | Success |
| 303 | Redirect (payment / paywall) |
| 400 | Bad request (invalid form data) |
| 404 | User/session not found |
| 500 | Internal server error |

---

## 12. UI/UX FLOW

### 12.1 User Journey Map

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Landing  │───▶│  Quiz    │───▶│ Paywall  │───▶│ Results  │
│  Page    │    │  (20 Q)  │    │ (Payment)│    │  (PAID)  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
      │              │               │               │
      │              │               │               │
      ▼              ▼               ▼               ▼
  IndexPage     QuizPage       PaywallPage      HasilPage
  (Static)      (Alpine.js)    (Static)         (Narratives)
```

### 12.2 Page Descriptions

#### 12.2.1 Landing Page (`GET /`)

| Element | Description |
|---------|-------------|
| **Hero** | Tagline: "Kenali Dirimu Lebih Dalam" with DM Serif Display |
| **CTA** | "Mulai Tes Gratis" — primary button linking to `/quiz` |
| **Trust Pills** | "Anonim · Gratis · Hasil Instan" |
| **Features** | 3-card section explaining the assessment |
| **FAQ** | Accordion-style frequently asked questions |
| **Footer** | Links, copyright, brand info |

#### 12.2.2 Quiz Page (`GET /quiz`)

| Element | Description |
|---------|-------------|
| **Progress Bar** | Top of page, shows completion (X/20) |
| **Question Counter** | "Pertanyaan N dari 20" |
| **Question Text** | Dynamic, loaded via Alpine.js |
| **Likert Scale** | 6 radio buttons (1–6) with labels |
| **Navigation** | "Sebelumnya" / "Selanjutnya" buttons |
| **Submit Button** | Appears only on question 20 |
| **Form Fields** | Email and Name (shown before quiz starts) |

**Alpine.js State Management:**

```javascript
Alpine.data('quizApp', () => ({
    // State
    step: 'identity',       // 'identity' | 'quiz' | 'submitting' | 'done'
    nama: '',
    email: '',
    currentQuestion: 0,
    answers: {},            // { questionIndex: answerValue }
    questions: [...],       // Array of 20 question objects
    
    // Computed
    get progress() {
        return Object.keys(this.answers).length;
    },
    
    // Methods
    startQuiz() { ... },
    selectAnswer(index, value) { ... },
    nextQuestion() { ... },
    prevQuestion() { ... },
    submitQuiz() { ... }
}));
```

#### 12.2.3 Paywall Page (`GET /paywall/:id`)

| Element | Description |
|---------|-------------|
| **Greeting** | "Halo, {nama}!" |
| **Value Prop** | Explanation of premium results |
| **Pricing** | IDR 14.900 (one-time payment) |
| **Payment Instructions** | Manual transfer to bank account |
| **Confirmation Button** | "Saya sudah bayar" — sends POST to `/konfirmasi-bayar/:id` |
| **Error State** | "belum_bayar" query param → information message |

#### 12.2.4 Results Page (`GET /hasil/:id`)

| Section | Content |
|---------|---------|
| **Header** | "Hasil Asesmen {nama}" |
| **MBTI Badge** | 4-letter type with temperament color |
| **Cognitive Stack** | Dominant, Auxiliary, Tertiary, Inferior |
| **Executive Summary** | AI-generated personalized narrative |
| **Relationship Profile** | 5-axis deep analysis |
| **Strengths** | 3–5 bullet points |
| **Growth Areas** | 3–5 bullet points |
| **Compatibility** | Partner type analysis |
| **Reflection Questions** | 3 guided questions |
| **Share/Print** | Action buttons (future) |

#### 12.2.5 Admin Dashboard (`GET /admin/dashboard`)

| Element | Description |
|---------|-------------|
| **Statistics** | Total users, paid/unpaid counts, total revenue |
| **User Table** | ID, Name, Email, MBTI, Payment Status, Scores |
| **Search/Filter** | By name, email, MBTI type |
| **User Detail** | Link to `/admin/user/:id` |
| **Logout** | Clear session cookie |

### 12.3 Design System Integration

All pages follow the design system defined in `DESIGN.md`:

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
| Mobile | < 640px | Single column, stacked |
| Tablet | 640–1024px | 2-column grid on cards |
| Desktop | > 1024px | Full layout with max-width container |

### 12.5 Accessibility Requirements

- WCAG AA minimum on all text elements
- Focus-visible outlines on interactive elements
- `prefers-reduced-motion` respected for all animations
- Semantic HTML structure (nav, main, section, footer)
- ARIA labels on all interactive components
- Color is never the sole differentiator (paired with labels)

---

## 13. FUTURE IMPROVEMENTS

### 13.1 Scoring & Algorithm

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Dynamic question bank** | High | Store questions in DB, support A/B testing | 2 weeks |
| **Forced-choice format** | Medium | Add alternate question format for higher discrimination | 1 week |
| **Cognitive function scoring** | Medium | Score each of 8 cognitive functions directly (not just dichotomies) | 3 weeks |
| **Anchor pair inconsistency** | Medium | Implement anchor pairs for reliability scoring | 1 week |
| **Confidence scoring** | Medium | Implement `CalculateConfidence` for result reliability | 3 days |
| **IQ-style scoring module** | Low | Add normative data-based IQ conversion (requires population data) | 2 weeks |

### 13.2 Question Bank

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Expand to 40+ questions** | High | Larger pool for randomized tests | 2 weeks |
| **Multilingual support** | Medium | English, Japanese, Bahasa Indonesia | 3 weeks |
| **Question categories** | Medium | Tag by domain (work, social, cognitive, emotional) | 1 week |
| **Adaptive testing** | Low | CAT-style dynamic question selection based on prior answers | 4 weeks |

### 13.3 Timer & Anti-Cheating

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Per-question timer** | High | 60s countdown per question with auto-advance | 1 week |
| **Total test timer** | Medium | 20-minute overall timer with warning at 5 min | 3 days |
| **Tab-switch detection** | Medium | Log and flag tab switches during test | 2 days |
| **Straight-line detection** | Medium | Detect all-same-answer patterns | 2 days |
| **IP rate limiting** | Medium | Prevent multiple submissions from same IP | 2 days |
| **CAPTCHA integration** | Low | Google reCAPTCHA v3 on submit | 2 days |

### 13.4 Payment & Monetization

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Automated payment gateway** | High | Midtrans/Xendit integration for instant payment | 4 weeks |
| **Multiple price tiers** | Medium | Basic (free) / Premium (detailed) / Pro (consultation) | 2 weeks |
| **Discount codes** | Medium | Admin-managed promo codes | 1 week |
| **Group/team pricing** | Low | Bulk purchase for organizations | 1 week |

### 13.5 User Experience

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Email delivery of results** | High | Send PDF result to user email after payment | 1 week |
| **Progress save/restore** | Medium | Save partial progress and allow resume | 2 weeks |
| **Social sharing** | Medium | Share results on social media (link only, no raw scores) | 3 days |
| **Comparison tool** | Low | Compare results with friends/partners | 2 weeks |
| **Result history** | Low | Track changes over time with retest capability | 2 weeks |

### 13.6 Admin & Operations

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Admin notification** | High | Notify admin on new payment confirmations | 3 days |
| **CSV export** | Medium | Export user data to CSV | 2 days |
| **Analytics dashboard** | Medium | Charts for MBTI type distribution, conversion rates | 2 weeks |
| **A/B test framework** | Low | Test different question wordings and formats | 3 weeks |

### 13.7 Technical Infrastructure

| Improvement | Priority | Description | Effort |
|-------------|----------|-------------|--------|
| **Redis caching** | Medium | Cache question bank, reduce DB load | 1 week |
| **CDN for assets** | Medium | Serve CSS/JS from CDN for faster load | 2 days |
| **CI/CD pipeline** | Medium | Automated testing and deployment | 1 week |
| **Load testing** | Medium | Benchmark for concurrent users (target: 10,000) | 1 week |
| **API rate limiting** | Medium | Per-IP and per-user rate limits | 3 days |
| **Docker compose for production** | Low | Multi-service setup (app, db, redis, nginx) | 3 days |

### 13.8 Recommended Roadmap

```
Phase 1 (Q3 2026) — Foundation
├── Dynamic question bank
├── Per-question timer
├── Anchor pair inconsistency detection
├── Confidence scoring

Phase 2 (Q4 2026) — Monetization
├── Automated payment gateway
├── Email delivery of results
├── Admin notification system
├── Discount codes

Phase 3 (Q1 2027) — Scale
├── Expand to 40+ questions
├── Multilingual support
├── Adaptive testing (beginner)
├── CSV export + analytics

Phase 4 (Q2 2027) — Advanced Features
├── Cognitive function direct scoring
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
| `handlers/admin.go` | Admin login, dashboard, user detail handlers |
| `helpers/render.go` | Templ component rendering helper |
| `models/user.go` | Data models (User, DikotomiScore, CognitiveStack, MBTIResult, etc.) |
| `repositories/user.go` | User data access (insert, query, update payment) |
| `repositories/admin.go` | Admin data access (list all users, get by ID) |
| `services/quiz.go` | Question bank, scoring algorithm, cognitive stack derivation |
| `services/narasi.go` | Narrative generation engine |
| `templ/components/` | Reusable UI components (head, navbar, footer) |
| `templ/layouts/` | Page layouts (public, quiz, auth, dashboard) |
| `templ/pages/` | Page templates (index, quiz, paywall, hasil, admin) |
| `templ/types/` | Type definitions for template data |
| `assets/css/` | Stylesheets |
| `assets/js/` | JavaScript (Alpine.js modules) |
| `assets/images/` | Static images |

### Appendix B — Scoring Reference Card

```
┌──────────────────────────────────────────────────────────────────┐
│                    SCORING REFERENCE CARD                        │
├────────────┬────────┬──────────┬────────┬─────────┬─────────────┤
│ Dichotomy  │ Pole A │ Pole B   │ Count  │ Weights │ Max Score   │
├────────────┼────────┼──────────┼────────┼─────────┼─────────────┤
│ E/I        │ E      │ I        │ 5      │ 2,2,1.5,│ 8.5         │
│            │        │          │        │ 1.5,1.5 │             │
├────────────┼────────┼──────────┼────────┼─────────┼─────────────┤
│ S/N        │ S      │ N        │ 6      │ 2,2,1.5,│ 10.0        │
│            │        │          │        │ 1.5,1.5,│             │
│            │        │          │        │ 1.5     │             │
├────────────┼────────┼──────────┼────────┼─────────┼─────────────┤
│ T/F        │ T      │ F        │ 5      │ 2,2,1.5,│ 8.5         │
│            │        │          │        │ 1.5,1.5 │             │
├────────────┼────────┼──────────┼────────┼─────────┼─────────────┤
│ J/P        │ J      │ P        │ 4      │ 2,2,1.5,│ 7.0         │
│            │        │          │        │ 1.5     │             │
├────────────┴────────┴──────────┴────────┴─────────┴─────────────┤
│ Total Questions: 20  │  Reverse-scored: 3 (Q_EI_005, Q_SN_005,  │
│                      │   Q_SN_006, Q_TF_005, Q_JP_004)         │
└──────────────────────────────────────────────────────────────────┘
```

### Appendix C — Cognitive Function Axis Map

```
Extraverted (E) ←──────────────────→ Introverted (I)

  Se (Extraverted Sensing)       ←→  Ni (Introverted Intuition)
  Si (Introverted Sensing)       ←→  Ne (Extraverted Intuition)
  Te (Extraverted Thinking)      ←→  Fi (Introverted Feeling)
  Ti (Introverted Thinking)      ←→  Fe (Extraverted Feeling)

Law of Opposites:
  Dominant function is always polar opposite of Inferior function
  Auxiliary function is always polar opposite of Tertiary function
```

### Appendix D — Glossary

| Term | Definition |
|------|------------|
| **Dichotomy** | A pair of opposing preferences (E/I, S/N, T/F, J/P) |
| **PCI** | Preference Clarity Index — measures strength of preference (0–100%) |
| **Cognitive Stack** | The hierarchy of 4 cognitive functions (Dominant → Auxiliary → Tertiary → Inferior) |
| **Temperament** | Keirsey temperament grouping: NT, NF, SJ, SP |
| **Likert Scale** | Rating scale measuring agreement level (1–6 in this system) |
| **Reverse-Scored** | Questions where the scale direction is inverted to prevent response bias |
| **Anchor Pair** | Two questions measuring the same construct, used for consistency checking |
| **Narrative Engine** | AI-like text generator that produces personalized analysis from scores |
| **Dark Triad** | Narcissism, Machiavellianism, and Psychopathy — used here as scoring framework for narratives |

### Appendix E — Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | HTTP server port | `8080` |
| `ADMIN_USERNAME` | Admin login username | `admin` |
| `ADMIN_PASSWORD` | Admin login password | `admin360` |
| `GIN_MODE` | Gin framework mode | `release` |

---

*End of QUIZ.md — Complete Specification Document*