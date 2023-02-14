package main

import (
	"log"
	"net/http"
	"todoWeb/app"
)

func main() {
	dbConn := "postgresql://postgres:1234@192.168.56.101/?sslmode=disable"
	m := app.MakeHandler(dbConn)
	defer m.Close()

	log.Println("Started App")
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
