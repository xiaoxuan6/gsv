package spinner

import (
	"github.com/briandowns/spinner"
	"time"
)

type any = interface{}

func RunF0(prefix string, f func()) {
	s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	s.Prefix = prefix
	s.FinalMSG = "done"
	s.Start()

	f()

	s.Stop()
	return
}

func RunF[T any](prefix string, f func() T) T {
	s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	s.Prefix = prefix
	s.FinalMSG = "done"
	s.Start()

	result := f()

	s.Stop()

	return result
}

func RunF2[T any, R any](prefix string, f func() (T, R)) (T, R) {
	s := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	s.Prefix = prefix
	s.FinalMSG = "done"
	s.Start()

	result1, result2 := f()

	s.Stop()

	return result1, result2
}
