-- Улучшенная версия с проверками
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'tasks') THEN
        CREATE TABLE tasks (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            date DATE NOT NULL,
            done BOOLEAN NOT NULL DEFAULT FALSE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX idx_tasks_date ON tasks(date);
        RAISE NOTICE 'Table "tasks" created';
    ELSE
        RAISE NOTICE 'Table "tasks" already exists';
    END IF;
END $$;