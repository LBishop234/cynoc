package packet

import "fmt"

func packetID(trafficFlowID, packetIndex string) string {
	return fmt.Sprintf("%s-%s", trafficFlowID, packetIndex)
}

func flitID(trafficFlowID, packetIndex string, flitIndex int) string {
	return fmt.Sprintf("%s-%d", packetID(trafficFlowID, packetIndex), flitIndex)
}
