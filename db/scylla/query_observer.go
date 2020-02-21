package scylla

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/metrics"
	"github.com/gocql/gocql"
	"time"
)

/*
   Creation Time: 2019 - Nov - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type queryObserver struct{}

func (q queryObserver) ObserveQuery(ctx context.Context, oq gocql.ObservedQuery) {
	metrics.Histogram(metrics.HistDbQueryLatenciesMS).Observe(float64(oq.End.Sub(oq.Start) / time.Millisecond))
}
