/*
Copyright 2021-2022 Red Hat, Inc.

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotValidatingWebhook(t *testing.T) {

	// Create initial Snapshot CR.
	originalSnapshot := Snapshot{
		Spec: SnapshotSpec{
			Application: "test-app-a",
			Components: []SnapshotComponent{
				{
					Name:           "test-component-a",
					ContainerImage: "test-container-image-a",
				},
			},
		},
	}

	tests := []struct {
		testName      string   // Name of test
		testData      Snapshot // Test data to be passed to webhook function
		expectedError string   // Expected error message from webhook function
	}{
		{
			testName: "No error when Spec is same.",
			testData: Snapshot{
				Spec: SnapshotSpec{
					Application: "test-app-a",
					Components: []SnapshotComponent{
						{
							Name:           "test-component-a",
							ContainerImage: "test-container-image-a",
						},
					},
				},
			},
			expectedError: "",
		},

		{
			testName: "Error occurs when Spec.Application name is changed.",
			testData: Snapshot{
				Spec: SnapshotSpec{
					Application: "test-app-a-changed",
					Components: []SnapshotComponent{
						{
							Name:           "test-component-a",
							ContainerImage: "test-container-image-a",
						},
					},
				},
			},
			expectedError: "application cannot be updated to test-app-a-changed",
		},

		{
			testName: "Error occurs when Spec.Components.Name is changed.",
			testData: Snapshot{
				Spec: SnapshotSpec{
					Application: "test-app-a",
					Components: []SnapshotComponent{
						{
							Name:           "test-component-a-changed",
							ContainerImage: "test-container-image-a",
						},
					},
				},
			},
			expectedError: "components cannot be updated to [{Name:test-component-a-changed ContainerImage:test-container-image-a}]",
		},

		{
			testName: "Error occurs when Spec.Components.ContainerImage is changed.",
			testData: Snapshot{
				Spec: SnapshotSpec{
					Application: "test-app-a",
					Components: []SnapshotComponent{
						{
							Name:           "test-component-a",
							ContainerImage: "test-container-image-a-changed",
						},
					},
				},
			},
			expectedError: "components cannot be updated to [{Name:test-component-a ContainerImage:test-container-image-a-changed}]",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualError := test.testData.ValidateUpdate(&originalSnapshot)

			if test.expectedError == "" {
				assert.Nil(t, actualError)
			} else {
				assert.Contains(t, actualError.Error(), test.expectedError)
			}
		})
	}
}
