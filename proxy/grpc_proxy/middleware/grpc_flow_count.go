package middleware

import (
	"gateway/enity"
	"gateway/globals"
	"log"

	"google.golang.org/grpc"
)

func GrpcFlowCountMiddleware(serviceDetail *enity.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		totalCounter, err := globals.FlowCounter.GetCounter(globals.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		serviceCounter, err := globals.FlowCounter.GetCounter(serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
