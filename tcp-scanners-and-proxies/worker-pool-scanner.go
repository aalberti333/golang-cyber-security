package main

// port scanner using worker pools
// fully optimized example can be found at the url below
// https://github.com/blackhat-go/bhg/blob/master/ch-2/scanner-port-format/
import (
	"fmt"
	"net"
	"sort"
)

func worker(ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	// 100 values per channel, so we should have 11 channels
	// 1024/100 = 10.24, so 100 values in 10 channels, and 24 in the 11th
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	// Create the workers. This won't do anything until we pass a value
	// into the "ports" channel
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// Pass value into "ports" channel concurrently
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// whenever results is assigned a value in worker(), append the results
	// to our slice of ints "openports"
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// close channels to clean up
	close(ports)
	close(results)

	//sort the openport results
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
