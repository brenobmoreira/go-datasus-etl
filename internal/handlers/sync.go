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
)

func SyncDatabase(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		rootDir, err := os.Getwd()
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		})
		parser.CadastroParser(archive_cd, blast_path, rootDir, cadastroChan)
		close(cadastroChan)
		wg.Wait()
		fmt.Println("Cadastro processado com sucesso!")
		c.JSON(http.StatusOK, "Cadastro processado com sucesso!")

		estabChan := make(chan entities.Estabelecimento)
		archive_st := "ST/STSC2501"
		wg.Go(func() {
			if err := repo.SalvarEstabelecimento(estabChan); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		})
		parser.EstabelecimentoParser(archive_st, blast_path, rootDir, competencia, estabChan)
		close(estabChan)
		wg.Wait()
		fmt.Println("Estabelecimentos processados com sucesso!")
		c.JSON(http.StatusOK, "Estabelecimentos processados com sucesso!")

		equipChan := make(chan entities.Equipamentos)
		archive_eq := "EQ/EQSC2501"
		wg.Go(func() {
			if err := repo.SalvarEquipamento(equipChan); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		})
		parser.EquipamentoParser(archive_eq, blast_path, rootDir, competencia, equipChan)
		close(equipChan)
		wg.Wait()
		fmt.Println("Equipamentos processados com sucesso!")
		c.JSON(http.StatusOK, "Equipamentos processados com sucesso!")

		descricaoChan := make(chan entities.EquipamentoDescricao)
		archive_desc := "TP_EQUIPAM"
		wg.Go(func() {
			if err := repo.SalvarDescricao(descricaoChan); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		})
		parser.DescricaoParser(archive_desc, blast_path, rootDir, descricaoChan)
		close(descricaoChan)
		wg.Wait()
		fmt.Println("Descrições processadas com sucesso!")
		c.JSON(http.StatusOK, "Descrições processadas com sucesso!")
	}
}
