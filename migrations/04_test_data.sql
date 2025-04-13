-- Очищаем таблицу перед наполнением (опционально)
TRUNCATE TABLE tasks RESTART IDENTITY CASCADE;

-- Вставляем тестовые данные
INSERT INTO tasks (title, description, date, done)
VALUES 
    ('Купить продукты', 'Молоко, хлеб, яйца', CURRENT_DATE + INTERVAL '1 day', false),
    ('Сделать ДЗ', 'Проект по Go', CURRENT_DATE, false),
    ('Позвонить маме', 'Обсудить праздники', CURRENT_DATE - INTERVAL '2 days', true),
    ('Записаться к врачу', 'Стоматолог на следующей неделе', CURRENT_DATE + INTERVAL '7 days', false),
    ('Прочитать книгу', '"Clean Code" - глава 5', CURRENT_DATE + INTERVAL '3 days', false);

-- Проверяем
SELECT * FROM tasks ORDER BY date;