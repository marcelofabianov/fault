Com certeza, Marcelo.

Com base em todas as nossas discussões, melhorias no código e o plano que traçamos para a documentação, preparei a versão final e aprimorada do README.md.

Este documento agora segue a "jornada do desenvolvedor", começando com um Quick Start prático e avançando para padrões de uso do mundo real, como havíamos planejado.

Aqui está o README.md completo em um único bloco, pronto para ser copiado.

Markdown

# fault

[![Go Report Card](https://goreportcard.com/badge/github.com/marcelofabianov/fault)](https://goreportcard.com/report/github.com/marcelofabianov/fault)
[![Go Reference](https://pkg.go.dev/badge/github.com/marcelofabianov/fault.svg)](https://pkg.go.dev/github.com/marcelofabianov/fault)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fault` é uma biblioteca Go leve, porém poderosa, para a criação de erros estruturados e ricos em contexto. Enquanto o pacote `errors` padrão do Go é excelente para o encadeamento de erros, ele não oferece uma maneira nativa de transportar dados estruturados como códigos de erro e metadados, que são cruciais para APIs modernas, logging e depuração. `fault` preenche essa lacuna, permitindo que as aplicações modelem suas falhas de forma clara e consistente através de todas as camadas da arquitetura.

## Principais Funcionalidades

* **Erros Estruturados:** Crie erros com códigos, mensagens e um mapa de contexto customizável.
* **Error Wrapping Idiomático:** Totalmente compatível com o pacote `errors` do Go, incluindo `errors.Is` e `errors.As`.
* **API Fluida:** Use o padrão *Functional Options* para construir erros de forma declarativa e legível.
* **Erros Aninhados:** Suporte para múltiplos erros detalhados, ideal para cenários de validação.
* **Utilitários HTTP:** Converta erros internos em respostas de API HTTP padrão com facilidade.


## Instalação

```bash
go get github.com/marcelofabianov/fault
```
### 🚀 Quick Start: Em 30 Segundos

Este exemplo mostra o ciclo de vida completo: criar um erro `fault`, encapsulá-lo, e lidar com ele de forma robusta ao lado de erros padrão do Go.

```go
package main

import (
	"errors"
	"fmt"
	"io" // Apenas para simular um erro padrão do Go

	"github.com/marcelofabianov/fault"
)

// Em sua aplicação, uma função central para processar erros precisa lidar com qualquer tipo de `error`.
func handleFinalError(err error) {
	var fErr *fault.Error

	// Use errors.As para verificar se o erro é do nosso tipo `fault.Error`
	if errors.As(err, &fErr) {
		// SIM! É um erro estruturado. Podemos acessar seus dados.
		fmt.Println("--> Erro estruturado detectado!")
		fmt.Printf("    Código: %s\n", fErr.Code)
		fmt.Printf("    Contexto: %v\n", fErr.Context)
	} else {
		// NÃO! É um erro genérico do Go. Tratamos de forma padrão.
		fmt.Println("--> Erro genérico detectado!")
		fmt.Printf("    Mensagem: %s\n", err.Error())
	}
}

func main() {
	// 1. Crie um erro de domínio com código e contexto
	originalErr := fault.New(
		"user balance is insufficient",
		fault.WithCode(fault.DomainViolation),
		fault.WithContext("user_id", "usr_123"),
	)

	// 2. Encapsule-o com mais contexto em uma camada superior
	wrappedErr := fault.Wrap(originalErr, "failed to process payment")

	// 3. Crie um erro padrão do Go para comparação
	standardErr := io.EOF

	// --- Processando os erros ---
	fmt.Println("--- Lidando com um erro 'fault' ---")
	handleFinalError(wrappedErr)

	fmt.Println("\n--- Lidando com um erro padrão do Go ---")
	handleFinalError(standardErr)
}
```

### Conceitos Fundamentais

#### 1. Criando e Encapsulando Erros

Use `fault.New` para criar um novo erro e `fault.Wrap` para encapsular um erro existente, adicionando mais contexto.

```go
// Erro simples com um código
notFoundErr := fault.New("product not found", fault.WithCode(fault.NotFound))

// Erro com código e contexto para logging
domainErr := fault.New(
    "account is suspended",
    fault.WithCode(fault.DomainViolation),
    fault.WithContext("account_id", "acc_456"),
)

// Encapsulando um erro de banco de dados
dbErr := errors.New("connection refused")
infraErr := fault.Wrap(
    dbErr,
    "failed to query database",
    fault.WithCode(fault.InfraError),
)
```

#### 2. Verificando Erros

Use as funções `Is[Code]` para verificações semânticas ou `IsCode` para checar um código específico. Para extrair o erro `*fault.Error`, use o idiomático `errors.As` ou a função de conveniência `fault.AsFault`.

```go
err := fault.New(
    "access denied",
    fault.WithCode(fault.Forbidden),
    fault.WithContext("user_role", "guest"),
)

// Usando as funções auxiliares (preferencial)
if fault.IsForbidden(err) {
    fmt.Println("Access is forbidden.")
}

// Verificando um código específico
if fault.IsCode(err, fault.Forbidden) {
    fmt.Println("This also works.")
}

// --- Extraindo o erro para ler seu conteúdo ---

// Opção 1: Usando o padrão da biblioteca `errors.As` (Recomendado)
var fErr *fault.Error
if errors.As(err, &fErr) {
    fmt.Printf("Contexto do erro: %v\n", fErr.Context)
}

// Opção 2: Usando a função de conveniência `fault.AsFault`
if fErr, ok := fault.AsFault(err); ok {
    fmt.Printf("Contexto do erro (usando AsFault): %v\n", fErr.Context)
}
```

### Padrões de Uso e Receitas

#### Receita 1: Tratamento de Erros em uma API RESTful

`fault` simplifica a tradução de erros internos para respostas HTTP consistentes.

Imagine que sua camada de serviço possa retornar um erro NotFound como este:

```go
// Em sua camada de serviço...
func (s *Service) GetEntity(id string) (*Entity, error) {
    // ... lógica para buscar a entidade ...
    if entityNotFound {
        return nil, fault.New(
            "entity not found",
            fault.WithCode(fault.NotFound),
            fault.WithContext("entity_id", id),
        )
    }
    return &entity, nil
}
```

O seu handler HTTP pode então tratar esse erro de forma agnóstica:

```go
func GetEntityHandler(w http.ResponseWriter, r *http.Request) {
    entityID := r.URL.Query().Get("id")

    // A camada de serviço retorna um erro (fault ou não)
    entity, err := service.GetEntity(entityID)
    if err != nil {
        // Converte qualquer erro em uma resposta estruturada e padronizada
        response := fault.ToResponse(err)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(response.StatusCode) // O status code será 404
        json.NewEncoder(w).Encode(response)
        return
    }

    // ... lógica de sucesso ...
    json.NewEncoder(w).Encode(entity)
}
```

**Exemplo de Saída JSON (para o erro NotFound acima):**

A resposta HTTP teria o status 404 Not Found e o seguinte corpo JSON:

```json
{
  "message": "entity not found",
  "code": "not_found",
  "context": {
    "entity_id": "some-id"
  }
}
```

#### Receita 2: Integração com Logging Estruturado

O contexto dos erros `fault` permitefault é ideal para logging estruturado, como o `slog`.

```go
import "log/slog"

// ...
_, err := service.DoSomething()
if err != nil {
    var fErr *fault.Error
    if errors.As(err, &fErr) {
        slog.Error(
            fErr.Message,
            "error_code", fErr.Code,
            "error_context", fErr.Context,
        )
    } else {
        slog.Error(err.Error())
    }
}
```

#### Receita 3: Validação de Requisições com go-playground/validator

Converta os erros detalhados da biblioteca `go-playground/validator` em um único `fault.Error` com múltiplos detalhes, um para cada campo inválido.

Imagine uma requisição com o seguinte `RequestBody`:

```json
type RequestBody struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=18"`
}
```

Seu handler pode usar `fault` para criar uma resposta de erro padronizada:

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req RequestBody
    // ... código para decodificar o JSON do corpo da requisição ...

    validate := validator.New()
    if errs := validate.Struct(req); errs != nil {
        // Converte os erros do validador para um fault.Error com múltiplos detalhes
        faultErr := fault.NewValidationErrorFromValidator(errs.(validator.ValidationErrors))

        // Agora `faultErr` pode ser tratado como qualquer outro erro `fault`
        response := fault.ToResponse(faultErr)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(response.StatusCode) // O status code será 400
        json.NewEncoder(w).Encode(response)
        return
    }

    // ... lógica de sucesso ...
}
```

**Exemplo de Saída JSON (para uma requisição vazia):**

A resposta HTTP teria o status `400 Bad Request` e um corpo JSON detalhando cada campo que falhou na validação:

```json
{
  "message": "Request validation failed",
  "code": "invalid_input",
  "details": [
    {
      "message": "validation failed on field 'Name'",
      "code": "invalid_input",
      "context": {
        "field": "Name",
        "tag": "required",
        "param": ""
      }
    },
    {
      "message": "validation failed on field 'Email'",
      "code": "invalid_input",
      "context": {
        "field": "Email",
        "tag": "required",
        "param": ""
      }
    },
    {
      "message": "validation failed on field 'Age'",
      "code": "invalid_input",
      "context": {
        "field": "Age",
        "tag": "gte",
        "param": "18"
      }
    }
  ]
}
```

#### Receita 4: Modelagem de Domínio com Value Objects (Padrão DDD)

O `fault` se integra perfeitamente com bibliotecas de Value Objects (como a [wisp](https://github.com/marcelofabianov/wisp)), um padrão comum em Domain-Driven Design (DDD). A ideia é usar tipos fortes para garantir as invariantes do seu domínio (ex: um nome não pode ser vazio, uma quantidade deve ser positiva).

Quando a criação de um Value Object falha, `fault` pode ser usado para encapsular esse erro de validação de baixo nível, adicionando um código de erro padronizado (`Invalid`) e um contexto rico que inclui o valor original que causou a falha. Isso cria uma fronteira de tradução de erros muito clara entre seu domínio e sua camada de aplicação.

```go
package model

import (
  "github.com/marcelofabianov/fault"
  "github.com/marcelofabianov/wisp"
)

type NewCourseInput struct {
  Name           string
  Description    string
  MaxEnrollments int
}

// ...

func NewCourse(input NewCourseInput, createdBy wisp.AuditUser) (*Course, error) {
  // ... (criação de ID) ...

  name, err := wisp.NewNonEmptyString(input.Name)
  if err != nil {
    // Encapsula o erro de validação do Value Object com contexto rico
    return nil, fault.Wrap(err,
      "Invalid name",
      fault.WithCode(fault.Invalid),
      fault.WithContext("name", input.Name), // Adiciona o valor inválido ao contexto
    )
  }

  description, err := wisp.NewNonEmptyString(input.Description)
  if err != nil {
    return nil, fault.Wrap(err,
      "Invalid description",
      fault.WithCode(fault.Invalid),
      fault.WithContext("description", input.Description),
    )
  }
  // ... resto da lógica ...
}
```

**Exemplo: Construtor de uma Entidade Course**

No exemplo abaixo, a função `NewCourse` tenta criar Value Objects (`wisp.NonEmptyString`, `wisp.PositiveInt`). Se qualquer uma dessas validações falhar, o erro retornado pela `wisp` é encapsulado por `fault` para criar um erro de aplicação estruturado.

---

#### Referência da API

| Código (`fault.Code`) | Status HTTP Padrão     | Descrição                                 |
| :-------------------- | :--------------------- | :---------------------------------------- |
| `NotFound`            | 404 Not Found          | O recurso solicitado não foi encontrado.  |
| `Invalid`             | 400 Bad Request        | A entrada fornecida é inválida.           |
| `DomainViolation`     | 422 Unprocessable Entity| Uma regra de negócio foi violada.         |
| `Forbidden`           | 403 Forbidden          | Acesso negado à funcionalidade.           |
| `Unauthorized`        | 401 Unauthorized       | Autenticação necessária ou falhou.        |
| `Conflict`            | 409 Conflict           | Conflito de estado, ex: recurso já existe.|
| `InfraError`          | 502 Bad Gateway        | Falha em um serviço externo ou infra.     |
| `Internal`            | 500 Internal Server Error| Erro inesperado e não tratado.          |



## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma *issue* para discutir novas funcionalidades ou reportar bugs.

## Licença

Este projeto é distribuído sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
