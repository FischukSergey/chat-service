package clientv1

import (
	"context"
	"fmt"

	"github.com/FischukSergey/chat-service/internal/usecases/client/get-history"
)


//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clientv1mocks
type getHistoryUseCase interface {
	Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error)
}

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	getHistoryUseCase getHistoryUseCase `option:"mandatory" validate:"required"`
	// Ждут своего часа.
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	// валидация опций
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate options: %v", err)
	}
	return Handlers{Options: opts}, nil
}
