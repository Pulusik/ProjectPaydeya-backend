package handlers

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
    //"github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)


// Временная функция для тестирования БД
func GetUsersTest(db *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        rows, err := db.Query(context.Background(), "SELECT id, email, full_name, role FROM users")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        type User struct {
            ID       int    `json:"id"`
            Email    string `json:"email"`
            FullName string `json:"fullName"`
            Role     string `json:"role"`
        }

        var users []User
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.Email, &user.FullName, &user.Role); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            users = append(users, user)
        }

        c.JSON(http.StatusOK, gin.H{
            "users": users,
            "total": len(users),
        })
    }
}