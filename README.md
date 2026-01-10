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
- ✅ API REST com Gin para consulta de dados
- ✅ Endpoint de sincronização sob demanda (`/sync`)
- ✅ Consulta de estabelecimentos por município
- ✅ Listagem de descrições de equipamentos
- ✅ Exportação para CSV (opcional/debug)
- ✅ Suporte a múltiplas competências e UFs

## 🏗️ Arquitetura

```
go-datasus-etl/
├── main.go                      # Entry point, API e orquestração
├── internal/
│   ├── datasus/
│   │   └── ftp_download.go      # Download FTP de arquivos .dbc
│   ├── entities/
│   │   └── model.go             # Modelos de dados
│   ├── handlers/
│   │   ├── descricoes.go        # Handler GET /descricoes
│   │   └── sync.go              # Handler GET /sync
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

equipamento (PK: cnes, codigo_equipamento, competencia)
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
go get github.com/gin-gonic/gin
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

### Executar a aplicação
```bash
go run .
```

A aplicação irá:
1. Baixar e processar os dados do DATASUS
2. Iniciar o servidor HTTP na porta 8080

### API Endpoints

#### `GET /descricoes`
Lista todas as descrições de equipamentos cadastrados.

**Exemplo:**
```bash
curl http://localhost:8080/descricoes
```

**Resposta:**
```json
[
  {
    "CodigoEquipamento": "01",
    "Descricao": "RAIO X ATÉ 100 mA"
  },
  ...
]
```

#### `GET /estabelecimento/:codigo_municipio`
Lista estabelecimentos de um município específico com CNES, nome e competência.

**Exemplo:**
```bash
curl http://localhost:8080/estabelecimento/420005
```

**Resposta:**
```json
[
  ["2675412", "SECRETARIA DE SAUDE DE ABDON BATISTA", "01/01/2025"],
  ["7656033", "POSTO DE COLETA ABDON BATISTA", "01/01/2025"],
  ...
]
```

#### `GET /sync`
Força uma nova sincronização dos dados do DATASUS.

**Exemplo:**
```bash
curl http://localhost:8080/sync
```

### Consultar dados via PostgreSQL
```bash
psql "host=localhost port=5432 user=postgres password=0000 dbname=godatabase sslmode=disable"
```

Exemplos de queries:
```sql
-- Total de estabelecimentos cadastrados
SELECT COUNT(*) FROM estabelecimento_cadastro;

-- Estabelecimentos por município com nome
SELECT e.codigo_municipio, ec.nome, e.competencia
FROM estabelecimento e
JOIN estabelecimento_cadastro ec ON e.cnes = ec.cnes
WHERE e.codigo_municipio = '420005'
ORDER BY ec.nome;

-- Equipamentos por estabelecimento
SELECT ec.nome, eq.codigo_equipamento, ed.descricao, 
       eq.quantidade_existente, eq.quantidade_uso
FROM equipamento eq
JOIN estabelecimento_cadastro ec ON eq.cnes = ec.cnes
LEFT JOIN equipamento_descricao ed ON eq.codigo_equipamento = ed.codigo
WHERE eq.cnes = '2675412';

-- Top 10 equipamentos mais comuns
SELECT ed.descricao, SUM(eq.quantidade_existente) as total
FROM equipamento eq
LEFT JOIN equipamento_descricao ed ON eq.codigo_equipamento = ed.codigo
GROUP BY ed.descricao
ORDER BY total DESC
LIMIT 10;
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
5. Adicione handler em `internal/handlers`
6. Registre rota no `main.go`

## 🐛 Troubleshooting

### Erro: "Nao foi possivel abrir o arquivo de entrada .dbc"
- Verifique se os arquivos foram baixados em `data/dbc/`
- Confira permissões do binário `blast-dbf`: `chmod +x internal/parser/blast-dbf`

### Erro: "pq: duplicate key value violates unique constraint"
- Dados já foram inseridos anteriormente
- Solução: O código já usa `ON CONFLICT DO UPDATE` para atualizar registros existentes

### Erro: "connection refused" (PostgreSQL)
- Verifique se o container está rodando: `docker ps`
- Inicie o container: `docker start godatabase`
- Teste conexão: `psql "host=localhost port=5432 user=postgres password=0000 dbname=godatabase"`

### API retorna erro 404
- Verifique se o servidor está rodando na porta 8080
- Confirme a URL: `http://localhost:8080/`

## 📝 Roadmap

### 🎯 Prioridade Alta
- [ ] **Parâmetros de sincronização**: Permitir escolher UF, ano e mês via endpoint `/sync?uf=SC&ano=25&mes=01`
- [ ] **Variáveis de ambiente**: Mover credenciais do banco para `.env`
- [ ] **Validação de entrada**: Validar códigos de município e parâmetros da API
- [ ] **Tratamento de erros**: Melhorar mensagens de erro da API com códigos HTTP apropriados
- [ ] **Logging estruturado**: Implementar zerolog ou zap para logs padronizados

### 🚀 Médio Prazo
- [ ] **Suporte multi-UF**: Processar todas as UFs do Brasil, não apenas SC
- [ ] **Range de competências**: Baixar múltiplos meses/anos em uma única execução
- [ ] **Cache de dados**: Implementar cache Redis para queries frequentes
- [ ] **Paginação**: Adicionar paginação nos endpoints que retornam muitos registros
- [ ] **Filtros avançados**: Permitir filtrar por tipo de equipamento, período, etc.
- [ ] **Autenticação**: Implementar JWT para proteger endpoints sensíveis

### 🔮 Longo Prazo
- [ ] **Dashboard web**: Interface visual com gráficos e métricas (React/Vue)
- [ ] **Notificações**: Alertas quando novos dados estiverem disponíveis
- [ ] **Exportação de relatórios**: PDF/Excel com dados agregados
- [ ] **API GraphQL**: Alternativa ao REST para queries complexas
- [ ] **Suporte a outros datasets DATASUS**: SIA, SIH, SINASC, etc.
- [ ] **Machine Learning**: Análise preditiva de demanda de equipamentos

### 🧪 Qualidade e DevOps
- [ ] **Testes unitários**: Cobertura mínima de 80%
- [ ] **Testes de integração**: Validar fluxo completo do ETL
- [ ] **CI/CD**: GitHub Actions para build, test e deploy
- [ ] **Docker Compose**: Ambiente completo com um comando
- [ ] **Kubernetes**: Manifests para deploy em produção
- [ ] **Monitoramento**: Prometheus + Grafana para métricas
- [ ] **Documentação API**: OpenAPI/Swagger

### 🐛 Correções Conhecidas
- [ ] **Sincronização de canais**: Corrigir WaitGroup e garantir que todos os goroutines finalizem
- [ ] **Retry logic FTP**: Implementar tentativas automáticas em caso de falha no download
- [ ] **Limpeza de arquivos temporários**: Deletar `.dbc` e `.dbf` após processamento

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:
1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/NovaFuncionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/NovaFuncionalidade`)
5. Abra um Pull Request

**Áreas que precisam de ajuda:**
- Testes automatizados
- Documentação de código
- Suporte a novas UFs
- Performance do parser DBF
- Interface web

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🔗 Links Úteis

- [DATASUS FTP](ftp://ftp.datasus.gov.br/dissemin/publicos/CNES/200508_/Dados/)
- [Documentação CNES](http://cnes.datasus.gov.br/)
- [Dicionário de Dados DATASUS](http://tabnet.datasus.gov.br/tabdata/cadernos/cnes.htm)
- [Documentação Go](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)

## 👤 Autor

**Breno Moreira**
- GitHub: [@brenobmoreira](https://github.com/brenobmoreira)

---

⭐ Se este projeto foi útil, considere dar uma estrela no repositório!

💡 **Sugestões?** Abra uma [issue](https://github.com/brenobmoreira/go-datasus-etl/issues) ou contribua com código!
