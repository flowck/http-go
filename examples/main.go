package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/flowck/http-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	router := http_go.NewServerDefaultNaiveRouter()
	router.GET("/", func(r *http_go.Request, w *http_go.Response) error {
		w.Headers.Set("Content-Type", "text/html; charset=UTF-8")
		return errors.New("something") // w.Write([]byte("Hello world"))
	})

	router.POST("/", func(r *http_go.Request, w *http_go.Response) error {
		return w.Write([]byte("request received"))
	})

	router.GET("/peoples", func(r *http_go.Request, w *http_go.Response) error {
		peoples := make([]map[string]string, 100)

		for i := 0; i < 100; i++ {
			peoples[i] = map[string]string{
				"id":         gofakeit.UUID(),
				"first_name": gofakeit.FirstName(),
				"last_name":  gofakeit.LastName(),
				"email":      gofakeit.Email(),
			}
		}

		payload, err := json.Marshal(peoples)
		if err != nil {
			w.WriteStatus(500)
			return w.Write([]byte("Internal Server Error"))
		}

		w.Headers.Set("Content-Type", "application/json; charset=UTF-8")
		return w.Write(payload)
	})

	s := http_go.Server{
		Addr:   ":8080",
		Router: router,
		Ctx:    ctx,
	}

	go func() {
		log.Println("server is running")
		log.Println(s.ListenAndServe())
		log.Println("server got shutdown")
	}()

	<-done

	log.Println(s.Shutdown())
	time.Sleep(time.Second * 2)
}
