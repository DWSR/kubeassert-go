package assertion

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

type listOptionsFunc func(*commonAssertion, *envconf.Config) metav1.ListOptions

func defaultListOptions(ca *commonAssertion, cfg *envconf.Config) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(ca.assertLabels)).String(),
		FieldSelector: fields.SelectorFromSet(fields.Set(ca.assertFields)).String(),
	}
}

func listOptionsWithNamespaceFromEnv(ca *commonAssertion, cfg *envconf.Config) metav1.ListOptions {
	selectorFields := fields.Set(ca.assertFields)
	selectorFields["metadata.namespace"] = cfg.Namespace()

	return metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(ca.assertLabels)).String(),
		FieldSelector: fields.SelectorFromSet(selectorFields).String(),
	}
}
