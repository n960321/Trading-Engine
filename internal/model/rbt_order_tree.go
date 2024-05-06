package model

import (
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type PriceTree struct {
	tree *rbt.Tree
}

func NewPriceTree(side OrderSide) *PriceTree {
	var tree *rbt.Tree
	switch side {
	case Buy:
		tree = rbt.NewWith(BuySideComparator)
	case Sell:
		tree = rbt.NewWith(SellSideComparator)
	}
	return &PriceTree{
		tree: tree,
	}
}

func BuySideComparator(a, b interface{}) int {
	aPrice := a.(decimal.Decimal)
	bPrice := b.(decimal.Decimal)

	if aPrice.Equal(bPrice) {
		return 0
	}

	if aPrice.GreaterThan(bPrice) {
		return 1
	}

	return -1
}

func SellSideComparator(a, b interface{}) int {
	aPrice := a.(decimal.Decimal)
	bPrice := b.(decimal.Decimal)

	if aPrice.Equal(bPrice) {
		return 0
	}

	if aPrice.GreaterThan(bPrice) {
		return -1
	}

	return 1
}

func (p *PriceTree) Size() int64 {
	result := int64(0)
	it := p.tree.Iterator()

	for it.Next() {
		timeTree := it.Value().(*TimeTree)
		result += int64(timeTree.tree.Size())
	}

	return result
}

func (p *PriceTree) AddOrder(order *Order) {
	var timeTree *TimeTree
	if tree, ok := p.tree.Get(order.Price); ok {
		timeTree = tree.(*TimeTree)
	} else {
		timeTree = NewTimeTree()
		p.tree.Put(order.Price, timeTree)
	}

	timeTree.tree.Put(order.CreatedAt, order)
}

func (p *PriceTree) RemoveOrder(order *Order) {
	if tree, ok := p.tree.Get(order.Price); ok {
		timeTree := tree.(*TimeTree)
		timeTree.tree.Remove(order.CreatedAt)

		if timeTree.tree.Empty() {
			p.tree.Remove(order.Price)
		}
	}
}

func (p *PriceTree) GetAllOrders() []*Order {
	orders := make([]*Order, 0, p.Size())
	it := p.tree.Iterator()
	for it.Next() {
		timeTree := it.Value().(*TimeTree)
		timeTree.GetAllOrders(&orders)
	}
	return orders
}

type TimeTree struct {
	tree *rbt.Tree
}

func NewTimeTree() *TimeTree {
	return &TimeTree{
		tree: rbt.NewWith(TimeComparator),
	}
}

func TimeComparator(a, b interface{}) int {
	aTime := a.(time.Time)
	bTime := b.(time.Time)
	return aTime.Compare(bTime)
}

func (t *TimeTree) GetAllOrders(orders *[]*Order) {
	it := t.tree.Iterator()
	for it.Next() {
		order := it.Value().(*Order)
		*orders = append(*orders, order)
	}
}
