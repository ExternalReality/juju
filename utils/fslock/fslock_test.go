package fslock_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	. "launchpad.net/gocheck"
	"launchpad.net/juju-core/utils/fslock"
)

func Test(t *testing.T) {
	TestingT(t)
}

type fslockSuite struct{}

var _ = Suite(fslockSuite{})

func (fslockSuite) SetUpSuite(c *C) {
	fslock.SetLockWaitDelay(1 * time.Millisecond)
}

func (fslockSuite) TearDownSuite(c *C) {
	fslock.SetLockWaitDelay(1 * time.Second)
}

// This test also happens to test that locks can get created when the parent
// lock directory doesn't exist.
func (fslockSuite) TestValidNamesLockDir(c *C) {

	for _, name := range []string{
		"a",
		"longer",
		"longer-with.special-characters",
	} {
		dir := c.MkDir()
		_, err := fslock.NewLock(dir, name)
		c.Assert(err, IsNil)
	}
}

func (fslockSuite) TestInvalidNames(c *C) {

	for _, name := range []string{
		"NoCapitals",
		"no+plus",
		"no/slash",
		"no\\backslash",
		"no$dollar",
	} {
		dir := c.MkDir()
		_, err := fslock.NewLock(dir, name)
		c.Assert(err, ErrorMatches, "Invalid lock name .*")
	}
}

func (fslockSuite) TestNewLockWithExistingDir(c *C) {
	dir := c.MkDir()
	err := os.MkdirAll(dir, 0755)
	c.Assert(err, IsNil)
	_, err = fslock.NewLock(dir, "special")
	c.Assert(err, IsNil)
}

func (fslockSuite) TestNewLockWithExistingFileInPlace(c *C) {
	dir := c.MkDir()
	err := os.MkdirAll(dir, 0755)
	c.Assert(err, IsNil)
	path := path.Join(dir, "locks")
	err = ioutil.WriteFile(path, []byte("foo"), 0644)
	c.Assert(err, IsNil)

	_, err = fslock.NewLock(path, "special")
	c.Assert(err, ErrorMatches, `.* not a directory`)
}

func (fslockSuite) TestIsLockHeldBasics(c *C) {
	dir := c.MkDir()
	lock, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	c.Assert(lock.IsLockHeld(), Equals, false)

	err = lock.Lock("")
	c.Assert(err, IsNil)
	c.Assert(lock.IsLockHeld(), Equals, true)

	err = lock.Unlock()
	c.Assert(err, IsNil)
	c.Assert(lock.IsLockHeld(), Equals, false)
}

func (fslockSuite) TestIsLockHeldTwoLocks(c *C) {
	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock1.Lock("")
	c.Assert(err, IsNil)
	c.Assert(lock2.IsLockHeld(), Equals, false)
}

func (fslockSuite) TestLockBlocks(c *C) {

	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	acquired := make(chan struct{})
	err = lock1.Lock("")
	c.Assert(err, IsNil)

	go func() {
		lock2.Lock("")
		acquired <- struct{}{}
		close(acquired)
	}()

	// Waiting for something not to happen is inherently hard...
	select {
	case <-acquired:
		c.Fatalf("Unexpected lock acquisition")
	case <-time.After(50 * time.Millisecond):
		// all good
	}

	err = lock1.Unlock()
	c.Assert(err, IsNil)

	select {
	case <-acquired:
		// all good
	case <-time.After(50 * time.Millisecond):
		c.Fatalf("Expected lock acquisition")
	}

	c.Assert(lock2.IsLockHeld(), Equals, true)
}

func (fslockSuite) TestLockWithTimeoutUnlocked(c *C) {
	dir := c.MkDir()
	lock, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock.LockWithTimeout(10*time.Millisecond, "")
	c.Assert(err, IsNil)
}

func (fslockSuite) TestLockWithTimeoutLocked(c *C) {
	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock1.Lock("")
	c.Assert(err, IsNil)

	err = lock2.LockWithTimeout(10*time.Millisecond, "")
	c.Assert(err, Equals, fslock.ErrTimeout)
}

func (fslockSuite) TestUnlock(c *C) {
	dir := c.MkDir()
	lock, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock.Unlock()
	c.Assert(err, Equals, fslock.ErrLockNotHeld)
}

func (fslockSuite) TestIsLocked(c *C) {
	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock1.Lock("")
	c.Assert(err, IsNil)

	c.Assert(lock1.IsLocked(), Equals, true)
	c.Assert(lock2.IsLocked(), Equals, true)
}

func (fslockSuite) TestBreakLock(c *C) {
	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock1.Lock("")
	c.Assert(err, IsNil)

	err = lock2.BreakLock()
	c.Assert(err, IsNil)
	c.Assert(lock2.IsLocked(), Equals, false)

	// Normally locks are broken due to client crashes, not duration.
	err = lock1.Unlock()
	c.Assert(err, Equals, fslock.ErrLockNotHeld)

	// Breaking a non-existant isn't an error
	err = lock2.BreakLock()
	c.Assert(err, IsNil)
}

func (fslockSuite) TestMessage(c *C) {
	dir := c.MkDir()
	lock, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "")

	err = lock.SetMessage("my message")
	c.Assert(err, Equals, fslock.ErrLockNotHeld)

	err = lock.Lock("")
	c.Assert(err, IsNil)

	err = lock.SetMessage("my message")
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "my message")

	// Messages can be changed while the lock is held.
	err = lock.SetMessage("new message")
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "new message")

	// Unlocking removes the message.
	err = lock.Unlock()
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "")
}

func (fslockSuite) TestMessageAcrossLocks(c *C) {
	dir := c.MkDir()
	lock1, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)
	lock2, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock1.Lock("")
	c.Assert(err, IsNil)
	err = lock1.SetMessage("very busy")
	c.Assert(err, IsNil)

	c.Assert(lock2.Message(), Equals, "very busy")
}

func (fslockSuite) TestInitialMessageWhenLocking(c *C) {
	dir := c.MkDir()
	lock, err := fslock.NewLock(dir, "testing")
	c.Assert(err, IsNil)

	err = lock.Lock("initial message")
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "initial message")

	err = lock.Unlock()
	c.Assert(err, IsNil)

	err = lock.LockWithTimeout(10*time.Millisecond, "initial timeout message")
	c.Assert(err, IsNil)
	c.Assert(lock.Message(), Equals, "initial timeout message")
}

func (fslockSuite) TestStress(c *C) {
	const lockAttempts = 100
	const concurrentLocks = 3

	var counter = new(int64)
	// Use atomics to update lockState to make sure the lock isn't held by
	// someone else. A value of 1 means locked, 0 means unlocked.
	var lockState = new(int32)
	var done = make(chan struct{})
	defer close(done)

	dir := c.MkDir()

	var stress = func(name string) {
		defer func() { done <- struct{}{} }()
		lock, err := fslock.NewLock(dir, "testing")
		if err != nil {
			return
		}
		for i := 0; i < lockAttempts; i++ {
			lock.Lock(name)
			state := atomic.AddInt32(lockState, 1)
			c.Check(state, Equals, int32(1))
			// Tell the go routine scheduler to give a slice to someone else
			// while we have this locked.
			runtime.Gosched()
			// need to decrement prior to unlock to avoid the race of someone
			// else grabbing the lock before we decrement the state.
			_ = atomic.AddInt32(lockState, -1)
			lock.Unlock()
			// increment the general counter
			_ = atomic.AddInt64(counter, 1)
		}
	}

	for i := 0; i < concurrentLocks; i++ {
		go stress(fmt.Sprintf("Lock %d", i))
	}
	for i := 0; i < concurrentLocks; i++ {
		<-done
	}
	c.Assert(*counter, Equals, int64(lockAttempts*concurrentLocks))
}
