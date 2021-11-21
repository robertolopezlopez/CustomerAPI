package customer

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockCustomer struct {
	mock.Mock
}

func TestCustomer_Validate(t *testing.T) {
	tests := map[string]struct {
		c         Customer
		withError *regexp.Regexp
	}{
		"ok": {
			c: Customer{
				Email: "oroparece@platano.es",
				Title: "a title",
			},
		},
		"long content": {
			c: Customer{
				Email: "oroparece@platano.es",
				Title: "a title",
				Content: "contentcontentcontentcontentcontentcontentcontentcontent" +
					"contentcontentcontentcontentcontentcontentcontentcontentcontent" +
					"contentcontentcontentcontentcontentcontentcontentcontentcontent" +
					"contentcontentcontentcontentcontentcontentcontentcontentcontent",
			},
			withError: regexp.MustCompile("content: the length must be no more than 150"),
		},
		"no title": {
			c: Customer{
				Email: "oroparece@platano.es",
			},
			withError: regexp.MustCompile("title: cannot be blank"),
		},
		"no email": {
			c: Customer{
				Title: "a title",
			},
			withError: regexp.MustCompile("email: cannot be blank"),
		},
		"long email": {
			c: Customer{
				Email: "oropareceoropareceoropareceoropareceoropareceoropareceoropareceoropareceoroparece@platano.es",
				Title: "a title",
			},
			withError: regexp.MustCompile("email: the length must be no more than 50"),
		},
		"long title": {
			c: Customer{
				Email: "oroparece@platano.es",
				Title: "oropareceoropareceoropareceoropareceoropareceoropareceoropareceoropareceoroparece@platano.es",
			},
			withError: regexp.MustCompile("title: the length must be no more than 50"),
		},
		"invalid email": {
			c: Customer{
				Email: "not an email",
				Title: "a title",
			},
			withError: regexp.MustCompile("email: must be a valid email address"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.c.Validate()
			if test.withError != nil {
				require.Error(t, err)
				assert.Regexp(t, test.withError, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
