module github.com/marcelofabianov/example

go 1.25.1

require (
	github.com/marcelofabianov/fault v1.4.0
	github.com/marcelofabianov/wisp v1.9.0
)

// Adicione esta linha para mapear o módulo para o diretório local
replace github.com/marcelofabianov/example => ./

require (
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)
