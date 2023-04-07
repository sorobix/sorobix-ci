package main

import (
	"fmt"
	"github.com/genjidb/genji"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"net/http"
	"os"
)

func main() {
	_, err := os.Stat(SCRIPT)
	if err != nil {
		//todo: log saying that script not found
		panic(err)
	}
	db := initDb()
	repo := NewRepo(db)
	err = repo.CreateTable()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	runRestServer(repo)
}

func initDb() *genji.DB {
	db, err := genji.Open(DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func runRestServer(repo Repository) {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "true",
			"details": "sorobix-ci-daemon",
			"author":  "Hemanth Krishna <@DarthBenro008>",
		})
	})

	app.Post("/update/:commit/:id", func(ctx *fiber.Ctx) error {
		key := ctx.Params("id", "")
		sha := ctx.Params("commit", "")
		if key != SECRET {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(fiber.Map{"status": false,
				"message": "invalid key"})
		}
		uuid, err := deployer(repo)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(fiber.Map{"status": false,
				"message": err.Error()})
		}
		return ctx.JSON(fiber.Map{"status": true,
			"message": fmt.Sprintf(
				"deployed commit %s succesfully with deployment id: %s",
				sha, uuid)})
	})

	app.Post("/fetchdeps", func(ctx *fiber.Ctx) error {
		deps, err := repo.FetchDeployments()
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(fiber.Map{"status": false,
				"message": err.Error()})
		}
		return ctx.JSON(fiber.Map{"status": true,
			"data": deps})
	})

	log.Fatal(app.Listen(":6969"))
}
