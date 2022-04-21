package models

type Url struct {
	BaseModel
	Owner             uint   `gorm:"owner"`
	OriginalURL       string `gorm:"column:original_url"`
	ShortenedURLParam string `gorm:"column:shortened_url_param"`
	VisitCount        *int   `gorm:"column:visit_count"`
}
