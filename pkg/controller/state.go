package controller

import (
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
)

func getCurrentState(client runtime.Client) (result []model.Pod, err error) {
	result = []model.Pod{}
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return result, err
	}

	for _, namespace := range namespaces {
		containers, err := client.GetContainers(namespace)
		if err != nil {
			return result, err
		}

		result = append(result, constructPodsFromContainerInfo(client, namespace, containers)...)
	}
	return result, nil
}

func constructPodsFromContainerInfo(client runtime.Client, namespace string, containersByPods map[string][]model.Container) (result []model.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, model.Pod{
			Metadata: model.NewMetadata(
				podName,
				namespace,
			),
			Spec: model.PodSpec{
				Containers: containers,
			},
			Status: model.PodStatus{
				ContainerStatuses: resolveContainerStatuses(client, containers),
			},
		})
	}

	return result
}

func resolveContainerStatuses(client runtime.Client, containers []model.Container) []model.ContainerStatus {
	containerStatuses := []model.ContainerStatus{}
	for _, container := range containers {
		containerStatuses = append(containerStatuses, resolveContainerStatus(client, container))
	}
	return containerStatuses
}

func resolveContainerStatus(client runtime.Client, container model.Container) model.ContainerStatus {
	return model.ContainerStatus{
		ContainerID: container.ID,
		Image:       container.Image,
		State:       client.GetContainerTaskStatus(container.ID),
	}
}

func getValues(podsByName map[string]model.Pod) []model.Pod {
	values := []model.Pod{}
	for _, pod := range podsByName {
		values = append(values, pod)
	}
	return values
}
