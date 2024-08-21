package gapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"github.com/lyb88999/Go-SimpleBank/pb"
	"github.com/lyb88999/Go-SimpleBank/token"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/lyb88999/Go-SimpleBank/worker"
)

// Server serves HTTP request for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}
	return server, nil
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
