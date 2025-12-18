package database

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"  // ← ИЗМЕНИТЬ ИМПОРТ
)

var DB *pgxpool.Pool  // ← ИЗМЕНИТЬ ТИП

func Init(cfg *Config) error {
    connString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
    )

    log.Printf("🔗 Connecting to: %s@%s:%d/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

    // ИЗМЕНИТЬ: используем pgxpool вместо pgx.Connect
    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        log.Printf("❌ Connection failed: %v", err)
        return fmt.Errorf("unable to connect to database: %w", err)
    }

    // Test connection
    if err := pool.Ping(context.Background()); err != nil {
        log.Printf("❌ Ping failed: %v", err)
        pool.Close()
        return fmt.Errorf("unable to ping database: %w", err)
    }

    DB = pool  // ← ТЕПЕРЬ POOL
    log.Printf("✅ Database connected successfully with connection pool! DB pointer: %p", DB)
    return nil
}

func Close() {
    if DB != nil {
        DB.Close()  // ← УБРАТЬ CONTEXT
    }
}

// Config остается без изменений
type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
}