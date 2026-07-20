package database

// initSchemaSQL berisi DDL untuk semua 7 tabel + indexes per IQTEST.md §10.2
const initSchemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) UNIQUE NOT NULL,
    nama              VARCHAR(255) NOT NULL,
    phone             VARCHAR(20),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS test_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID REFERENCES users(id),
    session_token     VARCHAR(64) UNIQUE NOT NULL,
    started_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at      TIMESTAMPTZ,
    device_type       VARCHAR(20),
    ip_address        VARCHAR(45),
    is_completed      BOOLEAN NOT NULL DEFAULT FALSE,
    metadata          JSONB
);

CREATE TABLE IF NOT EXISTS questions (
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

CREATE TABLE IF NOT EXISTS session_responses (
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

CREATE TABLE IF NOT EXISTS iq_results (
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

CREATE TABLE IF NOT EXISTS admins (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username          VARCHAR(50) UNIQUE NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
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

CREATE INDEX IF NOT EXISTS idx_sessions_user ON test_sessions(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_token ON test_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_responses_session ON session_responses(session_id);
CREATE INDEX IF NOT EXISTS idx_results_session ON iq_results(session_id);
CREATE INDEX IF NOT EXISTS idx_payments_user ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_questions_domain ON questions(domain) WHERE is_active = TRUE;
`

// SeedQuestionsSQL mengisi tabel questions dari data default (dijalankan terpisah)
const SeedQuestionsSQL = `
INSERT INTO questions (question_code, domain, difficulty, weight, image_url, option_a_image, option_b_image, option_c_image, option_d_image, correct_option)
VALUES
  ('Q_MTX_001', 'MTX', 'easy', 1.0, '/assets/images/q_mtx_001.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', 'A'),
  ('Q_MTX_002', 'MTX', 'easy', 1.0, '/assets/images/q_mtx_002.svg', '/assets/images/opt_a2.svg', '/assets/images/opt_b2.svg', '/assets/images/opt_c2.svg', '/assets/images/opt_d2.svg', 'B'),
  ('Q_MTX_003', 'MTX', 'medium', 1.5, '/assets/images/q_mtx_003.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', 'C'),
  ('Q_MTX_004', 'MTX', 'medium', 1.5, '/assets/images/q_mtx_004.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', 'D'),
  ('Q_MTX_005', 'MTX', 'hard', 2.0, '/assets/images/q_mtx_005.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', 'A'),
  ('Q_MTX_006', 'MTX', 'very_hard', 2.5, '/assets/images/q_mtx_006.svg', '/assets/images/opt_a2.svg', '/assets/images/opt_b2.svg', '/assets/images/opt_c2.svg', '/assets/images/opt_d2.svg', 'B'),
  ('Q_SEQ_001', 'SEQ', 'easy', 1.0, '/assets/images/q_seq_001.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', 'A'),
  ('Q_SEQ_002', 'SEQ', 'medium', 1.5, '/assets/images/q_seq_002.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', 'B'),
  ('Q_SEQ_003', 'SEQ', 'medium', 1.5, '/assets/images/q_seq_003.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', 'C'),
  ('Q_SEQ_004', 'SEQ', 'hard', 2.0, '/assets/images/q_seq_004.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', 'D'),
  ('Q_SEQ_005', 'SEQ', 'hard', 2.0, '/assets/images/q_seq_005.svg', '/assets/images/opt_a2.svg', '/assets/images/opt_b2.svg', '/assets/images/opt_c2.svg', '/assets/images/opt_d2.svg', 'A'),
  ('Q_SPA_001', 'SPA', 'medium', 1.5, '/assets/images/q_spa_001.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', 'B'),
  ('Q_SPA_002', 'SPA', 'medium', 1.5, '/assets/images/q_spa_002.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', 'C'),
  ('Q_SPA_003', 'SPA', 'hard', 2.0, '/assets/images/q_spa_003.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', 'D'),
  ('Q_SPA_004', 'SPA', 'very_hard', 2.5, '/assets/images/q_spa_004.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', 'A'),
  ('Q_SPA_005', 'SPA', 'very_hard', 2.5, '/assets/images/q_spa_005.svg', '/assets/images/opt_a2.svg', '/assets/images/opt_b2.svg', '/assets/images/opt_c2.svg', '/assets/images/opt_d2.svg', 'B'),
  ('Q_ANL_001', 'ANL', 'easy', 1.0, '/assets/images/q_anl_001.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', 'C'),
  ('Q_ANL_002', 'ANL', 'medium', 1.5, '/assets/images/q_anl_002.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', 'D'),
  ('Q_ANL_003', 'ANL', 'medium', 1.5, '/assets/images/q_anl_003.svg', '/assets/images/opt_c.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', 'A'),
  ('Q_ANL_004', 'ANL', 'hard', 2.0, '/assets/images/q_anl_004.svg', '/assets/images/opt_d.svg', '/assets/images/opt_a.svg', '/assets/images/opt_b.svg', '/assets/images/opt_c.svg', 'B')
ON CONFLICT (question_code) DO NOTHING;
`
