package expor

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Display struct { 
  Name string
  Connected bool
  IsExternal bool
  Options DisplayOptions
}

type Vec2 struct {
  X int
  Y int
}

type DisplayOptions struct {
  Resolution Vec2 // Eg. 1920x1080
  RefreshRate int `default:"60"` // Hz
  Position Vec2 // x,y
  Rotation string `default:"normal"`// normal, left, right, inverted
  Primary bool
  Active bool
}

func DefaultDisplayOptions() DisplayOptions {
  return DisplayOptions{
    Resolution: Vec2{X: 1920, Y: 1080},
    RefreshRate: 60,
    Position: Vec2{X: 0, Y: 0},
    Rotation: "normal",
    Primary: false,
    Active: false,
  }
}

func (d *Display) ToStr() string {
  if d.Connected {
    return fmt.Sprintf("%v connected", d.Name)
  }
  return fmt.Sprintf("%v disconnected", d.Name)
}

func (d *Display) SetResolution(resolution Vec2) *Display { 
  d.Options.Resolution = resolution
  return d
}

func (d *Display) SetRefreshRate(refreshRate int) *Display {
  d.Options.RefreshRate = refreshRate
  return d
}

func (d *Display) SetPosition(position Vec2) *Display {
  d.Options.Position = position
  return d
}

func (d *Display) SetRotation(rotation string) *Display {
  d.Options.Rotation = rotation
  return d
}

func (d *Display) SetPrimary(primary bool) *Display {
  d.Options.Primary = primary
  return d
}

func (d *Display) SetActive(active bool) *Display {
  d.Options.Active = active
  return d
}

func (d *Display) WithOptions(options DisplayOptions) *Display {
  d.Options = options
  return d
}

func ListDisplays() []Display {
  cmd := exec.Command("xrandr", "-q")
  out, err := cmd.Output()
  if err != nil {
    println(err)
  }
  outStr := string(out)
  outLst := strings.Split(outStr, "\n")
  r := regexp.MustCompile(`((eDP|HDMI|VGA|DP|DVI)-\d)\s((connected|disconnected))`)
  displays := []Display{}
  for _, line := range outLst {
    matches := r.FindAllStringSubmatch(line, -1)
    if len(matches) > 0 {
      if (len(matches[0]) < 4) {
        continue
      }
      display := Display{
        Name: matches[0][1],
        Connected: matches[0][3] == "connected",
        IsExternal: !strings.Contains(matches[0][1], "eDP"),
      }      
      display.WithOptions(DefaultDisplayOptions())
      displays = append(displays, display)
    }
  }
  return displays
}

func ListConnectedDisplays() []Display {
  connectedDisplays := []Display{}
  displays := ListDisplays()
  for _, display := range displays {
    if display.Connected {
      connectedDisplays = append(connectedDisplays, display)
    }
  }
  return connectedDisplays
}

func GenerateXrandrCmd(displays []Display) string {
  cmd := "xrandr"
  for _, display := range displays {
    cmd += fmt.Sprintf(" --output %v", display.Name)
    if display.Options.Active {
      cmd += fmt.Sprintf(" --mode %vx%v", display.Options.Resolution.X, display.Options.Resolution.Y)
      cmd += fmt.Sprintf(" --pos %vx%v", display.Options.Position.X, display.Options.Position.Y)
      cmd += fmt.Sprintf(" --rotate %v", display.Options.Rotation)
      cmd += fmt.Sprintf(" --rate %v", display.Options.RefreshRate)
      if display.Options.Primary {
        cmd += " --primary"
      }
    } else {
      cmd += " --off"
    }
  }
  return cmd
}
