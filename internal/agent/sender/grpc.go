package sender

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"aprokhorov-praktikum/internal/ccrypto"
	pb "aprokhorov-praktikum/internal/server/grpc/proto"
)

const (
	counter = "counter"
	gauge   = "gauge"
	base    = 10
	bitSize = 64
)

type GRPCAgentSender struct {
	Conn   *grpc.ClientConn
	Client pb.MetricsClient
}

func NewGRPCAgentSender(address string) (*GRPCAgentSender, error) {
	var err error
	grpcClient := GRPCAgentSender{}

	grpcClient.Conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &GRPCAgentSender{}, nil
	}

	grpcClient.Client = pb.NewMetricsClient(grpcClient.Conn)
	return &grpcClient, nil
}

func (s *GRPCAgentSender) Close() {
	s.Conn.Close()
}

func (s *GRPCAgentSender) SendMetricSingle(
	ctx context.Context,
	mtype string,
	name string,
	value string,
	hashKey string,
	cryptoKey *ccrypto.PublicKey,
) error {
	var metric pb.Metric
	switch mtype {
	case counter:
		intValue, err := strconv.ParseInt(value, base, bitSize)
		if err != nil {
			return err
		}
		metric = pb.Metric{
			Type:  mtype,
			Name:  name,
			Delta: intValue,
		}

	case gauge:
		floatValue, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return err
		}
		metric = pb.Metric{
			Type:  mtype,
			Name:  name,
			Value: floatValue,
		}
	default:
		return errors.New("GRPC Sender Error: Unknown metric type")
	}

	resp, err := s.Client.UpdateMetric(ctx, &pb.UpdateMetricRequest{Metric: &metric})
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

func (s *GRPCAgentSender) SendMetricBatch(
	ctx context.Context,
	metrics map[string]map[string]string,
	hashKey string,
	cryptoKey *ccrypto.PublicKey,
) error {
	// To implement
	return nil
}
