package repositories

import (
    "context"
    "fmt"
    "strings"
    "log"

    "paydeya-backend/internal/models"

    //"github.com/jackc/pgx/v5"
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
    var materials []models.CatalogMaterial
    var total int

    // Базовый запрос
    baseQuery := `
        SELECT m.id, m.title, m.subject,
               u.id as author_id, u.full_name as author_name,
               COALESCE(rm.rating, 0) as rating,
               COALESCE(rm.students_count, 0) as students_count
        FROM materials m
        JOIN users u ON m.author_id = u.id
        LEFT JOIN (
            SELECT material_id,
                   AVG(rating) as rating,
                   COUNT(*) as students_count
            FROM material_ratings
            GROUP BY material_id
        ) rm ON m.id = rm.material_id
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
        conditions = append(conditions, fmt.Sprintf("m.subject = $%d", argIndex))
        args = append(args, filters.Subject)
        argIndex++
    }

    if filters.Level != "" {
        conditions = append(conditions, fmt.Sprintf("m.level = $%d", argIndex))
        args = append(args, filters.Level)
        argIndex++
    }

    // Добавляем условия в запрос
    if len(conditions) > 0 {
        baseQuery += " AND " + strings.Join(conditions, " AND ")
    }

    // Запрос для общего количества
    countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as filtered"
    err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    // Добавляем пагинацию и сортировку
    baseQuery += " ORDER BY rating DESC, m.updated_at DESC"

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
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var material models.CatalogMaterial
        var author models.Author

        err := rows.Scan(
            &material.ID, &material.Title, &material.Subject,
            &author.ID, &author.Name, &material.Rating, &material.StudentsCount,
        )
        if err != nil {
            return nil, 0, err
        }

        material.Author = author
        materials = append(materials, material)
    }

    return materials, total, nil
}

// GetSubjects возвращает список предметов
func (r *CatalogRepository) GetSubjects(ctx context.Context) ([]models.Subject, error) {
    query := `SELECT id, name, icon FROM subjects ORDER BY name`

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var subjects []models.Subject
    for rows.Next() {
        var subject models.Subject
        if err := rows.Scan(&subject.ID, &subject.Name, &subject.Icon); err != nil {
            return nil, err
        }
        subjects = append(subjects, subject)
    }

    return subjects, nil
}

// SearchTeachers поиск преподавателей
func (r *CatalogRepository) SearchTeachers(ctx context.Context, filters models.TeacherFilters) ([]models.Teacher, error) {
    log.Printf("Building query with filters: %+v", filters)
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

    log.Printf("Final SQL query: %s", query) // ← ДОБАВЬТЕ
    log.Printf("Query args: %v", args)       // ← ДОБАВЬТЕ

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teachers []models.Teacher
    for rows.Next() {
        var teacher models.Teacher
        var avatarURL *string
        if err := rows.Scan(&teacher.ID, &teacher.Name, &avatarURL, &teacher.MaterialsCount, &teacher.Rating); err != nil {
            return nil, err
        }

        // Получаем специализации учителя
        specializations, err := r.getTeacherSpecializations(ctx, teacher.ID)
        if err != nil {
            log.Printf("SQL query error: %v", err) // ← ДОБАВЬТЕ
            return nil, err
        }
        teacher.Specializations = specializations

        teachers = append(teachers, teacher)
    }

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

    return specializations, nil
}