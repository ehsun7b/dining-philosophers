package main

import (
	"fmt"
	"sync"
)

type ChopStick struct { sync.Mutex }

type Philosopher struct {
	leftCS, rightCS *ChopStick
	index int
}

type Host struct {
	count int
	m sync.Mutex
}

func (host Host) getPermission() bool {
	result := false
	host.m.Lock()

	if host.count < 2 {
		host.count++
		result = true
	}

	host.m.Unlock()
	return result
}

func (host Host) release() {
	host.m.Lock()
	host.count--
	host.m.Unlock()
}

func (philosopher Philosopher) eat(host *Host, group *sync.WaitGroup)  {
	for i := 0; i < 3; i++ {
		if host.getPermission() {

			philosopher.leftCS.Lock()
			philosopher.rightCS.Lock()

			fmt.Printf("Started eating \t\t<%d>\n", philosopher.index)
			fmt.Printf("Finishing eating \t<%d>\n", philosopher.index)

			philosopher.leftCS.Unlock()
			philosopher.rightCS.Unlock()

			host.release()
		}
	}

	group.Done()
}

func main() {
	chopsticks := make([]*ChopStick, 5)

	for i := 0; i < 5; i++  {
		chopsticks[i] = new(ChopStick)
	}

	philosophers := make([]*Philosopher, 5)

	for i := 0; i < 5; i++  {
		philosophers[i] = &Philosopher{chopsticks[i], chopsticks[(i + 1) % 5], i + 1}
	}

	host := new(Host)

	group := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		group.Add(1)
		go philosophers[i].eat(host, &group)
	}

	group.Wait()
}
