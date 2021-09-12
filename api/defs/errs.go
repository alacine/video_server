package defs

import "net/http"

// Err 错误信息
type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"` // system error code, not http status code
}

// ErrResponse response 错误信息
type ErrResponse struct {
	HTTPSC int
	Error  Err
}

var (
	// ErrorRequestBodyParseFailed 请求体错误
	ErrorRequestBodyParseFailed = ErrResponse{
		HTTPSC: http.StatusBadRequest, // 400
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}

	// ErrorFileSize 上传文件太大
	ErrorFileSize = ErrResponse{
		HTTPSC: http.StatusRequestEntityTooLarge,
		Error: Err{
			Error:     "File size error",
			ErrorCode: "002",
		},
	}

	// ErrorNotAuthUser 未登录，登录过期
	ErrorNotAuthUser = ErrResponse{
		HTTPSC: http.StatusUnauthorized, // 401
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "003",
		},
	}

	// ErrorDBError 数据库操作错误
	ErrorDBError = ErrResponse{
		HTTPSC: http.StatusInternalServerError, // 500
		Error: Err{
			Error:     "DB ops failed",
			ErrorCode: "004",
		},
	}

	// ErrorInternalFaults 内部错误
	ErrorInternalFaults = ErrResponse{
		HTTPSC: http.StatusInternalServerError, // 500
		Error: Err{
			Error:     "Internal service error",
			ErrorCode: "005",
		},
	}
)
