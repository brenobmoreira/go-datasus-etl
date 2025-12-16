# DATASUS ETL - Go

Sistema ETL (Extract, Transform, Load) em Go para processar dados do DATASUS/CNES (Cadastro Nacional de Estabelecimentos de Saúde).

## 📋 Sobre o Projeto

Este projeto automatiza a extração, transformação e carga de dados públicos do DATASUS, focando em:
- **Cadastro de estabelecimentos de saúde** (CNES)
- **Equipamentos médicos** por estabelecimento
- **Dados mensais de estabelecimentos** por competência

Os dados são baixados via FTP do servidor público do DATASUS, convertidos do formato proprietário `.dbc` para `.dbf`, processados e armazenados em PostgreSQL.

## 🚀 Funcionalidades

- ✅ Download automatizado de arquivos do FTP DATASUS
- ✅ Conversão de arquivos `.dbc` → `.dbf` (formato dBase)
- ✅ Parse de arquivos DBF para estruturas Go
- ✅ Carga concorrente em PostgreSQL via canais
- ✅ Exportação para CSV (opcional/debug)
- ✅ Suporte a múltiplas competências e UFs

## 🏗️ Arquitetura

```
go-datasus-etl/
├── main.go                      # Entry point e orquestração
├── internal/
│   ├── datasus/
│   │   └── ftp_download.go      # Download FTP de arquivos .dbc
│   ├── entities/
│   │   └── model.go             # Modelos de dados
│   ├── parser/
│   │   ├── blast-dbf            # Binário conversor dbc→dbf
│   │   ├── dbase_helpers.go     # Helpers conversão/leitura DBF
│   │   ├── cadastro_parser.go   # Parser cadastro geral
│   │   ├── eq_mensal_parser.go  # Parser equipamentos mensais
│   │   ├── st_mensal_parser.go  # Parser estabelecimentos mensais
│   │   └── descricao_parser.go  # Parser descrições equipamentos
│   └── repository/
│       └── database.go          # Persistência PostgreSQL
├── migrations/
│   └── init.sql                 # Schema do banco de dados
├── data/
│   ├── dbc/                     # Arquivos .dbc baixados
│   ├── dbf/                     # Arquivos .dbf convertidos
│   └── csv/                     # CSVs exportados
└── assets/                      # Arquivos estáticos/base

```

## 🗄️ Schema do Banco de Dados

```sql
estabelecimento_cadastro (PK: cnes)
├── cnes (TEXT)
└── nome (TEXT)

estabelecimento (PK: cnes, competencia)
├── cnes (TEXT) → FK estabelecimento_cadastro
├── codigo_municipio (TEXT)
└── competencia (DATE)

equipamento (PK: cnes, competencia)
├── cnes (TEXT) → FK estabelecimento_cadastro
├── codigo_equipamento (TEXT)
├── quantidade_existente (INTEGER)
├── quantidade_uso (INTEGER)
└── competencia (DATE)

equipamento_descricao (PK: codigo)
├── codigo (TEXT)
└── descricao (TEXT)
```

## 🔧 Pré-requisitos

- **Go 1.21+**
- **PostgreSQL 13+**
- **Docker** (opcional, para rodar Postgres em container)

## 📦 Dependências

```bash
go get github.com/valentin-kaiser/go-dbase/dbase
go get github.com/jlaffaye/ftp
go get github.com/lib/pq
```

## ⚙️ Instalação e Configuração

### 1. Clone o repositório
```bash
git clone https://github.com/brenobmoreira/go-datasus-etl.git
cd go-datasus-etl
```

### 2. Configure o banco de dados

#### Opção A: Docker (recomendado)
```bash
docker run -d \
  --name godatabase \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=0000 \
  -e POSTGRES_DB=godatabase \
  -p 5432:5432 \
  postgres
```

#### Opção B: PostgreSQL local
Certifique-se que o PostgreSQL está rodando na porta 5432.

### 3. Execute as migrations
```bash
psql "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable" \
  -f migrations/init.sql
```

### 4. Ajuste a configuração (main.go)
Edite a string de conexão e os parâmetros de download conforme necessário:
```go
connection := "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable"

var downloadArchives = []datasus.Info{
    {UF: "SC", Ano: "25", Mes: "01"}, // Santa Catarina, Janeiro/2025
}
```

## 🎯 Uso

### Executar o ETL completo
```bash
go run .
```

### Processo executado
1. Download de arquivos `.dbc` do FTP DATASUS
2. Conversão `.dbc` → `.dbf` usando `blast-dbf`
3. Parse dos dados e envio via canais
4. Inserção concorrente no PostgreSQL
5. Geração de CSVs (opcional)

### Consultar dados processados
```bash
psql "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable"
```

Exemplos de queries:
```sql
-- Total de estabelecimentos cadastrados
SELECT COUNT(*) FROM estabelecimento_cadastro;

-- Estabelecimentos por município
SELECT codigo_municipio, COUNT(*) 
FROM estabelecimento 
GROUP BY codigo_municipio;

-- Equipamentos por tipo
SELECT e.codigo_equipamento, ed.descricao, SUM(e.quantidade_existente) as total
FROM equipamento e
LEFT JOIN equipamento_descricao ed ON e.codigo_equipamento = ed.codigo
GROUP BY e.codigo_equipamento, ed.descricao
ORDER BY total DESC;
```

## 📊 Tipos de Dados Processados

| Arquivo | Descrição | Frequência |
|---------|-----------|------------|
| `CADGERSC` | Cadastro geral de estabelecimentos SC | Base/Atualização |
| `STSC{AA}{MM}` | Estabelecimentos mensais SC | Mensal |
| `EQSC{AA}{MM}` | Equipamentos mensais SC | Mensal |
| `TP_EQUIPAM` | Descrição tipos de equipamento | Tabela auxiliar |

*AA = ano (25 = 2025), MM = mês (01-12)*

## 🛠️ Desenvolvimento

### Estrutura de Entities
```go
type EstabelecimentoCadastro struct {
    ID       string `dbase:"CNES"`
    Nome     string `dbase:"FANTASIA"`
    Excluido string `dbase:"EXCLUIDO"`
}

type Estabelecimento struct {
    ID              string    `dbase:"CNES"`
    CodigoMunicipio string    `dbase:"CODUFMUN"`
    Competencia     time.Time
}

type Equipamentos struct {
    ID                  string    `dbase:"CNES"`
    CodigoEquipamento   string    `dbase:"TIPEQUIP"`
    QuantidadeExistente int64     `dbase:"QT_EXIST"`
    QuantidadeUso       int64     `dbase:"QT_USO"`
    Competencia         time.Time
}
```

### Adicionar novos parsers
1. Crie arquivo em `internal/parser/` (ex: `novo_parser.go`)
2. Implemente função `NovoParser(archive_name, blast, dir string)`
3. Adicione entity correspondente em `internal/entities/model.go`
4. Crie método `SalvarNovo()` em `internal/repository/database.go`
5. Chame no `main.go`

## 🐛 Troubleshooting

### Erro: "Nao foi possivel abrir o arquivo de entrada .dbc"
- Verifique se os arquivos foram baixados em `data/dbc/`
- Confira permissões do binário `blast-dbf`: `chmod +x internal/parser/blast-dbf`

### Erro: "pq: duplicate key value violates unique constraint"
- Dados já foram inseridos anteriormente
- Solução: adicione `ON CONFLICT DO NOTHING` nas queries de inserção

### Erro: "connection refused" (PostgreSQL)
- Verifique se o container está rodando: `docker ps`
- Inicie o container: `docker start godatabase`
- Teste conexão: `psql "host=localhost port=5432 user=postgres password=0000 dbname=godatabase"`

## 📝 Roadmap

- [ ] Corrigir sincronização de canais e WaitGroup
- [ ] Adicionar logging estruturado (zerolog/zap)
- [ ] Implementar retry logic no FTP download
- [ ] Suporte a mais UFs além de SC
- [ ] API REST para consulta de dados
- [ ] Dashboard web com métricas
- [ ] Testes unitários e integração
- [ ] CI/CD com GitHub Actions
- [ ] Dockerização completa (multi-stage build)

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:
1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/NovaFuncionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/NovaFuncionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🔗 Links Úteis

- [DATASUS FTP](ftp://ftp.datasus.gov.br/dissemin/publicos/CNES/200508_/Dados/)
- [Documentação CNES](http://cnes.datasus.gov.br/)
- [Dicionário de Dados DATASUS](http://tabnet.datasus.gov.br/tabdata/cadernos/cnes.htm)

## 👤 Autor

**Breno Moreira**
- GitHub: [@brenobmoreira](https://github.com/brenobmoreira)

---

⭐ Se este projeto foi útil, considere dar uma estrela no repositório!
