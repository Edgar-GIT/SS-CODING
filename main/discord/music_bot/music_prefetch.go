package musicbot

import (
	"context"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func trackReady(t Track) bool {
	return t.FilePath != "" || t.URL != ""
}

func trackQuery(t Track) string {
	if q := t.Query; q != "" {
		return q
	}
	return t.Title
}

func sameTrackQuery(a Track, query string) bool {
	q := trackQuery(a)
	return q != "" && q == query
}

// prefetchQueueHead downloads only the next queued song (one ahead while playing).
func (gp *GuildPlayer) prefetchQueueHead(parent context.Context) {
	gp.mu.Lock()
	if gp.prefetching || len(gp.queue) == 0 {
		gp.mu.Unlock()
		return
	}
	head := gp.queue[0]
	if trackReady(head) {
		gp.mu.Unlock()
		return
	}
	query := trackQuery(head)
	if query == "" {
		gp.mu.Unlock()
		return
	}
	gp.prefetching = true
	gp.mu.Unlock()

	defer func() {
		gp.mu.Lock()
		gp.prefetching = false
		gp.mu.Unlock()
	}()

	if sessionStopped(parent) {
		return
	}

	botLogInfo("Prefetching: %s", query)
	track, err := downloadYouTubeCtx(parent, query)
	if err != nil {
		if sessionStopped(parent) {
			return
		}
		botLogWarn("Prefetch failed for %s: %v", query, err)
		return
	}
	track.Query = query

	gp.mu.Lock()
	defer gp.mu.Unlock()
	if len(gp.queue) == 0 || !sameTrackQuery(gp.queue[0], query) {
		return
	}
	gp.queue[0] = *track
}

func (gp *GuildPlayer) maybeStartPrefetch(parent context.Context) {
	go gp.prefetchQueueHead(parent)
}

// takeNextReadyTrack waits for the head of the queue to be downloaded, then dequeues it.
func (gp *GuildPlayer) takeNextReadyTrack(
	session *discordgo.Session,
	channelID string,
	sessionCtx context.Context,
	trackCtx context.Context,
) (Track, bool) {
	const poll = 200 * time.Millisecond
	waited := time.Duration(0)
	const maxWait = 45 * time.Second

	for {
		if sessionStopped(sessionCtx) || sessionStopped(trackCtx) {
			return Track{}, false
		}

		gp.mu.Lock()
		if len(gp.queue) == 0 {
			gp.mu.Unlock()
			return Track{}, false
		}
		head := gp.queue[0]
		if trackReady(head) {
			gp.queue = gp.queue[1:]
			gp.mu.Unlock()
			return head, true
		}
		query := trackQuery(head)
		prefetching := gp.prefetching
		gp.mu.Unlock()

		if query == "" {
			gp.mu.Lock()
			gp.queue = gp.queue[1:]
			gp.mu.Unlock()
			continue
		}

		gp.maybeStartPrefetch(sessionCtx)

		if !prefetching {
			sendPlain(session, channelID, "Downloading: **"+query+"**…")
			track, err := downloadYouTubeCtx(trackCtx, query)
			if err != nil {
				if sessionStopped(sessionCtx) || sessionStopped(trackCtx) {
					return Track{}, false
				}
				sendPlain(session, channelID, "Could not download: **"+query+"**")
				gp.mu.Lock()
				if len(gp.queue) > 0 && sameTrackQuery(gp.queue[0], query) {
					gp.queue = gp.queue[1:]
				}
				gp.mu.Unlock()
				return Track{}, true
			}
			track.Query = query
			gp.mu.Lock()
			if len(gp.queue) > 0 && sameTrackQuery(gp.queue[0], query) {
				gp.queue = gp.queue[1:]
			}
			gp.mu.Unlock()
			return *track, true
		}

		time.Sleep(poll)
		waited += poll
		if waited >= maxWait {
			sendPlain(session, channelID, "Timed out waiting for: **"+query+"**")
			gp.mu.Lock()
			if len(gp.queue) > 0 && sameTrackQuery(gp.queue[0], query) {
				gp.queue = gp.queue[1:]
			}
			gp.mu.Unlock()
			return Track{}, true
		}
	}
}

var (
	panelRefreshMu      sync.Mutex
	panelRefreshPending = map[string]bool{}
)

func refreshPanel(session *discordgo.Session, guildID string) {
	panelRefreshMu.Lock()
	if panelRefreshPending[guildID] {
		panelRefreshMu.Unlock()
		return
	}
	panelRefreshPending[guildID] = true
	panelRefreshMu.Unlock()

	go func() {
		time.Sleep(400 * time.Millisecond)
		panelRefreshMu.Lock()
		delete(panelRefreshPending, guildID)
		panelRefreshMu.Unlock()
		refreshPanelNow(session, guildID)
	}()
}

func refreshPanelNow(session *discordgo.Session, guildID string) {
	gp := getPlayer(guildID)
	channelID, messageID := gp.panelRef()
	if channelID == "" || messageID == "" {
		return
	}
	embed := buildPanelEmbed(gp)
	components := panelComponents()
	_, _ = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: channelID, ID: messageID, Embed: embed, Components: &components,
	})
}
