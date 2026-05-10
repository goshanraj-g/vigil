package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/goshanraj-g/vigil/internal/api"
	"github.com/goshanraj-g/vigil/internal/db"
)

func main() {

	/*
	Contexts: 
	Contexts are a blank, never called context which lives for the lifetime of a program
	Contexts carries two things through a call stack:
	1) cancellation signal - "stop what you're doing"
	2) deadline - "stop after this time"
	*/
	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	r := api.NewRouter(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}