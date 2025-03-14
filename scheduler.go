//go:build scheduler

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const AGING_FACTOR = 5
const TOP_N = 3

type Task struct {
	priority  int
	name      string
	init_time int
}

func (t *Task) EffectivePriority() int {
	waiting_time := int(time.Now().Unix()) - t.init_time
	effective_priority := t.priority + (waiting_time * AGING_FACTOR)
	return effective_priority
}

type TaskScheduler struct {
	tasks []Task
	m     sync.Mutex
}

func (ts *TaskScheduler) Add(task Task) {
	ts.m.Lock()
	ts.tasks = append(ts.tasks, task)
	ts.m.Unlock()
}

func (ts *TaskScheduler) Generate(task_stream chan Task, top_n int) {
	ts.m.Lock()

	// sort by EffectivePriority()
	sort.Slice(ts.tasks, func(i, j int) bool {
		return ts.tasks[i].EffectivePriority() > ts.tasks[j].EffectivePriority()
	})

	fmt.Println("init tasks count: ", len(ts.tasks))

	if len(ts.tasks) < top_n {
		// limit to all if top_n is more then number of items
		top_n = len(ts.tasks)
	}

	var i int
	for i = 0; i < top_n; i++ {
		task_stream <- ts.tasks[i]
		fmt.Println("dispatched ", ts.tasks[i].name, " with priority ", ts.tasks[i].EffectivePriority())
	}
	// remove what's dispatched
	copy(ts.tasks, ts.tasks[top_n:])
	ts.tasks = ts.tasks[:len(ts.tasks)-top_n]

	ts.m.Unlock()
	fmt.Println("left tasks count: ", len(ts.tasks))
}

func executor(executor_name string, task_stream chan Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range task_stream {
		fmt.Println(executor_name, " executing task ", task.name)
		time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup

	scheduler := TaskScheduler{
		tasks: make([]Task, 0),
		m:     sync.Mutex{},
	}
	task_stream := make(chan Task, 10)

	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go executor(fmt.Sprintf("#%d", i), task_stream, &wg)
	}

	for i := 0; i < 30; i++ {
		priority := rand.Intn(100)
		scheduler.Add(Task{
			priority:  priority,
			name:      fmt.Sprintf("task#{%d}", priority),
			init_time: int(time.Now().Unix()) - rand.Intn(1000),
		})
		fmt.Println("task added ", priority)
	}

	for len(scheduler.tasks) > 0 {
		scheduler.Generate(task_stream, TOP_N)
		time.Sleep(100 * time.Millisecond)
	}

	close(task_stream)

	wg.Wait()
}
