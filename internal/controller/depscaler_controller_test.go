package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	depscalev1 "github.com/yrs147/test-operator/api/v1"
)

var _ = Describe("Depscaler Controller", func() {

	const (
		depScalerName   = "test-depscaler"
		namespace       = "default"
		replicas        = 4
		startTime       = 0
		endTime         = 23
		deploytName     = "test-deploy"
		deployNamespace = "default"
	)

	Context("Reconcile", func() {
		It("Should reconcile Depscaler custom resources", func() {
			// fake client
			cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

			r := &DepScalerReconciler{
				Client: cl,
				Scheme: scheme.Scheme,
			}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      depScalerName,
					Namespace: namespace,
				},
			}

			res, err := r.Reconcile(context.Background(), req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(reconcile.Result{}))

			retrieved := &depscalev1.DepScaler{}
			Eventually(func() bool {

				if err = cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: depScalerName}, retrieved); err != nil {
					return false
				}
				return true
			})

			Expect(retrieved.Status).ToNot(BeNil())
		})
	})

	Context("When setting up manager", func() {
		It("Should set up manager successfully", func() {

			mgr, err := ctrl.NewManager(cfg, ctrl.Options{
				Scheme: scheme.Scheme,
			})
			Expect(err).ToNot(HaveOccurred())

			reconciler := &DepScalerReconciler{
				Client: mgr.GetClient(),
				Scheme: mgr.GetScheme(),
			}

			err = reconciler.SetupWithManager(mgr)
			Expect(err).ToNot(HaveOccurred())

		})
	})

	Context("When setting up test environment", func() {
		It("Should create Depscaler custom resources", func() {
			By("Creating a first Depscaler custom resource")
			ctx := context.Background()
			test_depscaler := &depscalev1.DepScaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      depScalerName,
					Namespace: namespace,
				},
				Spec: depscalev1.DepScalerSpec{
					Begin:    startTime,
					End:      endTime,
					Replicas: replicas,
					Deployments: []depscalev1.NamespacedName{
						{
							Name:      deploytName,
							Namespace: deployNamespace,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, test_depscaler)).To(Succeed())
			By("Retrieving the created Depscaler custom resource")
			retrieved := &depscalev1.DepScaler{}
			Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: depScalerName}, retrieved)).To(Succeed())
			Expect(retrieved.Spec.Begin).To(Equal(startTime))
			Expect(retrieved.Spec.End).To(Equal(endTime))
			Expect(retrieved.Spec.Deployments).To(HaveLen(1))
			// Expect(retrieved.Spec.Replicas).To(Equal(replicas))
			Expect(retrieved.Spec.Deployments[0].Name).To(Equal(deploytName))
			Expect(retrieved.Spec.Deployments[0].Namespace).To(Equal(deployNamespace))

			By("Deleting the created Depscaler custom resource")
			Expect(k8sClient.Delete(ctx, test_depscaler)).To(Succeed())
		})

	})
})
