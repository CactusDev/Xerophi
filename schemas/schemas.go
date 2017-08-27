package schemas

import (
	"strings"
)

// TODO: Remove `jsonapi:` tags. Should just be `json:`.
type Message struct {
	Text   []Component `jsonapi:"attr,text"`
	Action bool        `jsonapi:"attr,action"`
}

// TODO: Remove `jsonapi:` tags. Should just be `json:`.
type Component struct {
	Type string `jsonapi:"attr,type"` // HACK: need to validate type
	Data string `jsonapi:"attr,data"`
}

type Channel string
type User string
type Service string

// TODO: Remove `jsonapi:` tags. Should just be `json:`.
type Context struct {
	Packet  Message `jsonapi:"attr,packet"`
	Channel Channel `jsonapi:"attr,channel"`
	User    User    `jsonapi:"attr,user,omitempty"`
	Role    Role    `jsonapi:"attr,role,omitempty"` // FIXME: should return string ("owner"), not int (5). See UnmarshalJSON below.
	Target  User    `jsonapi:"attr,target,omitempty"`
	Service Service `jsonapi:"attr,service"`
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
