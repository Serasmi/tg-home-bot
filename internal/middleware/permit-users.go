package middleware

import (
	tele "gopkg.in/telebot.v3"
)

func PermitUsers(IDs []int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !hasPermit(c.Sender().ID, IDs) {
				return nil
			}

			return next(c) // continue execution chain
		}
	}
}

func hasPermit(userId int64, IDs []int64) bool {
	if len(IDs) == 0 {
		return true
	}

	for _, ID := range IDs {
		if ID == userId {
			return true
		}
	}

	return false
}
