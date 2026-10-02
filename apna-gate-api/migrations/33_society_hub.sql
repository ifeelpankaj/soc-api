-- +migrate Up
CREATE TABLE community_categories (
    id SMALLSERIAL PRIMARY KEY,

    code VARCHAR(30) NOT NULL,
    name VARCHAR(50) NOT NULL,

    display_order SMALLINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_community_categories_code UNIQUE (code),
    CONSTRAINT chk_community_categories_code_not_blank
        CHECK (BTRIM(code) <> ''),
    CONSTRAINT chk_community_categories_name_not_blank
        CHECK (BTRIM(name) <> ''),
    CONSTRAINT chk_community_categories_display_order
        CHECK (display_order >= 0)
);

CREATE INDEX idx_community_categories_active_order
    ON community_categories(is_active, display_order, id);

INSERT INTO community_categories (code, name, display_order)
VALUES
    ('general', 'General', 1),
    ('buy_sell', 'Buy & Sell', 2),
    ('lost_found', 'Lost & Found', 3),
    ('recommendation', 'Recommendations', 4),
    ('help_needed', 'Help Needed', 5),
    ('feedback', 'Feedback', 6),
    ('events_entertainment', 'Events & Entertainment', 7),
    ('services', 'Services', 8),
    ('carpool', 'Carpool', 9),
    ('other', 'Other', 99)
ON CONFLICT (code) DO NOTHING;


-- --------------------------------------------------------------------------
-- 2. Society channels
-- Exactly two permanent channels per society.
-- --------------------------------------------------------------------------

CREATE TABLE society_channels (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL
        REFERENCES societies(id)
        ON DELETE RESTRICT,

    type VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_society_channels_type
        CHECK (type IN ('announcement', 'community')),

    CONSTRAINT chk_society_channels_name_not_blank
        CHECK (BTRIM(name) <> ''),

    CONSTRAINT uq_society_channels_society_type
        UNIQUE (society_id, type),

    -- Supports composite FKs that guarantee a resource belongs to the same society.
    CONSTRAINT uq_society_channels_id_society
        UNIQUE (id, society_id)
);

CREATE INDEX idx_society_channels_society
    ON society_channels(society_id);


-- --------------------------------------------------------------------------
-- 3. Backfill channels for all existing societies.
-- --------------------------------------------------------------------------

INSERT INTO society_channels (
    society_id,
    type,
    name,
    description
)
SELECT
    s.id,
    'announcement',
    'Announcements',
    'Official society updates'
FROM societies s
ON CONFLICT (society_id, type) DO NOTHING;

INSERT INTO society_channels (
    society_id,
    type,
    name,
    description
)
SELECT
    s.id,
    'community',
    'Community',
    'Connect with neighbours'
FROM societies s
ON CONFLICT (society_id, type) DO NOTHING;


-- --------------------------------------------------------------------------
-- 4. Channel posts
-- Announcement and Community content both live here.
--
-- category_id:
--   - NULL for announcement posts
--   - required by application validation for community posts
--
-- The database cannot express that rule with a simple CHECK because channel type
-- lives in society_channels. Enforce it in the Go service layer.
-- --------------------------------------------------------------------------

CREATE TABLE channel_posts (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL,
    channel_id BIGINT NOT NULL,

    author_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    title VARCHAR(200),
    body TEXT NOT NULL,

    category_id SMALLINT
        REFERENCES community_categories(id)
        ON DELETE RESTRICT,

    comments_enabled BOOLEAN NOT NULL DEFAULT TRUE,

    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_important BOOLEAN NOT NULL DEFAULT FALSE,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    edited_at TIMESTAMPTZ,

    removal_reason TEXT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
        REFERENCES users(id)
        ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_channel_posts_channel_society
        FOREIGN KEY (channel_id, society_id)
        REFERENCES society_channels(id, society_id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_channel_posts_status
        CHECK (status IN ('active', 'deleted', 'removed')),

    CONSTRAINT chk_channel_posts_body_not_blank
        CHECK (BTRIM(body) <> ''),

    CONSTRAINT chk_channel_posts_title_not_blank
        CHECK (title IS NULL OR BTRIM(title) <> ''),

    CONSTRAINT chk_channel_posts_deleted_state
        CHECK (
            (status = 'active' AND deleted_at IS NULL)
            OR
            (status IN ('deleted', 'removed') AND deleted_at IS NOT NULL)
        ),

    -- Supports safer downstream composite references if needed later.
    CONSTRAINT uq_channel_posts_id_society
        UNIQUE (id, society_id),
    UNIQUE (id, channel_id)
);

CREATE INDEX idx_channel_posts_feed
    ON channel_posts(
        channel_id,
        created_at DESC,
        id DESC
    )
    WHERE status = 'active';

CREATE INDEX idx_channel_posts_society
    ON channel_posts(
        society_id,
        created_at DESC,
        id DESC
    )
    WHERE status = 'active';

CREATE INDEX idx_channel_posts_category
    ON channel_posts(
        channel_id,
        category_id,
        created_at DESC,
        id DESC
    )
    WHERE status = 'active'
      AND category_id IS NOT NULL;

CREATE INDEX idx_channel_posts_pinned
    ON channel_posts(
        channel_id,
        created_at DESC,
        id DESC
    )
    WHERE status = 'active'
      AND is_pinned = TRUE;

CREATE INDEX idx_channel_posts_author
    ON channel_posts(author_id, created_at DESC, id DESC);


-- --------------------------------------------------------------------------
-- 5. Channel comments
-- Supports top-level comments and optional one-level reply structure.
-- Enforce maximum nesting depth in Go service logic.
-- --------------------------------------------------------------------------

CREATE TABLE channel_comments (
    id BIGSERIAL PRIMARY KEY,

    post_id BIGINT NOT NULL
        REFERENCES channel_posts(id)
        ON DELETE RESTRICT,

    author_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    parent_id BIGINT,
    UNIQUE (id, post_id),
    FOREIGN KEY (parent_id, post_id) REFERENCES channel_comments(id, post_id) ON DELETE RESTRICT,

    body TEXT NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    edited_at TIMESTAMPTZ,

    removal_reason TEXT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
        REFERENCES users(id)
        ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_channel_comments_status
        CHECK (status IN ('active', 'deleted', 'removed')),

    CONSTRAINT chk_channel_comments_body_not_blank
        CHECK (BTRIM(body) <> ''),

    CONSTRAINT chk_channel_comments_not_self_parent
        CHECK (parent_id IS NULL OR parent_id <> id),

    CONSTRAINT chk_channel_comments_deleted_state
        CHECK (
            (status = 'active' AND deleted_at IS NULL)
            OR
            (status IN ('deleted', 'removed') AND deleted_at IS NOT NULL)
        )
);

CREATE INDEX idx_channel_comments_post
    ON channel_comments(post_id, created_at ASC, id ASC);

CREATE INDEX idx_channel_comments_parent
    ON channel_comments(parent_id, created_at ASC, id ASC)
    WHERE parent_id IS NOT NULL;

CREATE INDEX idx_channel_comments_author
    ON channel_comments(author_id, created_at DESC, id DESC);


-- --------------------------------------------------------------------------
-- 6. Channel reactions
-- A reaction belongs to exactly one target: a post OR a comment.
-- --------------------------------------------------------------------------

CREATE TABLE channel_reactions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    post_id BIGINT
        REFERENCES channel_posts(id)
        ON DELETE CASCADE,

    comment_id BIGINT
        REFERENCES channel_comments(id)
        ON DELETE CASCADE,

    reaction_type VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

    CONSTRAINT chk_channel_reactions_target
        CHECK (
            (post_id IS NOT NULL AND comment_id IS NULL)
            OR
            (post_id IS NULL AND comment_id IS NOT NULL)
        ),

    CONSTRAINT chk_channel_reactions_type
        CHECK (reaction_type IN ('like', 'love', 'helpful'))
);

CREATE UNIQUE INDEX uq_channel_reaction_post
    ON channel_reactions(user_id, post_id, reaction_type)
    WHERE post_id IS NOT NULL;

CREATE UNIQUE INDEX uq_channel_reaction_comment
    ON channel_reactions(user_id, comment_id, reaction_type)
    WHERE comment_id IS NOT NULL;

CREATE INDEX idx_channel_reactions_post
    ON channel_reactions(post_id, reaction_type)
    WHERE post_id IS NOT NULL;

CREATE INDEX idx_channel_reactions_comment
    ON channel_reactions(comment_id, reaction_type)
    WHERE comment_id IS NOT NULL;


-- --------------------------------------------------------------------------
-- 7. Private uploads and attachment claims. Provider metadata is stored once in
-- channel_uploads; channel_attachments references it with a unique claim.
-- Signed URLs are generated on authorized reads and never persisted.
-- --------------------------------------------------------------------------

CREATE TABLE channel_uploads (
 id BIGSERIAL PRIMARY KEY,
 society_id BIGINT NOT NULL REFERENCES societies(id),
 uploader_id BIGINT NOT NULL REFERENCES users(id),
 file_id TEXT NOT NULL UNIQUE,
 file_path TEXT NOT NULL,
 original_filename TEXT NOT NULL,
 mime_type TEXT NOT NULL CHECK (mime_type IN ('image/jpeg','image/png','application/pdf')),
 file_size BIGINT NOT NULL CHECK (file_size > 0),
 width INTEGER CHECK (width > 0), height INTEGER CHECK (height > 0),
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','claimed','cleanup','deleted')),
 expires_at TIMESTAMPTZ NOT NULL DEFAULT now()+interval '24 hours',
 lease_until TIMESTAMPTZ,
 lease_token UUID,
 created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
CREATE INDEX channel_uploads_cleanup ON channel_uploads(expires_at,id) WHERE status IN ('pending','cleanup');
CREATE TABLE channel_attachments (
 id BIGSERIAL PRIMARY KEY,
 upload_id BIGINT NOT NULL UNIQUE REFERENCES channel_uploads(id),
 post_id BIGINT REFERENCES channel_posts(id),
 comment_id BIGINT REFERENCES channel_comments(id),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 CHECK ((post_id IS NOT NULL)::int + (comment_id IS NOT NULL)::int = 1)
);
CREATE INDEX channel_attachments_post ON channel_attachments(post_id) WHERE post_id IS NOT NULL;
CREATE INDEX channel_attachments_comment ON channel_attachments(comment_id) WHERE comment_id IS NOT NULL;

-- --------------------------------------------------------------------------
-- 8. Read states
-- One row per user per channel.
--
-- last_read_post_id gives a deterministic cursor for unread calculations.
-- last_read_at is retained for audit/UI purposes.
-- --------------------------------------------------------------------------

CREATE TABLE channel_read_states (
    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    channel_id BIGINT NOT NULL
        REFERENCES society_channels(id)
        ON DELETE CASCADE,

    last_read_post_id BIGINT NOT NULL,
    last_read_created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY(last_read_post_id, channel_id) REFERENCES channel_posts(id, channel_id),

    last_read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, channel_id)
);

CREATE INDEX idx_channel_read_states_channel
    ON channel_read_states(channel_id);

CREATE INDEX idx_channel_read_states_last_post
    ON channel_read_states(last_read_post_id)
    WHERE last_read_post_id IS NOT NULL;



CREATE TABLE channel_reports (
 id BIGSERIAL PRIMARY KEY,
 post_id BIGINT NOT NULL REFERENCES channel_posts(id),
 reporter_id BIGINT NOT NULL REFERENCES users(id),
 reason TEXT NOT NULL CHECK (length(btrim(reason)) BETWEEN 1 AND 2000),
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','dismissed','actioned')),
 reviewer_id BIGINT REFERENCES users(id),
 reviewed_at TIMESTAMPTZ,
 resolution_note TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
 UNIQUE(post_id,reporter_id),
 CHECK ((status='pending' AND reviewer_id IS NULL AND reviewed_at IS NULL) OR
        (status<>'pending' AND reviewer_id IS NOT NULL AND reviewed_at IS NOT NULL))
);
CREATE INDEX channel_reports_pending ON channel_reports(status,created_at,id);
CREATE INDEX channel_comments_rate ON channel_comments(author_id,created_at DESC);

-- +migrate Down
DROP TABLE channel_reports;
DROP TABLE channel_read_states;
DROP TABLE channel_attachments;
DROP TABLE channel_uploads;
DROP TABLE channel_reactions;
DROP TABLE channel_comments;
DROP TABLE channel_posts;
DROP TABLE society_channels;
DROP TABLE community_categories;
