package handlers

import (
    "net/http"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type CatalogHandler struct {
    catalogService *services.CatalogService
}

func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
    return &CatalogHandler{catalogService: catalogService}
}

// SearchMaterials godoc
// @Summary Поиск материалов в каталоге
// @Description Возвращает материалы с фильтрацией и пагинацией
// @Tags catalog
// @Accept json
// @Produce json
// @Param search query string false "Поисковый запрос"
// @Param subject query string false "Фильтр по предмету"
// @Param level query string false "Фильтр по уровню сложности"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество материалов на странице" default(20)
// @Success 200 {object} MaterialsResponse "Список материалов"
// @Failure 400 {object} ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /catalog/materials [get]
func (h *CatalogHandler) SearchMaterials(c *gin.Context) {
    var filters models.CatalogFilters

    // Парсим query параметры
    if err := c.ShouldBindQuery(&filters); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Устанавливаем значения по умолчанию для пагинации
    if filters.Page == 0 {
        filters.Page = 1
    }
    if filters.Limit == 0 {
        filters.Limit = 20
    }

    materials, total, err := h.catalogService.SearchMaterials(c.Request.Context(), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search materials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "materials": materials,
        "total":     total,
        "page":      filters.Page,
        "limit":     filters.Limit,
        "hasMore":   (filters.Page * filters.Limit) < total,
    })
}

// GetSubjects godoc
// @Summary Получить список предметов
// @Description Возвращает все доступные учебные предметы
// @Tags catalog
// @Accept json
// @Produce json
// @Success 200 {object} SubjectsResponse "Список предметов"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /catalog/subjects [get]
func (h *CatalogHandler) GetSubjects(c *gin.Context) {
    subjects, err := h.catalogService.GetSubjects(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subjects"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "subjects": subjects,
    })
}

// SearchTeachers godoc
// @Summary Поиск преподавателей
// @Description Возвращает преподавателей с фильтрацией
// @Tags catalog
// @Accept json
// @Produce json
// @Param search query string false "Поисковый запрос"
// @Param subject query string false "Фильтр по предмету"
// @Success 200 {object} TeachersResponse "Список преподавателей"
// @Failure 400 {object} ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /catalog/teachers [get]
func (h *CatalogHandler) SearchTeachers(c *gin.Context) {
    var filters models.TeacherFilters

    if err := c.ShouldBindQuery(&filters); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    teachers, err := h.catalogService.SearchTeachers(c.Request.Context(), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search teachers"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "teachers": teachers,
    })
}

// Response models for Swagger

// MaterialsResponse represents materials search response
// @Description Ответ с результатами поиска материалов
type MaterialsResponse struct {
    Materials []models.CatalogMaterial `json:"materials"`
    Total     int                      `json:"total" example:"150"`
    Page      int                      `json:"page" example:"1"`
    Limit     int                      `json:"limit" example:"20"`
    HasMore   bool                     `json:"hasMore" example:"true"`
}

// SubjectsResponse represents subjects list response
// @Description Ответ со списком предметов
type SubjectsResponse struct {
    Subjects []models.Subject `json:"subjects"`
}

// TeachersResponse represents teachers search response
// @Description Ответ с результатами поиска преподавателей
type TeachersResponse struct {
    Teachers []models.Teacher `json:"teachers"`
}

