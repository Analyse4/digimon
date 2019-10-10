package errorhandler

import "fmt"

const (
	SUCESS = iota
	ERR_SERVICEBUSY
	ERR_USERNOTLOGIN
	ERR_PARAMETERINVALID
	ERR_SKILLPOINTNOTENOUGH
)

var (
	ERR_PARAMETERINVALID_MSG    = fmt.Errorf("parameter invalid")
	ERR_SKILLPOINTNOTENOUGH_MSG = fmt.Errorf("skill point not enough")
	ERR_SERVICEBUSY_MSG         = fmt.Errorf("service busy")
)

func GetErrMsg(errcode int64) string {
	switch errcode {
	case SUCESS:
		return fmt.Sprintf("success")
	case ERR_SERVICEBUSY:
		return ERR_SERVICEBUSY_MSG.Error()
	case ERR_USERNOTLOGIN:
		return fmt.Sprintf("user not login")
	case ERR_PARAMETERINVALID:
		return ERR_PARAMETERINVALID_MSG.Error()
	case ERR_SKILLPOINTNOTENOUGH:
		return ERR_SKILLPOINTNOTENOUGH_MSG.Error()
	default:
		return fmt.Sprintf("can't find err code")
	}
}
