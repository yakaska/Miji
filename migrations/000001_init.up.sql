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
    owner_id      BIGINT       NOT NULL REFERENCES miji.users (id),
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
    link_id     BIGINT      NOT NULL REFERENCES miji.links (id),
    visited_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address  INET,
    user_agent  TEXT,
    referer     TEXT
);

-- =============================================================
-- INDEXES
-- =============================================================

-- Основной путь редиректа: поиск по slug
CREATE INDEX idx_links_slug ON miji.links (slug);

-- Фоновая очистка истёкших ссылок
CREATE INDEX idx_links_expires_at ON miji.links (expires_at)
    WHERE expires_at IS NOT NULL;

-- Список ссылок конкретного пользователя
CREATE INDEX idx_links_owner_id ON miji.links (owner_id);

-- Аналитика: переходы по ссылке за период
CREATE INDEX idx_visits_link_id_visited_at ON miji.visits (link_id, visited_at DESC);

-- Топ источников трафика
CREATE INDEX idx_visits_referer ON miji.visits (link_id, referer)
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
    BEFORE UPDATE ON miji.users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- =============================================================
-- COMMENTS
-- =============================================================
COMMENT ON TABLE  miji.users                  IS 'Пользователи сервиса';
COMMENT ON COLUMN miji.users.role             IS 'user | admin';
COMMENT ON COLUMN miji.users.is_active        IS 'false = пользователь заблокирован, сессии инвалидируются';

COMMENT ON TABLE  miji.links                  IS 'Сокращённые ссылки';
COMMENT ON COLUMN miji.links.slug             IS 'Короткий код, base62, 3–20 символов';
COMMENT ON COLUMN miji.links.expires_at       IS 'NULL = бессрочно';
COMMENT ON COLUMN miji.links.is_active        IS 'false = soft-delete, редирект возвращает 404';
COMMENT ON COLUMN miji.links.visit_count      IS 'Денормализованный счётчик, обновляется Kafka-consumer';

COMMENT ON TABLE  miji.visits                 IS 'Лог переходов, пишется Kafka-consumer асинхронно';
COMMENT ON COLUMN miji.visits.visited_at      IS 'Фиксируется в producer до отправки в Kafka';