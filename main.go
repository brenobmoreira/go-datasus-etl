package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/brenobmoreira/go-datasus-etl/internal/datasus"
	"github.com/brenobmoreira/go-datasus-etl/internal/parser"
)

func main() {
	var downloadArchives = []datasus.Info{
		{UF: "SC", Ano: "25", Mes: "01"},
	}
	fmt.Println("Iniciando download...")
	if err := datasus.DownloadDBC(downloadArchives); err != nil {
		log.Fatalf("Erro no download: %v", err)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	archive_name := "ST/STSC2501"
	archive_eq := "EQ/EQSC2501"
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")

	parser.EstabelecimentoParser(archive_name, blast_path, rootDir)
	parser.EquipamentoParser(archive_eq, blast_path, rootDir)
}
