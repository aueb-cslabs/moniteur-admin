// +build linux darwin

package main

import (
	"encoding/json"
	"fmt"
	"github.com/aueb-cslabs/moniteur/types"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Lesson struct {
	Semester        string `json:"semester"`
	LessonComments  string `json:"Lesson_comments"`
	Room            string `json:"Room"`
	LessonTitle     string `json:"Lesson_title"`
	Professor       string `json:"professor"`
	Time            string `json:"time"`
	Day             string `json:"Day"`
	DepartmentTitle string `json:"Department_title"`
}

type RoomMap struct {
	Rooms map[string]string `yaml:"rooms,omitempty"`
}

// PLUGIN The plugin to be read by the moniteur agent.
var PLUGIN = Plugin{}
var mapping = &RoomMap{}

type Plugin struct {
}

// Initialize Method that initializes crucial functions for the plugin
func (Plugin) Initialize() {
	if len(mapping.Rooms) == 0 {
		mapping, _ = loadMapping("mapping.yml")
	}
}

// Schedule Method that returns current schedule from Schedule Master
func (Plugin) Schedule() (*types.Schedule, error) {
	return retriever(), nil
}

// ScheduleRoom Method that returns current schedule and the room that corresponds to it
func (Plugin) ScheduleRoom(room string) (*types.Schedule, error, string) {
	room, changed := checkMapping(room)
	if !changed {
		room = convertChars(room)
	}
	return retriever(), nil, room
}

func (Plugin) ExamsSchedule() (*types.Schedule, error) {
	return getExamsSchedule(), nil
}

// retriever Method that converts Schedule Master json to our json format
func retriever() *types.Schedule {
	resp := &types.Schedule{}
	for _, lesson := range getEntireSchedule() {
		subject := &types.ScheduleSlot{}
		lessonTime := strings.Split(lesson.Time, "-")
		start, _ := strconv.ParseInt(lessonTime[0], 10, 64)
		end, _ := strconv.ParseInt(lessonTime[1], 10, 64)
		subject.Start = start * int64(3600)
		subject.End = end * int64(3600)
		subject.Room = lesson.Room
		subject.Day = determineDay(lesson.Day)
		subject.Title = lesson.LessonTitle
		subject.Host = lesson.Professor
		subject.Semester = lesson.Semester
		resp.Slots = append(resp.Slots, subject)
	}
	return resp
}

// convertChars Method that converts english characters to greek in order to parse Schedule Master Room Name
func convertChars(room string) string {
	re := regexp.MustCompile("[0-9]+")

	if re.MatchString(room) {
		if strings.Contains(room, "a") {
			room = strings.ReplaceAll(room, "a", "Α")
		}
		if strings.Contains(room, "d") {
			room = strings.ReplaceAll(room, "d", "Δ")
		}
		if strings.Contains(room, "t") {
			room = strings.ReplaceAll(room, "t", "Τ")
		}
		return room
	}

	return room
}

// getEntireSchedule Method that retrieves the schedule for Schedule Master API
func getEntireSchedule() []*Lesson {
	resp, _ := http.Get("http://schedule.aueb.gr/mobile/")
	bts, _ := ioutil.ReadAll(resp.Body)
	var slots []*Lesson
	_ = json.Unmarshal(bts, &slots)
	return slots
}

// determineDay Method that converts the Day (from the greek language) to an int
func determineDay(day string) int {
	switch day {
	case "Δευτέρα", "ΔΕΥΤΕΡΑ":
		return 1
	case "Τρίτη", "ΤΡΙΤΗ":
		return 2
	case "Τετάρτη", "ΤΕΤΑΡΤΗ":
		return 3
	case "Πέμπτη", "ΠΕΜΠΤΗ":
		return 4
	case "Παρασκευή", "ΠΑΡΑΣΚΕΥΗ":
		return 5
	case "Σάββατο", "ΣΑΒΒΑΤΟ":
		return 6
	}
	return 0
}

// loadMapping Method that loads the room mapping from english to greek names
func loadMapping(file string) (*RoomMap, error) {
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	rooms := &RoomMap{}
	if err := yaml.Unmarshal(byt, rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

// checkMapping Method that checks the room map in order to retrieve the right room name
func checkMapping(room string) (string, bool) {
	rooms := mapping.Rooms

	if rooms != nil && rooms[room] != "" {
		return rooms[room], true
	}

	return room, false
}

func determineMonth(month string) int {
	switch month {
	case "ΙΑΝΟΥΑΡΙΟΥ":
		return 1
	case "ΦΕΒΡΟΥΑΡΙΟΥ":
		return 2
	case "ΜΑΡΤΙΟΥ":
		return 3
	case "ΑΠΡΙΛΙΟΥ":
		return 4
	case "ΜΑΪΟΥ":
		return 5
	case "ΙΟΥΝΙΟΥ":
		return 6
	case "ΙΟΥΛΙΟΥ":
		return 7
	case "ΑΥΓΟΥΣΤΟΥ":
		return 8
	case "ΣΕΠΤΕΜΒΡΙΟΥ":
		return 9
	case "ΟΚΤΩΒΡΙΟΥ":
		return 10
	case "ΝΟΕΜΒΡΙΟΥ":
		return 11
	case "ΔΕΚΕΜΒΡΙΟΥ":
		return 12
	}
	return 0
}

func getExamsSchedule() *types.Schedule {
	t := time.Now()
	date := t.Format("20060102")
	month := t.Month()
	year := t.Year()
	link := fmt.Sprintf("https://aueb.gr/sites/default/files/aueb/%s_Exams_%s.xlsx", month, date)
	resp, _ := http.Get(link)
	schedule := &types.Schedule{}

	if resp.StatusCode == 200 {
		bts, _ := ioutil.ReadAll(resp.Body)
		file, _ := xlsx.OpenBinary(bts)
		rows := file.Sheet[fmt.Sprintf("%d", year-1)].Rows
		var dayName string
		var day int
		var month int

		for i := 0; i < len(rows); i++ {
			var rooms []string
			var semester string
			var lessonName string
			var examiner string

			if strings.Contains(rows[i].Cells[0].Value, "*") {
				continue
			}

			if strings.Contains(rows[i].Cells[0].Value, "ΗΜΕΡΟΜΗΝΙΑ") {
				continue
			}

			if rows[i].Cells[0].Value == "" {
				continue
			}

			if strings.Contains(rows[i].Cells[4].Value, "ΠΡΥΤΑΝΕΙΑ") {
				break
			}

			if !strings.Contains(rows[i].Cells[1].Value, ":") {
				examsDate := strings.Split(rows[i].Cells[0].Value, " ")
				dayName = examsDate[0]
				day, _ = strconv.Atoi(examsDate[1])
				month = determineMonth(examsDate[2])
				fmt.Println(dayName, day, month)

			} else {

				rooms = strings.Split(rows[i].Cells[0].Value, ", ")
				timestamp := strings.Split(rows[i].Cells[1].Value, "-")
				semester = rows[i].Cells[2].Value
				lessonName = rows[i].Cells[3].Value
				examiner = rows[i].Cells[4].Value
				start := convertTime(timestamp[0])
				end := convertTime(timestamp[1])
				for j := range rooms {
					slot := &types.ScheduleSlot{}
					slot.Room = rooms[j]
					slot.Day = determineDay(dayName)
					slot.Start = start
					slot.End = end
					slot.Title = lessonName
					slot.Host = examiner
					slot.Semester = semester
					slot.DayNum = day
					slot.MonthNum = int(month)
					schedule.Slots = append(schedule.Slots, slot)
					fmt.Println(slot)
				}
			}
		}
	} else {
		return nil
	}
	return schedule
}

func convertTime(timestamp string) int64 {
	convert := strings.Split(timestamp, ":")
	var hourF float64

	hourInt, _ := strconv.Atoi(convert[0])

	hourF = float64(hourInt)
	if convert[1] == "30" {
		hourF += .5
	}

	return int64(hourF * 3600)
}
