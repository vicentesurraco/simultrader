-- name: SubStonk :exec
INSERT INTO stonks (user_id, symbol)
VALUES ($1, $2)
ON CONFLICT (user_id, symbol) 
DO UPDATE SET is_active = TRUE;

-- name: UnsubStonk :exec
UPDATE stonks
SET is_active = FALSE
WHERE user_id = $1 AND symbol = $2;

-- name: GetUserActiveStonks :many
select symbol
from stonks
where user_id = $1 and is_active = TRUE
;
