package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type WorkerPool struct {
	workers   int
	jobs      chan string
	waitGroup sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		workers: numWorkers,
		jobs:    make(chan string),
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.waitGroup.Add(1)
		go p.worker(i)
	}
}

func (p *WorkerPool) worker(id int) {
	defer p.waitGroup.Done()
	for url := range p.jobs {
		log.Printf("Worker %d: Starting request to %s\n", id, url)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Worker %d: Error while requesting %s: %v\n", id, url, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		log.Printf("Worker %d: Response from %s:\n%s\n", id, url, body)
	}
}

func (p *WorkerPool) Wait() {
	p.waitGroup.Wait()
}

func generateUrls(n int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < n; i++ {
			ch <- "http://example.com/" + string(i) // Пример URL
		}
		close(ch)
	}()
	return ch
}

func main() {
	pool := NewWorkerPool(3)
	pool.Start()
	for url := range generateUrls(5) {
		log.Println("Start request to:", url)
		pool.jobs <- url
	}
	close(pool.jobs)
	pool.Wait()
	log.Println("All requests completed.")
}
