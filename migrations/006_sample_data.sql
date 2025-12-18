-- migrations/006_sample_data_extended.sql
SET client_encoding = 'UTF8';

-- Очистка данных в правильном порядке
TRUNCATE TABLE 
    material_completions,
    favorite_materials,
    material_ratings,
    material_blocks,
    materials,
    teacher_specializations,
    users 
CASCADE;

-- Вставляем пользователей с красивыми ФИО
INSERT INTO users (id, email, password_hash, full_name, role, avatar_url, is_verified, created_at) VALUES
-- Администраторы
(1, 'admin@eduplatform.ru', '$2b$10$examplehash', 'Иванов Александр Сергеевич', 'admin', '/avatars/admin1.jpg', TRUE, '2024-01-01 10:00:00+03'),
(2, 'support@eduplatform.ru', '$2b$10$examplehash', 'Сидорова Ольга Владимировна', 'admin', '/avatars/admin2.jpg', TRUE, '2024-01-02 11:00:00+03'),

-- Учителя математики
(3, 'petrova.math@school.ru', '$2b$10$examplehash', 'Петрова Мария Ивановна', 'teacher', '/avatars/teacher1.jpg', TRUE, '2024-01-03 09:00:00+03'),
(4, 'sidorov.math@school.ru', '$2b$10$examplehash', 'Сидоров Игорь Петрович', 'teacher', '/avatars/teacher2.jpg', TRUE, '2024-01-04 10:00:00+03'),
(5, 'kozlov.math@school.ru', '$2b$10$examplehash', 'Козлов Андрей Викторович', 'teacher', '/avatars/teacher3.jpg', TRUE, '2024-01-05 11:00:00+03'),

-- Учителя программирования
(6, 'ivanova.code@school.ru', '$2b$10$examplehash', 'Иванова Елена Дмитриевна', 'teacher', '/avatars/teacher4.jpg', TRUE, '2024-01-06 12:00:00+03'),
(7, 'smirnov.code@school.ru', '$2b$10$examplehash', 'Смирнов Алексей Николаевич', 'teacher', '/avatars/teacher5.jpg', TRUE, '2024-01-07 13:00:00+03'),
(8, 'fedorov.code@school.ru', '$2b$10$examplehash', 'Фёдоров Денис Олегович', 'teacher', '/avatars/teacher6.jpg', TRUE, '2024-01-08 14:00:00+03'),

-- Учителя физики
(9, 'nikolaeva.physics@school.ru', '$2b$10$examplehash', 'Николаева Анна Сергеевна', 'teacher', '/avatars/teacher7.jpg', TRUE, '2024-01-09 15:00:00+03'),
(10, 'voronov.physics@school.ru', '$2b$10$examplehash', 'Воронов Павел Игоревич', 'teacher', '/avatars/teacher8.jpg', TRUE, '2024-01-10 16:00:00+03'),
(11, 'orlova.physics@school.ru', '$2b$10$examplehash', 'Орлова Светлана Михайловна', 'teacher', '/avatars/teacher9.jpg', TRUE, '2024-01-11 17:00:00+03'),

-- Студенты
(12, 'student1@edu.ru', '$2b$10$examplehash', 'Ковалев Алексей Дмитриевич', 'student', '/avatars/student1.jpg', TRUE, '2024-01-12 18:00:00+03'),
(13, 'student2@edu.ru', '$2b$10$examplehash', 'Новикова Марина Андреевна', 'student', '/avatars/student2.jpg', TRUE, '2024-01-13 19:00:00+03'),
(14, 'student3@edu.ru', '$2b$10$examplehash', 'Морозов Артем Сергеевич', 'student', '/avatars/student3.jpg', TRUE, '2024-01-14 20:00:00+03'),
(15, 'student4@edu.ru', '$2b$10$examplehash', 'Волкова Дарья Игоревна', 'student', '/avatars/student4.jpg', TRUE, '2024-01-15 21:00:00+03'),
(16, 'student5@edu.ru', '$2b$10$examplehash', 'Лебедев Павел Викторович', 'student', '/avatars/student5.jpg', TRUE, '2024-01-16 22:00:00+03'),
(17, 'student6@edu.ru', '$2b$10$examplehash', 'Соколова София Александровна', 'student', '/avatars/student6.jpg', TRUE, '2024-01-17 23:00:00+03'),
(18, 'student7@edu.ru', '$2b$10$examplehash', 'Попов Максим Олегович', 'student', '/avatars/student7.jpg', TRUE, '2024-01-18 08:00:00+03'),
(19, 'student8@edu.ru', '$2b$10$examplehash', 'Кузнецова Анастасия Денисовна', 'student', '/avatars/student8.jpg', TRUE, '2024-01-19 09:00:00+03');

-- Специализации учителей
INSERT INTO teacher_specializations (user_id, subject, created_at) VALUES
-- Петрова Мария Ивановна - математика
(3, 'mathematics', '2024-01-03 09:30:00+03'),
(3, 'algebra', '2024-01-03 09:35:00+03'),
(3, 'calculus', '2024-01-03 09:40:00+03'),

-- Сидоров Игорь Петрович - математика
(4, 'mathematics', '2024-01-04 10:30:00+03'),
(4, 'geometry', '2024-01-04 10:35:00+03'),
(4, 'trigonometry', '2024-01-04 10:40:00+03'),

-- Козлов Андрей Викторович - математика
(5, 'mathematics', '2024-01-05 11:30:00+03'),
(5, 'statistics', '2024-01-05 11:35:00+03'),
(5, 'probability', '2024-01-05 11:40:00+03'),

-- Иванова Елена Дмитриевна - программирование
(6, 'programming', '2024-01-06 12:30:00+03'),
(6, 'python', '2024-01-06 12:35:00+03'),
(6, 'datascience', '2024-01-06 12:40:00+03'),

-- Смирнов Алексей Николаевич - программирование
(7, 'programming', '2024-01-07 13:30:00+03'),
(7, 'javascript', '2024-01-07 13:35:00+03'),
(7, 'web', '2024-01-07 13:40:00+03'),
(7, 'frontend', '2024-01-07 13:45:00+03'),

-- Фёдоров Денис Олегович - программирование
(8, 'programming', '2024-01-08 14:30:00+03'),
(8, 'java', '2024-01-08 14:35:00+03'),
(8, 'algorithms', '2024-01-08 14:40:00+03'),
(8, 'backend', '2024-01-08 14:45:00+03'),

-- Николаева Анна Сергеевна - физика
(9, 'physics', '2024-01-09 15:30:00+03'),
(9, 'mechanics', '2024-01-09 15:35:00+03'),
(9, 'kinematics', '2024-01-09 15:40:00+03'),

-- Воронов Павел Игоревич - физика
(10, 'physics', '2024-01-10 16:30:00+03'),
(10, 'electrodynamics', '2024-01-10 16:35:00+03'),
(10, 'optics', '2024-01-10 16:40:00+03'),

-- Орлова Светлана Михайловна - физика
(11, 'physics', '2024-01-11 17:30:00+03'),
(11, 'thermodynamics', '2024-01-11 17:35:00+03'),
(11, 'quantum', '2024-01-11 17:40:00+03');

-- Образовательные материалы
INSERT INTO materials (id, title, subject, author_id, status, access, share_url, created_at, updated_at) VALUES
-- Математика - Петрова М.И.
(1, 'Основы алгебры: переменные и уравнения', 'algebra', 3, 'published', 'open', 'algebra-basics', '2024-01-15 09:00:00+03', '2024-01-15 09:00:00+03'),
(2, 'Дифференциальное исчисление для начинающих', 'calculus', 3, 'published', 'open', 'calculus-basics', '2024-01-16 10:00:00+03', '2024-01-16 10:00:00+03'),

-- Математика - Сидоров И.П.
(3, 'Геометрия: треугольники и их свойства', 'geometry', 4, 'published', 'open', 'geometry-triangles', '2024-01-17 11:00:00+03', '2024-01-17 11:00:00+03'),
(4, 'Тригонометрия: синусы, косинусы и тангенсы', 'trigonometry', 4, 'published', 'open', 'trigonometry-basics', '2024-01-18 12:00:00+03', '2024-01-18 12:00:00+03'),

-- Математика - Козлов А.В.
(5, 'Теория вероятностей: основы', 'probability', 5, 'published', 'open', 'probability-basics', '2024-01-19 13:00:00+03', '2024-01-19 13:00:00+03'),
(6, 'Математическая статистика', 'statistics', 5, 'published', 'open', 'statistics-basics', '2024-01-20 14:00:00+03', '2024-01-20 14:00:00+03'),

-- Программирование - Иванова Е.Д.
(7, 'Python: первые шаги в программировании', 'python', 6, 'published', 'open', 'python-first-steps', '2024-01-21 15:00:00+03', '2024-01-21 15:00:00+03'),
(8, 'Анализ данных с Pandas и NumPy', 'datascience', 6, 'published', 'open', 'data-analysis-python', '2024-01-22 16:00:00+03', '2024-01-22 16:00:00+03'),

-- Программирование - Смирнов А.Н.
(9, 'Веб-разработка на JavaScript', 'javascript', 7, 'published', 'open', 'javascript-web', '2024-01-23 17:00:00+03', '2024-01-23 17:00:00+03'),
(10, 'Современный фронтенд: React и Vue', 'frontend', 7, 'published', 'open', 'modern-frontend', '2024-01-24 18:00:00+03', '2024-01-24 18:00:00+03'),

-- Программирование - Фёдоров Д.О.
(11, 'Java: основы объектно-ориентированного программирования', 'java', 8, 'published', 'open', 'java-oop', '2024-01-25 19:00:00+03', '2024-01-25 19:00:00+03'),
(12, 'Алгоритмы и структуры данных', 'algorithms', 8, 'published', 'open', 'algorithms-data-structures', '2024-01-26 20:00:00+03', '2024-01-26 20:00:00+03'),

-- Физика - Николаева А.С.
(13, 'Механика: законы движения Ньютона', 'mechanics', 9, 'published', 'open', 'newton-laws', '2024-01-27 21:00:00+03', '2024-01-27 21:00:00+03'),
(14, 'Кинематика: движение тел', 'kinematics', 9, 'published', 'open', 'kinematics-motion', '2024-01-28 22:00:00+03', '2024-01-28 22:00:00+03'),

-- Физика - Воронов П.И.
(15, 'Электродинамика: основы', 'electrodynamics', 10, 'published', 'open', 'electrodynamics-basics', '2024-01-29 23:00:00+03', '2024-01-29 23:00:00+03'),
(16, 'Оптика: свет и линзы', 'optics', 10, 'published', 'open', 'optics-basics', '2024-01-30 08:00:00+03', '2024-01-30 08:00:00+03'),

-- Физика - Орлова С.М.
(17, 'Термодинамика: законы и применения', 'thermodynamics', 11, 'published', 'open', 'thermodynamics-laws', '2024-01-31 09:00:00+03', '2024-01-31 09:00:00+03'),
(18, 'Введение в квантовую физику', 'quantum', 11, 'published', 'open', 'quantum-physics', '2024-02-01 10:00:00+03', '2024-02-01 10:00:00+03');

-- Блоки материалов (расширенное содержание)
INSERT INTO material_blocks (material_id, block_id, type, content, position) VALUES
-- Материал 1: Основы алгебры
(1, 'title1', 'text', '{"text": "Основы алгебры", "level": "h1"}', 1),
(1, 'intro1', 'text', '{"text": "Алгебра - это раздел математики, который изучает математические операции и отношения.", "level": "p"}', 2),
(1, 'example1', 'text', '{"text": "Пример линейного уравнения: 2x + 3 = 7", "level": "p"}', 3),

-- Материал 7: Python программирование
(7, 'title7', 'text', '{"text": "Python для начинающих", "level": "h1"}', 1),
(7, 'intro7', 'text', '{"text": "Python - мощный и простой в изучении язык программирования.", "level": "p"}', 2),
(7, 'code1', 'text', '{"text": "print(\\\"Привет, мир!\\\")", "level": "code"}', 3),

-- Материал 13: Физика
(13, 'title13', 'text', '{"text": "Законы Ньютона", "level": "h1"}', 1),
(13, 'law1', 'text', '{"text": "Первый закон Ньютона: тело сохраняет состояние покоя или равномерного движения...", "level": "p"}', 2);

-- Рейтинги материалов
INSERT INTO material_ratings (material_id, user_id, rating, created_at) VALUES
-- Рейтинги для математических материалов
(1, 12, 5, '2024-01-15 14:00:00+03'),
(1, 13, 4, '2024-01-15 15:00:00+03'),
(1, 14, 5, '2024-01-15 16:00:00+03'),
(3, 12, 4, '2024-01-17 14:00:00+03'),
(3, 15, 5, '2024-01-17 15:00:00+03'),
(5, 13, 5, '2024-01-19 14:00:00+03'),

-- Рейтинги для программирования
(7, 12, 5, '2024-01-21 16:00:00+03'),
(7, 14, 4, '2024-01-21 17:00:00+03'),
(7, 16, 5, '2024-01-21 18:00:00+03'),
(9, 13, 4, '2024-01-23 16:00:00+03'),
(9, 17, 5, '2024-01-23 17:00:00+03'),
(11, 14, 5, '2024-01-25 16:00:00+03'),

-- Рейтинги для физики
(13, 15, 4, '2024-01-27 18:00:00+03'),
(13, 16, 5, '2024-01-27 19:00:00+03'),
(15, 17, 4, '2024-01-29 18:00:00+03'),
(17, 12, 5, '2024-01-31 19:00:00+03');

-- Избранные материалы
INSERT INTO favorite_materials (user_id, material_id, created_at) VALUES
(12, 1, '2024-01-15 16:00:00+03'),
(12, 7, '2024-01-21 17:00:00+03'),
(13, 3, '2024-01-17 15:00:00+03'),
(13, 9, '2024-01-23 16:00:00+03'),
(14, 5, '2024-01-19 15:00:00+03'),
(14, 11, '2024-01-25 17:00:00+03'),
(15, 13, '2024-01-27 18:00:00+03'),
(16, 7, '2024-01-21 18:00:00+03'),
(17, 9, '2024-01-23 17:00:00+03');

-- Завершенные материалы (прогресс обучения)
INSERT INTO material_completions (user_id, material_id, time_spent, grade, completed_at, last_activity) VALUES
(12, 1, 3600, 4.5, '2024-01-15 17:00:00+03', '2024-01-15 17:00:00+03'),
(12, 7, 4800, 4.8, '2024-01-21 18:00:00+03', '2024-01-21 18:00:00+03'),
(13, 3, 4200, 4.2, '2024-01-17 16:00:00+03', '2024-01-17 16:00:00+03'),
(14, 5, 3800, 4.7, '2024-01-19 16:00:00+03', '2024-01-19 16:00:00+03'),
(15, 13, 5200, 4.9, '2024-01-27 19:00:00+03', '2024-01-27 19:00:00+03');

-- Сбрасываем последовательности
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('materials_id_seq', (SELECT MAX(id) FROM materials));
SELECT setval('material_blocks_id_seq', (SELECT MAX(id) FROM material_blocks));
SELECT setval('material_ratings_id_seq', (SELECT MAX(id) FROM material_ratings));
SELECT setval('teacher_specializations_id_seq', (SELECT MAX(id) FROM teacher_specializations));
SELECT setval('favorite_materials_id_seq', (SELECT MAX(id) FROM favorite_materials));
SELECT setval('material_completions_id_seq', (SELECT MAX(id) FROM material_completions));