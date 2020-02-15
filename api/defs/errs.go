package defs

import "net/http"

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"` // system error code, not http status code
}

type ErrResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{
		HttpSC: http.StatusBadRequest, // 400
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}
	ErrorNotAuthUser = ErrResponse{
		HttpSC: http.StatusUnauthorized, // 401
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "002",
		},
	}
	ErrorDBError = ErrResponse{
		HttpSC: http.StatusInternalServerError, // 500
		Error: Err{
			Error:     "DB ops failed",
			ErrorCode: "003",
		},
	}
	ErrorInternalFaults = ErrResponse{
		HttpSC: http.StatusInternalServerError, // 500
		Error: Err{
			Error:     "Internal service error",
			ErrorCode: "004",
		},
	}
)
