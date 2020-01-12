package main

import (
	"fmt"
	"jobqueue"
	"sync"
	"time"
)

type Sleep struct {
	Id int
}
var count = 0
var mu =  sync.Mutex{}

func (s *Sleep)Todo() error {
	mu.Lock()
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("ID %v wake up\n", s.Id)
	count++
	mu.Unlock()
	return nil
}

func main()  {
	f := jobqueue.NewFactory(100)
	f.Run()
	for i := 0; i < 100000; i++ {
		var sleep jobqueue.Job = &Sleep{i}
		f.SetAJob(&sleep)
	}

	f.Wait()
	f.Close()
	fmt.Println(count)
}
