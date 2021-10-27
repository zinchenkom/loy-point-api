package producer

import (
	"sync"
	"time"
	"log"
	"github.com/zinchenkom/loy-point-api/internal/app/sender"
	"github.com/zinchenkom/loy-point-api/internal/model"
	"github.com/zinchenkom/loy-point-api/internal/app/repo"
	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	repo   repo.EventRepo
	events <-chan loyalty.PointEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan loyalty.PointEvent,
	workerPool *workerpool.WorkerPool,
) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						log.Printf("Kafka error. Failed to send event %v", event.ID)
						p.workerPool.Submit(func() {							
							if err = p.repo.Unlock([]uint64{event.ID}); err != nil {
								log.Fatalf("Kafka error. Failed to unlock event %v", event.ID)
							}
						})
					} else {
						p.workerPool.Submit(func() {				
							if err = p.repo.Remove([]uint64{event.ID}); err != nil {
								log.Fatalf("Kafka error. Failed to remove event %v after fail in sending", event.ID)
							}
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
