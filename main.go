package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type app struct {
	debounce debouncer
	cmd      string
}

var debounceTime = time.Second

func main() {
	a := app{
		debounce: debounce(debounceTime),
		cmd:      os.Getenv("FLOAT_CMD"),
	}

	log.Println("float started")

	log.Fatal(http.ListenAndServe("0.0.0.0:41232", http.HandlerFunc(a.handler)))
}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received")
	a.debounce(a.run)
}

func (a *app) run() {
	cmd := exec.Command("/bin/sh", "-c", a.cmd)
	cmd.Stdout = log.Writer()
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

type debouncer func(f func())

// credit to https://github.com/bep/debounce
func debounce(after time.Duration) debouncer {
	d := &struct {
		mu    sync.Mutex
		after time.Duration
		timer *time.Timer
	}{
		after: after,
	}

	return func(f func()) {
		d.mu.Lock()
		defer d.mu.Unlock()

		if d.timer != nil {
			d.timer.Reset(d.after)
			return
		}

		d.timer = time.AfterFunc(d.after, f)
	}
}
