BIN_DIR := .bin
MINIKUBE_VERSION := v1.25.2
MINIKUBE := $(BIN_DIR)/minikube

$(MINIKUBE):
	mkdir -p $(BIN_DIR)
	curl -Lo $(MINIKUBE) https://storage.googleapis.com/minikube/releases/$(MINIKUBE_VERSION)/minikube-linux-amd64
	chmod +x $(MINIKUBE)

.PHONY: install-minikube
install-minikube: $(MINIKUBE) alias-minikube

.PHONY: start-minikube
start-minikube:
	$(MINIKUBE) start --driver=docker --cpus=2 --memory=2048mb

.PHONY: dashboard-minikube
dashboard-minikube:
	$(MINIKUBE) dashboard

.PHONY: install-argocd
install-argocd:
	kubectl create namespace argocd || true
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

expose-argo:
	kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
	kubectl port-forward svc/argocd-server -n argocd 8080:80

.PHONY: install-grafana-operator
install-grafana-operator:
	kubectl delete -f https://github.com/grafana/grafana-operator/releases/latest/download/kustomize-cluster_scoped.yaml --force --wait --ignore-not-found
	kubectl create -f https://github.com/grafana/grafana-operator/releases/latest/download/kustomize-cluster_scoped.yaml

.PHONY: create-grafana
create-grafana: remove-grafana
	kubectl apply -f ./manifests/00-grafana.yaml -n default

.PHONY: remove-grafana
remove-grafana:
	kubectl delete -f ./manifests/00-grafana.yaml -n default --force --wait --ignore-not-found

.PHONY: expose-grafana
expose-grafana:
	kubectl port-forward svc/test-grafana-service -n default 3000:3000

sleep:
	sleep 10

.PHONY: grafana
grafana: remove-grafana create-grafana sleep expose-grafana