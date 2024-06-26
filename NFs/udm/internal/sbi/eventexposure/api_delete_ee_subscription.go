/*
 * Nudm_EE
 *
 * Nudm Event Exposure Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package eventexposure

import (
	"github.com/gin-gonic/gin"

	"github.com/free5gc/udm/internal/sbi/producer"
	"github.com/machi12/util/httpwrapper"
)

// DeleteEeSubscription - Unsubscribe
func HTTPDeleteEeSubscription(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueIdentity"] = c.Params.ByName("ueIdentity")
	req.Params["subscriptionID"] = c.Params.ByName("subscriptionId")

	rsp := producer.HandleDeleteEeSubscription(req)
	// only return 204 no content
	c.Status(rsp.Status)
}
