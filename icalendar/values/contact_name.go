package values

import (
	"fmt"
	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar/properties"
)

type ContactName struct {
	SimpleName                                      string
	FirstName, LastName, MiddleName, Prefix, Suffix string
}

func NewContactName(firstName, lastName, middleName, prefix, suffix string) *ContactName {
	return &ContactName{
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
		Prefix:     prefix,
		Suffix:     suffix,
	}
}

func NewSimpleContactName(simpleName string) *ContactName {
	return &ContactName{
		SimpleName: simpleName,
	}
}

func (c *ContactName) GetDisplayName() string {
	if c.SimpleName != "" {
		return c.SimpleName
	}

	var nameParts []string
	if c.Prefix != "" {
		nameParts = append(nameParts, c.Prefix)
	}

	if c.FirstName != "" {
		nameParts = append(nameParts, c.FirstName)
	}

	if c.LastName != "" {
		nameParts = append(nameParts, c.LastName)
	}

	if c.Suffix != "" {
		nameParts = append(nameParts, c.Suffix)
	}

	return strings.Join(nameParts, " ")
}

func (c *ContactName) EncodeICalValue() (string, error) {
	if c.SimpleName != "" {
		return c.SimpleName, nil
	} else {
		return fmt.Sprintf("%s;%s;%s;%s;%s", c.LastName, c.FirstName, c.MiddleName, c.Prefix, c.Suffix), nil
	}

}

func (c *ContactName) DecodeICalValue(value string) error {
	parts := strings.Split(value, ";")

	if len(parts) == 1 {
		c.SimpleName = value
		return nil
	}

	c.LastName = parts[0]

	if len(parts) > 1 {
		c.FirstName = parts[1]
	}

	if len(parts) > 2 {
		c.MiddleName = parts[2]
	}

	if len(parts) > 3 {
		c.Prefix = parts[3]
	}

	if len(parts) > 4 {
		c.Suffix = parts[4]
	}

	return nil
}

func (c *ContactName) EncodeICalName() (properties.PropertyName, error) {
	return properties.NamePropertyName, nil
}
