package models

// CatalogMaterial represents material in catalog
// @Description Материал в каталоге
type CatalogMaterial struct {
    ID            int     `json:"id" example:"1"`
    Title         string  `json:"title" example:"Основы алгебры"`
    Subject       string  `json:"subject" example:"math"`
    Author        Author  `json:"author"`
    Rating        float64 `json:"rating" example:"4.8"`
    StudentsCount int     `json:"studentsCount" example:"150"`
    Duration      int     `json:"duration,omitempty" example:"120"`
    Level         string  `json:"level,omitempty" example:"beginner"`
    ThumbnailURL  string  `json:"thumbnailUrl,omitempty" example:"https://example.com/thumbnail.jpg"`
}

// Author represents material author
// @Description Автор материала
type Author struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"Иван Иванов"`
}

// Teacher represents teacher in catalog
// @Description Преподаватель в каталоге
type Teacher struct {
    ID              int      `json:"id" example:"1"`
    Name            string   `json:"name" example:"Мария Петрова"`
    Specializations []string `json:"specializations" example:"math,physics"`
    Rating          float64  `json:"rating" example:"4.9"`
    MaterialsCount  int      `json:"materialsCount" example:"25"`
    AvatarURL       *string   `json:"avatarUrl,omitempty" example:"https://example.com/avatar.jpg"`
}

// Subject represents subject/course
// @Description Учебный предмет
type Subject struct {
    ID   string `json:"id" example:"math"`
    Name string `json:"name" example:"Математика"`
    Icon string `json:"icon" example:"📐"`
}

// CatalogFilters represents filters for materials search
// @Description Фильтры для поиска материалов
type CatalogFilters struct {
    Search  string `form:"search" example:"алгебра"`
    Subject string `form:"subject" example:"math"`
    Level   string `form:"level" example:"beginner"`
    Page    int    `form:"page" example:"1"`
    Limit   int    `form:"limit" example:"20"`
}

// TeacherFilters represents filters for teachers search
// @Description Фильтры для поиска преподавателей
type TeacherFilters struct {
    Search  string `form:"search" example:"математика"`
    Subject string `form:"subject" example:"math"`
}