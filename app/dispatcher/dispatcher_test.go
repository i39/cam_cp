package dispatcher

import (
	"cam_cp/app/watcher"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	fl := []watcher.ExData{
		{
			Name: "/test2/test3/test3.txt",
			Data: []byte("test3"),
		},
		{
			Name: "/test2/test2.txt",
			Data: []byte("test2"),
		},

		{
			Name: "/test1.txt",
			Data: []byte("test1"),
		},
	}

	inChan1 := make(watcher.ExChan)
	inChan2 := make(watcher.ExChan)
	inChan3 := make(watcher.ExChan)

	outChan1 := make(watcher.ExChan)
	outChan2 := make(watcher.ExChan)
	outChan3 := make(watcher.ExChan)

	var d Impl
	d.AddIn(inChan1)
	d.AddIn(inChan2)
	d.AddIn(inChan3)

	d.AddOut(outChan1)
	d.AddOut(outChan2)
	d.AddOut(outChan3)
	ctx, cancel := context.WithCancel(context.Background())
	go d.Run(ctx)

	go func() {
		time.Sleep(time.Second * 1)
		d.in[0] <- fl
		d.in[1] <- fl
		d.in[2] <- fl
	}()

	var o1, o2, o3 []watcher.ExData

	go func() {
		time.Sleep(time.Second * 2)
		o1 = <-d.out[0]
		o2 = <-d.out[1]
		o3 = <-d.out[2]
	}()

	time.Sleep(time.Second * 5)
	assert.Equal(t, fl, o1, "expected %v, got %v", fl, o1)
	assert.Equal(t, fl, o2, "expected %v, got %v", fl, o2)
	assert.Equal(t, fl, o3, "expected %v, got %v", fl, o3)

	cancel()

}
