package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	_ "github.com/lib/pq"
	db "github.com/meomeocoj/simplebank/db/sqlc"
	"github.com/meomeocoj/simplebank/token"
	"github.com/meomeocoj/simplebank/utils"
)

type Server struct {
	store               db.Store
	router              *gin.Engine
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	// Setup currency validator

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server := &Server{
		store:               store,
		tokenMaker:          maker,
		accessTokenDuration: config.AccessTokenDuration}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.login)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
