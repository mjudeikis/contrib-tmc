/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package synctarget

import (
	"context"

	utilserrors "k8s.io/apimachinery/pkg/util/errors"

	workloadv1alpha1 "github.com/faroshq/tmc/apis/workload/v1alpha1"
)

type reconcileStatus int

const (
	reconcileStatusStopAndRequeue reconcileStatus = iota
	reconcileStatusContinue
)

type reconciler interface {
	reconcile(ctx context.Context, syncTarget *workloadv1alpha1.SyncTarget) (reconcileStatus, error)
}

func (c *controller) reconcile(ctx context.Context, syncTarget *workloadv1alpha1.SyncTarget) (bool, error) {
	var errs []error

	requeue := false
	for _, r := range c.reconcilers {
		var err error
		var status reconcileStatus
		status, err = r.reconcile(ctx, syncTarget)
		if err != nil {
			errs = append(errs, err)
		}
		if status == reconcileStatusStopAndRequeue {
			requeue = true
			break
		}
	}

	return requeue, utilserrors.NewAggregate(errs)
}
