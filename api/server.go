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
	server := &Server{
		store:               store,
		tokenMaker:          maker,
		accessTokenDuration: config.AccessTokenDuration}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/login", server.login)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
