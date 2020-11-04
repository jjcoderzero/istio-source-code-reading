package config

// Config struct defines the Istio CNI installation options
type Config struct {
	CNINetDir        string
	MountedCNINetDir string
	CNIConfName      string
	ChainedCNIPlugin bool

	CNINetworkConfigFile string
	CNINetworkConfig     string

	LogLevel           string
	KubeconfigFilename string
	KubeconfigMode     int
	KubeCAFile         string
	SkipTLSVerify      bool

	K8sServiceProtocol string
	K8sServiceHost     string
	K8sServicePort     string
	K8sNodeName        string

	UpdateCNIBinaries bool
	SkipCNIBinaries   []string
}
