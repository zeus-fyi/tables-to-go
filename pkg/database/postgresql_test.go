package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeus-fyi/tables-to-go/pkg/settings"
)

func TestPostgresql_DSN(t *testing.T) {
	tests := []struct {
		desc     string
		settings func() *settings.Settings
		expected func(*settings.Settings) string
	}{
		{
			desc: "no username given, defaults to `postgres`",
			settings: func() *settings.Settings {
				s := settings.New()
				s.DbType = settings.DBTypePostgresql
				return s
			},
			expected: func(s *settings.Settings) string {
				return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
					s.Host, s.Port, "postgres", s.DbName, s.Pswd)
			},
		},
		{
			desc: "with given username, default gets overwritten",
			settings: func() *settings.Settings {
				s := settings.New()
				s.DbType = settings.DBTypePostgresql
				s.User = "my_custom_user"
				return s
			},
			expected: func(s *settings.Settings) string {
				return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
					s.Host, s.Port, "my_custom_user", s.DbName, s.Pswd)
			},
		},
		{
			desc: "with given username and socket, default gets overwritten",
			settings: func() *settings.Settings {
				s := settings.New()
				s.DbType = settings.DBTypePostgresql
				s.User = "my_custom_user"
				s.Pswd = "mysecretpassword"
				s.Socket = "/tmp"
				return s
			},
			expected: func(s *settings.Settings) string {
				return "host=/tmp user=my_custom_user dbname=postgres password=mysecretpassword"
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			db := NewPostgresql(test.settings())
			actual := db.DSN()
			assert.Equal(t, test.expected(db.Settings), actual)
		})
	}
}
