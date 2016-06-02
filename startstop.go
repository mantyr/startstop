package startstop

import (
    "sync"
    "errors"
    "time"
)

type StartStop struct {
    sync.RWMutex
    status string
    ch     chan struct{}
}

const (
    IsStart       string = "start"
    IsStop        string = "stop"
    IsFinish      string = "finish"
    IsAlternative string = "alternative"
)

func NewStartStop() (s *StartStop) {
    s = new(StartStop)
    return
}

func (s *StartStop) Start() *StartStop {
    s.Lock()
    defer s.Unlock()

    switch s.status {
        case "":
            s.status = IsStart
        case IsStop:
            s.status = IsStart
            s.ch <- struct{}{}
        case IsStart, IsFinish:
    }
    return s
}

func (s *StartStop) Stop() {
    s.Lock()
    defer s.Unlock()

    if s.status == IsStart {
        s.status = IsStop
        s.ch = make(chan struct{}, 1)
    }
}

func (s *StartStop) Finish() {
    s.Lock()
    defer s.Unlock()

    switch s.status {
        case "":
            s.status = IsFinish
        case IsStop:
            s.status = IsFinish
            s.ch <- struct{}{}
        case IsStart:
            s.status = IsFinish
    }
}

func (s *StartStop) GetStatus() string {
    s.RLock()
    defer s.RUnlock()
    return s.status
}

func (s *StartStop) Next(params ...<-chan time.Time) (status string, err error) {
    s.RLock()
    status   = s.status
    is_next := s.ch
    s.RUnlock()

    switch status {
        case IsStart, IsFinish:
            return status, nil
        case IsStop:
            if len(params) > 0 {
                select {
                    case <- is_next:
                        return s.Next()
                    case <- params[0]:
                        return IsAlternative, nil
                }
            } else {
                <- is_next
                return s.Next(params...)
            }
    }
    return "", errors.New("StartStop incorrect status")
}
