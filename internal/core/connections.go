package core

import (
	"errors"
	"net/netip"
	"time"

	"github.com/metacubex/mihomo/tunnel/statistic"
)

// ConnSnap is one live connection from the mihomo statistic manager.
// The UI copies it; the core does not keep the slice.
type ConnSnap struct {
	ID                string
	Host              string
	DestinationIP     string
	DestPort          uint16
	RemoteDestination string
	Process           string
	ProcessPath       string
	Network           string
	Type              string
	SourceIP          string
	SourcePort        uint16
	Rule              string
	RulePayload       string
	Chains            []string
	Upload            int64
	Download          int64
	Start             time.Time
}

// ConnectionSnapshot reads the process-wide statistic manager.
// ok is false when the core is not running, so a stopped kernel does not
// look like an empty traffic counter.
func (c *Core) ConnectionSnapshot() (up, down int64, items []ConnSnap, ok bool) {
	if c == nil {
		return 0, 0, nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return 0, 0, nil, false
	}
	snap := statistic.DefaultManager.Snapshot()
	if snap == nil {
		return 0, 0, nil, true
	}
	items = make([]ConnSnap, 0, len(snap.Connections))
	for _, info := range snap.Connections {
		if info == nil {
			continue
		}
		items = append(items, snapFrom(info))
	}
	return snap.UploadTotal, snap.DownloadTotal, items, true
}

// CloseAllConnections closes every tracked connection. It fails when the
// core is not running; the caller decides whether to surface that.
func (c *Core) CloseAllConnections() error {
	if c == nil {
		return errors.New("core is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return errors.New("core is not running")
	}
	statistic.DefaultManager.Range(func(t statistic.Tracker) bool {
		_ = t.Close()
		return true
	})
	return nil
}

// CloseConnection closes one tracked connection by id. A missing id is
// not an error: the flow may already have ended.
func (c *Core) CloseConnection(id string) error {
	if c == nil {
		return errors.New("core is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return errors.New("core is not running")
	}
	if t := statistic.DefaultManager.Get(id); t != nil {
		_ = t.Close()
	}
	return nil
}

func snapFrom(info *statistic.TrackerInfo) ConnSnap {
	item := ConnSnap{
		ID:          info.UUID.String(),
		Rule:        info.Rule,
		RulePayload: info.RulePayload,
		Start:       info.Start,
		Upload:      info.UploadTotal.Load(),
		Download:    info.DownloadTotal.Load(),
	}
	if len(info.Chain) > 0 {
		item.Chains = append([]string(nil), info.Chain...)
	}
	if md := info.Metadata; md != nil {
		item.Host = md.Host
		item.DestinationIP = addrString(md.DstIP)
		item.DestPort = md.DstPort
		item.RemoteDestination = md.RemoteDst
		item.Process = md.Process
		item.ProcessPath = md.ProcessPath
		item.Network = md.NetWork.String()
		item.Type = md.Type.String()
		item.SourceIP = addrString(md.SrcIP)
		item.SourcePort = md.SrcPort
	}
	return item
}

func addrString(a netip.Addr) string {
	if !a.IsValid() {
		return ""
	}
	return a.String()
}
