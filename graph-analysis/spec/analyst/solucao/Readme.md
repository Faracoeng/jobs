# Projeto: Graph Analysis - ETL + API + Neo4j

## 1. Introdução

Este projeto integra um processo de ETL com uma API REST, utilizando o banco de dados de grafos Neo4j. O objetivo é fornecer uma estrutura para análise de dados relacionados à COVID-19, incluindo casos, vacinação e aprovações de vacinas.

## 2. Arquitetura Geral

O projeto adota a arquitetura hexagonal. Essa abordagem promove a separação de preocupações, facilitando a manutenção e a escalabilidade do sistema.

- **Domínio**: Contém as entidades e interfaces que representam as regras de negócio.
- **Aplicação**: Implementa os casos de uso do sistema.
- **Infraestrutura**: Fornece implementações concretas para interfaces, como repositórios e serviços externos.
- **Interfaces de Entrada/Saída**: Incluem a API REST e o processo de ETL.

## 3. ETL

### 3.1. Descrição Geral

O processo de ETL (Extract, Transform, Load) é responsável por extrair dados de arquivos CSV (ou outras fontes futuras), transformá-los em entidades de domínio, e carregá-los no banco de dados Neo4j.  
Foi desenhado para rodar continuamente e com capacidade de adaptação a novas fontes de dados.

### 3.2. Fluxo de Execução

1. **Leitura dos dados**
   - Leitura desacoplada via interface `DataSource`.
   - Atualmente implementado com `CSVSource`, mas preparado para trocar para APIs, bancos de dados, etc.
   - Opções de leitura:
     - `ReadAll`: lê todo o conteúdo de uma vez.
     - `StreamData`: envia linha a linha via canais para processamentos de fluxo.

2. **Transformação**
   - Cada linha é convertida em uma entidade do domínio (`Country`, `CovidCase`, `VaccinationStat`, `Vaccine`, `CountryVaccine`, `VaccineApproval`).
   - Datas são parseadas de forma robusta com a função `ParseDate`, aceitando diferentes formatos.

3. **Carga**
   - Utiliza transações em lote (`batch`) para inserir no Neo4j.
   - Utiliza o comando `MERGE` para garantir idempotência e evitar duplicações.

4. **Controle de Ciclo**
   - O ETL aguarda o surgimento do arquivo `ready.flag` antes de processar.
   - Após o processamento, o `ready.flag` é removido para evitar reprocessamentos.
   - O ciclo de execução é repetido conforme `ETL_INTERVAL_SECONDS`.

### 3.3. Estratégias de Eficiência

- **Arquitetura de fonte genérica (`DataSource`)**
  - Permite que o ETL troque a origem dos dados sem impactar o core da aplicação.

- **Batch e transações controladas**
  - Dados são agrupados em lotes de tamanho configurável (`ETL_BATCH_SIZE`), otimizando performance no Neo4j.

- **Stream de dados**
  - O módulo `StreamData` possibilita processamento linha a linha sob demanda, facilitando futuras implementações com alto volume ou pipelines assíncronos.

- **Separação de loaders**
  - Cada tipo de dado tem um loader independente, favorecendo paralelismo e manutenção isolada.

- **Criação automática de constraints**
  - Constraints de unicidade são criadas no Neo4j antes da carga, garantindo integridade.

- **Leitura deduplicada**
  - A função `ReadCSVFile` remove duplicações antes do processamento, reduzindo conflitos de chave.

### 3.4. Decisões Técnicas Importantes

- **Interface `DataSource`**
  - Criada para desacoplar a leitura dos dados da origem (ex: CSV → API → Banco → etc).
  - Atual implementação: `CSVSource`, que lê via utilitário `ReadCSVFile`.

- **Uso de `MERGE` no Neo4j**
  - Escolhido para suportar reprocessamentos sem gerar inconsistências ou dados duplicados.

- **Separação de nós e relacionamentos**
  - Primeiramente são inseridos os nós (`Country`, `Vaccine`, `CovidCase`, etc), e depois os relacionamentos (`HAS_CASE`, `USES`, `APPROVED_ON`, `VACCINATED_ON`).

- **Uso de canais em `StreamData`**
  - Permite futuras melhorias como consumo paralelo, throttling, backpressure ou processamento progressivo em fluxos muito grandes.

- **Remoção de redundâncias**
  - Exemplo: campo redundante `vaccine` foi removido da entidade `VaccineApproval`, deixando a relação explícita somente no grafo.

- **Pronto para mudanças futuras**
  - O design modularizado permite trocar facilmente a fonte dos dados ou escalar o paralelismo do ETL sem reescrever o core.

### 3.5. Fluxo de Execução do ETL

``` text
1. Gerador de CSVs:
   - Cria os arquivos de dados (`countries.csv`, `covid_cases.csv`, etc.).
   - Gera o arquivo `ready.flag` sinalizando que os dados estão prontos.

2. Detecção do arquivo ready.flag:
   - O ETL monitora continuamente a existência do `ready.flag` no diretório de dados.

3. Leitura da fonte de dados:
   - Cada arquivo CSV é lido utilizando abstrações de DataSource (`CSVSource`).
   - Leitura em memória ou via streaming de dados linha a linha.

4. Transformação dos dados:
   - Cada linha lida é convertida em uma entidade do domínio (`Country`, `CovidCase`, etc.).
   - Aplicação de validações e formatações (ex: parse de datas).

5. Processamento por tipo de dado:
   - Carregamento dos dados ocorre de forma modular:
     - Países → Vacinas → Casos de COVID → Relações país-vacina → Vacinação → Aprovações.

6. Escrita em lote no banco Neo4j:
   - Inserções são agrupadas em batches (`ETL_BATCH_SIZE`).
   - Uso de transações otimizadas e comandos `MERGE` para garantir idempotência.

7. Conclusão do ciclo:
   - O arquivo `ready.flag` é removido.
   - O ETL retorna ao estado de monitoramento, aguardando nova geração de dados.

```



### 3.6. Variáveis de Ambiente

| Variável             | Descrição                          | Exemplo             |
|----------------------|------------------------------------|---------------------|
| OUTPUT_DIR           | Diretório dos arquivos CSV         | ./data              |
| ETL_INTERVAL_SECONDS | Intervalo entre execuções do ETL   | 300                 |
| ETL_BATCH_SIZE       | Tamanho do batch de inserção       | 500                 |

---


## 4. API

### 4.1. Descrição Geral

A API REST fornece endpoints para consulta de dados relacionados à COVID-19, permitindo acesso a informações respondendo aos questionamentos especificados para o case técnico:

1. Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma
data determinada?
2. Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em
uma data específica?
3. Quais vacinas foram usadas em um país específico?
4. Em quais datas as vacinas foram autorizadas para uso?
5. Quais países usaram uma vacina específica

### 4.2. Estrutura

A estrutura da API foi projetada seguindo os princípios da Arquitetura Hexagonal (Ports and Adapters), respeitando conceitos de Clean Architecture e SOLID, com foco em desacoplamento e escalabilidade.

- **Entities (Domínio)**
  - Representam os conceitos fundamentais do sistema: `Country`, `CovidCase`, `VaccinationStat`, `Vaccine`, `VaccineApproval`.
  - Define também as interfaces de repositórios (`CovidRepository`, `VaccinationRepository`, `VaccineRepository`), que descrevem o contrato da persistência sem se acoplar a nenhuma tecnologia.

- **UseCases (Aplicação)**
  - Implementam a lógica de orquestração da aplicação.
  - Cada caso de uso é responsável por receber dados de entrada (via DTOs), chamar o repositório adequado e retornar dados de saída também formatados como DTOs.
  - O UseCase depende apenas da abstração (`interface`) dos repositórios, permitindo a inversão de dependência.

  Exemplo de fluxo em um UseCase:
  ```
  InputDTO → UseCase → Chamada de interface de repositório → OutputDTO
  ```

- **Repositories (Infraestrutura)**
  - Implementam as interfaces definidas no domínio.
  - São responsáveis por interagir diretamente com o banco de dados Neo4j.
  - Isolam detalhes técnicos como consultas Cypher, permitindo que a aplicação não dependa da tecnologia usada.

- **Handlers (Interface Externa - API HTTP)**
  - São responsáveis por receber as requisições HTTP, extrair os parâmetros e chamar os UseCases.
  - Os Handlers não conhecem regras de negócio, entidades ou detalhes de banco de dados.
  - Eles apenas tratam a entrada da requisição, invocam o UseCase e devolvem a resposta adequada.

- **DTOs (Data Transfer Objects)**
  - São utilizados para definir os formatos de entrada e saída dos UseCases.
  - Permitem desacoplar a API das estruturas internas do domínio, garantindo que mudanças internas não afetem diretamente a comunicação externa.

### 4.3. Padrão de Execução da API

O fluxo completo de processamento de uma requisição HTTP na API segue a seguinte sequência controlada:

```text
1. Recepção da Requisição:
   - O Handler recebe a requisição e extrai os parâmetros de entrada.

2. Montagem do InputDTO:
   - Os dados extraídos são organizados em um objeto de transferência de dados (DTO).

3. Chamada ao UseCase:
   - O UseCase recebe o InputDTO e executa a lógica de negócio necessária.

4. Acesso ao Repositório:
   - O UseCase interage com a interface do repositório para consultar ou persistir informações no banco.

5. Montagem do OutputDTO:
   - O resultado da operação é encapsulado em um OutputDTO, isolando a estrutura interna do domínio.

6. Resposta HTTP:
   - O Handler retorna o OutputDTO convertido em resposta JSON para o cliente.
```

### 4.4. Justificativas da Arquitetura

- **Desacoplamento total entre camadas**
  - Handlers não sabem da estrutura do domínio.
  - UseCases não sabem quem implementa os repositórios.
  - Repositórios são definidos por interfaces e injetados no momento da inicialização.

- **Facilidade para testes**
  - Como tudo depende de abstrações, é possível mockar repositórios facilmente nos testes de UseCase.

- **Facilidade para trocar implementações**
  - Se amanhã for necessário trocar o Neo4j por outro banco de grafos, apenas a implementação do repositório precisaria mudar — o restante do sistema continuaria intacto.

- **Escalabilidade de código**
  - Fica fácil adicionar novos endpoints, novos casos de uso, ou novas fontes de dados (por exemplo, Kafka, filas, ou API externa) sem afetar a estrutura central da aplicação.

- **Organização de diretórios**
  - Separação em `internal/entity`, `internal/usecase`, `internal/infra/repository/neo4j`, e `internal/api/handler` reflete as fronteiras da arquitetura hexagonal.

---


### 4.5. Endpoints

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

### 4.6. Variáveis de Ambiente

| Variável       | Descrição     | Exemplo             |
|----------------|---------------|---------------------|
| API_HOST       | IP da API     | 0.0.0.0             |
| API_PORT       | Porta da API  | 8080                |
| NEO4J_URI      | URI do Neo4j  | bolt://neo4j:7687   |
| NEO4J_USERNAME | Usuário Neo4j | neo4j               |
| NEO4J_PASSWORD | Senha Neo4j   | testpassword        |

## 5. Banco de Dados (Neo4j)

### 5.1. Modelagem de Dados

- **Nós**: Representam entidades como países, vacinas e aprovações.
- **Relacionamentos**: Conectam os nós para indicar associações, como um país que utilizou uma vacina.

### 5.2. Estratégias de Persistência

- **MERGE**: Utilizado para garantir que nós e relacionamentos não sejam duplicados.
- **Constraints**: Definidos para assegurar a unicidade de determinados atributos.

## 6. Decisões Técnicas

### 6.1. Geração de CSVs baseada em tabelas relacionais

Optou-se por gerar os arquivos CSV simulando tabelas relacionais para facilitar a transformação dos dados em entidades do domínio e posterior inserção no Neo4j.

Essa estratégia visa:

- Permitir uma transição natural entre modelo relacional e modelo de grafos.
- Garantir clareza nos relacionamentos entre entidades.
- Facilitar a leitura dos dados para futuros loaders de grafos ou bancos relacionais.

#### 6.1.1. Explicação dos Relacionamentos

| Tabela Origem | Tabela Destino | Tipo de Relacionamento | Explicação |
|---------------|----------------|------------------------|------------|
| country       | covid_case      | 1:N                    | Um país pode ter muitos registros de casos de COVID (um por data). |
| country       | vaccination_stats | 1:N                  | Um país pode ter muitos registros de vacinação (um por data). |
| country       | vaccine (via country_vaccine) | N:N         | Um país pode usar várias vacinas, e uma vacina pode ser usada por vários países. |
| vaccine       | vaccine_approval | 1:1                  | Cada vacina tem uma única data de aprovação. |

#### 6.1.2. Resumo Simplificado

- `country` 1:N `covid_case`
- `country` 1:N `vaccination_stats`
- `country` N:N `vaccine` (através da tabela de junção `country_vaccine`)
- `vaccine` 1:1 `vaccine_approval`

#### 6.1.3. Controle da Ordem de Leitura e Carga

Para manter a coerência lógica dos dados ao construir o grafo no Neo4j, a ordem de ingestão dos CSVs foi definida da seguinte forma:

1. **Carregar Países (`countries.csv`)**
   - Os nós de países precisam existir antes de criar qualquer relacionamento ou evento associado a eles.

2. **Carregar Vacinas (`vaccines.csv`)**
   - As vacinas precisam estar presentes antes de associá-las aos países ou registrar aprovações.

3. **Carregar Casos de COVID (`covid_cases.csv`) e Relações País → Vacinas (`country_vaccines.csv`)**
   - Os registros de casos e as relações de uso de vacinas precisam dos nós de países e vacinas já presentes.

4. **Carregar Estatísticas de Vacinação (`vaccinations.csv`) e Aprovações de Vacinas (`vaccine_approvals.csv`)**
   - Estatísticas e aprovações são inseridas após os relacionamentos principais já estarem estabelecidos.

Essa ordem garante que:

- Nenhuma tentativa de criar uma relação ocorra sem os nós necessários já criados.
- O grafo final mantenha consistência sem dependências faltantes.
- O processo de carga seja robusto mesmo em execuções parciais ou reinícios.

Além disso, como cada tipo de entidade tem loaders separados, é possível **ler diferentes arquivos em paralelo**, respeitando a ordem de dependências lógicas.

---

### 6.2. Remoção do campo `vaccine` da entidade `VaccineApproval`

Identificou-se que o campo `vaccine` era redundante, pois a associação entre vacinas e aprovações já é representada por relacionamentos no grafo. Sua remoção simplificou a estrutura e evitou inconsistências.

### 6.3. Retorno de objetos completos ou dados primitivos nas interfaces de repositório

A decisão sobre o tipo de retorno das interfaces de repositório foi baseada na complexidade dos dados:

- **Objetos completos**: Utilizados quando múltiplos atributos relacionados são retornados.
- **Dados primitivos**: Utilizados para retornos simples, como contagens ou valores únicos.

### 6.4. Tratamento de dados repetidos com chaves primárias

Ao identificar registros com chaves primárias duplicadas, optou-se por atualizar os dados existentes com as novas informações. Essa abordagem garante que o banco reflita os dados mais recentes.



## 7. Instruções de Uso

### 7.1. Pré-requisitos

- Docker e Docker Compose instalados
- Make instalado (opcional, mas recomendado)

### 7.2. Subir todos os serviços

Para compilar e subir todos os serviços (Neo4j, CSV Generator, ETL, e API), basta rodar:

```
make run
```

Esse comando executa:

- Testes de unidade
- Build das imagens Docker (API, ETL, CSV Generator)
- Inicialização dos containers

> Caso prefira executar manualmente:
> 
> ```
> make build-all
> make up
> ```

### 7.3. Parar todos os serviços

Para parar e remover os containers:

```
make down
```

### 7.4. Visualizar logs dos containers

Para acompanhar os logs em tempo real:

```
make logs
```

### 7.5. Reiniciar todo o ambiente

Para reiniciar os containers:

```
make restart
```

### 7.6. Testar cobertura dos testes

Para rodar testes com análise de cobertura:

```
make test-cover
```

### 7.7. Acessar os serviços

- **Neo4j Browser**: `http://localhost:7474`
- **API REST**: `http://localhost:8080`
- **Documentação Swagger**: `http://localhost:8080/swagger/index.html`

### 7.8. Variáveis de Ambiente

Todas as variáveis são configuradas no arquivo `.env`:

| Variável             | Descrição                 | Exemplo             |
|----------------------|---------------------------|---------------------|
| NEO4J_URI             | URI do Neo4j               | bolt://neo4j:7687   |
| NEO4J_USERNAME        | Usuário do Neo4j           | neo4j               |
| NEO4J_PASSWORD        | Senha do Neo4j             | testpassword        |
| OUTPUT_DIR            | Diretório de saída dos CSVs | ./data              |
| ETL_INTERVAL_SECONDS  | Intervalo entre execuções do ETL | 300         |
| ETL_BATCH_SIZE        | Tamanho dos batches de carga | 500               |
| API_HOST              | IP para a API ouvir         | 0.0.0.0             |
| API_PORT              | Porta da API               | 8080                |

---



## 8. Melhorias Futuras

- Adicionar retentativas automáticas para batches com falhas no ETL.
- Integrar monitoramento e alertas para falhas no ETL e na API.
- Melhorar o sistema de logs para registrar operações críticas e falhas com mais detalhes.
- Aumentar a consistência e profundidade dos dados armazenados no Neo4j para suportar novos casos analíticos.
- Implementar testes unitários e mocks na camada de UseCases da API, garantindo cobertura de regras de negócio.
- Implementar estágios intermediários no ETL seguindo a arquitetura Medallion:
  - **Bronze**: armazenamento dos dados brutos recebidos (sem limpeza).
  - **Silver**: armazenamento dos dados validados, deduplicados e transformados.
  - **Gold**: dados agregados e preparados para análises de negócio.
  - Atualmente os dados são manipulados apenas em memória durante o processamento. A criação de CSVs intermediários nas camadas bronze/silver permitiria maior rastreabilidade, recuperação em caso de falhas e preparação para cenários de alta escala.

---

## 9. Gerador de CSVs

### 9.1. Descrição

O projeto inclui um gerador automático de arquivos CSV para simular dados relacionados à COVID-19, vacinação e aprovação de vacinas.  
Este gerador permite a alimentação contínua do processo de ETL, garantindo que haja novos dados disponíveis periodicamente.

### 9.2. Funcionamento

- Gera os seguintes arquivos no diretório de saída (`OUTPUT_DIR`):
  - `countries.csv`
  - `covid_cases.csv`
  - `vaccinations.csv`
  - `vaccines.csv`
  - `country_vaccines.csv`
  - `vaccine_approvals.csv`
- Cria o arquivo `ready.flag` para sinalizar que os arquivos estão prontos para processamento.
- Reescreve os arquivos a cada intervalo de tempo configurável.

Os dados são gerados de forma aleatória, mas respeitando uma estrutura consistente para simular cenários reais.

### 9.3. Estratégias Aplicadas

- **Escrita segura**: os arquivos são primeiro gravados como `.tmp` e depois renomeados, evitando arquivos corrompidos.
- **Controle de disponibilidade**: a criação de `ready.flag` permite que o ETL saiba quando os dados estão prontos.
- **Configurações flexíveis**: o número de países, intervalo de geração e diretório de saída podem ser controlados via variáveis de ambiente.

### 9.4. Variáveis de Ambiente

| Variável                      | Descrição                                 | Exemplo |
|--------------------------------|-------------------------------------------|---------|
| OUTPUT_DIR                     | Diretório onde os CSVs serão gerados      | ./data  |
| GENERATOR_COUNTRY_COUNT        | Número de países a serem gerados          | 30      |
| GENERATOR_INTERVAL_MINUTES     | Intervalo entre gerações de CSVs (minutos) | 5      |

### 9.5. Observações

- Todos os países gerados possuem dados contínuos desde 01/01/2021.
- Cada país pode estar associado a 1 até 3 vacinas aleatórias.
- A cada nova geração, os arquivos anteriores são substituídos.
- O gerador depende apenas das bibliotecas padrão do Python e é agendado usando `schedule`.
- O serviço `csv-generator` no `docker-compose.yml` é iniciado junto com o ambiente e depende do serviço `neo4j` estar saudável.

---
## 10. Desempenho

### 10.1. API

#### Pontos fortes

- **Alta eficiência**: A API utiliza o framework Gin Gonic, conhecido pela sua alta performance para requisições HTTP.
- **Baixo acoplamento**: A arquitetura hexagonal permite que Handlers, UseCases e Repositories operem de forma independente, otimizando o tempo de resposta.
- **Consultas otimizadas**: As queries no Neo4j são específicas e utilizam constraints (`UNIQUE`) para buscas rápidas.
- **Escalabilidade**: A API é stateless, podendo ser replicada horizontalmente com balanceamento de carga.

#### Limitações atuais

- **Conexão padrão com Neo4j**: O pool de conexões ainda não foi customizado para cenários de alta concorrência.
- **Ausência de cache**: Todas as requisições resultam em consultas diretas ao banco.
- **Sem paginação**: Endpoints podem potencialmente retornar muitos resultados sem limites, aumentando a carga.

#### Melhorias futuras

- Implementar configuração customizada de pool de conexões para Neo4j.
- Introduzir cache de respostas para consultas mais frequentes.
- Implementar paginação e limites de resposta para grandes volumes de dados.

---

### 10.2. ETL

#### Pontos fortes

- **Inserções em batch**: Dados são agrupados em lotes (`ETL_BATCH_SIZE`) para reduzir overhead de transações no Neo4j.
- **Idempotência com MERGE**: Permite reprocessamentos seguros, sem risco de dados duplicados.
- **Separação modular**: Cada tipo de entidade (`Country`, `CovidCase`, `VaccinationStat`, etc.) possui seu loader específico, facilitando manutenção e evolução.
- **Arquitetura preparada para streaming**: A estrutura `StreamData` permite processamentos linha a linha no futuro.

#### Limitações atuais

- **Processamento completo a cada ciclo**: Todo o conjunto de dados é lido e reprocessado a cada execução do ETL.
- **Ausência de carga incremental**: O ETL ainda não identifica e carrega apenas registros novos.
- **Execução sequencial**: Os diferentes loaders operam de forma sequencial, sem paralelismo entre tipos de dados.

#### Melhorias futuras

- Implementar carga incremental (delta ingestion) para otimizar tempo e consumo de recursos.
- Explorar paralelismo entre loaders para acelerar o tempo total de ingestão.
- Aproveitar o streaming linha a linha (`StreamData`) para processar grandes volumes sem sobrecarregar a memória.

---
#### Testes de Carga Realizados

Foram realizados testes de ingestão configurando:

- `GENERATOR_COUNTRY_COUNT=10`
- `GENERATOR_CASES_PER_COUNTRY=10000`

Esse cenário gerou aproximadamente **100.000 registros de casos de COVID** e **100.000 registros de vacinações** para simular um volume considerável de dados.  
O ETL se comportou de forma estável neste volume, validando a estratégia de inserções em batch e o uso de transações controladas para escrita no Neo4j.
