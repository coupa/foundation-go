package health

import (
	"net/http"
	"time"
)

/*
 * Service Dependency information struct.
 */
type ServiceDependencyInfo struct {
	name           string
	version        string
	revision       string
	healthCheckUrl string
}

func NewServiceDependency(name string, version string, revision string, healthCheckUrl string) ServiceDependencyInfo {
	serviceDependencyInfo := ServiceDependencyInfo{name, version, revision, healthCheckUrl}
	return serviceDependencyInfo
}

type HTTPEndPointHealthCheckService struct {
	dependentServices []ServiceDependencyInfo
}

func NewHTTPEndPointHealthCheckService(dependentServices []ServiceDependencyInfo) HTTPEndPointHealthCheckService {
	return HTTPEndPointHealthCheckService{dependentServices}
}

/*
 * Returns false is service does not have any dependent services
 */
func (service HTTPEndPointHealthCheckService) HasDependencies() bool {
	return service.dependentServices != nil && len(service.dependentServices) > 0
}

func (service HTTPEndPointHealthCheckService) CheckHttpServiceStatus() []DependencyInfo {
	var dependencies []DependencyInfo
	for i := range service.dependentServices {
		var stateInfo StateInfo
		serviceInfo := service.dependentServices[i]

		sTime := time.Now()
		resp, err := http.Get(serviceInfo.healthCheckUrl)
		responseTime := time.Since(sTime).Seconds()
		if err == nil && resp.StatusCode == http.StatusOK {
			stateInfo.Status = OK
		} else {
			stateInfo.Status = CRIT
		}
		dependencies = append(dependencies, DependencyInfo{serviceInfo.name, "http", serviceInfo.version, serviceInfo.revision, stateInfo, responseTime})
	}
	return dependencies
}
