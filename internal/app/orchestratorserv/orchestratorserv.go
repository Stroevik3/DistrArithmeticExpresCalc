package orchestratorserv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store/memstore"
	"github.com/sirupsen/logrus"
)

func Start(cfg *Config) error {
	store := memstore.New()
	s := newServer(store)

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to get log lvl, %v", err)
	}
	s.logger.SetLevel(level)
	ctx, cancel := context.WithCancel(context.Background())

	go runHTTPServer(ctx, cfg, s)

	shutDown(cancel)

	return nil
}

func runHTTPServer(ctx context.Context, cfg *Config, s *server) {

	s.mux.Handle("POST /api/v1/calculate", s.SetEpressionHandler(ctx))
	s.mux.Handle("GET /api/v1/expressions", s.GetEpressionsHandler(ctx))
	s.mux.Handle("GET /api/v1/expressions/{id}", s.GetEpressionHandler(ctx))
	s.mux.Handle("GET /api/v1/tasks", s.GetTasksHandler(ctx))
	s.mux.Handle("GET /internal/task", s.GetTaskToCompleteHandler(ctx))
	s.mux.Handle("POST /internal/task", s.PostTaskResultHandler(ctx))

	log.Printf("starting listening http server at %s", cfg.HttpAddr)
	if err := http.ListenAndServe(cfg.HttpAddr, s); err != nil {
		log.Fatalf("error service http server %v", err)
	}
}

func shutDown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	log.Println(errorMessage)
	cancel()
}
