package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"github.com/miekg/dns"
)

var (
	p string // Define variable for protocol type
	mu sync.Mutex // Mutex for synchronization
)

func init() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go [tcp/udp]")
		os.Exit(1)
	}

	// Check if the argument is either "tcp" or "udp"
	protocol := strings.ToLower(os.Args[1])
	if protocol != "tcp" && protocol != "udp" {
		fmt.Println("Invalid protocol specified. Please use 'tcp' or 'udp'.")
		os.Exit(1)
	}

	p = protocol
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	mu.Lock() // Lock to ensure that only one goroutine is writing at a time
	defer mu.Unlock()

	var client *dns.Client

	// Choose the protocol based on user configuration
	switch p {
	case "tcp":
		client = &dns.Client{Net: "tcp-tls"}
	case "udp":
		client = &dns.Client{Net: "udp"}
	default:
		log.Printf("Invalid protocol specified: %s. Using default protocol (tcp)\n", p)
		client = &dns.Client{Net: "tcp-tls"}
	}

	// Use standard DNS UDP port (53) and TCP port (853)
	var serverAddr string
	switch client.Net {
	case "udp":
		serverAddr = "1.1.1.1:53"
	case "tcp-tls":
		serverAddr = "1.1.1.1:853"
	}

	// Exchange the DNS query with the DNS server
	reply, _, err := client.Exchange(r, serverAddr)
	if err != nil {
		log.Println("Error occurred querying DNS server:", err)
		return
	}

	// Log all answers received from the DNS server
	for _, answer := range reply.Answer {
		log.Printf("Answer: %s", answer.String())
	}

	// Send the DNS response back to the client
	w.WriteMsg(reply)
}

func main() {
	// Handle DNS queries using the handleRequest function
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		// Launch a new goroutine for each incoming request
		go handleRequest(w, r)
	})

	// Set up the DNS daemon to listen on the specified protocol and port
	server := &dns.Server{
		Addr:    ":53",
		Net:     p,
		Handler: dns.DefaultServeMux,
	}

	// Start the DNS server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Failed to set up DNS server: %v\n", err)
			os.Exit(1)
		}
	}()

	log.Printf("DNS proxy listening on :53 (%s)\n", p)

	// CLI command output based on the protocol
	var cliCommand string
	switch p {
	case "tcp":
		cliCommand = fmt.Sprintf("CLI: dig +short +tcp google.com @localhost\n")
	case "udp":
		cliCommand = fmt.Sprintf("CLI: dig +short udp google.com @localhost\n")
	default:
		cliCommand = fmt.Sprintf("Invalid protocol specified: %s\n", p)
	}

	log.Print(cliCommand)

	// Blocking to keep the app running indefinitely.
	select {}
}