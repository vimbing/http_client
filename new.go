package http

import (
	"time"

	"github.com/vimbing/retry"
)

func New(options ...any) (*Client, error) {
	cfg := parseOptions(options...)

	c := &Client{
		cfg: cfg,
	}

	err := c.reinitFhttpClient()

	return c, err
}

func MustNew(options ...any) *Client {
	client := &Client{}

	err := retry.Retrier{Max: 15, Delay: time.Millisecond * 15}.Retry(func() error {
		var err error
		client, err = New(options...)
		return err
	})

	if err != nil {
		panic(err)
	}

	return client
}
