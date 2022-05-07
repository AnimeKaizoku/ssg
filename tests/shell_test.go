package tests

import (
	"os"
	"sync"
	"testing"
	"time"

	ws "github.com/AnimeKaizoku/ssg/ssg"
)

func TestShell01(t *testing.T) {
	result := ws.RunCommand("go version")
	if result == nil {
		t.Error("result is nil")
		return
	}

	if result.Error != nil {
		t.Error(result.Error)
		return
	}

	if result.Stdout == "" {
		t.Error("unexpected empty stdout string from result")
		return
	}
}

func TestShellAsync01(t *testing.T) {
	if os.PathSeparator != '/' {
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	result := ws.RunCommandAsync("sleep 100 && echo hello")
	if result == nil {
		t.Error("result is nil")
		return
	}

	result.WaitAndRun(time.Second, 3*time.Second, func(r *ws.ExecuteCommandResult) {
		if r.IsDone() {
			t.Error("unexpected true value returned from IsDone method")
			return
		}

		// kill the process
		err := r.Kill()
		if err != nil {
			t.Error("when tried to kill the proess: ", err)
			return
		}

		//log.Println("killed the proccess")
		//time.Sleep(300 * time.Second)

		wg.Done()
	})

	wg.Wait()
}
