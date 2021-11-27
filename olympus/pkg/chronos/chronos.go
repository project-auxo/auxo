package chronos

import (
	"time"

	ntp "github.com/beevik/ntp"
)

// The pool below will change every hour.
const ntpServer = "0.pool.ntp.org"

type Chronos struct{}

func (chronos *Chronos) GetTime() (time time.Time, err error) {
	time, err = ntp.Time(ntpServer)
	return
}

func (chronos *Chronos) GetTimeWithQuery(queryOptions ntp.QueryOptions) (response *ntp.Response, err error) {
	response, err = ntp.QueryWithOptions(ntpServer, queryOptions)
	return
}
