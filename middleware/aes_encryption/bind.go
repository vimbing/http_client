package aesencryption

import (
	"bytes"
	"io"

	httpv3 "github.com/vimbing/httpv2"
)

func (m *Middleware) Bind(c *httpv3.Client) {
	c.UseRequest(func(r *httpv3.Request) error {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			return err
		}

		encrypted, err := m.encrypt(body)

		if err != nil {
			return err
		}

		r.Body = bytes.NewBuffer(encrypted)

		return nil
	})

	c.UseResponse(func(r *httpv3.Response) error {
		decrypted, err := m.decrypt(r.BodyBytes())

		if err != nil {
			return err
		}

		r.Body = decrypted

		return nil
	})
}
