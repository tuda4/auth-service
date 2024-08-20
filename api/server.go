package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	config "github.com/tuda4/mb-backend/configs"
	db "github.com/tuda4/mb-backend/db/sqlc"
	"github.com/tuda4/mb-backend/internal/val"
	"github.com/tuda4/mb-backend/token"
	"github.com/tuda4/mb-backend/worker"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/tuda4/mb-backend/docs"
)

type Server struct {
	store           db.Store
	token           token.Maker
	router          *echo.Echo
	config          config.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config config.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:           store,
		token:           tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()
	return server, nil
}

// @title Mobile Project API
// @version 1.0
// @description Document API for Mobile Project

// @contact.name API Support
// @contact.url https://www.instagram.com/imtuda4
// @contact.email imtuda4@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:8080
// @BasePath /api/web/v1
// @Accept		json
// @Produce	json
func (server *Server) setupRouter() {

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		fmt.Println("health check")
		return c.String(200, "OK")
	})
	e.Validator = &val.CustomValidator{Validator: validator.New()}
	// e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(CustomLoggerConfig))
	// health check

	web := e.Group("/api/web/v1")

	web.POST("/login", server.login)
	web.POST("/signup", server.createAccount)
	web.POST("/refresh", server.refreshToken)
	web.GET("/verify", server.verifyEmail)

	authRouter := web.Group("/")
	authRouter.Use(server.authMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	server.router = e
}

func (server *Server) Start(address string) error {
	return server.router.Start(address)
}
