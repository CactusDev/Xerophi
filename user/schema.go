package user

// Database is all the data stored in the DB for a user record
type Database struct {
	ID        string // RethinkDB Record UUID
	Hash      string // Argon2 password hash
	Token     string // The user's token
	DeletedAt int    // A unix timestamp value for soft deletion, 0 if active
	CreatedAt string // A timestamp string for when the record was created
	UserID    int    // The internal static numeric ID for the user
	UserName  string // The username for the user, can be changed unlike UserID
	Service   string // The default service this user uses
}

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID string `jsonapi:"primary,command"`

	Service  string `jsonapi:"attr,service"`
	UserName string `jsonapi:"attr,username"`
	UserID   int    `jsonapi:"attr,userId"`

	Token     string `jsonapi:"meta,token"`
	CreatedAt string `jsonapi:"meta,createdAt"`
}
