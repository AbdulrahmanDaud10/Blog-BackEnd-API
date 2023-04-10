package controllers

import (
	"net/http"

	"github.com/AbdulrahmanDaud10/fullstack-project/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to the API")
}
