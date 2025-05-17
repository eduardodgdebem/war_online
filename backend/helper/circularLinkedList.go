package helper

import (
	"fmt"
)

type Node struct {
	data uint32
	next *Node
}

type CicularLinkedList struct {
	head   *Node
	tail   *Node
	length int
}

func (list *CicularLinkedList) Display() {

	if list.head == nil {
		fmt.Printf("No Data Present in Linked List.\n")
	} else {
		temp := list.head
		for temp != nil {
			fmt.Printf("%v -> ", temp.data)
			temp = temp.next
			if temp == list.tail.next {
				break
			}
		}
		fmt.Println("END")
	}
}

func (list *CicularLinkedList) ListLength() int {
	return list.length
}

func (list *CicularLinkedList) InsertBeginning(n *Node) {

	if list.head == nil {
		list.head = n
		list.tail = list.head
		list.tail.next = list.head
	} else {
		n.next = list.head
		list.head = n
		list.tail.next = list.head
	}
	list.length++

}

func (list *CicularLinkedList) InsertEnd(n *Node) {

	if list.head == nil {
		list.InsertBeginning(n)
	} else {
		n.next = list.head
		list.tail.next = n
		list.tail = n
		list.length++
	}

}

func (list *CicularLinkedList) DeleteFromBegining() {

	if list.head == nil {
		fmt.Printf("No Data Present in Linked List.\n")
	} else {
		list.head = list.head.next
		list.tail.next = list.head
		list.length--
	}

}

func (list *CicularLinkedList) DeleteFromEnd() {

	if list.head == nil {
		fmt.Printf("No Data Present in Linked List.\n")
	} else {
		temp := list.head
		var prev *Node = nil

		for temp.next != list.head {
			prev = temp
			temp = temp.next
		}

		prev.next = list.head
		list.tail = prev
		list.length--
	}

}

func (list *CicularLinkedList) DeleteLinkedList() {

	if list.head != nil {
		list.tail = nil
		list.head = nil
		list.length = 0
	}

}
