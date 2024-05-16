package middleware

import (
	tele "gopkg.in/telebot.v3"
)

func PermitUsers(ids []int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !hasPermit(c.Sender().ID, ids) {
				return nil
			}

			return next(c) // continue execution chain
		}
	}
}

func hasPermit(userID int64, ids []int64) bool {
	if len(ids) == 0 {
		return true
	}

	for _, ID := range ids {
		if ID == userID {
			return true
		}
	}

	return false
}
