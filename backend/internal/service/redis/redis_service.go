package redis

import (
	"context"
	"errors"
	"fmt"
	"mychat-backend/internal/config"
	"mychat-backend/pkg/zaplog"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// 空上下文
// 传递ctx可以向redis操作设置超时，取消，传递元数据等
var ctx = context.Background()

func init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.RedisConfig.Db
	addr := host + ":" + strconv.Itoa(port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// 设置带过期时间的key-value
func SetKeyEx(key string, value string, timeout time.Duration) error {
	if err := redisClient.Set(ctx, key, value, timeout).Err(); err != nil {
		return err
	}
	return nil
}

// 查询键值，key不存在视为正常
func GetKey(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zaplog.Info("该key不存在")
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// 查询键值，key不存在视为错误，直接返回错误值
func GetKeyNilIsErr(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 前缀查询key，前缀匹配的key必须是唯一的
func GetKeyWithPrefixNilIsErr(prefix string) (string, error) {
	var keys []string
	var err error

	//该模块可用for循环 + scan 分批获取键来优化性能
	keys, err = redisClient.Keys(ctx, prefix+"*").Result()
	if err != nil {
		return "", err
	}
	if len(keys) == 0 {
		zaplog.Info("没有找到相关前缀key")
		return "", redis.Nil
	}
	if len(keys) == 1 {
		zaplog.Info(fmt.Sprintln("成功找到了相关前缀key", keys))
		return keys[0], nil
	} else {
		zaplog.Info("前缀匹配到的key数量大于1，查找异常")
		return "", errors.New("前缀匹配到的key数量大于1，查找异常")
	}
}

// 与前缀了类似，只不过是查询后缀
func GetKeyWithSuffixNilIsErr(suffix string) (string, error) {
	var keys []string
	var err error

	keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
	if err != nil {
		return "", err
	}
	if len(keys) == 0 {
		zaplog.Info("没有找到相关后缀key")
		return "", redis.Nil
	}
	if len(keys) == 1 {
		zaplog.Info(fmt.Sprintln("成功找到了相关后缀key", keys))
		return keys[0], nil
	} else {
		zaplog.Info("后缀匹配到的key数量大于1，查找异常")
		return "", errors.New("后缀匹配到的key数量大于1，查找异常")
	}
}

func DelKeyIfExists(key string) error {
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	//键存在
	if exists == 1 {
		delErr := redisClient.Del(ctx, key).Err()
		if delErr != nil {
			//删除失败需报错
			return delErr
		}
	}
	//删除的键不存在，也不需要报错
	return nil
}

// 删除包含pattern的所有键，模式匹配
func DelKeysWithPattern(pattern string) error {
	var keys []string
	var err error

	keys, err = redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		zaplog.Info("未找到对应的key")
	} else {
		delErr := redisClient.Del(ctx, keys...).Err()
		if delErr != nil {
			return delErr
		}
		zaplog.Info("删除成功")
	}
	return nil
}

func DelKeysWithPrefix(prefix string) error {
	//var cursor uint64 = 0
	var keys []string
	var err error

	// 使用 Keys 命令迭代匹配的键
	keys, err = redisClient.Keys(ctx, prefix+"*").Result()
	if err != nil {
		return err
	}

	// 如果没有更多的键，则跳出循环
	if len(keys) == 0 {
		zaplog.Info("未找到对应的键")
	}

	// 删除找到的键
	if len(keys) > 0 {
		_, err = redisClient.Del(ctx, keys...).Result()
		if err != nil {
			return err
		}
		zaplog.Info("成功删除相关前缀key")
	}

	return nil
}

func DelKeysWithSuffix(suffix string) error {
	//var cursor uint64 = 0
	var keys []string
	var err error

	// 使用 Keys 命令迭代匹配的键
	keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
	if err != nil {
		return err
	}

	// 如果没有更多的键，则跳出循环
	if len(keys) == 0 {
		zaplog.Info("没有找到相关后缀key")
	}

	// 删除找到的键
	if len(keys) > 0 {
		_, err = redisClient.Del(ctx, keys...).Result()
		if err != nil {
			return err
		}
		zaplog.Info("成功删除相关后缀key")
	}

	return nil
}

func DeleteAllRedisKeys() error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		cursor = nextCursor

		if len(keys) > 0 {
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
