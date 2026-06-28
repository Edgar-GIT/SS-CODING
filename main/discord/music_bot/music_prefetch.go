package musicbot

import "os"

func (gp *GuildPlayer) prefetchNextTrack(gen uint64) {
	if botHalted() {
		return
	}

	gp.mu.Lock()
	if gp.playGen != gen || len(gp.queue) == 0 {
		gp.mu.Unlock()
		return
	}
	next := gp.queue[0]
	if !trackNeedsDownload(next) {
		gp.mu.Unlock()
		return
	}
	query := trackSearchQuery(next)
	if gp.prefetchQuery == query {
		gp.mu.Unlock()
		return
	}
	gp.prefetchQuery = query
	gp.mu.Unlock()

	defer func() {
		gp.mu.Lock()
		if gp.prefetchQuery == query {
			gp.prefetchQuery = ""
		}
		gp.mu.Unlock()
	}()

	botLogInfo("Prefetching: %s", query)
	resolved, err := searchYouTube(query)
	if err != nil || resolved == nil || botHalted() {
		if err != nil {
			botLogWarn("Prefetch failed for %s: %v", query, err)
		}
		return
	}

	gp.mu.Lock()
	defer gp.mu.Unlock()

	if gp.playGen != gen || len(gp.queue) == 0 {
		removeOrphan(resolved.FilePath)
		return
	}
	if trackSearchQuery(gp.queue[0]) != query {
		removeOrphan(resolved.FilePath)
		return
	}

	resolved.Query = query
	gp.queue[0] = *resolved
	botLogInfo("Prefetched: %s -> %s", query, resolved.FilePath)
}

func removeOrphan(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}
