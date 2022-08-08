package dispatcher

import (
	"cam_cp/app/watcher"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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
	ctx := context.Background()
	go d.Run(ctx)

	go func() {
		d.in[0] <- fl
		d.in[1] <- fl
		d.in[2] <- fl
	}()

	res := func() [][]watcher.ExData {
		var res [][]watcher.ExData
		res = append(res, <-d.out[0])
		res = append(res, <-d.out[1])
		res = append(res, <-d.out[2])
		return res
	}()

	assert.Equal(t, fl, res[0], "expected %v, got %v", fl, res[0])
	assert.Equal(t, fl, res[1], "expected %v, got %v", fl, res[1])
	assert.Equal(t, fl, res[2], "expected %v, got %v", fl, res[2])
}
