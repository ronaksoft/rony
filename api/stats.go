package api

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/db/redis"
	"git.ronaksoftware.com/ronak/rony/pools"
	"github.com/mediocregopher/radix/v3"
	"time"
)

/*
   Creation Time: 2020 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	dayInSec  = 24 * 60 * 60
	weekInSec = 7 * dayInSec
)

func SetUniqueUserVisit(userID int64) {
	t := time.Now().Unix()
	dKey := t / dayInSec
	wKey := t / weekInSec
	waitGroup := pools.AcquireWaitGroup()
	waitGroup.Add(2)
	go func() {
		_ = redis.TempCache().Do(radix.FlatCmd(nil, "PFADD", fmt.Sprintf("UVD.%d", dKey), userID))
		waitGroup.Done()
	}()
	go func() {
		_ = redis.TempCache().Do(radix.FlatCmd(nil, "PFADD", fmt.Sprintf("UVW.%d", wKey), userID))
		waitGroup.Done()
	}()
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)
}

func CountUniqueVisits() (week, day, weekBefore, dayBefore int64) {
	t := time.Now().Unix()
	dKey := t / dayInSec
	wKey := t / weekInSec
	_ = redis.TempCache().Do(radix.FlatCmd(&day, "PFCOUNT", fmt.Sprintf("UVD.%d", dKey)))
	_ = redis.TempCache().Do(radix.FlatCmd(&week, "PFCOUNT", fmt.Sprintf("UVW.%d", wKey)))
	_ = redis.TempCache().Do(radix.FlatCmd(&dayBefore, "PFCOUNT", fmt.Sprintf("UVD.%d", dKey-1)))
	_ = redis.TempCache().Do(radix.FlatCmd(&weekBefore, "PFCOUNT", fmt.Sprintf("UVW.%d", wKey-1)))
	return
}
