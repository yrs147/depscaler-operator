package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	depscalev1 "github.com/yrs147/test-operator/api/v1"
)

var _ = Describe("Depscaler Controller", func() {

	const (
		depScalerName   = "test-depscaler"
		namespace       = "default"
		replicas        = 4
		startTime       = 13
		endTime         = 15
		deploytName     = "test-deploy"
		deployNamespace = "default"
	)

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
			Expect(k8sClient.Create(ctx, test_depscaler)).To(Succeed()) // Check the error returned by Create()
		})

	})
})
