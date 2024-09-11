package envkit

import (
	"errors"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

var (
	// registration
	regEnv         = EnvNone
	regPodName     = ""
	regServiceName = ""
	regProjectName = ""
	regJwtSecret   = ""
	reqEnvirement  = ""

	// registerOnce limits registering once
	registerOnce = sync.Once{}
)

func Register(cfg Config) {
	registerOnce.Do(func() {
		env, err := ParseEnv(cfg.EnvNamespace)
		if err != nil || env == EnvNone {
			panic(errors.New("no such environment namespace"))
		}

		regEnv = env
		regPodName = cfg.PodName
		regServiceName = cfg.ServiceName
		regProjectName = cfg.ProjectName
	})
}

func Registered() bool {
	return Namespace() != EnvNone
}

// ResetRegister is used in test case only.
func ResetRegister() {
	regEnv = EnvNone
	regPodName = ""
	regServiceName = ""
	regProjectName = ""
	registerOnce = sync.Once{}
}

// ServiceName named by functionality in the microservice.
// e.g. amazing-chatroom-rpc, at-notification-rpc, ...
func ServiceName() string {
	if regServiceName != "" {
		return regServiceName
	}

	return os.Getenv("ENV_SERVICE_NAME")
}

// ProjectName named by project which the microservice belongs to
// e.g. amazing-chatroom, at-notification, ...
func ProjectName() string {
	if regProjectName != "" {
		return regProjectName
	}

	return os.Getenv("ENV_PROJECT_NAME")
}

// PodName pod name in k8s.
// It's used to identify different instances with the same service name.
func PodName() string {
	if regPodName != "" {
		return regPodName
	}

	// customized pod name in the environment variable
	// - name: PODNAME
	//   valueFrom:
	//     fieldRef:
	//     fieldPath: metadata.name
	if name := os.Getenv("ENV_POD_NAME"); name != "" {
		return name
	}

	// https://stackoverflow.com/questions/58101598/kubernetes-get-the-full-pod-name-as-environment-variable
	return os.Getenv("HOSTNAME")
}

func JwtSecret() string {
	if regJwtSecret != "" {
		return regJwtSecret
	}

	return os.Getenv("ENV_JWT_SECRET")
}

// EnvNamespace environment namespace.
// e.g. Production, Staging, Development, ...
func EnvNamespace() string {
	env := Namespace()
	if env != EnvNone {
		return env.String()
	}

	return ""
}

// Namespace environment namespace
func Namespace() Env {
	if regEnv != EnvNone {
		return regEnv
	}

	env, err := ParseEnv(os.Getenv("ENV_NAMESPACE"))
	if err != nil {
		return EnvNone
	}

	return env
}

// Namespace environment namespace
func Environment() string {
	if reqEnvirement != "" {
		return reqEnvirement
	}

	return os.Getenv("ENV_ENVIRONMENT")
}

// CallerInfo presents the detail in the runtime
type CallerInfo struct {
	FullPkgName  string
	PkgName      string
	FileName     string
	FuncName     string
	FullFuncName string
	Line         int
}

// RetrieveCallerInfo returns the info related to the caller in the runtime.
// ref: https://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func RetrieveCallerInfo(options ...EnvOptions) CallerInfo {
	// load options
	o := loadEnvOptions(options...)

	pc, file, line, _ := runtime.Caller(o.skip)
	_, fileName := path.Split(file)

	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	funcName := parts[pl-1]
	fullFuncName := funcName
	if parts[pl-2][0] == '(' {
		fullFuncName = parts[pl-2] + "." + funcName
	}

	fullPkgName := ""
	if parts[pl-2][0] == '(' {
		fullPkgName = strings.Join(parts[0:pl-2], ".")
	} else {
		fullPkgName = strings.Join(parts[0:pl-1], ".")
	}

	// package name takes the latest part of full package name
	// ex: gs-mono/example/package
	//   - FullPkgName: gs-mono/example/package
	//   - PkgName: package
	ss := strings.Split(fullPkgName, "/")
	pkgName := ss[len(ss)-1]

	return CallerInfo{
		FullPkgName:  fullPkgName,
		PkgName:      pkgName,
		FileName:     fileName,
		FuncName:     funcName,
		FullFuncName: fullFuncName,
		Line:         line,
	}
}
