package ceph

const ClusterOperator = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: rook-ceph-operator
  namespace: rook-ceph
  labels:
    operator: rook
    storage-backend: ceph
spec:
  selector:
    matchLabels:
      app: rook-ceph-operator
  replicas: 1
  template:
    metadata:
      labels:
        app: rook-ceph-operator
    spec:
      serviceAccountName: rook-ceph-system
      containers:
      - name: rook-ceph-operator
        image: rook/ceph:master
        args: ["ceph", "operator"]
        volumeMounts:
        - mountPath: /var/lib/rook
          name: rook-config
        - mountPath: /etc/ceph
          name: default-config-dir
        env:
        - name: ROOK_CURRENT_NAMESPACE_ONLY
          value: "false"
        - name: ROOK_ALLOW_MULTIPLE_FILESYSTEMS
          value: "false"
        # The logging level for the operator: INFO | DEBUG
        - name: ROOK_LOG_LEVEL
          value: "INFO"
        # The interval to check the health of the ceph cluster and update the status in the custom resource.
        - name: ROOK_CEPH_STATUS_CHECK_INTERVAL
          value: "60s"
        # The interval to check if every mon is in the quorum.
        - name: ROOK_MON_HEALTHCHECK_INTERVAL
          value: "45s"
        # The duration to wait before trying to failover or remove/replace the
        # current mon with a new mon (useful for compensating flapping network).
        - name: ROOK_MON_OUT_TIMEOUT
          value: "600s"
        # The duration between discovering devices in the rook-discover daemonset.
        - name: ROOK_DISCOVER_DEVICES_INTERVAL
          value: "60m"
        # Whether to start pods as privileged that mount a host path, which includes the Ceph mon and osd pods.
        # This is necessary to workaround the anyuid issues when running on OpenShift.
        # For more details see https://github.com/rook/rook/issues/1314#issuecomment-355799641
        - name: ROOK_HOSTPATH_REQUIRES_PRIVILEGED
          value: "false"
        # In some situations SELinux relabelling breaks (times out) on large filesystems, and doesn't work with cephfs ReadWriteMany volumes (last relabel wins).
        # Disable it here if you have similar issues.
        # For more details see https://github.com/rook/rook/issues/2417
        - name: ROOK_ENABLE_SELINUX_RELABELING
          value: "true"
        # In large volumes it will take some time to chown all the files. Disable it here if you have performance issues.
        # For more details see https://github.com/rook/rook/issues/2254
        - name: ROOK_ENABLE_FSGROUP
          value: "true"
        # Disable automatic orchestration when new devices are discovered
        - name: ROOK_DISABLE_DEVICE_HOTPLUG
          value: "false"

        # Whether to enable the flex driver. By default it is enabled and is fully supported, but will be deprecated in some future release
        # in favor of the CSI driver.
        - name: ROOK_ENABLE_FLEX_DRIVER
          value: "true"

        # Whether to start the discovery daemon to watch for raw storage devices on nodes in the cluster.
        # This daemon does not need to run if you are only going to create your OSDs based on StorageClassDeviceSets with PVCs.
        - name: ROOK_ENABLE_DISCOVERY_DAEMON
          value: "true"

        # Enable the default version of the CSI driver. To start another version of the CSI driver, see image properties below.
        - name: ROOK_CSI_ENABLE_CEPHFS
          value: "true"
        - name: ROOK_CSI_ENABLE_RBD
          value: "true"
        - name: ROOK_CSI_ENABLE_GRPC_METRICS
          value: "true"

        # The name of the node to pass with the downward API
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        # The pod name to pass with the downward API
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        # The pod namespace to pass with the downward API
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      volumes:
      - name: rook-config
        emptyDir: {}
      - name: default-config-dir
        emptyDir: {}`
