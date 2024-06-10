package app

type Commission struct {
	ID          int
	Name        string
	Description string
	Budget      string
	Currency    string
	Location    string
	Deadline    string
	Status      string
	Winner      string
}

func (c *Commission) Create() error {
	// create a commission
	return nil
}
