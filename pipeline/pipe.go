package pipeline

// Pipe passes values from an input channel to a Stage for handling, then
// sends the output of the Stage onto the output channel. The Stage func
// will be run in its own go routine.
func Pipe(stage Stage) Pipeline {
	return ParallelPipe(stage, 1)
}
