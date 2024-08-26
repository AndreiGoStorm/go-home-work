package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push front item", func(t *testing.T) {
		l := NewList()
		l.PushFront(10) // [10]

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())
	})

	t.Run("push back item", func(t *testing.T) {
		l := NewList()
		l.PushBack(20) // [20]

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())
	})

	t.Run("push front and back items", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				l.PushFront(i * 10)
			} else {
				l.PushBack(i * 100)
			}
		} // [80,60,40,20,0,100,300,500,700,900]
		l.PushBack(15) // [80,60,40,20,0,100,300,500,700,900,15]

		require.Equal(t, 11, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 15, l.Back().Value)
	})

	t.Run("remove front item", func(t *testing.T) {
		l := NewList()
		l.PushFront(10) // [10]
		l.Remove(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("remove back item", func(t *testing.T) {
		l := NewList()
		l.PushBack(20) // [20]
		l.Remove(l.Back())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("remove front and back items", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 10; i++ {
			l.PushBack(i * 10) // [0,10,20,30,40,50,60,70,80,90]
		}
		item := l.PushBack(100) // [0,10,20,30,40,50,60,70,80,90,100]
		l.PushBack(200)         // [0,10,20,30,40,50,60,70,80,90,100,200]
		l.Remove(l.Front())
		l.Remove(item)
		l.Remove(l.Back())

		require.Equal(t, 9, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 90, l.Back().Value)
	})

	t.Run("move to front item", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)        // [10]
		item := l.PushBack(20) // [10, 20]
		l.PushFront(30)        // [30, 10, 20]
		l.PushBack(40)         // [30, 10, 20, 40]
		l.MoveToFront(item)

		require.Equal(t, 4, l.Len())
		require.Equal(t, item.Value, l.Front().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
