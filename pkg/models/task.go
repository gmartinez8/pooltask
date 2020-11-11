package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

//timeFormat = "2006.01.02 15:04:05.000000"
const (
	timeFormat = "2006.01.02 15:04:05.000000"
)

//Task struct
//can be formated to JSON
//also omitempty allows us to handle zero values on
//Status 0: not executed, 1: executing, 2: finished sucessfully
type Task struct {
	ID            string `json:"taskID"`
	Status        int    `json:"status,omitempty"`
	ExecutionTime int    `json:"processMeForThisMuchSeconds"`
	Detail        string `json:"detail,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	ExecutedAt    string `json:"executed_at,omitempty"`
	FinishedAt    string `json:"finished_at,omitempty"`
}

//NewTask creates a new user
//ExecutionTime in seconds
func NewTask(ExecutionTime int) *Task {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	t := time.Now()
	return &Task{
		ID:            hex.EncodeToString(buf),
		Status:        0,
		ExecutionTime: ExecutionTime,
		CreatedAt:     t.Format(timeFormat),
	}
}

//SetID Creates a random ID for the user using crypto/rand library
//uuid will work better
func (t *Task) SetID() {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	t.ID = hex.EncodeToString(buf)
}

//SetCreatedAt set CreatedAt datetime in format 2006.01.02 15:04:05.000000
func (t *Task) SetCreatedAt() {
	tm := time.Now()
	t.CreatedAt = tm.Format(timeFormat)
}

//SetExecutedAt set ExecutedAt datetime in format 2006.01.02 15:04:05.000000
func (t *Task) SetExecutedAt() {
	tm := time.Now()
	t.ExecutedAt = tm.Format(timeFormat)
}

//SetFinishedAt set FinishedAt datetime in format 2006.01.02 15:04:05.000000
func (t *Task) SetFinishedAt() {
	tm := time.Now()
	t.FinishedAt = tm.Format(timeFormat)
}
