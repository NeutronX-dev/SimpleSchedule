package ConfigLoader

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Schedule struct {
	Title string `json:"title"`
	Time  int64  `json:"time"`
}

func (schedule *Schedule) HasPassed() bool {
	return time.Now().UnixMilli() >= schedule.Time
}

type Schedules struct {
	data []*Schedule
	path string
}

func (schedules *Schedules) ScheduleAmounts() int {
	return len(schedules.data)
}

func (schedules *Schedules) Access(index int) *Schedule {
	return schedules.data[index]
}

func (schedules *Schedules) HasPassed() []*Schedule {
	var res []*Schedule = make([]*Schedule, 0)
	for _, v := range schedules.data {
		if v.HasPassed() {
			res = append(res, v)
		}
	}
	return res
}

func (schedules *Schedules) RemovePassed() error {
	var res []*Schedule = make([]*Schedule, 0)
	for _, v := range schedules.data {
		if time.Now().UnixMilli() <= v.Time {
			res = append(res, v)
		}
	}
	schedules.data = res
	return schedules.Save()
}

func (schedules *Schedules) AddSchedule(Title string, UnixTime int64) error {
	var NewSchedule *Schedule = &Schedule{
		Title: Title,
		Time:  UnixTime,
	}
	schedules.data = append(schedules.data, NewSchedule)

	return schedules.Save()
}

func (schedules *Schedules) Save() error {
	data, err := json.MarshalIndent(schedules.data, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(schedules.path, data, 0777)

	if err != nil {
		return err
	}

	return nil
}

func ReadCFG(path string) (Schedules, error) {
	var res Schedules
	res.path = path
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res.data)
	if err != nil {
		return res, err
	}

	return res, nil
}
