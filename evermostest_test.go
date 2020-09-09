package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
)

func Test_healthz(t *testing.T) {
  tests := []struct {
    name     string
    status   int
    expected string
  }{
    {"success", http.StatusOK, "{\"status\":\"success\",\"message\":\"OK\",\"data\":null}\n"},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      req, err := http.NewRequest("GET", "/healthz", nil)
      if err != nil {
        t.Fatal(err)
      }

      rr := httptest.NewRecorder()
      handler := http.HandlerFunc(healthz)

      handler.ServeHTTP(rr, req)

      if status := rr.Code; status != tt.status {
        t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
      }

      if rr.Body.String() != tt.expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expected)
      }
    })
  }
}
