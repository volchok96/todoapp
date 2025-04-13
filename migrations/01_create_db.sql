-- Проверяем существование БД и создаем если нет
SELECT 'CREATE DATABASE todoapp'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'todoapp')\gexec

-- Подключаемся к новой БД
\c todoapp

-- Создаем таблицу с проверкой на существование
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    date DATE NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для поиска по дате
CREATE INDEX IF NOT EXISTS idx_tasks_date ON tasks(date);