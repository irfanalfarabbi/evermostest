package api

import (
  "encoding/json"
  "fmt"
  "math/rand"
  "net/http"
  "sync"
  "time"
)

const (
  MINSTOCK = 2
  MAXSTOCK = 5
)

type Product struct {
  Name  string `json:"name"`
  Stock int    `json:"stock_remaining"`
}

var product Product

var mutex = &sync.Mutex{}

func init() {
  initProduct()

}

func initProduct() {
  product = Product{
    Name: "My Product",
  }
  restock()
  fmt.Println("INIT InitProduct - Product:", product)
}

func Order(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := ApiResponse{
    Status: "success",
    Message: "Successfully order the product.",
  }

  if !reduceStock() {
    statusCode = http.StatusUnprocessableEntity
    apiResponse.Status = "error"
    apiResponse.Message = "Failed to order the product, stock is not enough"
  }

  apiResponse.Data = product

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /shops/order - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}

func MassOrder(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := ApiResponse{
    Status: "success",
    Message: "Successfully mass order.",
  }

  orders := []string{}
  var wg sync.WaitGroup
  wg.Add(MAXSTOCK * 2)


  for i := 0; i < MAXSTOCK * 2; i++ {
    go func() {
      if !reduceStock() {
        orders = append(orders, "Failed to order the product, stock is not enough")
      } else {
        orders = append(orders, "Successfully order the product.")
      }
      wg.Done()
    }()
  }
  wg.Wait()

  apiResponse.Data = struct{
    Product Product  `json:"product"`
    Orders  []string `json:"orders"`
  }{ product, orders}

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /shops/order - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}

func reduceStock() bool {
  mutex.Lock()
  if product.Stock > 0 {
    product.Stock--
    mutex.Unlock()
    return true
  } else {
    mutex.Unlock()
    return false
  }
}

func restock() {
  rand.Seed(time.Now().UnixNano())

  total := MINSTOCK
  if MINSTOCK < MAXSTOCK {
    total = rand.Intn(1 + MAXSTOCK - MINSTOCK) + MINSTOCK
  }

  product.Stock = total
}
