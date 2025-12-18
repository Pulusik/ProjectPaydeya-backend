package repositories

import (
    "context"
    "strings"


    "paydeya-backend/internal/models"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, full_name, role, avatar_url, is_verified)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

    err := r.db.QueryRow(ctx, query,
        user.Email, user.PasswordHash, user.FullName, user.Role, user.AvatarURL, user.IsVerified,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    return err
}

// GetUserByEmail возвращает пользователя по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    var blockReason *string

    query := `
        SELECT id, email, password_hash, full_name, role, avatar_url, is_verified, is_blocked, block_reason, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    err := r.db.QueryRow(ctx, query, email).Scan(
           &user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
           &user.AvatarURL, &user.IsVerified, &user.IsBlocked, &blockReason,
           &user.CreatedAt, &user.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    user.BlockReason = blockReason
    return &user, err
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    var user models.User

    query := `
        SELECT id, email, password_hash, full_name, role, avatar_url, is_verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
        &user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    return &user, err
}

// EmailExists проверяет, существует ли email
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    err := r.db.QueryRow(ctx, query, email).Scan(&exists)
    return exists, err
}
// GetUserProfile возвращает профиль пользователя по ID
func (r *UserRepository) GetUserProfile(ctx context.Context, userID int) (*models.User, error) {
    var user models.User

    query := `
        SELECT id, email, full_name, role, avatar_url, is_verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, userID).Scan(
        &user.ID, &user.Email, &user.FullName, &user.Role,
        &user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    return &user, err
}

// UpdateUserProfile обновляет данные пользователя
func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID int, fullName string, specializations []string) error {
    // Начинаем транзакцию
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Обновляем основную информацию
    _, err = tx.Exec(ctx,
        "UPDATE users SET full_name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
        fullName, userID,
    )
    if err != nil {
        return err
    }

    // Если пользователь - учитель, обновляем специализации
    var userRole string
    err = tx.QueryRow(ctx, "SELECT role FROM users WHERE id = $1", userID).Scan(&userRole)
    if err != nil {
        return err
    }

    if userRole == "teacher" {
        // Удаляем старые специализации
        _, err = tx.Exec(ctx, "DELETE FROM teacher_specializations WHERE user_id = $1", userID)
        if err != nil {
            return err
        }

        // Добавляем новые специализации
        for _, subject := range specializations {
            _, err = tx.Exec(ctx,
                "INSERT INTO teacher_specializations (user_id, subject) VALUES ($1, $2)",
                userID, strings.TrimSpace(subject),
            )
            if err != nil {
                return err
            }
        }
    }

    return tx.Commit(ctx)
}

// UpdateUserAvatar обновляет аватар пользователя
func (r *UserRepository) UpdateUserAvatar(ctx context.Context, userID int, avatarURL string) error {
    query := `
        UPDATE users
        SET avatar_url = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `

    _, err := r.db.Exec(ctx, query, avatarURL, userID)
    return err
}
// GetUserSpecializations возвращает специализации пользователя
func (r *UserRepository) GetUserSpecializations(ctx context.Context, userID int) ([]string, error) {
    query := `
        SELECT subject
        FROM teacher_specializations
        WHERE user_id = $1
        ORDER BY subject
    `

    rows, err := r.db.Query(ctx, query, userID)
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

// UpdateUserSpecializations обновляет специализации пользователя
func (r *UserRepository) UpdateUserSpecializations(ctx context.Context, userID int, specializations []string) error {
    // Начинаем транзакцию
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Удаляем старые специализации
    _, err = tx.Exec(ctx, "DELETE FROM teacher_specializations WHERE user_id = $1", userID)
    if err != nil {
        return err
    }

    // Добавляем новые специализации
    for _, subject := range specializations {
        _, err = tx.Exec(ctx,
            "INSERT INTO teacher_specializations (user_id, subject) VALUES ($1, $2)",
            userID, strings.TrimSpace(subject),
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit(ctx)
}