# Hyperledger Fabric Chaincode

## Introduction
Hyperledger Fabric Chaincode is a run time environment to run/endorse the transacion before commiting it to the ledger.

Learn more about it by visiting https://hyperledger-fabric.readthedocs.io/en/release-2.5

## Prerequisites

- Kubernetes 1.9+
- PV provisioner support in the underlying infrastructure.

## Installing the chart

```
helm install enoc-chaincode ./chaincode --dry-run --debug > final.yaml
helm install enoc-chaincode ./chaincode

```

## Updating the chart
To upgrade the c deployment:
```
helm upgrade enoc-chaincode ./chaincode --install
```

## Uninstalling the chart
To uninstall/delete the ca deployment:
```
helm delete enoc-chaincode
```



