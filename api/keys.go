package api

import (
  "fmt"
  "net/http"
)

type point struct {
  y int
  x int
}

type step struct {
  a int
  b int
  c int
}

const (
  HEIGHT = 6
  WIDTH  = 8
)

var theMap = [HEIGHT][WIDTH]int{
  {-1, -1, -1, -1, -1, -1, -1, -1},
  {-1,  0,  0,  0,  0,  0,  0, -1},
  {-1,  0, -1, -1, -1,  0,  0, -1},
  {-1,  0,  0,  0, -1,  0, -1, -1},
  {-1,  1, -1,  0,  0,  0,  0, -1},
  {-1, -1, -1, -1, -1, -1, -1, -1},
}
var start = point{4, 1}

func KeySimulation(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Find A Joni's HomeKey\n\n", mapString(theMap), "\n\n")
  resultMap := simulate()
  fmt.Fprint(w, "Finished!\n", mapString(resultMap))
  fmt.Printf("GET /keys/simulation - StatusCode: %v", http.StatusOK)
}

func simulate() [HEIGHT][WIDTH]int {
  steps := make(map[point]step)
  queue := []point{}
  resultMap := copyMap(theMap)

  push(start, &queue, &resultMap)

  for len(queue) > 0 {
    var next point
    currentPoint := pop(&queue)
    currentStep := steps[currentPoint];

    // move north
    next = point{currentPoint.y - 1, currentPoint.x}
    if push(next, &queue, &resultMap) {
      steps[next] = step{currentStep.a + 1, currentStep.b, currentStep.c}
    }

    // move east
    next = point{currentPoint.y, currentPoint.x + 1}
    if push(next, &queue, &resultMap) {
      steps[next] = step{currentStep.a, currentStep.b + 1, currentStep.c}
    }

    // move south
    next = point{currentPoint.y + 1, currentPoint.x}
    if push(next, &queue, &resultMap) {
      steps[next] = step{currentStep.a, currentStep.b, currentStep.c + 1}
    }

    // move west
    next = point{currentPoint.y, currentPoint.x - 1}
    if push(next, &queue, &resultMap) {
      steps[next] = currentStep
    }
  }

  for key, step := range steps {
    if step.a > 0 && step.b > 0 && step.c > 0 {
      resultMap[key.y][key.x] = 2
    }
  }

  return resultMap
}

func push(point point, queue *[]point, mark *[HEIGHT][WIDTH]int) bool {
  if mark[point.y][point.x] > -1 {
    *queue = append(*queue, point)
    mark[point.y][point.x] = -2
    return true
  }
  return false
}

func pop(queue *[]point) point {
  point := (*queue)[0]
  *queue = (*queue)[1:]
  return point
}

func copyMap(m [HEIGHT][WIDTH]int) [HEIGHT][WIDTH]int {
  cm := [HEIGHT][WIDTH]int{}

  for i := 0; i < HEIGHT; i++ {
    for j := 0; j < WIDTH; j++ {
      cm[i][j] = m[i][j]
    }
  }

  return cm
}

func mapString(m [HEIGHT][WIDTH]int) string {
  var mapString string

  for i := 0; i < HEIGHT; i++ {
    for j := 0; j < WIDTH; j++ {
      key := "."
      if m[i][j] == -1 {
        key = "#"
      } else if m[i][j] == 1 {
        key = "X"
      } else if m[i][j] == -2 {
        key = "x"
      } else if m[i][j] == 2 {
        key = "k"
      }
      mapString += key
    }
    mapString += "\n"
  }

  return mapString
}
