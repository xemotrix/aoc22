package evilgo

type Dequeue[T any] interface {
	PushHead(items ...T)
	PushTail(items ...T)
	PopHead() T
	PopTail() T
	ApplyToAll(fun func(T) T)
	Empty() bool
	Len() int
	SetCursor()
}

type DequeueNode[T any] struct {
	item T
	next *DequeueNode[T]
	prev *DequeueNode[T]
}

type dequeue[T any] struct {
	head   *DequeueNode[T]
	tail   *DequeueNode[T]
	cursor *DequeueNode[T]
}

func BuildDequeue[T any](items ...T) Dequeue[T] {
	if len(items) == 0 {
		return &dequeue[T]{}
	}

	node := DequeueNode[T]{item: items[0]}

	dq := dequeue[T]{
		head: &node,
		tail: &node,
	}

	dq.PushTail(items[1:]...)

	return &dq
}

func (q *dequeue[T]) pushHead(item T) {
	if q.Len() == 0 {
		node := DequeueNode[T]{item: item}
		q.head = &node
		q.tail = &node
		return
	}
	node := &DequeueNode[T]{
		item: item,
		next: q.head,
	}
	q.head.prev = node
	q.head = node
}

func (q *dequeue[T]) PushHead(items ...T) {
	for _, it := range items {
		q.pushHead(it)
	}
}

func (q *dequeue[T]) SetCursor() {
	q.cursor = q.head
}

func (q *dequeue[T]) Next() (T, bool) {
	if q.cursor == nil {
		var a T
		return a, false
	}
	res := q.cursor.item
	q.cursor = q.cursor.next
	return res, true
}

func (q *dequeue[T]) pushTail(item T) {
	if q.Len() == 0 {
		node := DequeueNode[T]{item: item}
		q.head = &node
		q.tail = &node
		return
	}
	node := &DequeueNode[T]{
		item: item,
		prev: q.tail,
	}
	q.tail.next = node
	q.tail = node
}

func (q *dequeue[T]) PushTail(items ...T) {
	for _, it := range items {
		q.pushTail(it)
	}
}

func (q *dequeue[T]) ApplyToAll(fun func(T) T) {
	cursor := q.head
	for cursor != nil {
		cursor.item = fun(cursor.item)
		cursor = cursor.next
	}
}

func (q *dequeue[T]) PopHead() T {
	node := q.head
	if node == nil {
		panic("can't pop from empty queue")
	}

	next := q.head.next
	if next != nil {
		next.prev = nil
	}
	q.head = q.head.next
	return node.item
}

func (q *dequeue[T]) PopTail() T {
	node := q.tail
	if node == nil {
		panic("can't pop from empty queue")
	}

	prev := node.prev
	if prev != nil {
		prev.next = nil
	}
	q.tail = q.tail.prev
	return node.item
}

func (q *dequeue[T]) Empty() bool {
	return q.head == nil
}

func (q *dequeue[T]) Len() int {
	if q.Empty() {
		return 0
	}
	count := 0
	cursor := q.head
	for cursor != nil {
		count++
		cursor = cursor.next
	}
	return count
}
