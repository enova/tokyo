package spawn

import (
	"fmt"
	"github.com/enova/tokyo/src/alert"
	"os"
	"os/exec"
	"sync"
)

// Spawner spawns a set of commands only N at a time. If N is zero, it
// spawns all the commands at once (i.e. 0 => Infinity).
type Spawner struct {
	max   int
	procs chan int
	alive sync.WaitGroup
	cmds  []string
	errs  []string
}

// NewSpawner returns a new Spawner instance and sets the max
func NewSpawner(max int) *Spawner {
	return &Spawner{max: max}
}

// Add adds new commands to the Spawner
func (s *Spawner) Add(cmds ...string) {

	for _, cmd := range cmds {
		s.cmds = append(s.cmds, cmd)
		s.errs = append(s.errs, "")
	}
}

// Run begins executing commands and returns after all have completed.
func (s *Spawner) Run() {
	if s.max == 0 {
		s.max = len(s.cmds)
	}

	s.procs = make(chan int, s.max)
	s.alive.Add(len(s.cmds))

	// Spawn Commands
	for i := 0; i < len(s.cmds); i++ {
		go s.spawn(i)
	}

	s.alive.Wait()
}

// Err returns an error containing all errors that may have occurred
// during the Run() call. If there were no errors, it returns nil.
func (s *Spawner) Err() error {
	var msg string

	for i, e := range s.errs {
		if len(e) > 0 {
			msg += fmt.Sprintf("Error Cmd #%d: %s\n", i, e)
		}
	}

	// No Errors!
	if len(msg) == 0 {
		return nil
	}

	return fmt.Errorf("%s", msg)
}

// Spawn is a convenience function for spawning a single command
func Spawn(cmd string) error {
	s := NewSpawner(0)
	s.Add(cmd)
	s.Run()
	return s.Err()
}

// Spawn a process
func (s *Spawner) spawn(i int) {
	s.procs <- i

	cmd := s.cmds[i]
	fmt.Fprintf(os.Stderr, "Spawning: %s\n", cmd)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		msg := fmt.Sprintf("Spawner could not run command: %s, %s\n", cmd, err.Error())
		s.errs[i] = msg
		alert.Cerr(msg)
	}

	_ = <-s.procs
	s.alive.Done()
}
