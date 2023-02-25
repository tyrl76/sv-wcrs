package formats

const (
	Login  = 1000
	Logout = 1001
)

type RecvHeader struct {
	CmdType   int    `json:"CmdType"`
	RequestID string `json:"RequestID"`
}

type SendHeader struct {
	CmdType int    `json:"CmdType"`
	ErrCode int    `json:"ErrCode"`
	ErrMsg  string `json:"ErrMsg"`
}
