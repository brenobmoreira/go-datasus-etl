package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/brenobmoreira/go-datasus-etl/internal/datasus"
	"github.com/brenobmoreira/go-datasus-etl/internal/entities"
	"github.com/brenobmoreira/go-datasus-etl/internal/parser"
	"github.com/brenobmoreira/go-datasus-etl/internal/repository"
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

	connection := "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable"
	repo, err := repository.OpenConn(connection)
	if err != nil {
		panic(err)
	}

	// competencia := time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
	var wg sync.WaitGroup
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")

	cadastroChan := make(chan entities.EstabelecimentoCadastro)
	archive_cd := "CADGERSC"
	wg.Go(func() {
		if err := repo.SalvarCadastro(cadastroChan); err != nil {
			panic(err)
		}
	})
	parser.CadastroParser(archive_cd, blast_path, rootDir, cadastroChan)
	fmt.Println("Cadastro updated")

	// estabChan := make(chan entities.Estabelecimento)
	// archive_st := "ST/STSC2501"
	// wg.Go(func() {
	// 	if err := repo.SalvarEstabelecimento(estabChan); err != nil {
	// 		panic(err)
	// 	}
	// })
	// parser.EstabelecimentoParser(archive_st, blast_path, rootDir, competencia, estabChan)

	// archive_eq := "EQ/EQSC2501"
	// archive_desc := "TP_EQUIPAM"
	// parser.EquipamentoParser(archive_eq, blast_path, rootDir)
	// parser.DescricaoParser(archive_desc, blast_path, rootDir)
}
