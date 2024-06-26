package k8sctrl

/*
Copyright 2022 The k8gb Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/external-dns/endpoint"
)

type geo map[string]interface{}

type LocalDNSEndpoint struct {
	Targets []string
	TTL     endpoint.TTL
	Labels  map[string]string
	DNSName string
	RecordType string
}

func (lep LocalDNSEndpoint) String() string {
	return fmt.Sprintf("%s: %v, Targets: %v, Labels: %v", lep.DNSName, lep.TTL, lep.Targets, lep.Labels)
}

func (lep LocalDNSEndpoint) extractGeo(endpoint *endpoint.Endpoint, clientIP net.IP, geoDataFilePath string, geoDataFieldPath ...string) (result []string) {
	if geoDataFilePath == "" {
		return nil
	}

	db, err := maxminddb.Open(geoDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // nolint:errcheck

	var clientGeo geo
	err = db.Lookup(clientIP, &clientGeo)
	if err != nil {
		return nil
	}

	log.Infof("extracted client geo data: %+v", clientGeo)

	if len(geoDataFieldPath) == 0 {
		log.Info("no geo data field specified")
		return result
	}

	clientGeoData, found, err := unstructured.NestedString(clientGeo, geoDataFieldPath...)
	if err != nil {
		log.Infof("error retrieving client geo data for field %+v: %v", geoDataFieldPath, err)
		return result
	}

	if !found || clientGeoData == "" {
		log.Infof("client geo data field %+v not found", geoDataFieldPath)
		return result
	}

	log.Infof("client geo data field value for %+v: %+v", geoDataFieldPath, clientGeoData)

	for _, ip := range endpoint.Targets {
		var endpointGeo geo
		log.Infof("processing IP %+v", ip)
		err = db.Lookup(net.ParseIP(ip), &endpointGeo)
		if err != nil {
			log.Error(err)
			continue
		}
		endpointGeoData, found, err := unstructured.NestedString(endpointGeo, geoDataFieldPath...)
		if err != nil {
			log.Infof("error retrieving endpoint geo data for field %+v: %v", geoDataFieldPath, err)
			return result
		}

		if !found || endpointGeoData == "" {
			log.Infof("endpoint geo data field %+v not found", geoDataFieldPath)
			return result
		}

		log.Infof("endpoint data field value for %+v: %+v", geoDataFieldPath, clientGeoData)
		if clientGeoData == endpointGeoData {
			result = append(result, ip)
		}
	}
	return result
}
