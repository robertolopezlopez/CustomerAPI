package customer

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"
)

type (
	Customer struct {
		gorm.Model
		Email     string `json:"email,omitempty" gorm:"uniqueIndex:idx_multi"`
		Title     string `json:"title,omitempty" gorm:"uniqueIndex:idx_multi"`
		Content   string `json:"content,omitempty" gorm:"uniqueIndex:idx_multi"`
		MailingID int64  `json:"mailing_id,omitempty" gorm:"uniqueIndex:idx_multi"`
	}
)

func (c *Customer) Validate() error {
	return validation.Errors{
		"email":   validation.Validate(c.Email, validation.Required, is.Email, validation.Length(0, 50)),
		"title":   validation.Validate(c.Title, validation.Required, validation.Length(0, 50)),
		"content": validation.Validate(c.Content, validation.Length(0, 150)),
	}.Filter()
}
