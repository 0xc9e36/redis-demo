/**
 * @Author: wei.tan
 * @Description: 最近最少使用算法, 使用 map + 双链表实现  http://bikong0411.github.io/2016/06/29/lru-go.html
 * @File:  lru
 * @Version: 1.0.0
 * @Date: 2019-10-18 21:15
 */

package lru

type Node struct {
	key, value interface{}
	prev, next *Node
}

type LRUCache struct {
	capacity   int
	cacheMap   map[interface{}]*Node
	head, tail *Node
}

func NewLRUCache(capacity int) *LRUCache {
	lru := &LRUCache{}
	lru.tail = &Node{}
	lru.head = &Node{}
	lru.capacity = capacity
	lru.cacheMap = make(map[interface{}]*Node)
	lru.head.next = lru.tail
	lru.head.prev = nil
	lru.tail.next = nil
	lru.tail.prev = lru.head
	return lru
}

// h == a == b  == a
func (l *LRUCache) delete(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (l *LRUCache) attach(node *Node) {
	node.next = l.head.next
	node.next.prev = node
	l.head.next = node
	node.prev = l.head
}

func (l *LRUCache) Get(key interface{}) interface{} {
	node, ok := l.cacheMap[key]
	if !ok {
		return -1
	}

	l.delete(node)
	l.attach(node)
	return node.value
}

func (l *LRUCache) Put(key, value interface{}) {

	oldNode, ok := l.cacheMap[key]

	//如果存在 key
	if ok {
		l.delete(oldNode)
		l.attach(oldNode)
		oldNode.value = value
	} else {
		var node *Node
		//没有剩余空间, 删除最近最少使用节点
		if len(l.cacheMap) >= l.capacity {
			node = l.tail.prev
			l.delete(node)
			delete(l.cacheMap, node.key)
		} else {
			node = new(Node)
		}

		node.key = key
		node.value = value

		l.cacheMap[key] = node

		//放到头结点
		l.attach(node)
	}
}
