<template>
  <div class="relative">
    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :class="[
        'w-full px-4 py-2.5 text-sm border rounded-lg transition-all duration-200',
        'focus:outline-none focus:ring-2 focus:ring-offset-0',
        disabled 
          ? 'bg-gray-50 text-gray-500 cursor-not-allowed border-gray-200' 
          : 'bg-white text-gray-900 border-gray-300 focus:border-primary-500 focus:ring-primary-500',
        error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : '',
        'placeholder:text-gray-400'
      ]"
      @input="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      @focus="$emit('focus')"
    />
    <div v-if="error" class="absolute -bottom-5 left-0 text-xs text-red-600 mt-1">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue', 'blur', 'focus'])
</script>

