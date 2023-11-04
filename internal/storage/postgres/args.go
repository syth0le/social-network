package postgres

import "fmt"

const (
	UserTable  = "user_table"
	TokenTable = "token_table"
)

const (
	returning = "RETURNING"
	separator = ","
)

const (
	fieldUserID = "user_id"

	fieldUsername   = "username"
	fieldFirstName  = "first_name"
	fieldSecondName = "second_name"
	fieldAge        = "age"
	fieldSex        = "sex"
	fieldBirthdate  = "birthdate"
	fieldBiography  = "biography"
	fieldCity       = "city"

	fieldCreatedAt = "created_at"
	fieldDeletedAt = "deleted_at"

	fieldToken    = "token"
	fieldAlivedAt = "alived_at"
)

var (
	userFields = []string{
		fieldUserID, fieldUsername, fieldFirstName, fieldSecondName,
		fieldAge, fieldSex, fieldBirthdate, fieldBiography, fieldCity,
		fieldCreatedAt, fieldDeletedAt,
	}
	tokenFields = []string{fieldUserID, fieldToken, fieldCreatedAt, fieldDeletedAt, fieldAlivedAt}
)

func tableField(table, field string) string {
	return fmt.Sprintf("%s.%s", table, field)
}
