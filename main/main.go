package main

import (
	"fmt"
	"jobqueue"
	"time"
)

type Sleep struct {
	Id int
}

func (s *Sleep)Todo() error {
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("ID %v wake up\n", s.Id)
	return nil
}

func main()  {
	f := jobqueue.NewFactory(100)
	f.Run()
	for i := 0; i < 10000; i++ {
		var sleep jobqueue.Job = &Sleep{i}
		f.SetAJob(&sleep)
	}
	f.Wait()
}
