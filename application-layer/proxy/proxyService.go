package proxy

import (
    "fmt"
    "net/http"
)

// Placeholder: Can All Be Replaced
func StartProxyServer() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Proxy server is running!")
    })

    port := "8080"
    fmt.Printf("Proxy server running on http://localhost:%s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Printf("Failed to start server: %v\n", err)
    }
}
