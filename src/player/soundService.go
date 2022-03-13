package player

import (
    "bytes"
    "encoding/binary"
    "fmt"

    "github.com/gordonklaus/portaudio"
    "github.com/bobertlo/go-mpg123/mpg123"

    "sync"
)

type SoundPlayer struct {
    control         int
    filepath        string
    controller      chan int
    done            chan error
    lock            sync.Mutex
}

func (player *SoundPlayer) WaitDone () error {
    for player.control == CONTROLLER_PLAY {
        select {
            case err:= <- player.done:
                return err
            default:
                continue
        }
    }

    return nil
}

func NewSoundPlayer (filepath string) (SoundPlayer) {
    player := SoundPlayer {
    }

    player.init (filepath)
    
    return player
}

func (player *SoundPlayer) init (filepath string) {
    player.control = CONTROLLER_STOP
    player.filepath = filepath
    player.controller = make(chan int)
    player.done = make(chan error)

}

func (player *SoundPlayer) Set (filepath string) {
    if player.control == CONTROLLER_STOP {
        player.Stop()
    }
    player.lock.Lock()
    player.filepath =  filepath
    player.lock.Unlock()

}

func (player *SoundPlayer) setControl (control int) {
    player.lock.Lock()
    player.control = control
    player.controller <- player.control

    if control == CONTROLLER_STOP {
        player.done = make(chan error)
    }
    player.lock.Unlock()
}

func (player *SoundPlayer) Start () {
    if player.control == CONTROLLER_STOP {
        decoder := OpenDecoder(player.filepath)
        go PlaySound (decoder, player.controller, player.done)
        player.setControl(CONTROLLER_PLAY)
        
    } else if player.control == CONTROLLER_PAUSE {
        player.setControl (CONTROLLER_PLAY)
    }
}

func (player *SoundPlayer) Pause () {
    if player.control == CONTROLLER_PLAY {
        player.setControl (CONTROLLER_PAUSE)
    }
}

func (player *SoundPlayer) Stop () {
    if player.control != CONTROLLER_STOP {
        player.setControl (CONTROLLER_STOP)
       
    }
}

func (player *SoundPlayer) Restart () {
    player.Stop()
    player.Start()
}

func (player *SoundPlayer) IsPlay () bool {
    return player.control == CONTROLLER_PLAY
}

func (player *SoundPlayer) IsStop () bool {
    return player.control == CONTROLLER_STOP
}

func (player *SoundPlayer) IsPause () bool {
    return player.control == CONTROLLER_PAUSE
}


const OUT_BUFFER_SIZE = 8192

const CONTROLLER_PLAY = 0
const CONTROLLER_PAUSE = 1
const CONTROLLER_STOP = 2

func PlaySound (decoder *mpg123.Decoder, controller chan int, done chan error) {
    control := CONTROLLER_STOP

    // get audio format information
    rate, channels, _ := decoder.GetFormat()

    // make sure output format does not change
    decoder.FormatNone()
    decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

    portaudio.Initialize()
    defer portaudio.Terminate()
    out := make([]int16, OUT_BUFFER_SIZE)

    fmt.Printf("%d - %d" , channels, rate)
    stream, err := portaudio.OpenDefaultStream (0, channels, float64(rate), len(out), &out)
    handleError(err)

    defer stream.Close()

    err = stream.Start()
    handleError(err)

    for {
        select {
            case newSignal := <- controller:
                control = newSignal
            default:
                //Do nothing
        }

        audio := make([]byte, 2*len(out))
        if control == CONTROLLER_PLAY {
            _, err = decoder.Read(audio)

            if err == mpg123.EOF {
                done <- nil
                return
            }

        } else if control == CONTROLLER_STOP {
            done <- nil
            return
        }

        handleError(err)

        err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out)
        handleError(err)

        err = stream.Write()
        handleError(err)

    }

}

func OpenDecoder (filename string) (*mpg123.Decoder) {
    //TODO: make customize decoder -> no control over decoder
    //TODO: research c-mpg123 package
    decoder, err := mpg123.NewDecoder ("")
    handleError(err)  
    //err = decoder.OpenFeed()
    err = decoder.Open(filename)
    handleError(err)

    return decoder
}

func handleError (err error) {
    if err != nil {
        panic(err)
    }
}
