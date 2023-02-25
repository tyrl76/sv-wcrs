package formats

var messages map[int]string

const (
	Success = 0

	ErrorCmdType = 9999
)

func init() {
	messages = make(map[int]string)
	messages[Success] = "정상"
	messages[ErrorCmdType] = "시스템 오류(CmdType)"
}

func GetMsg(ErrCode int) string {
	return messages[ErrCode]
}

func (self *SendHeader) SetHeader(errCode int, errMsg string) {
	self.ErrCode = errCode
	if errMsg == "" {
		self.ErrMsg = GetMsg(errCode)
	} else {
		self.ErrMsg = errMsg
	}
}
