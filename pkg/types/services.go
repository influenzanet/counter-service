package types

import(
	"context"
)

type RegistryService interface {
	 Handle(ctx context.Context, input <-chan[]Counter ) error 
	 Read()[]Counter
}

type CollectorService interface {
	 Run(ctx context.Context, out chan<-[]Counter) error 
}