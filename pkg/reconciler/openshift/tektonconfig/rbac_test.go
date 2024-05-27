package tektonconfig

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	fakesecurity "github.com/openshift/client-go/security/clientset/versioned/fake"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	fakeoperator "github.com/tektoncd/operator/pkg/client/injection/client/fake"
	"github.com/tektoncd/operator/pkg/reconciler/common"
	"github.com/tektoncd/pipeline/test/diff"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	fakek8s "k8s.io/client-go/kubernetes/fake"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	ts "knative.dev/pkg/reconciler/testing"
)

func TestGetNamespacesToBeReconciled(t *testing.T) {
	var deletionTime = metav1.Now()
	for _, c := range []struct {
		desc           string
		wantNamespaces []corev1.Namespace
		objs           []runtime.Object
		ctx            context.Context
	}{
		{
			desc: "ignore system namespaces",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "openshift-test"}},
			},
			wantNamespaces: nil,
			ctx:            context.Background(),
		},
		{
			desc: "ignore namespaces with deletion timestamp",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "openshift-test", DeletionTimestamp: &deletionTime}},
			},
			wantNamespaces: nil,
			ctx:            context.Background(),
		},
		{
			desc: "add namespace to reconcile list if it has openshift scc operator.tekton.dev/scc annotation set ",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test", Annotations: map[string]string{"operator.tekton.dev/scc": "restricted"}}},
			},
			wantNamespaces: []corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test",
						Annotations: map[string]string{"operator.tekton.dev/scc": "restricted"},
					},
				},
			},
			ctx: context.Background(),
		},
		{
			desc: "add namespace to reconcile list if it has bad label openshift-pipelines.tekton.dev/namespace-reconcile-version",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
			},
			wantNamespaces: []corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "test",
						Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""},
					},
				},
			},
			ctx: context.Background(),
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			kubeclient := fakek8s.NewSimpleClientset(c.objs...)
			r := rbac{
				kubeClientSet: kubeclient,
				version:       "devel",
			}
			namespaces, err := r.getNamespacesToBeReconciled(c.ctx)
			if err != nil {
				t.Fatalf("getNamespacesToBeReconciled: %v", err)
			}
			if d := cmp.Diff(c.wantNamespaces, namespaces); d != "" {
				t.Fatalf("Diff %s", diff.PrintWantGot(d))
			}
		})
	}
}

func TestCreateResources(t *testing.T) {
	ctx, _, _ := ts.SetupFakeContextWithCancel(t)
	os.Setenv(common.KoEnvKey, "testdata")

	for _, c := range []struct {
		desc string
		objs []runtime.Object
	}{
		{
			desc: "reconcile single namespace",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
			},
		},
		{
			desc: "reconcile multiple namepsace",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test1", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test2", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test3", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
			},
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			kubeclient := fakek8s.NewSimpleClientset(c.objs...)
			fakesecurityclient := fakesecurity.NewSimpleClientset()
			fakeoperatorclient := fakeoperator.Get(ctx)

			informers := informers.NewSharedInformerFactory(kubeclient, 0)
			nsInformer := informers.Core().V1().Namespaces()
			rbacinformer := informers.Rbac().V1().ClusterRoleBindings()
			tc := &v1alpha1.TektonConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: v1alpha1.ConfigResourceName,
				},
				Spec: v1alpha1.TektonConfigSpec{
					CommonSpec: v1alpha1.CommonSpec{
						TargetNamespace: "foo",
					},
				},
				Status: v1alpha1.TektonConfigStatus{
					Status:             duckv1.Status{},
					TektonInstallerSet: map[string]string{},
				},
			}

			r := rbac{
				kubeClientSet:     kubeclient,
				operatorClientSet: fakeoperatorclient,
				securityClientSet: fakesecurityclient,
				rbacInformer:      rbacinformer,
				nsInformer:        nsInformer,
				version:           "devel",
				tektonConfig:      tc,
			}
			err := r.createResources(ctx)
			if err != nil {
				t.Fatalf("createResources: %v", err)
			}
		})
	}
}
