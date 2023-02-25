package factory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"silvernote/utils"

	"github.com/sirupsen/logrus"
	cronowriter "github.com/utahta/go-cronowriter"
)

type Config struct {
	Database struct {
		Driver     string `json:"Driver"`
		Connection string `json:"Connection"`
		PoolSize   int    `json:"PoolSize"`
	} `json:"Database"`
	Port            string `json:"Port"`
	LogerFilePath   string `json:"LogerFilePath"`
	MessageFilePath string `json:"MessageFilePath"`
	Debug           bool   `json:"Debug"`
	AppName         string `json:"AppName"`
}

type Factory struct {
	db             *sql.DB
	logger         *logrus.Logger
	JSONConfigPath string
	JSONConfigURL  string
	property       Config
	ConfigMap      map[string]interface{}
	MessageMap     map[string]string
	AppEnv         *string
}

func (_self *Factory) InitDB(driver, connection string, poolSize int) (*sql.DB, error) {
	db, err := sql.Open(driver, connection)
	db.SetMaxIdleConns(poolSize)
	db.SetMaxOpenConns(poolSize)

	if err != nil {
		panic(err)
		// return nil, err
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	return db, nil
}

func (_self *Factory) loadConfiguration(file string) {

	if _self.JSONConfigURL != "" {
		utils.DownloadFile(_self.JSONConfigURL, file)
	}

	c, _ := ioutil.ReadFile(file)
	fmt.Print(string(c))

	json.Unmarshal(c, &_self.property)

	_self.ConfigMap = make(map[string]interface{})
	json.Unmarshal(c, &_self.ConfigMap)

	m, _ := ioutil.ReadFile(_self.property.MessageFilePath)
	fmt.Print(string(m))

	_self.MessageMap = make(map[string]string)
	json.Unmarshal(m, &_self.MessageMap)

}

func (_self *Factory) Initialize() {

	writer := cronowriter.MustNew(_self.property.LogerFilePath)
	mw := io.MultiWriter(os.Stdout, writer)
	logrus.SetOutput(mw)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"

	_self.logger = logrus.New()
	_self.logger.Formatter = new(logrus.TextFormatter)
	_self.logger.SetFormatter(customFormatter)
	_self.logger.Level = logrus.DebugLevel
	_self.logger.Out = mw

}

func (_self *Factory) Logger() *logrus.Logger {
	return _self.logger
}

func (_self *Factory) Print(ReqID string, v ...interface{}) {
	v = append(v[:1], v[0+1:]...)
	_self.logger.Print(v)
}
