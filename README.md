# Restock Priority Service

Microserviço em Go para gestão de peças de autopeças e cálculo automático
de priorização de reposição de estoque, considerando estoque atual,
criticidade, padrão de vendas e lead time de fornecedores.

## Sumário

- [Arquitetura](#arquitetura)
- [Decisões de projeto](#decisões-de-projeto)
- [Como rodar localmente](#como-rodar-localmente)
- [Endpoints e exemplos de requisição](#endpoints-e-exemplos-de-requisição)
- [Regras de negócio](#regras-de-negócio)
- [Testes](#testes)
- [Trocando o banco de dados no futuro](#trocando-o-banco-de-dados-no-futuro)

## Arquitetura

O projeto segue uma arquitetura em camadas, com a direção de dependência
sempre apontando para dentro (HTTP → Service → Repository/Domain), nunca o
contrário:

```
cmd/
└── api/
    └── main.go

internal/
├── domain/
├── repository/
├── service/
└── http/
    ├── handler/
    └── router/
```

## Decisões de projeto

| Decisão                                                                                                | Motivo                                                                                                                                                                                  |
| ------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Go puro + Chi**                                                                                      | API pequena, leve e sem dependências desnecessárias                                                                                                                                     |
| **Persistência in-memory** (`map` + `sync.RWMutex`)                                                    | Isola a camada de armazenamento atrás de uma interface, sem exigir banco externo para rodar/avaliar o desafio                                                                           |
| **IDs gerados com `google/uuid`**                                                                      | IDs únicos sem depender do banco                                                                                                                                                        |
| **`PartService` e `PriorityService` separados**                                                        | Responsabilidades diferentes (CRUD vs. cálculo de priorização)                                                                                                                          |
| **Arredondamento de 2 casas decimais** (`expectedConsumption`, `projectedStock`, `urgencyScore`)       | Evita ruído de ponto flutuante na exibição (ex: `44.99999999999999`). Aplicado apenas no valor de saída — `needsRestock` e `urgencyScore` usam os valores não arredondados internamente |
| **Campos adicionais no JSON de resposta** (`category` na priorização; `createdAt`/`updatedAt` no CRUD) | O desafio não restringe a resposta a exatamente os campos do exemplo; esses campos agregam contexto sem remover nenhum dos exemplificados                                               |

## Como rodar localmente

Pré-requisitos: [Go 1.22+](https://go.dev/dl/) instalado.

```bash
# 1. Entrar na pasta do projeto
cd backend-test

# 2. Baixar as dependências e gerar/validar o go.sum
go mod tidy

# 3. Rodar os testes
go test ./... -v

# 4. Rodar o servidor (porta padrão 8080, configurável via env PORT)
go run ./cmd/api
```

O servidor inicia em `http://localhost:8080`.

## Endpoints e exemplos de requisição

### Health check

```bash
curl http://localhost:8080/health
```

### Criar peça

```bash
curl -X POST http://localhost:8080/parts/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Filtro de Óleo X",
    "category": "engine",
    "currentStock": 15,
    "minimumStock": 20,
    "averageDailySales": 4,
    "leadTimeDays": 5,
    "unitCost": 18.50,
    "criticalityLevel": 3
  }'
```

### Listar todas as peças

```bash
curl http://localhost:8080/parts/
```

### Listar peças por categoria

```bash
curl "http://localhost:8080/parts/?category=engine"
```

### Buscar peça por ID

```bash
curl http://localhost:8080/parts/{id}
```

### Atualizar peça

```bash
curl -X PUT http://localhost:8080/parts/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Filtro de Óleo X",
    "category": "engine",
    "currentStock": 40,
    "minimumStock": 20,
    "averageDailySales": 4,
    "leadTimeDays": 5,
    "unitCost": 18.50,
    "criticalityLevel": 3
  }'
```

### Remover peça

```bash
curl -X DELETE http://localhost:8080/parts/{id}
```

### Priorização de reposição

```bash
curl http://localhost:8080/restock/priorities
```

Resposta (apenas peças que precisam de reposição, ordenadas por urgência):

```json
{
  "priorities": [
    {
      "partId": "uuid-1",
      "name": "Filtro de Óleo X",
      "category": "engine",
      "currentStock": 15,
      "projectedStock": -5,
      "minimumStock": 20,
      "urgencyScore": 75
    }
  ]
}
```

## Regras de negócio

Implementadas em `internal/domain/priority.go`:

1. `expectedConsumption = averageDailySales * leadTimeDays`
2. `projectedStock = currentStock - expectedConsumption`
3. `needsRestock = projectedStock < minimumStock`
4. `urgencyScore = (minimumStock - projectedStock) * criticalityLevel`

Critérios de desempate (em ordem): maior `urgencyScore` → maior
`criticalityLevel` → maior `averageDailySales` → ordem alfabética pelo nome.

### Tratamento de casos extremos

- **Estoque atual negativo**: tratado como cenário de negócio válido (não
  é rejeitado na validação de entrada). Resulta em `projectedStock` ainda
  mais negativo e, consequentemente, `urgencyScore` mais alto.
- **Venda média zero**: `expectedConsumption = 0`, logo
  `projectedStock = currentStock`. A peça só precisa de reposição se o
  estoque atual já estiver abaixo do mínimo.
- **Lead time alto**: aumenta proporcionalmente o consumo esperado e,
  portanto, a urgência — sem teto artificial.

## Testes

```bash
go test ./... -v
```

Para incluir o detector de condição de corrida (relevante para a camada
de persistência, que precisa suportar acesso concorrente):

```bash
go test ./... -race -v
```

Cobertura:

- **`internal/domain`**: testes unitários puros do cálculo de prioridade
  — cenários extremos (estoque negativo, venda média zero, lead time
  alto), critérios de desempate isolados, e teste de limite
  (`projectedStock == minimumStock`). A validação de invariantes da
  entidade (`Part.Validate()`) é exercitada via `internal/service` e
  nos testes de integração HTTP.
- **`internal/repository`**: testes funcionais do CRUD em memória, e
  testes de concorrência (`-race`) que validam o `RWMutex` sob criação e
  leitura simultâneas — relevante para o requisito de suportar centenas
  ou milhares de peças.
- **`internal/service`**: testes com um repositório _fake_ (implementação
  de teste da interface `PartRepository`), validando a orquestração entre
  domínio e persistência sem depender de infraestrutura real.
- **`internal/http/handler`**: testes de integração via `httptest`, com a
  aplicação completa montada (repositório real + services reais +
  handlers + router), validando o fluxo HTTP de ponta a ponta.

## Trocando o banco de dados no futuro

A interface `PartRepository` (`internal/repository/part_repository.go`)
define o contrato. Para trocar a persistência:

1. Criar uma nova struct (ex: `PostgresPartRepository`) implementando os
   mesmos métodos da interface.
2. Trocar a linha de instanciação em `cmd/api/main.go`:

```go
   // antes
   partRepo := repository.NewInMemoryPartRepository()
   // depois
   partRepo := repository.NewPostgresPartRepository(db)
```

3. Nenhuma outra camada (`service`, `handler`, `domain`) precisa de
   qualquer alteração.
