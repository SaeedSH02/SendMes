package middle

import "gopkg.in/telebot.v4"

var allowedUsers = map[int64]bool{
	//set of allowed users with TEL IDs
	1234567890: true,
	987654321:  true,
	1122334455: true,
}

func UserAllowed(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		if !allowedUsers[userID] {
			return c.Send("🚫 شما دسترسی به این ربات ندارید 🚫")
		}
		return next(c)
	}
}
