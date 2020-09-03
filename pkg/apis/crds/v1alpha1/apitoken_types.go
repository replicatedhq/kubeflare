/*
Copyright 2019 Replicated, Inc.

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

package v1alpha1

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func (a APIToken) GetTokenValue(ctx context.Context) (string, error) {
	if a.Spec.Value != "" {
		return a.Spec.Value, nil
	}

	if a.Spec.ValueFrom != nil {
		cfg, err := config.GetConfig()
		if err != nil {
			return "", errors.Wrap(err, "failed to get config")
		}

		clientset, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return "", errors.Wrap(err, "failed to get clientset")
		}

		if a.Spec.ValueFrom.SecretKeyRef != nil {
			secret, err := clientset.CoreV1().Secrets(a.Namespace).Get(ctx, a.Spec.ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
			if err != nil {
				return "", errors.Wrap(err, "failed to get secret")
			}

			return string(secret.Data[a.Spec.ValueFrom.SecretKeyRef.Key]), nil
		}
	}

	return "", errors.New("unable to read api token")
}

type ValueFrom struct {
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// APITokenSpec defines the desired state of APIToken
type APITokenSpec struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Value     string     `json:"value,omitempty"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
}

// APITokenStatus defines the observed state of APIToken
type APITokenStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIToken is the Schema for the APITokens API
// +k8s:openapi-gen=true
type APIToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APITokenSpec   `json:"spec,omitempty"`
	Status APITokenStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APITokenList contains a list of APIToken
type APITokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIToken{}, &APITokenList{})
}
