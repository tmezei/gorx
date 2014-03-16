package rx

import (
	"testing"
	"time"
)

func TestItShouldDisposeAsyncronously(t *testing.T) {
	d := NewDisposable(nil)
	if d.IsDisposed() == true {
		t.Error("Expect Disposable not to be disposed")
	}
	d.Dispose()
	if d.IsDisposed() == true {
		t.Error("Expect Disposable not to be disposed before DispositionChan triggers")
	}
	<-d.DispositionChan()
	if d.IsDisposed() == false {
		t.Error("Expect Disposable to be disposed")
	}
}

func TestItShouldCallCallback(t *testing.T) {
	callbackCalled := 0
	d := NewDisposable(func() {
		callbackCalled += 1
	})
	d.Dispose()
	<-d.DispositionChan()
	if callbackCalled != 1 {
		t.Error("Expect callback to be called on disposal")
	}
}

func TestItShouldCallMultipleCallbacks(t *testing.T) {
	calledCallbacks := []int{}
	d := NewDisposable(func() {
		calledCallbacks = append(calledCallbacks, 1)
	})
	d.AddCallback(func() {
		calledCallbacks = append(calledCallbacks, 2)
	})
	d.Dispose()
	<-d.DispositionChan()
	if len(calledCallbacks) != 2 || calledCallbacks[0] != 1 || calledCallbacks[1] != 2 {
		t.Error("Expect all callbacks to be called on disposal")
	}
}

func TestMultipleDispositionChan(t *testing.T) {
	d := NewDisposable(nil)
	receivedChanDefault := false
	receivedChan1 := false
	receivedChan3 := false
	chan1 := make(chan bool, 1)
	chan2 := make(chan bool)
	chan3 := make(chan bool)
	testTermChan := make(chan bool)
	go func() {
		for {
			select {
			case <-d.DispositionChan():
				receivedChanDefault = true
			case <-chan1:
				receivedChan1 = true
			case <-chan3:
				receivedChan3 = true
			case <-time.After(1000):
				testTermChan <- true
				return
			default:
				if receivedChanDefault && receivedChan1 && receivedChan3 {
					testTermChan <- true
					return
				}
			}
		}
	}()
	d.AddDispositionChan(chan1)
	d.AddDispositionChan(chan2)
	d.Dispose()
	d.AddDispositionChan(chan3)
	<-testTermChan
	if !d.IsDisposed() {
		t.Error("Expect disposable to dispose")
	}
	if !receivedChanDefault || !receivedChan1 || !receivedChan3 {
		t.Error("Expecting multiple disposition channels to be signaled properly")
	}
}
