package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()
var redisClient *redis.Client

type User struct {
	ID   string
	Name string
	Age  int
}

func (u User) UserIdCacheKey() string {
	return fmt.Sprintf("UserID_%d", u.ID)
}

var MyDataBase = map[int]User{
	1: {ID: "656", Name: "John", Age: 30},
	2: {ID: "657", Name: "Mary", Age: 25},
	3: {ID: "658", Name: "Peter", Age: 40},
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Endereço do seu servidor Redis
		Password: "",               // Senha (se necessário)
		DB:       0,                // Número do banco de dados Redis a ser usado
	})
}

func getDataFromCacheOrDB(key string) (*User, error) {
	userFound := &User{}

	cachedData, err := redisClient.Get(ctx, key).Result()
	if err == nil && cachedData != "" {
		if err = json.Unmarshal([]byte(cachedData), &userFound); err != nil {
			return nil, err
		}
		return userFound, nil
	}

	//BANCO DE DADOS POSTGRES
	for _, user := range MyDataBase {
		if user.ID == key {
			userFound = &user
		}
	}

	jsonUser, err := json.Marshal(userFound)
	if err != nil {
		return nil, err
	}

	err = redisClient.SetEX(ctx, key, jsonUser, time.Minute).Err()
	if err != nil {
		fmt.Println("Erro ao definir o valor no cache:", err)
	}

	return userFound, nil
}

func main() {
	initRedis()

	data, err := getDataFromCacheOrDB("657")
	if err != nil {
		fmt.Println("Erro ao obter dados:", err)
		return
	}

	fmt.Println("Dados finais:", data)
}
