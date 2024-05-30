package joint

import (
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"sync"
)

var (
	mutexW sync.Mutex
	mutexR sync.Mutex
)

type writersLocker struct {
	// Канал для блокировки/разблокировки писателей
	channel chan struct{}
	// Счетчик ожидающих писателей
	counter int
}

type StateLocker struct {
	writers map[loginVO.Login]writersLocker
	// Карта с массивом каналов, куда будут приходить сигналы о возможности чтения данных
	readers map[loginVO.Login][]chan struct{}
}

// CreateStateLocker конструктор для блокировщика записи/чтения статуса.
func CreateStateLocker() StateLocker {
	return StateLocker{
		writers: make(map[loginVO.Login]writersLocker),
		readers: make(map[loginVO.Login][]chan struct{}),
	}
}

// Lock блокирует разрешение на запись и чтение статуса для переданного логина.
func (s *StateLocker) Lock(login loginVO.Login) {
	s.lock(login) <- struct{}{}
}

// lock вызывается для получения канала, в который необходимо считать сигнал о разрешении записи. Блокирует
// разрешение на запись для других писателей.
func (s *StateLocker) lock(login loginVO.Login) chan struct{} {
	mutexW.Lock()
	defer mutexW.Unlock()

	if _, ok := s.writers[login]; !ok {
		// для писателей данных создается канал с буфером равным одному, чтобы одновременно мог писать только один
		// писатель
		c := make(chan struct{}, 1)
		s.writers[login] = writersLocker{channel: c, counter: 0}
		return c
	} else {
		s.writers[login] = writersLocker{channel: s.writers[login].channel, counter: s.writers[login].counter + 1}
		return s.writers[login].channel
	}
}

// Unlock вызывается писателем при окончании записи данных. Если писатель был последним в очереди, читателям открывается
// доступ для чтения данных.
func (s *StateLocker) Unlock(login loginVO.Login) {
	if _, ok := s.writers[login]; !ok {
		s.releaseReaders(login)
		return
	}

	<-s.writers[login].channel

	mutexW.Lock()
	defer mutexW.Unlock()

	s.writers[login] = writersLocker{channel: s.writers[login].channel, counter: s.writers[login].counter - 1}

	if s.writers[login].counter == 0 {
		close(s.writers[login].channel)
		delete(s.writers, login)
	}
}

// WantRead задерживает выполнение кода, пока идёт запись статуса для переданного логина.
func (s *StateLocker) WantRead(login loginVO.Login) {
	if c := make(chan struct{}); !s.wantRead(login, c) {
		<-c
	}
}

// wantRead возвращает true, если в данных момент не происходит запись статуса для переданного логина и чтение значения
// статуса разрешено. В противном случае заносит переданный канал в очередь на рассылку сигнала о разрешении на чтение.
func (s *StateLocker) wantRead(login loginVO.Login, c chan struct{}) bool {
	if _, ok := s.writers[login]; !ok {
		return true
	} else {
		mutexR.Lock()
		s.readers[login] = append(s.readers[login], c)
		mutexR.Unlock()
	}
	return false
}

// releaseReaders производит рассылку в каналы сигнала о разрешении на чтение данных.
func (s *StateLocker) releaseReaders(login loginVO.Login) {
	mutexR.Lock()
	mutexR.Unlock()

	for _, c := range s.readers[login] {
		c <- struct{}{}
		close(c)
	}
}
