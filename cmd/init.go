package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/ryotarai/cmstore/k8s"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	initCmd = &cobra.Command{
		Use:  "init",
		RunE: runInit,
	}
	initFlags = struct {
		createIfNotFound bool
	}{}
)

func init() {
	initCmd.Flags().BoolVarP(&initFlags.createIfNotFound, "create-if-not-found", "c", false, "")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	clientset, err := k8s.BuildClientset()
	if err != nil {
		return err
	}

	log.Printf("Creating a directory %s", rootFlags.dir)
	if err := os.MkdirAll(rootFlags.dir, 0777); err != nil {
		return xerrors.Errorf("mkdir: %w", err)
	}

	log.Printf("Getting a ConfigMap %s/%s", rootFlags.namespace, rootFlags.name)
	cm, err := clientset.CoreV1().ConfigMaps(rootFlags.namespace).Get(ctx, rootFlags.name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) && initFlags.createIfNotFound {
			log.Printf("Creating a ConfigMap %s/%s", rootFlags.namespace, rootFlags.name)
			_, err := clientset.CoreV1().ConfigMaps(rootFlags.namespace).Create(ctx, &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootFlags.namespace,
					Name:      rootFlags.name,
				},
			}, metav1.CreateOptions{})
			if err != nil {
				return xerrors.Errorf("create configmap: %w", err)
			}
			return nil
		}
		return xerrors.Errorf("get configmap: %w", err)
	}

	for k, v := range cm.Data {
		p := filepath.Join(rootFlags.dir, k)
		log.Printf("Writing %s", p)
		if err := os.WriteFile(p, []byte(v), 0666); err != nil {
			return xerrors.Errorf("write file: %w", err)
		}
	}
	for k, v := range cm.BinaryData {
		p := filepath.Join(rootFlags.dir, k)
		log.Printf("Writing %s", p)
		if err := os.WriteFile(p, v, 0666); err != nil {
			return xerrors.Errorf("write file: %w", err)
		}
	}

	return nil
}
