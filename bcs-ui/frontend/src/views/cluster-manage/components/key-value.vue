<template>
  <div>
    <span
      class="add-btn" v-if="!labels.length"
      @click="handleAddLabel(0)">
      <i class="bk-icon icon-plus-circle-shape mr5"></i>
      {{$t('generic.button.add')}}
    </span>
    <div class="key-value" v-for="(item, index) in labels" :key="index">
      <Validate class="w-[100%]" :value="item.key" :rules="keyRules">
        <bcs-input
          v-model="item.key"
          :placeholder="$t('generic.label.key')"
          @change="handleLabelKeyChange">
        </bcs-input>
      </Validate>
      <span class="ml8 mr8">=</span>
      <Validate class="w-[100%]" :value="item.value" :rules="valueRules">
        <bcs-input
          v-model="item.value"
          :placeholder="$t('generic.label.value')"
          @change="handleLabelValueChange">
        </bcs-input>
      </Validate>
      <i class="bk-icon icon-plus-circle ml15" @click="handleAddLabel(index)"></i>
      <i
        :class="['bk-icon icon-minus-circle ml10', { disabled: disabledDelete }]"
        @click="handleDeleteLabel(index)">
      </i>
    </div>
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs, watch } from 'vue';

import Validate from '@/components/validate.vue';
import $i18n from '@/i18n/i18n-setup';

interface ILabel {
  key: string;
  value: string;
}
export default defineComponent({
  name: 'KeyValue',
  components: { Validate },
  model: {
    prop: 'value',
    event: 'change',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    minItem: {
      type: Number,
      default: 1,
    },
    disableDeleteItem: {
      type: Boolean,
      default: true,
    },
    keyRules: {
      type: Array as PropType<Array<any>>,
      default: () => [
        {
          message: $i18n.t('generic.validate.labelKey'),
          validator: '^[A-Za-z0-9._/-]+$',
        },
      ],
    },
    valueRules: {
      type: Array as PropType<Array<any>>,
      default: () => [
        {
          message: $i18n.t('generic.validate.labelKey'),
          validator: '^[A-Za-z0-9._/-]+$',
        },
      ],
    },
  },
  setup(props, ctx) {
    const labels = ref<ILabel[]>([]);
    const { value, disableDeleteItem } = toRefs(props);
    const unWatchValue = watch(value, () => {
      labels.value = Object.keys(props.value).map(key => ({
        key,
        value: props.value[key],
      }));
      if (labels.value.length < (props.minItem || 0)) {
        const reset: ILabel[] = new Array(props.minItem - labels.value.length).fill({ key: '', value: '' });
        labels.value.push(...reset);
      }
    }, { immediate: true, deep: true });
    watch(labels, () => {
      unWatchValue();
      emitChange();
    }, { deep: true });
    const disabledDelete = computed(() => labels.value.length <= props.minItem && disableDeleteItem.value);
    const emitChange = () => {
      const keyValues = labels.value.reduce((pre, item) => {
        if (!item.key) return pre;

        pre[item.key] = item.value;
        return pre;
      }, {});
      ctx.emit('change', keyValues);
    };
    const handleLabelKeyChange = (newValue, oldValue) => {
      ctx.emit('key-change', newValue, oldValue);
    };
    const handleLabelValueChange = (newValue, oldValue) => {
      ctx.emit('value-change', newValue, oldValue);
    };
    const handleAddLabel = (index = 0) => {
      labels.value.splice(index + 1, 0, { key: '', value: '' });
    };
    const handleDeleteLabel = (index) => {
      if (disabledDelete.value) return;
      labels.value.splice(index, 1);
    };
    return {
      labels,
      disabledDelete,
      handleLabelKeyChange,
      handleLabelValueChange,
      handleAddLabel,
      handleDeleteLabel,
    };
  },
});
</script>
<style lang="postcss" scoped>
.add-btn {
    font-size: 14px;
    color: #3a84ff;
    cursor: pointer;
    display: flex;
    align-items: center;
    height: 32px;
    max-width: 100px;
}
.key-value {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    .bk-icon {
        font-size: 20px;
        color: #979bA5;
        cursor: pointer;
        &.disabled {
            color: #DCDEE5;
            cursor: not-allowed;
        }
    }
    .ml8 {
        margin-left: 8px;
    }
    .mr8 {
        margin-right: 8px;
    }
}
</style>
