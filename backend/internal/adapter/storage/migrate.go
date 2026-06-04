package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// 스키마 정의를 dialect 별로 보관한다. 운영 환경에서는 별도 마이그레이션 도구
// (golang-migrate, sql-migrate 등)로 대체하는 것을 권장하지만, 부트스트랩 단계에서는
// 애플리케이션 기동 시 idempotent 하게 생성한다.

const schemaSQLite = `
CREATE TABLE IF NOT EXISTS categories (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    parent_id   INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        TEXT NOT NULL,
    slug         TEXT NOT NULL UNIQUE,
    excerpt      TEXT NOT NULL DEFAULT '',
    content      TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'draft',
    category_id  INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    published_at DATETIME,
    created_at   DATETIME NOT NULL,
    updated_at   DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id      INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id    INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    author_name  TEXT NOT NULL,
    author_email TEXT NOT NULL DEFAULT '',
    content      TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'pending',
    created_at   DATETIME NOT NULL,
    updated_at   DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
`

const schemaPostgres = `
CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(120) NOT NULL,
    slug        VARCHAR(140) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    parent_id   BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(80) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id           BIGSERIAL PRIMARY KEY,
    title        VARCHAR(200) NOT NULL,
    slug         VARCHAR(220) NOT NULL UNIQUE,
    excerpt      VARCHAR(500) NOT NULL DEFAULT '',
    content      TEXT NOT NULL DEFAULT '',
    status       VARCHAR(20) NOT NULL DEFAULT 'draft',
    category_id  BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id           BIGSERIAL PRIMARY KEY,
    post_id      BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id    BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    author_name  VARCHAR(120) NOT NULL,
    author_email VARCHAR(160) NOT NULL DEFAULT '',
    content      TEXT NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
`

// Migrate 는 dialect 에 맞는 스키마를 생성한다.
func Migrate(db *sqlx.DB, driver string) error {
	schema := schemaSQLite
	if driver == DriverPostgres {
		schema = schemaPostgres
	}
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("스키마 마이그레이션 실패: %w", err)
	}
	return nil
}
