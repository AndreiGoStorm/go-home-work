package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(take(out, done))
	}
	return take(out, done)
}

func take(in In, done In) Out {
	outStream := make(Bi)
	go func() {
		defer close(outStream)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				outStream <- v
			case <-done:
				return
			}
		}
	}()
	return outStream
}
