<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import {
  CaretDownSmallIcon,
  CaretRightSmallIcon,
} from 'tdesign-icons-vue-next';
import { Tree } from 'tdesign-vue-next';

import { getDeptTree } from '#/api/system/dept';

const emit = defineEmits(['select']);

const treeData = ref<any[]>([]);

const treeKeys = {
  value: 'id',
  label: 'name',
  children: 'children',
};

// 展开的节点
const expanded = ref<Array<number | string>>([-1]);

async function fetchDeptTree() {
  try {
    const res = await getDeptTree();
    treeData.value = [
      {
        id: -1,
        name: '所有部门',
        children: res || [],
      },
    ];
  } catch (error) {
    console.error('Failed to fetch dept tree', error);
  }
}

function handleNodeClick(context: any) {
  const { node } = context;
  emit('select', node.value === -1 ? '' : node.value);
}

onMounted(() => {
  fetchDeptTree();
});
</script>

<template>
  <div class="h-full w-56 flex flex-col bg-white">
    <div class="p-3 flex items-center justify-between border-b border-gray-100">
      <div class="text-[14px] text-gray-700 font-medium">搜索部门</div>
      <div class="text-[14px] text-blue-500 cursor-pointer hover:text-blue-600">
        折叠
      </div>
    </div>
    <div class="flex-1 overflow-auto p-2 custom-tree-wrap">
      <Tree
        :data="treeData"
        :keys="treeKeys"
        hover
        :expanded="expanded"
        line
        @click="handleNodeClick"
      >
        <template #icon="{ node }">
          <CaretDownSmallIcon v-if="node.expanded" />
          <CaretRightSmallIcon
            v-else-if="node.children && node.children.length > 0"
          />
          <span v-else class="w-4 h-4 inline-block"></span>
        </template>
      </Tree>
    </div>
  </div>
</template>

<style scoped>
.custom-tree-wrap :deep(.t-tree) {
  --td-tree-hover-color: #f3f4f6;
  --td-tree-active-color: #f3f4f6;

  color: #333;
}

.custom-tree-wrap :deep(.t-tree__item) {
  padding: 4px 8px;
  border-radius: 4px;
}

.custom-tree-wrap :deep(.t-tree__item--active) {
  font-weight: 500;
  color: #165dff;
}
</style>
