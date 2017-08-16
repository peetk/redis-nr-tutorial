package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
}

const (
	integrationName    = "com.NewRelic.redis"
	integrationVersion = "0.1.0"
)

var args argumentList

func populateInventory(inventory sdk.Inventory) error {
	// Insert here the logic of your integration to get the inventory data
	// Ex: inventory.SetItem("softwareVersion", "value", "1.0.1")
	// --
	return nil
}

func populateMetrics(ms *metric.MetricSet) error {
	cmd := exec.Command("/bin/sh", "-c", "redis-cli info | grep instantaneous_ops_per_sec:")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	splitedLine := strings.Split(string(output), ":")
	if len(splitedLine) != 2 {
		return fmt.Errorf("Cannot split the output line")
	}
	metricValue, err := strconv.ParseFloat(strings.TrimSpace(splitedLine[1]), 64)
	if err != nil {
		return err
	}
	ms.SetMetric("query.instantaneousOpsPerSecond", metricValue, metric.GAUGE)

	cmd = exec.Command("/bin/sh", "-c", "redis-cli info | grep total_connections_received:")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	splitedLine = strings.Split(string(output), ":")
	if len(splitedLine) != 2 {
		return fmt.Errorf("Cannot split the output line")
	}
	metricValue, err = strconv.ParseFloat(strings.TrimSpace(splitedLine[1]), 64)
	if err != nil {
		return err
	}

	ms.SetMetric("net.connectionsReceivedPerSecond", metricValue, metric.RATE)

	return nil
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	if args.All || args.Inventory {
		fatalIfErr(populateInventory(integration.Inventory))
	}

	if args.All || args.Metrics {
		ms := integration.NewMetricSet("NrRedisSample")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
