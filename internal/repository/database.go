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
	sql := `INSERT INTO estabelecimento (cnes, codigo_municipio, competencia) VALUES ($1, $2, $3)`
	for wt := range st {
		_, err := r.db.Exec(sql, wt.ID, wt.CodigoMunicipio, wt.Competencia)
		if err != nil {
			fmt.Printf("Erro ao inserir %s, %v", wt.ID, wt.CodigoMunicipio)
		}
	}
	return nil
}

func (r *Repo) SalvarEquipamento(eq entities.Equipamentos) error {
	sql := `INSERT INTO equipamento (cnes, codigo_equipamento, quantidade_existente, quantidade_uso, competencia) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (cnes, competencia) DO NOTHING`
	_, err := r.db.Exec(sql, eq.ID, eq.CodigoEquipamento, eq.QuantidadeExistente, eq.QuantidadeUso)
	return err
}

func (r *Repo) SalvarCadastro(cd chan entities.EstabelecimentoCadastro) error {
	sql := `INSERT INTO estabelecimento_cadastro (cnes, nome) VALUES ($1, $2) ON CONFLICT (cnes) DO UPDATE SET cnes = EXCLUDED.cnes, nome = EXCLUDED.nome;`
	for wt := range cd {
		_, err := r.db.Exec(sql, wt.ID, wt.Nome)
		if err != nil {
			fmt.Println("Erro ao inserir: ", wt.ID, wt.Nome, err)
		}
	}
	return nil
}

func (r *Repo) SalvarDescricao(st entities.EquipamentoDescricao) error {
	sql := `INSERT INTO equipamento_descricao (codigo_equipamento, descricao) VALUES ($1, $2)`
	_, err := r.db.Exec(sql, st.CodigoEquipamento, st.Descricao)
	return err
}
