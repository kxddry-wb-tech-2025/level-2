package handlers

type Storage interface {
	Create(ctx context.Context, event models.Event) error
}
