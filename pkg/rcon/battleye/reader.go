package battleye

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	be_proto "github.com/playnet-public/battleye/battleye"
	"github.com/playnet-public/gorcon/pkg/rcon"
	context "github.com/seibert-media/golibs/log"
)

// HandlePacket received from UDP connection
func (c *Connection) HandlePacket(ctx context.Context, p be_proto.Packet) (err error) {
	defer func() {
		err = errors.Wrap(err, "handling packet")
	}()
	err = c.Protocol.Verify(p)
	if err != nil {
		ctx.Error("handling packet", zap.Error(err))
		return err
	}
	data, err := c.Protocol.Data(p)
	if err != nil {
		ctx.Error("handling packet", zap.Error(err))
		return err
	}

	// Handle KeepAlive Pingback
	if len(data) < 1 {
		c.AddPingback()
		ctx.Debug("pingback", zap.Int64("count", c.Pingback()))
		return nil
	}

	t, err := c.Protocol.Type(p)
	if err != nil {
		ctx.Error("handling packet", zap.Error(err))
		return err
	}

	switch t {
	case be_proto.Command | be_proto.MultiCommand:
		return c.HandleResponse(ctx, p)

	case be_proto.ServerMessage:
		return c.HandleServerMessage(ctx, p)
	}

	return nil
}

// HandleResponse by retrieving the corresponding transmission and updating it
func (c *Connection) HandleResponse(ctx context.Context, p be_proto.Packet) error {
	s, err := c.Protocol.Sequence(p)
	if err != nil {
		return errors.Wrap(err, "handling response")
	}

	trm := c.GetTransmission(s)
	if trm == nil {
		return errors.New("no transmission for response")
	}

	t, err := c.Protocol.Type(p)
	if err != nil {
		return errors.Wrap(err, "handling response")
	}

	data, err := c.Protocol.Data(p)
	if err != nil {
		return errors.Wrap(err, "handling response")
	}

	last := true
	if t == be_proto.MultiCommand {
		count, index, single := c.Protocol.Multi(p)
		if !single {
			trm.multiBuffer[int(index)] = data
			last = (index+1 >= count)
		}
	} else {
		trm.multiBuffer[0] = data
	}

	if last {
		select {
		case trm.done <- true:
			return nil
		case <-time.After(time.Second):
			ctx.Debug("timeout on done transmission", zap.Uint32("seq", trm.Key()), zap.String("request", trm.Request()))
			return nil
		}
	}

	return nil
}

// HandleServerMessage containing chat and events
func (c *Connection) HandleServerMessage(ctx context.Context, p be_proto.Packet) error {
	s, err := c.Protocol.Sequence(p)
	if err != nil {
		return errors.Wrap(err, "handling server message")
	}

	var channels = []string{
		"(Group)",
		"(Vehicle)",
		"(Unknown)",
	}
	var t = rcon.TypeEvent
	for _, c := range channels {
		if strings.HasPrefix(string(p), c) {
			t = rcon.TypeChat
		}
	}

	event := &rcon.Event{
		Timestamp: time.Now(),
		Type:      t,
		Payload:   string(p),
	}

	_, err = c.UDP.Write(c.Protocol.BuildMsgAckPacket(s))
	if err != nil {
		return errors.Wrap(err, "handling server message")
	}

	c.subscriptionsMutex.RLock()
	defer c.subscriptionsMutex.RUnlock()
	for _, l := range c.subscriptions {
		go func(l chan *rcon.Event) { l <- event }(l)
	}

	return nil
}
