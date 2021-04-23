package generators

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	argoprojiov1alpha1 "github.com/argoproj-labs/applicationset/api/v1alpha1"
	"github.com/argoproj-labs/applicationset/pkg/services/repo_host"
)

func TestRepoHostGetSecretRef(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "test"},
		Data: map[string][]byte{
			"my-token": []byte("secret"),
		},
	}
	gen := &RepoHostGenerator{client: fake.NewClientBuilder().WithObjects(secret).Build()}
	ctx := context.Background()

	cases := []struct {
		name, namespace, token string
		ref                    *argoprojiov1alpha1.SecretRef
		hasError               bool
	}{
		{
			name:      "valid ref",
			ref:       &argoprojiov1alpha1.SecretRef{Name: "test-secret", Key: "my-token"},
			namespace: "test",
			token:     "secret",
			hasError:  false,
		},
		{
			name:      "nil ref",
			ref:       nil,
			namespace: "test",
			token:     "",
			hasError:  false,
		},
		{
			name:      "wrong name",
			ref:       &argoprojiov1alpha1.SecretRef{Name: "other", Key: "my-token"},
			namespace: "test",
			token:     "",
			hasError:  true,
		},
		{
			name:      "wrong key",
			ref:       &argoprojiov1alpha1.SecretRef{Name: "test-secret", Key: "other-token"},
			namespace: "test",
			token:     "",
			hasError:  true,
		},
		{
			name:      "wrong namespace",
			ref:       &argoprojiov1alpha1.SecretRef{Name: "test-secret", Key: "my-token"},
			namespace: "other",
			token:     "",
			hasError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := gen.getSecretRef(ctx, c.ref, c.namespace)
			if c.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, c.token, token)

		})
	}
}

func TestRepoHostGenerateParams(t *testing.T) {
	mockHost := &repo_host.MockRepoHost{
		Repos: []*repo_host.HostedRepo{
			{
				Organization: "myorg",
				Repository:   "repo1",
				URL:          "git@github.com:myorg/repo1.git",
				Branch:       "main",
				Labels:       []string{"prod", "staging"},
			},
			{
				Organization: "myorg",
				Repository:   "repo2",
				URL:          "git@github.com:myorg/repo2.git",
				Branch:       "main",
			},
		},
	}
	gen := &RepoHostGenerator{overrideHost: mockHost}
	params, err := gen.GenerateParams(&argoprojiov1alpha1.ApplicationSetGenerator{
		RepoHost: &argoprojiov1alpha1.RepoHostGenerator{},
	}, nil)
	assert.Nil(t, err)
	assert.Len(t, params, 2)
	assert.Equal(t, "myorg", params[0]["organization"])
	assert.Equal(t, "repo1", params[0]["repository"])
	assert.Equal(t, "git@github.com:myorg/repo1.git", params[0]["url"])
	assert.Equal(t, "main", params[0]["branch"])
	assert.Equal(t, "prod,staging", params[0]["labels"])
	assert.Equal(t, "repo2", params[1]["repository"])
}
