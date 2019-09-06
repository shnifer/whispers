package juggler

import (
	"errors"
	"github.com/shnifer/whispers/rid"
	"sync"
)

var (
	ErrNotFound   = errors.New("juggler: not found")
	ErrDropTaken  = errors.New("juggler: can't drop taken object")
	ErrNotTaken   = errors.New("juggler: can't free not taken")
	ErrWrongOwner = errors.New("juggler: wrong owner")
)

type transID = rid.U32

type object struct {
	owner   transID
	waiter  transID
	taken   bool
	waited  bool
	deleted bool
}

type Objects struct {
	interrupt func(transID)

	*sync.Cond
	objects []object
	frees   []int
}

func New(interrupt func(transID)) *Objects {
	if interrupt == nil {
		interrupt = func(transID) {}
	}
	return &Objects{
		Cond:      sync.NewCond(&sync.Mutex{}),
		interrupt: interrupt,
		objects:   make([]object, 0),
		frees:     make([]int, 0),
	}
}

func (j *Objects) Add() int {
	j.L.Lock()
	defer j.L.Unlock()

	for len(j.frees) > 0 {
		L := len(j.frees)
		ind := j.frees[L-1]
		if ind > len(j.objects) {
			j.frees = j.frees[:L-1]
			continue
		}
		j.frees = j.frees[:L-1]
		j.objects[ind] = object{}
		return ind
	}

	j.objects = append(j.objects, object{})
	return len(j.objects) - 1
}

func (j *Objects) Drop(ind int) error {
	j.L.Lock()
	defer j.L.Unlock()

	if err := j.checkInd(ind); err != nil {
		return err
	}
	if j.objects[ind].taken || j.objects[ind].waited {
		return ErrDropTaken
	}

	if ind == len(j.objects)-1 {
		j.objects = j.objects[:len(j.objects)-1]
	} else {
		j.objects[ind].deleted = true
		j.frees = append(j.frees, ind)
	}
	return nil
}

//Take control of object [ind].
// If it is already taken by owner with smaller rid -- wait
// If it is already taken by owner with bigger rid -- throw owner away and gain control
func (j *Objects) Take(transaction transID, ind int) error {
	j.L.Lock()
	defer j.L.Unlock()

	for {
		if err := j.checkInd(ind); err != nil {
			return err
		}
		obj := j.objects[ind]

		if !obj.taken {
			if obj.waited && obj.waiter.Less(transaction) {
				//object is not taken, but have earlier waiter, wait for him
				j.Wait()
				continue
			} else {
				//object is free for us , just grab it
				j.objTaken(transaction, ind)
				j.Broadcast()
				return nil
			}
		} else {
			if obj.owner == transaction {
				//you could take twice idempotent
				return nil
			} else if obj.owner.Less(transaction) {
				//object is taken by earlier transaction, stay in queue and wait
				j.objWaiter(transaction, ind)
				j.Wait()
				continue
			} else {
				//object was taken by later transaction, we have to interrupt it and wait for free
				j.interrupt(obj.owner)
				j.objWaiter(transaction, ind)
				j.Wait()
				continue
			}
		}
	}
}

func (j *Objects) Free(transaction transID, ind int) error {
	j.L.Lock()
	defer j.L.Unlock()
	if err := j.checkInd(ind); err != nil {
		return err
	}
	obj := j.objects[ind]
	if !obj.taken {
		return ErrNotTaken
	}
	if obj.owner != transaction {
		return ErrWrongOwner
	}
	obj.taken = false
	j.Broadcast()
	return nil
}

func (j *Objects) checkInd(ind int) error {
	if ind < 0 || ind >= len(j.objects) || j.objects[ind].deleted {
		return ErrNotFound
	}
	return nil
}

func (j *Objects) objTaken(transaction transID, ind int) {
	if j.objects[ind].waited && j.objects[ind].waiter == transaction {
		j.objects[ind].waited = false
	}
	j.objects[ind].taken = true
	j.objects[ind].owner = transaction
}

func (j *Objects) objWaiter(transaction transID, ind int) {
	if !j.objects[ind].waited {
		j.objects[ind].waited = true
		j.objects[ind].waiter = transaction
	} else if transaction.Less(j.objects[ind].waiter) {
		j.objects[ind].waiter = transaction
	}
}
