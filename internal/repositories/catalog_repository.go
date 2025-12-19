package repositories

import (
    "context"
    "fmt"
    "strings"
    "log"

    "paydeya-backend/internal/models"
    "github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
    db *pgxpool.Pool
}

func NewCatalogRepository(db *pgxpool.Pool) *CatalogRepository {
    return &CatalogRepository{db: db}
}

// SearchMaterials поиск материалов с фильтрацией
func (r *CatalogRepository) SearchMaterials(ctx context.Context, filters models.CatalogFilters) ([]models.CatalogMaterial, int, error) {
    // Объявляем переменные здесь, в начале функции
    var materials []models.CatalogMaterial
    var total int

    // Базовый запрос
    baseQuery := `
        SELECT
            m.id,
            m.title,
            m.subject_id as subject,
            u.id as author_id,
            u.full_name as author_name,
            COALESCE(AVG(mr.rating), 0) as rating,
            COUNT(DISTINCT mr.user_id) as students_count
        FROM materials m
        JOIN users u ON m.author_id = u.id
        LEFT JOIN material_ratings mr ON m.id = mr.material_id
        WHERE m.status = 'published'
    `

    var conditions []string
    var args []interface{}
    argIndex := 1

    // Добавляем условия фильтрации
    if filters.Search != "" {
        conditions = append(conditions, fmt.Sprintf("(m.title ILIKE $%d OR u.full_name ILIKE $%d)", argIndex, argIndex))
        args = append(args, "%"+filters.Search+"%")
        argIndex++
    }

    if filters.Subject != "" {
        conditions = append(conditions, fmt.Sprintf("m.subject_id = $%d", argIndex))
        args = append(args, filters.Subject)
        argIndex++
    }

    // УБЕРИТЕ фильтр по level - его нет в таблице!
    // if filters.Level != "" {
    //     conditions = append(conditions, fmt.Sprintf("m.level = $%d", argIndex))
    //     args = append(args, filters.Level)
    //     argIndex++
    // }

    // Добавляем условия в запрос
    if len(conditions) > 0 {
        baseQuery += " AND " + strings.Join(conditions, " AND ")
    }

    // Добавляем GROUP BY
    baseQuery += " GROUP BY m.id, m.title, m.subject_id, u.id, u.full_name"

    // Запрос для общего количества - исправленная версия
    countQuery := `
        SELECT COUNT(DISTINCT m.id)
        FROM materials m
        JOIN users u ON m.author_id = u.id
        LEFT JOIN material_ratings mr ON m.id = mr.material_id
        WHERE m.status = 'published'
    `

    if len(conditions) > 0 {
        countQuery += " AND " + strings.Join(conditions, " AND ")
    }

    // Выполняем запрос количества
    err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        log.Printf("❌ Error counting materials: %v", err)
        return nil, 0, err
    }

    // Добавляем пагинацию и сортировку
    baseQuery += " ORDER BY rating DESC NULLS LAST, m.updated_at DESC"

    if filters.Limit > 0 {
        baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
        args = append(args, filters.Limit)
        argIndex++

        if filters.Page > 0 {
            offset := (filters.Page - 1) * filters.Limit
            baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
            args = append(args, offset)
        }
    }

    // Выполняем основной запрос
    rows, err := r.db.Query(ctx, baseQuery, args...)
    if err != nil {
        log.Printf("❌ Error querying materials: %v", err)
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var material models.CatalogMaterial
        var author models.Author

        err := rows.Scan(
            &material.ID,
            &material.Title,
            &material.Subject,
            &author.ID,
            &author.Name,
            &material.Rating,
            &material.StudentsCount,
        )
        if err != nil {
            log.Printf("❌ Error scanning material row: %v", err)
            return nil, 0, err
        }

        material.Author = author
        materials = append(materials, material)
    }

    // Проверяем ошибки после итерации
    if err := rows.Err(); err != nil {
        log.Printf("❌ Error after rows iteration: %v", err)
        return nil, 0, err
    }

    log.Printf("✅ SearchMaterials: found %d materials (total: %d)", len(materials), total)
    return materials, total, nil
}

// GetSubjects возвращает список предметов
func (r *CatalogRepository) GetSubjects(ctx context.Context) ([]models.Subject, error) {
    // Получаем уникальные предметы из материалов
    query := `
        SELECT DISTINCT
            m.subject_id as id,
            COALESCE(s.name, m.subject_id) as name
        FROM materials m
        LEFT JOIN subjects s ON m.subject_id = s.id
        WHERE m.status = 'published'
        ORDER BY name
    `

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        log.Printf("❌ Error querying subjects: %v", err)
        return nil, err
    }
    defer rows.Close()

    var subjects []models.Subject
    for rows.Next() {
        var subject models.Subject
        if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
            log.Printf("❌ Error scanning subject row: %v", err)
            return nil, err
        }
        subjects = append(subjects, subject)
    }

    if err := rows.Err(); err != nil {
        log.Printf("❌ Error after subjects rows iteration: %v", err)
        return nil, err
    }

    log.Printf("✅ GetSubjects: found %d subjects", len(subjects))
    return subjects, nil
}

// SearchTeachers поиск преподавателей
func (r *CatalogRepository) SearchTeachers(ctx context.Context, filters models.TeacherFilters) ([]models.Teacher, error) {
    log.Printf("🔍 SearchTeachers called with filters: %+v", filters)

    query := `
        SELECT u.id, u.full_name, u.avatar_url,
               COUNT(DISTINCT m.id) as materials_count,
               COALESCE(AVG(mr.rating), 0) as rating
        FROM users u
        LEFT JOIN materials m ON u.id = m.author_id AND m.status = 'published'
        LEFT JOIN material_ratings mr ON m.id = mr.material_id
        WHERE u.role = 'teacher'
    `

    var conditions []string
    var args []interface{}
    argIndex := 1

    if filters.Search != "" {
        conditions = append(conditions, fmt.Sprintf("u.full_name ILIKE $%d", argIndex))
        args = append(args, "%"+filters.Search+"%")
        argIndex++
    }

    if filters.Subject != "" {
        conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM teacher_specializations ts WHERE ts.user_id = u.id AND ts.subject = $%d)", argIndex))
        args = append(args, filters.Subject)
        argIndex++
    }

    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    query += " GROUP BY u.id, u.full_name, u.avatar_url ORDER BY rating DESC NULLS LAST, materials_count DESC"

    log.Printf("📝 SQL query: %s", query)
    log.Printf("📝 Query args: %v", args)

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        log.Printf("❌ Error querying teachers: %v", err)
        return nil, err
    }
    defer rows.Close()

    var teachers []models.Teacher
    for rows.Next() {
        var teacher models.Teacher
        var avatarURL *string
        if err := rows.Scan(&teacher.ID, &teacher.Name, &avatarURL, &teacher.MaterialsCount, &teacher.Rating); err != nil {
            log.Printf("❌ Error scanning teacher row: %v", err)
            return nil, err
        }

        // Получаем специализации учителя
        specializations, err := r.getTeacherSpecializations(ctx, teacher.ID)
        if err != nil {
            log.Printf("❌ Error getting specializations for teacher %d: %v", teacher.ID, err)
            // Продолжаем без специализаций
        }
        teacher.Specializations = specializations

        teachers = append(teachers, teacher)
    }

    if err := rows.Err(); err != nil {
        log.Printf("❌ Error after teachers rows iteration: %v", err)
        return nil, err
    }

    log.Printf("✅ SearchTeachers: found %d teachers", len(teachers))
    return teachers, nil
}

func (r *CatalogRepository) getTeacherSpecializations(ctx context.Context, teacherID int) ([]string, error) {
    query := `SELECT subject FROM teacher_specializations WHERE user_id = $1 ORDER BY subject`

    rows, err := r.db.Query(ctx, query, teacherID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var specializations []string
    for rows.Next() {
        var subject string
        if err := rows.Scan(&subject); err != nil {
            return nil, err
        }
        specializations = append(specializations, subject)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return specializations, nil
}