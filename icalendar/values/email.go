package values

import (
	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar/properties"
)

const (
	preferredTypeValue = "pref"
)

type Email struct {
	Mail        string
	Types       []string
	Label       string
	IsPreferred bool
}

func NewEmail(mail string, preferred bool, types ...string) *Email {
	return &Email{
		Mail:        mail,
		IsPreferred: preferred,
		Types:       types,
	}
}

func (e *Email) ValidateICalValue() error {
	return nil
}

func (e *Email) EncodeICalParams() (properties.Params, error) {
	params := make(properties.Params, len(e.Types))
	for i, str := range e.Types {
		params[i] = properties.Param{
			Name:  properties.ParameterType,
			Value: str,
		}
	}

	if e.IsPreferred {
		params = append(params, properties.Param{
			Name:  properties.ParameterType,
			Value: preferredTypeValue,
		})
	}

	return params, nil
}

func (e *Email) DecodeICalParams(params properties.Params) error {
	for _, param := range params {
		if strings.EqualFold(string(param.Name), properties.ParameterType) {
			if param.Value == preferredTypeValue {
				e.IsPreferred = true
			} else {
				e.Types = append(e.Types, param.Value)
			}
		} else if strings.EqualFold(string(param.Name), properties.ABLabel) {
			e.Label = param.Value
		}
	}
	return nil
}

func (e *Email) EncodeICalValue() (string, error) {
	return e.Mail, nil
}

func (e *Email) DecodeICalValue(value string) error {
	e.Mail = value
	return nil
}

func (e *Email) EncodeICalName() (properties.PropertyName, error) {
	return properties.EmailPropertyName, nil
}

func (e *Email) EncodeAdditionalICalProperties() ([]*properties.Property, error) {
	var props []*properties.Property
	if e.Label != "" {
		props = append(props, properties.NewProperty("X-ABLabel", e.Label))
	}

	return props, nil
}
