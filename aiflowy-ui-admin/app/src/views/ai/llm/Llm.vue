<script setup>

import CategoryPanel from "#/components/categoryPanel/CategoryPanel.vue";
import {onMounted, ref} from "vue";
import {getLlmBrandList} from "#/api/ai/llm.js";

const brandListData = ref([])
onMounted(() =>{
  getLlmBrandList().then(res => {
    console.log('res')
    console.log(res)
    brandListData.value = res.data
  })
})

const handleCategoryClick = (category) => {
  console.log('category')
  console.log(category.key)
}
</script>

<template>
  <CategoryPanel :categories="brandListData" title-key="title" :use-img-for-svg="true" :expandWidth="150" @click="handleCategoryClick"/>
</template>

<style scoped>

</style>
