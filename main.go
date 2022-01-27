package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"

	"main/src/AudioPlayer"
	"main/src/ConfigLoader"
	"main/src/EventManager"
)

const (
	SS_VERSION = "1.0.1"

	GUI_HEIGHT = 360
	GUI_WIDTH  = 540
)

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Check(Close chan bool, Error chan error, Events *EventManager.EventList, Table *widget.Table, Player *AudioPlayer.AudioPlayer, cfg *ConfigLoader.Config, requestFocus func(), DisplayError func(error), EventPassed func(*EventManager.Event)) {
	for {
		select {
		case <-Close:
			os.Exit(4)
			Player.Close()
		default:
			Passed := Events.PassedEvents()
			if len(Passed) > 0 {
				err := Events.RemovePassed()
				if err != nil {
					DisplayError(err)
				}
				Table.Refresh()
				for _, v := range Passed {
					if time.Now().UnixMilli() >= v.UnixTimestamp {
						if cfg.Instructive && strings.HasPrefix(v.Title, "i:") {
							if strings.HasPrefix(v.Title, "i:OPEN[") && v.Title[len(v.Title)-1] == ']' {
								v.Title = v.Title[7 : len(v.Title)-1]
								fmt.Println(v.Title)
								browser.OpenURL(v.Title)
							} else {
								DisplayError(fmt.Errorf("Error analyzing Instructive Command"))
							}
						} else {
							EventPassed(v)
						}
						Player.Play()
					}
				}
			}
		}
		time.Sleep(time.Second * cfg.CheckIntervalDuration())
	}
}

func main() {
	Events, err := EventManager.ReadEvents("./schedules.json")
	if err != nil {
		os.Exit(3)
	}

	Config, err := ConfigLoader.ReadConfig("./config.json")
	if err != nil {
		os.Exit(3)
	}

	Player, err := AudioPlayer.New(Config.AudioPath)
	if err != nil {
		os.Exit(3)
	}

	Application := app.New()
	Window := Application.NewWindow("SimpleSchedule")
	Window.Resize(fyne.NewSize(GUI_WIDTH, GUI_HEIGHT))
	Window.SetFixedSize(true)

	// Header
	GUI_LOGO := canvas.NewImageFromFile("./logos/noBG-1000x1000-SimpleSchedule.png")
	GUI_LOGO.SetMinSize(fyne.NewSize(40, 40))
	GUI_LOGO.Resize(fyne.NewSize(40, 40))
	GUI_LOGO.FillMode = canvas.ImageFillContain

	GUI_TABLE := widget.NewTable(func() (int, int) {
		return len(Events.Events) + 1, 2
	}, func() fyne.CanvasObject {
		return widget.NewLabel("---- Unable To Load ----")
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		label := cell.(*widget.Label)

		if id.Row == 0 {
			switch id.Col {
			case 0:
				label.SetText("Title")
			case 1:
				label.SetText("Time")
			}
			return
		}
		Evnt, _ := Events.GetEvent(id.Row - 1)
		switch id.Col {
		case 0:
			if strings.HasPrefix(Evnt.Title, "i:") && Config.Instructive {
				label.SetText("(Instructive)")
			} else {
				label.SetText(Evnt.Title)
			}
		case 1:
			dat := time.UnixMilli(Evnt.UnixTimestamp)
			_, month, day := dat.Date()
			hour, min, sec := dat.Clock()
			var min_str string = fmt.Sprintf("%v", min)
			var suffix string = "A.M."
			if hour > 12 {
				hour -= 12
				suffix = "P.M."
			}
			if min < 10 {
				min_str = fmt.Sprintf("0%v", min)
			}
			label.SetText(fmt.Sprintf("%v %v   at   %v:%v:%v %v", month.String(), day, hour, min_str, sec, suffix))

		}
	})
	GUI_ADD_BUTTON := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		title := widget.NewEntry()
		unix := widget.NewEntry()
		unix.Validator = validation.NewRegexp(`^[0-9]*$`, "Unix Timestamp ONLY have Numbers")
		items := []*widget.FormItem{
			widget.NewFormItem("Title", title),
			widget.NewFormItem("Unix Timestamp", unix),
		}

		form := dialog.NewForm("Add Event", "Add Event", "Cancel", items, func(b bool) {
			if b {
				i, err := strconv.Atoi(unix.Text)
				if err != nil {
					dialog.ShowError(err, Window)
					return
				}
				Events.AddEvent(title.Text, int64(i))
				GUI_TABLE.Refresh()
			}
		}, Window)
		form.Resize(fyne.NewSize(GUI_WIDTH, GUI_HEIGHT))

		form.SetDismissText("Nevermind")
		form.Show()

	})
	GUI_SETTING_BUTTON := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {

		NotificationSoundLabel := widget.NewEntry()
		NotificationSoundLabel.SetText(Config.AudioPath)

		Instructive := widget.NewCheck("", func(b bool) {
			Config.Instructive = b
			err := Config.Save()
			if err != nil {
				dialog.ShowError(err, Window)
			}
		})

		Instructive.Checked = Config.Instructive

		ConfigForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "Sound", Widget: container.NewAdaptiveGrid(2, NotificationSoundLabel, widget.NewButton("Set", func() {
					Config.AudioPath = NotificationSoundLabel.Text
					err := Config.Save()
					if err != nil {
						dialog.ShowError(err, Window)
					}
				}))},
				{Text: "Instructive", Widget: Instructive},
			},
		}

		InfoTab := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "Version", Widget: widget.NewLabel(SS_VERSION)},
			},
		}

		ConfigForm.Resize(fyne.NewSize(GUI_WIDTH-300, GUI_HEIGHT))

		SettingsTab := container.NewVBox(ConfigForm, widget.NewSeparator(), InfoTab)

		ConfigDialog := dialog.NewCustom("Config & Info", "Close", container.NewAdaptiveGrid(2, SettingsTab), Window)
		ConfigDialog.Resize(fyne.NewSize(GUI_WIDTH-300, GUI_HEIGHT))
		ConfigDialog.Show()

	})
	GUI_HEADER := container.NewAdaptiveGrid(2, container.NewHBox(GUI_LOGO), container.NewHBox(layout.NewSpacer(), GUI_ADD_BUTTON, GUI_SETTING_BUTTON))

	Window.SetContent(container.NewBorder(GUI_HEADER, nil, nil, nil, GUI_TABLE))
	var CloseChannel chan bool = make(chan bool)
	var ErrorChannel chan error = make(chan error)

	go Check(CloseChannel, ErrorChannel, Events, GUI_TABLE, Player, Config, Window.RequestFocus, func(e error) {
		dialog.ShowError(e, Window)
	}, func(e *EventManager.Event) {
		dialog.ShowInformation("Event", fmt.Sprintf("Looks like you have an event named \"%v\".", e.Title), Window)
	})

	Window.SetOnClosed(func() {
		CloseChannel <- true
	})

	Window.ShowAndRun()
}
