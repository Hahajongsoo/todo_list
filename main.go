package main

import (
	"log"
	"net/http"
	"todoWeb/app"
)

func main() {
	m := app.MakeHandler("./test.db")
	defer m.Close()

	log.Println("Started App")
	err := http.ListenAndServe(":80", m)
	if err != nil {
		panic(err)
	}
}
