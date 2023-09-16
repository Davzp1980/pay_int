package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"payint/handler"
	"payint/repository"
	"payint/service"

	_ "github.com/lib/pq"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("BD is not opend")
	}
	defer db.Close()

	repository.MigrateDB(db)

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := handler.NewHandler(service)

	log.Fatal(http.ListenAndServe(":8000", handler.InitRouters()))

}
