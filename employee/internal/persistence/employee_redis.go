package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/model"

	"github.com/go-redis/redis"
)

type redisEmployeeRepository struct {
	redisCache *redis.Client
}

func newRedisEmployeeRepository() (repo EmployeeRepository, err error) {
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

func (repo redisEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.EmployeePost, error) {
	var employee model.Employee

	val, err := repo.redisCache.Get(uid).Result()
	if err != nil {
		return employee.ToEmployeePost(), err
	}
	err = json.Unmarshal([]byte(val), &employee)
	return employee.ToEmployeePost(), err
}

func (repo redisEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	j, err := json.Marshal(employee)
	if err != nil {
		fmt.Println(err.Error())
	}

	return repo.redisCache.Set(employee.UID, j, 10*time.Minute).Err()
}

func (repo redisEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	j, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	return repo.redisCache.Set(uid, j, 10*time.Minute).Err()
}

func (repo redisEmployeeRepository) Remove(ctx context.Context, uid string) error {
	return repo.redisCache.Del(uid).Err()
}

func (repo redisEmployeeRepository) AddToTeam(ctx context.Context, employeeId string, teamId string) error {
	return fmt.Errorf("none method from redis")
}

func (repo redisEmployeeRepository) FindByEmployeeId(ctx context.Context, employeeId string) ([]string, error) {
	return nil, fmt.Errorf("none method from redis")
}

func (repo redisEmployeeRepository) DeleteByEmployeeId(ctx context.Context, employeeId string) error {
	return fmt.Errorf("none method from redis")
}

func (repo redisEmployeeRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error) {
	return []model.EmployeePost{}, fmt.Errorf("no get all from redis")
}

func (repo redisEmployeeRepository) DeleteFromTeam(ctx context.Context, employeeId string, teamId string) error {
	return fmt.Errorf("none method from redis")
}
