package queue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"protohackers/9_queue/queue"
	"reflect"
	"strings"
	"testing"
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

	clientNum := 100
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
}
