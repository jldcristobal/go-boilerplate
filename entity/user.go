package entity

type User struct {
	ID       uint64  `gorm:"primary_key:auto_increment" json:"-"`
	Name     string  `gorm:"type:varchar(255)" json:"-"`
	Email    string  `gorm:"uniqueIndex;type:varchar(255)" json:"-"`
	Password string  `gorm:"->;<-;not null" json:"-"`
	Token    string  `gorm:"-" json:"token,omitempty"`
	Books    *[]Book `json:"books,omitempty"`
}
