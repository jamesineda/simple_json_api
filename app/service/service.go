package service

import (
	"fmt"
	"github.com/jamesineda/simple_json_api/app/models"
	"log"
	"time"
)

type Service struct {
	PhotoProcessChannel chan<- models.Photos
	Logger              ContextLogger
}

func NewService(pc chan<- models.Photos) *Service {
	return &Service{PhotoProcessChannel: pc, Logger: NewContextLogger()}
}

// StartPhotoProcessorRoutine Just an example of asynchronous processing. We could validate the payment_url here, hit a car registration endpoint, etc.
// In practice, something that would require processing in the background would be a separate service all togather,
// processing a redis queue or something.
func (s *Service) StartPhotoProcessorRoutine(photos <-chan models.Photos, shutdown <-chan bool) {
	go func(photos <-chan models.Photos, shutdown <-chan bool) {
		for {
			select {
			case p := <-photos:

				p, err := p.GetFileTypes()
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] error processing photos %s", err))
				}
				log.Printf("PHOTO IMAGE TYPES %s", p)

			case <-shutdown:
				log.Println("shutting down photo processor!")
				return
			}

			// just so the CPU isn't at 100%
			time.Sleep(1 * time.Second)
		}

	}(photos, shutdown)
}
