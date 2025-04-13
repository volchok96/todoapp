-- Безопасное удаление таблицы
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'tasks') THEN
        DROP TABLE tasks CASCADE;
        RAISE NOTICE 'Table "tasks" dropped';
    ELSE
        RAISE NOTICE 'Table "tasks" does not exist';
    END IF;
END $$;