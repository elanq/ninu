package ninu

import (
	"fmt"
	"os/exec"
	"time"
)

func ProcessURL(url string, filename string) error {
	execCommand := exec.Command("youtube-dl", "--format", "mp4", "-o", filename, url)
	fmt.Println("Running youtube-dl")

	startTime := time.Now()
	if err := execCommand.Start(); err != nil {
		return err
	}
	fmt.Println("Waiting youtube-dl to complete...")
	if err := execCommand.Wait(); err != nil {
		return err
	}
	finishTime := time.Since(startTime).Milliseconds()
	fmt.Printf("Download completed in %d milliseconds\n", finishTime)

	return nil
}
