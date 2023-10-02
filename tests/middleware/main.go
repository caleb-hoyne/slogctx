package main

import (
	"fmt"
	"github.com/caleb-hoyne/slogctx"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	h := &slogctx.Handler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	}
	slog.SetDefault(slog.New(h))

	go func() {
		time.Sleep(1 * time.Second)
		count := 1
		for range time.Tick(2 * time.Second) {
			req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			req.Header.Set("X-Request-Id", fmt.Sprintf("%d", count))
			_, _ = http.DefaultClient.Do(req)
			count++
		}
	}()
	s := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			ctx = slogctx.AddValues(ctx, slog.String("request_id", request.Header.Get("X-Request-Id")))
			slog.InfoContext(ctx, "Hello, world!")
		}),
	}

	panic(s.ListenAndServe())
}
