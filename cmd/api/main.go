package main

import (
	"fmt"
	"go-app-template/internal/server"
)

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
