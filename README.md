# Deployment Scaler Operator

The Deployment Scaler Operator is a Kubernetes controller designed to dynamically scale deployments based on predefined time periods.

## Overview

This operator utilizes a custom resource named `DepScaler` to specify the scaling schedule for one or more deployments. It continuously monitors the current time and scales the specified deployments to the desired number of replicas if the current time falls within the specified time period (in hours).

### Prerequisites
1. Go version 1.20 (only required if running locally)
2. A Running Kubernetes Cluster

### Running on the cluster


1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/depscaler:tag
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/depscaler:tag
```

3. Create a sample deployment 

```bash
kubectl create deployment nginx --image=nginx
```

4. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```
   
## Running Locally against Cluster

1. Clone the repository:

   ```bash
   git clone https://github.com/yrs147/deployment-scaler-operator.git
   cd deployment-scaler-operator
   ```

2. Run the `make install` command

    ```bash
    make install
    ```

3. Create a sample deployment 

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

### Testing the Controller
Use the `make test` command to run tests

```
make test
```

![image](https://github.com/yrs147/depscaler-operator/assets/98258627/baa3850d-ebaf-4259-b3df-842ce3c92101)



### ScreenShots 
Operator in action 

 ![image](https://github.com/yrs147/depscaler-operator/assets/98258627/67ff7d91-aeb2-41ae-a652-0dfc059db142)


![image](https://github.com/yrs147/depscaler-operator/assets/98258627/8f9ed5bf-2d30-480b-ae6e-731ee70d3def)



To adjust the deployment scale for a particular time period, modify the begin and end hours in the spec section of the DepScaler manifest. Specify the desired number of replicas for that time period.

