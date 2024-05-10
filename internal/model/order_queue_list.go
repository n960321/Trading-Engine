package model

import (
	al "github.com/emirpasic/gods/lists/arraylist"
	"github.com/shopspring/decimal"
)

var QueueTypeArrayList = QueueType("ArrayList")
type ArrayList struct {
	list       *al.List
	side       OrderSide
	comparator func(a, b interface{}) int
}

func NewQueueByList(side OrderSide) OrderQueueior {
	var list *al.List
	list = al.New()
	var c func(a, b interface{}) int

	if side == OrderSideBuy {
		c = BuyPriceTimeComparator
	} else {
		c = SellPriceTimeComparator
	}

	return &ArrayList{
		list:       list,
		comparator: c,
		side:       side,
	}
}

func BuyPriceTimeComparator(a, b interface{}) int {
	aKey := a.(*Order)
	bKey := b.(*Order)

	c := aKey.Price.Compare(bKey.Price)
	if c != 0 {
		return c*-1
	}
	return aKey.CreatedAt.Compare(bKey.CreatedAt)
}

func SellPriceTimeComparator(a, b interface{}) int {
	aKey := a.(*Order)
	bKey := b.(*Order)

	c := aKey.Price.Compare(bKey.Price)
	if c != 0 {
		return c 
	}
	return aKey.CreatedAt.Compare(bKey.CreatedAt)
}

func (l *ArrayList) Size() int64 {
	return int64(l.list.Size())
}

func (l *ArrayList) AddOrder(order *Order) {
	l.list.Add(order)
	l.list.Sort(l.comparator)
}

func (l *ArrayList) RemoveOrder(order *Order) {
	find := func(index int, value interface{}) bool {
		o := value.(*Order)
		return o.ID == order.ID
	}
	i, _ := l.list.Find(find)

	if i == -1 {
		return
	}

	l.list.Remove(i)
}

func (l *ArrayList) GetOrdersBetweenPirceWithAmount(price, laveAmount decimal.Decimal) []*Order {
	orders := make([]*Order, 0)
	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(*Order)

		if l.side == OrderSideBuy && order.Price.LessThan(price) || (l.side == OrderSideSell && order.Price.GreaterThan(price)) {
			break
		}

		laveAmount = laveAmount.Sub(order.GetLaveAmount())
		orders = append(orders, order)
		if laveAmount.LessThanOrEqual(decimal.Zero) {
			break
		}
	}

	return orders
}

func (l *ArrayList) GetOrdersWithAmount(laveAmount decimal.Decimal) []*Order {
	orders := make([]*Order, 0)
	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(*Order)
		laveAmount = laveAmount.Sub(order.GetLaveAmount())
		orders = append(orders, order)
		if laveAmount.LessThanOrEqual(decimal.Zero) {
			break
		}
	}

	return orders
}

func (l *ArrayList) GetAllOrders() []*Order {
	orders := make([]*Order, 0, l.list.Size())

	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(*Order)
		orders = append(orders, order)
	}
	return orders
}
