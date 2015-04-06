package main

import (
	"encoding/json"
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

type AppById struct {
	App Application `json:"app"`
}

type Applications struct {
	Apps []Application `json:"apps"`
}

type Tasks struct {
	Tasks []*Task `json:"tasks"`
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

type Versions struct {
	Versions []string
}

type Update struct {
	DeploymentID string `json:"deploymentId"`
	Version      string `json:"version"`
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

type QueuedTask struct {
	App   *Application    `json:"app"`
	Delay map[string]bool `json:"delay"`
}

type Queue struct {
	Queue []QueuedTask `json:"queue"`
}

type Which struct {
	Leader string `json:"leader"`
}

type Message struct {
	Message string `json:"message"`
}

type Deploys []Deploy

type Deploy struct {
	AffectedApps   []string `json:"affectedApps"`
	DeployID       string   `json:"id"`
	Steps          [][]Step `json:"steps"`
	CurrentActions []Step   `json:"currentActions"`
	Version        string   `json:"version"`
	CurrentStep    int      `json:"currentStep"`
	TotalSteps     int      `json:"totalSteps"`
}

type Step struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

type Group struct {
	GroupID      string         `json:"id"`
	Version      string         `json:"version"`
	Apps         []*Application `json:"apps"`
	Dependencies []string       `json:"dependencies"`
	Groups       []*Group       `json:"groups"`
}
