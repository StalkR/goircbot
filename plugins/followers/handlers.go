// Package followers implements a plugin which invites specified people to join channels the bot joins.
// Note: it assumes that RPL_TOPIC is a response to joining a channel; if the bot is
// manually asking for a topic it will send the invites anyway.
package followers

import (
	"sync"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

// Register creates a Guru (plugin's instance) and registers it with the bot.
// Use Close to stop the instance.
func Register(b bot.Bot) *Guru {
	g := &Guru{
		followers: make(map[string]struct{}),
	}
	g.remover = b.HandleFunc(replyTopic, func(_ *client.Conn, line *client.Line) {
		g.handleJoin(b, line)
	})
	return g
}

// Guru is a structure which holds list of bot's followers.
// Use Register to create one.
type Guru struct {
	sync.RWMutex

	followers map[string]struct{}
	remover   client.Remover
}

// Add causes the bot to send invites to the follower.
func (g *Guru) Add(follower string) {
	g.Lock()
	defer g.Unlock()
	g.followers[follower] = struct{}{}
}

// Del stops the bot from sending invites to the follower.
func (g *Guru) Del(follower string) {
	g.Lock()
	defer g.Unlock()
	delete(g.followers, follower)
}

// Close stops the Guru from inviting the followers.
func (g *Guru) Close() {
	g.remover.Remove()
}

const replyTopic = "332"

// handleJoin invites followers to channel joined by the bot.
func (g *Guru) handleJoin(b bot.Bot, line *client.Line) {
	g.RLock()
	defer g.RUnlock()
	channel := line.Args[0]
	for follower := range g.followers {
		b.Invite(follower, channel)
	}
}
