package handler

import (
	"io/ioutil"
	"silvernote/factory"
	"silvernote/formats"
	"silvernote/utils"
	"time"

	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	Fac *factory.Factory
}

func (_self *ServiceHandler) RequestHandle(con echo.Context) error {

	var rspData []byte
	var rspHeader formats.SendHeader
	var rspBody interface{}

	body := con.Request().Body
	x, _ := ioutil.ReadAll(body)
	json := string(x)

	var reqHeader formats.RecvHeader

	utils.JsonToHeader(x, &reqHeader)
	reqStartTime := time.Now()

	_self.Fac.Print(reqHeader.RequestID, "Request Start", json)

	rspHeader = formats.SendHeader{
		CmdType: reqHeader.CmdType,
		ErrCode: formats.Success,
		ErrMsg:  formats.GetMsg(formats.Success),
	}

	switch reqHeader.CmdType {
	default:
		_self.Fac.Print(reqHeader.RequestID, "Not Allowed CmdType", json)
		rspHeader = formats.SendHeader{
			CmdType: reqHeader.CmdType,
			ErrCode: formats.ErrorCmdType,
			ErrMsg:  formats.GetMsg(formats.ErrorCmdType),
		}
	}

	rspData = utils.MakeJsonData(reqHeader, rspBody)

	timedu := time.Since(reqStartTime)

	_self.Fac.Print(reqHeader.RequestID, "소요시간 :", timedu, "Result :", rspHeader)

	con.String(200, string(rspData))

	return nil
}
