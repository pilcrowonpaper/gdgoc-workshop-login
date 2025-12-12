package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	_ "embed"

	"golang.org/x/sync/semaphore"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed schema.sql
var schema string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Commands:")
		fmt.Println("- init: Set up database.")
		fmt.Println("- start [port=3000]: Start server (default port of 3000).")
		return
	}

	if os.Args[1] == "init" {
		err := removeFileIfExists("main.db")
		if err != nil {
			log.Fatalf("Failed to remove db file: %s\n", err.Error())
		}

		conn, err := sqlite.OpenConn("main.db", sqlite.OpenReadWrite, sqlite.OpenCreate)
		if err != nil {
			log.Fatalf("Failed to open db connection: %s\n", err.Error())
		}

		err = sqlitex.ExecuteScript(conn, schema, nil)
		if err != nil {
			log.Fatalf("Failed to execute schema set up script: %s\n", err.Error())
		}

		conn.Close()

		return
	}

	if os.Args[1] == "start" {
		port := 3000
		if len(os.Args) >= 3 {
			parsedPort, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalln("Failed to parse port argument")
			}
			port = parsedPort
		}

		conn, err := sqlite.OpenConn("main.db", sqlite.OpenReadWrite)
		if err != nil {
			log.Fatalf("Failed to open db connection: %s\n", err.Error())
		}

		server := &serverStruct{
			conn:                          conn,
			passwordHashingSemaphore:      semaphore.NewWeighted(int64(runtime.NumCPU())),
			passwordVerificationRateLimit: newRateLimit(5, time.Minute),
		}

		fmt.Printf("Starting server at port %d...\n", port)

		address := fmt.Sprintf(":%d", port)
		err = http.ListenAndServe(address, server)
		if err != nil {
			log.Fatalf("Failed to start server on port %d: %s\n", port, err.Error())
		}
	}

	log.Fatalf("Unknown command %s\n", os.Args[1])
}

func removeFileIfExists(name string) error {
	err := os.Remove(name)
	if err == nil {
		return nil
	}
	if _, ok := err.(*os.PathError); ok {
		return nil
	}
	return err
}
