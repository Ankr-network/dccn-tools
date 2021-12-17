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
		topicInfo *nsqadmin.TopicInfo
	)

	for range tc.C {
		topics, err = nsq.GetTopic()
		if err != nil {
			log.Error().Err(err).Msg("failed to get topics")
			continue
		}

		log.Debug().Strs("topics", topics).Msg("get topics")
		for _, topic := range topics {
			topicInfo, err = nsq.GetTopicDepth(topic)
			if err != nil {
				log.Error().Err(err).Str("topic", topic).Msg("failed to get topic depth")
				continue
			}

			for _, v := range topicInfo.Channels {
				log.Debug().Str("topic", topic).Str("channel", v.ChannelName).Uint64("depth", v.Depth).Msg("get depth for chan")
				if v.Depth > threshold {
					err = nsq.EmptyQueue(topic, v.ChannelName)
					if err != nil {
						log.Error().Err(err).Str("topic", topic).Str("channel", v.ChannelName).Msg("failed to empty queue for chan")
					} else {
						log.Info().Str("topic", topic).Str("channel", v.ChannelName).Uint64("depth", v.Depth).Msg("empty queue for chan")
					}
				}
			}

			log.Debug().Str("topic", topic).Uint64("depth", topicInfo.Depth).Msg("get depth for topic")
			if topicInfo.Depth > threshold {
				err = nsq.EmptyQueue(topic, "")
				if err != nil {
					log.Error().Err(err).Str("topic", topic).Msg("failed to empty queue for topic")
				} else {
					log.Info().Str("topic", topic).Uint64("depth", topicInfo.Depth).Msg("empty queue for topic")
				}
			}
		}
	}
	return nil
}
