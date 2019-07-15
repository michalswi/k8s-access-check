package checker

import (
	"fmt"
	"log"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
	authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubectl/cmd/auth/cani.go#L234
type canIOptions struct {
	authClient authorizationv1client.AuthorizationV1Interface
	namespace  string
}

func GetKubeVersion(kclient *kubernetes.Clientset) {
	version, err := kclient.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("K8s version: %s", version)
}

// Lists all the accesses that current user has
// $ kubectl auth can-i --list --namespace=<namespace>
func WhatCanIdoList(kclient *kubernetes.Clientset, ns string) {
	// add to import:
	// authorizationv1 "k8s.io/api/authorization/v1"
	// authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"

	o := &canIOptions{}

	o.authClient = kclient.AuthorizationV1()

	// SelfSubjectAccessReview (SSAR)
	ssar := &authorizationv1.SelfSubjectRulesReview{
		Spec: authorizationv1.SelfSubjectRulesReviewSpec{
			Namespace: ns,
		},
	}

	response, err := o.authClient.SelfSubjectRulesReviews().Create(ssar)
	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Println(response.Status)
}

// Checks if user has access to a certain resource
// $ kubectl auth can-i get deployments
// $ kubectl auth can-i get deployments -n kube-system
func WhatCanIdo(kclient *kubernetes.Clientset, ns string) {
	// add to import:
	// authorizationv1 "k8s.io/api/authorization/v1"
	// authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"

	var actionGroup string
	var actionGroupVer string

	// /apis/apps/v1/deployments
	actionVerb := "get"
	actionRsc := "deployments"
	// actionRsc := "pods"

	switch actionRsc {
	case "deployments":
		actionGroup = "apps"
		actionGroupVer = "v1"
	case "pods":
		actionGroup = ""
		actionGroupVer = "v1"
	}

	o := &canIOptions{}

	o.authClient = kclient.AuthorizationV1()

	// SelfSubjectAccessReview (SSAR)
	var ssar *authorizationv1.SelfSubjectAccessReview
	ssar = &authorizationv1.SelfSubjectAccessReview{
		// ssar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: ns,
				Verb:      actionVerb,
				Resource:  actionRsc,
				Group:     actionGroup,
				Version:   actionGroupVer,
			},
			// NonResourceAttributes: &authorizationv1.NonResourceAttributes{
			// 	Verb: actionVerb,
			// },
		},
	}

	response, err := o.authClient.SelfSubjectAccessReviews().Create(ssar)
	// response, err := kclient.AuthorizationV1().SelfSubjectAccessReviews().Create(ssar)
	if err != nil {
		log.Printf("%v", err)
	}

	// fmt.Println(response)
	// fmt.Println(response.Spec)
	// fmt.Println(response.Status.Allowed)

	if response.Status.Allowed {
		log.Printf("User can /%s/ a /%s/, status: ALLOWED", actionVerb, actionRsc)
	} else {
		log.Printf("User can /%s/ a /%s/, status: NOTALLOWED", actionVerb, actionRsc)
		if len(response.Status.Reason) > 0 {
			log.Printf("%v", response.Status.Reason)
		}
		if len(response.Status.EvaluationError) > 0 {
			log.Printf("%v", response.Status.EvaluationError)
		}
	}

}
