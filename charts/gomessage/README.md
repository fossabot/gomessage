
gomessage
===========

Sample application for producing and consuming data with a message queue and TSDB.

## Configuration

The following table lists the configurable parameters of the gomessage chart and their default values.

| Parameter                | Description             | Default        |
| ------------------------ | ----------------------- | -------------- |
| `nameOverride` |  | `""` |
| `fullnameOverride` |  | `""` |
| `imagePullSecrets` |  | `[]` |
| `serviceAccount.create` |  | `true` |
| `serviceAccount.annotations` |  | `{}` |
| `serviceAccount.name` |  | `"gomessage"` |
| `producer.replicaCount` |  | `1` |
| `producer.image.repository` |  | `"nginx"` |
| `producer.image.pullPolicy` |  | `"IfNotPresent"` |
| `producer.image.tag` |  | `"latest"` |
| `producer.podAnnotations` |  | `{}` |
| `producer.podSecurityContext` |  | `{}` |
| `producer.securityContext` |  | `{}` |
| `producer.service.type` |  | `"ClusterIP"` |
| `producer.service.port` |  | `80` |
| `producer.ingress.enabled` |  | `false` |
| `producer.ingress.annotations` |  | `{}` |
| `producer.ingress.hosts` |  | `[]` |
| `producer.ingress.tls` |  | `[]` |
| `producer.autoscaling.enabled` |  | `false` |
| `producer.autoscaling.minReplicas` |  | `1` |
| `producer.autoscaling.maxReplicas` |  | `100` |
| `producer.autoscaling.targetCPUUtilizationPercentage` |  | `80` |
| `producer.autoscaling.targetMemoryUtilizationPercentage` |  | `80` |
| `producer.resources` |  | `{}` |
| `producer.nodeSelector` |  | `{}` |
| `producer.tolerations` |  | `[]` |
| `producer.affinity` |  | `{}` |
| `consumer.replicaCount` |  | `1` |
| `consumer.image.repository` |  | `"nginx"` |
| `consumer.image.pullPolicy` |  | `"IfNotPresent"` |
| `consumer.image.tag` |  | `"latest"` |
| `consumer.podAnnotations` |  | `{}` |
| `consumer.podSecurityContext` |  | `{}` |
| `consumer.securityContext` |  | `{}` |
| `consumer.service.type` |  | `"ClusterIP"` |
| `consumer.service.port` |  | `80` |
| `consumer.ingress.enabled` |  | `false` |
| `consumer.ingress.annotations` |  | `{}` |
| `consumer.ingress.hosts` |  | `[]` |
| `consumer.ingress.tls` |  | `[]` |
| `consumer.autoscaling.enabled` |  | `false` |
| `consumer.autoscaling.minReplicas` |  | `1` |
| `consumer.autoscaling.maxReplicas` |  | `100` |
| `consumer.autoscaling.targetCPUUtilizationPercentage` |  | `80` |
| `consumer.autoscaling.targetMemoryUtilizationPercentage` |  | `80` |
| `consumer.resources` |  | `{}` |
| `consumer.nodeSelector` |  | `{}` |
| `consumer.tolerations` |  | `[]` |
| `consumer.affinity` |  | `{}` |

