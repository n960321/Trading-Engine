package engine

import (
	"Trading-Engine/internal/model"

	al "github.com/emirpasic/gods/lists/arraylist"
	"github.com/shopspring/decimal"
)

var QueueTypeArrayList = model.QueueType("ArrayList")

type ArrayList struct {
	list       *al.List
	side       model.OrderSide
	comparator func(a, b interface{}) int
}

func NewQueueByList(side model.OrderSide) OrderQueueior {

	list := al.New()
	var c func(a, b interface{}) int

	if side == model.OrderSideBuy {
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
	aKey := a.(model.Order)
	bKey := b.(model.Order)

	c := aKey.Price.Compare(bKey.Price)
	if c != 0 {
		return c * -1
	}
	return aKey.CreatedAt.Compare(bKey.CreatedAt)
}

func SellPriceTimeComparator(a, b interface{}) int {
	aKey := a.(model.Order)
	bKey := b.(model.Order)

	c := aKey.Price.Compare(bKey.Price)
	if c != 0 {
		return c
	}
	return aKey.CreatedAt.Compare(bKey.CreatedAt)
}

func (l *ArrayList) Size() int64 {
	return int64(l.list.Size())
}

func (l *ArrayList) AddOrder(order model.Order) {

	find := func(index int, value interface{}) bool {
		o := value.(model.Order)
		return o.ID == order.ID
	}
	i, _ := l.list.Find(find)

	if i != -1 {
		l.list.Remove(i)
	}
	l.list.Add(order)

	l.list.Sort(l.comparator)
}

func (l *ArrayList) RemoveOrder(order model.Order) {
	find := func(index int, value interface{}) bool {
		o := value.(model.Order)
		return o.ID == order.ID
	}
	i, _ := l.list.Find(find)

	if i == -1 {
		return
	}

	l.list.Remove(i)
}

func (l *ArrayList) GetOrdersBetweenPirceWithAmount(price, laveAmount decimal.Decimal) []model.Order {
	orders := make([]model.Order, 0)
	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(model.Order)

		if l.side == model.OrderSideBuy && order.Price.LessThan(price) || (l.side == model.OrderSideSell && order.Price.GreaterThan(price)) {
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

func (l *ArrayList) GetOrdersWithAmount(laveAmount decimal.Decimal) []model.Order {
	orders := make([]model.Order, 0)
	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(model.Order)
		laveAmount = laveAmount.Sub(order.GetLaveAmount())
		orders = append(orders, order)
		if laveAmount.LessThanOrEqual(decimal.Zero) {
			break
		}
	}

	return orders
}

func (l *ArrayList) GetAllOrders() []model.Order {
	orders := make([]model.Order, 0, l.list.Size())

	it := l.list.Iterator()
	for it.Next() {
		order := it.Value().(model.Order)
		orders = append(orders, order)
	}
	return orders
}
