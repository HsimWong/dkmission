package utils

import "sync"

func ThreadBlock() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
