package BaseToken

//GetFilter 脱敏设计
func GetFilter(token FieldsTokenType) FieldsTokenType {
	token.Key = "***"
	return token
}
