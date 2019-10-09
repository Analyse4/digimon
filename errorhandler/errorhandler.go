package errorhandler

import "fmt"

const (
	SUCESS = iota
	ERR_SERVICEBUSY
	ERR_USERNOTLOGIN
	ERR_PARAMETERINVALID
)

func GetErrMsg(errcode int64) string {
	switch errcode {
	case SUCESS:
		return fmt.Sprintf("success")
	case ERR_SERVICEBUSY:
		return fmt.Sprintf("service busy")
	case ERR_USERNOTLOGIN:
		return fmt.Sprintf("user not login")
	case ERR_PARAMETERINVALID:
		return fmt.Sprintf("parameter invalid")
	default:
		return fmt.Sprintf("can't find err code")
	}
}
