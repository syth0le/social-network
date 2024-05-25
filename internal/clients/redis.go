package clients

import "github.com/redis/go-redis/v9"

type RedisClient struct {
	*redis.Client
	// endpoints
}

func NewRedisClient(endpoints []string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Network:               "",
			Addr:                  "",
			ClientName:            "",
			Dialer:                nil,
			OnConnect:             nil,
			Protocol:              0,
			Username:              "",
			Password:              "",
			CredentialsProvider:   nil,
			DB:                    0,
			MaxRetries:            0,
			MinRetryBackoff:       0,
			MaxRetryBackoff:       0,
			DialTimeout:           0,
			ReadTimeout:           0,
			WriteTimeout:          0,
			ContextTimeoutEnabled: false,
			PoolFIFO:              false,
			PoolSize:              0,
			PoolTimeout:           0,
			MinIdleConns:          0,
			MaxIdleConns:          0,
			MaxActiveConns:        0,
			ConnMaxIdleTime:       0,
			ConnMaxLifetime:       0,
			TLSConfig:             nil,
			Limiter:               nil,
			DisableIndentity:      false,
			IdentitySuffix:        "",
		}),
		endpoints: endpoints,
	}
}

func (s *RedisClient) Close() error {
	return s.Client.Close()
}
