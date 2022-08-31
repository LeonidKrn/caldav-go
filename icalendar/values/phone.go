package values

import (
	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar/properties"
)

type Phone struct {
	Number      string
	Types       []string
	IsPreferred bool
}

func NewPhone(number string, preferred bool, types ...string) *Phone {
	return &Phone{
		Number:      number,
		IsPreferred: preferred,
		Types:       types,
	}
}

func (p *Phone) ValidateICalValue() error {
	return nil
}

func (p *Phone) EncodeICalParams() (properties.Params, error) {
	params := make(properties.Params, len(p.Types))
	for i, str := range p.Types {
		params[i] = properties.Param{
			Name:  properties.ParameterType,
			Value: str,
		}
	}

	if p.IsPreferred {
		params = append(params, properties.Param{
			Name:  properties.ParameterType,
			Value: preferredTypeValue,
		})
	}

	return params, nil
}

func (p *Phone) DecodeICalParams(params properties.Params) error {
	for _, param := range params {
		if strings.EqualFold(string(param.Name), properties.ParameterType) {
			if param.Value == preferredTypeValue {
				p.IsPreferred = true
			} else {
				p.Types = append(p.Types, param.Value)
			}
		}
	}

	return nil
}

func (p *Phone) EncodeICalValue() (string, error) {
	return p.Number, nil
}

func (p *Phone) DecodeICalValue(value string) error {
	p.Number = value
	return nil
}

func (p *Phone) EncodeICalName() (properties.PropertyName, error) {
	return properties.PhonePropertyName, nil
}
