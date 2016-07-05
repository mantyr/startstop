package startstop

import (
    "testing"
    "time"
    "sync"
)

func TestStartStop(t *testing.T) {
    s := NewStartStop()

    status, err := s.Next() // 1
    if err == nil {
        t.Errorf("Error status, %q, %q", status, err)
    }

    s.Start()
    s.Start()
    s.Start()
    s.Start()
    s.Start()
    s.Start()
    s.Start()
    s.Start()

    status, err = s.Next()  // 2
    if status != IsStart {
        t.Errorf("Error status, %q, %q", status, err)
    }

    go func(){
         time.Sleep(10*time.Second)
         s.Start()
    }()
    s.Stop()
    status, err = s.Next(time.After(30 * time.Second)) // 3
    if status == IsAlternative {
        t.Errorf("Error Stop() and Next()")
    }

    go func(){
         time.Sleep(10*time.Second) // for 5
         s.Start()
    }()
    s.Stop()
    if s.GetStatus() != IsStop {
        t.Errorf("Error GetStatus()")
    }
    status, err = s.Next(time.After(3 * time.Second)) // 4
    if status != IsAlternative {
        t.Errorf("Error Stop() and Next() - alternative")
    }
    status, err = s.Next(time.After(30 * time.Second)) // 5
    if status == IsAlternative {
        t.Errorf("Error Stop() and Next()")
    }

    s.Finish()
    if s.GetStatus() != IsFinish {
        t.Errorf("Error GetStatus()")
    }

    status, err = s.Next(time.After(3 * time.Second)) // 6
    if status == IsAlternative {
        t.Errorf("Error Finish() - alternative")
    }
    if status != IsFinish || err != nil {
        t.Errorf("Error Finish(), %q, %q", status, err)
    }
    if s.GetStatus() != IsFinish {
        t.Errorf("Error GetStatus()")
    }
}

func TestConcurrentNext(t *testing.T) {
    s := NewStartStop()
    s.Start()
    s.Stop()

    var wg sync.WaitGroup

    var status_1 string
    var status_2 string

    s.AddChan()
    s.AddChan()
    s.AddChan()

    wg.Add(1)
    go func(){
        defer wg.Done()
        var err error
        status_1, err = s.Next(time.After(30 * time.Second))
        if err != nil {
            status_1 = err.Error()
        }
    }()
    wg.Add(1)
    go func(){
        defer wg.Done()
        var err error
        status_2, err = s.Next(time.After(15 * time.Second))
        if err != nil {
            status_2 = err.Error()
        }
    }()
    go func(){
        time.Sleep(7 * time.Second)
        s.Start()
    }()
    wg.Wait()

    if status_1 != status_2 || status_1 != IsStart {
        t.Errorf("Error concurrent Next, %q, %q", status_1, status_2)
    }

    s.Stop()

    wg.Add(1)
    go func(){
        defer wg.Done()
        var err error
        status_1, err = s.Next(time.After(30 * time.Second))
        if err != nil {
            status_1 = err.Error()
        }
    }()
    wg.Add(1)
    go func(){
        defer wg.Done()
        var err error
        status_2, err = s.Next(time.After(15 * time.Second))
        if err != nil {
            status_2 = err.Error()
        }
    }()
    go func(){
        time.Sleep(7 * time.Second)
        s.Start()
    }()
    wg.Wait()

    if status_1 != status_2 || status_1 != IsStart {
        t.Errorf("Error concurrent Next (phase 2), %q, %q", status_1, status_2)
    }

}