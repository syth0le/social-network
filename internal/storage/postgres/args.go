package postgres

import (
	"fmt"
	"strings"
)

const (
	UserTable   = "user_table"
	TokenTable  = "token_table"
	FriendTable = "friend_table"
	PostTable   = "post_table"
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
	fieldUpdatedAt = "updated_at"
	fieldDeletedAt = "deleted_at"

	fieldToken    = "token"
	fieldAlivedAt = "alived_at"

	fieldText = "text"

	fieldFirstUserID  = "first_user_id"
	fieldSecondUserID = "second_user_id"
)

var (
	userFields = []string{
		fieldID, fieldUsername, fieldHashedPassword, fieldFirstName, fieldSecondName,
		fieldSex, fieldBirthdate, fieldBiography, fieldCity, fieldCreatedAt,
	}
	tokenFields  = []string{fieldID, fieldUserID, fieldToken, fieldCreatedAt, fieldAlivedAt}
	friendFields = []string{fieldID, fieldFirstUserID, fieldSecondUserID, fieldCreatedAt}
	postFields   = []string{fieldID, fieldUserID, fieldText, fieldCreatedAt, fieldUpdatedAt}

	returningUser  = returning + strings.Join(userFields, separator)
	returningToken = returning + strings.Join(tokenFields, separator)
	returningPost  = returning + strings.Join(postFields, separator)
)

func tableField(table, field string) string {
	return fmt.Sprintf("%s.%s", table, field)
}

func tableFields(table string, fields []string) []string {
	var respFields []string
	for _, field := range fields {
		respFields = append(respFields, tableField(table, field))
	}
	return respFields
}

func mergeFields(firstFields []string, secondFields ...string) []string {
	for _, field := range secondFields {
		firstFields = append(firstFields, field)
	}
	return firstFields
}

func joinString(sourceTable, sourceField, joinTable, joinField string) string {
	return fmt.Sprintf("%[1]s ON %[2]s.%[3]s = %[1]s.%[4]s", joinTable, sourceTable, sourceField, joinField)
}
