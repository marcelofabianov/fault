# fault

[![Go Report Card](https://goreportcard.com/badge/github.com/marcelofabianov/fault)](https://goreportcard.com/report/github.com/marcelofabianov/fault)
[![Go Reference](https://pkg.go.dev/badge/github.com/marcelofabianov/fault.svg)](https://pkg.go.dev/github.com/marcelofabianov/fault)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`fault` é uma biblioteca Go leve, porém poderosa, para a criação de erros estruturados e ricos em contexto. Ela foi projetada para permitir que as aplicações modelem suas falhas de forma clara e consistente através de todas as camadas da arquitetura, desde o domínio até a apresentação.

## ✨ Principais Funcionalidades

* **Erros Estruturados:** Crie erros com códigos, mensagens e um mapa de contexto customizável para facilitar a depuração e o logging.
* **Error Wrapping Idiomático:** Totalmente compatível com o pacote `errors` do Go, incluindo `errors.Is` e `errors.As`.
* **API Fluida:** Use o padrão *Functional Options* para construir erros de forma declarativa e legível.
* **Erros Aninhados:** Suporte para múltiplos erros detalhados, ideal para cenários complexos como a validação de formulários.
* **Desacoplado de Protocolos:** O núcleo do `fault` é agnóstico a protocolos. Utilitários para APIs web (HTTP) estão disponíveis em um sub-pacote (`httputil`) para manter as responsabilidades separadas.

## 🚀 Instalação

```bash
go get [github.com/marcelofabianov/fault](https://github.com/marcelofabianov/fault)
```

## 💡 Uso e Conceitos

### Criação Básica de Erros

Use `fault.New` para criar um erro simples e o padrão *Functional Options* para enriquecê-lo.

```go
import "[github.com/marcelofabianov/fault](https://github.com/marcelofabianov/fault)"

// Um erro simples com um código
err := fault.New(
    "user not found",
    fault.WithCode(fault.NotFound),
    fault.WithContext("user_id", "usr-123"),
)

fmt.Println(err)
// Saída: user not found
```

### Embrulhando (Wrapping) Erros Existentes

Use `fault.Wrap` para adicionar contexto de negócio a um erro técnico de uma camada inferior.

```go
import (
    "os"
    "[github.com/marcelofabianov/fault](https://github.com/marcelofabianov/fault)"
)

func readFile() error {
    f, err := os.Open("meu-arquivo.txt")
    if err != nil {
        // Envolve o erro original 'os.PathError' com nosso erro estruturado
        return fault.Wrap(err,
            "failed to open configuration file",
            fault.WithCode(fault.Internal),
            fault.WithContext("filename", "meu-arquivo.txt"),
        )
    }
    defer f.Close()
    return nil
}
```

### Erros de Validação com Detalhes

O campo `Details` é perfeito para retornar múltiplos erros de uma só vez.

```go
import "[github.com/marcelofabianov/fault](https://github.com/marcelofabianov/fault)"

func validateRequest(email, password string) error {
    var details []*fault.Error

    if email == "" {
        details = append(details, fault.New("email is required", fault.WithContext("field", "email")))
    }
    if password == "" {
        details = append(details, fault.New("password is required", fault.WithContext("field", "password")))
    }

    if len(details) > 0 {
        return fault.New(
            "validation failed",
            fault.WithCode(fault.Invalid),
            fault.WithDetails(details...),
        )
    }

    return nil
}
```

### Integração com APIs HTTP (`httputil`)

O sub-pacote `httputil` ajuda a traduzir um `*fault.Error` em uma resposta JSON padronizada.

```go
import (
    "encoding/json"
    "net/http"

    "[github.com/marcelofabianov/fault](https://github.com/marcelofabianov/fault)"
    "[github.com/marcelofabianov/fault/httputil](https://github.com/marcelofabianov/fault/httputil)"
)

// Em sua camada de serviço:
func findUser(userID string) *fault.Error {
    if userID == "" {
        return fault.New("user ID cannot be empty", fault.WithCode(fault.Invalid))
    }
    // ... lógica para buscar o usuário ...
    return fault.New("user not found", fault.WithCode(fault.NotFound), fault.WithContext("searched_id", userID))
}


// Em seu HTTP handler:
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    // A lógica de negócio retorna um *fault.Error
    err := findUser("user-abc")

    if err != nil {
        // Converte o erro em uma resposta JSON com o código de status correto
        response := httputil.ToResponse(err)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(response.StatusCode)
        json.NewEncoder(w).Encode(response)
        return
    }

    // ... lógica de sucesso ...
}
```

#### Exemplo de Resposta JSON de Erro

O handler acima, ao receber o erro `NotFound`, produziria a seguinte resposta JSON com status `404 Not Found`:

```json
{
    "message": "user not found",
    "code": "not_found",
    "context": {
        "searched_id": "user-abc"
    }
}
```

## 🤝 Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma *issue* para discutir novas funcionalidades ou reportar bugs.

## 📄 Licença

Este projeto é distribuído sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
