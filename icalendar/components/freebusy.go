package components

import (
	"time"

	"github.com/LeonidKrn/caldav-go/icalendar/values"
	"github.com/LeonidKrn/caldav-go/utils"
)

type FreeBusy struct {

	// defines the persistent, globally unique identifier for the calendar component.
	UID string `ical:",required"`

	// indicates the date/time that the instance of the iCalendar object was created.
	DateStamp *values.DateTime `ical:"dtstamp,required"`

	// specifies when the calendar component begins.
	DateStart *values.DateTime `ical:"dtstart,required"`

	// specifies the date and time that a calendar component ends.
	DateEnd *values.DateTime `ical:"dtend,omitempty"`

	// specifies a positive duration of time.
	Duration *values.Duration `ical:",omitempty"`

	// defines the organizer for a calendar component.
	Organizer *values.OrganizerContact `ical:",omitempty"`

	// defines an "Attendee" within a calendar component.
	Attendees []*values.AttendeeContact `ical:"attendee,omitempty"`

	// free busy entries
	FreeBusyItems []*values.FreeBusyItem `ical:"freebusy,omitempty"`
}

// validates the FreeBusy internals
func (e *FreeBusy) ValidateICalValue() error {

	if e.UID == "" {
		return utils.NewError(e.ValidateICalValue, "the UID value must be set", e, nil)
	}

	if e.DateStart == nil {
		return utils.NewError(e.ValidateICalValue, "event start date must be set", e, nil)
	}

	if e.DateEnd == nil && e.Duration == nil {
		return utils.NewError(e.ValidateICalValue, "event end date or duration must be set", e, nil)
	}

	if e.DateEnd != nil && e.Duration != nil {
		return utils.NewError(e.ValidateICalValue, "event end date and duration are mutually exclusive fields", e, nil)
	}

	return nil
}

// creates a new iCalendar event with no end time
func NewFreeBusy(uid string, start time.Time) *FreeBusy {
	e := new(FreeBusy)
	e.UID = uid
	e.DateStamp = values.NewDateTime(time.Now().UTC())
	e.DateStart = values.NewDateTime(start)
	return e
}

// creates a new iCalendar event that lasts a certain duration
func NewFreeBusyWithDuration(uid string, start time.Time, duration time.Duration) *FreeBusy {
	e := NewFreeBusy(uid, start)
	e.Duration = values.NewDuration(duration)
	return e
}

// creates a new iCalendar event that has an explicit start and end time
func NewFreeBusyWithEnd(uid string, start time.Time, end time.Time) *FreeBusy {
	e := NewFreeBusy(uid, start)
	e.DateEnd = values.NewDateTime(end)
	return e
}
