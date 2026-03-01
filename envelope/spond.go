package envelope

type ErrorDetail struct {
	Title    string
	Message  string
	Solution string
}

type AppError struct {
	Code   int
	Detail ErrorDetail
}
