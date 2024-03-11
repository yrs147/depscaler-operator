# Deployment Scaler Operator

The Deployment Scaler Operator is a Kubernetes controller designed to dynamically scale deployments based on predefined time periods.

## Overview

This operator utilizes a custom resource named `DepScaler` to specify the scaling schedule for one or more deployments. It continuously monitors the current time and scales the specified deployments to the desired number of replicas if the current time falls within the specified time period (in hours).


## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yrs147/deployment-scaler-operator.git
   cd deployment-scaler-operator
   ```

2. Apply the `Depscaler` CRD

    ```bash
    kubectl apply -f ./config/crd/bases/depscale.yrs.scaler_depscalers.yaml
    ```
3. Verify if the `crd` is created
    ```bash
    kubectl get crd
    ```

4. Create a sample deployment 

    ```bash
    kubectl create deployment nginx --image=nginx
    ```
5. Now create your custom `Depscaler`  using the template below

    ```
    apiVersion: depscale.yrs.scaler/v1
    kind: DepScaler
    metadata:
      labels:
      name: depscaler-sample
    spec:
      begin: 10 # start time hour(in 24hr format)
      end: 17   # end time hour(in 24hr format)
      replicas: 6
      deployments:
        - name: nginx
          namespace: default
    ```

To adjust the deployment scale for a particular time period, modify the begin and end hours in the spec section of the DepScaler manifest. Specify the desired number of replicas for that time period.

