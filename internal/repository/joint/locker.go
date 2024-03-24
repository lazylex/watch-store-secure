package joint

import (
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
)

type StateLocker struct {
	waiting map[loginVO.Login][]chan bool
}

// CreateStateLocker конструктор для блокировщика
func CreateStateLocker() StateLocker {
	return StateLocker{waiting: make(map[loginVO.Login][]chan bool)}
}

// Lock блокировка чтения состояния для переданного логина
func (s *StateLocker) Lock(login loginVO.Login) {
	s.waiting[login] = make([]chan bool, 0)
}

// Unlock разблокировка чтения состояния для переданного логина
func (s *StateLocker) Unlock(login loginVO.Login) {
	for _, c := range s.waiting[login] {
		c <- true
		close(c)
	}

	delete(s.waiting, login)
}

// ReadyToRead запрашивает, возможно ли в данный момент получить значение для переданного логина. На случай, если в
// данный момент это не возможно, передается канал, в который необходимо отправить сообщение, когда эта возможность
// появится
func (s *StateLocker) ReadyToRead(login loginVO.Login, c chan bool) bool {
	if _, ok := s.waiting[login]; !ok {
		return true
	} else {
		s.waiting[login] = append(s.waiting[login], c)
	}
	return false
}
