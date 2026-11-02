package core

import (
	"sync"
	"time"

	"github.com/metacubex/mihomo/log"
)

// maxLogLines is Verge's MAX_LOG_NUM. Older lines are dropped.
const maxLogLines = 1000

// LogLine is one kernel log event. Level is mihomo's type string
// (debug, info, warning, error, silent). Time is when this process
// received the event; the kernel event itself has no timestamp.
type LogLine struct {
	Time    time.Time
	Level   string
	Payload string
}

type storedLog struct {
	seq  uint64
	line LogLine
}

// logBook is the ring the UI drains. seq is monotonic so a cursor can
// ask for only the lines it has not seen, including across Clear.
type logBook struct {
	mu    sync.Mutex
	lines []storedLog
	next  uint64
}

func (b *logBook) append(line LogLine) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, storedLog{seq: b.next, line: line})
	b.next++
	if extra := len(b.lines) - maxLogLines; extra > 0 {
		b.lines = append([]storedLog(nil), b.lines[extra:]...)
	}
}

func (b *logBook) since(seq uint64) ([]LogLine, uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []LogLine
	for _, item := range b.lines {
		if item.seq >= seq {
			out = append(out, item.line)
		}
	}
	return out, b.next
}

func (b *logBook) clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = nil
}

// startLogFeed subscribes to the process-wide mihomo logger. The
// goroutine must keep reading: a full subscriber blocks every log line
// in the kernel. Caller holds c.mu.
func (c *Core) startLogFeed() {
	if c == nil || c.logDone != nil {
		return
	}
	sub := log.Subscribe()
	done := make(chan struct{})
	c.logSub = sub
	c.logDone = done
	go func() {
		defer close(done)
		for ev := range sub {
			c.logBook.append(LogLine{
				Time:    time.Now(),
				Level:   ev.Type(),
				Payload: ev.Payload,
			})
		}
	}()
}

// stopLogFeed unsubscribes and waits until the reader has exited.
// Caller holds c.mu. The reader never takes c.mu.
func (c *Core) stopLogFeed() {
	if c == nil || c.logDone == nil {
		return
	}
	log.UnSubscribe(c.logSub)
	c.logSub = nil
	done := c.logDone
	c.logDone = nil
	<-done
}

// LogsSince returns lines recorded at seq or later, and the cursor to
// pass next time. ok is false when the core is not running.
func (c *Core) LogsSince(seq uint64) (lines []LogLine, next uint64, ok bool) {
	if c == nil {
		return nil, seq, false
	}
	c.mu.Lock()
	running := c.running
	c.mu.Unlock()
	if !running {
		return nil, seq, false
	}
	lines, next = c.logBook.since(seq)
	return lines, next, true
}

// ClearLogs drops the buffered lines. The sequence cursor keeps
// increasing, so a later LogsSince does not replay what was cleared.
// It reports false when the core is not running.
func (c *Core) ClearLogs() bool {
	if c == nil {
		return false
	}
	c.mu.Lock()
	running := c.running
	c.mu.Unlock()
	if !running {
		return false
	}
	c.logBook.clear()
	return true
}
