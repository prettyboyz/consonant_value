package buyorderscontroller

import (
	"encoding/json"
	"time"

	"github.com/royvandewater/trading-post/ordersservice"
)

func formatCreateResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(toOrderResponse(order), "", "  ")
}

func formatGetResponse(order ordersservice.BuyOrder) ([]byte, error) {
	return json.MarshalIndent(toOrderResponse(order), "", "  ")
}

func formatListResponse(orders []ordersservice.BuyOrder) ([]byte, error) {
	orderResponses := make([]_OrderResponse, len(orders)