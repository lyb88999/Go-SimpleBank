-- name: CreateVerifyEmail :one
INSERT INTO verify_emails(
    username,email,secret_code
) values($1,$2,$3)
RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = @id AND secret_code = @secret_code AND is_used=false AND expired_at > now()
RETURNING *;