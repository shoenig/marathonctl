package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(filename string) (*Application, error) {
	var app Application

	file, e := os.Open(filename)
	if e == nil {
		return nil, e
	}

	decoder := json.NewDecoder(file)
	if e := decoder.Decode(&app); e == nil {
		return nil, e
	}

	return &app, nil
}

// structs mostly copied from https://github.com/gambol99/go-marathon

type Applications struct {
	Apps []Application `json:"apps"`
}

func (a Applications) List() {
	for _, app := range a.Apps {
		// todo, pretty format columns
		fmt.Println(app.ID, app.Version, app.User)
	}
}

type Application struct {
	ID              string              `json:"id"`
	Cmd             string              `json:"cmd,omitempty"`
	Args            []string            `json:"args,omitempty"`
	Constraints     [][]string          `json:"constraints,omitempty"`
	Container       *Container          `json:"container,omitempty"`
	CPUs            float64             `json:"cpus,omitempty"`
	Disk            float64             `json:"disk,omitempty"`
	Env             map[string]string   `json:"env,omitempty"`
	Executor        string              `json:"executor,omitempty"`
	HealthChecks    []*HealthCheck      `json:"healthChecks,omitempty"`
	Instances       int                 `json:"instances,omitemptys"`
	Mem             float64             `json:"mem,omitempty"`
	Tasks           []*Task             `json:"tasks,omitempty"`
	Ports           []int               `json:"ports,omitempty"`
	RequirePorts    bool                `json:"requirePorts,omitempty"`
	BackoffFactor   float64             `json:"backoffFactor,omitempty"`
	DeploymentID    []map[string]string `json:"deployments,omitempty"`
	Dependencies    []string            `json:"dependencies,omitempty"`
	TasksRunning    int                 `json:"tasksRunning,omitempty"`
	TasksStaged     int                 `json:"tasksStaged,omitempty"`
	User            string              `json:"user,omitempty"`
	UpgradeStrategy *UpgradeStrategy    `json:"upgradeStrategy,omitempty"`
	Uris            []string            `json:"uris,omitempty"`
	Version         string              `json:"version,omitempty"`
}

type Container struct {
	Type    string    `json:"type,omitempty"`
	Docker  *Docker   `json:"docker,omitempty"`
	Volumes []*Volume `json:"volumes,omitempty"`
}

type PortMapping struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort"`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol"`
}

type Volume struct {
	ContainerPath string `json:"containerPath,omitempty"`
	HostPath      string `json:"hostPath,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

type Docker struct {
	Image        string         `json:"image,omitempty"`
	Network      string         `json:"network,omitempty"`
	PortMappings []*PortMapping `json:"portMappings,omitempty"`
}

type UpgradeStrategy struct {
	MinimumHealthCapacity float64 `json:"minimumHealthCapacity,omitempty"`
}

type HealthCheck struct {
	Protocol               string `json:"protocol,omitempty"`
	Path                   string `json:"path,omitempty"`
	GracePeriodSeconds     int    `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds        int    `json:"intervalSeconds,omitempty"`
	PortIndex              int    `json:"portIndex,omitempty"`
	MaxConsecutiveFailures int    `json:"maxConsecutiveFailures,omitempty"`
	TimeoutSeconds         int    `json:"timeoutSeconds,omitempty"`
}

type HealthCheckResult struct {
	Alive               bool   `json:"alive"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
	FirstSuccess        string `json:"firstSuccess"`
	LastFailure         string `json:"lastFailure"`
	LastSuccess         string `json:"lastSuccess"`
	TaskID              string `json:"taskId"`
}

type Task struct {
	AppID             string               `json:"appId"`
	Host              string               `json:"host"`
	ID                string               `json:"id"`
	HealthCheckResult []*HealthCheckResult `json:"healthCheckResults"`
	Ports             []int                `json:"ports"`
	ServicePorts      []int                `json:"servicePorts"`
	StagedAt          string               `json:"stagedAt"`
	StartedAt         string               `json:"startedAt"`
	Version           string               `json:"version"`
}
