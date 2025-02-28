<template>
  <div class="use-package-apps">
    <bk-select :value="''" :disabled="tplCounts === 0" :filterable="true" :input-search="false">
      <template #trigger>
        <div
            :class="['select-app-trigger', { disabled: tplCounts === 0 }]"
            v-bk-tooltips="{
              disabled: tplCounts > 0,
              content: t('该套餐中没有可用配置文件，无法被导入到服务配置中')
            }">
          <Plus class="plus-icon" />
          {{ t('新服务中使用') }}
        </div>
      </template>
      <bk-option v-for="app in unBoundApps" :key="app.id" :id="app.id" :label="app.spec.name">
        <div class="app-option" @click="goToConfigPageImport(app.id as number)">
          <div class="name-text">{{ app.spec.name }}</div>
          <LinkToApp class="link-icon" :id="app.id as number"  />
        </div>
      </bk-option>
    </bk-select>
    <div class="table-wrapper">
      <bk-loading :loading="boundAppsLoading">
        <div class="refresh-header">
          <span class="text">{{ t('当前使用此套餐的服务') }}</span>
          <right-turn-line class="refresh-button" :class="{rotate:boundAppsLoading}"  @click="getBoundApps"/>
        </div>
        <bk-table :border="['outer']" :data="boundApps" :thead="{isShow:false}" :empty-text="t('暂无数据')">
          <bk-table-column label="">
            <template #default="{ row }">
              <div v-if="row.app_id" class="app-info" @click="goToConfigPageImport(row.app_id)">
                <div v-overflow-title class="name-text">{{ row.app_name }}</div>
                <LinkToApp class="link-icon" :id="row.app_id" />
              </div>
            </template>
          </bk-table-column>
        </bk-table>
        <bk-pagination class="table-pagination" small align="center" :show-limit="false" :show-total-count="false">
        </bk-pagination>
      </bk-loading>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { Plus, RightTurnLine } from 'bkui-vue/lib/icon';
import useGlobalStore from '../../../../../store/global';
// import useUserStore from '../../../../../store/user';
import useTemplateStore from '../../../../../store/template';
import { getAppList } from '../../../../../api/index';
import { getUnNamedVersionAppsBoundByPackage } from '../../../../../api/template';
import { IAppItem } from '../../../../../../types/app';
import { IPackageCitedByApps } from '../../../../../../types/template';
import LinkToApp from '../components/link-to-app.vue';

const router = useRouter();
const { t } = useI18n();

const { spaceId } = storeToRefs(useGlobalStore());
// const { userInfo } = storeToRefs(useUserStore());
const templateStore = useTemplateStore();
const { currentTemplateSpace, currentPkg } = storeToRefs(templateStore);

defineProps<{
  tplCounts: number;
}>();

const userApps = ref<IAppItem[]>([]);
const userAppListLoading = ref(false);
const boundApps = ref<IPackageCitedByApps[]>([]);
const boundAppsLoading = ref(false);

const unBoundApps = computed(() => {
  const res = userApps.value.filter(app => boundApps.value.findIndex(item => item.app_id === app.id) === -1);
  return res;
});

watch(
  () => currentPkg.value,
  () => {
    boundApps.value = [];
    getBoundApps();
  },
);

onMounted(() => {
  getBoundApps();
  getUserApps();
});

const getUserApps = async () => {
  userAppListLoading.value = true;
  const params = {
    start: 0,
    all: true,
  };
  const res = await getAppList(spaceId.value, params);
  userApps.value = res.details;
  userAppListLoading.value = false;
};

const getBoundApps = async () => {
  if (typeof currentPkg.value !== 'number') return;
  boundAppsLoading.value = true;
  const params = {
    start: 0,
    all: true,
  };
  const res = await getUnNamedVersionAppsBoundByPackage(
    spaceId.value,
    currentTemplateSpace.value,
    currentPkg.value as number,
    params,
  );
  boundApps.value = res.details;
  boundAppsLoading.value = false;
};

const goToConfigPageImport = (id: number) => {
  const { href } = router.resolve({
    name: 'service-config',
    params: { appId: id },
    query: { pkg_id: currentPkg.value },
  });
  window.open(href, '_blank');
};
</script>
<style lang="scss" scoped>
.use-package-apps {
  padding: 16px 24px;
  width: 240px;
  height: 100%;
  background: #ffffff;
}
.select-app-trigger {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5px;
  width: 192px;
  height: 32px;
  line-height: 22px;
  border: 1px solid #c4c6cc;
  border-radius: 2px;
  color: #63656e;
  font-size: 14px;
  overflow: hidden;
  cursor: pointer;
  &.disabled {
    color: #dcdee5;
    border-color: #dcdee5;
    cursor: not-allowed;
  }
  .plus-icon {
    font-size: 20px;
  }
}
.app-option,
.app-info {
  display: flex;
  align-items: center;
  overflow: hidden;
  cursor: pointer;
}
.table-wrapper {
  margin-top: 16px;
  .app-info {
    display: flex;
    align-items: center;
    overflow: hidden;
  }
  .table-pagination {
    margin-top: 16px;
  }
  .refresh-header {
    display: flex;
    align-items: center;
    padding: 0 16px;
    font-size: 12px;
    height: 41px;
    border: 1px solid #DCDEE5;
    border-bottom: none;
    &:hover {
      background-color: #f0f1f5;
    }
    .text {
      margin-right: 16px;
    }
    .refresh-button {
      color:#3a84ff;
      font-size: 16px;
      cursor: pointer;
    }
  }
  :deep(.bk-exception-img) {
    height: 80px;
  }
}
.name-text {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.link-icon {
  flex-shrink: 0;
  margin-left: 10px;
}
.rotate {
  animation: rotate 0.5s infinite linear;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
