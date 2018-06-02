package internal_tests

import (
	"testing"
	"urls_checker/internal/app/urls_checker"
)

func Test_WhenLockExists_GetOrAdd_ReturnsSameLock(t *testing.T) {
	sl := urls_checker.NewStringLocks()
	
	lock1 := sl.GetOrAdd("bla")
	lock2 := sl.GetOrAdd("bla")

	if lock1 != lock2 { t.Error("Should be same lock") }
}

func Test_GetOrAdd_ReturnsDifferentLocks_ForDifferentKeys(t *testing.T) {
	sl := urls_checker.NewStringLocks()
	
	lock1 := sl.GetOrAdd("bla1")
	lock2 := sl.GetOrAdd("bla2")

	if lock1 == lock2 { t.Error("Should be different locks") }
}

type LockSpy struct { 
	WasLocked bool
	WasUnlocked bool
}
func (ls *LockSpy) Lock()   { ls.WasLocked   = true }
func (ls *LockSpy) Unlock() { ls.WasUnlocked = true }
func (ls *LockSpy) Reset() {
	ls.WasLocked   = false
	ls.WasUnlocked = false
}

func Test_GetOrAdd_EntersLock_IfNewKey(t *testing.T) {
	spy := &LockSpy{}
	sl  := urls_checker.NewStringLocks(spy)

	sl.GetOrAdd("bla1")

	assertEnteredLock(spy, t)
}

func Test_GetOrAdd_DoesNotEnterLock_IfKeyExists(t *testing.T) {
	spy := &LockSpy{}
	sl  := urls_checker.NewStringLocks(spy)

	sl.GetOrAdd("bla1")
	spy.Reset()
	sl.GetOrAdd("bla1")

	assertDidNotEnterLock(spy, t)
}

func assertEnteredLock(spy *LockSpy, t *testing.T) {
	if !spy.WasLocked || !spy.WasUnlocked {
		t.Error("Should have entered critical section")
	}
}

func assertDidNotEnterLock(spy *LockSpy, t *testing.T) {
	if spy.WasLocked || spy.WasUnlocked {
		t.Error("Shouldn't have entered critical section")
	}
}
