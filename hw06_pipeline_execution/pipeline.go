package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := take(in, done)
	for _, stage := range stages {
		out = stage(take(out, done))
	}
	return out
}

func take(in In, done In) Out {
	outStream := make(Bi)
	go func() {
		defer func() {
			close(outStream)
			for range in { //nolint:revive
			}
		}()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case outStream <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return outStream
}
