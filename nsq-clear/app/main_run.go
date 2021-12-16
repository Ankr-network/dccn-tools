package app

import (
	"github.com/Ankr-network/dccn-tools/nsq-clear/pkg/nsqadmin"
	"github.com/Ankr-network/dccn-tools/nsq-clear/pkg/slog"
	"github.com/Ankr-network/dccn-tools/nsq-clear/share"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"time"
)

func MainServe(c *cli.Context) error {
	slog.SetLogLevel(c.String(share.LogLevel))
	var (
		nsq       = nsqadmin.NewNsqAdmin(c.String(share.NsqAdmin))
		tc        = time.NewTicker(c.Duration(share.Schedule))
		threshold = c.Uint64(share.Threshold)
		topics    []string
		err       error
		depth     uint64
	)

	for range tc.C {
		topics, err = nsq.GetTopic()
		if err != nil {
			log.Error().Err(err).Msg("failed to get topics")
			continue
		}
		log.Debug().Strs("topics", topics).Msg("get topics")
		for _, topic := range topics {
			depth, err = nsq.GetTopicDepth(topic)
			if err != nil {
				log.Error().Err(err).Str("topic", topic).Msg("failed to get topic depth")
				continue
			}
			log.Debug().Uint64("depth", depth).Str("topic", topic).Msg("get depth")
			if depth > threshold {
				err = nsq.EmptyQueue(topic)
				if err != nil {
					log.Error().Err(err).Str("topic", topic).Msg("failde empty queue")
				} else {
					log.Info().Uint64("depth", depth).Str("topic", topic).Msg("empty queue")
				}
			}
		}
	}
	return nil
}
