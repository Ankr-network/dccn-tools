package ceph

const ClusterStorageClass = `apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
   name: rook-ceph-block
   annotations: 
     storageclass.kubernetes.io/is-default-class: "true"
provisioner: ceph.rook.io/block
parameters:
  blockPool: replicapool
  clusterNamespace: rook-ceph
  fstype: xfs
reclaimPolicy: Retain`
