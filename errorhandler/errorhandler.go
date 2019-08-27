package errorhandler

import "fmt"

const (
	SUCESS = iota
	ERR_SERVICEBUSY
)

func GetErrMsg(errcode int64) string {
	switch errcode {
	case SUCESS:
		return fmt.Sprintf("success")
	case ERR_SERVICEBUSY:
		return fmt.Sprintf("service busy")
	default:
		return fmt.Sprintf("can't find err code")
	}
}
