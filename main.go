// package main

// import (
// 	"log"
// 	"net/http"
// 	"sarah/api"
// )

// func main() {
// 	http.HandleFunc("/calls/create", api.CreateCall)
// 	http.HandleFunc("/calls/list", api.ListCalls)
// 	http.HandleFunc("/calls/call/", api.GetCall)

// 	http.HandleFunc("/campaigns/create", api.CreateCampaign)
// 	http.HandleFunc("/campaigns/org/", api.GetCampaignViaOrgID)
// 	http.HandleFunc("/campaigns/campaign/", api.GetCampaignViaCampaignID)

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

package main

import (
	"fmt"
	mongo "sarah/mongodb"
)

func main() {
	orgId := "org_2zvQt8zx3QQdJPGVHYmqBajnnK1"

	campaigns := mongo.GetCampaignByOrgId(orgId)
	contacts := mongo.GetContactByOrgId(orgId)
	assistants := mongo.GetOrganizationAssistants(orgId)
	phoneNumbers := mongo.GetPhoneNumberByOrgId(orgId)

	fmt.Println(campaigns)
	fmt.Println(contacts)
	fmt.Println(assistants)
	fmt.Println(phoneNumbers)
}
