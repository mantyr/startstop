# Golang Start/Stop goroutine

[![Build Status](https://travis-ci.org/mantyr/startstop.svg?branch=master)](https://travis-ci.org/mantyr/startstop)
[![GoDoc](https://godoc.org/github.com/mantyr/startstop?status.png)](http://godoc.org/github.com/mantyr/startstop)
[![Software License](https://img.shields.io/badge/license-The%20Not%20Free%20License,%20Commercial%20License-brightgreen.svg)](LICENSE.md)

This stable version

## Installation

    $ go get github.com/mantyr/startstop

## Example
```GO
package main

import (
    "github.com/mantyr/startstop"
    "time"
)
func main() {
    s := NewStartStop()

    status, err := s.Next()       // status = "",      err = errors.New("StartStop incorrect status")
    s.Start()                     // status = "start"
    s.Finish()                    // status = "finish"

    go func(){
        time.Sleep(10 * time.Second)
        s.Start()                 // status = "start", send for Next()
    }()
    status, err = s.Next()        // sleep 10 sesond, status = "start"

    s.Stop()                      // status = "stop"
    go func(){
        time.Sleep(10 * time.Second)
        s.Finish()                // status = "finish", send for Next()
    }()
    status, err = s.Next(time.After(100 * time.Second))        // sleep 10 or 100 sesond, status = "finish" or status = "alternative"

    if err != nil {
        // error
    }
    if status == startstop.IsFinish {
        // finish
        return
    }
    // next
```

## Author

[Oleg Shevelev][mantyr]

[mantyr]: https://github.com/mantyr
