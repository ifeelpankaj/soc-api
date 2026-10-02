-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    timezone,
    language,
    metadata
)
VALUES (
    sqlc.narg('first_name'),
    sqlc.narg('last_name'),
    sqlc.arg('full_name'),
    sqlc.narg('email'),
    sqlc.narg('phone_number'),
    sqlc.narg('password_hash'),
    sqlc.arg('auth_provider')::auth_provider,
    sqlc.arg('global_role')::global_role,
    sqlc.arg('email_verified'),
    sqlc.arg('phone_verified'),
    sqlc.arg('is_active'),
    sqlc.arg('is_blocked'),
    sqlc.arg('timezone'),
    sqlc.arg('language'),
    sqlc.arg('metadata')
)
RETURNING
    id,
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    provider_id,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    blocked_reason,
    avatar_url,
    date_of_birth,
    gender,
    timezone,
    language,
    last_login_at,
    password_changed_at,
    deleted_at,
    metadata,
    created_at,
    updated_at,
    avatar_imagekit_file_id,
    avatar_imagekit_file_path,
    session_version;

-- name: GetUserByID :one
SELECT
    id,
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    provider_id,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    blocked_reason,
    avatar_url,
    date_of_birth,
    gender,
    timezone,
    language,
    last_login_at,
    password_changed_at,
    deleted_at,
    metadata,
    created_at,
    updated_at,
    avatar_imagekit_file_id,
    avatar_imagekit_file_path,
    session_version
FROM users
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
    id,
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    provider_id,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    blocked_reason,
    avatar_url,
    date_of_birth,
    gender,
    timezone,
    language,
    last_login_at,
    password_changed_at,
    deleted_at,
    metadata,
    created_at,
    updated_at,
    avatar_imagekit_file_id,
    avatar_imagekit_file_path,
    session_version
FROM users
WHERE LOWER(email) = LOWER($1)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByPhoneNumber :one
SELECT
    id,
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    provider_id,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    blocked_reason,
    avatar_url,
    date_of_birth,
    gender,
    timezone,
    language,
    last_login_at,
    password_changed_at,
    deleted_at,
    metadata,
    created_at,
    updated_at,
    avatar_imagekit_file_id,
    avatar_imagekit_file_path,
    session_version
FROM users
WHERE phone_number = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: EmailExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE LOWER(email) = LOWER($1)
      AND deleted_at IS NULL
);

-- name: PhoneExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE phone_number = $1
      AND deleted_at IS NULL
);

-- name: UpdateUserPasswordHash :execrows
UPDATE users
SET
    password_hash = $2,
    session_version = session_version + 1,
    password_changed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND session_version = $3
  AND deleted_at IS NULL;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
    last_login_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateUserProfile :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    full_name = TRIM(BOTH FROM CONCAT(
        COALESCE(sqlc.narg('first_name'), first_name, ''),
        ' ',
        COALESCE(sqlc.narg('last_name'), last_name, '')
    )),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    date_of_birth = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    gender = COALESCE(sqlc.narg('gender'), gender),
    timezone = COALESCE(sqlc.narg('timezone'), timezone),
    language = COALESCE(sqlc.narg('language'), language),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
  AND deleted_at IS NULL
RETURNING
    id,
    first_name,
    last_name,
    full_name,
    email,
    phone_number,
    password_hash,
    auth_provider,
    provider_id,
    global_role,
    email_verified,
    phone_verified,
    is_active,
    is_blocked,
    avatar_url,
    date_of_birth,
    gender,
    timezone,
    language,
    last_login_at,
    password_changed_at,
    metadata,
    created_at,
    updated_at,
    deleted_at,
    session_version;
