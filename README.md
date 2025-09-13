# fault

[![Go Report Card](https://goreportcard.com/badge/github.com/marcelofabianov/fault)](https://goreportcard.com/report/github.com/marcelofabianov/fault)
[![Go Reference](https://pkg.go.dev/badge/github.com/marcelofabianov/fault.svg)](https://pkg.go.dev/github.com/marcelofabianov/fault)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fault` é uma biblioteca Go leve, porém poderosa, para a criação de erros estruturados e ricos em contexto. Ela foi projetada para permitir que as aplicações modelem suas falhas de forma clara e consistente através de todas as camadas da arquitetura, desde o domínio até a apresentação.

## Principais Funcionalidades

* **Erros Estruturados:** Crie erros com códigos, mensagens e um mapa de contexto customizável para facilitar a depuração e o logging.
* **Error Wrapping Idiomático:** Totalmente compatível com o pacote `errors` do Go, incluindo `errors.Is` e `errors.As`.
* **API Fluida:** Use o padrão *Functional Options* para construir erros de forma declarativa e legível.
* **Erros Aninhados:** Suporte para múltiplos erros detalhados, ideal para cenários complexos como a validação de formulários.

## Instalação

```bash
go get github.com/marcelofabianov/fault
```

## Uso e Conceitos

### 1. Construtores

Funções para criar novos erros e adicionar informações contextuais.

```go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/marcelofabianov/fault"
)

func main() {
	// Criação de um erro de validação com contexto e um erro original
	err := fault.NewValidationError(
		errors.New("email is already in use"),
		"User data is invalid",
		map[string]any{"field": "email"},
	)

	// Criação de um erro de domínio simples
	domainErr := fault.New("Account is suspended", fault.WithCode(fault.DomainViolation))

	// Criando um erro que encapsula outro
	wrappedErr := fault.Wrap(
		err,
		"Failed to create new user account",
		fault.WithDetails(domainErr),
	)

	fmt.Println(wrappedErr.Error())
	// Saída: Failed to create new user account: User data is invalid: email is already in use
}
```

### 2. Verificadores

Funções para verificar o tipo ou código de um erro de forma segura, percorrendo a cadeia de erros.

```go
package main

import (
	"errors"
	"fmt"

	"github.com/marcelofabianov/fault"
)

func main() {
	// Exemplo de um erro complexo
	validationErr := fault.New("email is invalid", fault.WithCode(fault.Invalid))
	wrappedErr := fault.Wrap(validationErr, "could not process request")

	// Usando a função IsCode para verificar o código em qualquer nível da cadeia
	if fault.IsCode(wrappedErr, fault.Invalid) {
		fmt.Println("Error is of type 'invalid_input'")
	}

	// Usando as funções de verificação específicas
	if fault.IsInvalid(wrappedErr) {
		fmt.Println("Using the specific checker for 'invalid_input'")
	}

	// Cenário negativo
	if !fault.IsNotFound(wrappedErr) {
		fmt.Println("Error is not of type 'not_found'")
	}

	// Verificando um erro genérico
	if !fault.IsCode(errors.New("a simple error"), fault.Internal) {
		fmt.Println("Generic error is not a 'fault' error")
	}
}
```

### 3. Utilitários HTTP

Utilitários para converter códigos de erro internos em códigos de status HTTP padrão, ideal para a camada de API.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/marcelofabianov/fault"
)

func main() {
	// Mapeando um código de erro para um status HTTP
	statusCode := fault.GetHTTPStatusCode(fault.NotFound)
	fmt.Printf("HTTP Status for 'not_found' is: %d\n", statusCode)
	// Saída: HTTP Status for 'not_found' is: 404

	// Mapeando um código de erro de negócio para 422 Unprocessable Entity
	domainStatus := fault.GetHTTPStatusCode(fault.DomainViolation)
	fmt.Printf("HTTP Status for 'domain_violation' is: %d\n", domainStatus)
	// Saída: HTTP Status for 'domain_violation' is: 422

	// Um código de erro desconhecido retorna 500
	unknownStatus := fault.GetHTTPStatusCode("unknown_code")
	fmt.Printf("HTTP Status for an unknown code is: %d\n", unknownStatus)
	// Saída: HTTP Status for an unknown code is: 500
}
```

### 4. Respostas HTTP

A seguir, um exemplo de como converter um erro de acesso negado em uma resposta HTTP serializável. A estrutura ErrorResponse encapsula todas as informações relevantes, incluindo o código de erro (forbidden) e contexto, que podem ser usados pela aplicação cliente para exibir uma mensagem adequada.

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/marcelofabianov/fault"
)

func main() {
	// Simulação de um erro de acesso negado, sem permissão para executar a operação.
	err := fault.New(
		"Access denied: you do not have permission to perform this action.",
		fault.WithCode(fault.Forbidden),
		fault.WithContext("user_id", "12345"),
		fault.WithContext("required_role", "admin"),
	)

	// Convertendo o erro para uma estrutura de resposta HTTP
	response := fault.ToResponse(err)

	// Imprimindo o status code para a camada de framework
	fmt.Printf("HTTP Status Code: %d\n", response.StatusCode)
	// Saída: HTTP Status Code: 403

	// Serializando a resposta para JSON
	encoder.Encode(response)
	/*
	{
	  "message": "Access denied: you do not have permission to perform this action.",
	  "code": "forbidden",
	  "context": {
	    "required_role": "admin",
	    "user_id": "12345"
	  }
	}
	*/
}
```

### 5. Integrando com `go-playground/validator`

O pacote `fault` oferece uma maneira fluida de integrar bibliotecas de validação, como o `go-playground/validator/v10`, convertendo seus erros específicos em um formato estruturado e consistente. Essa abordagem simplifica a camada de API, garantindo que todos os erros de validação sigam um único padrão.

Para isso, o pacote `fault` expõe o erro sentinela `ErrValidation` e uma função de conveniência que faz toda a conversão para você.

```go
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/go-playground/validator/v10"
	"github.com/marcelofabianov/fault"
)

// RequestBody represents a request payload to be validated.
type RequestBody struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=18"`
}

func main() {
	validate := validator.New()

	// Simula um handler de API
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RequestBody
		// Simula um erro de validação
		if errs := validate.Struct(req); errs != nil {
			// Converte os erros do validador para um fault.Error
			faultErr := fault.NewValidationErrorFromValidator(errs.(validator.ValidationErrors))

			// O desenvolvedor pode agora verificar o erro de forma idiomática
			if errors.Is(faultErr, fault.ErrValidation) {
				response := fault.ToResponse(faultErr)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(response.StatusCode)
				json.NewEncoder(w).Encode(response)
			}
		}
	})

	// Executa o handler e captura a resposta
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Imprime a resposta completa para fins de demonstração
	// Status HTTP: 400 Bad Request
	// Body: {"message":"Request validation failed", ...}
	fmt.Printf("HTTP Status Code: %d\n", rr.Code)
	fmt.Println("---")
	fmt.Println("HTTP Response Body (JSON):")
	fmt.Println(rr.Body.String())
}
```

A saída do console seria:

```sh
HTTP Status Code: 400
---
HTTP Response Body (JSON):
{"message":"Request validation failed","code":"invalid_input","details":[{"message":"validation failed on field 'Name'","code":"invalid_input","context":{"field":"Name","param":"","tag":"required"}},{"message":"validation failed on field 'Email'","code":"invalid_input","context":{"field":"Email","param":"","tag":"required"}},{"message":"validation failed on field 'Age'","code":"invalid_input","context":{"field":"Age","param":"18","tag":"gte"}}]}
```

O json de saída formato:

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

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma *issue* para discutir novas funcionalidades ou reportar bugs.

## Licença

Este projeto é distribuído sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
