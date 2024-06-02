package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/agent/model"
	"github.com/sirupsen/logrus"
)

type client struct {
	logger *logrus.Logger
	cl     *http.Client
}

func newClient() *client {
	c := &client{
		logger: logrus.New(),
		cl:     &http.Client{},
	}

	return c
}

func (c *client) getTask(ctx context.Context) (*model.TaskGetResp, error) {
	c.logger.Debugln("getTask")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL_TASK, nil)
	if err != nil {
		c.logger.Error("NewRequestWithContext err - ", err.Error())
		return nil, err
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		c.logger.Error("Do err - ", err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("response status =", resp.StatusCode)
		return nil, errors.New("response status not 200")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("ReadAll err - ", err)
		return nil, err
	}

	type TaskRow struct {
		Task model.TaskGetResp `json:"task"`
	}
	var t TaskRow
	err = json.Unmarshal(body, &t)
	if err != nil {
		c.logger.Error("Unmarshal err - ", err)
		return nil, err
	}
	return &t.Task, nil
}

func (c *client) postTask(ctx context.Context, t *model.TaskReq) error {
	c.logger.Debugln("postTask")
	taskJsonByt, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL_TASK, bytes.NewReader(taskJsonByt))
	if err != nil {
		c.logger.Error("NewRequestWithContext err - ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cl.Do(req)
	if err != nil {
		c.logger.Error("Do err - ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("response status =", resp.StatusCode)
		return errors.New("response status not 200")
	}

	return nil
}
