//go:generate go-enum -f=$GOFILE --nocase

package envkit

import "strings"

// Env is an enumeration of environment namespace.
/*
ENUM(
None // Not registered namespace by default.
Development // Development environment namespace.
Staging // Staging environment namespace.
Production // Production environment namespace.
Staging-CS // Staging for CS environment namespace.
Staging-RT // Staging for RT environment namespace.
Staging-ST // Staging for ST environment namespace.
Staging-TT // Staging for TT environment namespace.
)
*/
type Env int32

func (x Env) LowerString() string {
	return strings.ToLower(x.String())
}

type Config struct {
	EnvNamespace string
	PodName      string
	ServiceName  string
	ProjectName  string
	JwtSecret    string
}
