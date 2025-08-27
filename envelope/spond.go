package envelope

type ErrorDetail struct {
	Title    string
	Message  string
	Solution string
}

type Error struct {
	Code   StatusCode
	Detail ErrorDetail
}
