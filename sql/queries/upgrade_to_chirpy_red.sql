-- name: UpgradeUserToChirpyRed :one
UPDATE users
SET
    is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red;
