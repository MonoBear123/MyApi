package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Operation представляет операцию для вычислителя
type Operation struct {
	Name      string    `json:"name"`
	Execution time.Time `json:"execution"`
}

// Task представляет задачу для вычислителя
type Task struct {
	ExpressionID string    `json:"expressionID"`
	Operation    Operation `json:"operation"`
}

// Result представляет результат выполнения задачи
type Result struct {
	ExpressionID string  `json:"expressionID"`
	Result       float64 `json:"result"`
}

// Agent представляет вычислитель (агент)
type Agent struct {
	completedTasks map[string]Result
	mu             sync.Mutex
}

func NewAgent() *Agent {
	return &Agent{
		completedTasks: make(map[string]Result),
	}
}

func (a *Agent) handleTask(w http.ResponseWriter, r *http.Request) {
	// Получение задачи от оркестратора
	resp, err := http.Get("http://localhost:8080/get-task-for-agent")
	if err != nil {
		fmt.Println("Ошибка при получении задачи:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка при получении задачи, статус:", resp.Status)
		w.WriteHeader(resp.StatusCode)
		return
	}

	var task Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		fmt.Println("Ошибка при разборе задачи:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Выполнение операции
	// "Вычисление" может занять время, симулируем это задержкой

	// Сохранение результата в мапе
	result := Result{
		ExpressionID: task.ExpressionID,
		Result:       42.0, // Здесь нужно реальное значение результата
	}

	a.mu.Lock()
	a.completedTasks[task.ExpressionID] = result
	a.mu.Unlock()

	// Отправка результата обратно оркестратору
	completedTask := Task{
		ExpressionID: task.ExpressionID,
		Operation:    task.Operation,
	}

	resp, err = http.Post("http://localhost:8080/receive-completed-task", "application/json", encodeJSON(completedTask))
	if err != nil {
		fmt.Println("Ошибка при отправке результата:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка при отправке результата, статус:", resp.Status)
		w.WriteHeader(resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func encodeJSON(data interface{}) *bytes.Buffer {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	return buf
}

func main() {
	agent := NewAgent()

	http.HandleFunc("/handle-task", agent.handleTask)

	fmt.Println("Агент слушает на порту 8090...")
	http.ListenAndServe(":8090", nil)
}
