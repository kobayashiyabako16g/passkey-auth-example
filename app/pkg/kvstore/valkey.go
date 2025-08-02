package kvstore

import (
	"context"

	"github.com/valkey-io/valkey-go"
)

type valKeyClient struct {
	client valkey.Client
	config ValKeyConfig
}

type ValKeyConfig struct {
	InitAddress []string `env:"INIT_ADDRESS"`
	Username    string   `env:"USERNAME"`
	Password    string   `env:"PASSWORD"`
	ClientName  string   `env:"CLIENTNAME"`
	Expiration  int64    `env:"EXPIRATION"`
}

func NewValKeyClient(config ValKeyConfig) (Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: config.InitAddress,
		Username:    config.Username,
		Password:    config.Password,
		ClientName:  config.ClientName,
	})
	if err != nil {
		return nil, err
	}

	return &valKeyClient{
		client: client,
		config: config,
	}, nil
}

func (c *valKeyClient) Get(ctx context.Context, key string) (string, error) {
	resp := c.client.Do(ctx, c.client.B().Get().Key(key).Build())
	if err := resp.Error(); err != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return "", nil //  not found
		}
		return "", err
	}

	val, err := resp.ToString()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *valKeyClient) Set(ctx context.Context, key, value string, opts ...SetOptions) error {
	var opt SetOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = SetOptions{
			Expiration: c.config.Expiration, // default
		}
	}

	resp := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Build())
	if err := resp.Error(); err != nil {
		return err
	}
	_, err := c.client.Do(ctx, c.client.B().Expire().Key(key).Seconds(opt.Expiration).Build()).AsInt64()
	if err != nil {
		return err
	}

	return nil
}
