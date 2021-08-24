// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package event

import (
	"fmt"
	"sync"
	"time"
)

// Target - event target interface
type Target interface {
	ID() TargetID
	IsActive() (bool, error)
	Save(Event) error
	Send(string) error
	Close() error
	HasQueueStore() bool
}

// TargetList - holds list of targets indexed by target ID.
type TargetList struct {
	sync.RWMutex
	targets map[TargetID]Target
}

// Add - adds unique target to target list.
func (list *TargetList) Add(targets ...Target) error {
	list.Lock()
	defer list.Unlock()

	for _, target := range targets {
		if _, ok := list.targets[target.ID()]; ok {
			return fmt.Errorf("target %v already exists", target.ID())
		}
		list.targets[target.ID()] = target
	}

	return nil
}

// Exists - checks whether target by target ID exists or not.
func (list *TargetList) Exists(id TargetID) bool {
	list.RLock()
	defer list.RUnlock()

	_, found := list.targets[id]
	return found
}

// TargetIDResult returns result of Remove/Send operation, sets err if
// any for the associated TargetID
type TargetIDResult struct {
	// ID where the remove or send were initiated.
	ID TargetID
	// Stores any error while removing a target or while sending an event.
	Err error
}

// Remove - closes and removes targets by given target IDs.
func (list *TargetList) Remove(targetIDSet TargetIDSet) {
	list.Lock()
	defer list.Unlock()

	for id := range targetIDSet {
		target, ok := list.targets[id]
		if ok {
			target.Close()
			delete(list.targets, id)
		}
	}
}

// Targets - list all targets
func (list *TargetList) Targets() []Target {
	if list == nil {
		return []Target{}
	}

	list.RLock()
	defer list.RUnlock()

	targets := []Target{}
	for _, tgt := range list.targets {
		targets = append(targets, tgt)
	}

	return targets
}

// List - returns available target IDs.
func (list *TargetList) List() []TargetID {
	list.RLock()
	defer list.RUnlock()

	keys := []TargetID{}
	for k := range list.targets {
		keys = append(keys, k)
	}

	return keys
}

// TargetMap - returns available targets.
func (list *TargetList) TargetMap() map[TargetID]Target {
	list.RLock()
	defer list.RUnlock()
	return list.targets
}

// Send - sends events to targets identified by target IDs.
func (list *TargetList) Send(event Event, targetIDset TargetIDSet, resCh chan<- TargetIDResult) {
	go func() {
		var wg sync.WaitGroup
		for id := range targetIDset {
			list.RLock()
			target, ok := list.targets[id]
			list.RUnlock()
			if ok {
				wg.Add(1)
				go func(id TargetID, target Target) {
					defer wg.Done()
					tgtRes := TargetIDResult{ID: id}
					if err := target.Save(event); err != nil {
						tgtRes.Err = err
					}
					resCh <- tgtRes
				}(id, target)
			} else {
				resCh <- TargetIDResult{ID: id}
			}
		}
		wg.Wait()
	}()
}

// 同步发送事件到指定目标
// done - 通过channel返回事件是否发送完成
// Synchronization Send - sends events to targets identified by target IDs.
func (list *TargetList) SendSync(event Event, targetIDset TargetIDSet, resCh chan<- TargetIDResult, done chan bool) {
	// 声明一个等待组
	var wg sync.WaitGroup
	for id := range targetIDset {
		fmt.Printf("SendSync: id=%s, %d\n", id, time.Now().Unix())
		list.RLock()
		target, ok := list.targets[id]
		list.RUnlock()
		if ok {
			// 每个任务开始时，将等待组加1
			wg.Add(1)
			go func(id TargetID, target Target) {
				// 使用defer, 使方法完成时将等待组减1
				defer wg.Done()
				tgtRes := TargetIDResult{ID: id}
				if err := target.Save(event); err != nil {
					tgtRes.Err = err
				}
				resCh <- tgtRes
			}(id, target)
		} else {
			resCh <- TargetIDResult{ID: id}
		}
	}
	// 等待所有任务完成
	wg.Wait()
	// 向channel中写入完成标志
	done <- true
}

// NewTargetList - creates TargetList.
func NewTargetList() *TargetList {
	return &TargetList{targets: make(map[TargetID]Target)}
}
