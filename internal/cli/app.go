package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/example/teleport-cluster-digital-twin/internal/analysis"
	"github.com/example/teleport-cluster-digital-twin/internal/collectors"
	"github.com/example/teleport-cluster-digital-twin/internal/diff"
	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/storage"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// Runner implements CLI commands for collect, show, explain, analyze and diff.
type Runner struct {
	Factory      teleport.ClientFactory
	Orchestrator *collectors.Orchestrator
	Store        storage.SnapshotStore
	Out          io.Writer
}

// Run executes CLI commands.
func (r *Runner) Run(ctx context.Context, args []string) error {
	if r.Out == nil {
		return fmt.Errorf("output writer is nil")
	}
	if len(args) == 0 {
		_, _ = fmt.Fprintln(r.Out, "usage: twin <collect|show|explain|analyze|diff>")
		return nil
	}
	switch args[0] {
	case "collect":
		return r.collect(ctx)
	case "show":
		return r.show(ctx, args[1:])
	case "explain":
		return r.explain(ctx, args[1:])
	case "analyze":
		return r.analyze(ctx, args[1:])
	case "diff":
		if len(args) < 3 {
			return fmt.Errorf("usage: twin diff <old-id> <new-id>")
		}
		return r.runDiff(ctx, args[1], args[2])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (r *Runner) collect(ctx context.Context) error {
	client, auth, err := r.Factory.New(ctx)
	if err != nil {
		return err
	}
	s := r.Orchestrator.Collect(ctx, client, auth)
	if err := r.Store.Save(ctx, s); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(r.Out, s.SnapshotID)
	return nil
}

func (r *Runner) show(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: twin show <summary|users|roles|resources>")
	}
	s, err := r.Store.LoadCurrent(ctx)
	if err != nil {
		return err
	}
	switch args[0] {
	case "summary":
		return json.NewEncoder(r.Out).Encode(analysis.New(graph.FromSnapshot(s)).Topology())
	case "users":
		return printEntities(r.Out, s.Entities, model.EntityUser)
	case "roles":
		return printEntities(r.Out, s.Entities, model.EntityRole)
	case "resources":
		for _, t := range []model.EntityType{model.EntityNode, model.EntityDatabase, model.EntityKubernetesCluster, model.EntityWindowsDesktop, model.EntityApplication} {
			if err := printEntities(r.Out, s.Entities, t); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown show subcommand %q", args[0])
	}
}

func (r *Runner) explain(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: twin explain <user|resource> <value>")
	}
	s, err := r.Store.LoadCurrent(ctx)
	if err != nil {
		return err
	}
	a := analysis.New(graph.FromSnapshot(s))
	switch args[0] {
	case "user":
		out, err := a.ExplainUserAccess(args[1])
		if err != nil {
			return err
		}
		return json.NewEncoder(r.Out).Encode(out)
	case "resource":
		out, err := a.ExplainResourceExposure(args[1])
		if err != nil {
			return err
		}
		return json.NewEncoder(r.Out).Encode(out)
	default:
		return fmt.Errorf("unknown explain subcommand %q", args[0])
	}
}

func (r *Runner) analyze(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] != "risks" {
		return fmt.Errorf("usage: twin analyze risks")
	}
	s, err := r.Store.LoadCurrent(ctx)
	if err != nil {
		return err
	}
	return json.NewEncoder(r.Out).Encode(analysis.New(graph.FromSnapshot(s)).AnalyzeRisks())
}

func (r *Runner) runDiff(ctx context.Context, oldID, newID string) error {
	oldSnap, err := r.Store.Load(ctx, oldID)
	if err != nil {
		return err
	}
	newSnap, err := r.Store.Load(ctx, newID)
	if err != nil {
		return err
	}
	return json.NewEncoder(r.Out).Encode(diff.Diff(oldSnap, newSnap))
}

func printEntities(out io.Writer, entities []model.Entity, typ model.EntityType) error {
	for _, e := range entities {
		if e.Type == typ {
			if _, err := fmt.Fprintln(out, strings.TrimSpace(string(e.Type))+":", e.ID, e.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
