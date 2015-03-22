package admin

type Show struct {
	ID         int
	CategoryID int
	ChannelID  int
	Title      string
	Thumbnail  string
	Detail     string
	IsActive   bool
	IsOnlive   bool
}
