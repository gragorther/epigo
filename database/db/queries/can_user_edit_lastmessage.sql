SELECT
  EXISTS (
    SELECT 1
    FROM last_messages lm
    WHERE lm.id = $1 AND lm.user_id = $2
  )
  AND
  (CASE
     WHEN $3::int[] IS NULL OR cardinality($3::int[]) = 0 THEN TRUE
     ELSE (
       (SELECT COUNT(*) FROM groups g WHERE g.id = ANY($3::int[]) AND g.user_id = $2)
       = cardinality($3)
     )
   END);
