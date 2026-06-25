package res

// BuildSpec is a config spec materialized by builder.Build.
type BuildSpec interface {
	Build() (any, error)
}

// NewResourceer is a wire type whose resource is materialized in builder.Build.
type NewResourceer interface {
	NewResource() (any, error)
}

// BuildConfiger is a wire type that yields a config spec for env apply and build.
type BuildConfiger interface {
	BuildConfig() (BuildSpec, error)
}
