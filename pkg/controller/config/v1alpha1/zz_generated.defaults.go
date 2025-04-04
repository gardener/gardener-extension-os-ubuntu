//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by defaulter-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&ExtensionConfig{}, func(obj interface{}) { SetObjectDefaults_ExtensionConfig(obj.(*ExtensionConfig)) })
	return nil
}

func SetObjectDefaults_ExtensionConfig(in *ExtensionConfig) {
	SetDefaults_ExtensionConfig(in)
	if in.NTP != nil {
		SetDefaults_NTPConfig(in.NTP)
		if in.NTP.NTPD != nil {
			SetDefaults_NTPDConfig(in.NTP.NTPD)
		}
	}
}
