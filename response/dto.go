package response

import (
	"fmt"

	"github.com/wxnacy/bdpan"
)

type ApiError struct {
	StatusCode       int    `json:"status_code"`
	Errno            int32  `json:"errno,omitempty"`
	ErrMsg           string `json:"errmsg"`
	Erro             string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
}

func (e *ApiError) Error() string {
	// return fmt.Sprintf("%d: %s(%s)", e.StatusCode, e.Erro, e.ErrorDescription)
	return e.String()
}

func (e *ApiError) String() string {
	if e.ErrorCode > 0 {
		return fmt.Sprintf("%d[%s]", e.ErrorCode, e.ErrorMsg)
	} else if e.Erro != "" {
		return fmt.Sprintf("%s[%s]", e.Erro, e.ErrorDescription)
	} else {
		switch e.Errno {
		case -9:
			return bdpan.ErrPathNotFound.Error()
		case -6:
			return bdpan.ErrAccessFail.Error()
		case 2:
			return bdpan.ErrParamError.Error()
		case 6:
			return bdpan.ErrUserNoUse.Error()
		case 12:
			return bdpan.ErrPathExists.Error()
		case 111:
			return bdpan.ErrAccessTokenFail.Error()
		case 31034, 9013, 9019:
			return bdpan.ErrApiFrequent.Error()
		default:
			return fmt.Sprintf("未知错误: %d", e.Errno)
		}
	}
}
