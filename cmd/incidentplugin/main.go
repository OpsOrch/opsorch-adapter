package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	adapter "github.com/opsorch/opsorch-adapter/incident"
	coreincident "github.com/opsorch/opsorch-core/incident"
	"github.com/opsorch/opsorch-core/schema"
)

type rpcRequest struct {
	Method  string          `json:"method"`
	Config  map[string]any  `json:"config"`
	Payload json.RawMessage `json:"payload"`
}

type rpcResponse struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var provider coreincident.Provider

func main() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	for {
		var req rpcRequest
		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			writeErr(enc, err)
			return
		}

		prov, err := ensureProvider(req.Config)
		if err != nil {
			writeErr(enc, err)
			continue
		}

		ctx := context.Background()
		switch req.Method {
		case "incident.query":
			var query schema.IncidentQuery
			if err := json.Unmarshal(req.Payload, &query); err != nil {
				writeErr(enc, err)
				continue
			}
			res, err := prov.Query(ctx, query)
			write(enc, res, err)
		case "incident.list":
			res, err := prov.List(ctx)
			write(enc, res, err)
		case "incident.get":
			var payload struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				writeErr(enc, err)
				continue
			}
			res, err := prov.Get(ctx, payload.ID)
			write(enc, res, err)
		case "incident.create":
			var in schema.CreateIncidentInput
			if err := json.Unmarshal(req.Payload, &in); err != nil {
				writeErr(enc, err)
				continue
			}
			res, err := prov.Create(ctx, in)
			write(enc, res, err)
		case "incident.update":
			var payload struct {
				ID    string                     `json:"id"`
				Input schema.UpdateIncidentInput `json:"input"`
			}
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				writeErr(enc, err)
				continue
			}
			res, err := prov.Update(ctx, payload.ID, payload.Input)
			write(enc, res, err)
		case "incident.timeline.get":
			var payload struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				writeErr(enc, err)
				continue
			}
			res, err := prov.GetTimeline(ctx, payload.ID)
			write(enc, res, err)
		case "incident.timeline.append":
			var payload struct {
				ID    string                     `json:"id"`
				Input schema.TimelineAppendInput `json:"input"`
			}
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				writeErr(enc, err)
				continue
			}
			err := prov.AppendTimeline(ctx, payload.ID, payload.Input)
			write(enc, map[string]string{"status": "ok"}, err)
		default:
			writeErr(enc, fmt.Errorf("unknown method: %s", req.Method))
		}
	}
}

func ensureProvider(cfg map[string]any) (coreincident.Provider, error) {
	if provider != nil {
		return provider, nil
	}
	prov, err := adapter.New(cfg)
	if err != nil {
		return nil, err
	}
	provider = prov
	return provider, nil
}

func write(enc *json.Encoder, result any, err error) {
	if err != nil {
		writeErr(enc, err)
		return
	}
	_ = enc.Encode(rpcResponse{Result: result})
}

func writeErr(enc *json.Encoder, err error) {
	_ = enc.Encode(rpcResponse{Error: err.Error()})
}
