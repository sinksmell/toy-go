<template>
  <n-space vertical class="header">
    <n-tooltip placement="bottom" trigger="hover">
      <template #trigger>
        <div class="logo" @click="jsonToGo">
          <img :src="logoSVG" height="100" alt="logo"/>
        </div>
      </template>
      <span> Click Me to Generate Go Struct!</span>
    </n-tooltip>
  </n-space>

  <n-config-provider :hljs="hljs" :locale="zhCN" :date-locale="dateZhCN">
    <n-grid x-gap="24" :cols="2" class="main">
      <n-gi>
        <n-space vertical>
          <n-input
              v-model:value="jsonStr"
              @change="jsonStrChange"
              type="textarea"
              placeholder="paste json to here，click the logo to generate..."
              :autosize="{
        minRows:24,
      }"
          />
        </n-space>
      </n-gi>

      <n-gi>

        <div v-show="goStruct === ''">
          <n-space vertical align="center">
            <n-result status="404" title="Let`s Go!"
                      size="large"
                      description="input json then click the gopher icon, you will get go code!">
            </n-result>
            <template #icon>
              <n-icon><LogoGithub /></n-icon>
            </template>
          </n-space>
        </div>

        <n-card v-if="goStruct !== ''">
          <div style="overflow: auto;">
            <n-space vertical>
              <n-code
                  v-model:code="goStruct"
                  language="go"
                  :autosize="{minRows:24,}"
              />
            </n-space>
          </div>
        </n-card>
      </n-gi>
    </n-grid>
  </n-config-provider>

</template>


<script>
import {defineComponent,} from 'vue'
import {LogoGithub} from '@vicons/ionicons4'
// theme
// locale & dateLocale
import {
  createTheme,
  datePickerDark,
  dateZhCN,
  inputDark,
  NConfigProvider,
  NInput,
  NSpace,
  useMessage,
  zhCN,
} from 'naive-ui'
import hljs from 'highlight.js/lib/core'
import cpp from 'highlight.js/lib/languages/cpp'
import go from 'highlight.js/lib/languages/go'
import js from 'highlight.js/lib/languages/javascript'
import axios from "axios";

hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('go', go)
hljs.registerLanguage('js', js)

export default defineComponent({
  name: "naive",
  data() {
    return {
      logoSVG: require('../assets/logo.svg'),
      jsonStr: '',
      goStruct: '',
    }
  },
  components: {
    NConfigProvider,
    NInput,
    // NDatePicker,
    NSpace,
    LogoGithub
  },
  methods: {
    jsonStrChange(v) {
      if (v === '') {
        this.goStruct = ''
      }
    },
  }
  ,
  setup() {
    const msg = useMessage()

    return {
      darkTheme: createTheme([inputDark, datePickerDark]),
      zhCN,
      dateZhCN,
      hljs,

      jsonToGo() {
        console.log()
        // console.log(this.data.jsonStr)
        console.log(this.jsonStr)
        console.log(this.jsonStr)

        axios
            .post(`${process.env.VUE_APP_API_URL}/api/json/to/go`, {
              "json_str": this.jsonStr,
            })

            .then(response => {
              console.log(response.status)
              console.log(response.data.message)
              if (response.status !== 200) {
                msg.error(response.data.message)
                return
              }

              msg.success("generate success")
              this.goStruct = response.data.go_struct_str
            }).catch(err => {

          if (err.response) {
            msg.error(err.response.data.message)
          } else if (err.request) {
            msg.error(err.request.json_str)
          } else {
            console.log(err)
          }
        });
      },
    }

  },

})

</script>

<style>

@media screen and (min-width: 800px) {
  body {
    padding-right: 15%;
    padding-left: 15%;
  }
}

@media screen and (min-width: 1600px) {
  body {
    padding-right: 30%;
    padding-left: 30%;
  }
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, "Courier New",
  monospace;
}

.title {
  text-align: center;
}

.header {
    margin-top: 5%;
    padding-right: 5%;
    padding-left: 5%;
}

.main {
  margin-top: 2%;
  padding-right: 5%;
  padding-left: 5%;
}

.logo {
  text-align: center;
}
</style>

