package httpv3

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	c, _ := New(
		WithCustomTimeout(time.Second*25),
		WithDisallowedRedirects(),
	)

	res, err := c.Get("https://tls.peet.ws/api/all")

	if err != nil {
		panic(err)
	}

	fmt.Println(res.BodyString())
}
