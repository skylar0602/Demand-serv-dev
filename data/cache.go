package data

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	log "github.com/cihub/seelog"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/smarterwallet/demand-abstraction-serv/config"
	"github.com/smarterwallet/demand-abstraction-serv/model"
)

type Cache struct {
	client *redis.Client
}

func keyConversation(cid string) string {
	return "smart-wallet-conversation:" + cid
}

func NewCache(redisCfg *config.RedisCfg) (*Cache, error) {
	if redisCfg == nil || redisCfg.Addr == "" {
		panic("NewRedisClient nil redisCfg or empty redis Addr")
	}
	opt := &redis.Options{}
	if redisCfg.Addr != "" {
		opt.Addr = redisCfg.Addr
	}
	if redisCfg.Password != "" {
		opt.Password = redisCfg.Password
	}
	if redisCfg.DB != 0 {
		opt.DB = redisCfg.DB
	}
	if redisCfg.MinIdle != 0 {
		opt.MinIdleConns = redisCfg.MinIdle
	}
	if redisCfg.PoolSize != 0 {
		opt.PoolSize = redisCfg.PoolSize
	}
	if redisCfg.DialTimeout != 0 {
		opt.DialTimeout = redisCfg.DialTimeout
	}
	if redisCfg.ReadTimeout != 0 {
		opt.ReadTimeout = redisCfg.ReadTimeout
	}
	if redisCfg.WriteTimeout != 0 {
		opt.WriteTimeout = redisCfg.WriteTimeout
	}
	if redisCfg.PoolTimeout != 0 {
		opt.PoolTimeout = redisCfg.PoolTimeout
	}
	client := redis.NewClient(opt)
	rCmd := client.Ping(context.Background())
	if err := rCmd.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &Cache{client: client}, nil
}

func (c *Cache) ChatHistory(ctx context.Context, cid string) ([]model.Dialogue, error) {
	res, err := c.client.SMembers(ctx, keyConversation(cid)).Result()
	if err != nil {
		return nil, err
	}
	dialogues := make([]model.Dialogue, 0)
	for _, s := range res {
		var dialogue model.Dialogue
		if err := json.Unmarshal([]byte(s), &dialogue); err != nil {
			return nil, err
		}
		dialogues = append(dialogues, dialogue)
	}
	sort.Slice(dialogues, func(i, j int) bool {
		return dialogues[i].Timestamp < dialogues[j].Timestamp
	})
	return dialogues, nil
}

func (c *Cache) AppendChat(ctx context.Context, cid, content, role string) error {
	var dialogue = &model.Dialogue{
		Type:      "text",
		Role:      role,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
	res, err := c.client.SAdd(ctx, keyConversation(cid), dialogue).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("duplicate conversation")
	}
	return nil
}

func keyCtx(key string) string {
	return "smart-wallet-context:" + key
}

func (c *Cache) SetCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if _, err := c.client.Set(ctx, keyCtx(key), value, expiration).Result(); err != nil {
		log.Errorf("SetCtx key=%s err=%s\n", key, err)
		return err
	}
	return nil
}

func (c *Cache) GetCtx(ctx context.Context, key string) (string, error) {
	res, err := c.client.Get(ctx, keyCtx(key)).Result()
	if err != nil {
		log.Errorf("GetCtx key=%s err=%s\n", key, err)
		return "", err
	}
	return res, nil
}

func (c *Cache) Invalid(ctx context.Context, cid string) error {
	return c.client.Del(ctx, keyConversation(cid)).Err()
}
