package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/ryotarai/cmstore/k8s"
	"github.com/spf13/cobra"
	"gopkg.in/fsnotify.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:  "watch",
	RunE: runWatch,
}

func runWatch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	clientset, err := k8s.BuildClientset()
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	errCh := make(chan error)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					errCh <- fmt.Errorf("fsnotify: events channel closed")
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) > 0 {
					log.Printf("Updating a ConfigMap")
					err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
						cm, err := clientset.CoreV1().ConfigMaps(rootFlags.namespace).Get(ctx, rootFlags.name, metav1.GetOptions{})
						if err != nil {
							return err
						}
						cm.Data = map[string]string{}
						cm.BinaryData = map[string][]byte{}

						entries, err := os.ReadDir(rootFlags.dir)
						if err != nil {
							return err
						}
						for _, e := range entries {
							if e.IsDir() {
								continue
							}
							b, err := os.ReadFile(filepath.Join(rootFlags.dir, e.Name()))
							if err != nil {
								return err
							}
							if utf8.Valid(b) {
								cm.Data[e.Name()] = string(b)
							} else {
								cm.BinaryData[e.Name()] = b
							}
						}

						_, err = clientset.CoreV1().ConfigMaps(rootFlags.namespace).Update(ctx, cm, metav1.UpdateOptions{})
						if err != nil {
							return err
						}

						return nil
					})
					if err != nil {
						log.Printf("WARN: Error on updating a ConfigMap: %s", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					errCh <- fmt.Errorf("fsnotify: errors channel closed")
					return
				}
				log.Printf("WARN: fsnotify error: %s", err)
			}
		}
	}()

	if err := watcher.Add(rootFlags.dir); err != nil {
		return err
	}
	log.Printf("Watching files in %s", rootFlags.dir)
	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
