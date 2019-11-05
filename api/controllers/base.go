package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	//server.DB.Debug().AutoMigrate(&models.User{}) //, &models.Post{} Liviu No te olvide aqui es donde tiened que añadir los siguentes modelos

	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8090")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		//AllowedHeaders:     []string{"X-Requested-With,Origin, Accept, Content-Type, Content-Length, X-Requested-With, Accept-Encoding, X-CSRF-Token, Authorization, X-PINGOTHER"},
		AllowedMethods:     []string{"GET", "HEAD", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept", "Authorization"},
		Debug:              true,
		OptionsPassthrough: true,
	})

	handler := c.Handler(server.Router)
	log.Fatal(http.ListenAndServe(addr, handler))
}
