package components

import (
	"fmt"
	"testing"
	"time"

	"github.com/LeonidKrn/caldav-go/icalendar"
	"github.com/LeonidKrn/caldav-go/icalendar/values"
	. "gopkg.in/check.v1"
)

type FreeBusySuite struct{}

var _ = Suite(new(FreeBusySuite))

func TestFreeBusy(t *testing.T) { TestingT(t) }

func (s *FreeBusySuite) TestMissingEndMarshal(c *C) {
	now := time.Now().UTC()
	freeBusy := NewFreeBusy("test", now)
	_, err := icalendar.Marshal(freeBusy)
	c.Assert(err, ErrorMatches, "(?s).*end date or duration must be set.*")
}

func (s *FreeBusySuite) TestBasicWithDurationMarshal(c *C) {
	now := time.Now().UTC()
	freeBusy := NewFreeBusyWithDuration("test", now, time.Hour)
	enc, err := icalendar.Marshal(freeBusy)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VFREEBUSY\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDURATION:PT1H\r\nEND:VFREEBUSY"
	fdate := now.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, fdate, fdate))
}

func (s *FreeBusySuite) TestBasicWithEndMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	freeBusy := NewFreeBusyWithEnd("test", now, end)
	enc, err := icalendar.Marshal(freeBusy)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VFREEBUSY\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nEND:VFREEBUSY"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate))
}

func (s *FreeBusySuite) TestFullFreeBusyMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	freeBusy := NewFreeBusyWithEnd("1:2:3", now, end)
	freeBusy.Attendees = []*values.AttendeeContact{
		values.NewAttendeeContact("Jon Azoff", "jon@dolanor.com"),
		values.NewAttendeeContact("Matthew Davie", "matthew@dolanor.com"),
	}
	freeBusy.Organizer = values.NewOrganizerContact("Jon Azoff", "jon@dolanor.com")
	enc, err := icalendar.Marshal(freeBusy)
	if err != nil {
		c.Fatal(err.Error())
	}
	tmpl := "BEGIN:VFREEBUSY\r\nUID:1:2:3\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\n" +
		"ORGANIZER;CN=\"Jon Azoff\":mailto:jon@dolanor.com\r\n" +
		"ATTENDEE;CN=\"Jon Azoff\":mailto:jon@dolanor.com\r\nATTENDEE;CN=\"Matthew Davie\":mailto:matthew@dolanor.com\r\n" +
		"END:VFREEBUSY"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate))
}

func (s *FreeBusySuite) TestUnmarshalMultipleLines(c *C) {
	// freeBusy that has an ATTENDEE that spans 3 lines
	raw := `BEGIN:VFREEBUSY
DTSTART;TZID=America/Los_Angeles:20150511T140000
DTEND;TZID=America/Los_Angeles:20150511T150000
DTSTAMP:20150511T204516Z
ORGANIZER;CN=Fakebiz Shared:mailto:fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikm
 aeq8@group.calendar.google.com
FREEBUSY:20200214T022955Z/20200214T032955Z
FREEBUSY:20200216T071500Z/20200216T081500Z
UID:na9njgloe10sch3h0uootli104@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;CN=Fakebiz
  Shared;X-NUM-GUESTS=0:mailto:fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikm
 aeq8@group.calendar.google.com
END:VFREEBUSY`

	e := FreeBusy{}
	err := icalendar.Unmarshal(raw, &e)
	c.Assert(err, IsNil)
	c.Assert(len(e.FreeBusyItems), Equals, 2)
	c.Assert(len(e.Attendees), Equals, 1)
	start1 := values.DateTime{}
	start1.DecodeICalValue("20200214T022955Z")
	end1 := values.DateTime{}
	end1.DecodeICalValue("20200214T032955Z")
	c.Assert(e.FreeBusyItems[0].Start, Equals, start1)
	c.Assert(e.FreeBusyItems[0].End, Equals, end1)
	start2 := values.DateTime{}
	start2.DecodeICalValue("20200216T071500Z")
	end2 := values.DateTime{}
	end2.DecodeICalValue("20200216T081500Z")
	c.Assert(e.FreeBusyItems[1].Start, Equals, start2)
	c.Assert(e.FreeBusyItems[1].End, Equals, end2)
	c.Assert(e.Attendees[0].Entry.Address, Equals, "fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikmaeq8@group.calendar.google.com")
	c.Assert(e.Attendees[0].Entry.Name, Equals, "Fakebiz Shared")
}
