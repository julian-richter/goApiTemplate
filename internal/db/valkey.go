package db

import (
	"context"
	"fmt"
	"time"

	"github.com/julian-richter/ApiTemplate/internal/config"
	"github.com/valkey-io/valkey-go"
)

type ValkeyClient struct {
	client valkey.Client
}

// ValkeyClientInterface defines the contract that repositories expect.
type ValkeyClientInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Close()
}

func (v *ValkeyClient) Get(ctx context.Context, key string) (string, error) {
	resp := v.client.Do(ctx, v.client.B().Get().Key(key).Build())
	if err := resp.Error(); err != nil {
		return "", err
	}
	val, err := resp.ToString()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (v *ValkeyClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	cmdSet := v.client.B().Set().Key(key).Value(value).Build()
	resp := v.client.Do(ctx, cmdSet)
	if err := resp.Error(); err != nil {
		return err
	}

	if ttl > 0 {
		cmdExpire := v.client.B().Expire().Key(key).Seconds(int64(ttl.Seconds())).Build()
		resp2 := v.client.Do(ctx, cmdExpire)
		if err2 := resp2.Error(); err2 != nil {
			return err2
		}
	}

	return nil
}

// Close implements ValkeyClientInterface.
func (v *ValkeyClient) Close() {
	v.client.Close()
}

func NewValkeyClient(cfg config.Config) (*ValkeyClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port)

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
		Username:    cfg.Cache.User,
		Password:    cfg.Cache.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect Valkey: %w", err)
	}

	// Use a timeout context instead of context.Background.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := client.Do(ctx, client.B().Ping().Build())
	if err := resp.Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("valkey ping failed: %w", err)
	}

	return &ValkeyClient{client: client}, nil
}
