package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Nikita527/todo/internal/api"
	"github.com/Nikita527/todo/internal/todo"
)

func exitErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	file := flag.String("file", "todos.json", "file to store tasks")
	flag.Parse()
	store := todo.New()
	exitErr(store.Load(*file))
	handler := api.New(store, *file).Routes()
	fmt.Printf("listening on %s\n", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
