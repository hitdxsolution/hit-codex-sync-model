package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type packet struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string { return fmt.Sprintf("%d %s", e.Code, e.Message) }

type rpc struct {
	cmd       *exec.Cmd
	input     io.WriteCloser
	done      chan struct{}
	readDone  chan struct{}
	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[int]chan packet
	closeOnce sync.Once
	nextID    int
}

func startRPC(ctx context.Context, command string, args []string) (*rpc, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	input, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	output, err := cmd.StdoutPipe()
	if err != nil {
		input.Close()
		return nil, err
	}
	r := &rpc{cmd: cmd, input: input, pending: make(map[int]chan packet), done: make(chan struct{}), readDone: make(chan struct{})}
	if err = cmd.Start(); err != nil {
		input.Close()
		return nil, err
	}
	go func() {
		defer close(r.readDone)
		scanner := bufio.NewScanner(output)
		scanner.Buffer(make([]byte, 64*1024), 64*1024*1024)
		for scanner.Scan() {
			var p packet
			if json.Unmarshal(scanner.Bytes(), &p) != nil || len(p.ID) == 0 {
				continue
			}
			if p.Method != "" {
				if r.send(map[string]any{"id": p.ID, "error": map[string]any{
					"code": -32601, "message": "Model sync client does not execute tools",
				}}) != nil {
					return
				}
				continue
			}
			var id int
			if json.Unmarshal(p.ID, &id) != nil {
				continue
			}
			r.pendingMu.Lock()
			response := r.pending[id]
			r.pendingMu.Unlock()
			// Each request owns a buffered mailbox, so replies can arrive out of order.
			if response != nil {
				select {
				case response <- p:
				default:
				}
			}
		}
	}()
	var initialized json.RawMessage
	err = r.call(ctx, "initialize", map[string]any{
		"clientInfo":   map[string]string{"name": "hit_codex_sync_model", "version": "1.0"},
		"capabilities": map[string]bool{"experimentalApi": true},
	}, &initialized)
	if err == nil {
		err = r.send(map[string]any{"method": "initialized"})
	}
	if err != nil {
		r.close()
		return nil, err
	}
	return r, nil
}

func (r *rpc) send(value any) error {
	r.writeMu.Lock()
	defer r.writeMu.Unlock()
	return json.NewEncoder(r.input).Encode(value)
}

func (r *rpc) call(ctx context.Context, method string, params any, result any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.pendingMu.Lock()
	r.nextID++
	id := r.nextID
	response := make(chan packet, 1)
	r.pending[id] = response
	r.pendingMu.Unlock()
	defer func() { r.pendingMu.Lock(); delete(r.pending, id); r.pendingMu.Unlock() }()
	if err := r.send(map[string]any{"id": id, "method": method, "params": params}); err != nil {
		return err
	}
	timer := time.NewTimer(60 * time.Second)
	defer timer.Stop()
	var p packet
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.done:
		return fmt.Errorf("%s: App Server 연결이 종료됐습니다", method)
	case <-timer.C:
		return fmt.Errorf("%s 응답 시간 초과", method)
	case <-r.readDone:
		select {
		case p = <-response:
		default:
			return fmt.Errorf("%s: App Server 연결이 종료됐습니다", method)
		}
	case p = <-response:
	}
	if p.Error != nil {
		return fmt.Errorf("%s: %w", method, p.Error)
	}
	if result == nil {
		return nil
	}
	return json.Unmarshal(p.Result, result)
}

func pages[T any](ctx context.Context, r *rpc, method string, params map[string]any) ([]T, error) {
	var rows []T
	seen := map[string]bool{}
	for {
		var page struct {
			Data []T    `json:"data"`
			Next string `json:"nextCursor"`
		}
		if err := r.call(ctx, method, params, &page); err != nil {
			return nil, err
		}
		rows = append(rows, page.Data...)
		if page.Next == "" {
			return rows, nil
		}
		if seen[page.Next] {
			return nil, fmt.Errorf("%s: 페이지 커서가 반복됩니다", method)
		}
		seen[page.Next] = true
		params["cursor"] = page.Next
	}
}

func (r *rpc) close() {
	r.closeOnce.Do(func() {
		close(r.done)
		r.input.Close()
		finished := make(chan struct{})
		go func() { r.cmd.Wait(); close(finished) }()
		select {
		case <-finished:
		case <-time.After(2 * time.Second):
			r.cmd.Process.Kill()
			<-finished
		}
	})
}
