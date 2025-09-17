SELECT
  EXISTS (
    SELECT 1
    FROM groups g
    WHERE g.id = $1 AND g.user_id = $2
  )
  AND
  (CASE
     WHEN $3::int[] IS NULL OR cardinality($3::int[]) = 0 THEN TRUE
     ELSE (
       (SELECT COUNT(*) FROM last_messages lm WHERE lm.id = ANY($3::int[]) AND lm.user_id = $2)
       = cardinality($3)
     )
   END);
