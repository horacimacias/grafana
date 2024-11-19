package provisioning

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafana/grafana/pkg/apimachinery/identity"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/grafana/grafana/pkg/apis/dashboard/v1alpha1"
	"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1"
)

type helloWorldSubresource struct {
	getter        rest.Getter
	statusUpdater rest.Updater
	// for our hacky test access right now
	parent *ProvisioningAPIBuilder
}

func (*helloWorldSubresource) New() runtime.Object {
	// This is added as the "ResponseType" regardless what ProducesObject() returns
	return &v0alpha1.HelloWorld{}
}

func (*helloWorldSubresource) Destroy() {}

func (*helloWorldSubresource) NamespaceScoped() bool {
	return true
}

func (*helloWorldSubresource) GetSingularName() string {
	return "HelloWorld"
}

func (*helloWorldSubresource) ProducesMIMETypes(verb string) []string {
	return []string{"application/json"}
}

func (*helloWorldSubresource) ProducesObject(verb string) any {
	return &v0alpha1.HelloWorld{}
}

func (*helloWorldSubresource) ConnectMethods() []string {
	return []string{http.MethodGet}
}

func (*helloWorldSubresource) NewConnectOptions() (runtime.Object, bool, string) {
	return nil, false, ""
}

func (s *helloWorldSubresource) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := s.getter.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	repo, ok := obj.(*v0alpha1.Repository)
	if !ok {
		return nil, fmt.Errorf("expected repository, but got %t", obj)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		whom := r.URL.Query().Get("whom")
		if whom == "" {
			whom = "World"
		}

		newCommit := r.URL.Query().Get("commit")
		if newCommit != "" {
			repo.Status.CurrentGitCommit = newCommit
			obj, b, err := s.statusUpdater.Update(
				ctx,
				name, // resource name
				rest.DefaultUpdatedObjectInfo(obj, func(ctx context.Context, newObj, oldObj runtime.Object) (transformedNewObj runtime.Object, err error) {
					newObj.(*v0alpha1.Repository).Status.CurrentGitCommit = newCommit
					slog.InfoContext(ctx, "updated the commit", "newObj", newObj, "newCommit", newCommit)
					return newObj, nil
				}),
				func(ctx context.Context, obj runtime.Object) error { return nil },      // createValidation
				func(ctx context.Context, obj, old runtime.Object) error { return nil }, // updateValidation
				false,                   // forceAllowCreate
				&metav1.UpdateOptions{}, // options
			)
			if err != nil {
				responder.Error(err)
				return
			}
			repo = obj.(*v0alpha1.Repository)
			slog.InfoContext(ctx, "the conspicuous boolean", "bool", b)
		}

		write := r.URL.Query().Get("list")
		if write != "" {
			requester, err := identity.GetRequester(r.Context())
			if err != nil {
				responder.Error(err)
				return
			}
			// Propagate user's identity even as we use a service account
			client := s.parent.resourceClientProvider.Client(requester.GetNamespace(), requester.GetIDToken())
			resource, err := client.Resource(schema.GroupVersionKind{
				Group:   v1alpha1.GROUP,   // dashboard.grafana.app
				Version: v1alpha1.VERSION, // "v0alpha1",
				Kind:    "Dashboard",
			})
			if err != nil {
				responder.Error(err)
				return
			}

			ttt, err := resource.List(ctx, metav1.ListOptions{})
			if err != nil {
				responder.Error(err)
				return
			}
			responder.Object(200, ttt) // proxy the list
			return
		}

		slog.InfoContext(ctx, "Got a repository",
			"name", repo.Name,
			"type", repo.Spec.Type,
			"status", repo.Status,
			"local", repo.Spec.Local,
			"github", repo.Spec.GitHub,
			"s3", repo.Spec.S3)

		responder.Object(http.StatusOK, &v0alpha1.HelloWorld{Whom: whom})
	}), nil
}

var (
	_ rest.Storage              = (*helloWorldSubresource)(nil)
	_ rest.Connecter            = (*helloWorldSubresource)(nil)
	_ rest.Scoper               = (*helloWorldSubresource)(nil)
	_ rest.SingularNameProvider = (*helloWorldSubresource)(nil)
	_ rest.StorageMetadata      = (*helloWorldSubresource)(nil)
)
