package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/object-storage-proxy/pkg/env"
	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
	"github.com/dimaskiddo/object-storage-proxy/pkg/proxy"
)

type svrOptions struct {
	ListenAddress string
	TLSCertFile   string
	TLSKeyFile    string
	Scheme        string
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Region        string
	UpstreamStyle string
	LocalStyle    string
	Insecure      bool
	Verbose       bool
}

// Proxy Variable Structure
var Proxy = &cobra.Command{
	Use:   "proxy",
	Short: "Start Object Storage Proxy",
	Long:  "Start Object Storage Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var opts svrOptions
		var svc http.Handler

		opts.ListenAddress, err = env.GetEnvString("OBJECT_STORAGE_PROXY_LISTEN_ADDRESS")
		if err != nil {
			opts.ListenAddress, err = cmd.Flags().GetString("listen-address")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.TLSCertFile, err = env.GetEnvString("OBJECT_STORAGE_PROXY_TLS_CERT_FILE")
		if err != nil {
			opts.TLSCertFile, err = cmd.Flags().GetString("tls-cert-file")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.TLSKeyFile, err = env.GetEnvString("OBJECT_STORAGE_PROXY_TLS_KEY_FILE")
		if err != nil {
			opts.TLSKeyFile, err = cmd.Flags().GetString("tls-key-file")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.Endpoint, err = env.GetEnvString("OBJECT_STORAGE_PROXY_ENDPOINT")
		if err != nil {
			opts.Endpoint, err = cmd.Flags().GetString("endpoint")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if opts.Endpoint == "" {
					log.Println(log.LogLevelFatal, "Object Storage Endpoint is Required!")
				}
			}
		}

		opts.AccessKey, err = env.GetEnvString("OBJECT_STORAGE_PROXY_ACCESS_KEY")
		if err != nil {
			opts.AccessKey, err = cmd.Flags().GetString("access-key")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if opts.AccessKey == "" {
					log.Println(log.LogLevelFatal, "Object Storage Access Key is Required!")
				}
			}
		}

		opts.SecretKey, err = env.GetEnvString("OBJECT_STORAGE_PROXY_SECRET_KEY")
		if err != nil {
			opts.SecretKey, err = cmd.Flags().GetString("secret-key")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			} else {
				if opts.SecretKey == "" {
					log.Println(log.LogLevelFatal, "Object Storage Secret Key is Required!")
				}
			}
		}

		opts.Region, err = env.GetEnvString("OBJECT_STORAGE_PROXY_REGION")
		if err != nil {
			opts.Region, err = cmd.Flags().GetString("region")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.UpstreamStyle, err = env.GetEnvString("OBJECT_STORAGE_PROXY_UPSTREAM_STYLE")
		if err != nil {
			opts.UpstreamStyle, err = cmd.Flags().GetString("upstream-style")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.LocalStyle, err = env.GetEnvString("OBJECT_STORAGE_PROXY_LOCAL_STYLE")
		if err != nil {
			opts.LocalStyle, err = cmd.Flags().GetString("local-style")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.Insecure, err = env.GetEnvBool("OBJECT_STORAGE_PROXY_INSECURE")
		if err != nil {
			opts.Insecure, err = cmd.Flags().GetBool("insecure")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.Verbose, err = env.GetEnvBool("OBJECT_STORAGE_PROXY_VERBOSE")
		if err != nil {
			opts.Verbose, err = cmd.Flags().GetBool("verbose")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		opts.Scheme = "https"
		if opts.Insecure {
			opts.Scheme = "http"
		}

		svc, err = proxy.NewObjectStorageProxy(opts.Scheme, opts.Endpoint, opts.AccessKey, opts.SecretKey, opts.Region, opts.UpstreamStyle, opts.LocalStyle, opts.Verbose)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		log.Println(log.LogLevelInfo, "Starting Object Storage Proxy")

		if opts.Verbose {
			log.Println(log.LogLevelInfo, "Object Storage Proxy Endpoint          : "+opts.Scheme+"://"+opts.Endpoint)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Access Key        : "+opts.AccessKey)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Secret Key        : "+opts.SecretKey)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Region            : "+opts.Region)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Upstream Style    : "+opts.UpstreamStyle)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Local Style       : "+opts.LocalStyle)
		}

		if len(opts.TLSCertFile) > 0 && len(opts.TLSKeyFile) > 0 {
			log.Println(log.LogLevelInfo, "Object Storage Proxy is Listening on "+opts.ListenAddress+" with HTTPS protocol")
			err = http.ListenAndServeTLS(opts.ListenAddress, opts.TLSCertFile, opts.TLSKeyFile, svc)
		} else {
			log.Println(log.LogLevelInfo, "Object Storage Proxy is Listening on "+opts.ListenAddress+" with HTTP protocol")
			err = http.ListenAndServe(opts.ListenAddress, svc)
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
