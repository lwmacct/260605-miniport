package handler

type BodyDTO[T any] struct {
	Body T
}

type BodyInputDTO[T any] struct {
	Body T
}

type DeleteDTO struct {
	Deleted bool `json:"deleted" example:"true"`
}
