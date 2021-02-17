Activate SSL/TLS certificates for HTTPS access? Important note: this will NOT work in a development environment. [y/N] n
Configuration saved to /Users/rosecrisp/Library/Application Support/tutor/config.yml
================================================
        Updating the current environment
================================================
Environment generated in /Users/rosecrisp/Library/Application Support/tutor/env
=====================================
        Starting the platform
=====================================
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --wait --selector app.kubernetes.io/component=namespace
Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
namespace/openedx configured
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --wait --selector app.kubernetes.io/component=volume
persistentvolumeclaim/caddy created
persistentvolumeclaim/elasticsearch created
persistentvolumeclaim/mongodb created
persistentvolumeclaim/mysql created
persistentvolumeclaim/redis created
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --selector app.kubernetes.io/component!=job
namespace/openedx unchanged
configmap/caddy-config-tfkdm2749k created
configmap/nginx-config-9mkd96879c created
configmap/openedx-config-6b7kc8gk9c created
configmap/openedx-settings-cms-b8bgf47ck2 created
configmap/openedx-settings-lms-287ft5cmbg created
configmap/redis-config-5m6b9th9d6 created
service/caddy created
service/cms created
service/elasticsearch created
service/forum created
service/lms created
service/mongodb created
service/mysql created
service/nginx created
service/redis created
service/smtp created
deployment.apps/caddy created
deployment.apps/cms-worker created
deployment.apps/cms created
deployment.apps/elasticsearch created
deployment.apps/forum created
deployment.apps/lms-worker created
deployment.apps/lms created
deployment.apps/mongodb created
deployment.apps/mysql created
deployment.apps/nginx created
deployment.apps/redis created
deployment.apps/smtp created
persistentvolumeclaim/caddy unchanged
persistentvolumeclaim/elasticsearch unchanged
persistentvolumeclaim/mongodb unchanged
persistentvolumeclaim/mysql unchanged
persistentvolumeclaim/redis unchanged
================================================
        Database creation and migrations
================================================
Waiting for a mysql pod to be ready...
kubectl wait --namespace openedx --selector=app.kubernetes.io/instance=openedx-xmWN7pgOQHszibWniaFBoui7,app.kubernetes.io/name=mysql --for=condition=ContainersReady --timeout=600s pod
pod/mysql-c666b9948-2msgn condition met
Waiting for a elasticsearch pod to be ready...
kubectl wait --namespace openedx --selector=app.kubernetes.io/instance=openedx-xmWN7pgOQHszibWniaFBoui7,app.kubernetes.io/name=elasticsearch --for=condition=ContainersReady --timeout=600s pod
pod/elasticsearch-7fc9477cb9-klqzp condition met
Waiting for a mongodb pod to be ready...
kubectl wait --namespace openedx --selector=app.kubernetes.io/instance=openedx-xmWN7pgOQHszibWniaFBoui7,app.kubernetes.io/name=mongodb --for=condition=ContainersReady --timeout=600s pod
pod/mongodb-8577d8c76f-g2nxc condition met
Initialising all services...
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --selector app.kubernetes.io/name=mysql-job-20210217094027
job.batch/mysql-job-20210217094027 created
Job mysql-job-20210217094027 is running. To view the logs from this job, run:

    kubectl logs --namespace=openedx --follow $(kubectl get --namespace=openedx pods --selector=job-name=mysql-job-20210217094027 -o=jsonpath="{.items[0].metadata.name}")

Waiting for job completion...
Job mysql-job-20210217094027 successful.
Initialising lms...
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --selector app.kubernetes.io/name=lms-job-20210217094253
job.batch/lms-job-20210217094253 created
Job lms-job-20210217094253 is running. To view the logs from this job, run:

    kubectl logs --namespace=openedx --follow $(kubectl get --namespace=openedx pods --selector=job-name=lms-job-20210217094253 -o=jsonpath="{.items[0].metadata.name}")

Waiting for job completion...
Job lms-job-20210217094253 successful.
Initialising cms...
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --selector app.kubernetes.io/name=cms-job-20210217094938
job.batch/cms-job-20210217094938 created
Job cms-job-20210217094938 is running. To view the logs from this job, run:

    kubectl logs --namespace=openedx --follow $(kubectl get --namespace=openedx pods --selector=job-name=cms-job-20210217094938 -o=jsonpath="{.items[0].metadata.name}")

Waiting for job completion...
Job cms-job-20210217094938 successful.
Initialising forum...
kubectl apply --kustomize /Users/rosecrisp/Library/Application Support/tutor/env --selector app.kubernetes.io/name=forum-job-20210217095005
job.batch/forum-job-20210217095005 created
Job forum-job-20210217095005 is running. To view the logs from this job, run:

    kubectl logs --namespace=openedx --follow $(kubectl get --namespace=openedx pods --selector=job-name=forum-job-20210217095005 -o=jsonpath="{.items[0].metadata.name}")

Waiting for job completion...
Job forum-job-20210217095005 successful.
All services initialised.
Your Open edX platform is ready and can be accessed at the following urls:

    http://www.lms-openedx.apps.courses.operatortrain.me
    http://studio.www.lms-openedx.apps.courses.operatortrain.me
    