package expor

import (
	"flag"
	"fmt"
)

func RunCli() {
  flag.String("lc", "", "List connected displays")
  flag.String("l", "", "List all displays")
  flag.Parse()

  switch {
    case flag.Arg(0) == "lc": {
      displays := ListConnectedDisplays()
      for _, display := range displays {
        fmt.Println(display.Name)
      }
    }
    case flag.Arg(0) == "l": {
      displays := ListDisplays()
      for _, display := range displays {
        fmt.Println(display.ToStr())
      }
    }
  }
}
