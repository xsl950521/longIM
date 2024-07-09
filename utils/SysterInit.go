package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))
}

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func InitMysql() {
	// 自定义日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢sql阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	//if err != nil {
	//	panic("failed to connect database")
	//}

}

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		DB:           0,
		Password:     "",
		PoolSize:     viper.GetInt("redis.pool_size"),
		MinIdleConns: viper.GetInt("redis.min_idle_conn"),
	})
	pong, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis ping err:", err)
		panic(err)
	}
	fmt.Println("redis ping success ", pong)
}

const (
	PublishKey = "websocket"
)

// 发布消息
func Publish(ctx context.Context, ch string, msg string) error {
	var err error
	err = RDB.Publish(ctx, ch, msg).Err()
	return err
}

// 订阅消息
func Subscribe(ctx context.Context, ch string) (string, error) {
	sub := RDB.Subscribe(ctx, ch)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("subscribe:", msg.Payload)
	return msg.Payload, err
}
