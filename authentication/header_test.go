package authentication

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type mockHeader struct {
	mock.Mock
}

func (m *mockHeader) Login() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestHeaderAuth_Login(t *testing.T) {
	tests := map[string]struct {
		input  HeaderAuth
		output bool
	}{
		"happy auth path": {
			input: HeaderAuth{
				Name:  "X-Token",
				Value: "test",
			},
			output: true,
		},
		"wrong value": {
			input: HeaderAuth{
				Name:  "X-Token",
				Value: "-",
			},
		},
		"no x-token": {
			input: HeaderAuth{
				Name:  "-",
				Value: "test",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.input.Login()
			assert.Equal(t, test.output, got)
		})
	}
}
