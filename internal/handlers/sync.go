package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/brenobmoreira/go-datasus-etl/internal/entities"
	"github.com/brenobmoreira/go-datasus-etl/internal/parser"
	"github.com/brenobmoreira/go-datasus-etl/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func ReturnRootDir() (rootDir string) {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return rootDir
}

func SyncDatabase(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var g errgroup.Group

		err := SyncCadastro(repo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		g.Go(func() error {
			return SyncEstabelecimento(repo)
		})

		g.Go(func() error {
			return SyncEquipamento(repo)
		})

		g.Go(func() error {
			return SyncDescricao(repo)
		})

		if err := g.Wait(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Sync completo com sucesso!"})

	}
}

func SyncCadastro(repo *repository.Repo) error {
	rootDir := ReturnRootDir()
	var wg sync.WaitGroup
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")

	cadastroChan := make(chan entities.EstabelecimentoCadastro)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.SalvarCadastro(cadastroChan); err != nil {
			fmt.Printf("Erro ao salvar cadastro: %v\n", err)
		}
	}()

	parser.CadastroParser("CADGERSC", blast_path, rootDir, cadastroChan)
	close(cadastroChan)
	wg.Wait()

	return nil
}

func SyncEstabelecimento(repo *repository.Repo) error {
	rootDir := ReturnRootDir()
	var wg sync.WaitGroup
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")
	competencia := time.Now()

	estabChan := make(chan entities.Estabelecimento)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.SalvarEstabelecimento(estabChan); err != nil {
			fmt.Printf("Erro ao salvar estabelecimento: %v\n", err)
		}
	}()

	parser.EstabelecimentoParser("ST/STSC2501", blast_path, rootDir, competencia, estabChan) // TODO: change path
	close(estabChan)
	wg.Wait()

	return nil
}

func SyncEquipamento(repo *repository.Repo) error {
	rootDir := ReturnRootDir()
	var wg sync.WaitGroup
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")
	competencia := time.Now()

	equipChan := make(chan entities.Equipamentos)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.SalvarEquipamento(equipChan); err != nil {
			fmt.Printf("Erro ao salvar equipamento: %v\n", err)
		}
	}()

	parser.EquipamentoParser("EQ/EQSC2501", blast_path, rootDir, competencia, equipChan) // TODO: change path
	close(equipChan)
	wg.Wait()

	return nil
}

func SyncDescricao(repo *repository.Repo) error {
	rootDir := ReturnRootDir()
	var wg sync.WaitGroup
	blast_path := filepath.Join(rootDir, "internal", "parser", "blast-dbf")

	descChan := make(chan entities.EquipamentoDescricao)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.SalvarDescricao(descChan); err != nil {
			fmt.Printf("Erro ao salvar descrição: %v\n", err)
		}
	}()

	parser.DescricaoParser("TP_EQUIPAM", blast_path, rootDir, descChan)
	close(descChan)
	wg.Wait()

	return nil
}
