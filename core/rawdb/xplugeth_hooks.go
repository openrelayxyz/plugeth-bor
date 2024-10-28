package rawdb

import (
	"sync"
	"reflect"
	
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/modifyancients"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	
)

var (
	freezerUpdates          map[uint64]map[string]interface{}
	lock                    sync.Mutex
)

func PluginTrackUpdate(num uint64, kind string, value interface{}) {

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
		if len(xplugeth.GetModules[modifyancients.ModifyAncientsPlugin]()) > 0 {
			var keys []string
			for k := range update {
				keys = append(keys, k)
			}
			if headeri, ok := update[ChainFreezerHeaderTable]; ok {
				var h types.Header
				switch v := headeri.(type) {
				case []byte:
					err := rlp.DecodeBytes(v, &h)
					if err != nil {
						log.Error("error decoding header, commit update", "err", err)
						continue
					}
				case *types.Header:
					log.Error("second case")
					h = *v
				default:
					log.Error("header of unknown type", "type", reflect.TypeOf(headeri))
					continue
				}
				for _, m := range xplugeth.GetModules[modifyancients.ModifyAncientsPlugin]() {
					m.ModifyAncients(num, &h)
				}
			}
		}
	}
}
