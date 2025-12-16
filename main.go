package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

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

	competencia := time.Now()
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
	close(cadastroChan)
	wg.Wait()
	fmt.Println("Cadastro processado com sucesso!")

	estabChan := make(chan entities.Estabelecimento)
	archive_st := "ST/STSC2501"
	wg.Go(func() {
		if err := repo.SalvarEstabelecimento(estabChan); err != nil {
			panic(err)
		}
	})
	parser.EstabelecimentoParser(archive_st, blast_path, rootDir, competencia, estabChan)
	close(estabChan)
	wg.Wait()
	fmt.Println("Estabelecimentos processados com sucesso!")

	equipChan := make(chan entities.Equipamentos)
	archive_eq := "EQ/EQSC2501"
	wg.Go(func() {
		if err := repo.SalvarEquipamento(equipChan); err != nil {
			panic(err)
		}
	})
	parser.EquipamentoParser(archive_eq, blast_path, rootDir, competencia, equipChan)
	close(equipChan)
	wg.Wait()
	fmt.Println("Equipamentos processados com sucesso!")

	descricaoChan := make(chan entities.EquipamentoDescricao)
	archive_desc := "TP_EQUIPAM"
	wg.Go(func() {
		if err := repo.SalvarDescricao(descricaoChan); err != nil {
			panic(err)
		}
	})
	parser.DescricaoParser(archive_desc, blast_path, rootDir, descricaoChan)
	close(descricaoChan)
	wg.Wait()
	fmt.Println("Descrições processadas com sucesso!")
}
