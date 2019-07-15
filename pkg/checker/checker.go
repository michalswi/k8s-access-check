package checker

import (
	"fmt"
	"log"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
	authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

type canIOptions struct {
	authClient authorizationv1client.AuthorizationV1Interface
	namespace  string
}

type Params struct {
	Kclient *kubernetes.Clientset
	Decjson map[string]interface{}
	Ns      string
}

func (p *Params) Runner() {
	p.GetKubeVersion()
	p.WhatCanIdo()
	// p.WhatCanIdoList()
}

// func GetKubeVersion(kclient *kubernetes.Clientset) {
func (p *Params) GetKubeVersion() {
	version, err := p.Kclient.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("K8s version: %s", version)
}

// Lists all the accesses that current user has
// $ kubectl auth can-i --list --namespace=<namespace>
func (p *Params) WhatCanIdoList() {
	// add to import:
	// authorizationv1 "k8s.io/api/authorization/v1"
	// authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"

	// SelfSubjectAccessReview (SSAR)
	ssar := &authorizationv1.SelfSubjectRulesReview{
		Spec: authorizationv1.SelfSubjectRulesReviewSpec{
			Namespace: p.Ns,
		},
	}

	o := &canIOptions{}
	o.authClient = p.Kclient.AuthorizationV1()
	response, err := o.authClient.SelfSubjectRulesReviews().Create(ssar)
	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Println(response.Status)
}

// Checks if user has access to a certain resource
// $ kubectl auth can-i get deployments
// $ kubectl auth can-i get deployments -n kube-system
func (p *Params) WhatCanIdo() {
	// add to import:
	// authorizationv1 "k8s.io/api/authorization/v1"
	// authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"

	var actionVerb string
	var actionRsc string
	var actionGroup string
	var actionGroupVer string

	for verb, rsrc := range p.Decjson {
		actionVerb = verb
		for _, r := range rsrc.([]interface{}) {
			rFormat := fmt.Sprintf("%v", r)
			actionRsc = rFormat

			// /apis/apps/v1/deployments
			switch actionRsc {
			case "deployments":
				actionGroup = "apps"
				actionGroupVer = "v1"
			case "pods":
				actionGroup = ""
				actionGroupVer = "v1"
			}

			// SelfSubjectAccessReview (SSAR)
			// ssar := &authorizationv1.SelfSubjectAccessReview{
			var ssar *authorizationv1.SelfSubjectAccessReview
			ssar = &authorizationv1.SelfSubjectAccessReview{
				Spec: authorizationv1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: &authorizationv1.ResourceAttributes{
						Namespace: p.Ns,
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

			o := &canIOptions{}
			o.authClient = p.Kclient.AuthorizationV1()

			response, err := o.authClient.SelfSubjectAccessReviews().Create(ssar)
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
	}
}
