// +build !ignore_autogenerated

/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	core "github.com/gardener/gardener/pkg/apis/core"
	corev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	operations "github.com/gardener/gardener/pkg/apis/operations"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Bastion)(nil), (*operations.Bastion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Bastion_To_operations_Bastion(a.(*Bastion), b.(*operations.Bastion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*operations.Bastion)(nil), (*Bastion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_operations_Bastion_To_v1alpha1_Bastion(a.(*operations.Bastion), b.(*Bastion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BastionIngressPolicy)(nil), (*operations.BastionIngressPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_BastionIngressPolicy_To_operations_BastionIngressPolicy(a.(*BastionIngressPolicy), b.(*operations.BastionIngressPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*operations.BastionIngressPolicy)(nil), (*BastionIngressPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_operations_BastionIngressPolicy_To_v1alpha1_BastionIngressPolicy(a.(*operations.BastionIngressPolicy), b.(*BastionIngressPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BastionList)(nil), (*operations.BastionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_BastionList_To_operations_BastionList(a.(*BastionList), b.(*operations.BastionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*operations.BastionList)(nil), (*BastionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_operations_BastionList_To_v1alpha1_BastionList(a.(*operations.BastionList), b.(*BastionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BastionSpec)(nil), (*operations.BastionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_BastionSpec_To_operations_BastionSpec(a.(*BastionSpec), b.(*operations.BastionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*operations.BastionSpec)(nil), (*BastionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_operations_BastionSpec_To_v1alpha1_BastionSpec(a.(*operations.BastionSpec), b.(*BastionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BastionStatus)(nil), (*operations.BastionStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_BastionStatus_To_operations_BastionStatus(a.(*BastionStatus), b.(*operations.BastionStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*operations.BastionStatus)(nil), (*BastionStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_operations_BastionStatus_To_v1alpha1_BastionStatus(a.(*operations.BastionStatus), b.(*BastionStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Bastion_To_operations_Bastion(in *Bastion, out *operations.Bastion, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_BastionSpec_To_operations_BastionSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_BastionStatus_To_operations_BastionStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Bastion_To_operations_Bastion is an autogenerated conversion function.
func Convert_v1alpha1_Bastion_To_operations_Bastion(in *Bastion, out *operations.Bastion, s conversion.Scope) error {
	return autoConvert_v1alpha1_Bastion_To_operations_Bastion(in, out, s)
}

func autoConvert_operations_Bastion_To_v1alpha1_Bastion(in *operations.Bastion, out *Bastion, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_operations_BastionSpec_To_v1alpha1_BastionSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_operations_BastionStatus_To_v1alpha1_BastionStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_operations_Bastion_To_v1alpha1_Bastion is an autogenerated conversion function.
func Convert_operations_Bastion_To_v1alpha1_Bastion(in *operations.Bastion, out *Bastion, s conversion.Scope) error {
	return autoConvert_operations_Bastion_To_v1alpha1_Bastion(in, out, s)
}

func autoConvert_v1alpha1_BastionIngressPolicy_To_operations_BastionIngressPolicy(in *BastionIngressPolicy, out *operations.BastionIngressPolicy, s conversion.Scope) error {
	out.IPBlock = in.IPBlock
	return nil
}

// Convert_v1alpha1_BastionIngressPolicy_To_operations_BastionIngressPolicy is an autogenerated conversion function.
func Convert_v1alpha1_BastionIngressPolicy_To_operations_BastionIngressPolicy(in *BastionIngressPolicy, out *operations.BastionIngressPolicy, s conversion.Scope) error {
	return autoConvert_v1alpha1_BastionIngressPolicy_To_operations_BastionIngressPolicy(in, out, s)
}

func autoConvert_operations_BastionIngressPolicy_To_v1alpha1_BastionIngressPolicy(in *operations.BastionIngressPolicy, out *BastionIngressPolicy, s conversion.Scope) error {
	out.IPBlock = in.IPBlock
	return nil
}

// Convert_operations_BastionIngressPolicy_To_v1alpha1_BastionIngressPolicy is an autogenerated conversion function.
func Convert_operations_BastionIngressPolicy_To_v1alpha1_BastionIngressPolicy(in *operations.BastionIngressPolicy, out *BastionIngressPolicy, s conversion.Scope) error {
	return autoConvert_operations_BastionIngressPolicy_To_v1alpha1_BastionIngressPolicy(in, out, s)
}

func autoConvert_v1alpha1_BastionList_To_operations_BastionList(in *BastionList, out *operations.BastionList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]operations.Bastion)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_BastionList_To_operations_BastionList is an autogenerated conversion function.
func Convert_v1alpha1_BastionList_To_operations_BastionList(in *BastionList, out *operations.BastionList, s conversion.Scope) error {
	return autoConvert_v1alpha1_BastionList_To_operations_BastionList(in, out, s)
}

func autoConvert_operations_BastionList_To_v1alpha1_BastionList(in *operations.BastionList, out *BastionList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Bastion)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_operations_BastionList_To_v1alpha1_BastionList is an autogenerated conversion function.
func Convert_operations_BastionList_To_v1alpha1_BastionList(in *operations.BastionList, out *BastionList, s conversion.Scope) error {
	return autoConvert_operations_BastionList_To_v1alpha1_BastionList(in, out, s)
}

func autoConvert_v1alpha1_BastionSpec_To_operations_BastionSpec(in *BastionSpec, out *operations.BastionSpec, s conversion.Scope) error {
	out.ShootRef = in.ShootRef
	out.SeedName = (*string)(unsafe.Pointer(in.SeedName))
	out.ProviderType = (*string)(unsafe.Pointer(in.ProviderType))
	out.SSHPublicKey = in.SSHPublicKey
	out.Ingress = *(*[]operations.BastionIngressPolicy)(unsafe.Pointer(&in.Ingress))
	return nil
}

// Convert_v1alpha1_BastionSpec_To_operations_BastionSpec is an autogenerated conversion function.
func Convert_v1alpha1_BastionSpec_To_operations_BastionSpec(in *BastionSpec, out *operations.BastionSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_BastionSpec_To_operations_BastionSpec(in, out, s)
}

func autoConvert_operations_BastionSpec_To_v1alpha1_BastionSpec(in *operations.BastionSpec, out *BastionSpec, s conversion.Scope) error {
	out.ShootRef = in.ShootRef
	out.SeedName = (*string)(unsafe.Pointer(in.SeedName))
	out.ProviderType = (*string)(unsafe.Pointer(in.ProviderType))
	out.SSHPublicKey = in.SSHPublicKey
	out.Ingress = *(*[]BastionIngressPolicy)(unsafe.Pointer(&in.Ingress))
	return nil
}

// Convert_operations_BastionSpec_To_v1alpha1_BastionSpec is an autogenerated conversion function.
func Convert_operations_BastionSpec_To_v1alpha1_BastionSpec(in *operations.BastionSpec, out *BastionSpec, s conversion.Scope) error {
	return autoConvert_operations_BastionSpec_To_v1alpha1_BastionSpec(in, out, s)
}

func autoConvert_v1alpha1_BastionStatus_To_operations_BastionStatus(in *BastionStatus, out *operations.BastionStatus, s conversion.Scope) error {
	out.Ingress = (*v1.LoadBalancerIngress)(unsafe.Pointer(in.Ingress))
	out.Conditions = *(*[]core.Condition)(unsafe.Pointer(&in.Conditions))
	out.LastHeartbeatTimestamp = (*metav1.Time)(unsafe.Pointer(in.LastHeartbeatTimestamp))
	out.ExpirationTimestamp = (*metav1.Time)(unsafe.Pointer(in.ExpirationTimestamp))
	out.ObservedGeneration = (*int64)(unsafe.Pointer(in.ObservedGeneration))
	return nil
}

// Convert_v1alpha1_BastionStatus_To_operations_BastionStatus is an autogenerated conversion function.
func Convert_v1alpha1_BastionStatus_To_operations_BastionStatus(in *BastionStatus, out *operations.BastionStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_BastionStatus_To_operations_BastionStatus(in, out, s)
}

func autoConvert_operations_BastionStatus_To_v1alpha1_BastionStatus(in *operations.BastionStatus, out *BastionStatus, s conversion.Scope) error {
	out.Ingress = (*v1.LoadBalancerIngress)(unsafe.Pointer(in.Ingress))
	out.Conditions = *(*[]corev1alpha1.Condition)(unsafe.Pointer(&in.Conditions))
	out.LastHeartbeatTimestamp = (*metav1.Time)(unsafe.Pointer(in.LastHeartbeatTimestamp))
	out.ExpirationTimestamp = (*metav1.Time)(unsafe.Pointer(in.ExpirationTimestamp))
	out.ObservedGeneration = (*int64)(unsafe.Pointer(in.ObservedGeneration))
	return nil
}

// Convert_operations_BastionStatus_To_v1alpha1_BastionStatus is an autogenerated conversion function.
func Convert_operations_BastionStatus_To_v1alpha1_BastionStatus(in *operations.BastionStatus, out *BastionStatus, s conversion.Scope) error {
	return autoConvert_operations_BastionStatus_To_v1alpha1_BastionStatus(in, out, s)
}
