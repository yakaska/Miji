-- =============================================================
-- Migration: 0001_init.up.sql
-- URL Shortener — начальная схема данных
-- =============================================================

CREATE SCHEMA miji;

-- =============================================================
-- TABLE: users
-- =============================================================
CREATE TABLE miji.users (
    id            BIGSERIAL    PRIMARY KEY,
    email         VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(10)  NOT NULL DEFAULT 'user',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_role_check   CHECK  (role IN ('user', 'admin'))
);

-- =============================================================
-- TABLE: links
-- =============================================================
CREATE TABLE miji.links (
    id            BIGSERIAL    PRIMARY KEY,
    owner_id      BIGINT       NOT NULL REFERENCES users (id),
    slug          VARCHAR(100)  NOT NULL,
    original_url  TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    visit_count   BIGINT       NOT NULL DEFAULT 0,

    CONSTRAINT links_slug_unique UNIQUE (slug)
);

-- =============================================================
-- TABLE: visits
-- =============================================================
CREATE TABLE miji.visits (
    id          BIGSERIAL   PRIMARY KEY,
    link_id     BIGINT      NOT NULL REFERENCES links (id),
    visited_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address  INET,
    user_agent  TEXT,
    referer     TEXT
);

-- =============================================================
-- INDEXES
-- =============================================================

-- Основной путь редиректа: поиск по slug
CREATE INDEX miji.idx_links_slug ON links (slug);

-- Фоновая очистка истёкших ссылок
CREATE INDEX miji.idx_links_expires_at ON links (expires_at)
    WHERE expires_at IS NOT NULL;

-- Список ссылок конкретного пользователя
CREATE INDEX miji.idx_links_owner_id ON links (owner_id);

-- Аналитика: переходы по ссылке за период
CREATE INDEX miji.idx_visits_link_id_visited_at ON visits (link_id, visited_at DESC);

-- Топ источников трафика
CREATE INDEX miji.idx_visits_referer ON visits (link_id, referer)
    WHERE referer IS NOT NULL;

-- =============================================================
-- TRIGGER: updated_at на таблице users
-- =============================================================
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- =============================================================
-- COMMENTS
-- =============================================================
COMMENT ON TABLE  users                  IS 'Пользователи сервиса';
COMMENT ON COLUMN users.role             IS 'user | admin';
COMMENT ON COLUMN users.is_active        IS 'false = пользователь заблокирован, сессии инвалидируются';

COMMENT ON TABLE  links                  IS 'Сокращённые ссылки';
COMMENT ON COLUMN links.slug             IS 'Короткий код, base62, 3–20 символов';
COMMENT ON COLUMN links.expires_at       IS 'NULL = бессрочно';
COMMENT ON COLUMN links.is_active        IS 'false = soft-delete, редирект возвращает 404';
COMMENT ON COLUMN links.visit_count      IS 'Денормализованный счётчик, обновляется Kafka-consumer';

COMMENT ON TABLE  visits                 IS 'Лог переходов, пишется Kafka-consumer асинхронно';
COMMENT ON COLUMN visits.visited_at      IS 'Фиксируется в producer до отправки в Kafka';