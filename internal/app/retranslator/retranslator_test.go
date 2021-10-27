package retranslator

import (
	"testing"
	"time"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/zinchenkom/loy-point-api/internal/mocks"

	loyalty "github.com/zinchenkom/loy-point-api/internal/model"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := getStdConfig(repo, sender)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	retranslator.Close()
}


func TestLockSendAndRemoveProducer(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := getStdConfig(repo, sender)

	events := getPointEvent();

	gomock.InOrder(
		repo.EXPECT().Lock(uint64(10)).Return(events, nil).MinTimes(1).MaxTimes(1),
		sender.EXPECT().Send(&events[0]).Return(nil).MinTimes(1).MaxTimes(1),
		repo.EXPECT().Remove([]uint64{events[0].ID}).Return(nil).MinTimes(1).MaxTimes(1),
	)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second)
	retranslator.Close()
}

func TestLockSendErrorAndUnlock(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := getStdConfig(repo, sender)

	events := getPointEvent();

	gomock.InOrder(
		repo.EXPECT().Lock(uint64(10)).Return(events, nil).MinTimes(1).MaxTimes(1),
		sender.EXPECT().Send(&events[0]).Return(errors.New("sending error")).MinTimes(1).MaxTimes(1),
		repo.EXPECT().Unlock([]uint64{events[0].ID}).Return(nil).MinTimes(1).MaxTimes(1),
	)

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second)
	retranslator.Close()
}


func getStdConfig(repo *mocks.MockEventRepo, sender *mocks.MockEventSender) Config {
	return Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 500 * time.Millisecond,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}
}


func getPointEvent() []loyalty.PointEvent {
	return []loyalty.PointEvent{
		{
			ID:     1,
			Type:   loyalty.Created,
			Status: loyalty.Processed,
			Entity: &loyalty.Point{
				ID:            1,
				Name:          "PointOne",
				Description:   "DescriptionOne",
			},
		},
	}
}