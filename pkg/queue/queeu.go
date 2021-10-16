package queue

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

// FreeSpaceWatcher is an endless loop checking how many free space's left
func FreeSpaceWatcher(unlocked chan bool) {
	// for {
		
		var stat unix.Statfs_t
		
		wd, err := os.Getwd()
		//TODO: handle error better
		if err != nil {
			log.Fatal(err)
		}

		unix.Statfs(wd, &stat)
		// if stat.Bavail * uint64(stat.Bsize) < 74532102144 {
		if stat.Bavail * uint64(stat.Bsize) < 745321021499 {
			fmt.Println("Writing to channel")
			unlocked <- true
		}

	// }
}
