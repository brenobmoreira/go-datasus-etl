package repository

import (
	"database/sql"
	"fmt"

	"github.com/brenobmoreira/go-datasus-etl/internal/entities"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func OpenConn(connection string) (Repo, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	return Repo{db: db}, err
}

func (r *Repo) SalvarEstabelecimento(st chan entities.Estabelecimento) error {
	sql := `INSERT INTO estabelecimento (cnes, codigo_municipio, competencia) VALUES ($1, $2, $3) ON CONFLICT (cnes, competencia) DO UPDATE SET codigo_municipio = EXCLUDED.codigo_municipio`
	for wt := range st {
		_, err := r.db.Exec(sql, wt.ID, wt.CodigoMunicipio, wt.Competencia)
		if err != nil {
			fmt.Printf("Erro ao inserir estabelecimento %s (%s): %v\n", wt.ID, wt.CodigoMunicipio, err)
		}
	}
	return nil
}

func (r *Repo) SalvarEquipamento(eq chan entities.Equipamentos) error {
	sql := `INSERT INTO equipamento (cnes, codigo_equipamento, quantidade_existente, quantidade_uso, competencia) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (cnes, codigo_equipamento, competencia) DO UPDATE SET quantidade_existente = EXCLUDED.quantidade_existente, quantidade_uso = EXCLUDED.quantidade_uso`
	for wt := range eq {
		_, err := r.db.Exec(sql, wt.ID, wt.CodigoEquipamento, wt.QuantidadeExistente, wt.QuantidadeUso, wt.Competencia)
		if err != nil {
			fmt.Printf("Erro ao inserir equipamento %s (%s): %v\n", wt.ID, wt.CodigoEquipamento, err)
		}
	}
	return nil
}

func (r *Repo) SalvarCadastro(cd chan entities.EstabelecimentoCadastro) error {
	sql := `INSERT INTO estabelecimento_cadastro (cnes, nome) VALUES ($1, $2) ON CONFLICT (cnes) DO UPDATE SET nome = EXCLUDED.nome`
	for wt := range cd {
		_, err := r.db.Exec(sql, wt.ID, wt.Nome)
		if err != nil {
			fmt.Printf("Erro ao inserir cadastro %s (%s): %v\n", wt.ID, wt.Nome, err)
		}
	}
	return nil
}

func (r *Repo) SalvarDescricao(dc chan entities.EquipamentoDescricao) error {
	sql := `INSERT INTO equipamento_descricao (codigo, descricao) VALUES ($1, $2) ON CONFLICT (codigo) DO UPDATE SET descricao = EXCLUDED.descricao`
	for wt := range dc {
		_, err := r.db.Exec(sql, wt.CodigoEquipamento, wt.Descricao)
		if err != nil {
			fmt.Printf("Erro ao inserir descrição %s (%s): %v\n", wt.CodigoEquipamento, wt.Descricao, err)
		}
	}
	return nil
}
