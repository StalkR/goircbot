// Package followers implements a plugin which invites specified people to join channels the bot joins.
// It assumes that servers reply to a succesful join with RPL_CHANNELMODEIS (324).
// It's not in RFC but it's pretty popular behaviour.
// Note: it WILL invite the followers if RPL_CHANNELMODEIS is sent because of other event.
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
	g.remover = b.HandleFunc(rplChannelModeIs, func(_ *client.Conn, line *client.Line) {
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

const rplChannelModeIs = "324"

// handleJoin invites followers to channel joined by the bot.
func (g *Guru) handleJoin(b bot.Bot, line *client.Line) {
	g.RLock()
	defer g.RUnlock()
	channel := line.Args[1]
	for follower := range g.followers {
		b.Invite(follower, channel)
	}
}
