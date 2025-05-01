package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"buf.build/gen/go/wakabaseisei/ms-protobuf/connectrpc/go/ms/user/v1/userv1connect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wakabaseisei/ms-user/internal/config"
	"github.com/wakabaseisei/ms-user/internal/domain/repository"
	"github.com/wakabaseisei/ms-user/internal/driver/grpc"
	infraRepo "github.com/wakabaseisei/ms-user/internal/repository"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, cerr := config.NewConfig(ctx)
	if cerr != nil {
		log.Fatalf("New Config: %v", cerr)
	}

	dbConn, dberr := infraRepo.NewDatabase(ctx, cfg.DBConfig, cfg.AWSDefaultConfig)
	if dberr != nil {
		log.Fatalf("New database: %v", dberr)
	}
	defer closeDBConn(dbConn)

	sgCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	services := repository.NewServices(infraRepo.NewUserRepository(dbConn))
	service := grpc.NewUserService(services)
	mux := http.NewServeMux()

	path, handler := userv1connect.NewUserServiceHandler(service)

	mux.Handle(path, handler)
	mux.HandleFunc("/", healthCheckHandler)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	done := make(chan error, 1)
	go func() {
		done <- server.ListenAndServe()
	}()

	select {
	case err := <-done:
		if err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	case <-sgCtx.Done():
		log.Println("Server stopping")
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(c); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
		log.Println("Server gracefully stopped")
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func closeDBConn(db io.Closer) {
	if cerr := db.Close(); cerr != nil {
		log.Printf("closing db connection: %v", cerr)
	} else {
		log.Println("db connection gracefully closed")
	}
}
