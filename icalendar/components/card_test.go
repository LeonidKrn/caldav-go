package components

import (
	"fmt"
	"testing"

	"regexp"

	"strings"

	"github.com/LeonidKrn/caldav-go/icalendar"
	"github.com/LeonidKrn/caldav-go/icalendar/values"
	. "gopkg.in/check.v1"
)

const (
	exampleGroupRaw = `BEGIN:VCARD
VERSION:3.0
X-ADDRESSBOOKSERVER-KIND:group
PRODID:-//Apple Inc.//iOS 10.0.1//EN
N:IMMOMIG
FN:IMMOMIG
UID:D7C763B7-39C2-4330-A0A7-CAD859F8C297
END:VCARD`
)

type CardSuite struct{}

var _ = Suite(new(CardSuite))

func TestCard(t *testing.T) { TestingT(t) }

func (s *CardSuite) TestMarshalCard(c *C) {
	card := NewCard()
	enc, err := icalendar.Marshal(card)
	c.Assert(err, IsNil)
	fmt.Println(enc)
	c.Assert(enc, Equals,
		"BEGIN:VCARD\r\nVERSION:3.0\r\nUID:229CD09F-7FCB-4873-88DC-E16D568D8B50\r\nPRODID:-//iPaladinLLC/caldav-go//NONSGML v1.0.0//EN\r\nN:Doe;Frank;;;\r\nORG:DOE Enterprise;Management\r\nFN:Frank Doe\r\nTEL;TYPE=CELL;TYPE=VOICE;TYPE=pref:111 222 333\r\nTEL;TYPE=WORK;TYPE=VOICE:111 333 444\r\nTEL;TYPE=MAIN:111 444 555\r\nEMAIL;TYPE=WORK;TYPE=INTERNET;TYPE=pref:frank.doe@example.com\r\nEMAIL;TYPE=WORK;TYPE=INTERNET:frank.doo@example.com\r\nEND:VCARD")
}

func (s *CardSuite) TestUnmarshalCard(c *C) {
	raw := `BEGIN:VCARD
VERSION:3.0
UID:2A32F9A6-DCB9-4F7B-85F5-75C1195B1450
N:Doe;Frank;;;
FN:Frank Doe
ORG:COX Enterprise;Management
TITLE:CEO
item1.ADR;TYPE=HOME;TYPE=pref:;;Street 10\nStreet 11;Prague;Prague county;1
1000;Czech Republic
item1.X-ABADR:cz
item2.ADR;TYPE=WORK:;;Street 20;Prague;;11000;Czech Republic
item2.X-ABADR:cz
TEL;TYPE=CELL;TYPE=pref;TYPE=VOICE:111 222 333
TEL;TYPE=MAIN:111 333 444
TEL;TYPE=WORK;TYPE=VOICE:111 444 555
TEL;TYPE=HOME;TYPE=VOICE:111 555 666
EMAIL;TYPE=WORK;TYPE=pref;TYPE=INTERNET:frank.doe@example.com
EMAIL;TYPE=WORK;TYPE=INTERNET:frank.duo@example.com
PRODID:-//Apple Inc.//iCloud Web Address Book 17A37//EN
REV:2017-02-01T02:51:01Z
END:VCARD`

	card := Card{}
	err := icalendar.Unmarshal(raw, &card)
	c.Assert(err, IsNil)
	c.Assert(len(card.Emails), Equals, 2)
	c.Assert(card.Organization, NotNil)
	c.Assert(card.Organization.Company, Equals, "COX Enterprise")
}

func NewCard() Card {
	return Card{
		UID:          "229CD09F-7FCB-4873-88DC-E16D568D8B50",
		Name:         values.NewContactName("Frank", "Doe", "", "", ""),
		DisplayName:  "Frank Doe",
		Organization: values.NewOrganization("DOE Enterprise", "Management"),
		Phones: []*values.Phone{
			values.NewPhone("111 222 333", true, "CELL", "VOICE"),
			values.NewPhone("111 333 444", false, "WORK", "VOICE"),
			values.NewPhone("111 444 555", false, "MAIN"),
		},
		Emails: []*values.Email{
			values.NewEmail("frank.doe@example.com", true, "WORK", "INTERNET"),
			values.NewEmail("frank.doo@example.com", false, "WORK", "INTERNET"),
		},
	}
}

func (s *CardSuite) TestUnmarshalGroup(c *C) {
	card := Card{}
	err := icalendar.Unmarshal(exampleGroupRaw, &card)
	c.Assert(err, IsNil)
	c.Assert(card.IsGroup(), Equals, true)
}

func (s *CardSuite) TestUnmarshalAnotherCards(c *C) {
	example := `
BEGIN:VCARD
VERSION:3.0
UID:123412
PRODID:-//Apple Inc.//iPhone OS 10.2//EN
N:Appleseed;John;;;
FN:John  Appleseed
TEL;type=IPHONE;type=CELL;type=VOICE;type=pref:(408) 555-0126
item1.EMAIL;type=INTERNET;type=pref:john@example.com
item1.X-ABLabel:Home Email
item2.EMAIL;type=INTERNET:j.appleseed@icloud.com
item2.X-ABLabel:Work Email
END:VCARD
BEGIN:VCARD
VERSION:3.0
UID:5673128
PRODID:-//Apple Inc.//iPhone OS 10.2//EN
N:Second;Dude;;;
FN:Dude  Second
TEL;type=IPHONE;type=CELL;type=VOICE;type=pref:(212) 333-4455
item1.EMAIL;type=INTERNET;type=pref:sd@example.com
item1.X-ABLabel:Home Email
item2.EMAIL;type=INTERNET:s.dude@icloud.com
item2.X-ABLabel:Work Email
END:VCARD`

	var contacts []Card
	err := icalendar.Unmarshal(example, &contacts)
	c.Assert(err, IsNil)
	c.Assert(len(contacts), Equals, 2)
	c.Assert(len(contacts[0].Emails), Equals, 2)
	c.Assert(contacts[0].Emails[0].IsPreferred, Equals, true)
	c.Assert(contacts[0].Emails[0].Mail, Equals, "john@example.com")
	c.Assert(contacts[0].Emails[0].Label, Equals, "Home Email")
	c.Assert(len(contacts[1].Emails), Equals, 2)
	c.Assert(contacts[1].Emails[1].IsPreferred, Equals, false)
	c.Assert(contacts[1].Emails[1].Mail, Equals, "s.dude@icloud.com")
	c.Assert(contacts[1].Emails[1].Label, Equals, "Work Email")

	cc, err := icalendar.Marshal(contacts)
	c.Assert(err, IsNil)
	c.Assert(strings.ToLower(strings.Replace(cc, "\r\n", "\n", -1)), Equals, strings.ToLower(strings.TrimSpace(example)))
}

func (s *CardSuite) TestUnmarshalAnother2Cards(c *C) {
	example := `
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 10.2//EN
N:Appleseed;John;;;
FN:John  Appleseed
item1.EMAIL;type=INTERNET;type=pref:john@example.com
item1.X-ABLabel:Home Email
item2.EMAIL;type=INTERNET:j.appleseed@icloud.com
item2.X-ABLabel:Work Email
TEL;type=IPHONE;type=CELL;type=VOICE;type=pref:(408) 555-0126
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 10.2//EN
N:Second;Dude;;;
FN: Dude  Second
item1.EMAIL;type=INTERNET;type=pref:sd@example.com
item1.X-ABLabel:Home Email
item2.EMAIL;type=INTERNET:s.dude@icloud.com
item2.X-ABLabel:Work Email
TEL;type=IPHONE;type=CELL;type=VOICE;type=pref:(212) 333-4455
END:VCARD`

	rr := regexp.MustCompile(`(BEGIN:VCARD(?:.|\n)*?END:VCARD)`)
	sub := rr.FindAllStringSubmatch(example, -1)
	for _, ss := range sub {
		var contacts Card
		err := icalendar.Unmarshal(ss[1], &contacts)
		c.Assert(err, IsNil)
		c.Log(contacts.Emails)
	}

}

func (s *CardSuite) TestLargeAmountCars(c *C) {
	example := `BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:bb;aa;;;
FN: aa  bb
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:Contact;Contact;;;
FN: Contact  Contact
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:Contact;Test iCloudd;;;
FN: Test iCloudd  Contact
TEL;type=CELL;type=VOICE;type=pref:24 0220 17
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:dd;cc;;;
FN: cc  dd
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:Doe;Carl;;;
FN: Carl  Doe
ORG:Doe Enterprises;
EMAIL;type=INTERNET;type=HOME;type=pref:carl.doe@example.com
TEL;type=CELL;type=VOICE;type=pref:+421121244555
TEL;type=HOME;type=VOICE:+420 603 287 934
TEL;type=WORK;type=VOICE:603 823 444
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:Dude;Some;;;
FN: Some  Dude
ORG:Dude's Company;
EMAIL;type=INTERNET;type=HOME;type=pref:dude@mail.com
EMAIL;type=INTERNET;type=WORK:some@icloud.com
TEL;type=HOME;type=VOICE;type=pref:(313) 333-3333
TEL;type=CELL;type=VOICE:(414) 444-4444
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:ff;ee;;;
FN: ee  ff
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:Group;Account;;;
FN: Account  Group
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:hh;gg;;;
FN: gg  hh
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:IMLast;IMFirst;;;
FN: IMFirst  IMLast
ORG:Immomig;
EMAIL;type=INTERNET;type=HOME;type=pref:immomig@mail.com
TEL;type=HOME;type=VOICE;type=pref:(212) 111-1111
TEL;type=CELL;type=VOICE:(212) 222-2222
ADR;type=HOME;type=pref:;;215 Broaway Ave\nApt 25;Manhattan;NY;11111;United States
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:marian 99;test;;;
FN: test  marian 99
TEL;type=HOME;type=VOICE;type=pref:776427542
END:VCARD
BEGIN:VCARD
VERSION:3.0
PRODID:-//Apple Inc.//iPhone OS 9.3.5//EN
N:14000;test name;;;
FN: test name  14000
EMAIL;type=INTERNET;type=HOME;type=pref:temail@realpad.eu
TEL;type=HOME;type=VOICE;type=pref:123456
END:VCARD`

	var contacts []Card
	err := icalendar.Unmarshal(example, &contacts)
	c.Assert(err, IsNil)
	c.Log(len(contacts))
	for _, cnt := range contacts {
		c.Log(cnt.DisplayName)
	}
}
