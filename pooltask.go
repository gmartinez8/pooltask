package pooltask

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//maxWorkers is the number of go routines that can be executed concurrently
var maxWorkers = 10

//activeWorkers is the number of go routines that are executing now
var activeWorkers int = 0
var currentJobs = make(chan *Task)
var processedJobs = make(chan *Task)

//currentWorkers is a map with all execution time of current tasks key is the ID of the task
var currentWorkers map[string]int = make(map[string]int)

//tasks of the system, will work better with a DB connection
//but for the task purpose not a requirement asked
var tasks = make(map[string]*Task)

//mutex to protect activeWorkers and tasks from data race
var mutex = &sync.Mutex{}

//CallbackURL is exported if you whant to change the config of the Callback after processing a task
//Upon execution of the task, service would issue a callback to CallbackURL
var CallbackURL = "http://localhost:8080/callback"

//ErrorResponse allows us to send Error message with error status
type ErrorResponse struct {
	Message string
}

//CreateResponse for expected response after creating a task
type CreateResponse struct {
	ID string `json:"taskID"`
}

//CallbackRequest for expected response after creating a task
type CallbackRequest struct {
	ID      string `json:"taskID"`
	Success bool   `json:"success"`
}

//MaxWorkers allows us to get maxWorkers value
func MaxWorkers() int {
	return maxWorkers
}

//SetMaxWorkers allows us to set maxWorkers value avoiding data race
func SetMaxWorkers(mw int) {
	maxWorkers = mw
}

//HandleHome Handles Home Route and /
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" HandleHome Called \n"))
}

//HandleListTasks Handles Route
func HandleListTasks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	response, err := json.Marshal(tasks)
	if len(tasks) == 0 {
		response, _ = json.Marshal(make([]Task, 0))
	}
	mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		e, _ := json.Marshal(ErrorResponse{err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(e)
		return
	}
	w.Write(response)
}

//HandleCreateTask Handles Route
func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Task
	err := decoder.Decode(&t)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		e, _ := json.Marshal(ErrorResponse{err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(e)
		return
	}
	t.SetID()
	t.SetCreatedAt()
	cr := &CreateResponse{
		ID: t.ID,
	}
	response, err := json.Marshal(cr)
	if err != nil {
		e, _ := json.Marshal(ErrorResponse{err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(e)
		return
	}
	//if activeWorkers  == maxWorkers return HTTP code 503 with "Retry-After" header with calculated time when at least one worker will become ready
	mutex.Lock()
	if activeWorkers == maxWorkers {
		freeIn := strconv.Itoa(minIntMap(currentWorkers))
		e, _ := json.Marshal(ErrorResponse{"Retry-After " + freeIn + " secs"})
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(e)
		mutex.Unlock()
		return
	}
	activeWorkers++
	tasks[t.ID] = &t
	go addTask(currentJobs, tasks[t.ID])
	go processTask(currentJobs, processedJobs)
	go workFinished(processedJobs)
	mutex.Unlock()
	w.Write(response)
}

//work process the task
func addTask(jobs chan *Task, t *Task) {
	//log.Println("activeWorkers: ", activeWorkers)
	jobs <- t
	mutex.Lock()
	currentWorkers[t.ID] = t.ExecutionTime
	mutex.Unlock()
}

//processTask process the task
//jobs <-chan only reads on jobs chanel
//results chan<- only sends on results chanel
func processTask(jobs <-chan *Task, results chan<- *Task) {
	t := <-jobs
	//log.Println("Start executing: ", t)
	//log.Println("For this much seconds: ", t.ExecutionTime)
	t.SetExecutedAt()
	time.Sleep(time.Duration(t.ExecutionTime) * time.Second)
	t.Status = 1
	t.SetFinishedAt()
	results <- t
}

//Task executed, will issue a Callback to CallbackURL with this fields:
//BODY{ "taskID: string, "success": bool }
func workFinished(results chan *Task) {
	t := <-results
	mutex.Lock()
	t.Status = 2
	activeWorkers--
	delete(currentWorkers, t.ID)
	mutex.Unlock()
	//log.Println("Finished executing: ", t.ID, t.CreatedAt, t.ExecutedAt, t.FinishedAt)
	cr := &CallbackRequest{
		t.ID,
		true,
	}
	body, err := json.Marshal(cr)
	if err != nil {
		log.Fatalln(err)
	}
	r, err := http.NewRequest("POST", CallbackURL, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		log.Printf("Callback for TaskID:%s on %s = %s", t.ID, CallbackURL, resp.Status)
	}
	//log.Println("activeWorkers on finished: ", activeWorkers)
}

//calc min int value of a map[string]int
func minIntMap(w map[string]int) int {
	var min int
	for _, min = range w {
		break
	}
	for _, e := range w {
		if e < min {
			min = e
		}
	}
	return min
}
