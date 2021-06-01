package k8s

import (
	"context"
	"os"
	"path/filepath"
	"unicode/utf8"

	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type ConfigMapUpdater struct {
	Clientset *kubernetes.Clientset
}

func (u *ConfigMapUpdater) Update(ctx context.Context, namespace, name, dir string, allowDelete bool) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cm, err := u.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return xerrors.Errorf("get configmap: %w", err)
		}
		if cm.Data == nil || allowDelete {
			cm.Data = map[string]string{}
		}
		if cm.BinaryData == nil || allowDelete {
			cm.BinaryData = map[string][]byte{}
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return xerrors.Errorf("read dir: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			b, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return xerrors.Errorf("read file: %w", err)
			}
			if utf8.Valid(b) {
				cm.Data[e.Name()] = string(b)
			} else {
				cm.BinaryData[e.Name()] = b
			}
		}

		_, err = u.Clientset.CoreV1().ConfigMaps(namespace).Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			return xerrors.Errorf("update configmap: %w", err)
		}

		return nil
	})
}
