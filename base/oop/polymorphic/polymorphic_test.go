package polymorphic

import "testing"

func TestIncome(t *testing.T) {
	project1 := &FixedBilling{projectName: "Project 1", biddedAmount: 5000}
	project2 := &FixedBilling{projectName: "Project 2", biddedAmount: 10000}
	project3 := &TimeAndMaterial{projectName: "Project 3", noOfHours: 160, hourlyRate: 25}
	incomeStreams := []Income{project1, project2, project3}
	calculateNetIncome(incomeStreams)
}
