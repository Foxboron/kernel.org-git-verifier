package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

var (
	Signed   = 0
	Unsigned = 0
)

func ScrapeTLog() {
	r, err := GetRepo()
	if err != nil {
		log.Fatal(err)
	}
	r.Pull()
	t, err := LastTimestamp()
	if err != nil {
		// This should be the start of the log.. or something
		t = time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC)
	}
	counter := 0
	var segment []TLogCommit
	r.LogSince(&t,
		func(c *object.Commit, r io.Reader) error {
			files, err := ParseMail(r)
			if err != nil {
				log.Fatal(err)
			}
			checksum := fmt.Sprintf("%s", c.Hash)
			tlog := WorkTLog(checksum, files)
			tlog.CommitMsg = c.Message
			log.Debug("added commit")
			log.WithFields(log.Fields{
				"message":  c.Message,
				"checksum": checksum,
			}).Info("Added commit")
			tlog.CommitDate = c.Author.When
			segment = append(segment, *tlog)
			counter += 1
			if counter%10 == 0 {
				AddCommits(segment)
				segment = []TLogCommit{}
			}
			return nil
		},
	)
	AddCommits(segment) // Add remaining
	t = time.Now()
	AddScan(&Scan{t, counter})
	log.WithFields(log.Fields{
		"timestamp": t,
		"counter":   counter,
	}).Info("scan completed")
}

func CancelSignal(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalCh:
			fmt.Println("Got ^C")
			cancel()
		case <-ctx.Done():
			return
		}
	}()
}

func Timers(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		ScrapeTLog()
		c := time.Tick(1 * time.Hour)
		for {
			select {
			case <-c:
				ScrapeTLog()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	CancelSignal(ctx, cancel)
	Timers(ctx, cancel)
	webpage(ctx)
	<-ctx.Done()
	// CheckCommit(os.Args[1])
}
