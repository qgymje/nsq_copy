package clusterinfo

import (
	"github.com/qgymje/nsq_copy/internal/http_api"
)

type logger interface {
	Output(maxdepth int, s string) error
}

type ClusterInfo struct {
	log    logger
	client *http_api.Client
}
