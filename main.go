package main

import (
	"fmt"
	"os/exec"
	"sort"

	"github.com/sudo-nick/expor/expor"
)

func main() {
  startPos := expor.Vec2{X: 0, Y: 0}
  res := expor.Vec2{X: 1920, Y: 1080}
  displays := expor.ListConnectedDisplays()
  sort.Slice(displays, func(i, j int) bool {
    return displays[i].IsExternal
  })
  for i := 0; i < len(displays); i++ {
    displays[i].
      SetPosition(startPos).
      SetResolution(res).
      SetActive(true)
    if displays[i].IsExternal {
      displays[i].SetRefreshRate(100)
    } else {
      displays[i].SetPrimary(true)
    }
    startPos.X += res.X
  }
  xrandrCmd := expor.GenerateXrandrCmd(displays)
  fmt.Println("Running", xrandrCmd)
  exec.Command("bash", "-c", xrandrCmd).Run()
}
