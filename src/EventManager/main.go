package EventManager

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Event struct {
	Title         string `json:"title"`
	UnixTimestamp int64  `json:"timestamp"`
}

type EventList struct {
	Path   string
	Events []*Event
}

func (Schedule *EventList) PassedEvents() (res []*Event) {
	res = make([]*Event, 0)
	for _, v := range Schedule.Events {
		if time.Now().UnixMilli() >= v.UnixTimestamp {
			res = append(res, v)
		}
	}
	return
}

func (Schedule *EventList) RemovePassed() error {
	var res []*Event = make([]*Event, 0)
	for _, v := range Schedule.Events {
		if time.Now().UnixMilli() <= v.UnixTimestamp {
			res = append(res, v)
		}
	}
	Schedule.Events = res
	return Schedule.Save()
}

func (Schedule *EventList) AddEvent(title string, unix_timestamp int64) {
	Schedule.Events = append(Schedule.Events, &Event{Title: title, UnixTimestamp: unix_timestamp})
	Schedule.Save()
}

func (Schedule *EventList) GetEvent(index int) (*Event, bool) {
	if len(Schedule.Events)-1 >= index {
		return Schedule.Events[index], true
	} else {
		return &Event{}, false
	}
}

func (Schedule *EventList) Save() error {
	data, err := json.MarshalIndent(Schedule.Events, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(Schedule.Path, data, 0777)

	if err != nil {
		return err
	}

	return nil
}

func ReadEvents(JSON_path string) (*EventList, error) {
	dat, err := ioutil.ReadFile(JSON_path)
	if err != nil {
		return &EventList{}, err
	}

	var res *EventList = &EventList{Events: make([]*Event, 0), Path: JSON_path}

	err = json.Unmarshal(dat, &res.Events)
	if err != nil {
		return res, err
	}

	return res, nil
}
