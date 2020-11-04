package destinationrule

import (
	"istio.io/api/networking/v1alpha3"

	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

// CaCertificateAnalyzer checks if CaCertificate is set in case mode is SIMPLE/MUTUAL
type CaCertificateAnalyzer struct{}

var _ analysis.Analyzer = &CaCertificateAnalyzer{}

func (c *CaCertificateAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "destinationrule.CaCertificateAnalyzer",
		Description: "Checks if caCertificates is set when TLS mode is SIMPLE/MUTUAL",
		Inputs: collection.Names{
			collections.IstioNetworkingV1Alpha3Destinationrules.Name(),
		},
	}
}

func (c *CaCertificateAnalyzer) Analyze(ctx analysis.Context) {
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), func(r *resource.Instance) bool {
		c.analyzeDestinationRule(r, ctx)
		return true
	})
}

func (c *CaCertificateAnalyzer) analyzeDestinationRule(r *resource.Instance, ctx analysis.Context) {
	dr := r.Message.(*v1alpha3.DestinationRule)
	drNs := r.Metadata.FullName.Namespace
	drName := r.Metadata.FullName.String()
	mode := dr.GetTrafficPolicy().GetTls().GetMode()

	if mode == v1alpha3.ClientTLSSettings_SIMPLE || mode == v1alpha3.ClientTLSSettings_MUTUAL {
		if dr.GetTrafficPolicy().GetTls().GetCaCertificates() == "" {
			ctx.Report(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), msg.NewNoServerCertificateVerificationDestinationLevel(r, drName,
				drNs.String(), mode.String(), dr.GetHost()))
		}
	}
	portSettings := dr.TrafficPolicy.GetPortLevelSettings()

	for _, p := range portSettings {
		mode = p.GetTls().GetMode()
		if mode == v1alpha3.ClientTLSSettings_SIMPLE || mode == v1alpha3.ClientTLSSettings_MUTUAL {
			if p.GetTls().GetCaCertificates() == "" {
				ctx.Report(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), msg.NewNoServerCertificateVerificationPortLevel(r, drName,
					drNs.String(), mode.String(), dr.GetHost(), p.GetPort().String()))
			}
		}
	}
}
