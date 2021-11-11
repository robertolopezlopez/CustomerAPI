package customer

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Email      string `json:"email,omitempty"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	MailingID  int64  `json:"mailing_id,omitempty"`
	InsertTime string `json:"insert_time,omitempty"`
}

// todo replace ozzo-validation
// https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go
func (c *Customer) Validate() error {
	return validation.Errors{
		"email":   validation.Validate(c.Email, validation.Length(0, 50)),
		"title":   validation.Validate(c.Title, validation.Length(0, 50)),
		"content": validation.Validate(c.Content, validation.Length(0, 150)),
	}.Filter()
}
