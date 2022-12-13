package evilgo

type PQueue[T comparable] interface {
	Push(items ...T)
	Pop() T
	Empty() bool
	Len() int
	Update(item T)
}

type pqueueNode[T comparable] struct {
	item T
	next *pqueueNode[T]
}

type pqueue[T comparable] struct {
	head *pqueueNode[T]
	lt   func(a, b T) bool
}

func BuildPQueue[T comparable](lt func(a, b T) bool, items ...T) PQueue[T] {
	pq := &pqueue[T]{
		lt: lt,
	}
	if len(items) == 0 {
		return pq
	}

	node := pqueueNode[T]{item: items[0]}

	pq.head = &node

	pq.Push(items[1:]...)

	return pq
}

func (q *pqueue[T]) Update(item T) {
	currentNode := q.head
	var lastNode *pqueueNode[T]
	for currentNode != nil {
		if currentNode.item == item {
			if lastNode == nil { // first item in list
				q.Push(q.Pop())
				break
			}
			toPush := currentNode.item
			lastNode.next = currentNode.next
			q.Push(toPush)
			break

		}
		lastNode = currentNode
		currentNode = currentNode.next
	}
}

func (q *pqueue[T]) push(item T) {
	if q.Empty() {
		node := pqueueNode[T]{item: item}
		q.head = &node
		return
	}
	itemNode := &pqueueNode[T]{
		item: item,
	}
	currentNode := q.head
	for {
		nextNode := currentNode.next

		if q.lt(item, currentNode.item) { //only possible if item < head
			itemNode.next = q.head
			q.head = itemNode
			break
		} else if nextNode == nil {
			currentNode.next = itemNode
			break
		} else if q.lt(item, nextNode.item) {
			currentNode.next = itemNode
			itemNode.next = nextNode
			break
		}
		currentNode = nextNode
	}
}

func (q *pqueue[T]) Push(items ...T) {
	for _, it := range items {
		q.push(it)
	}
}

func (q *pqueue[T]) Pop() T {
	node := q.head
	if node == nil {
		panic("can't pop from empty queue")
	}
	q.head = node.next
	return node.item
}

func (q *pqueue[T]) Empty() bool {
	return q.head == nil
}

func (q *pqueue[T]) Len() int {
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
