package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"redis",
	ConfigModule,
	fx.Provide(NewClient),
	fx.Invoke(HookRedis),
)

type Redis struct {
	logger *zap.SugaredLogger
	client *redis.Client
	ctx    context.Context
}

func NewClient(c *Config, logger *zap.SugaredLogger) *Redis {
	logger.Debug("Starting Redis connection...")
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     *c.ConnectionString(),
		Password: *c.Password,
		DB:       *c.DB,
	})

	return &Redis{
		logger: logger,
		client: client,
		ctx:    ctx,
	}
}

func (r *Redis) SetNX(ctx context.Context, key *string, value interface{}, expiration *time.Duration) (*bool, error) {
	locked, err := r.client.SetNX(ctx, *key, value, *expiration).Result()
	if err != nil {
		return utils.PBool(false), err
	}

	return utils.PBool(locked), nil
}

func (r *Redis) Set(ctx context.Context, key *string, value interface{}, expiration *time.Duration) error {
	res, err := r.client.Set(ctx, *key, value, *expiration).Result()
	if err != nil {
		return fmt.Errorf("error saving message to key: %s (%v)", *key, err)
	}

	r.logger.Debugf("saved key %s successfully: %s", *key, res)
	return nil
}

func (r *Redis) GetMsgByKey(key *string) (*string, error) {
	val, err := r.client.Get(r.ctx, *key).Result()
	if err != nil {
		r.logger.Errorw("error retrieving message with key", "key", key, "error", err)
		return nil, fmt.Errorf(fmt.Sprint("error retrieving message with key: ", key, " (", err, ")"))
	}

	return utils.PString(val), nil
}

func (r *Redis) Keys(prefix *string) (*[]string, error) {
	result, err := r.client.Keys(r.ctx, *prefix+"*").Result()
	if err != nil {
		r.logger.Errorw("error scanning messages", "prefix", *prefix, "error", err)
		return nil, fmt.Errorf("error scanning messages. Prefix: %v. Err: %v", *prefix, err.Error())
	}
	return &result, nil
}

func (r *Redis) Del(matchingKeys *string) error {
	_, err := r.client.Del(r.ctx, *matchingKeys).Result()
	if err != nil {
		r.logger.Errorw("error deleting messages", "matchingKeys", *matchingKeys, "error", err)
		return fmt.Errorf("error deleting messages. Matching keys: %v. Err: %v", *matchingKeys, err.Error())
	}
	return nil
}

func (r *Redis) DeleteMsgByKey(key *string) error {
	err := r.client.Del(r.ctx, *key).Err()
	if err != nil {
		r.logger.Errorw("error deleting key", "key", *key, "error", err)
		return fmt.Errorf("error deleting key: %v", err)
	}

	return nil
}

func HookRedis(lc fx.Lifecycle, r *Redis, logger *zap.SugaredLogger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			pong, err := r.client.Ping(r.ctx).Result()
			if err != nil {
				logger.Panic(err)
			}
			logger.Infof("redis connection: %s", pong)
			return nil
		},
		OnStop: func(context.Context) error {
			err := r.client.Close()
			if err != nil {
				logger.Error(err)
			}
			logger.Info("redis connection Closed!")
			return nil
		},
	})
}
