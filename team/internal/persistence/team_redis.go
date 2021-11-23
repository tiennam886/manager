package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/model"
)

type redisEmployeeRepository struct {
	redisCache *redis.Client
}

func newRedisTeamRepository() (repo TeamRepository, err error) {
	redisCache := redis.NewClient(&redis.Options{
		Addr:     config.Get().RedisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	repo = &redisEmployeeRepository{
		redisCache: redisCache,
	}
	return
}

func (repo redisEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	var employee model.Team

	val, err := repo.redisCache.Get(uid).Result()
	if err != nil {
		return employee, err
	}
	err = json.Unmarshal([]byte(val), &employee)
	return employee, err
}

func (repo redisEmployeeRepository) Save(ctx context.Context, employee model.Team) error {
	j, err := json.Marshal(employee)
	if err != nil {
		fmt.Println(err.Error())
	}

	return repo.redisCache.Set(employee.UID, j, 10*time.Minute).Err()
}

func (repo redisEmployeeRepository) Update(ctx context.Context, uid string, employee model.Team) error {
	j, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	return repo.redisCache.Set(uid, j, 10*time.Minute).Err()
}

func (repo redisEmployeeRepository) Remove(ctx context.Context, uid string) error {
	return repo.redisCache.Del(uid).Err()
}
