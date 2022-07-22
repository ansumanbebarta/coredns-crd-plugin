package main

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
	_ "github.com/AbsaOSS/k8s_crd"
	"github.com/AbsaOSS/k8s_crd/common/directives"
	_ "github.com/coredns/coredns/core/plugin"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	p := directives.NewDirectivesManager(dnsserver.Directives)
	_ = p.InsertBefore("k8s_crd", "kubernetes")
	p.Remove("kubernetes")
	p.Remove("k8s_external")
	dnsserver.Directives = p.Get()
}

func main() {
	coremain.Run()
}
