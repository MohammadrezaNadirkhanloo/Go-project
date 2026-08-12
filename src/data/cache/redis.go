package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) {
	redisClient = redis.NewClient(&redis.Options{
		// آدرس سرور Redis (مثلاً localhost:6379)
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),

		// پسورد Redis (اگر تنظیم نشده باشد خالی می‌ماند)
		Password: cfg.Redis.Password,

		// شماره دیتابیس (0 تا 15، پیش‌فرض 0 است)
		DB: 0,

		// حداکثر زمان مجاز برای برقراری اتصال اولیه به سرور Redis
		DialTimeout: cfg.Redis.DialTimeout * time.Second,

		// حداکثر زمان مجاز برای خواندن داده از Redis
		ReadTimeout: cfg.Redis.ReadTimeout * time.Second,

		// حداکثر زمان مجاز برای نوشتن داده در Redis
		WriteTimeout: cfg.Redis.WriteTimeout * time.Second,

		// حداکثر تعداد اتصال‌های همزمان که می‌تواند به Redis برقرار شود
		PoolSize: cfg.Redis.PoolSize,

		// حداکثر زمان انتظار برای دریافت اتصال از pool
		// اگر pool پر باشد و اتصال خالی نباشد، چقدر صبر کند
		PoolTimeout: cfg.Redis.PoolTimeout * time.Second,

		// حداقل تعداد اتصال‌های idle که همیشه باز نگه داشته می‌شوند
		// باعث می‌شود در زمان پیک درخواست‌ها تاخیر کمتری داشته باشیم
		MinIdleConns: 10,

		// حداکثر تعداد اتصال‌های idle که می‌تواند وجود داشته باشد
		// اگر بیشتر از این تعداد idle باشد، اتصال‌های اضافی بسته می‌شوند
		MaxIdleConns: 50,

		// حداکثر زمانی که یک اتصال می‌تواند idle (بیکار) بماند
		// بعد از این زمان، اتصال بسته می‌شود تا منابع آزاد شوند
		ConnMaxIdleTime: 5 * time.Minute,

		// حداکثر عمر کل یک اتصال از زمان ایجاد
		// حتی اگر فعال باشد، بعد از این زمان بسته و دوباره ایجاد می‌شود
		// برای prevent از memory leak و مشکلات طولانی‌مدت
		ConnMaxLifetime: 1 * time.Hour,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger := logging.NewLogger(cfg)
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	} else {
		logger.Info(logging.Redis, logging.Startup, "✅ Successfully connected to Redis", nil)
	}
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	if redisClient == nil {
		log.Println("⚠️ Redis client is nil, nothing to close")
		return
	}

	err := redisClient.Close()
	if err != nil {
		log.Printf("❌ Error closing Redis connection: %v", err)
	} else {
		log.Println("✅ Redis connection closed successfully")
	}
}
