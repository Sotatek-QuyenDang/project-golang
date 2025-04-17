package models

type User struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;unique_index"`
	Image    string `json:"image"`
	Status   int    `json:"status"`
	Possible int    `json:"possible"`
}
