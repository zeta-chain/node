package main

import (
	"fmt"
	"time"
)

func (sm *SmokeTest) TestMyTest() {
	LoudPrintf("My test\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
}
