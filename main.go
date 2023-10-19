package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/uptrace/bunrouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))
	httpLn, err := net.Listen("tcp", appPort)
	if err != nil {
		panic(fmt.Errorf("failed to listen to port %v: %v", appPort, err))
	}

	db, err := gorm.Open(
		sqlite.Open(
			os.Getenv("DB_PATH"),
		),
		&gorm.Config{},
	)
	if err != nil {
		panic(fmt.Errorf("failed to init db: %v", err))
	}

	r := bunrouter.New()

	validate := validator.New(validator.WithRequiredStructEnabled())

	r.WithGroup("/user", func(g *bunrouter.Group) {
		g.POST("", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
			body, err := io.ReadAll(bunReq.Body)
			if err != nil {
				return err
			}

			var a Admin

			err = json.Unmarshal(body, &a)
			if err != nil {
				return err
			}

			err = validate.Struct(&a)
			if err != nil {
				return err
			}

			err = db.Create(&a).Error
			if err != nil {
				return err
			}

			bunrouter.JSON(w, bunrouter.H{
				"a": a,
			})
			return nil
		})

		g.GET("", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
			data := []*Admin{}
			err := db.Find(&data).Error
			if err != nil {
				return err
			}

			bunrouter.JSON(w, bunrouter.H{
				"data": data,
			})
			return nil
		})

		g.GET("/:id", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
			id := bunReq.Param("id")

			m := Admin{}
			err := db.Where(map[string]interface{}{
				"id": id,
			}).Take(&m).Error
			if err != nil {
				return err
			}

			bunrouter.JSON(w, bunrouter.H{
				"m": m,
			})
			return nil
		})

		g.PUT("/:id", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
			a := Admin{}
			err := db.Create(&a).Error
			if err != nil {
				return err
			}

			bunrouter.JSON(w, bunrouter.H{
				"a": a,
			})
			return nil
		})
	})

	r.GET("/", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
		bunrouter.JSON(w, bunrouter.H{
			"message": "pong",
		})
		return nil
	})

	r.POST("/migrate", func(w http.ResponseWriter, bunReq bunrouter.Request) error {
		db.AutoMigrate(&Admin{})
		bunrouter.JSON(w, bunrouter.H{
			"message": "automigrate",
		})
		return nil
	})

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      r,
	}

	go func() {
		log.Printf("running on port %v", appPort)
		err := httpServer.Serve(httpLn)
		if err != nil {
			log.Println(fmt.Errorf("failed to serve http: %v", err))
		}
	}()

	log.Printf("got signal: %v", waitExitSignal().String())

	ctx := context.Background()
	// Graceful shutdown using sleep bcs I don't trust kube to properly shutdown
	time.Sleep(2 * time.Second)
	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(fmt.Errorf("trying to shutdown with error: %v", err))
	}
}

func waitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}

type Admin struct {
	ID          int       `json:"id"`
	Firstname   string    `validate:"required" json:"firstname"`
	Lastname    string    `validate:"required" json:"lastname"`
	Email       string    `validate:"required,email" json:"email"`
	DateOfBirth time.Time `validate:"required" json:"date_of_birth"`
	Gender      string    `validate:"oneof=male female prefer_not_to" json:"gender"`
	Password    string    `validate:"required" json:"password,omitempty"`
}
