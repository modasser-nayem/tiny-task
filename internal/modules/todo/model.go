package todo

import "time"

type Todo struct {
	ID 						string
	UserID 				string
	Title 				string
	Description 	string
	Completed 	  bool
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}