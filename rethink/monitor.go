package rethink

import (
	"net/http"
	"time"

	"github.com/Google/uuid"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Status keeps track of the status
type Status struct {
	Tables      map[string]struct{} // All tables we're monitoring
	DBs         map[string]struct{} // All databases we're monitoring
	Issues      []Issue             // Any active issues
	LastUpdated time.Time           // The last time any value was updated
}

// NewMonitor initializes a new Status monitor for the tables and DBs given
func NewMonitor(tables []string, dbs []string) Status {
	var mappedTables = map[string]struct{}{}
	var mappedDBs = map[string]struct{}{}

	for _, table := range tables {
		mappedTables[table] = struct{}{}
	}

	for _, db := range dbs {
		mappedDBs[db] = struct{}{}
	}

	return Status{
		Tables: mappedTables,
		DBs:    mappedDBs,
	}
}

// Monitor monitors the connection status for the DB and can reconnect
func (s *Status) Monitor(c *Connection) {
	go func() {
		for {
			issues, err := c.Status()
			if err != nil {
				log.Error(err.Error())
				s.Issues = []Issue{
					{
						ID:          uuid.New().String(),
						Type:        "internal",
						Description: "tl;dr we suck at making APIs",
					},
				}
				continue
			}
			s.Issues = issues

			s.LastUpdated = time.Now()
			// Delay next update for 30 seconds
			time.Sleep(30 * time.Second)
		}
	}()
}

// APIStatusHandler returns the current status of the API, as of
func (s *Status) APIStatusHandler(ctx *gin.Context) {
	var code = http.StatusOK
	if len(s.Issues) > 0 {
		code = http.StatusInternalServerError
	}
	ctx.JSON(code, gin.H{
		"updated": s.LastUpdated.Format(time.UnixDate),
		"issues":  s.Issues,
	})
}
