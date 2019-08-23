# GDB

>##  Mongo
> use package of open source： gopkg.in/mgo.v2  
> go get github.com/shanzhaiwukong/gdb/mongo
 ``` go
 //Demo 示例

func main() {
	test()
	select {}
}

func test() {
	db, err := mongo.New(&mongo.Config{
		Host:      "127.0.0.1:27017",
		User:      "test",
		Pwd:       "123456",
		DbName:    "blog",
		Source:    "admin",
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
 ```