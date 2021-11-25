package console

import (
	"bytes"
	"io"
	"sync"
	"sync/atomic"
)

type chunk struct {
	a []byte
	n int
}

type chunkReader struct {
	chunks	   []chunk
 	pendingErr error
	pipe       io.Reader
	stopLoop   uintptr
	m          sync.Mutex
}

func (c *chunkReader) loop() {
	// Used to stop racing before the start of the loop.
	c.m.Unlock()

	for atomic.LoadUintptr(&c.stopLoop) == 0 {
		// Read from the pipe.
		a := make([]byte, 3)
		n, err := c.pipe.Read(a)
		if atomic.LoadUintptr(&c.stopLoop) == 1 {
			// The loop was stopped in the time since last run.
			return
		}

		// Handle setting the values in the struct.
		c.m.Lock()
		if err == nil {
			// Append to the chunks.
			doAppend := true
			for _, v := range c.chunks {
				if bytes.Equal(v.a, a) && v.n == n {
					doAppend = false
					break
				}
			}
			if doAppend {
				c.chunks = append(c.chunks, chunk{a: a, n: n})
			}
			c.m.Unlock()
		} else {
			// In this case, put the error and stop the loop.
			c.pendingErr = err
			c.m.Unlock()
			return
		}
	}
}

func (c *chunkReader) close() {
	atomic.StoreUintptr(&c.stopLoop, 1)
}

func (c *chunkReader) flush() []chunk {
	c.m.Lock()
	defer c.m.Unlock()
	a := c.chunks
	c.chunks = []chunk{}
	return a
}

func newChunkReader(r io.Reader) *chunkReader {
	c := chunkReader{chunks: []chunk{}, pipe: r}
	c.m.Lock()
	go c.loop()
	return &c
}
