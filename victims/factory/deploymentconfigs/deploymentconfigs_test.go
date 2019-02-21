package deploymentconfigs

import (
	"testing"

	"github.com/asobti/kube-monkey/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	IDENTIFIER = "kube-monkey-id"
	NAME       = "deploymentconfig_name"
	NAMESPACE  = metav1.NamespaceDefault
)

func newDeploymentConfig(name string, labels map[string]string) v1.DeploymentConfig {

	return v1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: NAMESPACE,
			Labels:    labels,
		},
	}
}

func TestNew(t *testing.T) {

	v1depl := newDeploymentConfig(
		NAME,
		map[string]string{
			config.IdentLabelKey: IDENTIFIER,
			config.MtbfLabelKey:  "1",
		},
	)
	depl, err := New(&v1depl)

	assert.NoError(t, err)
	assert.Equal(t, "v1.DeploymentConfig", depl.Kind())
	assert.Equal(t, NAME, depl.Name())
	assert.Equal(t, NAMESPACE, depl.Namespace())
	assert.Equal(t, IDENTIFIER, depl.Identifier())
	assert.Equal(t, 1, depl.Mtbf())
}

func TestInvalidIdentifier(t *testing.T) {
	v1depl := newDeploymentConfig(
		NAME,
		map[string]string{
			config.MtbfLabelKey: "1",
		},
	)
	_, err := New(&v1depl)

	assert.Errorf(t, err, "Expected an error if "+config.IdentLabelKey+" label doesn't exist")
}

func TestInvalidMtbf(t *testing.T) {
	v1depl := newDeploymentConfig(
		NAME,
		map[string]string{
			config.IdentLabelKey: IDENTIFIER,
		},
	)
	_, err := New(&v1depl)

	assert.Errorf(t, err, "Expected an error if "+config.MtbfLabelKey+" label doesn't exist")

	v1depl = newDeploymentConfig(
		NAME,
		map[string]string{
			config.IdentLabelKey: IDENTIFIER,
			config.MtbfLabelKey:  "string",
		},
	)
	_, err = New(&v1depl)

	assert.Errorf(t, err, "Expected an error if "+config.MtbfLabelKey+" label can't be converted a Int type")

	v1depl = newDeploymentConfig(
		NAME,
		map[string]string{
			config.IdentLabelKey: IDENTIFIER,
			config.MtbfLabelKey:  "0",
		},
	)
	_, err = New(&v1depl)

	assert.Errorf(t, err, "Expected an error if "+config.MtbfLabelKey+" label is lower than 1")
}
