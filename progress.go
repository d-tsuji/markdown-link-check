package markdownlink

import (
	"github.com/cheggaaa/pb/v3"
)

var DefaultProgressBar = &progressBar{}

type progressBar struct {
	progress *pb.ProgressBar
}

func (p *progressBar) Start(count int) {
	p.progress = pb.Simple.New(count)
	p.progress.Set("prefix", "Checking... ")
	p.progress.SetWidth(50)
	p.progress.Start()
}

func (p *progressBar) Increment() {
	p.progress.Increment()
}

func (p *progressBar) Finish() {
	p.progress.Finish()
}
