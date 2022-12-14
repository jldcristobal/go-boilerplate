package entity

type Book struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"-"`
	Title       string `gorm:"type:varchar(255)" json:"-"`
	Description string `gorm:"type:text" json:"-"`
	UserID      uint64 `gorm:"not null" json:"-"`
	User        User   `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE, onDelete:CASCADE" json:"-"`
}
