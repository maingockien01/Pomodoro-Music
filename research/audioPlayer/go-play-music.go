
/*
 * References: https://gist.github.com/declank/7ed8926383bd971cd1b7
 */

package main

import (
    "github.com/gordonklaus/portaudio"
    "encoding/binary"
    "log"
    "io"
    "os"
    "os/exec"
)

func getAudioFileArg() (filename string) {
    if len(os.Args) < 2 {
        log.Fatal("Missing argument: input file name")
    }
    filename = os.Args[1]
    return
}

func chk(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func createFfmpegPipe(filename string) (output io.ReadCloser) {
    cmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-")
    output, err := cmd.StdoutPipe()
    chk(err)

    err = cmd.Start()
    chk(err)

    return
}

func playAudioFile(filename string) {
    output := createFfmpegPipe(filename)

    portaudio.Initialize()
    defer portaudio.Terminate()
    
    audiobuf := make([]int16, 16384)
    stream, err := portaudio.OpenDefaultStream(0, 2, 44100, len(audiobuf), &audiobuf)
    chk(err)    
    defer stream.Close()
    
    chk(stream.Start())
    defer stream.Stop()

    for err = binary.Read(output, binary.LittleEndian, &audiobuf); err == nil; err = binary.Read(output, binary.LittleEndian, &audiobuf) {
        chk(stream.Write())
    }

    chk(err) 
}

func main() {
    filename := getAudioFileArg()    
    playAudioFile(filename)
}
