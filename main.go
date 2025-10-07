package main

import (
	"fmt"
	"log"

	"github.com/brenobmoreira/go-datasus-etl/internal/datasus"
)

func main() {
	var downloadArchives = []datasus.Info{
		{UF: "SC", Ano: "25", Mes: "01"},
	}
	fmt.Println("Iniciando download...")
	if err := datasus.DownloadDBC(downloadArchives); err != nil {
		log.Fatalf("Erro no download: %v", err)
	}
}
