package prof

import (
	"context"
	"fmt"

	"github.com/DOSNetwork/core/log"
	"github.com/DOSNetwork/core/p2p"
)

type pprof struct {
	p      p2p.P2PInterface
	logger log.Logger
}

// NewPprof creates a pdkg struct
func NewPprof(p p2p.P2PInterface) *pprof {
	d := &pprof{
		p:      p,
		logger: log.New("module", "pprof"),
	}
	return d
}

func (d *pprof) Loop() {
	defer d.logger.Info("End Loop")
	peersToBuf, _ := d.p.SubscribeMsg(400, ProfMsg{})
	for {
		select {
		case msg, ok := <-peersToBuf:
			if !ok {
				d.logger.Info("End peersToBuf")
				return
			}
			switch content := msg.Msg.Message.(type) {
			case *ProfMsg:
				err := d.p.Reply(context.Background(), msg.Sender, msg.RequestNonce, content)
				fmt.Printf("接收到节点消息: msg.Sender=%x, content=%s\n", msg.Sender, content.Content)
				if err != nil {
					fmt.Printf("回复hello信息失败, err=%v\n", err)
				}
			}
		}
	}
}
