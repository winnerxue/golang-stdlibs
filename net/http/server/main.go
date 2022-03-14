package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	go func() {
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie("name"); err != nil {
				fmt.Println("new")
				http.SetCookie(w, &http.Cookie{Name: "name", Value: "aaa"})
			} else {
				fmt.Println("old", cookie.Name, cookie.Value)
				cookie.Value = "123"
				http.SetCookie(w, cookie)
			}

		}))
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)

	defer cancel()

	for {
		select {
		case <-ctx.Done():
			break
		default:
		}
	}

}
