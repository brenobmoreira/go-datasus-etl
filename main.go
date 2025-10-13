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

	archive_st := "ST/STSC2501"
	archive_eq := "EQ/EQSC2501"
	archive_desc := "TP_EQUIPAM"
	archive_cd := "CADGERSC"
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")

	parser.EstabelecimentoParser(archive_st, blast_path, rootDir)
	parser.EquipamentoParser(archive_eq, blast_path, rootDir)
	parser.CadastroParser(archive_cd, blast_path, rootDir)
	parser.DescricaoParser(archive_desc, blast_path, rootDir)
}
