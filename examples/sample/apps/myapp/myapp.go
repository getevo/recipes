package myapp

import (
	"context"
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/apps/redis"
	"github.com/getevo/evo-ng/lib/generic"
	"github.com/getevo/evo-ng/websocket"
	"github.com/getevo/examples/sample/http"
	"log"
	"os"
	"time"
)



func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	evo.RegisterView("myapp", "./apps/myapp/views")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			message <- fmt.Sprint(time.Now().Unix())
		}
	}()


	return nil
}

var group = http.Group("/a")
var message = make(chan string)

func Router() error {
	http.Get("/", func(context *http.Context) error {
		context.WriteResponse("Hey!")
		return nil
	})

	http.Use("/test", func(context *http.Context) error {

		return fmt.Errorf("you ask for /test/*")
	})
	http.Asset("/asset", "./assets")

	group.Get("/b", func(context *http.Context) error {
		fmt.Println("something reaches here")
		context.WriteResponse("Hey")
		return nil
	})

	http.Get("/view", func(context *http.Context) error {
		context.Message.Error("Test error")
		return context.View("myapp", "test", "name", "John Doe", map[string]interface{}{
			"a": "A",
			"b": "B",
		})
	})

	http.Get("/panic", func(context *http.Context) error {
		var m map[string]interface{}
		return m["1"].(error)
	})

	http.WebSocket("/ws", func(context *http.Context, c *websocket.Conn) error {

		for {
			msg := <-message

			if err := c.WriteMessage(1, []byte(msg)); err != nil {
				log.Println("write:", err)
				break
			}
		}
		return nil
	})
	return nil
}

func Ready() error {
	//TestDistLock()
	//TestCache()
	var db =  evo.GetDBO()
	db.AutoMigrate(&MyModel{})
	var x = MyModel{
		Name: generic.Parse(evo.Config),LastName: "Last Name",
	}
	db.Debug().Create(&x)


	fmt.Println(x)
	var find MyModel
	fmt.Println(db.Where("id = ?",x.ID).Take(&find).Error)

	var cfg evo.Configuration
	fmt.Println(find.Name.ParseJSON(&cfg))
	fmt.Println(cfg)
	return nil
}

func TestCache()  {
	redis.Set("randomkey","test",1*time.Minute)
	for i := 0 ; i < 10; i++{
		redis.Cache.Set(fmt.Sprintf("key-%d", i),i,1*time.Minute)
	}
	fmt.Println("before flush:",redis.Cache.Keys())
	redis.Cache.Flush()
	fmt.Println("after flush:",redis.Cache.Keys())
}
func TestDistLock() {
	var ctx = context.Background()

	// Try to obtain lock.
	lock, err := redis.Locker.Obtain("my-key", 10*time.Second, nil,nil)
	if err == redis.ErrLockNotObtained {
		fmt.Println("Could not obtain lock!")
	} else if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			_, err := redis.Locker.Obtain("my-key", 10*time.Second, nil, nil)
			if err == redis.ErrLockNotObtained {
				fmt.Println("Could not obtain lock from thread!")
			}
			if err == nil{
				fmt.Println("Thread also got lock")
			}
			time.Sleep(1*time.Second)
		}
	}()
	// Don't forget to defer Release.
	defer lock.Release(ctx)
	fmt.Println("I have a lock!")
}
