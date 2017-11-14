package rethink

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
	mapstruct "github.com/mitchellh/mapstructure"
)

// Status keeps track of the status
type Status struct {
	Tables      map[string]struct{}      // All tables we're monitoring
	TableActive []string                 // All tables that are confirmed to be currently active
	TableDown   map[string]statusDetails // All tables that are confirmed to be currently down
	DBs         map[string]struct{}      // All databases we're monitoring
	DBActive    []string                 // All databases that are confirmed to be currently active
	DBDown      []string                 // All databases that are confirmed to be currently down
	LastUpdated time.Time                // The last time any value was updated
}

type statusDetails struct {
	AllReplicasReady      bool `mapstructure:"all_replicas_ready"`
	ReadyForOutdatedReads bool `mapstructure:"ready_for_outdated_reads"`
	ReadyForReads         bool `mapstructure:"ready_for_reads"`
	ReadyForWrites        bool `mapstructure:"ready_for_writes"`
}

type statusInfo struct {
	DB   string // DB table is stored in
	ID   string // Table UUID
	Name string // Table name
	// Shards []struct{} // Not needed currently, for replicas
	Status statusDetails // Subfields show status of the table
}

func (s statusInfo) isReady() bool {
	return (s.Status.AllReplicasReady && s.Status.ReadyForOutdatedReads &&
		s.Status.ReadyForReads && s.Status.ReadyForWrites)
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
			// Reset the table lists
			s.TableActive = []string{}
			s.TableDown = map[string]statusDetails{}
			// Iterate over the defined tables and check their status
			for table := range s.Tables {
				var status statusInfo

				raw, err := c.Status(table)
				if err != nil {
					// The status lookup failed, do something with that.
					log.Error(err.Error())
					continue
				}

				err = mapstruct.Decode(raw[0], &status)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				// Check that the table is ready
				if status.isReady() {
					s.TableActive = append(s.TableActive, status.Name)
				} else {
					s.TableDown[status.Name] = status.Status
				}

				// Delay next update for 30 seconds
				time.Sleep(30 * time.Second)
			}
		}
	}()
}

// APIStatusHandler returns the current status of the API, as of
func (s *Status) APIStatusHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"updated": s.LastUpdated.Format(time.UnixDate),
		"tables": map[string][]string{
			"active":  s.TableActive,
			"down":    s.TableDown,
			"unknown": diffMap(s.Tables, s.TableActive, s.TableDown),
		},
		"databases": map[string][]string{
			"active":  s.DBActive,
			"down":    s.DBDown,
			"unknown": diffMap(s.DBs, s.DBActive, s.DBDown),
		},
	})
}

func diffMap(source map[string]struct{}, maps ...[]string) []string {
	var undefined []string
	// Iterate over all the arrays of strings sent, we'll be checking the values
	// in these individually to see if there are any values defined in source not
	// in any of the sub-maps
	for _, currentMap := range maps {
		// Iterate over the values in the current array
		for _, val := range currentMap {
			// Does the current value exist in source?
			if _, ok := source[val]; !ok {
				undefined = append(undefined, val)
			}
		}
	}

	return undefined
}
