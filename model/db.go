package model

import (
	"OuterChat/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Cache *redis.Client
)

func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
	}
	_ = DB.AutoMigrate(&UserBasic{}, &GroupBasic{}, &Message{}, &Contact{})
}

func InitCache() {
	Cache = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConn,
	})
	_, err := Cache.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis Init Failed:", err)
	}
	fmt.Println("Redis Init Success")
}

const (
	PublishKey = "websocket"
)

func Publish(ctx context.Context, channel string, data interface{}) error {
	pub := Cache.Publish(ctx, channel, data)
	fmt.Println("Publish:", pub)
	return pub.Err()
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Cache.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	fmt.Println("Subscribe:", msg)
	return msg.Payload, nil
}
