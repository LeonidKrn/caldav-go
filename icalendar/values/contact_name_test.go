package values

import (
	"testing"

	"github.com/LeonidKrn/caldav-go/icalendar"
	. "gopkg.in/check.v1"
)

type ContactNameSuite struct{}

var _ = Suite(new(ContactNameSuite))

func TestContactName(t *testing.T) { TestingT(t) }

func (s *ContactNameSuite) TestMarshalName(c *C) {
	n := &ContactName{
		FirstName:  "Frank",
		LastName:   "Doe",
		MiddleName: "Francis",
		Prefix:     "Mr.",
		Suffix:     "jr.",
	}

	enc, err := icalendar.Marshal(n)
	c.Assert(err, IsNil)
	c.Assert(enc, Equals, "N:Doe;Frank;Francis;Mr.;jr.")
}
