package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/agent/model"
	"github.com/sirupsen/logrus"
)

const (
	URL_TASK        string = "http://localhost:8080/internal/task"
	OPER_SYMB_ADDIT string = "+"
	OPER_SYMB_SUBTR string = "-"
	OPER_SYMB_MULTP string = "*"
	OPER_SYMB_DIVIS string = "/"
)

func Start(cfg *Config) error {

	ctx, cancel := context.WithCancel(context.Background())
	client := newClient()

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to get log lvl, %v", err)
	}
	client.logger.SetLevel(level)

	for i := 0; i < cfg.ComputingPower; i++ {
		go runHTTPCLServer(ctx, cfg, client)
	}

	shutDown(cancel)

	return nil
}

func runHTTPCLServer(ctx context.Context, cfg *Config, c *client) {
	c.logger.Debugln("runHTTPCLServer")
	for {
		time.Sleep(cfg.PollingInterval)
		select {
		case <-ctx.Done():
			return
		default:
			var (
				res float64
			)
			t, err := c.getTask(ctx)
			if err != nil {
				continue
			}
			time.Sleep(time.Duration(t.OperTime) * time.Millisecond)
			switch t.Oper {
			case OPER_SYMB_ADDIT:
				res = t.ArgOne + t.ArgTwo
			case OPER_SYMB_SUBTR:
				res = t.ArgOne - t.ArgTwo
			case OPER_SYMB_MULTP:
				res = t.ArgOne * t.ArgTwo
			case OPER_SYMB_DIVIS:
				res = t.ArgOne / t.ArgTwo
			}
			tReq := &model.TaskReq{
				Id:     t.Id,
				Result: res,
			}
			err = c.postTask(ctx, tReq)
			if err != nil {
				continue
			}
		}
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
