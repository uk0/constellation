/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: BUSL-1.1
*/

package resources

import (
	"github.com/edgelesssys/constellation/v2/internal/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

// AuditPolicy defines rulesets for what should be logged in the kube-apiserver audit log.
// reference: https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/ .
type AuditPolicy struct {
	Policy auditv1.Policy
}

// NewDefaultAuditPolicy create a new default Constellation audit policty.
func NewDefaultAuditPolicy() *AuditPolicy {
	return &AuditPolicy{
		Policy: auditv1.Policy{
			TypeMeta: v1.TypeMeta{
				APIVersion: "audit.k8s.io/v1",
				Kind:       "Policy",
			},
			Rules: []auditv1.PolicyRule{
				{
					Level: auditv1.LevelMetadata,
				},
			},
		},
	}
}

// Marshal marshals the audit policy as a YAML document.
func (p *AuditPolicy) Marshal() ([]byte, error) {
	return kubernetes.MarshalK8SResources(p)
}
