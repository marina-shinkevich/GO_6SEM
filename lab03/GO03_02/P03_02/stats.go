package P03_02

import "fmt"

type RequestStats struct {
	GetCount  int
	PostCount int
}

func (s *RequestStats) PlusGet() {
	s.GetCount++
}

func (s *RequestStats) PlusPost() {
	s.PostCount++
}

func (s *RequestStats) GenStr() string {
	return fmt.Sprintf("Get-request count = %d, Post-request count = %d", s.GetCount, s.PostCount)
}
