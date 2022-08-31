package entities

import (
	"encoding/xml"
	"time"

	"github.com/LeonidKrn/caldav-go/caldav/values"
	"github.com/LeonidKrn/caldav-go/utils"
	"github.com/LeonidKrn/caldav-go/webdav/entities"
)

// a CalDAV calendar query object
type CalendarQuery struct {
	XMLName xml.Name          `xml:"urn:ietf:params:xml:ns:caldav calendar-query"`
	Prop    *Prop             `xml:",omitempty"`
	AllProp *entities.AllProp `xml:",omitempty"`
	Filter  *Filter           `xml:",omitempty"`
}

func NewEventQuery() *CalendarQuery {
	// construct the query object
	query := new(CalendarQuery)

	// request all calendar data
	query.Prop = new(Prop)
	query.Prop.CalendarData = new(CalendarData)

	// filter down calendar data to only iCalendar data
	query.Filter = new(Filter)
	query.Filter.ComponentFilter = new(ComponentFilter)
	query.Filter.ComponentFilter.Name = values.CalendarComponentName

	// filter down iCalendar data to only events
	query.Filter.ComponentFilter.ComponentFilter = new(ComponentFilter)
	query.Filter.ComponentFilter.ComponentFilter.Name = values.EventComponentName

	return query
}

// creates a new CalDAV query for iCalendar events from a particular time range
func NewEventRangeQuery(start, end time.Time, singleEvents bool) (*CalendarQuery, error) {

	var err error
	var dtstart, dtend *values.DateTime
	if dtstart, err = values.NewDateTime("start", start); err != nil {
		return nil, utils.NewError(NewEventRangeQuery, "unable to encode start time", start, err)
	} else if dtend, err = values.NewDateTime("end", end); err != nil {
		return nil, utils.NewError(NewEventRangeQuery, "unable to encode end time", end, err)
	}

	// construct the query object
	query := NewEventQuery()

	// expand recurring events
	if singleEvents {
		query.Prop.CalendarData.ExpandRecurrenceSet = new(ExpandRecurrenceSet)
		query.Prop.CalendarData.ExpandRecurrenceSet.StartTime = dtstart
		query.Prop.CalendarData.ExpandRecurrenceSet.EndTime = dtend
	}

	// filter down the events to only those that fall within the time range
	query.Filter.ComponentFilter.ComponentFilter.TimeRange = new(TimeRange)
	query.Filter.ComponentFilter.ComponentFilter.TimeRange.StartTime = dtstart
	query.Filter.ComponentFilter.ComponentFilter.TimeRange.EndTime = dtend

	// return the event query
	return query, nil

}

func NewSimpleEventRangeQuery(start, end time.Time) (*CalendarQuery, error) {
	var err error
	var dtstart, dtend *values.DateTime
	if dtstart, err = values.NewDateTime("start", start); err != nil {
		return nil, utils.NewError(NewEventRangeQuery, "unable to encode start time", start, err)
	} else if dtend, err = values.NewDateTime("end", end); err != nil {
		return nil, utils.NewError(NewEventRangeQuery, "unable to encode end time", end, err)
	}

	// construct the query object
	query := NewEventQuery()

	// filter down the events to only those that fall within the time range
	query.Filter.ComponentFilter.ComponentFilter.TimeRange = new(TimeRange)
	query.Filter.ComponentFilter.ComponentFilter.TimeRange.StartTime = dtstart
	query.Filter.ComponentFilter.ComponentFilter.TimeRange.EndTime = dtend

	// return the event query
	return query, nil
}
