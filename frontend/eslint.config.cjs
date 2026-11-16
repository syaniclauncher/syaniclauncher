const pluginVue = require('eslint-plugin-vue')

module.exports = [
  { ignores: ['dist/**'] },
  ...pluginVue.configs['flat/essential'],
  {
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
]
