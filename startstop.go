package startstop

import (
    "sync"
    "errors"
    "time"
)

type StartStop struct {
    sync.RWMutex
    status string
    ch     []chan struct{}
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

func (s *StartStop) AddChan() (ch chan struct{}, status string) {
    s.Lock()
    defer s.Unlock()

    status = s.status
    if status == IsStop {
        ch = make(chan struct{}, 1)
        s.ch = append(s.ch, ch)
    }
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
            for _, ch := range s.ch {
                ch <- struct{}{}
            }
            s.ch = []chan struct{}{}
        case IsStart, IsFinish:
    }
    return s
}

func (s *StartStop) Stop() {
    s.Lock()
    defer s.Unlock()

    switch s.status {
        case "":
            s.status = IsStop
        case IsStart:
            s.status = IsStop
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
            for _, ch := range s.ch {
                ch <- struct{}{}
            }
            s.ch = []chan struct{}{}
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
    status = s.GetStatus()
    var is_next chan struct{}

    if status == IsStop {
        is_next, status = s.AddChan()
    }

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
