package main

import (
	"fmt"
	"os"
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

	"main/src/AudioPlayer"
	"main/src/ConfigLoader"
	"main/src/EventManager"
	Inst "main/src/Instructive"
	SS_Main "main/src/SS_Custom"
)

const (
	SS_VERSION = "1.0.3"

	GUI_HEIGHT = 360
	GUI_WIDTH  = 540
)

func FirstAndLastDayOfMonth(now time.Time) (int, int) {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth.Day(), lastOfMonth.Day()
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
						SS_Main.Custom_Main(v.Title, v.UnixTimestamp)
						/*^^^^^^^^^^^^^^^^^^^^^^^^^^
						  |     Custom Callback    |
						  +------------------------+
						  |        Found on        |
						  | /src/SS_Custom/main.go |
						  +------------------------+ */
						if cfg.Instructive && strings.HasPrefix(v.Title, "i:") {
							ParsedInstruction, err := Inst.ParseInstructions(v.Title)
							if err != nil {
								DisplayError(err)
							} else {
								for _, Instruction := range ParsedInstruction {
									requestFocus()
									err = Instruction.Execute()
									if err != nil {
										DisplayError(err)
									}
								}
							}
						} else {
							EventPassed(v)
							requestFocus()
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
				label.SetText("[Instructive(s)]")
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

		var NOW time.Time = time.Now()
		var DATE_Day int = NOW.Day()
		var DATE_Month string
		var DATE_Year int = NOW.Year()
		var DATE_Hour int = NOW.Hour()
		var DATE_Minute int = NOW.Minute()

		ALL_MONTHS_ARR := []string{"01: January", "02: February", "03: March", "04: April", "05: May", "06: June", "07: July", "08: August", "09: September", "10: October", "11: November", "12: December"}
		FORM_DATE_Month := widget.NewSelect(ALL_MONTHS_ARR, func(s string) {
			DATE_Month = strings.Split(s, ": ")[0]
		})

		FORM_DATE_Month.SetSelected(ALL_MONTHS_ARR[NOW.Month()-1])

		FORM_DATE_Day := widget.NewEntry()
		FORM_DATE_Day.PlaceHolder = "Day"
		FORM_DATE_Day.Validator = validation.NewRegexp(`^[0-9]{1,2}$`, "Day can ONLY have TWO DIGITS.")
		FORM_DATE_Day.SetText(fmt.Sprintf("%v", DATE_Day))

		FORM_DATE_Year := widget.NewEntry()
		FORM_DATE_Year.PlaceHolder = "Year"
		FORM_DATE_Year.Validator = validation.NewRegexp(`^[0-9]{4}$`, "Year can ONLY have FOUR DIGITS.")
		FORM_DATE_Year.SetText(fmt.Sprintf("%v", DATE_Year))

		FORM_DATE_Hour := widget.NewEntry()
		FORM_DATE_Hour.PlaceHolder = "Hour"
		FORM_DATE_Hour.Validator = validation.NewRegexp(`^[0-9]{1,2}$`, "Hour can ONLY have TWO DIGITS.")
		FORM_DATE_Hour.SetText(fmt.Sprintf("%v", DATE_Hour))

		FORM_DATE_Minute := widget.NewEntry()
		FORM_DATE_Minute.PlaceHolder = "Minute"
		FORM_DATE_Minute.Validator = validation.NewRegexp(`^[0-9]{1,2}$`, "Minute can ONLY have TWO DIGITS.")
		FORM_DATE_Minute.SetText(fmt.Sprintf("%v", DATE_Minute))

		// DATE_Year := time.Now().Year()

		items := []*widget.FormItem{
			widget.NewFormItem("Title", title),
			widget.NewFormItem("", widget.NewLabel("")),
			widget.NewFormItem("Event Date", container.NewAdaptiveGrid(3, FORM_DATE_Month, FORM_DATE_Day, FORM_DATE_Year)),
			widget.NewFormItem("[24h] Event Time", container.NewHBox(FORM_DATE_Hour, widget.NewLabel(":"), FORM_DATE_Minute)),
		}

		form := dialog.NewForm("Add Event", "Add Event", "Cancel", items, func(b bool) {
			if b {

				DayInt, err := strconv.Atoi(FORM_DATE_Day.Text)
				if err != nil {
					dialog.ShowError(err, Window)
					return
				}

				Day := FORM_DATE_Day.Text
				if DayInt < 10 {
					Day = fmt.Sprintf("0%v", DayInt)
				}

				HourInt, err := strconv.Atoi(FORM_DATE_Hour.Text)
				if err != nil {
					dialog.ShowError(err, Window)
					return
				}

				Hour := FORM_DATE_Hour.Text
				if HourInt < 10 {
					Hour = fmt.Sprintf("0%v", HourInt)
				}

				MinInt, err := strconv.Atoi(FORM_DATE_Minute.Text)
				if err != nil {
					dialog.ShowError(err, Window)
					return
				}

				Min := FORM_DATE_Minute.Text
				if MinInt < 10 {
					Min = fmt.Sprintf("0%v", HourInt)
				}

				layout := "01/02/2006 15:04:05 MST-0700"
				date_string := fmt.Sprintf("%v/%v/%v %v:%v:00 GMT-0500", DATE_Month, Day, FORM_DATE_Year.Text, Hour, Min)

				EVNT_TIME, err := time.Parse(layout, date_string)

				if err != nil {
					dialog.ShowError(err, Window)
					return
				}

				fmt.Println(date_string)
				fmt.Println(EVNT_TIME)

				Events.AddEvent(title.Text, EVNT_TIME.UnixMilli())
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
				{Text: "", Widget: widget.NewButton("Reload Events", func() {
					Events, _ = EventManager.ReadEvents("./schedules.json")
					GUI_TABLE.Refresh()
				})},
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
