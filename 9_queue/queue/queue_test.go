package queue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"protohackers/9_queue/queue"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	in := "{\"request\":\"get\",\"queues\":[\"queue1\",\"queue2\"],\"wait\":true}"
	exp := queue.GetRequest{
		Request: "get",
		Queues:  []string{"queue1", "queue2"},
		Wait:    true,
	}

	var out queue.GetRequest

	b := bytes.NewBufferString(in)
	dec := json.NewDecoder(b)

	err := dec.Decode(&out)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exp, out) {
		t.Fatalf("decode fail exp %v got %v", exp, out)
	}
}

func intPtr(v int) *int {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func rawMsgPtr(v json.RawMessage) *json.RawMessage {
	return &v
}

func TestEncode(t *testing.T) {
	type encodeCases struct {
		in  any
		exp string
	}

	runTest := func(t *testing.T, cs []encodeCases) {
		for _, c := range cs {
			b := new(bytes.Buffer)
			enc := json.NewEncoder(b)
			err := enc.Encode(c.in)
			if err != nil {
				t.Fatal(err)
			}
			out := b.String()
			out = strings.TrimSpace(out)
			if out != c.exp {
				t.Fatalf("decode fail exp %v got %v", c.exp, out)
			}
		}
	}

	t.Run("get cases", func(t *testing.T) {
		cases := []encodeCases{
			{
				in:  queue.GetResponse{Status: queue.StatusNoJob},
				exp: "{\"status\":\"no-job\"}",
			},
			{
				in: queue.GetResponse{
					Status: queue.StatusOK,
					Id:     intPtr(12345),
					Pri:    intPtr(123),
					Queue:  strPtr("queue1"),
					Job:    rawMsgPtr([]byte("{\"f\":1}")),
				},
				exp: "{\"status\":\"ok\",\"id\":12345,\"job\":{\"f\":1},\"pri\":123,\"queue\":\"queue1\"}",
			},
		}
		runTest(t, cases)
	})
}

func TestJobCenter(t *testing.T) {
	ctx := context.Background()
	jc := queue.NewJobCenter(ctx)

	clientNum := jc.GetClientID()
	inPri := 123
	inQueue := "queue1"
	inJob := json.RawMessage([]byte("{\"f\":1}"))

	var jobId int

	t.Run("put test", func(t *testing.T) {
		req := &queue.PutRequest{
			Request: queue.RequestPut,
			Queue:   inQueue,
			Pri:     inPri,
			Job:     inJob,
		}
		resp := jc.Put(req)

		if resp.Status != queue.StatusOK {
			t.Fatalf("failed to put")
		}

		jobId = resp.Id
	})

	t.Run("get nonexistent test", func(t *testing.T) {
		req := &queue.GetRequest{
			Request:  queue.RequestGet,
			Queues:   []string{inQueue + "abc"},
			ClientID: clientNum,
		}
		resp := jc.Get(req)

		if resp.Status != queue.StatusNoJob {
			t.Fatalf("got nonexistent job")
		}
	})

	t.Run("get test", func(t *testing.T) {
		req := &queue.GetRequest{
			Request:  queue.RequestGet,
			Queues:   []string{inQueue},
			ClientID: clientNum,
		}
		resp := jc.Get(req)

		if resp.Status != queue.StatusOK {
			t.Fatalf("failed to get")
		}

		outId := *resp.Id
		outJob := *resp.Job
		outPri := *resp.Pri
		outQueue := *resp.Queue

		if outId != jobId {
			t.Fatalf("id mismatch with put exp %v got %v", jobId, outId)
		}
		if !bytes.Equal(outJob, inJob) {
			t.Fatalf("wrong job exp %v got %v", inJob, outJob)
		}
		if outPri != inPri {
			t.Fatalf("wrong prio exp %v got %v", outPri, inPri)
		}
		if outQueue != inQueue {
			t.Fatalf("wrong queue exp %v got %v", outQueue, inQueue)
		}
	})

	t.Run("abort by someone else", func(t *testing.T) {
		req := &queue.AbortRequest{
			Request:  queue.RequestAbort,
			Id:       jobId,
			ClientID: clientNum + 1,
		}
		resp := jc.Abort(req)

		if resp.Status != queue.StatusError {
			t.Fatalf("job aborted by someone else")
		}
	})

	t.Run("abort test", func(t *testing.T) {
		req := &queue.AbortRequest{
			Request:  queue.RequestAbort,
			Id:       jobId,
			ClientID: clientNum,
		}
		resp := jc.Abort(req)

		if resp.Status != queue.StatusOK {
			t.Fatalf("failed to abort")
		}
	})

	t.Run("abort not found test", func(t *testing.T) {
		req := &queue.AbortRequest{
			Request: queue.RequestAbort,
			Id:      jobId + 1,
		}
		resp := jc.Abort(req)

		if resp.Status != queue.StatusNoJob {
			t.Fatalf("aborted nonexistent request")
		}
	})

	t.Run("delete test", func(t *testing.T) {
		req := &queue.DeleteRequest{
			Request: queue.RequestDelete,
			Id:      jobId,
		}
		resp := jc.Delete(req)

		if resp.Status != queue.StatusOK {
			t.Fatalf("failed to delete")
		}
	})

	t.Run("disconnect", func(t *testing.T) {
		req := &queue.DisconnectRequest{
			ClientID: clientNum,
		}
		jc.DisconnectWorker(req)
	})
}

func TestJsonHandling(t *testing.T) {
	in := "{\"request\":\"put\",\"queue\":\"queue1\",\"job\":{\"title\":\"example-job\"},\"pri\":123}\n"

	val, err := queue.Decode(json.RawMessage(in))
	if err != nil {
		t.Fatal(err)
	}

	valCast, ok := val.(*queue.PutRequest)
	if !ok {
		t.Fatalf("wrong type %v %T", valCast, valCast)
	}
}

func TestJsonHandlingFail(t *testing.T) {
	in := "{\"request\":\"put\",\"queue\":\"queue1\",\"job\":{title:\"example-job\"},\"pri\":123}\n"
	val, err := queue.Decode(json.RawMessage(in))
	if err == nil {
		t.Fatalf("supposed to error %v", val)
	}
}

func TestWait(t *testing.T) {
	ctx := context.Background()
	jc := queue.NewJobCenter(ctx)

	inQueue := "test"
	inClient := 100
	inJob := json.RawMessage([]byte{})
	inPri := 1000

	putReq := queue.PutRequest{
		Request: queue.RequestPut,
		Queue:   inQueue,
		Pri:     inPri,
		Job:     inJob,
	}

	getReq := &queue.GetRequest{
		Request:  queue.RequestGet,
		Queues:   []string{inQueue},
		Wait:     true,
		ClientID: inClient,
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		jc.Put(&putReq)
	}()

	resp := jc.Get(getReq)

	outJob := *resp.Job
	outPri := *resp.Pri
	outQueue := *resp.Queue

	if !bytes.Equal(outJob, inJob) {
		t.Fatalf("wrong job exp %v got %v", inJob, outJob)
	}
	if outPri != inPri {
		t.Fatalf("wrong prio exp %v got %v", outPri, inPri)
	}
	if outQueue != inQueue {
		t.Fatalf("wrong queue exp %v got %v", outQueue, inQueue)
	}
}

func TestPerformance(t *testing.T) {
	ctx := context.Background()
	jc := queue.NewJobCenter(ctx)
	done := make(chan struct{})

	inQueue := "test"
	inJob := json.RawMessage([]byte{})
	reqs := 100000
	clientId := jc.GetClientID()

	go func() {
		defer close(done)

		for i := 0; i < reqs; i++ {
			putReq := queue.PutRequest{
				Request: queue.RequestPut,
				Queue:   inQueue,
				Pri:     i,
				Job:     inJob,
			}
			jc.Put(&putReq)
		}

		for i := 0; i < reqs; i++ {
			getReq := queue.GetRequest{
				Request:  queue.RequestGet,
				Queues:   []string{inQueue},
				ClientID: clientId,
			}
			res := jc.Get(&getReq)
			out_pri := *res.Pri
			exp_pri := (100000 - 1) - i
			if out_pri != exp_pri {
				err := fmt.Errorf("wrong prio exp %v got %v", exp_pri, out_pri)
				panic(err)
			}
		}
	}()

	select {
	case <-time.After(5 * time.Second):
		t.Fatalf("expected %v requests to finish under 5 secs.", reqs*2)
	case <-done:
	}

}
