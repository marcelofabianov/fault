package service

import (
	"fmt"

	"github.com/marcelofabianov/wisp"

	"github.com/marcelofabianov/example/domain"
)

// UserService lida com a lógica de negócio para usuários.
// Usamos um mapa em memória para simular um banco de dados.
type UserService struct {
	users map[string]*domain.User
}

func NewUserService() *UserService {
	// Pré-populamos o serviço com alguns dados de exemplo
	users := make(map[string]*domain.User)
	id1, _ := wisp.ParseUUID("0199910e-d028-7e10-aa14-3954e781a9bc")
	name1, _ := wisp.NewNonEmptyString("Marcelo Fabiano")
	email1, _ := wisp.NewNonEmptyString("marcelo.fabiano@example.com")
	users[id1.String()] = &domain.User{ID: id1, Name: name1, Email: email1, Active: true}

	id2, _ := wisp.ParseUUID("0199910f-4238-7882-a075-b7842b4d41f0")
	name2, _ := wisp.NewNonEmptyString("John Doe")
	email2, _ := wisp.NewNonEmptyString("john.doe@example.com")
	users[id2.String()] = &domain.User{ID: id2, Name: name2, Email: email2, Active: false} // Usuário inativo

	return &UserService{users: users}
}

// FindUserByID busca um usuário pelo ID.
// Retorna um erro de domínio específico se uma regra for violada.
func (s *UserService) FindUserByID(id string) (*domain.User, error) {
	user, exists := s.users[id]
	if !exists {
		// Regra de negócio: usuário deve existir.
		return nil, domain.NewUserNotFound(id)
	}

	if !user.Active {
		// Regra de negócio: usuário deve estar ativo.
		return nil, domain.NewUserInactive(id)
	}

	return user, nil
}

// RegisterUser simula o registro de um novo usuário.
func (s *UserService) RegisterUser(id, name, email string) error {
	if _, exists := s.users[id]; exists {
		// Regra de negócio: não pode registrar um usuário com ID duplicado.
		return domain.NewUserAlreadyExists("id", id)
	}

	// Lógica de registro...
	fmt.Printf("Usuário com ID %s registrado com sucesso!\n", id)
	return nil
}
