// Program: Download 100 images concurrently and save them to a local folder.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	totalImages = 100
	folder      = "images"
	maxParallel = 10
)

func main() {

	if err := os.MkdirAll(folder, 0755); err != nil {
		fmt.Println("Could not create folder:", err)
		return
	}

	start := time.Now()

	var wg sync.WaitGroup                     // waits for all goroutines to finish
	limit := make(chan struct{}, maxParallel) // limits how many run at once

	for i := 1; i <= totalImages; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			limit <- struct{}{}        // take a slot (waits if 10 are busy)
			defer func() { <-limit }() // free the slot when done

			url := fmt.Sprintf("https://picsum.photos/200/300?random=%d", id)
			fileName := filepath.Join(folder, fmt.Sprintf("image_%d.jpg", id))

			if err := downloadImage(url, fileName); err != nil {
				fmt.Printf("Image %d failed: %v\n", id, err)
				return
			}
			fmt.Printf("Image %d saved -> %s\n", id, fileName)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nDone! Downloaded %d images in %v\n", totalImages, time.Since(start))
}

func downloadImage(url, fileName string) error {
	client := http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
