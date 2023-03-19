package _type

import "fmt"

func FormatPrice(priceInKopeks uint64) string {
	rubles := priceInKopeks / 100
	kopeks := priceInKopeks % 100
	if kopeks == 0 {
		return fmt.Sprintf("%d ₽", rubles)
	} else {
		return fmt.Sprintf("%d.%02d ₽", rubles, kopeks)
	}
}
