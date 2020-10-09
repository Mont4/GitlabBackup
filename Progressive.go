package main

import (
	"fmt"
	"strings"
	"sync"
)

const BlankCount = 30

type Progressive struct {
	total      int
	current    int
	percentage float64
	done       bool

	mutex sync.Mutex
}

func NewProgressive(total int) *Progressive {
	p := &Progressive{}

	p.total = total

	return p
}

func (p *Progressive) SetTotal(total int) {
	p.total = total

	p.setPercentage()
}

func (p *Progressive) AddTotal() {
	p.AddTotalBy(1)
}

func (p *Progressive) AddTotalBy(count int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.total += count

	p.setPercentage()
	p.printProgressive()
}

func (p *Progressive) Start() {
	p.setPercentage()
	p.printProgressive()
}

func (p *Progressive) Advance() {
	p.AdvanceBy(1)
}

func (p *Progressive) AdvanceBy(count int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.current += count
	if p.current > p.total {
		p.total = p.current
	}

	p.setPercentage()
	p.printProgressive()
}

func (p *Progressive) Done() {
	p.done = true

	p.printProgressive()
}

func (p *Progressive) setPercentage() {
	p.percentage = float64(p.current) / float64(p.total) * 100
	if p.done {
		p.percentage = 100.
	}
}

func (p *Progressive) printProgressive() {

	// fmt.Println(strings.Repeat("=", ratio)+">")
	progressBar := strings.Repeat("=", BlankCount)
	if !p.done {
		progressBar = strings.Repeat("=", int(p.percentage/3.33334))
		progressBar += ">"
		progressBar += strings.Repeat(" ", int((100-p.percentage)/3.33334))
	}

	bar := fmt.Sprintf("\r[%s] %3.1f %d / %d", progressBar, p.percentage, p.current, p.total)

	fmt.Printf(bar)
}
