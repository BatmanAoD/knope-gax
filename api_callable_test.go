package gax

import (
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	testCallSettings = []CallOption{
		WithRetryCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}),
		// initial, max, multiplier
		WithDelayTimeoutSettings(100*time.Millisecond, 300*time.Millisecond, 1.5),
		WithRPCTimeoutSettings(50*time.Millisecond, 500*time.Millisecond, 3.0),
		WithTotalRetryTimeout(1000 * time.Millisecond),
	}
)

func TestInvokeWithTimeout(t *testing.T) {
	ctx := context.Background()
	var ok bool
	Invoke(ctx, func(childCtx context.Context) error {
		_, ok = childCtx.Deadline()
		return nil
	}, WithTimeout(10000*time.Millisecond))
	if !ok {
		t.Errorf("expected call to have an assigned timeout")
	}
}

func TestInvokeWithOKResponseWithTimeout(t *testing.T) {
	ctx := context.Background()
	var resp int
	err := Invoke(ctx, func(childCtx context.Context) error {
		resp = 42
		return nil
	}, WithTimeout(10000*time.Millisecond))
	if resp != 42 || err != nil {
		t.Errorf("expected call to return nil and set resp to 42")
	}
}

func TestInvokeWithDeadlineAfterRetries(t *testing.T) {
	ctx := context.Background()
	count := 0

	now := time.Now()
	expectedTimeout := []time.Duration{
		0,
		150 * time.Millisecond,
		450 * time.Millisecond,
	}

	err := Invoke(ctx, func(childCtx context.Context) error {
		t.Log("delta:", time.Now().Sub(now.Add(expectedTimeout[count])))
		if !time.Now().After(now.Add(expectedTimeout[count])) {
			t.Errorf("expected %s to pass before this call", expectedTimeout[count])
		}
		count += 1
		<-childCtx.Done()
		return grpc.Errorf(codes.DeadlineExceeded, "")
	}, testCallSettings...)
	if count != 3 || err == nil {
		t.Errorf("expected call to retry 3 times and return an error")
	}
}

func TestInvokeWithOKResponseAfterRetries(t *testing.T) {
	ctx := context.Background()
	count := 0

	var resp int
	err := Invoke(ctx, func(childCtx context.Context) error {
		count += 1
		if count == 3 {
			resp = 42
			return nil
		}
		<-childCtx.Done()
		return grpc.Errorf(codes.DeadlineExceeded, "")
	}, testCallSettings...)
	if count != 3 || resp != 42 || err != nil {
		t.Errorf("expected call to retry 3 times, return nil, and set resp to 42")
	}
}