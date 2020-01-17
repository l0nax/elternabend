package actions

import (
	"log"
	"strconv"

	"github.com/casbin/casbin"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/pop"
	suuid "github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
	"github.com/l0nax/elternabend/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator


// ======================================================
//		Global Configuration Variables
// ======================================================
var TEACHER_PASSWORD_LENGTH int

var PASSWORD_HASH_MEMORY int
var PASSWORD_HASH_ITERATIONS int
var PASSWORD_HASH_THREADS int
var PASSWORD_HASH_SALT_LEN int
var PASSWORD_HASH_KEY_LEN int

// getGlobalConf initializes all global (env) Variables.
// Panics if there was an error
func getGlobalConf() {
	intConv(&TEACHER_PASSWORD_LENGTH, "TEACHER_PASSWORD_LENGTH", "6",
		"Error while converting 'TEACHER_PASSWORD_LENGTH' to int")


	intConv(&PASSWORD_HASH_MEMORY, "PASSWORD_HASH_MEMORY", "512",
		"Error while converting 'PASSWORD_HASH_MEMORY' to int")

	intConv(&PASSWORD_HASH_ITERATIONS, "PASSWORD_HASH_ITERATIONS", "4",
		"Error while converting 'PASSWORD_HASH_ITERATIONS' to int")

	intConv(&PASSWORD_HASH_THREADS, "PASSWORD_HASH_THREADS", "4",
		"Error while converting 'PASSWORD_HASH_THREADS' to int")

	intConv(&PASSWORD_HASH_SALT_LEN, "PASSWORD_HASH_SALT_LEN", "16",
		"Error while converting 'PASSWORD_HASH_SALT_LEN' to int")

	intConv(&PASSWORD_HASH_KEY_LEN, "PASSWORD_HASH_KEY_LEN", "32",
		"Error while converting 'PASSWORD_HASH_KEY_LEN' to int")
}

// intConv converts an environment variable to an integer and puts it into @res.
// panics with error message @errMsg if conversion was not possible
func intConv(res *int, name string, defaultVal string, errMsg string) {
	var err error

	(*res), err = strconv.Atoi(envy.Get(name, defaultVal))
	if err != nil {
		panic(Wrap(err, errMsg))
	}
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_elternabend_session",
		})

		// first setup config (env) vars
		getGlobalConf()

		// setup casbin auth rules
		authEnforcer, err := casbin.NewEnforcerSafe(envy.Get("RBAC_AUTH_MODEL_PATH", "rbac_model.conf"),
			envy.Get("RBAC_POLICY_PATH", "rbac_policy.csv"))
		if err != nil {
			// TODO: Log this as Fatal
			panic(err)
		}

		// add some Variables
		app.Use(AddVariablesMiddleware)

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Check if the User is already authentificated or not
		app.Use(SetCurrentUser)

		// add the authorizer MW
		app.Use(NewRBACCheckMiddleware(authEnforcer))

		// Setup and use translations:
		app.Use(translations())

		app.GET("", HomeHandler)
		app.GET("/", HomeHandler)

		// ===> Dashboard <===
		dashboardRoute := app.Group("/d/")
		dashboardRoutes := []RouteResource{
			{
				"/",
				"GET",
				Dashboard,
			},
			{
				"/new",
				"GET",
				CreateEvent,
			},
			{
				"/new",
				"POST",
				CreateEventPost,
			},
		}

		for _, route := range dashboardRoutes {
			route.AddRoute(dashboardRoute)
		}

		// ===> Auth/ User <===
		userRoute := app.Group("/u/")
		userRoutes := []RouteResource{
			{
				Route:   "/login",
				Method:  "GET",
				Handler: NewLogin,
			},
			{
				Route:   "/login",
				Method:  "POST",
				Handler: Login,
			},
			{
				"/logout",
				"GET",
				LogOut,
			},
			{
				"/new",
				"GET",
				CreateUser,
			},
			{
				"/new",
				"POST",
				CreateUserPost,
			},
		}

		for _, route := range userRoutes {
			route.AddRoute(userRoute)
		}

		// ===> API <===
		// ==> v1 <==
		apiV1 := app.Group("/v1/")
		v1APIs := []RouteResource{
			{
				Route:   "/",
				Method:  "POST",
				Handler: RoutesHandler,
			},
		}

		for _, rrApi := range v1APIs {
			rrApi.AddRoute(apiV1)
		}

		// Teacher API
		apiV1_teacher := apiV1.Group("/teacher")
		apiV1_teachers := []RouteResource{
			{
				"/new",
				"CREATE",
				TeacherCreate,
			},
		}

		for _, rrApi := range apiV1_teachers {
			rrApi.AddRoute(apiV1_teacher)
		}

		// ===> Only for Developing <===
		routes := RouteResource{
			Route:   "/routes",
			Method:  "GET",
			Handler: RoutesHandler,
		}

		routes.AddRoute(app)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

// SetCurrentUser attempts to find a user based on the session_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uuid := c.Session().Get("session_uuid"); uuid != nil {
			log.Printf("FOUND: %s", c.Session().Get("session_uuid"))
			user := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if tx == nil {
				return errors.WithStack(errors.New("Error while getting 'tx'. 'tx' is nil"))
			}

			_uuid, err := suuid.Parse(uuid.(string))
			if err != nil {
				return errors.Wrap(errors.WithStack(err), "Error while converting string UUID into struct")
			}
			if err := tx.Where("session_id = ?", _uuid).First(user); err != nil {
				// clear cession storage of user
				c.Session().Clear()

				return errors.Wrap(errors.WithStack(err), "Error while searching UUID of user")
			}

			// // TODO: Check if this is really needed!
			c.Set("roles", user.Roles)
			c.Set("user", user)
		}

		// skip this if no valid session id was found
		return next(c)
	}
}

func AddVariablesMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// add the GO_ENV variable
		c.Set("ENV", ENV)

		return next(c)
	}
}
