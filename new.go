package httpv3

func New(options ...any) (*Client, error) {
	cfg := parseOptions(options...)

	c := &Client{
		cfg: cfg,
	}

	err := c.reinitFhttpClient()

	return c, err
}
