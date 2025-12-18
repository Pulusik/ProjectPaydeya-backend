package repositories

import (
    "context"
    "encoding/json"

    "paydeya-backend/internal/models"

    //"github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type BlockRepository struct {
    db *pgxpool.Pool
}

func NewBlockRepository(db *pgxpool.Pool) *BlockRepository {
    return &BlockRepository{db: db}
}

// SaveBlocks сохраняет блоки материала
func (r *BlockRepository) SaveBlocks(ctx context.Context, materialID int, blocks []models.Block) error {
    // Начинаем транзакцию
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Удаляем старые блоки
    _, err = tx.Exec(ctx, "DELETE FROM material_blocks WHERE material_id = $1", materialID)
    if err != nil {
        return err
    }

    // Сохраняем новые блоки
    for position, block := range blocks {
        contentJSON, err := json.Marshal(block.Content)
        if err != nil {
            return err
        }

        stylesJSON, err := json.Marshal(block.Styles)
        if err != nil {
            return err
        }

        var animationJSON []byte
        if block.Animation != nil {
            animationJSON, err = json.Marshal(block.Animation)
            if err != nil {
                return err
            }
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO material_blocks (material_id, block_id, type, content, styles, animation, position)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, materialID, block.ID, block.Type, contentJSON, stylesJSON, animationJSON, position)

        if err != nil {
            return err
        }
    }

    return tx.Commit(ctx)
}

// GetBlocks возвращает блоки материала
func (r *BlockRepository) GetBlocks(ctx context.Context, materialID int) ([]models.Block, error) {
    query := `
        SELECT block_id, type, content, styles, animation, position
        FROM material_blocks
        WHERE material_id = $1
        ORDER BY position
    `

    rows, err := r.db.Query(ctx, query, materialID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var blocks []models.Block
    for rows.Next() {
        var block models.Block
        var contentJSON, stylesJSON, animationJSON []byte

        err := rows.Scan(&block.ID, &block.Type, &contentJSON, &stylesJSON, &animationJSON, &block.Position)
        if err != nil {
            return nil, err
        }

        // Парсим JSON поля
        if err := json.Unmarshal(contentJSON, &block.Content); err != nil {
            return nil, err
        }
        if err := json.Unmarshal(stylesJSON, &block.Styles); err != nil {
            return nil, err
        }
        if len(animationJSON) > 0 {
            var animation models.BlockAnimation
            if err := json.Unmarshal(animationJSON, &animation); err != nil {
                return nil, err
            }
            block.Animation = &animation
        }

        blocks = append(blocks, block)
    }

    return blocks, nil
}