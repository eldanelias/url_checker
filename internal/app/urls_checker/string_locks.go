package urls_checker

import "sync"

type StringLocks struct {
	internalLock sync.Locker
	stringToLock map[string]*sync.Mutex
}

func NewStringLocks(maybeLock ...sync.Locker) StringLocks {
	lock := getLock(maybeLock)

	return StringLocks{
		internalLock: lock,
		stringToLock: make(map[string]*sync.Mutex),
	}
}

func getLock(maybeLock []sync.Locker) sync.Locker {
	if len(maybeLock) > 0 { return maybeLock[0] }
	return &sync.Mutex{}
}

func (stringLocks StringLocks) GetOrAdd(key string) *sync.Mutex {
	v, exists := stringLocks.stringToLock[key]
	if exists { return v }

	stringLocks.internalLock.Lock()
	defer stringLocks.internalLock.Unlock()

	v, exists = stringLocks.stringToLock[key]
	if exists { return v }

	stringLocks.stringToLock[key] = &sync.Mutex{}
	return stringLocks.stringToLock[key]
}
