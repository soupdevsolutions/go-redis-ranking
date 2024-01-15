package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"os"

	"github.com/go-redis/redis/v9"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var redisAddr = getEnv("REDIS_ADDR", "localhost:6379")
var redisPass = getEnv("REDIS_PASS", "")
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     redisAddr,
	Password: redisPass,
	DB:       0,
})

func addUser(name string) (user, error) {
	id := uuid.NewString()
	newUser := user{Name: name, ID: id, Score: len(name)}
	serializedUser, err := json.Marshal(newUser)
	if err != nil {
		return user{}, err
	}

	err = rdb.Set(ctx, id, serializedUser, 0).Err()
	if err != nil {
		return user{}, err
	}

	err = rdb.ZAdd(ctx, "rank", redis.Z{Score: float64(newUser.Score), Member: id}).Err()
	if err != nil {
		return user{}, err
	}

	rank, err := getUserRank(id)
	if err != nil {
		return user{}, err
	}
	newUser.Rank = int(rank)

	return newUser, nil
}

func getUserData(id string) (user, error) {
	userData, err := rdb.Get(ctx, id).Result()
	if err != nil {
		return user{}, err
	}

	var existingUser user
	err = json.Unmarshal([]byte(userData), &existingUser)
	if err != nil {
		return user{}, err
	}

	rank, err := getUserRank(id)
	if err != nil {
		return user{}, err
	}
	existingUser.Rank = int(rank)

	return existingUser, nil
}

func getUserRank(id string) (int, error) {
	rank, err := rdb.ZRevRank(ctx, "rank", id).Result()
	if err != nil {
		return -1, err
	}

	return int(rank + 1), nil
}

func getRedisRanks(offset, limit int) ([]user, error) {
	cmd := rdb.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:     "rank",
		Start:   offset - 1,
		Stop:    offset + limit - 2,
		Rev:     true,
		ByLex:   false,
		ByScore: false,
	})

	log.Printf("Command: %v", cmd.String())

	users, err := cmd.Result()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	log.Printf("Users: %v", users)
	var result []user
	for rank, userScore := range users {
		id := userScore.Member.(string)
		userData, err := rdb.Get(ctx, id).Result()
		if err != nil {
			return nil, err
		}
		var constructedUser user
		err = json.Unmarshal([]byte(userData), &constructedUser)
		if err != nil {
			log.Printf("Error: %v", err)
			return nil, err
		}

		result = append(result, user{ID: id, Name: constructedUser.Name, Score: int(userScore.Score), Rank: rank + offset})
	}

	return result, nil
}
