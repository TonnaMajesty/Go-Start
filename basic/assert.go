package main

import (
	"context"
	"fmt"
)

type Action = func(ctx context.Context) error

type WithPostActions interface {
	PostActions() []Action
}

type ScanIterator interface {
	// New a ptr value for scan
	New() interface{}
	// Next For receive scanned value
	Next(v interface{}) error
}

func List(ctx context.Context, list ScanIterator) error {
	if postActions, ok := list.(WithPostActions); ok {
		postActions.PostActions()
	}

	if length, ok := list.(interface{ Len() int }); ok {
		fmt.Println(length.Len())
	}

	return nil
}

type userModel struct {
}

func (ListDemo) New() interface{} {
	return &userModel{}
}

func (l *ListDemo) Next(v interface{}) error {
	c := v.(*userModel)
	l.Data = append(l.Data, c)
	return nil
}

type ListDemo struct {
	Data []*userModel
}

func (l *ListDemo) Len() int {
	return len(l.Data)
}

func (l *ListDemo) PostActions() []Action {
	return []Action{}
}

func main() {
	_ = List(context.Background(), &ListDemo{})
}
