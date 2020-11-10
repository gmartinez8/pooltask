package pooltask

import (
	"encoding/json"
	"net/http"
	"sync"

	"queueworker/pkg/models"
)

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

//Tasks of the system, will work better with a DB connection
//but for the task purpose not a requirement asked
var tasks = make(map[string]*models.Task)
var mutex = &sync.Mutex{}

//HandleHome Handles Home Route and /
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" HandleHome Called \n"))
}

//HandleListTasks Handles Route
func HandleListTasks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	response, err := json.Marshal(tasks)
	if len(tasks) == 0 {
		response, _ = json.Marshal(make([]models.Task, 0))
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
	var t models.Task
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
	//Add logic for go routine, check If all workers are busy processing and will be processing them for more than 1 second
	//if true return HTTP code 503 with "Retry-After" header with calculated time when at least one worker will become ready
	//else add task to process queue
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
	mutex.Lock()
	tasks[t.ID] = &t
	mutex.Unlock()
	w.Write(response)
}

//HandleCallback handles all finished tasks
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cr CallbackRequest
	err := decoder.Decode(&cr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		e, _ := json.Marshal(ErrorResponse{err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(e)
		return
	}
	response, err := json.Marshal(cr)
	if err != nil {
		e, _ := json.Marshal(ErrorResponse{err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(e)
		return
	}
	w.Write(response)
}
