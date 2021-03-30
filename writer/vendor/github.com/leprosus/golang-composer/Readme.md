# Script helps to pause goroutins synchronously

## Usage

```go
package main

import(
	composer "github.com/leprosus/golang-composer"
	"fmt"
	"time"
	"os"
)

func main() {
	for i := 1; i <= 3; i++ {
        go loop(i)
    }
 
    fmt.Println("All of goroutins are started")
 
    time.Sleep(5 * time.Second)
    composer.GetComposer().Pause()
 
    fmt.Println("All of goroutins are paused")
 
    time.Sleep(5 * time.Second)
    composer.GetComposer().Play()
 
    fmt.Println("All of goroutins are resumed")
 
    time.Sleep(5 * time.Second)
}

func loop(id int){
	for {
        fmt.Printf("Goroutin #%d\n", id)

        time.Sleep(time.Second)

        composer.GetComposer().NeedWait()
    }
}
```

## List all methods

* composer.GetComposer - returns composer
* composer.GetComposer.Play - lets to execute goroutins
* composer.GetComposer.Pause - lets to pause goroutins execution
* composer.GetComposer.NeedWait - if need to pause then wait resuming rather do nothing 