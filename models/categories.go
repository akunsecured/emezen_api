package models

type Category int64

const (
	ArtsCrafts Category = iota
	Books
	Electronics
	Fashion
)

func (c Category) String() string {
	switch c {
	case ArtsCrafts:
		return "Arts & Crafts"
	case Books:
		return "Books"
	case Electronics:
		return "Electronics"
	case Fashion:
		return "Fashion"
	}
	return "Unknown"
}
