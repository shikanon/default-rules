package utils

import (
	"context"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	klogv2 "k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlcache "sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	RESYNC_TIME = 6 * time.Hour
)

func AddScheme(scheme *runtime.Scheme) {
	clientgoscheme.AddToScheme(scheme)
}

func GetControllerClient() (c client.Client, err error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return
	}
	mapper, err := apiutil.NewDynamicRESTMapper(config)
	if err != nil {
		return
	}
	scheme := runtime.NewScheme()
	AddScheme(scheme)
	clientOptions := client.Options{Scheme: scheme, Mapper: mapper}

	// apiReader 不做缓存，所有操作直接操作 apiserver
	apiNoCacheClient, err := client.New(config, clientOptions)
	if err != nil {
		return
	}

	// cache 默认是启动一个watch&list
	// 重新list同步时间
	resyncTime := RESYNC_TIME
	cacheReader, err := ctrlcache.New(config, ctrlcache.Options{Scheme: scheme, Mapper: mapper, Resync: &resyncTime})
	if err != nil {
		return
	}

	// NewDelegatingClient 具有缓存能力，可以将部分缓存起来
	c, err = client.NewDelegatingClient(client.NewDelegatingClientInput{
		CacheReader: cacheReader,
		Client:      apiNoCacheClient,
	})
	// 启动cache
	ctx := context.Background()
	go func() {
		err := cacheReader.Start(ctx)
		if err != nil {
			klogv2.Error(err)
			os.Exit(1)
		}
	}()
	return
}
