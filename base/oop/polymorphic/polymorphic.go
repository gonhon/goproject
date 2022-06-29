package polymorphic

import "fmt"

type Income interface {
	calculate() int
	source() string
}

type FixedBilling struct {
	projectName  string
	biddedAmount int
}

func (fix *FixedBilling) calculate() int {
	return fix.biddedAmount
}

func (fix *FixedBilling) source() string {
	return fix.projectName
}

type TimeAndMaterial struct {
	projectName string
	noOfHours   int
	hourlyRate  int
}

func (t *TimeAndMaterial) calculate() int {
	return t.noOfHours * t.hourlyRate
}
func (t *TimeAndMaterial) source() string {
	return t.projectName
}

func calculateNetIncome(arrays []Income) {
	netincome := 0
	for _, v := range arrays {
		netincome += v.calculate()
		fmt.Printf("Income From %s = $%d\n", v.source(), v.calculate())
	}
	fmt.Printf("Net income of organisation = $%d", netincome)
}
