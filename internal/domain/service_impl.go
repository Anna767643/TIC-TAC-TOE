package domain

// gameService - реализация GameService.
type gameService struct {
	repo GameRepository
}

// NewGameService - конструктор для di. Принимает интерфейс репозитория.
// func NewGameService(repo GameRepository) GameService
func NewGameService(repo GameRepository) GameService {
	return &gameService{
		repo: repo,
	}
}

// userService - реализация UserService.
type userService struct {
	repo UserRepository
}

// NewUserService - конструктор для di. Принимает интерфейс репозитория.
func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
