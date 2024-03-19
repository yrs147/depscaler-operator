package controller

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
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

	createDeployment := func(ctx context.Context, cl client.Client) {
		deployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploytName,
				Namespace: deployNamespace,
			},
			Spec: v1.DeploymentSpec{
				Replicas: Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "nginx"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "nginx"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nginx",
								Image: "nginx:latest",
							},
						},
					},
				},
			},
		}
		Expect(cl.Create(ctx, deployment)).To(Succeed())
	}

	Context("Reconcile", func() {
		It("Should reconcile Depscaler custom resources and update replicas", func() {

			cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

			createDeployment(context.Background(), cl)

			testDepScaler := &depscalev1.DepScaler{
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

			Expect(cl.Create(context.Background(), testDepScaler)).To(Succeed(), "Should create DepScaler")

			retrievedDepScaler := &depscalev1.DepScaler{}
			Eventually(func() bool {
				err := cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: depScalerName}, retrievedDepScaler)
				return err == nil
			}).Should(BeTrue(), "DepScaler resource should be created")

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
			fmt.Print(err)

			Expect(res).To(Equal(reconcile.Result{}), "Reconcile result should be empty")

			retrievedDeployment := &v1.Deployment{}
			Eventually(func() int32 {
				err := cl.Get(context.Background(), client.ObjectKey{Namespace: deployNamespace, Name: deploytName}, retrievedDeployment)
				if err != nil {
					return 0
				}
				return *retrievedDeployment.Spec.Replicas
			}).Should(Equal(int32(replicas)), "Deployment replicas should be updated")

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
			Expect(retrieved.Spec.Deployments[0].Name).To(Equal(deploytName))
			Expect(retrieved.Spec.Deployments[0].Namespace).To(Equal(deployNamespace))
			By("Deleting the created Depscaler custom resource")
			Expect(k8sClient.Delete(ctx, test_depscaler)).To(Succeed())
		})

	})
})

func Int32Ptr(i int32) *int32 {
	return &i
}
