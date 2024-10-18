package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
    targetURL := "http://localhost:3000" // Change to external API

    // Parse the target URL
    proxyURL, err := url.Parse(targetURL)
    if err != nil {
        http.Error(w, "Invalid target URL", http.StatusInternalServerError)
        return
    }

    // Create a new reverse proxy
    proxy := httputil.NewSingleHostReverseProxy(proxyURL)

    // Modify the request to include the original URL path
    r.URL.Host = proxyURL.Host
    r.URL.Scheme = proxyURL.Scheme
    r.Header.Set("X-Forwarded-Host", r.Header.Get("Host")) // Forward the original host header
    r.Host = proxyURL.Host

    // Log the request URL for monitoring
    fmt.Printf("Proxying request: %s\n", r.URL.String())

    // Serve the HTTP request
    proxy.ServeHTTP(w, r)
}

func main() {
    http.HandleFunc("/", handleProxy)

    port := ":8080"
    fmt.Printf("Starting HTTP proxy server on port %s...\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
