package postgres

import (
	"fmt"
	"strings"
)

const (
	UserTable  = "user_table"
	TokenTable = "token_table"
)

const (
	returning = "RETURNING "
	separator = ","
)

const (
	fieldID     = "id"
	fieldUserID = "user_id"

	fieldUsername       = "username"
	fieldHashedPassword = "hashed_password"
	fieldFirstName      = "first_name"
	fieldSecondName     = "second_name"
	fieldSex            = "sex"
	fieldBirthdate      = "birthdate"
	fieldBiography      = "biography"
	fieldCity           = "city"

	fieldCreatedAt = "created_at"
	fieldDeletedAt = "deleted_at"

	fieldToken    = "token"
	fieldAlivedAt = "alived_at"
)

var (
	userFields = []string{
		fieldID, fieldUsername, fieldHashedPassword, fieldFirstName, fieldSecondName,
		fieldSex, fieldBirthdate, fieldBiography, fieldCity, fieldCreatedAt,
	}
	tokenFields = []string{fieldID, fieldUserID, fieldToken, fieldCreatedAt, fieldAlivedAt}

	returningUser  = returning + strings.Join(userFields, separator)
	returningToken = returning + strings.Join(tokenFields, separator)
)

func tableField(table, field string) string {
	return fmt.Sprintf("%s.%s", table, field)
}
