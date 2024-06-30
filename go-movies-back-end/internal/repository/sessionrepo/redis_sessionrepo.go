package sessionrepo

import "github.com/go-redis/redis"

type RedisSessionRepo struct {
	// Redis Connection
	RDB *redis.Client
}

func (r *RedisSessionRepo) Connection() *redis.Client {
	return r.RDB
}

// Implement Redis Get
func (r *RedisSessionRepo) Get(key string) (string, error) {
	return r.RDB.Get(key).Result()
}

// Implement Redis Set
func (r *RedisSessionRepo) Set(key string, value string) error {
	return r.RDB.Set(key, value, 0).Err()
}

// Implement Redis Delete
func (r *RedisSessionRepo) Delete(key string) error {
	return r.RDB.Del(key).Err()
}
