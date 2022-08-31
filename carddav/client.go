package carddav

import (
	"fmt"
	"log"
	"net/http"

	cont "github.com/LeonidKrn/caldav-go/carddav/entities"
	"github.com/LeonidKrn/caldav-go/icalendar/components"
	"github.com/LeonidKrn/caldav-go/utils"
	"github.com/LeonidKrn/caldav-go/webdav"
	"github.com/LeonidKrn/caldav-go/webdav/entities"
)

var _ = log.Print

// a client for making WebDAV requests
type Client webdav.Client

// downcasts the client to the WebDAV interface
func (c *Client) WebDAV() *webdav.Client {
	return (*webdav.Client)(c)
}

// returns the embedded CalDAV server reference
func (c *Client) Server() *Server {
	return (*Server)(c.WebDAV().Server())
}

// creates a new client for communicating with an WebDAV server
func NewClient(server *Server, native *http.Client) *Client {
	return (*Client)(webdav.NewClient((*webdav.Server)(server), native))
}

// creates a new client for communicating with a WebDAV server
// uses the default HTTP client from net/http
func NewDefaultClient(server *Server) *Client {
	return NewClient(server, http.DefaultClient)
}

// executes a CardDAV request
func (c *Client) Do(req *Request) (*Response, error) {
	if resp, err := c.WebDAV().Do((*webdav.Request)(req)); err != nil {
		return nil, utils.NewError(c.Do, "unable to execute CardDAV request", c, err)
	} else {
		return NewResponse(resp), nil
	}
}

// attempts to fetch an cards on the remote CalDAV server
func (c *Client) QueryCards(path string, query *cont.ContactQuery) (contacts []*components.ContactCard, oerr error) {
	ms := new(cont.Multistatus)
	if req, err := c.Server().WebDAV().NewRequest("REPORT", path, query); err != nil {
		oerr = utils.NewError(c.QueryCards, "unable to create request", c, err)
	} else if resp, err := c.WebDAV().Do(req); err != nil {
		oerr = utils.NewError(c.QueryCards, "unable to execute request", c, err)
	} else if resp.StatusCode == http.StatusNotFound {
		return // no events if not found
	} else if resp.StatusCode != webdav.StatusMulti {
		err := new(entities.Error)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		resp.Decode(err)
		oerr = utils.NewError(c.QueryCards, msg, c, err)
	} else if err := resp.Decode(ms); err != nil {
		msg := "unable to decode response"
		oerr = utils.NewError(c.QueryCards, msg, c, err)
	} else {
		for i, r := range ms.Responses {
			for j, p := range r.PropStats {
				if p.Prop == nil || p.Prop.AddressData == nil {
					continue
				} else if card, err := p.Prop.AddressData.Card(); err != nil {
					msg := fmt.Sprintf("unable to decode property %d of response %d", j, i)
					oerr = utils.NewError(c.QueryCards, msg, c, err)
					return
				} else {
					contacts = append(contacts, &components.ContactCard{Card: *card, Href: r.Href})
				}
			}
		}
	}
	return
}

// attempts to fetch an event on the remote CardDAV server
func (c *Client) GetCard(path string) (*components.ContactCard, error) {
	var crd components.Card
	if req, err := c.Server().NewRequest("GET", path); err != nil {
		return nil, utils.NewError(c.GetCard, "unable to create request", c, err)
	} else if resp, err := c.Do(req); err != nil {
		return nil, utils.NewError(c.GetCard, "unable to execute request", c, err)
	} else if resp.StatusCode != http.StatusOK {
		err := new(entities.Error)
		resp.WebDAV().Decode(err)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		return nil, utils.NewError(c.GetCard, msg, c, err)
	} else if err := resp.Decode(&crd); err != nil {
		return nil, utils.NewError(c.GetCard, "unable to decode response", c, err)
	} else {
		contactCard := &components.ContactCard{
			Card: crd,
			Href: path,
		}

		return contactCard, nil
	}
}

// creates or updates one or more cards on the remote CardDAV server
func (c *Client) PutCards(path string, cards ...*components.Card) error {
	if req, err := c.Server().NewRequest("PUT", path, cards); err != nil {
		return utils.NewError(c.PutCards, "unable to encode request", c, err)
	} else if resp, err := c.Do(req); err != nil {
		return utils.NewError(c.PutCards, "unable to execute request", c, err)
	} else if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		err := new(entities.Error)
		resp.WebDAV().Decode(err)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		return utils.NewError(c.PutCards, msg, c, err)
	}
	return nil
}

func (c *Client) DeleteCard(path string) error {
	req, err := c.Server().NewRequest("DELETE", path)
	if err != nil {
		return utils.NewError(c.DeleteCard, "unable to encode request", c, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return utils.NewError(c.DeleteCard, "unable to execute request", c, err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		err := new(entities.Error)
		resp.WebDAV().Decode(err)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		return utils.NewError(c.DeleteCard, msg, c, err)
	}

	return nil
}
