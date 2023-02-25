package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/guid"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func GenGUID() string {
	g := guid.New()
	return g.StringUpper()
}

func GenTransID() string {
	g := guid.New()
	return string(strings.ToLower(strings.Replace(g.StringUpper(), "-", "", -1))[0:16])
}

func TeleTypeToTelecomName(teleType int) string {
	switch teleType {
	case 0:
		return "SKT"
	case 1:
		return "KT"
	case 2:
		return "LGUP"
	}
	return ""
}

func TeleTypeNumToTelecomName(teleType string) string {
	switch teleType {
	case "0":
		return "SKT"
	case "1":
		return "KT"
	case "2":
		return "LGUP"
	}
	return ""
}

func EncodeFile(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

func SendCaptchaData(CaptchaURL string, token string, userIP string) string {

	fmt.Println(makeReqCaptchaData(token, userIP))
	resp, err := http.PostForm(CaptchaURL, makeReqCaptchaData(token, userIP))
	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
	var f []byte
	var strJson string
	if resp != nil {
		f, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println(string(f))

		strJson = strings.Replace(string(f), `"{`, "{", -1)
		strJson = strings.Replace(strJson, `}}"`, `}}`, -1)
		strJson = strings.Replace(strJson, `\`, ``, -1)
		fmt.Println(strJson)

	} else {
		fmt.Println("CAPTCHA RESP NULL")
	}

	return strJson
}

func SendCaptchaDataV3(CaptchaURL string, token string, userIP string) string {

	fmt.Println(makeReqCaptchaDataV3(token, userIP))
	resp, err := http.PostForm(CaptchaURL, makeReqCaptchaDataV3(token, userIP))
	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
	var f []byte
	var strJson string
	if resp != nil {
		f, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println(string(f))

		strJson = strings.Replace(string(f), `"{`, "{", -1)
		strJson = strings.Replace(strJson, `}}"`, `}}`, -1)
		strJson = strings.Replace(strJson, `\`, ``, -1)
		fmt.Println(strJson)

	} else {
		fmt.Println("CAPTCHA RESP NULL")
	}

	return strJson
}

func RestfulSendData(url string, inData []byte) []byte {
	reqData := bytes.NewBuffer(inData)
	resp, err := http.Post(url, "application/json", reqData)
	if err != nil {
		fmt.Println(err.Error())
	}
	var f []byte
	// var strJson strng
	if resp != nil {
		f, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	return f
}

func MakeJsonData(header interface{}, body interface{}) []byte {
	obj := make(map[string]interface{})
	obj["Header"] = header
	if body == nil {
		obj["Body"] = ""
	} else {
		obj["Body"] = body
	}

	retByte, err := json.Marshal(obj)
	if err != nil {
		fmt.Print("rest", err)
	}

	return retByte
}

func MakeJsonDataString(header interface{}, body interface{}) string {
	obj := make(map[string]interface{})
	obj["Header"] = header
	if body == nil {
		obj["Body"] = ""
	} else {
		obj["Body"] = body
	}

	retByte, err := json.Marshal(obj)
	if err != nil {
		fmt.Print("rest", err)
	}

	return string(retByte)
}

// MakeLgupRspData return interface to string
func MakeLgupRspData(body interface{}) string {
	retByte, err := json.Marshal(body)
	if err != nil {
		fmt.Print("rest", err)
	}

	return string(retByte)
}

func makeReqCaptchaData(token string, remoteIP string) url.Values {
	var formData = url.Values{}
	formData.Add("secret", "6LcmCn0UAAAAAFQ7jFK6rjSWDbUvwhLZvF7KVDA2")
	formData.Add("remoteip", remoteIP)
	formData.Add("response", token)
	return formData
}

func makeReqCaptchaDataV3(token string, remoteIP string) url.Values {
	var formData = url.Values{}
	formData.Add("secret", "6LfY-64UAAAAANXUj4McSFrMZ6tuX9VApo1JYl6Y")
	formData.Add("remoteip", remoteIP)
	formData.Add("response", token)
	return formData
}

func MakeReqDataURLEncode(header interface{}, body interface{}) *bytes.Buffer {
	req := make(map[string]interface{})
	req["Header"] = header
	if body == nil {
		req["Body"] = ""
	} else {
		req["Body"] = body
	}

	retByte, err := json.Marshal(req)
	if err != nil {
		fmt.Print("rest", err)
	}

	return bytes.NewBuffer([]byte(b64.URLEncoding.EncodeToString(retByte)))
}

func JsonToHaderBody(jsonOrg []byte, headerObj interface{}, bodyObj interface{}) {
	jsonData := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonOrg), &jsonData); err != nil {
		fmt.Print("ERROR line utilites225", err, jsonOrg)
	}

	headerBytes, _ := json.Marshal(jsonData["Header"])

	if err := json.Unmarshal(headerBytes, &headerObj); err != nil {
		fmt.Print(err)
	}

	bodyBytes, _ := json.Marshal(jsonData["Body"])

	if err := json.Unmarshal(bodyBytes, &bodyObj); err != nil {
		fmt.Print(err)
	}
}

func JsonToHeader(jsonOrg []byte, headerObj interface{}) {
	jsonData := make(map[string]interface{})

	strJson := ""
	strJson = strings.Replace(string(jsonOrg), `"{`, "{", -1)
	strJson = strings.Replace(strJson, `}}"`, `}}`, -1)
	strJson = strings.Replace(strJson, `\"`, `"`, -1)
	strJson = strings.Replace(strJson, `}"`, "}", -1)

	if err := json.Unmarshal([]byte(strJson), &jsonData); err != nil {
		fmt.Print(err)
	}

	headerBytes, _ := json.Marshal(jsonData["Header"])

	if err := json.Unmarshal(headerBytes, &headerObj); err != nil {
		fmt.Print(err)
	}

}

func JsonToBody(jsonOrg []byte, bodyObj interface{}) {
	jsonData := make(map[string]interface{})

	strJson := ""
	strJson = strings.Replace(string(jsonOrg), `"{`, "{", -1)
	strJson = strings.Replace(strJson, `}}"`, `}}`, -1)
	strJson = strings.Replace(strJson, `\"`, `"`, -1)
	strJson = strings.Replace(strJson, `}"`, "}", -1)

	if err := json.Unmarshal([]byte(strJson), &jsonData); err != nil {
		fmt.Print("ERROR line utilites260", err, strJson)
	}

	bodyBytes, _ := json.Marshal(jsonData["Body"])

	if err := json.Unmarshal(bodyBytes, &bodyObj); err != nil {
		fmt.Print("ERROR line utilites266", err, string(bodyBytes))
	}

}

// JSONToBuffer json to url encoded byte buffer
func ObjectToJson(jsonVal interface{}) []byte {
	b, _ := json.Marshal(jsonVal)
	return b
}

// JSONToBuffer json to url encoded byte buffer
func JSONToObject(jsonData []byte, obj interface{}) {

	if err := json.Unmarshal(jsonData, obj); err != nil {
		fmt.Print("ERROR line utilites266", err, string(jsonData))
	}
}

// RequestResByte request uri to tls 12 return byte arrays
func RequestResByte(uri string, data []byte) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}
	d := bytes.NewBuffer([]byte(b64.URLEncoding.EncodeToString(data)))
	req, err := http.NewRequest("POST", uri, d)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("request error : %v\n", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	b, _ := ioutil.ReadAll(res.Body)

	// resBody := commonformats.ResNhLifeInfoStruct{}
	// json.Unmarshal(b, &resBody)
	return b
}

// SetUtf8ToEucKr return euc-kr string
func SetUtf8ToEucKr(str string) string {
	b := bytes.NewReader([]byte(str))
	e := transform.NewReader(b, korean.EUCKR.NewEncoder())
	r, _ := ioutil.ReadAll(e)
	return string(r)
}

// GetBytes inverface to byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetAge input Ymd string to int age
func GetAge(str string) int {
	layout := "20060102"
	b, _ := time.Parse(layout, str)
	now := time.Now()
	age := now.Year() - b.Year()

	bd := b.YearDay()
	cd := now.YearDay()

	if isLeap(b) && !isLeap(now) && bd >= 60 {
		bd = bd - 1
	}
	if isLeap(now) && !isLeap(b) && cd >= 60 {
		bd = bd + 1
	}

	if now.YearDay() < bd {
		age--
	}

	return age
}

func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

// GetRandTranID length is max len return rand number with lpad
func GetRandTranID(length int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var lengthToRandSize string
	for i := 0; i < length; i++ {
		lengthToRandSize += "9"
	}

	randSize, _ := strconv.Atoi(lengthToRandSize)
	randNumber := r1.Intn(randSize)

	return randNumber

}

// SendTelegramMsg Send TelegramMsg ..
// return
// int http status code
// bool result
// int messageID
// if result false https statuscode 404 messageid 0
func SendTelegramMsg(APIKey, Channel, Message string) (int, bool, int, error) {
	type sTelegramBody struct {
		MessageID int    `json:"message_id"`
		Text      string `json:"text"`
		Date      int64  `json:"date"`
	}

	type sTelegramResult struct {
		Ok     bool          `json:"ok"`
		Result sTelegramBody `json:"result"`
	}

	URL := "https://api.telegram.org/bot[BOT_API_KEY]/sendMessage?chat_id=[MY_CHANNEL_NAME]&text=[MY_MESSAGE_TEXT]"

	URL = strings.Replace(URL, "[BOT_API_KEY]", APIKey, -1)
	URL = strings.Replace(URL, "[MY_CHANNEL_NAME]", Channel, -1)
	URL = strings.Replace(URL, "[MY_MESSAGE_TEXT]", Message, -1)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: transCfg}

	req, _ := http.NewRequest("GET", URL, nil)
	rsp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return 0, false, 0, err
	}
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()

	b, err := ioutil.ReadAll(rsp.Body)

	s := sTelegramResult{}
	json.Unmarshal(b, &s)
	return rsp.StatusCode, s.Ok, s.Result.MessageID, nil
}

func StrToTimeSeoul(ori, zone string) time.Time {
	layout := "2006-01-02 15:04:05 Etc/GMT"

	retDT, _ := time.Parse(layout, ori)
	loc, err := time.LoadLocation(zone)
	if err != nil {
		retDT = retDT.Add(time.Hour * 9)
	} else {
		retDT = retDT.In(loc)
	}

	return retDT
}

// GUIDGen hex generator
func GUIDGen(n int) string {
	g := guid.New()
	s := strings.Replace(g.StringUpper(), "-", "", -1)
	if n > len(s) {
		return ""
	}
	return string(strings.ToLower(s)[:n])
}

// UUIDVersion5Gen hex generator
func UUIDVersion5Gen(key string) string {

	id := uuid.NewV3(uuid.NamespaceURL, key)
	return strings.ToUpper(id.String())
}

// MatterRequestBody struct
type MatterRequestBody struct {
	Text string `json:"text"`
}

// JandiRequestBody struct
type JandiRequestBody struct {
	Body         string                `json:"body"`
	ConnectColor string                `json:"ConnectColor"`
	ConnectInfo  []JandiRequestConnect `json:"connectInfo"`
}

// JandiRequestConnect struct
type JandiRequestConnect struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
// SendSlackNotification("http://mm.datau.co.kr:8065/hooks/qzcho4t5s7n63m3gp6598k1mxy", Message)
func SendSlackNotification(msg string) error {
	webhookURL := "http://mm.datau.co.kr:8065/hooks/qzcho4t5s7n63m3gp6598k1mxy"
	slackBody, _ := json.Marshal(MatterRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

// SendSlackNotificationCS("http://mm.datau.co.kr:8065/hooks/bwpygwgsjjrt7yrthq9cbyhywa")
func SendSlackNotificationCS(msg string) error {
	webhookURL := "http://mm.datau.co.kr:8065/hooks/bwpygwgsjjrt7yrthq9cbyhywa"
	slackBody, _ := json.Marshal(MatterRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

// SendJandiNotification webhook for jandi
// https://wh.jandi.com/connect-api/webhook/19828284/f070f266603cbc23f8e83352e89ecd6c -- cs
// https://wh.jandi.com/connect-api/webhook/19828284/f39566db2baffb3d2f5dcd323d6e4d7c -- monitoring
func SendJandiNotification(jandiURL, msg string) error {
	slackBody, _ := json.Marshal(JandiRequestBody{Body: msg})
	req, err := http.NewRequest(http.MethodPost, jandiURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.tosslab.jandi-v2+json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// downloadXML from git server
func DownloadFile(downloadURL string, saveFilePath string) {

	rsp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Print("Download File 1", "", err)

		panic(err)
	}
	defer rsp.Body.Close()

	file, err := os.Create(saveFilePath)
	if err != nil {
		fmt.Print("Download File 2", "", err)

		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		fmt.Print("Download File 3", "", err)
	}
}
