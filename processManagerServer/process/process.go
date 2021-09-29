package process

import "time"

/* Process struct and methods  */
type Process struct {
	ID        uint64
	ICount    uint64
	Executing bool
}

func (p *Process) Execute() {
	/* Process execution */
	for {
		if p.Executing {
			p.ICount++
			time.Sleep(time.Millisecond * 500)
		} else {
			break
		}
	}
}

func (p *Process) Stop() {
	/* stops a process*/
	p.Executing = false
}
