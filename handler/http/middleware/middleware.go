package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	appCtx  *fiber.App
	redisCl *redis.Storage
	appLog  *logrus.Logger
	// another stuff , may be needed by middleware
}

func (m *GoMiddleware) RateLimiter() fiber.Handler {
	limiterCfg := limiter.Config{
		/*Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},*/
		Storage:    m.redisCl,
		Max:        10,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			authKey := c.Get("Authorization")
			if authKey == "" {
				return c.IP() + c.Get("User-Agent")
			} else {
				if len(authKey) > 7 && strings.HasPrefix(authKey, "Bearer ") {
					authKey = authKey[7:]
				}
				return authKey
			}
		},
		/*KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + c.Get("User-Agent")
		},*/
	}

	return limiter.New(limiterCfg)
}

func (m *GoMiddleware) RateLimiterFull() fiber.Handler {
	limiterCfg := limiter.Config{
		Max:        10,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			authKey := c.Get("Authorization")
			if authKey == "" {
				uniqueID := c.Cookies("sipd_penatausahaan_uk_")
				if uniqueID == "" {
					uniqueID = fmt.Sprintf("%d", time.Now().UnixNano())
					c.Cookie(&fiber.Cookie{
						Name:     "sipd_penatausahaan_uk_",
						Value:    uniqueID,
						Expires:  time.Now().Add(24 * time.Hour),
						HTTPOnly: true,
						Secure:   true,
					})
				}

				return c.IP() + c.Get("User-Agent") + uniqueID
			} else {
				return authKey
			}
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}

	return limiter.New(limiterCfg)
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS() fiber.Handler {
	crs := os.Getenv("SIPD_CORS_WHITELISTS")

	if crs == "*" {
		return cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowHeaders:  "Content-Type, Accept, Authorization",
			AllowMethods:  "GET, HEAD, PUT, PATCH, POST, DELETE",
			ExposeHeaders: "*", //"X-Pagination-Current-Page,X-Pagination-Next-Page,X-Pagination-Page-Count,X-Pagination-Page-Size,X-Pagination-Total-Count"
		})
	}

	return cors.New(cors.Config{
		AllowOrigins:     crs,
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Accept, Authorization",
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE",
		ExposeHeaders:    "*", //"X-Pagination-Current-Page,X-Pagination-Next-Page,X-Pagination-Page-Count,X-Pagination-Page-Size,X-Pagination-Total-Count"
	})
}

// LOGGER simple logger.
func (m *GoMiddleware) LOGGER() fiber.Handler {
	return logger.New()

	/*return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			log.Print("sdfsdsfdsddf")

			// Log incoming requests
			m.appLog.WithFields(logrus.Fields{
				"method":  c.Method(),
				"path":    c.Path(),
				"ip":      c.IP(),
				"message": err.Error(),
			}).Error("Error occurred")
			return err
		}

		// Proceed with the request
		return nil
	}*/
}

// JWT jwt.
func (m *GoMiddleware) JWT() fiber.Handler {
	// Create config for JWT authentication middleware.
	config := jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey: "jwt", // used in private routes
		// SuccessHandler: m.jwtSuccess,
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 400 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}

func (m *GoMiddleware) jwtSuccess(c *fiber.Ctx) error {
	var err error

	// log.Println("jwt success")

	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	tokenString := user.Raw
	pegId := int64(claims["id_pegawai"].(float64))
	redisKey := fmt.Sprintf("peg:%d", pegId)

	existingToken, err := m.redisCl.Get(redisKey)
	if err != nil {
		err = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to check token data. - " + err.Error(),
		})
		return err
	}
	if len(existingToken) == 0 {
		err = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "not authorized",
		})
		return err
	} else {
		if fmt.Sprintf("%s", existingToken) != tokenString {
			err = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "forbidden",
			})
			return err
		}
	}

	return c.Next()
}

// InitMiddleware initialize the middleware
func InitMiddleware(ctx *fiber.App, store *redis.Storage, appLog *logrus.Logger) *GoMiddleware {
	return &GoMiddleware{appCtx: ctx, redisCl: store, appLog: appLog}
}
