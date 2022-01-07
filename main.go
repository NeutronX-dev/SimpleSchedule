package main

import (
	"fmt"
	"main/src/ConfigLoader"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const (
	SS_VERSION = "1.0.0"
)

func Check(ErrorChannel chan error, CloseChannel chan bool, Schedules *ConfigLoader.Schedules, Table *widget.Table, requestFocus func(), DisplayError func(error), EventPassed func(*ConfigLoader.Schedule)) {
	f, err := os.Open("./assets/notification.mp3")
	if err != nil {
		ErrorChannel <- err
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		ErrorChannel <- err
		return
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ErrorChannel <- nil

	for {
		select {
		case <-CloseChannel:
			os.Exit(4)
			streamer.Close()
		default:
			Passed := Schedules.HasPassed()
			if len(Passed) >= 1 {
				err := Schedules.RemovePassed()
				if err != nil {
					DisplayError(err)
				}
				Table.Refresh()
				for _, v := range Passed {
					EventPassed(v)
				}
				streamer.Seek(0)
				requestFocus()
				speaker.Play(beep.Seq(streamer, beep.Callback(func() {})))
			}
		}
		time.Sleep(time.Second * 1)
	}

}

func main() {
	Application := app.New()
	Window := Application.NewWindow("SimpleSchedule")
	Window.Resize(fyne.NewSize(540, 360))

	Schedules, _ := ConfigLoader.ReadCFG("./schedules.json")

	// Header
	GUI_TITLE := widget.NewLabel("SimpleSchedule")
	GUI_VERSION := widget.NewLabel(SS_VERSION)
	GUI_HEADER := container.NewHBox(GUI_TITLE, GUI_VERSION)

	// Main Table
	GUI_TABLE := widget.NewTable(func() (int, int) {
		return Schedules.ScheduleAmounts() + 1, 4
	}, func() fyne.CanvasObject {
		return widget.NewLabel("............................")
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		label := cell.(*widget.Label)

		if id.Row == 0 {
			switch id.Col {
			case 0:
				label.SetText("ID")
			case 1:
				label.SetText("Title")
			case 2:
				label.SetText("Time")
			case 3:
				label.SetText("Has Passed")
			}
			return
		}

		schedule := Schedules.Access(id.Row - 1)
		switch id.Col {
		case 0:
			label.SetText(fmt.Sprintf("%v", id.Row-1))
		case 1:
			label.SetText(schedule.Title)
		case 2:
			label.SetText(fmt.Sprintf("%v", schedule.Time))
		case 3:
			label.SetText(fmt.Sprintf("%v", schedule.HasPassed()))
		}
	})

	// Buttons
	GUI_ADD_BTN := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		Title := widget.NewEntry()
		Title.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "Title can only contain letters, numbers, '_', and '-'")
		UnixTime := widget.NewEntry()
		UnixTime.Validator = validation.NewRegexp(`[0-9]+`, "Unix Timestamp can only be a Number")
		items := []*widget.FormItem{
			widget.NewFormItem("Title", Title),
			widget.NewFormItem("Unix Timestamp", UnixTime),
		}

		dialog.ShowForm("Add Event", "Add", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			i, err := strconv.Atoi(UnixTime.Text)
			if err != nil {
				dialog.ShowError(err, Window)
				return
			}
			Schedules.AddSchedule(Title.Text, int64(i))
			GUI_TABLE.Refresh()
		}, Window)
	})
	GUI_UPDATE_BTN := widget.NewButtonWithIcon("Update", theme.ViewRefreshIcon(), GUI_TABLE.Refresh)

	ErrorChannel := make(chan error)
	CloseChannel := make(chan bool)
	go Check(ErrorChannel, CloseChannel, &Schedules, GUI_TABLE, Window.RequestFocus, func(e error) {
		dialog.ShowError(e, Window)
	}, func(s *ConfigLoader.Schedule) {
		dialog.ShowInformation(s.Title, fmt.Sprintf("Look like you have an event named '%v'.", s.Title), Window)
	})
	err := <-ErrorChannel
	if err != nil {
		dialog.ShowError(err, Window)
	}

	Application.Lifecycle().SetOnStopped(func() {
		CloseChannel <- true
	})

	Window.SetContent(container.NewBorder(GUI_HEADER, container.NewAdaptiveGrid(1, container.NewHBox(GUI_ADD_BTN, GUI_UPDATE_BTN)), nil, nil, GUI_TABLE))
	Window.ShowAndRun()
}
