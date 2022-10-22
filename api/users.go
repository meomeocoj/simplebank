package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/meomeocoj/simplebank/db/sqlc"
	"github.com/meomeocoj/simplebank/utils"
)

type createUserRequest struct {
	UserName string `json:"username" binding:"required,min=6,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username         string    `json:"username"`
	Fullname         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func newUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		Username:         user.Username,
		Fullname:         user.Fullname,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.UserName,
		HashPassword: hashPassword,
		Email:        req.Email,
		Fullname:     req.FullName,
	})

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := newUserResponse(user)
	ctx.JSON(http.StatusCreated, res)

}

type loginRequest struct {
	Username string `json:"username" binding:"required,min=6,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReponse struct {
	User        CreateUserResponse `json:"user"`
	AccessToken string             `json:"accessToken"`
}

func (s *Server) login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckIsHashPassword(user.HashPassword, req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Username, s.accessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginReponse{
		User:        newUserResponse(user),
		AccessToken: accessToken,
	})

}
