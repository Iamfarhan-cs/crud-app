package storage

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Iamfarhan-cs/crud-app/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type MemoryUserStore struct {
	mu     sync.RWMutex
	nextID int
	users  map[string]models.User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		nextID: 1,
		users:  make(map[string]models.User),
	}
}

func (s *MemoryUserStore) List() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users
}

func (s *MemoryUserStore) Get(id string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

func (s *MemoryUserStore) Create(input models.UserInput) (models.User, error) {
	if err := validateUser(input); err != nil {
		return models.User{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprint(s.nextID)
	s.nextID++

	user := models.User{
		ID:    id,
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
	}
	s.users[id] = user

	return user, nil
}

func (s *MemoryUserStore) Update(id string, input models.UserInput) (models.User, error) {
	if err := validateUser(input); err != nil {
		return models.User{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return models.User{}, ErrUserNotFound
	}

	user := models.User{
		ID:    id,
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
	}
	s.users[id] = user

	return user, nil
}

func (s *MemoryUserStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return false
	}

	delete(s.users, id)
	return true
}

func validateUser(input models.UserInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return errors.New("email is required")
	}
	return nil
}
