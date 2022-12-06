package grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/grpc"

	pb "aprokhorov-praktikum/internal/server/grpc/proto"
	"aprokhorov-praktikum/internal/storage"
)

const (
	counter = "counter"
	gauge   = "gauge"
	base    = 10
	bitSize = 64
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage storage.Storage
}

func RegisterMetricsServer(server *grpc.Server, storage storage.Storage) {
	pb.RegisterMetricsServer(server, &MetricsServer{storage: storage})
}

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse

	metric, err := s.storage.Read(ctx, in.Type, in.Name)
	if err != nil {
		response.Error = fmt.Sprintf("Error: %s", err.Error())
	}

	response.Metric.Type = in.Type
	response.Metric.Name = in.Name

	switch data := metric.(type) {
	case storage.Counter:
		respond := int64(data)
		response.Metric.Delta = respond
	case storage.Gauge:
		respond := float64(data)
		response.Metric.Value = respond
	default:
		response.Error = "Error: Unknown metric type"
	}
	return &response, nil
}

func (s *MetricsServer) UpdateMetric(
	ctx context.Context,
	in *pb.UpdateMetricRequest,
) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse

	switch in.Metric.Type {
	case counter:
		value, err := s.storage.Read(ctx, in.Metric.Type, in.Metric.Name)
		if err != nil {
			// Если метрика не найдена, то устанавливаем счетчик в ноль
			value = storage.Counter(0)
		}

		newValue := in.Metric.Delta

		oldValue, ok := value.(storage.Counter)
		if !ok {
			return &response, errors.New("cannot make assertion (storage.Counter)")
		}

		resultValue := oldValue + storage.Counter(newValue)

		err = s.storage.Write(ctx, in.Metric.Name, resultValue)
		if err != nil {
			return &response, errors.New("cannot make assertion (storage.Counter)")
		}

	case gauge:
		newValue := in.Metric.Value

		err := s.storage.Write(ctx, in.Metric.Name, storage.Gauge(newValue))
		if err != nil {
			return &response, err
		}

	default:
		return &response, errors.New("not implemented metric type")
	}

	return &response, nil
}

func (s *MetricsServer) GetMetricsAll(
	ctx context.Context,
	in *pb.GetMetricsAllRequest,
) (*pb.GetMetricsAllResponse, error) {
	var response pb.GetMetricsAllResponse

	allMetrics, err := s.storage.ReadAll(ctx)
	if err != nil {
		return &response, err
	}

	for metricType, metrics := range allMetrics {
		switch metricType {
		case counter:
			for name, value := range metrics {
				responseValue, err := strconv.ParseInt(value, base, bitSize)
				if err != nil {
					return &response, err
				}
				response.Metrics = append(response.Metrics, &pb.Metric{
					Type:  counter,
					Name:  name,
					Delta: responseValue,
				})
			}
		case gauge:
			for name, value := range metrics {
				responseValue, err := strconv.ParseFloat(value, bitSize)
				if err != nil {
					return &response, err
				}
				response.Metrics = append(response.Metrics, &pb.Metric{
					Type:  counter,
					Name:  name,
					Value: responseValue,
				})
			}
		default:
			return &response, errors.New("error: storage return unknown metric type")
		}
	}

	return &response, nil
}
