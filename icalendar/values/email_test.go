package values

import (
	"testing"

	"github.com/LeonidKrn/caldav-go/icalendar"
	. "gopkg.in/check.v1"
)

type EmailSuite struct{}

type emailTestObj struct {
	*Email
}

var _ = Suite(new(EmailSuite))

func TestEmail(t *testing.T) { TestingT(t) }

func (s *EmailSuite) TestEncode(c *C) {
	eto := new(emailTestObj)
	eto.Email = &Email{
		Mail: "frank.doe@example.com",
		Types: []string{
			"WORK",
			"pref",
			"INTERNET",
		},
	}

	encoded, err := icalendar.Marshal(eto)
	c.Assert(err, IsNil)
	expected := "EMAIL;TYPE=WORK;TYPE=pref;TYPE=INTERNET:frank.doe@example.com"
	c.Assert(encoded, Equals, expected)
}
