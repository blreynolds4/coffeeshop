package models

// Test Helpers
func getTestGrinders() GrinderPool {
	return NewGrinderPool(&MockGrinder{})
}

func getTestBrewers() BrewerPool {
	return NewBrewerPool(&MockBrewer{})
}

func getTestMenuItem() MenuItem {
	return MenuItem{
		Name:        "Regular Coffee",
		Size:        8,
		CoffeeRatio: 2,
	}
}

type MockGrinder struct{}

func (mg *MockGrinder) Grind(b Beans) Beans {
	return b
}

type MockBrewer struct{}

func (mb *MockBrewer) Brew(finishedVolume int, beans Beans) *Coffee {
	return &Coffee{sizeOunces: finishedVolume}
}
