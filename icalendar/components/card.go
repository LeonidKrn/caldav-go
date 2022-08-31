package components

import (
	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar/values"
)

const (
	groupKind = "group"
)

type Card struct {
	Version string `ical:",3.0"`

	UID string `ical:",required"`

	ProductId string `ical:"prodid,-//iPaladinLLC/caldav-go//NONSGML v1.0.0//EN"`

	Name *values.ContactName `ical:"n,omitempty"`

	Organization *values.Organization `ical:"org,omitempty"`

	DisplayName string `ical:"fn,omitempty"`

	AddressBookKind string `ical:"x_addressbookserver_kind,omitempty"`

	AddressBookMembers []*values.AddressBookMember `ical:"x_addressbookserver_member,omitempty"`

	Categories string `ical:"categories,omitempty"`

	Phones []*values.Phone `ical:"tel,omitempty"`

	Emails []*values.Email `ical:"email,omitempty"`
}

func (c Card) IsGroup() bool {
	return strings.EqualFold(c.AddressBookKind, "group")
}

func (c *Card) AddAddressBookMember(m ...*values.AddressBookMember) {
	c.AddressBookMembers = append(c.AddressBookMembers, m...)
}

func NewCardGroup(uid string, name string) *Card {
	n := values.NewSimpleContactName(name)

	return &Card{
		UID:             uid,
		Name:            n,
		DisplayName:     n.GetDisplayName(),
		AddressBookKind: groupKind,
	}
}
