package money

import "fmt"

type Money int64

func New(amount int64) Money {
	return Money(amount)
}

func (m Money) Int64() int64 {
	return int64(m)
}

func (m Money) String() string {
	rubles := int64(m) / 100
	kopecks := int64(m) % 100
	if kopecks < 0 {
		kopecks = -kopecks
	}
	return fmt.Sprintf("%d.%02d", rubles, kopecks)
}