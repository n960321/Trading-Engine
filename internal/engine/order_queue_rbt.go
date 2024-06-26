package engine

import (
	"Trading-Engine/internal/model"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

var QueueTypePriceTree = model.QueueType("PriceTree")

type PriceTree struct {
	tree *rbt.Tree
	side model.OrderSide
}

func NewRBTPriceTree(side model.OrderSide) OrderQueueior {
	var tree *rbt.Tree
	switch side {
	case model.OrderSideBuy:
		tree = rbt.NewWith(BuySideComparator)
	case model.OrderSideSell:
		tree = rbt.NewWith(SellSideComparator)
	}
	return &PriceTree{
		tree: tree,
		side: side,
	}
}

func BuySideComparator(a, b interface{}) int {
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

func SellSideComparator(a, b interface{}) int {
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

func (p *PriceTree) Size() int64 {
	result := int64(0)
	it := p.tree.Iterator()

	for it.Next() {
		timeTree := it.Value().(*TimeTree)
		result += int64(timeTree.tree.Size())
	}

	return result
}

func (p *PriceTree) AddOrder(order model.Order) {
	var timeTree *TimeTree
	if tree, ok := p.tree.Get(order.Price); ok {
		timeTree = tree.(*TimeTree)
	} else {
		timeTree = NewTimeTree()
		p.tree.Put(order.Price, timeTree)
	}
	timeTree.AddOrder(order)
}

func (p *PriceTree) RemoveOrder(order model.Order) {
	if tree, ok := p.tree.Get(order.Price); ok {
		timeTree := tree.(*TimeTree)
		timeTree.RemoveOrder(order)

		if timeTree.GetLaveAmount().Equal(decimal.Zero) {
			p.tree.Remove(order.Price)
		}
	}
}

func (p *PriceTree) GetOrdersBetweenPirceWithAmount(price, laveAmount decimal.Decimal) []model.Order {
	orders := make([]model.Order, 0)
	it := p.tree.Iterator()
	for it.Next() {
		curPrice := it.Key().(decimal.Decimal)
		if p.side == model.OrderSideBuy && curPrice.LessThan(price) || (p.side == model.OrderSideSell && curPrice.GreaterThan(price)) {
			break
		}

		timeTree := it.Value().(*TimeTree)
		var subOrders []model.Order
		subOrders, laveAmount = timeTree.GetOrdersWithAmount(laveAmount)
		orders = append(orders, subOrders...)
		if laveAmount.LessThanOrEqual(decimal.Zero) {
			break
		}
	}

	return orders
}

func (p *PriceTree) GetOrdersWithAmount(laveAmount decimal.Decimal) []model.Order {
	orders := make([]model.Order, 0)
	it := p.tree.Iterator()
	for it.Next() {
		timeTree := it.Value().(*TimeTree)
		var subOrders []model.Order
		subOrders, laveAmount = timeTree.GetOrdersWithAmount(laveAmount)
		orders = append(orders, subOrders...)
		if laveAmount.LessThanOrEqual(decimal.Zero) {
			break
		}
	}

	return orders
}

func (p *PriceTree) GetAllOrders() []model.Order {
	orders := make([]model.Order, 0, p.Size())
	it := p.tree.Iterator()
	for it.Next() {
		timeTree := it.Value().(*TimeTree)
		timeTree.GetAllOrders(&orders)
	}
	return orders
}

type TimeTree struct {
	tree       *rbt.Tree
	laveAmount decimal.Decimal // 剩餘數量
}

func NewTimeTree() *TimeTree {
	return &TimeTree{
		tree:       rbt.NewWith(TimeComparator),
		laveAmount: decimal.Zero,
	}
}

func TimeComparator(a, b interface{}) int {
	aTime := a.(time.Time)
	bTime := b.(time.Time)
	return aTime.Compare(bTime)
}

func IDComparator(a, b interface{}) int {
	aID := a.(uint64)
	bID := b.(uint64)
	if aID <= bID {
		return 1
	}
	return -1
}

func (t *TimeTree) AddOrder(order model.Order) {
	t.tree.Put(order.CreatedAt, order)
	t.laveAmount = t.laveAmount.Add(order.GetLaveAmount())
}

func (t *TimeTree) RemoveOrder(order model.Order) {
	t.tree.Remove(order.CreatedAt)
	t.laveAmount = t.laveAmount.Sub(order.GetLaveAmount())
}

func (t *TimeTree) GetAllOrders(orders *[]model.Order) {
	it := t.tree.Iterator()
	for it.Next() {
		order := it.Value().(model.Order)
		*orders = append(*orders, order)
	}
}

func (t *TimeTree) GetOrdersWithAmount(laveAmount decimal.Decimal) ([]model.Order, decimal.Decimal) {
	orders := make([]model.Order, 0)
	it := t.tree.Iterator()
	for it.Next() {
		order := it.Value().(model.Order)
		laveAmount = laveAmount.Sub(order.GetLaveAmount())
		orders = append(orders, order)
		if !laveAmount.GreaterThan(decimal.Zero) {
			break
		}
	}
	return orders, laveAmount
}

func (t *TimeTree) GetLaveAmount() decimal.Decimal {
	return t.laveAmount
}
