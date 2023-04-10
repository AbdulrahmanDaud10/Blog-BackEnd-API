package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AbdulrahmanDaud10/fullstack-project/api/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) IntializeDB(DBDriver, DBUser, DBPassword, DBPort, DBHost, DBName string) {
	var err error

	if DBDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)
		server.DB, err = gorm.Open(DBDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to the %s database", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", DBDriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //DB Moigration

	server.Router = mux.NewRouter()

	server.DB.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
