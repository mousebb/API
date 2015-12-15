package contact

type ContactReceiver struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	ContactTypes []ContactType
}
