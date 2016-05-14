package values

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/taviti/caldav-go/icalendar/properties"
	"github.com/taviti/caldav-go/utils"
)

var _ = log.Print

const DateFormatString = "20060102"
const DateTimeFormatString = "20060102T150405"
const UTCDateTimeFormatString = "20060102T150405Z"

// a representation of a date and time for iCalendar
type DateTime struct {
	t time.Time
}

type DateTimeFullDay DateTime

type DateTimes []*DateTime

// The exception dates, if specified, are used in computing the recurrence set. The recurrence set is the complete set
// of recurrence instances for a calendar component. The recurrence set is generated by considering the initial
// "DTSTART" property along with the "RRULE", "RDATE", "EXDATE" and "EXRULE" properties contained within the iCalendar
// object. The "DTSTART" property defines the first instance in the recurrence set. Multiple instances of the "RRULE"
// and "EXRULE" properties can also be specified to define more sophisticated recurrence sets. The final recurrence set
// is generated by gathering all of the start date-times generated by any of the specified "RRULE" and "RDATE"
// properties, and then excluding any start date and times which fall within the union of start date and times
// generated by any specified "EXRULE" and "EXDATE" properties. This implies that start date and times within exclusion
// related properties (i.e., "EXDATE" and "EXRULE") take precedence over those specified by inclusion properties
// (i.e., "RDATE" and "RRULE"). Where duplicate instances are generated by the "RRULE" and "RDATE" properties, only
// one recurrence is considered. Duplicate instances are ignored.
//
// The "EXDATE" property can be used to exclude the value specified in "DTSTART". However, in such cases the original
// "DTSTART" date MUST still be maintained by the calendaring and scheduling system because the original "DTSTART"
// value has inherent usage dependencies by other properties such as the "RECURRENCE-ID".
type ExceptionDateTimes DateTimes

// The recurrence dates, if specified, are used in computing the recurrence set. The recurrence set is the complete set
// of recurrence instances for a calendar component. The recurrence set is generated by considering the initial
// "DTSTART" property along with the "RRULE", "RDATE", "EXDATE" and "EXRULE" properties contained within the iCalendar
// object. The "DTSTART" property defines the first instance in the recurrence set. Multiple instances of the "RRULE"
// and "EXRULE" properties can also be specified to define more sophisticated recurrence sets. The final recurrence set
// is generated by gathering all of the start date-times generated by any of the specified "RRULE" and "RDATE"
// properties, and then excluding any start date and times which fall within the union of start date and times
// generated by any specified "EXRULE" and "EXDATE" properties. This implies that start date and times within exclusion
// related properties (i.e., "EXDATE" and "EXRULE") take precedence over those specified by inclusion properties
// (i.e., "RDATE" and "RRULE"). Where duplicate instances are generated by the "RRULE" and "RDATE" properties, only
// one recurrence is considered. Duplicate instances are ignored.
type RecurrenceDateTimes DateTimes

// creates a new icalendar datetime representation
func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t: t.Truncate(time.Second)}
}

// creates a new icalendar datetime array representation
func NewDateTimes(dates ...*DateTime) DateTimes {
	return DateTimes(dates)
}

// creates a new icalendar datetime representation
func NewDateTimeFullDay(t time.Time) *DateTimeFullDay {
	datetimes := NewDateTime(t)
	return (*DateTimeFullDay)(datetimes)
}

// creates a new icalendar datetime array representation
func NewExceptionDateTimes(dates ...*DateTime) *ExceptionDateTimes {
	datetimes := NewDateTimes(dates...)
	return (*ExceptionDateTimes)(&datetimes)
}

// creates a new icalendar datetime array representation
func NewRecurrenceDateTimes(dates ...*DateTime) *RecurrenceDateTimes {
	datetimes := NewDateTimes(dates...)
	return (*RecurrenceDateTimes)(&datetimes)
}

// checks to see if two datetimes are equal
func (d *DateTime) Equals(test *DateTime) bool {
	return d.t.Equal(test.t)
}

// returns the native time for the datetime object
func (d *DateTime) NativeTime() time.Time {
	return d.t
}

// encodes the datetime value for the iCalendar specification
func (d *DateTime) EncodeICalValue() (string, error) {
	val := d.t.Format(DateTimeFormatString)
	loc := d.t.Location()
	if loc == time.UTC {
		val = fmt.Sprintf("%sZ", val)
	}
	return val, nil
}

// decodes the datetime value from the iCalendar specification
func (d *DateTime) DecodeICalValue(value string) error {
	layout := DateTimeFormatString
	if strings.HasSuffix(value, "Z") {
		layout = UTCDateTimeFormatString
	} else if len(value) == 8 {
		layout = DateFormatString
	}
	var err error
	d.t, err = time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		return utils.NewError(d.DecodeICalValue, "unable to parse datetime value", d, err)
	} else {
		return nil
	}
}

// encodes the datetime params for the iCalendar specification
func (d *DateTime) EncodeICalParams() (params properties.Params, err error) {
	loc := d.t.Location()
	if loc != time.UTC {
		params = properties.Params{properties.TimeZoneIdPropertyName: loc.String()}
	}
	return
}

// decodes the datetime params from the iCalendar specification
func (d *DateTime) DecodeICalParams(params properties.Params) error {
	layout := DateTimeFormatString
	value := d.t.Format(layout)
	if name, found := params[properties.TimeZoneIdPropertyName]; !found {
		return nil
	} else if loc, err := time.LoadLocation(name); err != nil {
		return utils.NewError(d.DecodeICalValue, "unable to parse timezone", d, err)
	} else if t, err := time.ParseInLocation(layout, value, loc); err != nil {
		return utils.NewError(d.DecodeICalValue, "unable to parse datetime value", d, err)
	} else {
		d.t = t
		return nil
	}
}

// validates the datetime value against the iCalendar specification
func (d *DateTime) ValidateICalValue() error {

	loc := d.t.Location()

	if loc == time.Local {
		msg := "DateTime location may not Local, please use UTC or explicit Location"
		return utils.NewError(d.ValidateICalValue, msg, d, nil)
	}

	if loc.String() == "" {
		msg := "DateTime location must have a valid name"
		return utils.NewError(d.ValidateICalValue, msg, d, nil)
	}

	return nil
}

// encodes the datetime value for the iCalendar specification
func (d *DateTime) String() string {
	if s, err := d.EncodeICalValue(); err != nil {
		panic(err)
	} else {
		return s
	}
}

// encodes the datetime for full day value for the iCalendar specification
func (d *DateTimeFullDay) EncodeICalValue() (string, error) {
	val := d.t.Format(DateFormatString)
	loc := d.t.Location()
	if loc == time.UTC {
		val = fmt.Sprintf("%sZ", val)
	}
	return val, nil
}

// encodes a list of datetime values for the iCalendar specification
func (ds *DateTimes) EncodeICalValue() (string, error) {
	var csv CSV
	for i, d := range *ds {
		if s, err := d.EncodeICalValue(); err != nil {
			msg := fmt.Sprintf("unable to encode datetime at index %d", i)
			return "", utils.NewError(ds.EncodeICalValue, msg, ds, err)
		} else {
			csv = append(csv, s)
		}
	}
	return csv.EncodeICalValue()
}

// encodes a list of datetime params for the iCalendar specification
func (ds *DateTimes) EncodeICalParams() (params properties.Params, err error) {
	if len(*ds) > 0 {
		params, err = (*ds)[0].EncodeICalParams()
	}
	return
}

// decodes a list of datetime params from the iCalendar specification
func (ds *DateTimes) DecodeICalParams(params properties.Params) error {
	for i, d := range *ds {
		if err := d.DecodeICalParams(params); err != nil {
			msg := fmt.Sprintf("unable to decode datetime params for index %d", i)
			return utils.NewError(ds.DecodeICalValue, msg, ds, err)
		}
	}
	return nil
}

// encodes a list of datetime values for the iCalendar specification
func (ds *DateTimes) DecodeICalValue(value string) error {
	csv := new(CSV)
	if err := csv.DecodeICalValue(value); err != nil {
		return utils.NewError(ds.DecodeICalValue, "unable to decode datetime list as CSV", ds, err)
	}
	for i, value := range *csv {
		d := new(DateTime)
		if err := d.DecodeICalValue(value); err != nil {
			msg := fmt.Sprintf("unable to decode datetime at index %d", i)
			return utils.NewError(ds.DecodeICalValue, msg, ds, err)
		} else {
			*ds = append(*ds, d)
		}
	}
	return nil
}

// encodes exception date times property name for icalendar
func (e *ExceptionDateTimes) EncodeICalName() (properties.PropertyName, error) {
	return properties.ExceptionDateTimesPropertyName, nil
}

// encodes recurrence date times property name for icalendar
func (r *RecurrenceDateTimes) EncodeICalName() (properties.PropertyName, error) {
	return properties.RecurrenceDateTimesPropertyName, nil
}

// encodes exception date times property value for icalendar
func (e *ExceptionDateTimes) EncodeICalValue() (string, error) {
	return (*DateTimes)(e).EncodeICalValue()
}

// encodes recurrence date times property value for icalendar
func (r *RecurrenceDateTimes) EncodeICalValue() (string, error) {
	return (*DateTimes)(r).EncodeICalValue()
}

// decodes exception date times property value for icalendar
func (e *ExceptionDateTimes) DecodeICalValue(value string) error {
	return (*DateTimes)(e).DecodeICalValue(value)
}

// decodes recurrence date times property value for icalendar
func (r *RecurrenceDateTimes) DecodeICalValue(value string) error {
	return (*DateTimes)(r).DecodeICalValue(value)
}

// encodes exception date times property params for icalendar
func (e *ExceptionDateTimes) EncodeICalParams() (params properties.Params, err error) {
	return (*DateTimes)(e).EncodeICalParams()
}

// encodes recurrence date times property params for icalendar
func (r *RecurrenceDateTimes) EncodeICalParams() (params properties.Params, err error) {
	return (*DateTimes)(r).EncodeICalParams()
}

// encodes exception date times property params for icalendar
func (e *ExceptionDateTimes) DecodeICalParams(params properties.Params) error {
	return (*DateTimes)(e).DecodeICalParams(params)
}

// encodes recurrence date times property params for icalendar
func (r *RecurrenceDateTimes) DecodeICalParams(params properties.Params) error {
	return (*DateTimes)(r).DecodeICalParams(params)
}
