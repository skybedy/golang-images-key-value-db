package main

import (
	"fmt"
	"log"
	"mustafar/config"
	"mustafar/routes"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/spf13/viper"
)

func main() {
	var configuration config.Configurations
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	routes.LoadTemplates("templates/*.html")
	port := configuration.Server.Port
	routes.PgName = configuration.Database.PgName
	router := routes.NewRouter()

	serverSocket := socketio.NewServer(nil)
	serverSocket.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	serverSocket.OnEvent("/", "get-image", func(s socketio.Conn, imgName string) {
		json := string(routes.GetImgHandlerSocket(imgName))
		s.Emit("image", json)
	})

	serverSocket.OnEvent("/", "images", func(s socketio.Conn, data string, nazev string) {
		json := string(routes.ImgUploadHandlerSocket(nazev, data))
		s.Emit("images", json)
	})

	serverSocket.OnEvent("/", "truncatedbpg", func(s socketio.Conn) {
		routes.TruncateDbPg()
	})

	go serverSocket.Serve()
	defer serverSocket.Close()
	router.Handle("/socket.io/", serverSocket)

	server := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Spouští se server na portu", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
