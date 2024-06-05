package hooks

import "github.com/pocketbase/pocketbase"

var app *pocketbase.PocketBase

func Initialize(pocketbaseInstance *pocketbase.PocketBase) {
	app = pocketbaseInstance

	app.OnBeforeServe().Add(BeforeServeHook)

	// fires only for "users" collections
	app.OnRecordAfterCreateRequest("users").Add(RecordAfterCreateRequestHook)

	//fires only for "users" collections
	app.OnRecordAuthRequest("users").Add(RecordAuthRequestHook)
}
