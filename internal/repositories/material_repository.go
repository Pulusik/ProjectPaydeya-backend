package repositories

import (
    "context"

    "paydeya-backend/internal/models"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type MaterialRepository struct {
    db *pgxpool.Pool
}

func NewMaterialRepository(db *pgxpool.Pool) *MaterialRepository {
    return &MaterialRepository{db: db}
}

// CreateMaterial создает новый материал
func (r *MaterialRepository) CreateMaterial(ctx context.Context, material *models.Material) error {
    query := `
        INSERT INTO materials (title, subject, author_id, status, access, share_url)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

    err := r.db.QueryRow(ctx, query,
        material.Title, material.Subject, material.AuthorID,
        material.Status, material.Access, material.ShareURL,
    ).Scan(&material.ID, &material.CreatedAt, &material.UpdatedAt)

    return err
}

// GetMaterial возвращает материал по ID
func (r *MaterialRepository) GetMaterial(ctx context.Context, id int) (*models.Material, error) {
    var material models.Material

    query := `
        SELECT id, title, subject, author_id, status, access, share_url, created_at, updated_at
        FROM materials
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &material.ID, &material.Title, &material.Subject, &material.AuthorID,
        &material.Status, &material.Access, &material.ShareURL,
        &material.CreatedAt, &material.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    return &material, err
}

// GetUserMaterials возвращает материалы пользователя
func (r *MaterialRepository) GetUserMaterials(ctx context.Context, userID int, status string) ([]*models.Material, error) {
    var query string
    var rows pgx.Rows
    var err error

    if status == "" {
        query = `SELECT id, title, subject, status, access, created_at, updated_at
                 FROM materials WHERE author_id = $1 ORDER BY updated_at DESC`
        rows, err = r.db.Query(ctx, query, userID)
    } else {
        query = `SELECT id, title, subject, status, access, created_at, updated_at
                 FROM materials WHERE author_id = $1 AND status = $2 ORDER BY updated_at DESC`
        rows, err = r.db.Query(ctx, query, userID, status)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var materials []*models.Material
    for rows.Next() {
        var material models.Material
        if err := rows.Scan(
            &material.ID, &material.Title, &material.Subject,
            &material.Status, &material.Access, &material.CreatedAt, &material.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        material.AuthorID = userID
        materials = append(materials, &material)
    }

    return materials, nil
}

// UpdateMaterial обновляет материал
func (r *MaterialRepository) UpdateMaterial(ctx context.Context, material *models.Material) error {
    query := `
        UPDATE materials
        SET title = $1, subject = $2, status = $3, access = $4, share_url = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6 AND author_id = $7
    `

    _, err := r.db.Exec(ctx, query,
        material.Title, material.Subject, material.Status, material.Access,
        material.ShareURL, material.ID, material.AuthorID,
    )
    return err
}

