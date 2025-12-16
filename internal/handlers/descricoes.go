package handlers

import (
	"net/http"

	"github.com/brenobmoreira/go-datasus-etl/internal/repository"
	"github.com/gin-gonic/gin"
)

func GetDescricoes(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		descricoes, err := repo.ListarDescricoes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, descricoes)
	}
}
