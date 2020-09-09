package main

import (
  "encoding/json"
  "fmt"
  "net/http"

  "./api"
)

func main() {
  http.HandleFunc("/healthz", healthz)

  http.ListenAndServe(":8080", nil)
}

func healthz(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := api.ApiResponse{
    Status: "success",
    Message: "OK",
  }

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /healthz - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}
