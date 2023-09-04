package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/object-storage-proxy/pkg/env"
	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
	"github.com/dimaskiddo/object-storage-proxy/pkg/proxy"
)

type ServerOptions struct {
	ListenAddress string
	TLSCertFile   string
	TLSKeyFile    string
	Insecure      bool
}

// Proxy Variable Structure
var Proxy = &cobra.Command{
	Use:   "proxy",
	Short: "Start Object Storage Proxy",
	Long:  "Start Object Storage Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var osp http.Handler
		var ospOpts proxy.ObjectStorageProxy
		var svrOpts ServerOptions

		svrOpts.ListenAddress, err = env.GetEnvString("OBJECT_STORAGE_PROXY_LISTEN_ADDRESS")
		if err != nil {
			svrOpts.ListenAddress, err = cmd.Flags().GetString("listen-address")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		svrOpts.TLSCertFile, err = env.GetEnvString("OBJECT_STORAGE_PROXY_TLS_CERT_FILE")
		if err != nil {
			svrOpts.TLSCertFile, err = cmd.Flags().GetString("tls-cert-file")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		svrOpts.TLSKeyFile, err = env.GetEnvString("OBJECT_STORAGE_PROXY_TLS_KEY_FILE")
		if err != nil {
			svrOpts.TLSKeyFile, err = cmd.Flags().GetString("tls-key-file")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		ospOpts.Endpoint, err = env.GetEnvString("OBJECT_STORAGE_PROXY_ENDPOINT")
		if err != nil {
			ospOpts.Endpoint, err = cmd.Flags().GetString("endpoint")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if ospOpts.Endpoint == "" {
					log.Println(log.LogLevelFatal, "Object Storage Endpoint is Required!")
				}
			}
		}

		ospOpts.AccessKey, err = env.GetEnvString("OBJECT_STORAGE_PROXY_ACCESS_KEY")
		if err != nil {
			ospOpts.AccessKey, err = cmd.Flags().GetString("access-key")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if ospOpts.AccessKey == "" {
					log.Println(log.LogLevelFatal, "Object Storage Access Key is Required!")
				}
			}
		}

		ospOpts.SecretKey, err = env.GetEnvString("OBJECT_STORAGE_PROXY_SECRET_KEY")
		if err != nil {
			ospOpts.SecretKey, err = cmd.Flags().GetString("secret-key")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if ospOpts.SecretKey == "" {
					log.Println(log.LogLevelFatal, "Object Storage Secret Key is Required!")
				}
			}
		}

		ospOpts.Region, err = env.GetEnvString("OBJECT_STORAGE_PROXY_REGION")
		if err != nil {
			ospOpts.Region, err = cmd.Flags().GetString("region")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		ospOpts.UpstreamStyle, err = env.GetEnvString("OBJECT_STORAGE_PROXY_UPSTREAM_STYLE")
		if err != nil {
			ospOpts.UpstreamStyle, err = cmd.Flags().GetString("upstream-style")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		ospOpts.LocalStyle, err = env.GetEnvString("OBJECT_STORAGE_PROXY_LOCAL_STYLE")
		if err != nil {
			ospOpts.LocalStyle, err = cmd.Flags().GetString("local-style")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		svrOpts.Insecure, err = env.GetEnvBool("OBJECT_STORAGE_PROXY_INSECURE")
		if err != nil {
			svrOpts.Insecure, err = cmd.Flags().GetBool("insecure")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		ospOpts.Verbose, err = env.GetEnvBool("OBJECT_STORAGE_PROXY_VERBOSE")
		if err != nil {
			ospOpts.Verbose, err = cmd.Flags().GetBool("verbose")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		ospOpts.Scheme = "https"
		if svrOpts.Insecure {
			ospOpts.Scheme = "http"
		}

		osp, err = proxy.NewObjectStorageProxy(ospOpts)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		log.Println(log.LogLevelInfo, "Starting Object Storage Proxy")

		if ospOpts.Verbose {
			log.Println(log.LogLevelInfo, "Object Storage Proxy Endpoint          : "+ospOpts.Scheme+"://"+ospOpts.Endpoint)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Access Key        : "+ospOpts.AccessKey)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Secret Key        : "+ospOpts.SecretKey)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Region            : "+ospOpts.Region)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Upstream Style    : "+ospOpts.UpstreamStyle)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Local Style       : "+ospOpts.LocalStyle)
		}

		if len(svrOpts.TLSCertFile) > 0 && len(svrOpts.TLSKeyFile) > 0 {
			log.Println(log.LogLevelInfo, "Object Storage Proxy is Listening on "+svrOpts.ListenAddress+" with HTTPS protocol")
			err = http.ListenAndServeTLS(svrOpts.ListenAddress, svrOpts.TLSCertFile, svrOpts.TLSKeyFile, osp)
		} else {
			log.Println(log.LogLevelInfo, "Object Storage Proxy is Listening on "+svrOpts.ListenAddress+" with HTTP protocol")
			err = http.ListenAndServe(svrOpts.ListenAddress, osp)
		}

		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}
	},
}

func init() {
	Proxy.Flags().String("listen-address", "0.0.0.0:9000", "Object Storage Proxy Listen Address (Default to 0.0.0.0:9000)")
	Proxy.Flags().String("tls-cert-file", "", "Object Storage Proxy Server TLS Certificate File")
	Proxy.Flags().String("tls-key-file", "", "Object Storage Proxy Server TLS Key File")
	Proxy.Flags().String("endpoint", "", "Object Storage Endpoint")
	Proxy.Flags().String("access-key", "", "Object Storage Credentials Access Key")
	Proxy.Flags().String("secret-key", "", "Object Storage Credentials Secret Key")
	Proxy.Flags().String("region", "us-east-1", "Object Storage Region (Default to us-east-1)")
	Proxy.Flags().String("upstream-style", "path", "Object Storage Upstream Access Style 'virtual' or 'path' (Default to path)")
	Proxy.Flags().String("local-style", "path", "Object Storage Local Access Style 'virtual' or 'path' (Default to path)")
	Proxy.Flags().Bool("insecure", false, "Object Storage Use Insecure Protocol (Default to false)")
	Proxy.Flags().Bool("verbose", false, "Activate Verbose Logging (Default to false)")
}
