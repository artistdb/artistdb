package location

import (
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type Location struct {
	ID   string
	Name string
}

func (l Location) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", l.ID)

	return nil
}

func New() *Location {
	return &Location{ID: uuid.New().String()}
}
