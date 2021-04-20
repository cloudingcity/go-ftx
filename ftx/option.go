package ftx

type Option func(*Client)

func WithAuth(key, secret string) Option {
	return func(c *Client) {
		c.SetAuth(key, secret)
	}
}

func WithSubAccount(account string) Option {
	return func(c *Client) {
		c.SetSubAccount(account)
	}
}
