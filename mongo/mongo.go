package mongo

import (
	"errors"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type (
	// Mongo 数据库
	Mongo struct {
		session *mgo.Session
		dbName  string
	}
	// Config 配置
	Config struct {
		Host      string
		User      string
		Pwd       string
		DbName    string
		Source    string
		PoolLimit int
		Timeout   time.Duration
	}
)

// New 实例化MongoDB
func New(config *Config) (*Mongo, error) {
	if config == nil {
		return nil, errors.New("config error")
	}
	if sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{config.Host},
		Username:  config.User,
		Password:  config.Pwd,
		Database:  config.DbName,
		Source:    config.Source,
		Timeout:   config.Timeout,
		PoolLimit: config.PoolLimit,
		Direct:    true,
	}); err != nil {
		return nil, err
	} else {
		return &Mongo{session: sess, dbName: config.DbName}, nil
	}
}

// Switch 切换数据库
func (m *Mongo) Switch(dbName string) *Mongo {
	m.dbName = dbName
	return m
}

// Session 底层连接
func (m *Mongo) Session() *mgo.Session {
	return m.session
}

// Close 关闭连接
func (m *Mongo) Close() {
	if m.session != nil {
		m.session.Close()
	}
}

// C Collection
func (m *Mongo) C(tab string, close chan bool) *mgo.Collection {
	cs := m.session.Copy()
	go func() {
		<-close
		cs.Close()
	}()
	c := cs.DB(m.dbName).C(tab)
	cs.SetMode(mgo.Monotonic, true)
	return c
}

// Strong
// session 的读写一直向主服务器发起并使用一个唯一的连接，因此所有的读写操作完全的一致。
// Monotonic
// session 的读操作开始是向其他服务器发起（且通过一个唯一的连接），只要出现了一次写操作，session 的连接就会切换至主服务器。由此可见此模式下，能够分散一些读操作到其他服务器，但是读操作不一定能够获得最新的数据。
// Eventual
// session 的读操作会向任意的其他服务器发起，多次读操作并不一定使用相同的连接，也就是读操作不一定有序。session 的写操作总是向主服务器发起，但是可能使用不同的连接，也就是写操作也不一定有序。
