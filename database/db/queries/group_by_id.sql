SELECT groups.name, groups.description, groups.user_id, ARRAY_AGG(group_last_messages.last_message_id) FILTER (WHERE group_last_messages.last_message_id IS NOT NULL) FROM groups LEFT JOIN
group_last_messages ON groups.id = group_last_messages.group_id WHERE groups.id = $1 GROUP BY groups.id
