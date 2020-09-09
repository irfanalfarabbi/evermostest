package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_initBallContainers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initBallContainers()
		})
	}
}

func TestLoadBall(t *testing.T) {
	ballContainersObj.TotalContainer = 1
	resetContainers()
	t.Log("BallContainers:", ballContainersObj)

	tests := []struct {
    name     string
    status   int
    expected string
  }{}
  for i := 0; i <= MAXBALL; i++ {
		name := "success"
		status := http.StatusOK
		expected := "{\"status\":\"success\",\"message\":\"Successfully loaded a ball into a random ball container.\",\"data\":{\"total_container\":%v,\"loaded_container_no\":%v}}\n"
		if i == MAXBALL {
			name = "failed"
			status = http.StatusUnprocessableEntity
			expected = "{\"status\":\"error\",\"message\":\"Failed to load a ball into a random ball container, as one of the ball container is already full.\",\"data\":{\"total_container\":%v,\"loaded_container_no\":%v}}\n"
		}
		test := struct{
			name     string
    	status   int
    	expected string
		}{ name, status, expected}
		tests = append(tests, test)
	}
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
    	t.Log("BallContainers:", ballContainersObj)
      req, err := http.NewRequest("GET", "/tennis/load_ball", nil)
      if err != nil {
        t.Fatal(err)
      }

      rr := httptest.NewRecorder()
      handler := http.HandlerFunc(LoadBall)

      handler.ServeHTTP(rr, req)

      if status := rr.Code; status != tt.status {
        t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
      }

      expected := fmt.Sprintf(tt.expected, ballContainersObj.TotalContainer, ballContainersObj.VerifiedContainer)
      if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
      }
      t.Log("BallContainers:", ballContainersObj)
    })
  }
}

func Test_loadBallToRandomContainer(t *testing.T) {
  resetContainers()

	tests := []struct {
		name string
		want int
	}{
		{"success", ballContainersObj.TotalContainer - 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadBallToRandomContainer();
			if got < 0 && got > tt.want {
				t.Errorf("loadBallToRandomContainer() = %v, want 0 - %v", got, tt.want)
			}
			if ballContainersObj.containers[got] != 1 {
				t.Errorf("ballContainer #%v is not loaded", got + 1)
			}
		})
	}
}

func Test_checkContainer(t *testing.T) {
	ballContainersObj.TotalContainer = MAXBALL + 1
	resetContainers()

	type args struct {
		container int
	}
	tests := []struct {
		name string
		args args
		want int
	}{}
	for i := 0; i < ballContainersObj.TotalContainer; i++ {
		ballContainersObj.containers[i] = i
		want := 0
		if i == MAXBALL {
			want = MAXBALL + 1
		}
		test := struct{
			name string
			args args
			want int
		}{ fmt.Sprintf("check container #%v", i + 1), args{i}, want }
		tests = append(tests, test)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkContainer(tt.args.container)
			if got := ballContainersObj.VerifiedContainer; got != tt.want {
				t.Errorf("ballContainersObj.VerifiedContainer = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isContainerFull(t *testing.T) {
	ballContainersObj.TotalContainer = MAXBALL + 1
	resetContainers()

	type args struct {
		container int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{}
	for i := 0; i < ballContainersObj.TotalContainer; i++ {
		ballContainersObj.containers[i] = i
		want := false
		if i == MAXBALL {
			want = true
		}
		test := struct{
			name string
			args args
			want bool
		}{ fmt.Sprintf("is container full #%v", i + 1), args{i}, want }
		tests = append(tests, test)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isContainerFull(tt.args.container); got != tt.want {
				t.Errorf("isContainerFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setVerifiedContainer(t *testing.T) {
	resetContainers()

	type args struct {
		container int
	}
	tests := []struct {
		name string
		args args
	}{
		{"container verified", args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setVerifiedContainer(tt.args.container)

			if ballContainersObj.VerifiedContainer != tt.args.container + 1 {
				t.Errorf("ballContainer #%v is not verified", tt.args.container + 1)
			}
		})
	}
}

func Test_resetContainers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"container reset"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetContainers()
		})
	}
}
