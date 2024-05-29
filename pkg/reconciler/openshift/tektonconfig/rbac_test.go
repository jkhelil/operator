package tektonconfig

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	fakesecurity "github.com/openshift/client-go/security/clientset/versioned/fake"
	"github.com/stretchr/testify/require"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	fakeoperator "github.com/tektoncd/operator/pkg/client/injection/client/fake"
	"github.com/tektoncd/operator/pkg/reconciler/common"
	"github.com/tektoncd/pipeline/test/diff"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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
	s := runtime.NewScheme()
	utilruntime.Must(v1alpha1.AddToScheme(s))
	fakek8s.AddToScheme(s)

	for _, c := range []struct {
		desc string
		objs []runtime.Object
		iSet *v1alpha1.TektonInstallerSet
		err  error
	}{
		/*{
			desc: "No existing installer set",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
			},
			err: v1alpha1.RECONCILE_AGAIN_ERR,
		},*/
		{
			desc: "existing installer set",
			objs: []runtime.Object{
				&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: map[string]string{"openshift-pipelines.tekton.dev/namespace-reconcile-version": ""}}},
			},
			iSet: &v1alpha1.TektonInstallerSet{ObjectMeta: metav1.ObjectMeta{Name: "rbac-resources", Labels: map[string]string{v1alpha1.CreatedByKey: createdByValue, v1alpha1.InstallerSetType: componentNameRBAC}}, Spec: v1alpha1.TektonInstallerSetSpec{}},
			err:  nil,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			fmt.Println("check debut test")
			kubeclient := fakek8s.NewSimpleClientset(c.objs...)
			fakesecurityclient := fakesecurity.NewSimpleClientset()
			fakeoperatorclient := fakeoperator.Get(ctx)
			/*fakeDiscoveryForOp, ok := fakeoperatorclient.Discovery().(*fakediscovery.FakeDiscovery)
			if !ok {
				t.Error("failed to convert clientset discovery to fake discovery")
			}
			fakeDiscoveryForOp.Fake.Resources = []*metav1.APIResourceList{
				{
					GroupVersion: "operator.tekton.dev/v1alpha1",
					APIResources: []metav1.APIResource{
						{
							Name:         "tektoninstallersets",
							SingularName: "tektoninstallerset",
							Kind:         "TektonInstallerSet",
							Group:        "operator.tekton.dev",
							Version:      "v1alpha1",
							Namespaced:   true,
							Verbs:        []string{"create", "delete", "get", "list", "patch", "update", "watch"},
						},
					},
				},
			}*/
			if c.iSet != nil {
				fakeoperatorclient.OperatorV1alpha1().TektonInstallerSets().Create(ctx, c.iSet, metav1.CreateOptions{})
				//labelSelector, _ := common.LabelSelector(rbacInstallerSetSelector)
				//existingInstallerSet, err := tektoninstallerset.CurrentInstallerSetName(ctx, fakeoperatorclient, labelSelector)
				//if err != nil {
				//	fmt.Errorf("failed to retreive existing InstallerSet with selector %v: %w", labelSelector, err)
				//}

				//iSets, _ := fakeoperatorclient.Operator1alpha1().TektonInstallerSets().List(ctx, v1.ListOptions{
				//	LabelSelector: labelSelector,
				//})
				//	fmt.Println(existingInstallerSet)
			}
			fmt.Println(fakeoperatorclient)
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
			require.ErrorIs(t, err, c.err)
		})
	}
}
