package types

// User is a user on Telegram.
type User struct {
	Id                      int     // required
	IsBot                   bool    // required
	FirstName               string  // required
	LastName                string 
	UserName                string 
	LanguageCode            string 
	CanJoioGroups           bool   
	CanReadAllGroupMessages bool   
	SupportsInlineQueries   bool   
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
func (u *User) String() string {
	if u.UserName != "" {
		return u.UserName
	}

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	return name
}
