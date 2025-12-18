-- migrations/002_add_specializations_table.sql

-- Таблица специализаций учителей (связь многие-ко-многим)
CREATE TABLE IF NOT EXISTS teacher_specializations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, subject) -- один пользователь не может иметь дублирующиеся специализации
);

-- Индекс для быстрого поиска специализаций по пользователю
CREATE INDEX IF NOT EXISTS idx_teacher_specializations_user_id ON teacher_specializations(user_id);
