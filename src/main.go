package main

import (
    "pomodoro/player"
    //"github.com/bobertlo/go-mpg123/mpg123"
    //"os"
    "time"
    "fmt"
)

func main () {
    filename := "NgayDauTien-DucPhuc.mp3"
    player := player.NewSoundPlayer(filename)
    player.Start()
    time.Sleep(2 * time.Second)
    player.Pause()
    fmt.Println("Waiting")
    player.WaitDone ()
    fmt.Println("Done")

    time.Sleep(2 * time.Second)
    player.Start()

    time.Sleep(2 * time.Second)
    player.Stop()

    time.Sleep(2 * time.Second)
    player.Start()

    time.Sleep(1 * time.Second)
    player.Restart()

    fmt.Println("Waiting")
    player.WaitDone ()
    fmt.Println("Done")
}

func handleError (err error) {
    if err != nil {
        panic(err)
    }
}
