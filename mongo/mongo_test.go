package mongo

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type News struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Title    string
	Author   string
	Isorigin bool `bson:"isOrigin"`
}

// 单元测试
func TestA(t *testing.T) {
	i := 3000
	for {
		if i < 1 {
			break
		}
		i--
		_, err1 := New(&Config{
			Host:    "127.0.0.1:27017",
			User:    "test",
			Pwd:     "123456",
			DbName:  "blog",
			Source:  "admin",
			Timeout: time.Second * 10,
		})
		fmt.Println(err1)
	}
}

// 测试并发写数据能力
func BenchmarkWrite(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		mg, _ := New(&Config{
			Host:    "127.0.0.1:27017",
			User:    "test",
			Pwd:     "123456",
			DbName:  "blog",
			Source:  "admin",
			Timeout: time.Second * 10,
		})
		for pb.Next() {
			isClose := make(chan bool)
			err := mg.C("news", isClose).Insert(&News{Title: bson.NewObjectId().Hex()})
			isClose <- true
			if err != nil {
				panic(err)
			}
		}
	})
}

//测试并发连接写数据
func BenchmarkConnWrite(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if mg, err := New(&Config{
				Host:    "127.0.0.1:27017",
				User:    "test",
				Pwd:     "123456",
				DbName:  "blog",
				Source:  "admin",
				Timeout: time.Second * 10,
			}); err == nil {
				defer mg.Close()
				isClose := make(chan bool)
				err := mg.C("news", isClose).Insert(&News{Title: bson.NewObjectId().Hex()})
				isClose <- true
				if err != nil {
					panic(err)
				}
			}
		}
	})
}

/*
执行TestX测试程序
-run 需要测试的函数名，注意和性能测试的文件名区分
-v 输出测试过程信息
-count 执行次数
go test -run TestX -v -count 1
执行所有单元测试
-count 执行次数
go test -v -count 1
*/

/*
执行BenchmarkX测试程序
-run 需要测试的.go文件名，如 main.go则文件名为main
-bench 需要测试的函数名
-count 执行次数
go test -run 文件名 -bench=BenchmarkX -count=3
运行所有的性能测试函数
-count 执行次数
go test -test.bench=".*" -count=3
*/
