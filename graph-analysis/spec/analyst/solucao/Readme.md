# README

## ETL - Ingestão de Dados no Neo4j

### Objetivo
O ETL lê periodicamente arquivos CSV que contêm informações de países, casos de Covid-19, vacinação e vacinas, e carrega os dados no banco de grafos Neo4j.

A execução é controlada pela presença do arquivo `ready.flag`, indicando que os CSVs estão prontos para leitura.

### Fluxo de execução

1. Verifica se o arquivo `ready.flag` existe.
2. Lê os seguintes arquivos CSV:
    - countries.csv
    - vaccines.csv
    - covid_cases.csv
    - vaccinations.csv
    - vaccine_approvals.csv
    - country_vaccines.csv
3. Carrega os nós no Neo4j:
    - Country
    - CovidCase
    - VaccinationStats
    - Vaccine
    - VaccineApproval
4. Cria os relacionamentos:
    - (:Country)-[:HAS_CASE]->(:CovidCase)
    - (:Country)-[:VACCINATED_ON]->(:VaccinationStats)
    - (:Country)-[:USES]->(:Vaccine)
    - (:Vaccine)-[:APPROVED_ON]->(:VaccineApproval)
5. Após a carga, o arquivo `ready.flag` é apagado para aguardar novos dados.

### Componentes Técnicos

- Linguagem: Go (Golang)
- Banco de Dados: Neo4j
- Controle de ciclo: ETL_INTERVAL_SECONDS
- Batch de inserções: ETL_BATCH_SIZE
- Sincronização: ready.flag
- Leitura: Streaming de CSV (linha por linha)

### Decisões Técnicas

- MERGE para garantir idempotência.
- Controle de sincronização via ready.flag.
- Operações em batch para reduzir sobrecarga de transações.
- Parsing de linhas com tratamento de erros sem parar o processo.

## API - Interface de Consulta aos Dados

### Estrutura

- Framework: Gin Gonic
- Divisão em handlers, usecases e repositories
- Documentação Swagger disponível em /swagger/index.html

### Endpoints

#### 1. Total acumulado de casos e mortes

- Método: GET
- URL: `/covid-stats?iso3={ISO3}&date={YYYY-MM-DD}`

Exemplo:
```
GET /covid-stats?iso3=BRA&date=2021-01-10
```

Resposta:
```json
{
  "iso3": "BRA",
  "date": "2021-01-10",
  "total_cases": 55726,
  "total_deaths": 1449
}
```

#### 2. Total de vacinados

- Método: GET
- URL: `/vaccination?iso3={ISO3}&date={YYYY-MM-DD}`

Exemplo:
```
GET /vaccination?iso3=BRA&date=2021-01-13
```

Resposta:
```json
{
  "iso3": "BRA",
  "date": "2021-01-13",
  "total_vaccinated": 839584
}
```

#### 3. Vacinas usadas em um país

- Método: GET
- URL: `/vaccines?iso3={ISO3}`

Exemplo:
```
GET /vaccines?iso3=DEU
```

Resposta:
```json
{
  "iso3": "DEU",
  "vaccines": [
    "Pfizer-BioNTech",
    "Sinovac",
    "Sputnik V"
  ]
}
```

#### 4. Datas de aprovação de vacinas

- Método: GET
- URL: `/approval-dates?name={VaccineName}`

Exemplo:
```
GET /approval-dates?name=Sinovac
```

Resposta:
```json
{
  "vaccine": "Sinovac",
  "approval_dates": [
    "2020-11-08"
  ]
}
```

#### 5. Países que usaram uma vacina

- Método: GET
- URL: `/countries-by-vaccine?name={VaccineName}`

Exemplo:
```
GET /countries-by-vaccine?name=Pfizer-BioNTech
```

Resposta:
```json
{
  "vaccine": "Pfizer-BioNTech",
  "countries": [
    "BRA",
    "DEU"
  ]
}
```

## Variáveis de Ambiente

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| NEO4J_URI | URI do Neo4j | bolt://neo4j:7687 |
| NEO4J_USERNAME | Usuário Neo4j | neo4j |
| NEO4J_PASSWORD | Senha Neo4j | testpassword |
| OUTPUT_DIR | Diretório dos CSVs | ./data |
| ETL_INTERVAL_SECONDS | Intervalo entre ETLs | 300 |
| ETL_BATCH_SIZE | Tamanho do batch | 500 |
| API_HOST | IP da API | 0.0.0.0 |
| API_PORT | Porta da API | 8080 |


