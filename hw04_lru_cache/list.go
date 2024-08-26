package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Head  *ListItem
	Tail  *ListItem
	Count int
}

func (l *list) Len() int {
	return l.Count
}

func (l *list) Front() *ListItem {
	if l.Count == 0 {
		return nil
	}
	return l.Head
}

func (l *list) Back() *ListItem {
	if l.Count == 0 {
		return nil
	}
	return l.Tail
}

func (l *list) PushFront(v any) *ListItem {
	item := NewListItem(v)
	if l.Head == nil {
		l.Head = item
		l.Tail = item
	} else {
		item.Next = l.Head
		l.Head.Prev = item
		l.Head = item
	}
	l.Count++
	return l.Head
}

func (l *list) PushBack(v any) *ListItem {
	item := NewListItem(v)
	if l.Head == nil {
		l.Head = item
		l.Tail = item
	} else {
		item.Prev = l.Tail
		l.Tail.Next = item
		l.Tail = item
	}
	l.Count++
	return l.Tail
}

func (l *list) Remove(e *ListItem) {
	if e.Prev != nil {
		e.Prev.Next = e.Next
	} else {
		l.Head = e.Next
	}

	if e.Next != nil {
		e.Next.Prev = e.Prev
	} else {
		l.Tail = e.Prev
	}

	e.Next = nil
	e.Prev = nil
	l.Count--
}

func (l *list) MoveToFront(e *ListItem) {
	l.Remove(e)
	l.PushFront(e.Value)
}

func NewListItem(v any) *ListItem {
	return &ListItem{Value: v}
}

func NewList() List {
	return new(list)
}
