SELECT last_messages.title,
last_messages.content,
ARRAY_AGG(DISTINCT recipients.email)
FROM last_messages
INNER JOIN group_last_messages ON group_last_messages.last_message_id = last_messages.id
INNER JOIN recipients ON recipients.group_id = group_last_messages.group_id
WHERE last_messages.user_id = $1
GROUP BY last_messages.id
