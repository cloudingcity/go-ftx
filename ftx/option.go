package ftx

import "net/url"

type Option func(*Client)

func WithAuth(key, secret string) Option {
	return func(c *Client) {
		c.key = key
		c.secret = []byte(secret)
	}
}

func WithSubAccount(account string) Option {
	return func(c *Client) {
		c.subAccount = url.QueryEscape(account)
	}
}
