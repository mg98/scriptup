package storage

type Order string

const (
	// OrderAsc for ascending order.
	OrderAsc Order = "ASC"

	// OrderDesc for descending order.
	OrderDesc Order = "DESC"
)

// Storage represents a stack-like (i.e., simple string values in a LIFO manner) persistence management interface.
type Storage interface {
	// Open a connection to the storage medium.
	Open() error

	// Close connection safely. This should be deferred after a call of Open.
	Close() error

	// Append (and persist) a new entry to the storage.
	Append(string) error

	// Pop removes the most recently added entry from the storage.
	Pop() error

	// All returns all items from the storage in ascending order.
	All(Order) ([]string, error)

	// Latest returns the latest/most recent item from the storage.
	Latest() (*string, error)
}
