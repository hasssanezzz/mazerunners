package server

import "fmt"

type ErrRoomNotFound struct {
	RoomID int
}

func (e ErrRoomNotFound) Error() string {
	return fmt.Sprintf("room not found: %d", e.RoomID)
}
