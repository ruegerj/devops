// Contains all middleware modules used by the application
package middleware

import "net/http"

// Type signature for a module on the middleware stack
type Middleware func(http.Handler) http.Handler
