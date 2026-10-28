package web

// // mockGameService - заглушка GameService,
// type mockGameService struct {
// 	NextMoveFunc   func(ctx context.Context, game domain.Game) (domain.Game, error)
// 	ValidateFunc   func(incoming domain.Board, stored domain.Board) error
// 	IsFinishedFunc func(board domain.Board) (bool, int)
// }

// func (m *mockGameService) NextMove(ctx context.Context, game domain.Game) (domain.Game, error) {
// 	return m.NextMoveFunc(ctx, game)
// }

// func (m *mockGameService) Validate(incoming domain.Board, stored domain.Board) error {
// 	return m.ValidateFunc(incoming, stored)
// }

// func (m *mockGameService) IsFinished(board domain.Board) (bool, int) {
// 	return m.IsFinishedFunc(board)
// }

// // mockUserRepository - заглушка для репозитория
// type mockGameRepository struct {
// 	SaveFunc func(ctx context.Context, game domain.Game) error
// 	GetFunc  func(ctx context.Context, id uuid.UUID) (domain.Game, error)
// }

// func (m *mockGameRepository) Save(ctx context.Context, game domain.Game) error {
// 	return m.SaveFunc(ctx, game)
// }

// func (m *mockGameRepository) Get(ctx context.Context, id uuid.UUID) (domain.Game, error) {
// 	return m.GetFunc(ctx, id)
// }
