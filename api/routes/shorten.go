package routes

import (
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"urlshortener/database/database"
	"urlshortener/helpers"
)

type Request struct {
    URL         string        `json:"url"`
    CustomShort string        `json:"custom_short"`
    Expiry      time.Duration `json:"expiry"`
}

type Response struct {
    URL             string        `json:"url"`
    CustomShort     string        `json:"custom_short"`
    Expiry          int64         `json:"expiry"`
    XRateRemaining  int64         `json:"x_rate_remaining"`
    XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
    body := new(Request)

    if err := c.BodyParser(body); err != nil {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"error": "cannot parse JSON"})
    }
	// Rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	val,err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		r2.Set(database.Ctx,c.IP(), os.Getenv("API_QUOTA"),30*60*time.Second)
	} else {
		valInt, _ := govalidator.ToInt(val)
		if valInt <= 0 {
			limit , _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{"error": "API rate limit exceeded",
					"x_rate_limit_reset": int(limit.Minutes())} )
		}
	}
    // URL validation
    if !govalidator.IsURL(body.URL) {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"error": "invalid URL"})
    }

    // Domain validation
    if !helpers.RemoveDomainError(body.URL) {
        return c.Status(fiber.StatusServiceUnavailable).
            JSON(fiber.Map{"error": "you are not allowed to shorten this domain"})
    }

    // Enforce http scheme
    body.URL = helpers.EnforceHTTP(body.URL)
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else{
		id = body.CustomShort
	}
	r := database.CreateClient(0)
	defer r.Close()
	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error":"Custom short URL is already in use"})
	} else {
		if body.Expiry == 0 {
			body.Expiry = 24
		}
		err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Internal Server Error"})
		}
	}
	resp := Response{
		URL: body.URL,
		CustomShort: "",
		Expiry: int64(body.Expiry),
		XRateRemaining: 0,
		XRateLimitReset: 0,
	}
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	remainingInt, _ := govalidator.ToInt(val)
	resp.XRateRemaining = remainingInt - 1
	limit , _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = time.Duration(limit.Minutes())
	r2.Decr(database.Ctx,c.IP())
    return c.Status(fiber.StatusOK).JSON(resp)
}
