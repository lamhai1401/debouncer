package debouncer

import (
	"sync"
	"time"
)

// debouncer linter
type debouncer struct {
	lastExcute      time.Time
	rootDuration    time.Duration // set debounce milisecond time
	currentDuration time.Duration // set debounce milisecond time
	timer           *time.Timer
	mutex           sync.RWMutex
}

// NewBouncer
// param: t milisecond
func NewBouncer(
	t time.Duration, // milisecond
) *debouncer {

	return &debouncer{
		rootDuration: t,
	}
}

func (d *debouncer) getLastExcute() time.Time {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.lastExcute
}

func (d *debouncer) setLastExcute(t time.Time) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.lastExcute = t
}

func (d *debouncer) getTimer() *time.Timer {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.timer
}

func (d *debouncer) setTimer(t *time.Timer) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	d.timer = t
}

func (d *debouncer) getRootDuration() time.Duration {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.rootDuration
}

func (d *debouncer) getCurrentDuration() time.Duration {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.currentDuration
}

func (d *debouncer) setCurrentDuration(t time.Duration) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	d.currentDuration = t
}

func (d *debouncer) stopTimer() {
	if timer := d.getTimer(); timer != nil {
		timer.Stop()
		d.setTimer(nil)
	}
}

// Add init timer
func (d *debouncer) Add(f func()) {
	excutefunc := func() {
		go f()
		d.setCurrentDuration(d.getRootDuration())
		d.setLastExcute(time.Now())
		d.stopTimer()
	}

	handleFunc := func() {
		d.stopTimer()
		d.setTimer(time.AfterFunc(d.getCurrentDuration(), excutefunc))
	}

	if d.getCurrentDuration() == 0 {
		d.setCurrentDuration(d.getRootDuration())
		excutefunc()
		return
	}

	// check time since last excute
	gapTime := time.Since(d.getLastExcute()).Milliseconds()
	if gapTime < int64(d.getCurrentDuration()) {
		d.setCurrentDuration(d.getRootDuration() - time.Duration(gapTime))
		handleFunc()
	} else {
		d.stopTimer()
		excutefunc()
	}
}
