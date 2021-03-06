// Toy scheduler is a mock scheduling output that simply reflects back whatever we tell it
package toy

import (
	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

var log = logging.MustGetLogger("mssscheduler")

type ToyScheduler struct {
}

func NewScheduler() *ToyScheduler {
	toy := ToyScheduler{}
	return &toy
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*ToyScheduler)(nil)

func (t *ToyScheduler) InitScheduler(task *demand.Task) error {
	log.Infof("Toy scheduler initialized task %s with %d initial demand", task.Name, task.Demand)
	return nil
}

// StopStartNTasks asks the scheduler to bring the number of running tasks up to task.Demand.
func (t *ToyScheduler) StopStartTasks(tasks *demand.Tasks) error {
	tasks.Lock()
	defer tasks.Unlock()

	for _, task := range tasks.Tasks {
		task.Requested = task.Demand
		log.Debugf("Toy scheduler setting Requested for %s to %d", task.Name, task.Requested)
	}

	return nil
}

// CountAllTasks for the Toy scheduler simply reflects back what has been requested
func (t *ToyScheduler) CountAllTasks(running *demand.Tasks) error {
	running.Lock()
	defer running.Unlock()

	for _, task := range running.Tasks {
		task.Running = task.Requested
	}
	return nil
}

func (t *ToyScheduler) Cleanup() error { return nil }
