/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package utils

import (
	"strings"
)

func IsPermissionDenied(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "PermissionDenied")
}

// IsContextCanceled confirm the error whether context canceled
func IsContextCanceled(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "context canceled")
}

// IsSecretAlreadyExist defines the secret whether exist
func IsSecretAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "already exists")
}

// IsArgoResourceNotFound defines the argo resource not found
func IsArgoResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "code = NotFound")
}

// IsArgoNotFoundAsPartOf defines the resource not found as port of application
func IsArgoNotFoundAsPartOf(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not found as part of")
}
