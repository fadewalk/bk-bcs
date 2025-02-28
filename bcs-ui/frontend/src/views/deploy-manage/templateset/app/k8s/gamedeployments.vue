<template>
  <div class="biz-content">
    <Header hide-back title="GameDeployments" />
    <div class="biz-content-wrapper" style="padding: 0;" v-bkloading="{ isLoading: isInitLoading, opacity: 0.1 }">
      <template v-if="!isInitLoading">
        <div class="biz-panel-header biz-event-query-query" style="padding-right: 0;">
          <div class="left !w-[260px]">
            <ClusterSelect
              :placeholder="$t('deploy.templateset.cluster')"
              cluster-type="all"
              :value="selectedClusterId"
              @change="handleChangeCluster" />
          </div>
          <div class="left">
            <bk-selector
              :placeholder="$t('dashboard.ns.validate.emptyNs')"
              :searchable="true"
              :allow-clear="true"
              :setting-key="'name'"
              :display-key="'name'"
              :selected.sync="selectedNamespaceName"
              :list="namespaceList"
              :is-loading="namespaceLoading"
              :search-placeholder="$t('deploy.templateset.searchNs')"
              @item-selected="handleChangeNamespace"
              @clear="handleClearNamespace">
            </bk-selector>
          </div>
          <div class="left">
            <bk-input
              v-model="searchKey"
              style="width: 240px;" clearable :placeholder="$t('generic.placeholder.searchName')" @clear="clearSearch" />
          </div>
          <div class="left">
            <bk-button type="primary" :title="$t('generic.button.query')" icon="search" @click="handleClick">
              {{$t('generic.button.query')}}
            </bk-button>
            <button class="bk-button" @click="batchDel">
              <span>{{$t('generic.button.batchDelete')}}</span>
            </button>
          </div>
        </div>
        <div v-bkloading="{ isLoading: isPageLoading && !isInitLoading }">
          <div class="biz-table-wrapper gamestatefullset-table-wrapper">
            <bk-table
              :data="curPageData"
              :page-params="pageConf"
              @page-change="handlePageChange"
              @page-limit-change="handlePageSizeChange"
              @selection-change="handleSelectionChange">
              <bk-table-column type="selection" width="60"></bk-table-column>
              <bk-table-column
                v-for="(column, index) in columnList"
                :label="(defaultColumnMap[column] && defaultColumnMap[column].label) || column"
                :min-width="defaultColumnMap[column] ? defaultColumnMap[column].minWidth : 'auto'"
                :key="index">
                <template slot-scope="{ row }">
                  <div class="cell">
                    <template v-if="column === 'name'">
                      <a
                        href="javascript:void(0);"
                        class="bk-text-button name-col bcs-ellipsis" style="font-weight: 700;"
                        @click="showDetail(row[column], row['namespace'])">{{row[column] || '--'}}</a>
                    </template>
                    <template v-else>
                      {{row[column] || '--'}}
                    </template>
                  </div>
                </template>
              </bk-table-column>
              <bk-table-column :label="$t('generic.label.action')" width="200">
                <template slot-scope="{ row }">
                  <a href="javascript:void(0);" class="bk-text-button" @click.stop="update(row, index)">{{$t('generic.button.update')}}</a>
                  <a href="javascript:void(0);" class="bk-text-button" @click.stop="scale(row, index)">{{$t('deploy.templateset.scale')}}</a>
                  <a href="javascript:void(0);" class="bk-text-button" @click.stop="del(row, index)">{{$t('generic.button.delete')}}</a>
                </template>
              </bk-table-column>
              <template #empty>
                <BcsEmptyTableStatus :type="searchKey ? 'search-empty' : 'empty'" @clear="handleClearSearchData" />
              </template>
            </bk-table>
          </div>
        </div>
      </template>
    </div>
    <!-- <gamedeployments-sideslider
      :is-show="isShowSideslider"
      :cluster-id="selectedClusterId"
      :namespace-name="curShowNamespace"
      :name="curShowName"
      @hide-sideslider="hideSideslider">
    </gamedeployments-sideslider> -->

    <gamedeployments-update
      :is-show="isShowUpdateDialog"
      :item="updateItem"
      @hide-update="hideGamedeploymentsUpdate"
      @update-success="gamedeploymentsUpdateSuccess">
    </gamedeployments-update>

    <gamedeployments-scale
      :is-show="isShowScale"
      :cluster-id="selectedClusterId"
      :item="scaleItem"
      @hide-scale="hideScale"
      @scale-success="gamedeploymentsScaleSuccess">
    </gamedeployments-scale>

    <bk-dialog
      :is-show="batchDelDialogConf.isShow"
      :width="500"
      :has-header="false"
      :quick-close="false"
      class="batch-delete-gamestatefulset"
      @cancel="hideBatchDelDialog">
      <div class="biz-batch-wrapper">
        <p class="batch-title">{{batchDelDialogConf.title}}</p>
        <ul class="batch-list">
          <li v-for="(item, index) of batchDelDialogConf.list" :key="index">{{item.name}}</li>
        </ul>
      </div>
      <template slot="footer">
        <div class="bk-dialog-outer">
          <template v-if="batchDelDialogConf.isDeleting">
            <bk-button theme="primary">
              {{$t('generic.status.deleting')}}
            </bk-button>
            <bk-button>
              {{$t('generic.button.cancel')}}
            </bk-button>
          </template>
          <template v-else>
            <bk-button theme="primary" @click="batchDelConfirm">
              {{$t('generic.button.confirm')}}
            </bk-button>
            <bk-button @click="hideBatchDelDialog">
              {{$t('generic.button.cancel')}}
            </bk-button>
          </template>
        </div>
      </template>
    </bk-dialog>
  </div>
</template>

<script>
import GamedeploymentsScale from './gamedeployments-scale';
import GamedeploymentsUpdate from './gamedeployments-update';

import { catchErrorHandler } from '@/common/util';
import ClusterSelect from '@/components/cluster-selector/cluster-select.vue';
import Header from '@/components/layout/Header.vue';
import { useNamespace } from '@/views/resource-view/namespace/use-namespace';

export default {
  name: 'GameDeployments',
  components: {
    GamedeploymentsUpdate,
    GamedeploymentsScale,
    Header,
    ClusterSelect,
  },
  data() {
    return {
      CATEGORY: 'gamedeployments',

      isInitLoading: true,
      isPageLoading: false,

      pageConf: {
        total: 0,
        totalPage: 1,
        pageSize: 10,
        curPage: 1,
        show: true,
      },
      defaultColumnMap: {
        name: {
          label: this.$t('generic.label.name'),
          minWidth: 150,
        },
        cluster_id: {
          label: this.$t('generic.label.cluster'),
          minWidth: 140,
        },
        namespace: {
          label: this.$t('k8s.namespace'),
          minWidth: 100,
        },
      },
      bkMessageInstance: null,
      selectedClusterId: '',
      selectedNamespaceName: '',
      namespaceList: [],
      namespaceLoading: false,
      columnList: ['name', 'cluster_id', 'namespace', 'Age'],
      renderList: [],
      renderListTmp: [],
      curPageData: [],
      isShowSideslider: false,
      curShowName: '',
      curShowNamespace: '',
      searchKey: '',
      isShowUpdateDialog: false,
      updateItem: null,
      isShowScale: false,
      scaleItem: null,
      isCheckAll: false,
      checkedNodeList: [],

      batchDelDialogConf: {
        title: '',
        isShow: false,
        list: [],
        isDeleting: false,
      },
    };
  },
  computed: {
    projectId() {
      return this.$route.params.projectId;
    },
    isEn() {
      return this.$store.state.isEn;
    },
    curClusterId() {
      return this.$store.getters.curClusterId;
    },
    clusterList() {
      return this.$store.state.cluster.clusterList;
    },
  },
  watch: {
    curClusterId: {
      handler(v) {
        this.selectedClusterId = v;
        this.curPageData = [];
      },
      immediate: true,
    },
  },
  async mounted() {
    await this.getClusters();
    await this.fetchData({
      projId: this.projectId,
      limit: this.pageConf.pageSize,
      offset: 0,
    });
  },
  destroyed() {
    this.bkMessageInstance?.close();
  },
  methods: {
    /**
             * 获取所有的集群
             */
    async getClusters() {
      try {
        if (this.clusterList.length) {
          if (!this.curClusterId) {
            this.selectedClusterId = this.clusterList[0].cluster_id;
          }
          this.getNameSpaceList();
        }
      } catch (e) {
        catchErrorHandler(e, this);
      }
    },

    /**
             * 获取命名空间列表
             */
    async getNameSpaceList() {
      try {
        const { getNamespaceData } = useNamespace();

        this.namespaceLoading = true;
        const res = await getNamespaceData({
          $clusterId: this.selectedClusterId,
        });
        const list = res || [];
        this.namespaceList.splice(0, this.namespaceList.length, ...list);
      } catch (e) {
        catchErrorHandler(e, this);
      } finally {
        this.namespaceLoading = false;
      }
    },

    /**
             * 集群下拉框 item-selected 事件
             *
             * @param {string} clusterId 集群 id
             * @param {Object} data 集群对象
             */
    async handleChangeCluster(clusterId) {
      this.selectedNamespaceName = '';
      this.selectedClusterId = clusterId;
      await this.getNameSpaceList();
    },

    /**
             * 命名空间下拉框 item-selected 事件
             *
             * @param {string} selectedNamespaceName 命名空间 name
             * @param {Object} data 命名空间对象
             */
    handleChangeNamespace(selectedNamespaceName) {
      this.selectedNamespaceName = selectedNamespaceName;
    },

    /**
             * 命名空间下拉框 clear 事件
             */
    handleClearNamespace() {
      this.selectedNamespaceName = '';
    },

    /**
             * 获取表格数据
             */
    async fetchData() {
      this.isPageLoading = true;
      try {
        if (!this.selectedClusterId) return;

        const params = {};
        if (this.selectedNamespaceName) {
          params.namespace = this.selectedNamespaceName;
        }
        const res = await this.$store.dispatch('app/getGameStatefulsetList', {
          projectId: this.projectId,
          clusterId: this.selectedClusterId,
          gamestatefulsets: 'gamedeployments.tkex.tencent.com',
          data: params,
        });

        const data = res.data || { td_list: [], th_list: [] };


        if (data.th_list.length) {
          this.columnList.splice(0, this.columnList.length, ...data.th_list);
        } else {
          this.columnList.splice(0, this.columnList.length, ...['name', 'cluster_id', 'namespace', 'Age']);
        }
        data.td_list.forEach((item) => {
          item.isChecked = false;
        });
        this.renderListTmp.splice(0, this.renderListTmp.length, ...data.td_list);

        if (this.searchKey.trim()) {
          this.renderList.splice(
            0,
            this.renderList.length,
            ...this.renderListTmp.filter(item => item.name.indexOf(this.searchKey) > -1),
          );
        } else {
          this.renderList.splice(0, this.renderList.length, ...this.renderListTmp);
        }

        this.pageConf.curPage = 1;
        this.initPageConf();
        this.curPageData = this.getDataByPage(this.pageConf.curPage);
      } catch (e) {
      } finally {
        this.isPageLoading = false;
        // 晚消失是为了防止整个页面loading和表格数据loading效果叠加产生闪动
        setTimeout(() => {
          this.isInitLoading = false;
        }, 200);
      }
    },

    /**
             * 初始化分页配置
             */
    initPageConf() {
      const total = this.renderList.length;
      this.pageConf.total = total;
      this.pageConf.totalPage = Math.ceil(total / this.pageConf.pageSize);
    },

    /**
             * 分页大小更改
             *
             * @param {number} pageSize pageSize
             */
    handlePageSizeChange(pageSize) {
      this.pageConf.pageSize = pageSize;
      this.pageConf.curPage = 1;
      this.initPageConf();
      this.handlePageChange();
    },
    handleSelectionChange(selection) {
      this.checkedNodeList = selection;
    },

    /**
             * 获取分页数据
             * @param  {number} page 第几页
             * @return {object} data 数据
             */
    getDataByPage(page) {
      let startIndex = (page - 1) * this.pageConf.pageSize;
      let endIndex = page * this.pageConf.pageSize;
      if (startIndex < 0) {
        startIndex = 0;
      }
      if (endIndex > this.renderList.length) {
        endIndex = this.renderList.length;
      }
      return this.renderList.slice(startIndex, endIndex);
    },

    /**
             * 页数改变回调
             * @param  {number} page 第几页
             */
    handlePageChange(page = 1) {
      this.pageConf.curPage = page;
      const data = this.getDataByPage(page);
      this.curPageData.splice(0, this.curPageData.length, ...data);

      // 当前页选中的
      const selectedNodeList = this.curPageData.filter(item => item.isChecked === true);
      this.isCheckAll = selectedNodeList.length === this.curPageData.length;
    },

    /**
             * 搜索框清除事件
             */
    clearSearch() {
      this.searchKey = '';
      this.handleClick();
    },

    /**
             * 搜索按钮点击
             *
             * @param {Object} e 时间对象
             */
    async handleClick() {
      await this.fetchData();
    },

    /**
             * 显示更新弹框
             *
             * @param {Object} item 当前行对象
             * @param {number} index 当前行对象索引
             *
             * @return {string} returnDesc
             */
    update(item) {
      this.isShowUpdateDialog = true;
      this.updateItem = item;
    },

    /**
             * 关闭更新弹框
             */
    hideGamedeploymentsUpdate() {
      this.isShowUpdateDialog = false;
      setTimeout(() => {
        this.updateItem = null;
      }, 300);
    },

    /**
             * 更新 gamestatefulset 成功回调
             */
    async gamedeploymentsUpdateSuccess() {
      this.hideGamedeploymentsUpdate();
      await this.fetchData();
    },

    /**
             * 删除当前行
             *
             * @param {Object} item 当前行对象
             * @param {number} index 当前行对象索引
             *
             * @return {string} returnDesc
             */
    async del(item) {
      // eslint-disable-next-line @typescript-eslint/no-this-alias
      const me = this;
      const boxStyle = {
        'margin-top': '-20px',
        'margin-bottom': '-20px',
      };
      // const titleStyle = {
      //     style: {
      //         'text-align': 'left',
      //         'font-size': '20px',
      //         'margin-bottom': '15px',
      //         'color': '#313238'
      //     }
      // }
      const itemStyle = {
        style: {
          'text-align': 'left',
          'font-size': '14px',
          'margin-bottom': '3px',
          color: '#71747c',
        },
      };

      const contexts = [
        // me.$createElement('h5', titleStyle, me.$t('deploy.templateset.delete')),
        me.$createElement('p', itemStyle, `${me.$t('generic.label.name')}：${item.name}`),
        me.$createElement('p', itemStyle, `${me.$t('generic.label.cluster1')}：${item.cluster_id}`),
      ];
      if (item.namespace) {
        contexts.push(me.$createElement('p', itemStyle, `${me.$t('k8s.namespace')}：${item.namespace}`));
      }
      me.$bkInfo({
        title: me.$t('generic.title.confirmDelete'),
        clsName: 'biz-remove-dialog',
        confirmLoading: true,
        content: me.$createElement('div', { class: 'biz-confirm-desc', style: boxStyle }, contexts),
        async confirmFn() {
          try {
            await me.$store.dispatch('app/deleteGameStatefulsetInfo', {
              projectId: me.projectId,
              clusterId: me.selectedClusterId,
              gamestatefulsets: 'gamedeployments.tkex.tencent.com',
              name: item.name,
              data: {
                namespace: item.namespace,
              },
            });

            me.bkMessageInstance?.close();
            me.bkMessageInstance = me.$bkMessage({
              theme: 'success',
              message: me.$t('generic.msg.success.delete'),
              delay: 1000,
            });
            await me.fetchData();
          } catch (e) {
            console.error(e);
          }
        },
      });
    },

    async showDetail(name, namespace) {
      this.$router.push({
        name: 'gamedeploymentsInstanceDetail',
        params: {
          category: 'custom_objects',
          name,
          namespace,
        },
        query: {
          kind: 'GameDeployment',
          crd: 'gamedeployments.tkex.tencent.com',
          name,
          namespace,
          cluster_id: this.selectedClusterId,
        },
      });
    },

    /**
             * 隐藏 sideslider
             */
    // hideSideslider() {
    //   this.curShowName = '';
    //   this.curShowNamespace = '';
    //   this.isShowSideslider = false;
    // },

    /**
             * 显示扩缩容弹框
             *
             * @param {Object} item 当前行对象
             * @param {number} index 当前行对象索引
             *
             * @return {string} returnDesc
             */
    scale(item) {
      this.isShowScale = true;
      this.scaleItem = item;
    },

    /**
             * 隐藏扩缩容弹框
             */
    hideScale() {
      this.isShowScale = false;
      setTimeout(() => {
        this.scaleItem = null;
      }, 300);
    },

    /**
             * gamestatefulset 扩缩容成功回调
             */
    async gamedeploymentsScaleSuccess() {
      this.hideScale();
      await this.fetchData();
    },

    /**
             * 列表全选
             */
    checkAllItem(e) {
      const isChecked = e.target.checked;
      this.curPageData.forEach((item) => {
        item.isChecked = isChecked;
      });

      const checkedNodeList = [];
      checkedNodeList.splice(0, 0, ...this.checkedNodeList);
      // 用于区分是否已经选择过
      const hasCheckedList = checkedNodeList.map(item => item.name + item.namespace + item.cluster_id);
      if (isChecked) {
        const checkedList = this.curPageData
          .filter(item => !hasCheckedList.includes(item.name + item.namespace + item.cluster_id));
        checkedNodeList.push(...checkedList);
        this.checkedNodeList.splice(0, this.checkedNodeList.length, ...checkedNodeList);
      } else {
        // 当前页所有合法的 node id 集合
        const validIdList = this.curPageData.map(item => item.name + item.namespace + item.cluster_id);

        const newCheckedNodeList = [];
        this.checkedNodeList.forEach((checkedNode) => {
          if (validIdList.indexOf(checkedNode.name + checkedNode.namespace + checkedNode.cluster_id) < 0) {
            newCheckedNodeList.push(JSON.parse(JSON.stringify(checkedNode)));
          }
        });
        this.checkedNodeList.splice(0, this.checkedNodeList.length, ...newCheckedNodeList);
      }
    },

    /**
             * 列表每一行的 checkbox 点击
             *
             * @param {Object} row 当前策略对象
             */
    checkItem(row) {
      this.$nextTick(() => {
        // 当前页选中的
        const selectedNodeList = this.curPageData.filter(item => item.isChecked === true);
        this.isCheckAll = selectedNodeList.length === this.curPageData.length;

        const checkedNodeList = [];
        if (row.isChecked) {
          checkedNodeList.splice(0, checkedNodeList.length, ...this.checkedNodeList);
          if (!this.checkedNodeList
            .filter(checkedNode => checkedNode.name + checkedNode.namespace + checkedNode.cluster_id
             === row.name + row.namespace + row.cluster_id).length) {
            checkedNodeList.push(row);
          }
        } else {
          this.checkedNodeList.forEach((checkedNode) => {
            if (checkedNode.name + checkedNode.namespace + checkedNode.cluster_id
            !== row.name + row.namespace + row.cluster_id) {
              checkedNodeList.push(JSON.parse(JSON.stringify(checkedNode)));
            }
          });
        }
        this.checkedNodeList.splice(0, this.checkedNodeList.length, ...checkedNodeList);
      });
    },

    /**
             * 批量删除
             */
    batchDel() {
      if (!this.checkedNodeList.length) {
        this.bkMessageInstance?.close();
        this.bkMessageInstance = this.$bkMessage({
          theme: 'error',
          message: this.$t('deploy.templateset.unselectedGameDeployments'),
        });
        return;
      }

      const len = this.checkedNodeList.length;
      const repeat = {};
      for (let i = 0; i < len; i++) {
        const item = this.checkedNodeList[i];
        repeat[item.namespace] = 1;
      }
      if (Object.keys(repeat).length > 1) {
        this.bkMessageInstance?.close();
        this.bkMessageInstance = this.$bkMessage({
          theme: 'error',
          message: this.$t('deploy.templateset.oneNS'),
          delay: 1000,
        });
        return;
      }

      this.batchDelDialogConf.isShow = true;
      this.batchDelDialogConf.title = this.isEn
        ? `Confirm to delete GameDeployments under namespace [${this.checkedNodeList[0].namespace}]`
        : `确定删除命名空间【${this.checkedNodeList[0].namespace}】下的GameDeployments？`;
      this.batchDelDialogConf.list.splice(0, this.batchDelDialogConf.list.length, ...this.checkedNodeList);
    },

    /**
             * 批量删除弹框取消按钮
             */
    hideBatchDelDialog() {
      this.batchDelDialogConf.isShow = false;
      setTimeout(() => {
        this.batchDelDialogConf.list.splice(0, this.batchDelDialogConf.list.length, ...[]);
        this.batchDelDialogConf.isDeleting = false;
        this.batchDelDialogConf.title = '';
      }, 300);
    },

    /**
             * 批量删除弹框确定按钮
             */
    async batchDelConfirm() {
      try {
        this.batchDelDialogConf.isDeleting = true;

        const cobjNameList = this.checkedNodeList.map(item => item.name);

        await this.$store.dispatch('app/batchDeleteGameStatefulset', {
          projectId: this.projectId,
          clusterId: this.selectedClusterId,
          gamestatefulsets: 'gamedeployments.tkex.tencent.com',
          data: {
            namespace: this.checkedNodeList[0].namespace,
            cobj_name_list: cobjNameList,
          },
        });

        this.hideBatchDelDialog();
        this.bkMessageInstance?.close();
        this.bkMessageInstance = this.$bkMessage({
          theme: 'success',
          message: this.$t('generic.msg.success.delete'),
          delay: 1000,
        });
        this.checkedNodeList.splice(0, this.checkedNodeList.length, ...[]);
        this.isCheckAll = false;
        await this.fetchData();
      } catch (e) {
        console.error(e);
      } finally {
        this.batchDelDialogConf.isDeleting = false;
      }
    },
    handleClearSearchData() {
      this.searchKey = '';
      this.handleClick();
    },
  },
};
</script>

<style scoped>
    @import '../gamestatefulset.css';
</style>
