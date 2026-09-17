INSERT INTO groups (name, description, created_at)
SELECT
    'Default',
    'Группа по умолчанию для всех пользователей',
    CURRENT_TIMESTAMP
WHERE NOT EXISTS (
    SELECT 1
    FROM groups
    WHERE name = 'Default'
);