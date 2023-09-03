package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/object-storage-proxy/pkg/env"
	"github.com/dimaskiddo/object-storage-proxy/pkg/log"
	"github.com/dimaskiddo/object-storage-proxy/pkg/proxy"
)

type Options struct {
	ListenAddress string
	Scheme        string
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Region        string
	Insecure      bool
	UpstreamStyle string
	LocalStyle    string
	Verbose       bool
}

// Proxy Variable Structure
var Proxy = &cobra.Command{
	Use:   "proxy",
	Short: "Start Object Storage Proxy",
	Long:  "Start Object Storage Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var opts Options
		var svc http.Handler

		opts.ListenAddress, err = env.GetEnvString("OBJECT_STORAGE_PROXY_LISTEN_ADDRESS")
		if err != nil {
			opts.ListenAddress, err = cmd.Flags().GetString("listen-address")
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

		opts.Insecure, err = env.GetEnvBool("OBJECT_STORAGE_PROXY_INSECURE")
		if err != nil {
			opts.Insecure, err = cmd.Flags().GetBool("insecure")
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
			log.Println(log.LogLevelInfo, "Object Storage Proxy Local Style       : "+opts.UpstreamStyle)
			log.Println(log.LogLevelInfo, "Object Storage Proxy Local Style       : "+opts.LocalStyle)
		}

		log.Println(log.LogLevelInfo, "Object Storage Proxy is Listening on "+opts.ListenAddress+" with HTTP protocol")
		http.ListenAndServe(opts.ListenAddress, svc)
	},
}

func init() {
	Proxy.Flags().String("endpoint", "", "Object Storage Endpoint")
	Proxy.Flags().String("access-key", "", "Object Storage Credentials Access Key")
	Proxy.Flags().String("secret-key", "", "Object Storage Credentials Secret Key")
	Proxy.Flags().String("region", "us-east-1", "Object Storage Region (Default to us-east-1)")
	Proxy.Flags().Bool("insecure", false, "Object Storage Use Insecure Protocol (Default to false)")
	Proxy.Flags().String("upstream-style", "path", "Object Storage Upstream Access Style 'virtual' or 'path' (Default to path)")
	Proxy.Flags().String("local-style", "path", "Object Storage Local Access Style 'virtual' or 'path' (Default to path)")
	Proxy.Flags().String("listen-address", "0.0.0.0:9000", "Object Storage Proxy Listen Address (Default to 0.0.0.0:9000)")
	Proxy.Flags().Bool("verbose", false, "Activate Verbose Logging (Default to false)")
}
