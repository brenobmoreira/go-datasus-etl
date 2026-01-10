package main

import (
	"fmt"
	"log"

	"github.com/brenobmoreira/go-datasus-etl/internal/datasus"
	"github.com/brenobmoreira/go-datasus-etl/internal/handlers"
	"github.com/brenobmoreira/go-datasus-etl/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	var downloadArchives = []datasus.Info{
		{UF: "SC", Ano: "25", Mes: "01"},
	}
	fmt.Println("Iniciando download...")
	if err := datasus.DownloadDBC(downloadArchives); err != nil {
		log.Fatalf("Erro no download: %v", err)
	}

	connection := "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable"
	repo, err := repository.OpenConn(connection)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/descricoes", handlers.GetDescricoes(&repo))
	router.GET("/sync", handlers.SyncDatabase(&repo))
	router.GET("/estabelecimento/:id", handlers.GetEstabelecimentosCidade(&repo))
	router.Run(":8080")
}
