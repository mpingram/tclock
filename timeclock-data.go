package main

import (
	"fmt"
	"time"
)

type projectStruct struct {
	name      string
	namespace string
}

type timeslotDataMessage struct {
	project      projectStruct
	clockOnTime  time.Time
	clockOffTime time.Time
}

type editTimeslotMessage struct {
	targetProject  projectStruct
	newProjectName projectStruct
}

type editProjectMessage struct {
	targetProject   projectStruct
	newClockOnTime  time.Time
	newClockOffTime time.Time
}

func sendData(data timeslotDataMessage) {
	fmt.Print(data.project.name)
}
