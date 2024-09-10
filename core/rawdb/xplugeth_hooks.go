package rawdb

import (
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/openrelayxyz/xplugeth"
)

var (
	freezerUpdates          map[uint64]map[string]interface{}
	lock                    sync.Mutex
	modifyAncientsInjection *bool
	appendRawInjection      *bool
	appendInjection         *bool
)

type modifyAncientsPlugin interface {
	ModifyAncients(uint64, map[string]interface{})
}

func init() {
	xplugeth.RegisterHook[modifyAncientsPlugin]()
}

func PluginTrackUpdate(num uint64, kind string, value interface{}) {

	if appendRawInjection != nil {
		called := true
		appendRawInjection = &called
	}

	if appendInjection != nil {
		called := true
		appendInjection = &called
	}

	lock.Lock()
	defer lock.Unlock()
	if freezerUpdates == nil {
		freezerUpdates = make(map[uint64]map[string]interface{})
	}
	update, ok := freezerUpdates[num]
	if !ok {
		update = make(map[string]interface{})
		freezerUpdates[num] = update
	}
	update[kind] = value
}

func pluginCommitUpdate(num uint64) {
	if modifyAncientsInjection != nil {
		called := true
		modifyAncientsInjection = &called
	}

	lock.Lock()
	defer lock.Unlock()
	if freezerUpdates == nil {
		freezerUpdates = make(map[uint64]map[string]interface{})
	}
	min := ^uint64(0)
	for i := range freezerUpdates {
		if min > i {
			min = i
		}
	}
	for i := min; i < num; i++ {
		update, ok := freezerUpdates[i]
		defer func(i uint64) { delete(freezerUpdates, i) }(i)
		if !ok {
			log.Warn("Attempting to commit untracked block", "num", i)
			continue
		}
		for _, m := range xplugeth.GetModules[modifyAncientsPlugin]() {
			m.ModifyAncients(i, update)
		}
	}
}
