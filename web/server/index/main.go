package index

import "github.com/jchavannes/jgo/web"

func GetRoutes() []web.Route {
	return []web.Route{
		indexRoute,
		guidesRoute,
		protocolRoute,
		disclaimerRoute,
		introducingMemoRoute,
		openSourcingMemoRoute,
		aboutRoute,
		activityRoute,
		needFundsRoute,
		newPostsRoute,
		statsRoute,
		chartsRoute,
		allRoute,
	}
}
