package main

import (
	"sync"

	"github.com/pkg/errors"
)

var placeholder = struct{}{}

var DefaultRoomHub = RoomHub{
	rooms: map[string]*Room{},
}

var (
	ErrClientIDConflict = errors.New("client id conflict")
	ErrRoomIDConflict   = errors.New("room id conflict")
	ErrRoomNotFound     = errors.New("room not found")
)

type RoomHub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func (h *RoomHub) GetRoom(ID string) (*Room, bool) {
	h.mu.RLock()
	r, ok := h.rooms[ID]
	h.mu.RUnlock()

	return r, ok
}

func (h *RoomHub) CreateRoom(ID string) (*Room, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.rooms[ID]
	if !ok {
		r := NewRoom(ID)
		h.rooms[ID] = r
		return r, nil
	}

	return nil, ErrRoomIDConflict
}

func (h *RoomHub) GetOrCreateRoom(ID string) (*Room, error) {
	r, ok := h.GetRoom(ID)
	if !ok {
		return h.CreateRoom(ID)
	}

	return r, nil
}
func (h *RoomHub) DestoryRoom(ID string) error {
	// 姑且只是把 room 从 hub 中移除
	h.mu.Lock()
	delete(h.rooms, ID)
	h.mu.Unlock()

	return nil
}

func (h *RoomHub) BoardCastAll(msg IMessage, opt *BoardCastOpt) error {
	h.mu.RLock()
	defer h.mu.Unlock()

	errs := &BundleErr{}

	for k, r := range h.rooms {
		err := r.BoardCast(msg, opt)
		if err != nil {
			errs.AddErr(errors.Wrapf(err, "boardCastAll fail : %s", k))
		}
	}

	return errs.Bundle()
}

func (h *RoomHub) BoardCast(msg IMessage, opt *BoardCastOpt, roomIds ...string) error {
	errs := &BundleErr{}

	// 这么大的读锁，可能会造成写阻塞
	h.mu.RLock()
	for _, roomid := range roomIds {
		r, ok := h.rooms[roomid]
		if !ok {
			errs.AddErr(errors.Wrap(ErrRoomNotFound, roomid))
			continue
		}

		err := r.BoardCast(msg, opt)
		if err != nil {
			errs.AddErr(errors.Wrapf(err, "board cast rooms fail: %s", r.id))
		}
	}
	h.mu.RUnlock()

	return errs.Bundle()
}

func NewRoom(ID string) *Room {
	return &Room{
		id:      ID,
		clients: map[string]*Client{},
	}
}

type Room struct {
	mu      sync.Mutex
	id      string
	clients map[string]*Client
}

func (r *Room) JoinRoom(c *Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.clients[c.ID()]
	if !ok {
		r.clients[c.ID()] = c
		return nil
	}

	return errors.Wrap(ErrClientIDConflict, c.ID())
}

func (r *Room) Leave(ID string) error {
	r.mu.Lock()
	delete(r.clients, ID)
	r.mu.Unlock()

	return nil
}

type BoardCastOpt struct {
	ExceptIDs []string
}

func (r *Room) BoardCast(msg IMessage, opt *BoardCastOpt) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var excepts map[string]struct{}
	if opt != nil && opt.ExceptIDs != nil {
		excepts = make(map[string]struct{})

		for _, id := range opt.ExceptIDs {
			excepts[id] = placeholder
		}
	}

	errs := &BundleErr{}

	for k, v := range r.clients {
		if _, ok := excepts[k]; ok {
			continue
		}

		err := v.Send(msg)
		if err != nil {
			errs.AddErr(errors.Wrapf(err, "boardcast fail of : %s", k))
		}
	}

	return errs.Bundle()
}
