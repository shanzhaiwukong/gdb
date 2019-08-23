package main

import (
	"fmt"
	"gdb/mongo"
	"time"
)

func main() {
	tests()
	select {}
}

func tests() {
	db, err := mongo.New(&mongo.Config{
		Host:      "192.168.80.130:27017",
		User:      "haha1",
		Pwd:       "123456",
		DbName:    "haha",
		Source:    "haha",
		PoolLimit: 2,
		Timeout:   time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
	close := make(chan bool)
	err = db.C("news", close).Insert(&struct{ Title string }{"123"})
	err = db.C("news", close).Insert(&struct{ Title string }{"DEF"})
	close <- true
	fmt.Println(err)
}
