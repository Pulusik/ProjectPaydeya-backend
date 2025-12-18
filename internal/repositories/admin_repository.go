package repositories

import (
    "context"
    "fmt"
    "strings"

    "paydeya-backend/internal/models"

    //"github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
    db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
    return &AdminRepository{db: db}
}

// GetPlatformStats возвращает статистику платформы
func (r *AdminRepository) GetPlatformStats(ctx context.Context) (*models.AdminStats, error) {
    var stats models.AdminStats

    // Общее количество пользователей
    query := `SELECT COUNT(*) FROM users`
    err := r.db.QueryRow(ctx, query).Scan(&stats.TotalUsers)
    if err != nil {
        return nil, err
    }

    // Общее количество материалов
    query = `SELECT COUNT(*) FROM materials`
    err = r.db.QueryRow(ctx, query).Scan(&stats.TotalMaterials)
    if err != nil {
        return nil, err
    }

    // Активные преподаватели (создавшие хотя бы 1 материал)
    query = `
        SELECT COUNT(DISTINCT u.id)
        FROM users u
        JOIN materials m ON u.id = m.author_id
        WHERE u.role = 'teacher'
    `
    err = r.db.QueryRow(ctx, query).Scan(&stats.ActiveTeachers)
    if err != nil {
        return nil, err
    }

    // Опубликованные материалы
    query = `SELECT COUNT(*) FROM materials WHERE status = 'published'`
    err = r.db.QueryRow(ctx, query).Scan(&stats.PublishedMaterials)
    if err != nil {
        return nil, err
    }

    return &stats, nil
}

// GetUsers возвращает список пользователей с фильтрацией
func (r *AdminRepository) GetUsers(ctx context.Context, role string, page, limit int) ([]models.UserManagement, int, error) {
    var users []models.UserManagement
    var total int

    // Базовый запрос
    baseQuery := `
        SELECT u.id, u.email, u.full_name, u.role, u.is_verified, u.is_blocked, u.block_reason, u.created_at,
               COUNT(m.id) as materials_count
        FROM users u
        LEFT JOIN materials m ON u.id = m.author_id
    `

    var conditions []string
    var args []interface{}
    argIndex := 1

    if role != "" {
        conditions = append(conditions, fmt.Sprintf("u.role = $%d", argIndex))
        args = append(args, role)
        argIndex++
    }

    if len(conditions) > 0 {
        baseQuery += " WHERE " + strings.Join(conditions, " AND ")
    }

    baseQuery += " GROUP BY u.id, u.email, u.full_name, u.role, u.is_verified, u.is_blocked, u.block_reason, u.created_at"

    // Получаем общее количество
    countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as filtered"
    err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    // Добавляем пагинацию
    baseQuery += " ORDER BY u.created_at DESC"

    if limit > 0 {
        baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
        args = append(args, limit)
        argIndex++

        if page > 0 {
            offset := (page - 1) * limit
            baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
            args = append(args, offset)
        }
    }

    // Выполняем запрос

    rows, err := r.db.Query(ctx, baseQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var user models.UserManagement
        // Убрали createdAt переменную - сканируем прямо в структуру

        err := rows.Scan(
            &user.ID, &user.Email, &user.FullName, &user.Role,
            &user.IsVerified, &user.IsBlocked, &user.BlockReason, &user.CreatedAt, &user.MaterialsCount, // ← сканируем прямо в time.Time поле
        )
        if err != nil {
            return nil, 0, err
        }

        users = append(users, user)
    }
    return users, total, nil
}

// BlockUser блокирует пользователя
func (r *AdminRepository) BlockUser(ctx context.Context, userID int, reason string) error {
    query := `UPDATE users SET is_blocked = true, block_reason = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
    _, err := r.db.Exec(ctx, query, reason, userID)
    return err
}

// CreateSubject создает новый предмет
func (r *AdminRepository) CreateSubject(ctx context.Context, req *models.CreateSubjectRequest) error {
    query := `INSERT INTO subjects (id, name, icon) VALUES ($1, $2, $3)`
    _, err := r.db.Exec(ctx, query, req.ID, req.Name, req.Icon)
    return err
}