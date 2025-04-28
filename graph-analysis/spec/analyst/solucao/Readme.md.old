## Resumo do Projeto

Este repositório apresenta a solução desenvolvida para o teste técnico da Neoway — Pessoa Desenvolvedora de Software Pleno (Graph Analysis). A proposta consiste na implementação de um pipeline completo de ingestão e consulta de dados com foco em análise de grafos utilizando Golang e Neo4j.

### Objetivos atendidos:

- Implementar uma ETL em Golang para leitura contínua de arquivos CSV, transformação e carga no banco de grafos Neo4j.
- Desenvolver uma API REST com suporte a múltiplos endpoints parametrizáveis que consultam o banco e respondem a perguntas específicas.
- Containerizar a solução com Docker, incluindo docker-compose e um Makefile com comandos úteis para desenvolvimento e execução.
- Garantir clareza de estrutura, separação de responsabilidades e justificativas técnicas para cada camada da aplicação (ETL, API, repositórios, etc).

### Decisões técnicas justificadas:

- A arquitetura foi baseada em Clean Architecture com separação por camadas (`handler`, `usecase`, `repository`) para facilitar manutenção, testes e reuso de código.
- A estratégia de ingestão segue o modelo Medallion (Bronze, Silver, Gold), conforme sugerido no livro *Understanding ETL (O'Reilly, 2024)*, o que facilita rastreabilidade, versionamento e reprocessamento de dados.
- A conexão com o banco Neo4j foi centralizada em um único client reaproveitável tanto pela API quanto pela ETL.
- A geração de dados simulados foi feita por um script Python embarcado no mesmo docker-compose, garantindo reprodutibilidade.

### O que poderia ser melhorado (fora do escopo por limitação de tempo)

Apesar da estrutura sólida, algumas melhorias que poderiam ser implementadas com mais tempo incluem:

- Paralelização da leitura de arquivos CSV com goroutines:
  Seria possível utilizar um Worker Pool com channels para paralelizar a leitura e processamento dos arquivos CSV, o que aumentaria a performance da ingestão, principalmente com arquivos grandes ou múltiplos simultâneos.

- Melhorias de performance na API:
  A API poderia utilizar cache em memória (ex: Redis) para respostas que não mudam com frequência, como datas de aprovação de vacinas. Além disso, middlewares como rate limiter e gzip poderiam ser adicionados para ambientes produtivos.

- Documentação automática da API com Swagger:
  A biblioteca `swaggo/swag` permite anotar os handlers em Go com comentários estruturados, e gerar a interface interativa do Swagger UI. Por falta de tempo, a documentação foi feita manualmente no README, mas o projeto já está estruturado para aceitar a integração com `gin-swagger`.



## Arquitetura da Solução

A estrutura do projeto foi dividida em três componentes principais:

- `csv-generator`: simula dumps de bancos relacionais gerando CSVs.
- `etl`: processo contínuo de leitura, transformação e carga no grafo.
- `api`: serviço REST com consultas ao banco.

### Organização Geral

```
.
├── cmd/
│   ├── api/               # Entrada principal da API REST
│   └── etl/               # Entrada principal da ETL contínua
├── internal/
│   ├── api/               # Handlers HTTP e roteamento
│   ├── config/            # Leitura de variáveis de ambiente
│   ├── loader/            # Lógica de escrita no Neo4j (Gold Layer)
│   ├── model/             # Definição das entidades
│   ├── reader/            # Leitura e parsing dos arquivos CSV (Silver Layer)
│   ├── repository/        # Conexão com banco (Neo4j client reutilizável)
│   ├── usecase/           # Lógica de negócio da API
│   └── util/              # Funções auxiliares (ex: parse de datas)
├── data/                  # Diretório onde os CSVs simulados são gerados
├── docker/                # Dockerfiles e configs
├── scripts/               # Gerador de CSVs em Python (simulando dumps SQL)
├── .env                   # Variáveis de ambiente
├── docker-compose.yml
└── Makefile
```

### Estratégia de Dados: Arquitetura Medallion

A arquitetura segue o modelo Medallion (Bronze, Silver, Gold), conforme descrito no livro *Understanding ETL (O’Reilly, 2024)*:

- **Bronze**: dados brutos gerados continuamente por um componente simulado chamado `csv-generator`, implementado em Python. Os dados são escritos em arquivos `.csv` com timestamps distintos.
- **Silver**: leitura e parsing dos dados, validação de tipos e padronização (ex: conversão de datas com `time.Parse`, nomes de países e vacinas, valores nulos como ponteiros).
- **Gold**: carga no banco de grafos Neo4j com garantias de unicidade (`MERGE`, constraints), relacionamentos entre entidades e modelagem orientada a queries analíticas.

### Decisões de Modelagem

Optou-se por usar **chaves de negócio** (ex: `iso3`, `name`, `date`) no lugar de IDs artificiais, pois são mais expressivas e já atendem os requisitos das queries do desafio. Isso facilita o uso de `MERGE` e evita duplicidade.

| Entidade           | ID fornecido | Chave usada no grafo | Justificativa técnica                         |
|--------------------|--------------|-----------------------|-----------------------------------------------|
| Country            | id           | iso3                  | É único por país e usado nas queries          |
| Vaccine            | id           | name                  | Nome é suficiente e referenciado diretamente |
| CovidCase          | id           | date                  | Unicidade por país + data                     |
| VaccinationStats   | id           | date                  | Igual ao anterior                             |
| VaccineApproval    | id           | vaccine               | Nome da vacina já é suficiente                |





### Reuso de Componentes

A conexão com o banco Neo4j foi encapsulada em um client reutilizável (`internal/repository/neo4j/client.go`), utilizado tanto pela API quanto pela ETL. Isso evita duplicação de código e facilita manutenção.

### Observações sobre Escalabilidade

A arquitetura atual permite evoluções como:

- Substituir a origem dos dados (CSV) por bancos relacionais, APIs ou filas, com mudanças mínimas no pacote `reader/`.
- Trocar o destino do grafo Neo4j por outro banco, alterando apenas o repositório.
- Paralelizar leitura e escrita com goroutines (não implementado por limitação de tempo).
- Criar múltiplas instâncias do ETL ou da API para suportar carga horizontal.

## Tratamento dos Dados Lidos

Durante a leitura dos CSVs, são aplicadas as seguintes validações e transformações:

- Conversão de datas com `time.Parse`, rejeitando valores inválidos.
- Campos nulos tratados com ponteiros (`*int`, `*string`).
- Validação de headers esperados.
- Ignora registros incompletos ou duplicados.

O tratamento de datas é especialmente importante para garantir consistência ao carregar informações temporais no grafo. Futuramente podem ser adicionadas:

- Validações de ranges temporais.
- Verificação de colunas obrigatórias por tipo de entidade.
- Normalização de nomes com equivalência semântica (ex: EUA, United States).

## Estrutura dos Arquivos CSV

Durante a entrevista técnica foi mencionado que, na prática, os dados da Neoway são frequentemente extraídos de sistemas relacionais (bancos SQL) e disponibilizados para processamento no formato CSV. A partir dessa informação, a geração dos dados deste projeto foi planejada para **simular dumps reais de bancos relacionais**, respeitando formatos, periodicidade e estruturas comuns a esse tipo de ambiente.

### Motivação

Ao invés de usar um único arquivo estático de entrada, foi criado um componente chamado `csv-generator`, que gera novos arquivos CSV periodicamente. Isso simula o comportamento de um sistema ETL upstream exportando dados continuamente, como ocorre em pipelines reais de produção.

Essa estratégia permite:

- Testar o comportamento do ETL de forma contínua e realista.
- Validar a resiliência do sistema a novos dados chegando em ciclos.
- Simular casos de versionamento de dumps (por data ou timestamp).

### Características da Estrutura

Os arquivos CSV seguem uma estrutura próxima à de tabelas normalizadas de bancos relacionais. Cada entidade foi separada em um arquivo distinto, com colunas bem definidas e compatíveis com as estruturas usadas na modelagem de grafos.

Exemplos:

#### `countries.csv`

```csv
id,name,iso3
1,Brazil,BRA
2,Germany,DEU
3,United States,USA
```

Simula a tabela `countries(id, name, iso3)` em bancos relacionais.

#### `covid_cases.csv`

```csv
id,country_iso3,date,total_cases,total_deaths
1,BRA,2021-01-01,1000000,30000
2,BRA,2021-02-01,1500000,45000
3,DEU,2021-01-01,500000,10000
```

Simula a tabela `covid_cases(country_id, date, total_cases, total_deaths)` com a chave estrangeira substituída por `iso3` para facilitar ingestão direta no Neo4j via `MERGE`.

#### `vaccinations.csv`

```csv
id,country_iso3,date,people_vaccinated
1,BRA,2021-01-10,500000
2,BRA,2021-02-10,1000000
3,DEU,2021-01-10,200000
```

#### `vaccines.csv`

```csv
id,name
1,Moderna
2,Pfizer-BioNTech
3,AstraZeneca
```

#### `country_vaccines.csv`

```csv
country_iso3,vaccine_name
BRA,Moderna
BRA,AstraZeneca
USA,Pfizer-BioNTech
```

Simula uma tabela associativa (n:n) entre país e vacina.

#### `vaccine_approvals.csv`

```csv
vaccine_name,approval_date
Moderna,2020-12-18
Pfizer-BioNTech,2020-12-11
AstraZeneca,2021-01-29
```

### Comportamento Realista

- Os arquivos são gerados com nomes distintos por timestamp (`cases_20240422_1500.csv`, etc.), simulando dumps versionados.
- A periodicidade é configurável via variável de ambiente: `GENERATOR_INTERVAL_MINUTES=5`
- O volume também é configurável: `GENERATOR_COUNTRY_COUNT=3`, `GENERATOR_CASES_PER_COUNTRY=10`
- O diretório `/data` é compartilhado entre o gerador e o ETL via volume Docker.

### Integração com o ETL

A ETL foi construída para:

- Processar múltiplos arquivos por entidade.
- Validar headers e campos esperados.
- Tolerar variações simples no conteúdo (linhas duplicadas, campos nulos).
- Eliminar duplicidade via uso de `MERGE` e `CREATE CONSTRAINT` no Neo4j.

### Justificativa

Essa estrutura foi desenhada intencionalmente para manter **coerência com dumps de bases SQL normalizadas**, e facilitar evolução futura do projeto, onde os CSVs poderiam ser substituídos por:

- Conexão direta com bancos relacionais (MySQL, PostgreSQL)
- Leitura via APIs internas da empresa
- Ingestão por mensageria (ex: Kafka ou SQS)

O uso de arquivos CSV com separação clara por entidade também permite reuso de código e modularização no pacote `reader/`.




## Fluxo de Execução Geral

1. O script Python `csv-generator.py` escreve arquivos CSV no diretório `data/`, simulando dumps periódicos.
2. O processo `cmd/etl/main.go` roda continuamente, lendo os arquivos e carregando os dados no Neo4j.
3. A API (em `cmd/api/main.go`) expõe endpoints que consultam o banco para responder às perguntas do desafio.

*Nota: por limitação de tempo, não foi incluído um diagrama de execução. Pretendo incluir isso futuramente na branch de melhorias.*

## Problema: Erro ao iniciar o Neo4j - `UnsupportedLogVersionException`

Esse erro ocorre quando a pasta `neo4j/data` foi criada com uma versão diferente do Neo4j e você tenta iniciar com uma versão inferior.

**Solução:**

```bash
rm -rf ./neo4j/data
docker-compose down -v
docker-compose up --build
```

Esse projeto permite reprocessamento completo, então o banco pode ser recriado com segurança.



## Execução com Docker e Makefile

A aplicação foi estruturada para ser executada facilmente em ambiente isolado, utilizando `Docker`, `docker-compose` e um `Makefile` com comandos padronizados. Isso permite replicar o ambiente de desenvolvimento e testes com consistência, independentemente do sistema operacional.

### Pré-requisitos

- Docker
- Docker Compose v2
- Make

### Variáveis de Ambiente

As configurações são centralizadas no arquivo `.env`, incluindo:

```
NEO4J_URI=bolt://neo4j:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=testpassword
GENERATOR_INTERVAL_MINUTES=1
GENERATOR_COUNTRY_COUNT=3
GENERATOR_CASES_PER_COUNTRY=10
ETL_INTERVAL_SECONDS=61
API_PORT=8080
OUTPUT_DIR=/data

```

### Subindo o ambiente completo

Para subir todos os serviços (API, ETL, Neo4j e gerador de CSVs):

```
make up
```

Para derrubar todos os containers:

```
make down
```

### Comandos disponíveis no Makefile

| Comando         | Descrição                                             |
|-----------------|-------------------------------------------------------|
| `make up`       | Sobe todos os serviços com docker-compose             |
| `make down`     | Derruba os containers                                 |
| `make api`      | Executa somente o serviço da API localmente           |
| `make etl`      | Executa somente o serviço de ETL localmente           |
| `make test`     | Executa os testes da aplicação                        |
| `make build`    | Compila os binários da API e ETL                      |
| `make logs`     | Mostra os logs dos containers                         |

### Execução isolada

Para rodar a ETL em modo contínuo:

```
go run cmd/etl/main.go
```

Para rodar a API localmente:

```
go run cmd/api/main.go
```

### Acesso ao Neo4j

Após subir o ambiente, o Neo4j estará disponível em:

```
http://localhost:7474
```

Login padrão:
- Usuário: `neo4j`
- Senha: `testpassword`

A senha pode ser alterada ou desabilitada via variável `NEO4J_AUTH=none` no `docker-compose.yml`, se necessário.

## Endpoints da API REST

A API foi implementada utilizando o framework `Gin` em Golang. Os endpoints foram desenvolvidos para responder às perguntas propostas no desafio técnico, de forma parametrizável via query params.

### 1. Total acumulado de casos e mortes por país e data

**Endpoint:**

```
GET /cases-and-deaths
```

**Parâmetros:**

- `country`: código ISO3 do país (ex: `BRA`)
- `date`: data no formato `YYYY-MM-DD`

**Exemplo:**

```
GET /cases-and-deaths?country=BRA&date=2021-02-01
```

**Resposta:**

```json
{
  "country": "Brazil",
  "date": "2021-02-01",
  "total_cases": 1500000,
  "total_deaths": 45000
}
```

---

### 2. Total de pessoas vacinadas por país e data

**Endpoint:**

```
GET /vaccinated
```

**Parâmetros:**

- `country`: código ISO3 do país
- `date`: data no formato `YYYY-MM-DD`

**Exemplo:**

```
GET /vaccinated?country=DEU&date=2021-01-10
```

**Resposta:**

```json
{
  "country": "Germany",
  "date": "2021-01-10",
  "total_vaccinated": 200000
}
```

---

### 3. Vacinas utilizadas por país

**Endpoint:**

```
GET /vaccines-by-country
```

**Parâmetros:**

- `country`: código ISO3 do país

**Exemplo:**

```
GET /vaccines-by-country?country=USA
```

**Resposta:**

```json
{
  "country": "United States",
  "vaccines": ["Pfizer-BioNTech", "Moderna"]
}
```

---

### 4. Datas de aprovação de uma vacina

**Endpoint:**

```
GET /approval-dates
```

**Parâmetros:**

- `vaccine`: nome da vacina (exato)

**Exemplo:**

```
GET /approval-dates?vaccine=AstraZeneca
```

**Resposta:**

```json
{
  "vaccine": "AstraZeneca",
  "approval_date": "2021-01-29"
}
```

---

### 5. Países que utilizaram uma vacina específica

**Endpoint:**

```
GET /countries-by-vaccine
```

**Parâmetros:**

- `vaccine`: nome da vacina (exato)

**Exemplo:**

```
GET /countries-by-vaccine?vaccine=Pfizer-BioNTech
```

**Resposta:**

```json
{
  "vaccine": "Pfizer-BioNTech",
  "countries": ["United States", "Germany"]
}
```

---

### Observações

- Todos os parâmetros são obrigatórios para retorno válido.
- As respostas retornam status HTTP 400 em caso de parâmetros inválidos ou ausentes.
- Em produção, endpoints poderiam ser protegidos por autenticação e cache, não implementados neste desafio.

## Testes Automatizados

O projeto inclui testes unitários para os principais componentes do ETL e utilitários auxiliares. A estratégia adotada foi validar funcionalidades críticas de leitura de arquivos, transformação de dados e utilitários que fazem parsing e normalização.

### Estratégia de Testes

- Foram priorizadas as camadas `reader/` e `util/`, que contêm lógica de parsing, validação e transformação.
- Os testes verificam o comportamento com dados válidos e inválidos, garantindo robustez em situações esperadas na ingestão de arquivos CSV simulados.
- Foram utilizados mocks simples e arquivos de entrada sintéticos para isolar a lógica testada.
- Os testes podem ser executados com `go test ./...`

### Exemplos de arquivos testados

- `reader/covid_case_reader_test.go`
- `reader/country_reader_test.go`
- `reader/vaccination_reader_test.go`
- `util/date_test.go`

### Como rodar os testes

A execução dos testes pode ser feita localmente com o comando:

```
make test
```

Ou diretamente com:

```
go test ./...
```

### Possíveis melhorias

Por limitação de tempo, **não foram implementados testes para os handlers da API nem testes de integração com o banco Neo4j**. Com mais tempo, seria possível:

- Criar testes de integração para os repositórios, validando Cypher executado.
- Mockar o driver do Neo4j nos testes da camada `usecase/`.
- Testar os handlers da API com `httptest.NewRecorder()` e chamadas simuladas.

A estrutura do projeto já está preparada para suportar testes mais avançados com injeção de dependência e interfaces desacopladas.



## Decisões Técnicas Justificadas

Esta seção descreve as principais decisões técnicas tomadas durante o desenvolvimento da solução, com justificativas baseadas em boas práticas, entrevistas e materiais de referência como o livro *Understanding ETL (O’Reilly, 2024)*.

### 1. Estrutura modular da aplicação

O projeto foi organizado de forma modular, com separação clara entre:

- `handler`: camada de entrada HTTP.
- `repository`: camada de acesso ao banco (Neo4j).
- `reader`: responsável pela leitura e parsing dos arquivos CSV.
- `loader`: lógica de transformação e escrita no banco.
- `util`: funções auxiliares para parse e validações.

Essa estrutura facilita manutenção e evolução do sistema, mesmo sem aplicar todos os elementos da Clean Architecture (como `usecase`, não utilizado aqui por simplicidade e escopo).

### 2. Uso do modelo Medallion (Bronze, Silver, Gold)

O pipeline ETL foi dividido em camadas:

- **Bronze**: dados brutos gerados periodicamente pelo `csv-generator`, simulando dumps SQL.
- **Silver**: parsing e validação dos dados em memória.
- **Gold**: escrita no Neo4j com relacionamento entre entidades, garantindo integridade.

A divisão entre camadas ajuda a organizar responsabilidades e permite rastrear erros com mais precisão.

### 3. Identificadores lógicos no grafo

A modelagem no Neo4j foi feita com identificadores de negócio, ao invés de IDs técnicos dos CSVs:

- País → `iso3`
- Vacina → `name`
- Casos → `date + iso3`
- Estatísticas de vacinação → `date + iso3`
- Aprovação → `vaccine_name`

Essa abordagem torna as consultas mais naturais e evita colisões ou duplicidade de dados.

### 4. Constraints no Neo4j

Foram criadas constraints para garantir unicidade de nós no grafo, permitindo ingestão idempotente:

```cypher
CREATE CONSTRAINT country_iso3_unique IF NOT EXISTS FOR (c:Country) REQUIRE c.iso3 IS UNIQUE
```

O uso de `MERGE` em conjunto com essas constraints garante que execuções repetidas do ETL não causem duplicações.

### 5. Tratamento de duplicidades nos CSVs

Na prática, arquivos CSV exportados de bancos podem conter registros repetidos. O ETL foi projetado para ser **tolerante a esse tipo de duplicidade**. A deduplicação ocorre por dois mecanismos:

- Durante o parsing, linhas duplicadas são ignoradas na etapa de leitura.
- No momento da escrita no Neo4j, o uso de `MERGE` e as `constraints` garantem que os nós e relacionamentos não sejam criados em duplicidade.

Isso torna a ingestão confiável mesmo em cenários onde os dados de entrada não são perfeitamente limpos.

### 6. Reuso da conexão com Neo4j

A lógica de conexão foi isolada no arquivo:

```
internal/repository/neo4j/client.go
```

Esse client é compartilhado entre a API e a ETL, seguindo o princípio de reuso e facilitando alterações futuras (ex: troca de banco ou configuração centralizada).

### 7. Simulação de dados com `csv-generator`

Foi desenvolvido um script em Python que gera arquivos CSV continuamente, simulando dumps de bancos relacionais. A escolha por múltiplos arquivos, gerados com timestamps distintos, ajuda a testar a ingestão incremental.

- Intervalo configurável: `GENERATOR_INTERVAL_MINUTES`
- Volume configurável: `GENERATOR_COUNTRY_COUNT`, `GENERATOR_CASES_PER_COUNTRY`
- Os arquivos são salvos em `./data`, volume compartilhado com a ETL.

### 8. Projeto preparado para evoluções

Mesmo com escopo controlado, a arquitetura permite evoluções futuras como:

- Substituição dos CSVs por conexões com bancos, APIs ou filas.
- Testes de integração com banco e rotas HTTP.
- Cache, autenticação e middlewares.
- Paralelização da ingestão com goroutines e workers (não implementado por falta de tempo).

## Limitações e Melhorias Futuras

Por conta de restrições de tempo relacionadas à minha carga de trabalho atual, algumas funcionalidades importantes não foram implementadas dentro da janela de entrega. No entanto, pretendo continuar evoluindo esta solução em uma branch separada, fora do escopo da avaliação inicial, com foco em torná-la mais robusta, escalável e próxima de um ambiente de produção real.

Abaixo estão listadas as melhorias previstas, com detalhamento técnico de como seriam implementadas.

### 1. Paralelismo com Workers na ETL

Atualmente, o processamento dos arquivos CSV é sequencial. Para ambientes com grande volume de dados, a escalabilidade pode ser significativamente melhorada utilizando paralelismo com goroutines e channels.

**Como seria implementado:**

- Criar um `channel` para envio de caminhos de arquivos CSV ou blocos de linhas.
- Iniciar múltiplos workers (`goroutines`) lendo do channel, cada um processando um arquivo ou lote de dados.
- Usar `sync.WaitGroup` para controlar a finalização dos workers.
- Garantir controle de concorrência na escrita no banco, respeitando limites de conexões do driver do Neo4j.

**Exemplo simplificado:**

```go
files := getCSVFiles()
lineChan := make(chan string)
wg := sync.WaitGroup{}

for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for line := range lineChan {
            _ = processLine(line)
        }
    }()
}

for _, file := range files {
    for _, line := range readLines(file) {
        lineChan <- line
    }
}
close(lineChan)
wg.Wait()
```

Essa abordagem permite escalar horizontalmente e usar melhor os recursos da máquina.

---

### 2. Otimizações com Workers e Cache na API

Em cenários com grande volume de chamadas simultâneas, a API pode se beneficiar de melhorias como:

- **Cache em memória** para queries pouco voláteis (ex: datas de aprovação de vacinas).
- **Fila interna de processamento** para chamadas intensivas que podem ser processadas em background.
- **Worker pools** para lidar com eventos concorrentes de forma previsível.

**Como seria feito:**

- Implementação de cache usando bibliotecas como `go-cache` ou `groupcache`.
- Uso de `middleware` para invalidação de cache condicional.
- Separação de handlers que exigem baixa latência dos que podem ser enfileirados (via `channel` + `worker`).

---

### 3. Tratamento mais robusto de erros

A versão atual da ETL apenas loga erros simples, mas não categoriza nem realiza retry. Em ambiente produtivo, é necessário:

**Como seria feito:**

- Uso de estrutura de erro enriquecida (`type EnrichedError struct { Op string; Err error; Line string; File string }`)
- Implementação de retry automático com política de backoff exponencial (ex: `time.Sleep(backoff)`).
- Envio de erros críticos para uma `dead-letter queue` (ex: arquivo `.errlog`, banco auxiliar ou tópico Kafka).
- Monitoramento com alertas para falhas repetidas (ex: integração futura com Prometheus/Grafana).

---

### 4. Testes com alto volume de dados

Os testes atuais utilizam arquivos pequenos e sintéticos. Em produção, seria necessário:

**Como seria feito:**

- Geração automática de milhares de linhas para cada CSV com um script Python ou Go.
- Simulação de ingestão em carga com medição de tempo e uso de memória.
- Análise de performance de ingestão por núcleo (CPU-bound vs IO-bound).
- Benchmark com `go test -bench` e medição do throughput de escrita no Neo4j.

---

### 5. Testes de integração e cobertura total da aplicação

**Como seria feito:**

- Implementação de testes da camada `repository` com banco Neo4j em container isolado.
- Testes da API com `httptest.NewRecorder`, validando rotas, parâmetros e status code.
- Uso de bibliotecas como `testcontainers-go` para simular ambiente real em CI.
- Criação de fixtures e seeds para cenários específicos de teste.
- Pipeline CI com `make test` e `go test ./...` executando a cada push.

---

### Considerações finais

Todas essas melhorias são tecnicamente viáveis e fazem parte do roadmap para evolução da aplicação. Embora não tenham sido incluídas nesta versão por limitações de agenda profissional, pretendo continuar a evolução do projeto em uma branch separada (`improvements`) após o prazo de entrega da avaliação técnica, com objetivo de transformá-lo em uma base reutilizável para projetos reais de ingestão e análise de dados em grafos.


