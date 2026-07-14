// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package fixtures

import _ "embed"

//go:embed UserDataBasic.yaml
var UserDataBasic string

//go:embed UserDataDisabledUnattendedUpgrades.yaml
var UserDataDisabledUnattendedUpgrades string

//go:embed UserDataCustomAptConfig.yaml
var UserDataCustomAptConfig string

//go:embed UserDataDockerRepo.yaml
var UserDataDockerRepo string

//go:embed UserDataCustomMirror.yaml
var UserDataCustomMirror string

//go:embed UserDataGpgKey.yaml
var UserDataGpgKey string

//go:embed UserDataPinnedDeps.yaml
var UserDataPinnedDeps string

//go:embed UserDataPinnedDepsWithFallback.yaml
var UserDataPinnedDepsWithFallback string

//go:embed UserDataDisabledDependency.yaml
var UserDataDisabledDependency string
