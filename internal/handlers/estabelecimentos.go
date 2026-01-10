package handlers

import (
	"net/http"

	"github.com/brenobmoreira/go-datasus-etl/internal/repository"
	"github.com/gin-gonic/gin"
)

func GetEstabelecimentosCidade(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		codigo := c.Param("id")
		estabelecimentos, err := repo.ListarEstabelecimentoCidade(codigo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, estabelecimentos)
	}
}
