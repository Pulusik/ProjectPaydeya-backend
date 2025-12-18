-- Таблица предметов (категорий) - используем множественное число для консистентности
CREATE TABLE IF NOT EXISTS subjects (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(200)
);



-- Таблица материалов
CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(1000) NOT NULL,
    subject_id VARCHAR(200) NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),
    access VARCHAR(20) NOT NULL DEFAULT 'open'
        CHECK (access IN ('open', 'link')),
    share_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_material_subject
        FOREIGN KEY (subject_id)
        REFERENCES subjects(id)
);
ALTER TABLE materials RENAME COLUMN subject TO subject_id;

-- Индексы
CREATE INDEX IF NOT EXISTS idx_materials_author_id ON materials(author_id);
CREATE INDEX IF NOT EXISTS idx_materials_subject_id ON materials(subject_id);
CREATE INDEX IF NOT EXISTS idx_materials_status ON materials(status);

-- Добавляем тестовые предметы
INSERT INTO subjects (id, name, icon) VALUES
    ('informatics', 'Информатика', '/icons/informatics.svg'),
    ('mathematics', 'Математика', '/icons/mathematics.svg'),
    ('physics', 'Физика', '/icons/physics.svg'),
    ('programming', 'Программирование', '/icons/programming.svg')
ON CONFLICT (id) DO NOTHING;

-- Таблица блоков материалов
CREATE TABLE IF NOT EXISTS material_blocks (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    block_id VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'image', 'video', 'formula', 'quiz')),
    content JSONB NOT NULL,
    styles JSONB,
    animation JSONB,
    position INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(material_id, block_id)
);

-- Индексы для блоков (убрал дублирующийся индекс)
CREATE INDEX IF NOT EXISTS idx_material_blocks_position ON material_blocks(material_id, position);