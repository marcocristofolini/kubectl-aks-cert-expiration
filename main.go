package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

type config struct {
	authorizer autorest.Authorizer
	subClient  subscriptions.Client
	aksClient  containerservice.ManagedClustersClient
}

var (
	clientID     string
	clientSecret string
	tenantID     string
)

func main() {
	cmd := &cobra.Command{
		Use:   "kubectl-aks-cert-expiration",
		Short: "Check AKS certificate expiration on all subscriptions",
		Run: func(cmd *cobra.Command, args []string) {
			flags := pflag.NewFlagSet("kubectl-aks-cert-expiration", pflag.ExitOnError)
			pflag.CommandLine = flags

			klog.InitFlags(nil)
			cmd.Flags().AddGoFlagSet(flag.CommandLine)

			ctx := context.Background()
			cfg := newConfig()

			// Retrieve all subscriptions
			subs, err := cfg.subClient.List(ctx)
			if err != nil {
				klog.Fatalf("Failed to list subscriptions: %v", err)
			}

			// Iterate through each subscription and check AKS certificate expiration
			for _, sub := range subs.Values() {
				subID := *sub.SubscriptionID
				klog.Infof("Checking subscription: %s", subID)

				cfg.aksClient = containerservice.NewManagedClustersClient(subID)
				cfg.aksClient.Authorizer = cfg.authorizer

				// List all AKS clusters in the subscription
				clusters, err := cfg.aksClient.ListComplete(ctx)
				if err != nil {
					klog.Fatalf("Failed to list AKS clusters for subscription %s: %v", subID, err)
				}

				// Check the certificate expiration for each AKS cluster
				for clusters.NotDone() {
					cluster := clusters.Value()
					apiServerURL := *cluster.ManagedClusterProperties.Fqdn
					checkCertificateExpiration(apiServerURL)
					clusters.NextWithContext(ctx)
				}
			}
		},
	}

	cmd.Flags().StringVar(&clientID, "client-id", "", "Azure Client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "Azure Client Secret")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Azure Tenant ID")

	if err := cmd.Execute(); err != nil {
		klog.Fatalf("Error executing command: %v", err)
	}
}

func newConfig() *config {
	cfg := &config{}

	// Create an authorizer using provided credentials
	credConfig := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
	authorizer, err := credConfig.Authorizer()
	if err != nil {
		klog.Fatalf("Failed to create authorizer: %v", err)
	}

	cfg.authorizer = authorizer

	// Create a subscriptions client
	cfg.subClient = subscriptions.NewClient()
	cfg.subClient.Authorizer = cfg.authorizer

	return cfg
}

func checkCertificateExpiration(apiServerURL string) {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", apiServerURL+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		klog.Errorf("Failed to connect to API server %s: %v", apiServerURL, err)
		return
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		serverCert := certs[0]
		expiration := serverCert.NotAfter
		now := time.Now()
		timeRemaining := expiration.Sub(now)

		fmt.Printf("Cluster API Server: %s - Certificate expiration: %s\n", apiServerURL, expiration.Format(time.RFC3339))
		fmt.Printf("Time remaining: %v\n", timeRemaining)
	} else {
		klog.Errorf("No certificates found for API server %s", apiServerURL)
	}
}
