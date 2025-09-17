WITH g AS (
  INSERT INTO groups (name, description, user_id)
    VALUES ($1, $2, $3)
  RETURNING id
),
ins AS (
  INSERT INTO group_last_messages (group_id, last_message_id)
  SELECT g.id, UNNEST($4::int[])
  FROM g
  RETURNING group_id
)
SELECT id AS group_id
FROM g;
