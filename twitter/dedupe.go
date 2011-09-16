package twitter

import (
	"container/list"
)

// Copy updates from S to C, dropping duplicates
type Deduper struct {
	Updates chan Update
	S       Stream
	recent  *list.List
	// How many updates to keep track of for duplicates (default 100)
	Track int
}

func NewDeduper(s Stream) *Deduper {
	d := &Deduper{Updates: make(chan Update, 100), S: s, recent: list.New(), Track: 100}
	go d.process()
	return d
}

func (d *Deduper) process() {
	var u Update
	var ok, seen bool
	s := d.S.C()
	for {
		if u, ok = <-s; ok {
			// check to see if we've already seen this one
			seen = false
			for i := d.recent.Back(); i != nil; i = i.Prev() {
				if u.Id == i.Value.(uint64) {
					seen = true
					break
				}
			}
			if !seen {
				d.Updates <- u
				d.recent.PushBack(u.Id)
				for d.recent.Len() > d.Track {
					d.recent.Remove(d.recent.Front())
				}
			}
		} else {
			// upstream is closed.
			close(d.Updates)
			return
		}
	}
}

func (d *Deduper) C() <-chan Update {
	return d.Updates
}

func (d *Deduper) Close() {
	// Close upstream stream.  process loop will close this stream
	d.S.Close()
}
