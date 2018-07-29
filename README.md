# kubernetes-pod-restarter

[![Build Status](https://travis-ci.org/rebuy-de/kubernetes-pod-restarter.svg?branch=master)](https://travis-ci.org/rebuy-de/kubernetes-pod-restarter)
[![license](https://img.shields.io/github/license/rebuy-de/kubernetes-pod-restarter.svg)]()
[![GitHub release](https://img.shields.io/github/release/rebuy-de/kubernetes-pod-restarter.svg)](https://github.com/rebuy-de/kubernetes-pod-restarter/releases)

Deletes targeted Kubernetes Pods to force a regular restart.

> **Development Status** *kubernetes-pod-restarter* is in an early development phase. Expect
> breaking changes any time.

## Use Cases

* *kubernetes-pod-restarter* can be used as a poor mans solution to memory leaks.
  * **Why don't you fix the leaks instead?** Memory leaks are often hard to
    find. This is especially the case, if it is about a third-party
    application. Probybly most of the leaks a fixable, but this might not worth
    the time for minor important services.
  * **Why don't you let Kubernetes OOMKiller do the job?** An OOMKill is not
    a planned event and indicates an issue with a service or within the system.
    Therefore we monitor OOMKills and getting regular notification about
    OOMKills, would shadow actual alerts.

## Usage

1. Deploy Custom Resource Definition (CRD):
   * `kubectl apply -f deploy/crd.yaml`
2. Deploy RBAC permissions:
   * `kubectl apply -f deploy/rbac.yaml`
3. Deploy operator:
   * TBD

## Developing

The *kubernetes-pod-restarter* is based on the [CoreOS Operator
SDK](https://github.com/operator-framework/operator-sdk). Take a look at the
docs for an introduction.

