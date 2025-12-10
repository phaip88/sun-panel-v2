package models

type Notepad struct {
	BaseModel
	UserID  uint   `json:"userId" gorm:"index;comment:用户ID"`
	Content string `json:"content" gorm:"type:text;comment:便签内容"`
}
