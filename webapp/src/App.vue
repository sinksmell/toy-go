<template>
  <div id="app" >
    <n-message-provider>
      <Naive/>
    </n-message-provider>
  </div>
</template>

<script>
import Naive from "./components/Naive.vue";

export default {
  name: 'App',
  components: {
    Naive,
  },
  data() {
    return {}
  }
};
</script>

