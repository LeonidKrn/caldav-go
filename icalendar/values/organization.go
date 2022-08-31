package values

import (
	"fmt"
	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar/properties"
)

type Organization struct {
	Company    string
	Department string
}

func NewOrganization(companyName string, department string) *Organization {
	return &Organization{
		Company:    companyName,
		Department: department,
	}
}

func (o *Organization) EncodeICalValue() (string, error) {
	return fmt.Sprintf("%s;%s", o.Company, o.Department), nil
}

func (o *Organization) DecodeICalValue(value string) error {
	parts := strings.Split(value, ";")
	o.Company = parts[0]
	if len(parts) > 2 {
		o.Department = parts[1]
	}

	return nil
}

func (o *Organization) EncodeICalName() (properties.PropertyName, error) {
	return properties.OrganizationPropertyName, nil
}
