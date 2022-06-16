package db

import (
	"aig-tech-okr/libs"
	"crypto/tls"
	"crypto/x509"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net"
)

var sessionMap map[string]*mgo.Session

type MongoLog struct {
}

func (MongoLog) Output(calldepth int, s string) error {
	log.SetFlags(log.Lshortfile)
	return log.Output(calldepth, s)
}

func RegisterMongo() {

	sessionMap = make(map[string]*mgo.Session, len(libs.Conf.Mongo))
	if libs.Conf.App.Env != "pro" {
		//mgo.SetDebug(true)
	}

	//mgo.SetLogger(MongoLog{})

	for key, mongoConf := range libs.Conf.Mongo {
		var session *mgo.Session
		if libs.Conf.App.Env == "pro" {
			session = mongo(key)
		} else {
			session = mongo(key)
		}

		session.SetPoolLimit(mongoConf.MaxPoolSize)
		session.SetMode(mgo.Monotonic, true)

		sessionMap[mongoConf.InsName] = session

	}
}
func mongo(index int) *mgo.Session {
	session, err := mgo.Dial(libs.Conf.Mongo[index].Addr)
	if err != nil {
		panic("连接mongo失败-" + libs.Conf.Mongo[index].InsName + ",err:" + err.Error())
	}
	return session
}
func certMongo(index int) *mgo.Session {

	b, err := ioutil.ReadFile(libs.Conf.Mongo[index].PemPath) // just pass the file name
	if err != nil {
		panic("连接cert_mongo失败-读取证书失败-" + libs.Conf.Mongo[index].InsName + ",err:" + err.Error())
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(b)
	if !ok {
		panic("连接cert_mongo失败-解析证书失败-" + libs.Conf.Mongo[index].InsName)
	}

	tlsConfig := &tls.Config{
		RootCAs:            roots,
		InsecureSkipVerify: true,
	}

	dialInfo, err := mgo.ParseURL(libs.Conf.Mongo[index].PemPath)
	if err != nil {
		panic("连接cert_mongo失败-解析dialInfo失败-" + libs.Conf.Mongo[index].InsName + ",err:" + err.Error())
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			panic("连接cert_mongo失败-连接失败-" + libs.Conf.Mongo[index].InsName + ",err:" + err.Error())
		}
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic("连接cert_mongo失败-连接失败-" + libs.Conf.Mongo[index].InsName + ",err:" + err.Error())
	}

	return session

}

// 获取实例
func GetMongoIns(insName string) *mgo.Session {

	if session, ok := sessionMap[insName]; ok {
		return session
	}

	panic("no such mongo instance conf [" + insName + "]")
}
