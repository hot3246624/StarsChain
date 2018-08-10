package child_chain

import (
	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/providers"
	"math/big"
	"github.com/hot3246624/StarsChain/plasma_core/utils"
)

type RootEventListener struct {
	rootChain []byte
	w3 *web3.Web3
	confirmations int
	seenEvents map[string]bool
	activeEvent map[string]bool
	subscribers map[string][]func(map[string]map[string]interface{})
}

func NewRE(rootChain []byte, w3 *web3.Web3, confirmations int) *RootEventListener {
	re := &RootEventListener{rootChain:rootChain, w3: w3, confirmations: confirmations}
	if w3.Provider == nil {
		re.w3 = web3.NewWeb3(providers.NewHTTPProvider("127.0.0.1:8545", 10, false))
	}
	if confirmations == 0 {
		re.confirmations = 6
	}
	re.listenForEvent("Deposit")
	re.listenForEvent("ExitStarted")
	return re
}

func (re *RootEventListener)on(eventName string, eventHandler func(map[string]map[string]interface{})) {
	re.subscribers[eventName] = append(re.subscribers[eventName], eventHandler)
}

func (re *RootEventListener)listenForEvent(eventName string){
	re.subscribers[eventName] = nil
	re.activeEvent[eventName] = true
	go re.filterLoop(eventName)

}

func (re *RootEventListener)stopListeningForEvent(eventName string){
	delete(re.activeEvent, eventName)
}

func (re *RootEventListener)stopAll() {
	for eventName := range re.activeEvent {
		re.stopListeningForEvent(eventName)
	}
}

func (re *RootEventListener)filterLoop(eventName string) error {
	for {
		if _, err := re.activeEvent[eventName]; !err {
			break
		}
		currentBlock, err := re.w3.Eth.GetBlockByNumber(big.NewInt(-1), false)
		if err != nil {
			return err
		}
		//TODO there is no eventfilter, so wait to implement one.
	}
}

func (re *RootEventListener)broadcastEvent(eventName string, event string){
	for _, subscriber := range re.subscribers[eventName] {
		//TODO how to define the handler.
		subscriber(event)
	}
}

func (re *RootEventListener)hashEvent(event map[string]string) string{
	//TODO how to judge, this is a temp method
	return utils.RlpHash(event).String()
}

