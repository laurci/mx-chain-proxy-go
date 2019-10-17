package process

import (
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// HeartBeatPath represents the path where an observer exposes his heartbeat status
const HeartBeatPath = "/node/heartbeatstatus"

// HeartbeatProcessor is able to process transaction requests
type HeartbeatProcessor struct {
	proc                  Processor
	cacher                HeartbeatCacheHandler
	cacheValidityDuration time.Duration
}

// NewHeartbeatProcessor creates a new instance of TransactionProcessor
func NewHeartbeatProcessor(
	proc Processor,
	cacher HeartbeatCacheHandler,
	cacheValidityDuration time.Duration,
) (*HeartbeatProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}
	if cacher == nil || cacher.IsInterfaceNil() {
		return nil, ErrNilHeartbeatCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}
	hbp := &HeartbeatProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return hbp, nil
}

// GetHeartbeatData will simply forward the heartbeat status from an observer
func (hbp *HeartbeatProcessor) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	heartbeatsToReturn, err := hbp.cacher.Heartbeats()
	if err == nil {
		return heartbeatsToReturn, nil
	}

	log.Info(fmt.Sprintf("heartbeat: cannot get heartbeats from cache: %s. Fetching from API...", err.Error()))

	return hbp.getHeartbeatsFromApi()
}

func (hbp *HeartbeatProcessor) getHeartbeatsFromApi() (*data.HeartbeatResponse, error) {
	observers, err := hbp.proc.GetAllObservers()
	if err != nil {
		return nil, err
	}

	var heartbeatResponse data.HeartbeatResponse
	for _, observer := range observers {
		err = hbp.proc.CallGetRestEndPoint(observer.Address, HeartBeatPath, &heartbeatResponse)
		if err == nil {
			log.Info("fetched heartbeats from API")
			return &heartbeatResponse, nil
		}
		log.Error("heartbeat: Observer " + observer.Address + " didn't respond to the heartbeat request")
	}
	return nil, ErrHeartbeatNotAvailable
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (hbp *HeartbeatProcessor) StartCacheUpdate() {
	go func() {
		for {
			hbts, err := hbp.getHeartbeatsFromApi()
			if err != nil {
				log.Warn("heartbeat: error while getting heartbeats from api: " + err.Error())
			}

			if hbts != nil {
				err = hbp.cacher.StoreHeartbeats(hbts)
				if err != nil {
					log.Warn("heartbeat: can't store heartbeats in cache: " + err.Error())
				}
			}

			time.Sleep(hbp.cacheValidityDuration)
		}
	}()
}
