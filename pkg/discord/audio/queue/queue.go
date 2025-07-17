package queue

import (
	"math/rand"
	"sync"
	"time"
)

type Queue struct {
	mutex  sync.RWMutex
	tracks []*Track
}

func NewQueue() *Queue {
	return &Queue{
		tracks: make([]*Track, 0),
	}
}

func (q *Queue) Add(track *Track) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.tracks = append(q.tracks, track)
}

func (q *Queue) Pop() *Track {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.tracks) == 0 {
		return nil
	}

	track := q.tracks[0]
	q.tracks = q.tracks[1:]

	return track
}

func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.tracks = q.tracks[:0]
}

func (q *Queue) Shuffle() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.tracks) < 2 {
		return
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(q.tracks)-1, func(i, j int) {
		q.tracks[i+1], q.tracks[j+1] = q.tracks[j+1], q.tracks[i+1]
	})
}

func (q *Queue) Peek() *Track {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.tracks) == 0 {
		return nil
	}

	track := q.tracks[0]

	return track
}

func (q *Queue) Tracks() []*Track {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	tracks := make([]*Track, len(q.tracks))
	copy(tracks, q.tracks)

	return tracks
}

func (q *Queue) Len() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.tracks)
}

func (q *Queue) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.tracks) == 0
}
