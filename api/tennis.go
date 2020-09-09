package api

import (
  "encoding/json"
  "fmt"
  "math/rand"
  "net/http"
  "time"
)

const (
  MINCONTAINER        = 2
  MAXCONTAINER        = 3
  MAXBALL             = 3
  NOVERIFIEDCONTAINER = 0
)

type ballContainers struct {
  TotalContainer    int   `json:"total_container"`
  VerifiedContainer int   `json:"loaded_container_no"`
  containers        []int
}

var ballContainersObj ballContainers

func init() {
  initBallContainers()
  resetContainers()
}

func initBallContainers() {
  rand.Seed(time.Now().UnixNano())

  total := MINCONTAINER
  if MINCONTAINER < MAXCONTAINER {
    total = rand.Intn(1 + MAXCONTAINER - MINCONTAINER) + MINCONTAINER
  }

  ballContainersObj = ballContainers{
    TotalContainer: total,
  }

  fmt.Println("INIT InitBallContainers - BallContainers:", ballContainersObj)
}

func LoadBall(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := ApiResponse{
    Status: "success",
    Message: "Successfully loaded a ball into a random ball container.",
  }

  if ballContainersObj.VerifiedContainer > 0 {
    statusCode = http.StatusUnprocessableEntity
    apiResponse.Status = "error"
    apiResponse.Message = "Failed to load a ball into a random ball container, as one of the ball container is already full."
  } else {
    container := loadBallToRandomContainer()
    checkContainer(container)
  }

  apiResponse.Data = ballContainersObj

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /tennis/load_ball - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}

func LoadBalls(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := ApiResponse{
    Status: "success",
    Message: "Successfully loaded some balls until one of ball container is full.",
  }

  if ballContainersObj.VerifiedContainer > 0 {
    statusCode = http.StatusUnprocessableEntity
    apiResponse.Status = "error"
    apiResponse.Message = "Failed to load some balls into ball containers, as one of the ball container is already full."
  } else {
    for ballContainersObj.VerifiedContainer == 0 {
      container := loadBallToRandomContainer()
      checkContainer(container)
    }
  }

  apiResponse.Data = ballContainersObj

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /tennis/load_balls - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}

func DumpBalls(w http.ResponseWriter, r *http.Request) {
  statusCode := http.StatusOK
  apiResponse := ApiResponse{
    Status: "success",
    Message: "Successfully dump all balls from the ball containers.",
  }

  resetContainers()

  apiResponse.Data = ballContainersObj

  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(apiResponse)
  fmt.Printf("GET /tennis/dump_balls - StatusCode: %v, ApiResponse: %v\n", statusCode, apiResponse)
}

func loadBallToRandomContainer() int {
  container := rand.Intn(ballContainersObj.TotalContainer)
  ballContainersObj.containers[container]++

  return container
}

func checkContainer(container int) {
  if isContainerFull(container) {
    setVerifiedContainer(container)
  }
}

func isContainerFull(container int) bool {
  return ballContainersObj.containers[container] >= MAXBALL
}

func setVerifiedContainer(container int) {
  ballContainersObj.VerifiedContainer = container + 1
}

func resetContainers() {
  ballContainersObj.containers = make([]int, ballContainersObj.TotalContainer)
  ballContainersObj.VerifiedContainer = NOVERIFIEDCONTAINER
}
