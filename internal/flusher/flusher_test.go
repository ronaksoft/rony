package flusher

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"sync"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - May - 05
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	// testEnv.Init()
}

func TestFlusherLegacy(t *testing.T) {
	f := New(1000, 100, 500*time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Ret != nil {
				items[idx].Ret <- items[idx].Key
			}
		}
	})
	starTime := time.Now()
	waitGroup := sync.WaitGroup{}
	for j := 0; j < 1000; j++ {
		waitGroup.Add(1)
		go func(j int) {
			for i := 0; i < 10; i++ {
				<-f.EnterWithChan(j, "v")
				// fmt.Println("Chan:", <-f.EnterWithChan(j, "v"))
			}
			waitGroup.Done()
		}(j)
	}
	waitGroup.Wait()
	fmt.Println("Legacy Finished", time.Now().Sub(starTime))
}

func TestFlusherLifo(t *testing.T) {
	f := NewLifo(1000, 100, 25*time.Millisecond, func(items []Entry) {
		time.Sleep(time.Duration(tools.RandomInt(100)) * time.Millisecond)
		for idx := range items {
			if items[idx].Ret != nil {
				items[idx].Ret <- items[idx].Key
			} else {
				items[idx].Callback(items[idx].Key)
			}
		}
	})
	// time.Sleep(5 * time.Second)
	for i := 0; i < 10; i++ {
		waitGroup := sync.WaitGroup{}
		for j := 0; j < 100; j++ {
			waitGroup.Add(1)
			go func(j int) {
				for i := 0; i < 10; i++ {
					// f.EnterWithResult(j, "v")
					<-f.EnterWithChan(j, "v")
				}
				waitGroup.Done()
			}(j)
		}
		waitGroup.Wait()
	}
}

func TestFlusherBugLifo(t *testing.T) {
	f := NewLifo(1000, 100, 25*time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Ret != nil {
				items[idx].Ret <- items[idx].Key
			}
		}
	})
	for i := 0; i < 10; i++ {
		<-f.EnterWithChan(i, "v")
	}
}

var maxBatchSize = int32(1000)

func benchLifoEnterWithChan(concurrency int32, b *testing.B) {
	f := NewLifo(maxBatchSize, concurrency, time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Ret != nil {
				items[idx].Ret <- items[idx].Key
			}
		}
	})
	b.StartTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			<-f.EnterWithChan(1, "v")
		}
	})
}
func benchLifoEnterWithResult(concurrency int32, b *testing.B) {
	f := NewLifo(maxBatchSize, concurrency, time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Callback != nil {
				items[idx].Callback(items[idx].Key)
			}
		}
	})
	b.StartTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = f.EnterWithResult(1, "v")
		}
	})
}
func benchEnterWithChan(concurrency int32, b *testing.B) {
	f := New(maxBatchSize, concurrency, time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Ret != nil {
				items[idx].Ret <- items[idx].Key
			}
		}
	})
	b.StartTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			<-f.EnterWithChan(1, "v")
		}
	})
}
func benchEnterWithResult(concurrency int32, b *testing.B) {
	f := New(maxBatchSize, concurrency, time.Millisecond, func(items []Entry) {
		for idx := range items {
			if items[idx].Callback != nil {
				items[idx].Callback(items[idx].Key)
			}
		}
	})
	b.StartTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = f.EnterWithResult(1, "v")
		}
	})
}
func BenchmarkFlusher(b *testing.B) {
	bs := []struct {
		name        string
		concurrency int32
	}{
		{"LifoWithChan", 100},
		{"LifoWithResult", 100},
		{"LegacyWithChan", 100},
		{"LegacyWithResult", 100},
	}
	for _, item := range bs {
		b.Run(item.name, func(b *testing.B) {
			switch item.name {
			case "LifoWithChan":
				benchLifoEnterWithChan(item.concurrency, b)
			case "LifoWithResult":
				benchLifoEnterWithResult(item.concurrency, b)
			case "LegacyWithChan":
				benchEnterWithChan(item.concurrency, b)
			case "LegacyWithResult":
				benchEnterWithResult(item.concurrency, b)
			}
		})
	}
}
