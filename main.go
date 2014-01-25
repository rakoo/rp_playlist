package main

import (
  "bytes"
  "fmt"
  "strconv"
  "io"
  "log"
  "net/http"
  "strings"
  "time"

  "launchpad.net/xmlpath"
)

const (
  RP_URL = "http://www.radioparadise.com/ajax_rp2_playlist.php"
)

var (
  // Works because a#nowplaying_title's first value is the right
  // one...
  current_playing_path = xmlpath.MustCompile("//a[@id='nowplaying_title']/b")
)

func main() {
  nextIn := 0 * time.Second
  for {
    <-time.After(nextIn)
    resp, err := http.Get(RP_URL)
    if err != nil {
      log.Println(err)
      nextIn = 1 * time.Second
      continue
    }

    var b bytes.Buffer
    _, err = io.Copy(&b, resp.Body)
    if err != nil {
      log.Println(err)
      nextIn = 1 * time.Second
      continue
    }
    defer resp.Body.Close()

    parts := strings.Split(string(b.Bytes()), "|")
    if len(parts) != 2 {
      log.Printf("Expected 2 parts, got %d", len(parts))
      nextIn = 1 * time.Second
      continue
    }
    resp.Body.Close()

    root, err := xmlpath.ParseHTML(bytes.NewReader([]byte(parts[1])))
    if err != nil {
      log.Println(err)
      nextIn = 1 * time.Second
      continue
    }

    current, ok := current_playing_path.String(root)
    if !ok {
      log.Println("Couldn't find currently playing")
      nextIn = 1 * time.Second
      continue
    }

    nextTick, err := strconv.Atoi(parts[0])
    if err != nil {
      log.Printf("Couldn't get int value out of %d", parts[0])
      nextIn = 1 * time.Second
      continue
    }

    nextIn = time.Duration(nextTick) * time.Millisecond
    fmt.Println(current)
  }
}
