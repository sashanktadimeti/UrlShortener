package routes
import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
	"urlshortener/database/database"
)
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()
	val,err := r.Get(database.Ctx,url).Result()
	if err == redis.Nil{
		return  c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"URL Not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errror":"Internal Server Error"})
	}else {
		rInr := database.CreateClient(1)
		defer rInr.Close()
		rInr.Incr(database.Ctx,url)
		return c.Redirect(val, fiber.StatusTemporaryRedirect)
		
	}

}