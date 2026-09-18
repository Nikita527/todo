package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Nikita527/todo/internal/todo"
)

func exitErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mustSave(store *todo.Store, path string) {
	exitErr(store.Save(path))
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	fs := flag.NewFlagSet("todo", flag.ContinueOnError)
	file := fs.String("file", "todos.json", "file to store tasks")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, "usage: todo <command> [arguments]\n")
		os.Exit(1)
	}
	store := todo.New()
	exitErr(store.Load(*file))
	switch args[0] {
	case "add":
		if len(args) < 2 {
			usage("usage: todo add <task>")
		}
		task := args[1]
		created, err := store.Add(task)
		exitErr(err)
		mustSave(store, *file)
		fmt.Printf("%d: %s, %t\n", created.ID, created.Text, created.Done)
	case "list":
		tasks := store.List()
		for _, task := range tasks {
			fmt.Printf("%d: %s %t\n", task.ID, task.Text, task.Done)
		}
	case "done":
		if len(args) < 2 {
			usage("usage: todo done <id>")
		}
		id, err := strconv.Atoi(args[1])
		exitErr(err)
		exitErr(store.Done(id))
		mustSave(store, *file)
	case "delete":
		if len(args) < 2 {
			usage("usage: todo delete <id>")
		}
		id, err := strconv.Atoi(args[1])
		exitErr(err)
		exitErr(store.Delete(id))
		mustSave(store, *file)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\nusage: todo <add|list|done|delete> ...\n", args[0])
		os.Exit(1)
	}
}
