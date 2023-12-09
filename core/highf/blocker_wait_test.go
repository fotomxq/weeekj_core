package CoreHighf

import (
	"testing"
	"time"
)

func TestBlockerWait(t *testing.T) {
	var waitBlocker BlockerWait
	waitBlocker.Init(3)
	go func() {
		time.Sleep(time.Second * 1)
		step := 0
		for {
			if step > 2 {
				return
			}
			waitBlocker.CheckWait(1, "a1", func(modID int64, modMark string) {
				t.Log("mod id: ", modID, ", mark: ", modMark)
			})
			step += 1
			time.Sleep(time.Second * 1)
		}
	}()
	waitBlocker.CheckWait(1, "a1", func(modID int64, modMark string) {
		t.Log("mod id: ", modID, ", mark: ", modMark)
	})
	waitBlocker.CheckWait(1, "a1", func(modID int64, modMark string) {
		t.Log("mod2 id: ", modID, ", mark: ", modMark)
	})
	go func() {
		time.Sleep(time.Second * 1)
		step := 0
		for {
			if step > 2 {
				return
			}
			waitBlocker.CheckWait(1, "a1", func(modID int64, modMark string) {
				t.Log("mod3 id: ", modID, ", mark: ", modMark)
			})
			step += 1
			time.Sleep(time.Second * 1)
		}
	}()
	waitBlocker.CheckWait(2, "a2", func(modID int64, modMark string) {
		t.Log("mod4 id: ", modID, ", mark: ", modMark)
	})
}
