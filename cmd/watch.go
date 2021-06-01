package cmd

import (
	"context"
	"log"

	"github.com/ryotarai/cmstore/k8s"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"gopkg.in/fsnotify.v1"
)

func init() {
	watchCmd.Flags().BoolVar(&watchFlags.allowDelete, "allow-delete", false, "Allow deleting data in ConfigMap")
	rootCmd.AddCommand(watchCmd)
}

var (
	watchCmd = &cobra.Command{
		Use:  "watch",
		RunE: runWatch,
	}
	watchFlags = struct {
		allowDelete bool
	}{}
)

func runWatch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	clientset, err := k8s.BuildClientset()
	if err != nil {
		return err
	}

	cmUpdater := &k8s.ConfigMapUpdater{
		Clientset: clientset,
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return xerrors.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	errCh := make(chan error)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					errCh <- xerrors.Errorf("fsnotify: events channel closed")
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) > 0 {
					log.Printf("Updating a ConfigMap")
					if err := cmUpdater.Update(ctx, rootFlags.namespace, rootFlags.name, rootFlags.dir, watchFlags.allowDelete); err != nil {
						log.Printf("WARN: Error on updating a ConfigMap: %s", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					errCh <- xerrors.Errorf("fsnotify: errors channel closed")
					return
				}
				log.Printf("WARN: fsnotify error: %s", err)
			}
		}
	}()

	if err := watcher.Add(rootFlags.dir); err != nil {
		return xerrors.Errorf("add dir to watcher: %w", err)
	}
	log.Printf("Watching files in %s", rootFlags.dir)
	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
