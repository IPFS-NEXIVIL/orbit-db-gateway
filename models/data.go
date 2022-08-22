package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Data struct {
	ID          string `mapstructure:"id" json:"-" validate:"uuid_rfc4122"`
	InReplyToID string `mapstructure:"in-reply-to-id" json:"-" validate:"omitempty,uuid_rfc4122"`
	Date        int64  `mapstructure:"date" json:"-" validate:"required,number"`
	Body        string `mapstructure:"body" json:"-" validate:"required,min=3,max=524288"`

	Replies     []*Data `mapstructure:"-" json:"-" validate:"-"`
	LatestReply int64   `mapstructure:"-" json:"-" validate:"-"`

	Read bool `mapstructure:"-" json:"read" validate:"-"`
}

func NewData() *Data {
	data := new(Data)

	id, _ := uuid.NewUUID()
	data.ID = id.String()

	data.Date = time.Now().UnixNano() / int64(time.Millisecond)

	return data
}

func (data *Data) IsValid() (bool, error) {
	validate := validator.New()
	errs := validate.Struct(data)
	if errs != nil {
		// validationErrors := errs.(validator.ValidationErrors)
		return false, errs
	}

	return true, nil
}
