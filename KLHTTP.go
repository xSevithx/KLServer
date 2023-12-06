package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//"golang.org/x/crypto/acme/autocert"
)

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Log the received data to a file on the server
	logToFile(body)

	// You can perform additional processing or storage here

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully"))
}

func logToFile(data []byte) {
	// Create or open the log file
	filePath := "server_dump.tmp"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	defer file.Close()

	// Write the received data to the log file
	if _, err := file.Write(data); err != nil {
		log.Fatal("Error writing to log file: ", err)
	}
}

func main() {
	port := "80" // Default port

	// Parse command line arguments
	for i := 1; i < len(os.Args); i += 2 {
		switch os.Args[i] {
		case "-p":
			if i+1 < len(os.Args) {
				port = os.Args[i+1]
			} else {
				fmt.Println("Usage: ./KLHTTP -p port")
				os.Exit(1)
			}
		default:
			fmt.Println("Usage: ./KLHTTP -p port")
			os.Exit(1)
		}
	}
	// Create a new autocert manager
	//certManager := &autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	Cache:      autocert.DirCache("certs"),
	//	HostPolicy: autocert.HostWhitelist("kltester.com"), // Add your domain here
	//}

	// Register the logHandler for the "/log" endpoint
	http.HandleFunc("/log", logHandler)
	// Start the HTTP server to handle ACME challenges
	//go func() {
	//	err := http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	//	if err != nil {
	//		log.Fatal("Error starting HTTP server: ", err)
	//	}
	//}()

	// Start the HTTPS server with automatic TLS certificate management
	//httpsServer := &http.Server{
	//	Addr:      ":443",
	//	TLSConfig: certManager.TLSConfig(),
	//}

	//err := httpsServer.ListenAndServeTLS("", "") // Empty paths for automatic certificate management
	fmt.Printf("Starting Server on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("Error starting HTTPS server: ", err)
	}
}
