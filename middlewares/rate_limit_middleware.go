package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RateLimiter struct {
	limiter *redis_rate.Limiter
}

var (
	contex = context.Background()
	// rdb     *redis.Client
	// once    sync.Once
)

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter: pkg.GetLimiter(),
	}
}

// func InitRateLimiter() {
// 	once.Do(func() {
// 		rdb = redis.NewClient(&redis.Options{
// 			Addr: pkg.Cfg.Redis.Host + ":" + fmt.Sprintf("%d", pkg.Cfg.Redis.Port), // "localhost:6379",
// 		})
// 		limiter = redis_rate.NewLimiter(rdb)
// 	})
// }

// func RateLimitMiddleware(limiter redis_rate.Limiter) func(c *fiber.Ctx) error {
func (r *RateLimiter) RateLimitMiddleware(maxRequests *int, window *time.Duration) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Implement rate limiting logic here
		// For example, you can use a token bucket or leaky bucket algorithm
		defaultMaxRequests := pkg.Cfg.Application.DefaultMaxRequestPerMinute
		defaultWindow := time.Minute

		if r.limiter == nil {
			return errors.InternalError("Rate limiter not initialized")
		}

		if maxRequests == nil {
			maxRequests = &defaultMaxRequests
		}
		if window == nil {
			window = &defaultWindow
		}

		claims := c.Locals("user").(jwt.MapClaims)
		userID := fmt.Sprintf("%v", claims["sub"])
		key := "rate_limit:" + userID + ":" + c.Path()

		// 50 request per menit per user
		// res, err := limiter.Allow(contex, key, redis_rate.PerMinute(50))

		limit := redis_rate.Limit{
			Rate:   *maxRequests,
			Period: *window,
			Burst:  *maxRequests,
		}

		res, err := r.limiter.Allow(contex, key, limit)
		if err != nil {
			return errors.InternalError(fmt.Sprintf("Rate limit error: %v", err))
		}

		if res.Allowed == 0 {
			return errors.TooManyRequests("Rate limit exceeded, please try again later")
		}

		return c.Next()
	}
}
