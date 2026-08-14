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

type ActionDTO struct {
	OK bool `json:"ok"`
}

type ActionOutputDTO struct {
	Body ActionDTO
}

type BatchIDsDTO struct {
	IDs []string `json:"ids" nullable:"false" minItems:"1"`
}
