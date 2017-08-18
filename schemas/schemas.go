package schemas

import (
	"strings"
)

type Message struct {
	Text   []Component `json:"text"` // FIXME: Anything besides [] breaks this. It is assumed that this is due to the interface being used for Component.
	Action bool        `json:"action"`
}

type Component interface {
	Text() string
}

func (t Text) Text() string {
	return t.Data
}

func (e Emoji) Text() string {
	return e.Data
}

func (t Tag) Text() string {
	return t.Data
}

func (u URL) Text() string {
	return u.Data
}

func (v Variable) Text() string {
	return v.Data
}

type Text struct {
	Data string `json:"data"`
}
type Emoji struct {
	Data string `json:"data"`
}
type Tag struct {
	Data string `json:"data"`
}
type URL struct {
	Data string `json:"data"`
}
type Variable struct {
	Data string `json:"data"`
}

type Channel string
type User string
type Service string

type Context struct {
	Packet  Message `json:"packet"`
	Channel Channel `json:"channel"`
	User    User    `json:"user,omitempty"`
	Role    Role    `json:"role,omitempty"`
	Target  User    `json:"target,omitempty"`
	Service Service `json:"service"`
}

type Role int

const (
	banned Role = iota
	user
	subscriber
	moderator
	owner
)

func (r *Role) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch {
	case str == "banned":
		*r = banned

	case str == "user":
		*r = user

	case str == "subscriber":
		*r = subscriber

	case str == "moderator":
		*r = moderator

	case str == "owner":
		*r = owner

	default:
		*r = user
	}

	return nil
}
