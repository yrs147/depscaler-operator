# Deployment Scaler Operator

The Deployment Scaler Operator is a Kubernetes controller designed to dynamically scale deployments based on predefined time periods.

## Overview

This operator utilizes a custom resource named `DepScaler` to specify the scaling schedule for one or more deployments. It continuously monitors the current time and scales the specified deployments to the desired number of replicas if the current time falls within the specified time period (in hours).

### Prerequisites
1. Go (version 1.20 or higher)
2. A Running Kubernetes Cluster
   
## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yrs147/deployment-scaler-operator.git
   cd deployment-scaler-operator
   ```

2. Run the `make install` command

    ```bash
    make install
    ```
3. Verify if the `crd` is created
    ```bash
    kubectl get crd
    ```

4. Create a sample deployment 

    ```bash
    kubectl create deployment nginx --image=nginx
    ```
4. Now in another terminal window run  the operator using

    ```bash
    make run
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

7. And see the operator in action 

 ![image](https://github.com/yrs147/depscaler-operator/assets/98258627/09e76a91-95af-4ca3-b9cc-72e8bbfed116)

![image](https://github.com/yrs147/depscaler-operator/assets/98258627/8f9ed5bf-2d30-480b-ae6e-731ee70d3def)



To adjust the deployment scale for a particular time period, modify the begin and end hours in the spec section of the DepScaler manifest. Specify the desired number of replicas for that time period.

