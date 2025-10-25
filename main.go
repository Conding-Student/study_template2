package main

import (
	kakfaconsumer "chatbot/pkg/kakfa_consumer"
	"chatbot/pkg/logs"
	"os"

	//"chatbot/pkg/realtime"
	"chatbot/pkg/routers"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{ //Modified Logs
		//Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white} | \n${yellow}${body} ${white} | ${responseData}\n",
		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white}\n",
		TimeFormat: "01/02/2006 3:04 PM",
		TimeZone:   "Asia/Manila",
	}))

	// Connect to the Database
	database.ConnectToCAGABAYDB()
	// database.ConnectToEsystemDB()

	// Declare & initialize routes
	routers.SetupPublicRoutes(app)
	routers.SetupPublicRoutesB(app)
	routers.SetupPrivateRoutes(app)

	go kakfaconsumer.ConsumedLoanFromKafka()

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			id := c.Params("id")
			c.Status(http.StatusMethodNotAllowed)

			currentTime := time.Now()
			errorMessage := fmt.Sprintf("Whitelabel Error Page\n"+
				"This application has no explicit mapping for error, so you are seeing this as a fallback.\n\n"+
				"%s\n"+
				"There was an unexpected error (type=Method Not Allowed, status=405).\n"+
				"Request method '%s' not supported",
				currentTime.Format("Mon Jan 2 15:04:05 MST 2006"), c.Method())

			logs.LOSLogs(c, "Main", id, "405", errorMessage)
			return c.SendString(errorMessage)
		}

		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			id := c.Params("id")
			c.Status(http.StatusMethodNotAllowed)

			currentTime := time.Now()
			errorMessage := fmt.Sprintf("Whitelabel Error Page\n"+
				"This application has no explicit mapping for error, so you are seeing this as a fallback.\n\n"+
				"%s\n"+
				"There was an unexpected error (type=Method Not Allowed, status=405).\n"+
				"Request method '%s' not supported",
				currentTime.Format("Mon Jan 2 15:04:05 MST 2006"), c.Method())

			logs.LOSLogs(c, "Main", id, "405", errorMessage)
			return c.SendString(errorMessage)
		}

		return c.Next()
	})

	// Start the server
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "LOCALIOS":

		port := utils.GetEnv("PORT")
		localip := utils.GetEnv("localip")
		localip, err := encryptDecrypt.Decrypt(localip, utils.GetEnv("SECRET_KEY"))
		if err != nil {
			log.Fatal("Failed to decrypt localip: ", err)
		}

		log.Printf("Server running on port %s", port)
		app.Listen(localip + ":" + port)

	default:

		if utils.GetEnv("SSL") == "disabled" {
			log.Fatal(app.ListenTLS(
				fmt.Sprintf(":%s", utils.GetEnv("PORT")),
				utils.GetEnv("SSL_CERTIFICATE"),
				utils.GetEnv("SSL_KEY"),
			))
		} else {
			err := app.Listen(fmt.Sprintf(":%s", utils.GetEnv("PORT")))
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

// func main() {

// 	app := fiber.New()

// 	app.Use(cors.New())

// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "*",
// 		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
// 	}))

// 	app.Use(logger.New(logger.Config{ //Modified Logs
// 		//Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white} | \n${yellow}${body} ${white} | ${responseData}\n",
// 		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white}\n",
// 		TimeFormat: "01/02/2006 3:04 PM",
// 		TimeZone:   "Asia/Manila",
// 	}))
// 	realtime.Register(app)
// 	// Connect to the Database
// 	database.ConnectToCAGABAYDB()
// 	database.ConnectToEsystemDB()

// 	// Declare & initialize routes
// 	routers.SetupPublicRoutes(app)
// 	routers.SetupPublicRoutesB(app)
// 	routers.SetupPrivateRoutes(app)

// 	go kakfaconsumer.ConsumedLoanFromKafka()

// 	app.Use(func(c *fiber.Ctx) error {
// 		if c.Method() != fiber.MethodPost {
// 			id := c.Params("id")
// 			c.Status(http.StatusMethodNotAllowed)

// 			currentTime := time.Now()
// 			errorMessage := fmt.Sprintf("Whitelabel Error Page\n"+
// 				"This application has no explicit mapping for error, so you are seeing this as a fallback.\n\n"+
// 				"%s\n"+
// 				"There was an unexpected error (type=Method Not Allowed, status=405).\n"+
// 				"Request method '%s' not supported",
// 				currentTime.Format("Mon Jan 2 15:04:05 MST 2006"), c.Method())

// 			logs.LOSLogs(c, "Main", id, "405", errorMessage)
// 			return c.SendString(errorMessage)
// 		}

// 		return c.Next()
// 	})

// 	app.Use(func(c *fiber.Ctx) error {
// 		if c.Method() != fiber.MethodGet {
// 			id := c.Params("id")
// 			c.Status(http.StatusMethodNotAllowed)

// 			currentTime := time.Now()
// 			errorMessage := fmt.Sprintf("Whitelabel Error Page\n"+
// 				"This application has no explicit mapping for error, so you are seeing this as a fallback.\n\n"+
// 				"%s\n"+
// 				"There was an unexpected error (type=Method Not Allowed, status=405).\n"+
// 				"Request method '%s' not supported",
// 				currentTime.Format("Mon Jan 2 15:04:05 MST 2006"), c.Method())

// 			logs.LOSLogs(c, "Main", id, "405", errorMessage)
// 			return c.SendString(errorMessage)
// 		}

// 		return c.Next()
// 	})

// 	if utils.GetEnv("SSL") == "enabled" {
// 		log.Fatal(app.ListenTLS(
// 			fmt.Sprintf(":%s", utils.GetEnv("PORT")),
// 			utils.GetEnv("SSL_CERTIFICATE"),
// 			utils.GetEnv("SSL_KEY"),
// 		))
// 	} else {
// 		err := app.Listen(fmt.Sprintf(":%s", utils.GetEnv("PORT")))
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}
// 	}

// 	// Start the server
// 	port := ":19000" // Update the port based on your preference

// 	// log.Printf("Server running on port %s", port)
// 	// // app.Listen(port)
// 	app.Listen("10.27.1.34" + port)
// 	//app.Listen(port)
// }
