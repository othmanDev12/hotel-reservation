package db

const (
	UriDb      = "mongodb://localhost:27017"
	Dbname     = "hotel-reservation"
	DbNameTest = "hotel-reservation-test"
)

type Store struct {
	HotelStore   HotelStore
	RoomStore    RoomStore
	UserStore    UserStore
	BookingStore BookingStore
}
