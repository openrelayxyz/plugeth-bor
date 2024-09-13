package plugins

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
)

var (
	stateControlData       = make(map[uint64]map[string]interface{})
	coreControlData        = make(map[uint64]map[string]interface{})
	node                   core.Node
	nodeInitialized        bool
	nodeInitMutex          sync.RWMutex
	earlyStateUpdates      []stateUpdateParams
	earlyStateUpdatesMutex sync.Mutex
)

type stateUpdateParams struct {
	blockRoot     core.Hash
	parentRoot    core.Hash
	coreDestructs map[core.Hash]struct{}
	coreAccounts  map[core.Hash][]byte
	coreStorage   map[core.Hash]map[core.Hash][]byte
	coreCode      map[core.Hash][]byte
}

func InitializeNode(stack core.Node, b restricted.Backend) {
	log.Error("Initialized node test plugin")
	node = stack
	initCoreData()

	go func() {
		time.Sleep(4 * time.Second)
		nodeInitMutex.Lock()
		nodeInitialized = true
		nodeInitMutex.Unlock()

		processEarlyStateUpdates()
	}()
}

func initCoreData() {
	for i := 0; i <= 2000; i++ {
		coreControlData[uint64(i)] = make(map[string]interface{})
	}
}

func getBlockNumber() (string, error) {
	client, err := node.Attach()
	if err != nil {
		return "", err
	}
	var num string
	if err = client.Call(&num, "eth_blockNumber"); err != nil {
		return "", err
	} else {
		return num, nil
	}
}

func importChain() {
	client, err := node.Attach()
	if err != nil {
		log.Error("error attaching to stack from test plugin, BlockChain()", "err", err)
		return
	}

	for {
		var result bool
		if err = client.Call(&result, "admin_importChain", "/Users/jesseakoh/Downloads/largeBorChain.gz"); err != nil {
			log.Error("error calling admin.importChain", "err", err)
			continue
		}
		if result {
			log.Info("blockchain successfully imported chain")
			return
		}
		log.Warn("blockchain failed to import chain, retrying...")
		time.Sleep(2 * time.Second)
	}
}

func store(data map[uint64]map[string]interface{}, filename string) error {
	for k, v := range data {
		if len(v) == 0 {
			delete(data, k)
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
	nodeInitMutex.RLock()
	initialized := nodeInitialized
	nodeInitMutex.RUnlock()

	if !initialized {
		earlyStateUpdatesMutex.Lock()
		earlyStateUpdates = append(earlyStateUpdates, stateUpdateParams{
			blockRoot:     blockRoot,
			parentRoot:    parentRoot,
			coreDestructs: coreDestructs,
			coreAccounts:  coreAccounts,
			coreStorage:   coreStorage,
			coreCode:      coreCode,
		})
		earlyStateUpdatesMutex.Unlock()
		log.Info("Queued early state update")
		return
	}

	processStateUpdate(blockRoot, parentRoot, coreDestructs, coreAccounts, coreStorage, coreCode)
}

func processEarlyStateUpdates() {
	earlyStateUpdatesMutex.Lock()
	defer earlyStateUpdatesMutex.Unlock()

	log.Info("Processing early state updates", "count", len(earlyStateUpdates))
	for _, params := range earlyStateUpdates {
		processStateUpdate(params.blockRoot, params.parentRoot, params.coreDestructs, params.coreAccounts, params.coreStorage, params.coreCode)
	}
	earlyStateUpdates = nil
}

func processStateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
	n, err := getBlockNumber()
	if err != nil {
		log.Error("error returned from getBlockNumber, test plugin, StateUpdate", "err", err)
		os.Exit(1)
	}
	nbr, err := hexutil.DecodeUint64(n)
	if err != nil {
		log.Error("number decodeing error, test plugin, StateUpdate", "err", err)
		os.Exit(1)
	}

	value := make(map[string]interface{})
	stateControlData[nbr] = value

	if len(blockRoot) > 0 {
		stateControlData[nbr]["blockRoot"] = blockRoot.String()
	}

	if len(coreDestructs) > 0 {
		if nbr%5 == 0 {
			desMap := make(map[string]struct{})
			for k, v := range coreDestructs {
				desMap[k.String()] = v
			}
			stateControlData[nbr]["destructs"] = desMap
		}
	}

	if len(coreAccounts) > 0 {
		if nbr%10 == 0 {
			acctMap := make(map[string][]byte)
			for k, v := range coreAccounts {
				acctMap[k.String()] = v
			}
			stateControlData[nbr]["accounts"] = acctMap
		}
	}

	if len(coreStorage) > 0 {
		if nbr%5 == 0 {
			strMap := make(map[string]map[string][]byte)
			for k, v := range coreStorage {
				innerMap := make(map[string][]byte)
				for iK, iV := range v {
					innerMap[iK.String()] = iV
				}
				strMap[k.String()] = innerMap
			}
			stateControlData[nbr]["storages"] = strMap
		}
	}

	if len(coreCode) > 0 {
		if nbr%5 == 0 {
			cdeMap := make(map[string][]byte)
			for k, v := range coreCode {
				cdeMap[k.String()] = v
			}
			stateControlData[nbr]["code"] = cdeMap
		}
	}
}

func NewHead(blockBytes []byte, hash core.Hash, logBytes [][]byte, td *big.Int) {
	n, err := getBlockNumber()
	if err != nil {
		log.Error("error returned from getBlockNumber, test plugin, NewHead", "err", err)
		os.Exit(1)
	}
	nbr, err := hexutil.DecodeUint64(n)
	if err != nil {
		log.Error("number decodeing error, test plugin, NewHead", "err", err)
		os.Exit(1)
	}

	if len(hash) > 0 {
		coreControlData[nbr]["hash"] = hash.String()
	}

	if td != nil {
		t := *td
		coreControlData[nbr]["totalDiff"] = t.String()
	}

	if len(blockBytes) > 0 && nbr%10 == 0 {
		coreControlData[nbr]["blockBytes"] = blockBytes
	}

	if len(logBytes) > 0 && nbr%10 == 0 {
		coreControlData[nbr]["logBytes"] = logBytes
	}
}

func GatherData() {
	importChain()
	if len(stateControlData) > 0 {
		if err := store(stateControlData, "../../plugins/state-control.json"); err != nil {
			log.Error("file store error, state package", "err", err)
			os.Exit(1)
		}
		log.Info("state package data collected")
	} else {
		log.Warn("no state package data collected")
	}

	if len(coreControlData) > 0 {
		if err := store(coreControlData, "../../plugins/core-control.json"); err != nil {
			log.Error("file store error, core package", "err", err)
			os.Exit(1)
		}
		log.Info("core package data collected")
	} else {
		log.Warn("no core package data collected")
	}

	os.Exit(0)
}
