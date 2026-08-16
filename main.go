package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

var debounceDefault = 15 * time.Minute

func main() {
	debounceStr := os.Getenv("FLOAT_DEBOUNCE_TIME")
	debounceDur, err := time.ParseDuration(debounceStr)
	if err != nil {
		log.Println("Failed to parse debounce duration; using default of", debounceDefault.String())
		debounceDur = debounceDefault
	}

	c := os.Getenv("FLOAT_CMD")
	debounced := debouncer(debounceDur, func() {
		cmd := exec.Command("/bin/sh", "-c", c)
		cmd.Stdout, cmd.Stderr = log.Writer(), log.Writer()
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	})

	log.Println("float started with debounce duration:", debounceDur.String(), "command:", c)
	log.Fatal(http.ListenAndServe("0.0.0.0:41232",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("request received")
			debounced()
			w.WriteHeader(http.StatusNoContent)
		})),
	)
}

// credit to https://github.com/bep/debounce
func debouncer(after time.Duration, f func()) func() {
	d := &struct {
		mu    sync.Mutex
		after time.Duration
		timer *time.Timer
	}{
		after: after,
	}

	return func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		if d.timer != nil {
			d.timer.Reset(d.after)
			return
		}

		d.timer = time.AfterFunc(d.after, f)
	}
}
