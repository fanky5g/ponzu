package analytics

import (
	"github.com/fanky5g/ponzu/internal/constants"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// RANGE determines the number of days ponzu request analytics and metrics are
// stored and displayed within the system
const RANGE = 14

var (
	requestChan chan HTTPRequestMetadata
)

// Record queues an apiRequest for metrics
func (s *Service) Record(req HTTPRequestMetadata) {
	// put r on buffered requestChan to take advantage of batch insertion in DB
	requestChan <- req
}

func (s *Service) StartRecorder() {
	requestChan = make(chan HTTPRequestMetadata, 1024*64*runtime.NumCPU())
	// make timer to notify select to batch request insert from requestChan
	// interval: 30 seconds
	apiRequestTimer := time.NewTicker(time.Second * 30)

	// make timer to notify select to remove analytics older than RANGE days
	// interval: RANGE/2 days
	pruneThreshold := time.Hour * 24 * RANGE
	pruneDBTimer := time.NewTicker(pruneThreshold / 2)

	for {
		select {
		case <-apiRequestTimer.C:
			var reqs []HTTPRequestMetadata
			batchSize := len(requestChan)

			for i := 0; i < batchSize; i++ {
				reqs = append(reqs, <-requestChan)
			}

			_, err := s.requestsRepository.Insert(reqs)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to insert analytics entry")
			}

		case <-pruneDBTimer.C:
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			maxTimeToUpdate := today.Add(pruneThreshold)

			err := s.requestsRepository.DeleteBy("timestamp", constants.LessThanOrEqualTo, maxTimeToUpdate)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to delete old analytics")
			}

		case <-time.After(time.Second * 30):
			continue
		}
	}
}
