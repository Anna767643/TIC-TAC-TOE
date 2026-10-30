package di_test

import (
	"src/internal/di"
	"testing"

	"go.uber.org/fx"
)

func TestDependencyInjection(t *testing.T) {
	err := fx.ValidateApp(
		fx.NopLogger,
		di.Module,
	)

	if err != nil {
		t.Fatalf("Module broke: %v", err)
	}
}
