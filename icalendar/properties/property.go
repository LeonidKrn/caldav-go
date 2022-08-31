package properties

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/LeonidKrn/caldav-go/utils"
)

var _ = log.Print

var propNameSanitizer = strings.NewReplacer(
	"_", "-",
	":", "\\:",
)

var propValueSanitizer = strings.NewReplacer(
	"\"", "'",
	"\\", "\\\\",
	"\n", "\\n",
)

var propNameDesanitizer = strings.NewReplacer(
	"-", "_",
	"\\:", ":",
)

var propValueDesanitizer = strings.NewReplacer(
	"'", "\"",
	"\\\\", "\\",
	"\\n", "\n",
)

type Property struct {
	Name                PropertyName
	Value, DefaultValue string
	Params              Params
	OmitEmpty, Required bool
	Prefix              string
}

func (p *Property) HasNameAndValue() bool {
	return p.Name != "" && p.Value != ""
}

func (p *Property) Merge(override *Property) {
	if override.Name != "" {
		p.Name = override.Name
	}
	if override.Value != "" {
		p.Value = override.Value
	}
	if override.Params != nil {
		p.Params = override.Params
	}
}

func PropertyFromStructField(fs reflect.StructField) (p *Property) {

	ftag := fs.Tag.Get("ical")
	if fs.PkgPath != "" || ftag == "-" {
		return
	}

	p = new(Property)

	// parse the field tag
	if ftag != "" {
		tags := strings.Split(ftag, ",")
		p.Name = PropertyName(tags[0])
		if len(tags) > 1 {
			if tags[1] == "omitempty" {
				p.OmitEmpty = true
			} else if tags[1] == "required" {
				p.Required = true
			} else {
				p.DefaultValue = tags[1]
			}
		}
	}

	// make sure we have a name
	if p.Name == "" {
		p.Name = PropertyName(fs.Name)
	}

	p.Name = PropertyName(strings.ToUpper(string(p.Name)))

	return

}

func MarshalProperty(p *Property) string {
	name := strings.ToUpper(propNameSanitizer.Replace(string(p.Name)))
	value := propValueSanitizer.Replace(p.Value)
	keys := []string{name}
	for _, param := range p.Params {
		name := ParameterName(strings.ToUpper(propNameSanitizer.Replace(string(param.Name))))
		value := propValueSanitizer.Replace(param.Value)
		if strings.ContainsAny(value, " :") {
			keys = append(keys, fmt.Sprintf("%s=\"%s\"", name, value))
		} else {
			keys = append(keys, fmt.Sprintf("%s=%s", name, value))
		}
	}

	name = strings.Join(keys, ";")
	if p.Prefix != "" {
		name = fmt.Sprintf("%s.%s", p.Prefix, name)
	}

	return fmt.Sprintf("%s:%s", name, value)
}

func PropertyFromInterface(target interface{}) (p *Property, adds []*Property, err error) {

	var ierr error
	if va, ok := target.(CanValidateValue); ok {
		if ierr = va.ValidateICalValue(); ierr != nil {
			err = utils.NewError(PropertyFromInterface, "interface failed validation", target, ierr)
			return
		}
	}

	p = new(Property)

	if enc, ok := target.(CanEncodeName); ok {
		if p.Name, ierr = enc.EncodeICalName(); ierr != nil {
			err = utils.NewError(PropertyFromInterface, "interface failed name encoding", target, ierr)
			return
		}
	}

	if enc, ok := target.(CanEncodeParams); ok {
		if p.Params, ierr = enc.EncodeICalParams(); ierr != nil {
			err = utils.NewError(PropertyFromInterface, "interface failed params encoding", target, ierr)
			return
		}
	}

	if enc, ok := target.(CanEncodeValue); ok {
		if p.Value, ierr = enc.EncodeICalValue(); ierr != nil {
			err = utils.NewError(PropertyFromInterface, "interface failed value encoding", target, ierr)
			return
		}
	}

	if enc, ok := target.(CanEncodeAdditionalProperties); ok {
		if adds, ierr = enc.EncodeAdditionalICalProperties(); ierr != nil {
			err = utils.NewError(PropertyFromInterface, "interface failed additional values encoding", target, ierr)
			return
		}
	}

	return

}

func UnmarshalProperty(line string) *Property {
	nvp := strings.SplitN(line, ":", 2)
	prop := new(Property)
	if len(nvp) > 1 {
		prop.Value = strings.TrimSpace(nvp[1])
	}
	npp := strings.Split(nvp[0], ";")
	if len(npp) > 1 {
		prop.Params = make(Params, 0)
		for i := 1; i < len(npp); i++ {
			var key, value string
			kvp := strings.Split(npp[i], "=")
			key = strings.TrimSpace(kvp[0])
			key = propNameDesanitizer.Replace(key)
			if len(kvp) > 1 {
				value = strings.TrimSpace(kvp[1])
				value = propValueDesanitizer.Replace(value)
				value = strings.Trim(value, "\"")
			}
			prop.Params = append(prop.Params, Param{
				Name:  ParameterName(key),
				Value: value,
			})
		}
	}
	prop.Name = PropertyName(strings.TrimSpace(npp[0]))
	prop.Name = PropertyName(propNameDesanitizer.Replace(string(prop.Name)))
	prop.Name, prop.Prefix = splitPropertyName(prop.Name)
	prop.Value = propValueDesanitizer.Replace(prop.Value)
	return prop
}

func NewProperty(name, value string) *Property {
	return &Property{Name: PropertyName(name), Value: value}
}

func splitPropertyName(name PropertyName) (nameOut PropertyName, prefix string) {
	strs := strings.Split(string(name), ".")
	nameOut = PropertyName(strs[len(strs)-1])
	if len(strs) > 1 {
		prefix = strs[0]
	}

	return
}
