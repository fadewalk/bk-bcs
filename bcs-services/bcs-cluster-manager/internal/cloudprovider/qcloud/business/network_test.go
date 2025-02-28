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
 */

package business

import (
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider"
)

var opt = &cloudprovider.CommonOption{}

func TestGetFreeIPNets(t *testing.T) {
	freeCidrs, err := GetFreeIPNets(opt, "vpc-xxx")
	if err != nil {
		t.Fatal(err)
	}

	for i := range freeCidrs {
		t.Log(freeCidrs[i].String())
	}
}

func TestAllocateSubnet(t *testing.T) {
	sub, err := AllocateSubnet(opt, "vpc-xxx", "ap-guangzhou-x", 24, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sub)
}
