package startstop

import (
    "testing"
    "time"
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