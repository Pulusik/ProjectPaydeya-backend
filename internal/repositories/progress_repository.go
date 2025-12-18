package repositories

import (
    "context"
    "time"

    "paydeya-backend/internal/models"

    //"github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository struct {
    db *pgxpool.Pool
}

func NewProgressRepository(db *pgxpool.Pool) *ProgressRepository {
    return &ProgressRepository{db: db}
}

// GetStudentProgress возвращает общую статистику ученика
func (r *ProgressRepository) GetStudentProgress(ctx context.Context, userID int) (*models.StudentProgress, error) {
    var progress models.StudentProgress

    // Получаем количество завершенных материалов
    query := `SELECT COUNT(*) FROM material_completions WHERE user_id = $1`
    err := r.db.QueryRow(ctx, query, userID).Scan(&progress.CompletedTopics)
    if err != nil {
        return nil, err
    }

    // Получаем среднюю оценку
    query = `SELECT COALESCE(AVG(grade), 0) FROM material_completions WHERE user_id = $1`
    err = r.db.QueryRow(ctx, query, userID).Scan(&progress.AverageGrade)
    if err != nil {
        return nil, err
    }

    // Получаем общее время обучения (в часах)
    query = `SELECT COALESCE(SUM(time_spent), 0) / 3600 FROM material_completions WHERE user_id = $1`
    err = r.db.QueryRow(ctx, query, userID).Scan(&progress.LearningHours)
    if err != nil {
        return nil, err
    }

    // Простой расчет успеваемости
    if progress.CompletedTopics > 0 {
        progress.SuccessRate = progress.AverageGrade / 5 * 100
    }

    // Получаем текущие материалы (последние 5)
    query = `
        SELECT m.id, m.title, m.subject, mc.last_activity
        FROM materials m
        JOIN material_completions mc ON m.id = mc.material_id
        WHERE mc.user_id = $1
        ORDER BY mc.last_activity DESC
        LIMIT 5
    `
    rows, err := r.db.Query(ctx, query, userID)
    if err == nil {
        defer rows.Close()

        for rows.Next() {
            var material models.ProgressMaterial
            var lastActivity time.Time

            if err := rows.Scan(&material.ID, &material.Title, &material.Subject, &lastActivity); err == nil {
                material.LastActivity = lastActivity
                material.Progress = 100.0 // если в completion, значит завершен
                progress.CurrentMaterials = append(progress.CurrentMaterials, material)
            }
        }
    }

    return &progress, nil
}

// MarkMaterialComplete отмечает материал как завершенный
func (r *ProgressRepository) MarkMaterialComplete(ctx context.Context, userID, materialID int, timeSpent int, grade float64) error {
    query := `
        INSERT INTO material_completions (user_id, material_id, time_spent, grade, completed_at, last_activity)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (user_id, material_id)
        DO UPDATE SET time_spent = EXCLUDED.time_spent, grade = EXCLUDED.grade, last_activity = EXCLUDED.last_activity
    `

    now := time.Now()
    _, err := r.db.Exec(ctx, query, userID, materialID, timeSpent, grade, now, now)
    return err
}

// GetFavoriteMaterials возвращает избранные материалы
func (r *ProgressRepository) GetFavoriteMaterials(ctx context.Context, userID int) ([]models.CatalogMaterial, error) {
    query := `
        SELECT m.id, m.title, m.subject,
               u.id as author_id, u.full_name as author_name,
               4.5 as rating, 10 as students_count
        FROM materials m
        JOIN users u ON m.author_id = u.id
        JOIN favorite_materials fm ON m.id = fm.material_id
        WHERE fm.user_id = $1 AND m.status = 'published'
        ORDER BY fm.created_at DESC
    `

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var materials []models.CatalogMaterial
    for rows.Next() {
        var material models.CatalogMaterial
        var author models.Author

        err := rows.Scan(
            &material.ID, &material.Title, &material.Subject,
            &author.ID, &author.Name, &material.Rating, &material.StudentsCount,
        )
        if err != nil {
            return nil, err
        }

        material.Author = author
        materials = append(materials, material)
    }

    return materials, nil
}

// ToggleFavorite добавляет/удаляет материал из избранного
func (r *ProgressRepository) ToggleFavorite(ctx context.Context, userID, materialID int, action string) error {
    if action == "add" {
        query := `INSERT INTO favorite_materials (user_id, material_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
        _, err := r.db.Exec(ctx, query, userID, materialID)
        return err
    } else {
        query := `DELETE FROM favorite_materials WHERE user_id = $1 AND material_id = $2`
        _, err := r.db.Exec(ctx, query, userID, materialID)
        return err
    }
}